package rustack_terraform

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pilat/rustack-go/rustack"
)

func resourceRustackVm() *schema.Resource {
	args := Defaults()
	args.injectCreateVm()
	args.injectContextVdcById()
	args.injectContextTemplateById() // override template_id

	return &schema.Resource{
		CreateContext: resourceRustackVmCreate,
		ReadContext:   resourceRustackVmRead,
		UpdateContext: resourceRustackVmUpdate,
		DeleteContext: resourceRustackVmDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: args,
	}
}

func resourceRustackVmCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting VDC: %s", err)
	}

	template, err := GetTemplateById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting template: %s", err)
	}

	vmName := d.Get("name").(string)
	cpu := d.Get("cpu").(int)
	ram := d.Get("ram").(int)
	userData := d.Get("user_data").(string)
	log.Printf("Vm details: name=%s, cpu: %d, ram: %d, user_data: %s, template name: %s",
		vmName, cpu, ram, userData, template.Name)

	// System disk creation
	systemDiskArgs := d.Get("system_disk").(string)
	var diskSize int
	var diskStorageProfile string
	_, err = fmt.Sscanf(systemDiskArgs, "%d-%s", &diskSize, &diskStorageProfile)
	if err != nil {
		return diag.Errorf("ERROR. Wrong system disk format: %s", err)
	}
	profiles, err := targetVdc.GetStorageProfiles()
	if err != nil {
		return diag.FromErr(err)
	}
	var storageProfile *rustack.StorageProfile
	for _, profile := range profiles {
		if strings.ToLower(profile.Name) == strings.ToLower(diskStorageProfile) {
			storageProfile = profile
		}
	}
	if storageProfile == nil {
		return diag.Errorf("ERROR. storage profile %s not found", diskStorageProfile)
	}

	var systemDisk []*rustack.Disk
	newDisk := rustack.NewDisk("System", diskSize, storageProfile)
	systemDisk = append(systemDisk, &newDisk)

	// Join additional disk
	disksIds := d.Get("disks").(*schema.Set).List()
	var disks []*rustack.Disk
	for _, diskId := range disksIds {
		disk, err := manager.GetDisk(diskId.(string))
		if err != nil {
			return diag.Errorf("ERROR. Disk not found: %s", err)
		}
		disks = append(disks, disk)
	}

	portsCount := d.Get("port.#").(int)
	ports := make([]*rustack.Port, portsCount)
	for i := 0; i < portsCount; i++ {
		portPrefix := fmt.Sprintf("port.%d", i)

		newPort, err := createPort(d, manager, &portPrefix)
		if err != nil {
			return diag.FromErr(err)
		}

		ports[i] = newPort

		log.Printf("Create port with network '%s'", newPort.Network.ID)
	}

	var floatingIp *string = nil
	if d.Get("floating").(bool) {
		floatingIpStr := "RANDOM_FIP"
		floatingIp = &floatingIpStr
	}

	newVm := rustack.NewVm(vmName, cpu, ram, template, nil, &userData, ports, systemDisk, floatingIp)

	err = targetVdc.CreateVm(&newVm)
	if err != nil {
		return diag.Errorf("Error creating vm: %s", err)
	}
	d.Set("system_disk_id", newVm.Disks[0].ID)

	syncDisks(d, manager, targetVdc, &newVm)

	d.SetId(newVm.ID)
	log.Printf("[INFO] VM created, ID: %s", d.Id())

	return resourceRustackVmRead(ctx, d, meta)
}

func resourceRustackVmRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	vm, err := manager.GetVm(d.Id())
	if err != nil {
		return diag.Errorf("Error getting vm: %s", err)
	}

	d.SetId(vm.ID)
	d.Set("name", vm.Name)
	d.Set("cpu", vm.Cpu)
	d.Set("ram", vm.Ram)
	d.Set("template_id", vm.Template.ID)

	// d.Set("user_data", vm.UserData)

	flattenDisks := make([]map[string]interface{}, len(vm.Disks))
	for i, disk := range vm.Disks {
		flattenDisks[i] = map[string]interface{}{
			"id":                 disk.ID,
			"name":               disk.Name,
			"size":               disk.Size,
			"storage_profile_id": disk.StorageProfile.ID,
		}
	}
	d.Set("disk", flattenDisks)

	flattenPorts := make([]map[string]interface{}, len(vm.Ports))
	for i, port := range vm.Ports {
		flattenFirewallTemplates := make([]string, len(port.FirewallTemplates))
		for j, firewallTemplate := range port.FirewallTemplates {
			flattenFirewallTemplates[j] = firewallTemplate.ID
		}
		sort.Strings(flattenFirewallTemplates)

		flattenPorts[i] = map[string]interface{}{
			"id":                 port.ID,
			"network_id":         port.Network.ID,
			"firewall_templates": flattenFirewallTemplates,
			"ip_address":         port.IpAddress,
		}
	}
	d.Set("port", flattenPorts)
	d.Set("floating", vm.Floating != nil)
	d.Set("floating_ip", "")
	if vm.Floating != nil {
		d.Set("floating_ip", vm.Floating.IpAddress)
	}

	return nil
}

func resourceRustackVmUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting VDC: %s", err)
	}

	hasFlavorChanged := false
	needUpdate := false

	vm, err := manager.GetVm(d.Id())
	if err != nil {
		return diag.Errorf("Error getting vm: %s", err)
	}

	// Detect vm changes
	if d.HasChange("name") {
		needUpdate = true
		vm.Name = d.Get("name").(string)
	}

	if d.HasChange("cpu") || d.HasChange("ram") {
		needUpdate = true
		hasFlavorChanged = true
		vm.Cpu = d.Get("cpu").(int)
		vm.Ram = d.Get("ram").(int)
	}

	needPowerOn := false
	if hasFlavorChanged && !vm.HotAdd && vm.Power {
		vm.PowerOff()
		needPowerOn = true
	}

	if d.HasChange("floating") {
		needUpdate = true
		if !d.Get("floating").(bool) {
			vm.Floating = nil
		} else {
			floatingIpStr := "RANDOM_FIP"
			vm.Floating = &rustack.Port{IpAddress: &floatingIpStr}
		}
		d.Set("floating", vm.Floating != nil)
	}

	if needUpdate {
		if err := vm.Update(); err != nil {
			return diag.Errorf("Error getting vm: %s", err)
		}
	}

	if needPowerOn {
		vm.PowerOn()
	}

	if diags := syncDisks(d, manager, targetVdc, vm); diags.HasError() {
		return diags
	}

	if diags := syncPorts(d, manager, targetVdc, vm); diags.HasError() {
		return diags
	}

	return resourceRustackVmRead(ctx, d, meta)
}

func resourceRustackVmDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	vm, err := manager.GetVm(d.Id())
	if err != nil {
		return diag.Errorf("Error getting vm: %s", err)
	}

	err = vm.Delete()
	if err != nil {
		return diag.Errorf("Error deleting vm: %s", err)
	}

	return nil
}

func createDisk(d *schema.ResourceData, manager *rustack.Manager, diskPrefix *string) (*rustack.Disk, error) {
	vdc, err := GetVdcById(d, manager)
	if err != nil {
		return nil, err
	}

	diskName := d.Get(MakePrefix(diskPrefix, "name")).(string)
	diskSize := d.Get(MakePrefix(diskPrefix, "size")).(int)
	storageProfile, err := GetStorageProfileById(d, manager, vdc, diskPrefix)
	if err != nil {
		return nil, err
	}

	newDisk := rustack.NewDisk(diskName, diskSize, storageProfile)
	return &newDisk, nil
}

func createPort(d *schema.ResourceData, manager *rustack.Manager, portPrefix *string) (*rustack.Port, error) {
	vdc, err := GetVdcById(d, manager)
	if err != nil {
		return nil, err
	}

	portNetwork, err := GetNetworkById(d, manager, portPrefix)
	if err != nil {
		return nil, err
	}

	firewallsCount := d.Get(MakePrefix(portPrefix, "firewall_templates.#")).(int)
	firewalls := make([]*rustack.FirewallTemplate, firewallsCount)
	for j := 0; j < firewallsCount; j++ {
		firewallName := MakePrefix(portPrefix, fmt.Sprintf("firewall_templates.%d", j))
		portFirewall, err := GetFirewallTemplateById(d, manager, vdc, &firewallName)
		if err != nil {
			return nil, err
		}

		firewalls[j] = portFirewall
	}

	newPort := rustack.NewPort(portNetwork, firewalls, nil)
	return &newPort, nil
}

func syncDisks(d *schema.ResourceData, manager *rustack.Manager, vdc *rustack.Vdc, vm *rustack.Vm) diag.Diagnostics {
	disksIds := d.Get("disks").(*schema.Set).List()

	disksIds = append(disksIds, d.Get("system_disk_id").(string))

	// Which disks are present on vm and not mentioned in the state?
	// Detach disks
	needReload := false
	for _, vmDisk := range vm.Disks {
		found := false
		for _, diskId := range disksIds {
			if diskId.(string) == vmDisk.ID {
				found = true
				break
			}
		}

		if !found {
			log.Printf("Disk %s found on vm and not mentioned in the state."+
				" Disk will be detached", vmDisk.ID)
			vm.WaitLock()
			vm.DetachDisk(vmDisk)
			needReload = true
		}
	}

	// List disks to join
	for _, diskId := range disksIds {
		found := false
		for _, vmDisk := range vm.Disks {
			if diskId.(string) == vmDisk.ID {
				found = true
				break
			}
		}

		if !found {
			disk, err := manager.GetDisk(diskId.(string))
			if err != nil {
				return diag.FromErr(err)
			}
			log.Printf("Disk `%s` will be Attached", disk.ID)
			vm.WaitLock()
			if err = vm.AttachDisk(disk); err != nil {
				return diag.Errorf("ERROR. Cannot attach disk `%s`: %s", disk.ID, err)
			}
			needReload = true
		}
	}

	// System disk resize
	if d.HasChange("system_disk") {
		systemDiskArgs := d.Get("system_disk").(string)
		var diskSize int
		var diskStorageProfile string
		_, err := fmt.Sscanf(systemDiskArgs, "%d-%s", &diskSize, &diskStorageProfile)
		if err != nil {
			return diag.Errorf("ERROR. Wrong system disk format: %s", err)
		}
		systemDisk, err := manager.GetDisk(d.Get("system_disk_id").(string))
		if err != nil {
			return diag.Errorf("ERROR. Wrong system disk id: %s", err)
		}
		systemDisk.Resize(diskSize)
	}

	if needReload {
		if err := vm.Reload(); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func syncPorts(d *schema.ResourceData, manager *rustack.Manager, vdc *rustack.Vdc, vm *rustack.Vm) diag.Diagnostics {
	portsCount := d.Get("port.#").(int)

	// Which ports are present on vm and not mentioned in the state?
	needReload := false
	for _, port := range vm.Ports {
		found := false
		for i := 0; i < portsCount; i++ {
			portPrefix := fmt.Sprintf("port.%d", i)
			portId := d.Get(MakePrefix(&portPrefix, "id"))
			if portId == port.ID {
				found = true
				break
			}
		}

		if !found {
			log.Printf("Port %s found on vm and not mentioned in the state. Port will be deleted", port.ID)
			port.Delete()
			needReload = true
		}
	}

	if needReload {
		if err := vm.Reload(); err != nil {
			return diag.FromErr(err)
		}
	}

	// Which ports are present in state and not in vm?
	for i := 0; i < portsCount; i++ {
		portPrefix := fmt.Sprintf("port.%d", i)
		portId := d.Get(MakePrefix(&portPrefix, "id"))

		found := false
		for _, port := range vm.Ports {
			if port.ID == portId {
				found = true
				break
			}
		}

		if !found {
			log.Printf("Port %s found in the state and is not present on vm. Port will be created", portPrefix)

			newPort, err := createPort(d, manager, &portPrefix)
			if err != nil {
				return diag.FromErr(err)
			}

			if err := vm.AddPort(newPort); err != nil {
				return diag.FromErr(err)
			}

			if err := vm.Reload(); err != nil {
				return diag.FromErr(err)
			}
		}

	}

	// Detect port changes for found ports with the same id
	for i := 0; i < portsCount; i++ {
		portPrefix := fmt.Sprintf("port.%d", i)
		portId, portExists := d.GetOk(MakePrefix(&portPrefix, "id"))
		if !portExists {
			// That case has been resolved above
			continue
		}

		var foundPort *rustack.Port = nil
		for _, port := range vm.Ports {
			if port.ID == portId {
				foundPort = port
				break
			}
		}

		if foundPort == nil {
			// That case has been resolved above
			continue
		}

		// Compare foundPort
		pseudoPort, err := createPort(d, manager, &portPrefix)
		if err != nil {
			return diag.FromErr(err)
		}

		// TODO: Compare firewall templates
		isEqual := true
		if len(pseudoPort.FirewallTemplates) != len(foundPort.FirewallTemplates) {
			isEqual = false
		} else {
			for _, l := range pseudoPort.FirewallTemplates {
				found := false
				for _, r := range foundPort.FirewallTemplates {
					if r.ID == l.ID {
						found = true
						break
					}
				}

				if !found {
					isEqual = false
					break
				}
			}
		}

		if !isEqual {
			if err = foundPort.UpdateFirewall(pseudoPort.FirewallTemplates); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return nil
}
