package rustack_terraform

import (
	"context"
	"fmt"
	"log"
	"sort"
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
	systemDiskArgs := d.Get("system_disk.0").(map[string]interface{})
	diskSize := systemDiskArgs["size"].(int)
	storageProfileId := systemDiskArgs["storage_profile_id"].(string)

	storageProfile, err := targetVdc.GetStorageProfile(storageProfileId)
	if err != nil {
		return diag.Errorf("ERROR. storage profile %s not found", storageProfileId)
	}

	systemDiskList := make([]*rustack.Disk, 1)
	newDisk := rustack.NewDisk("Основной диск", diskSize, storageProfile)
	systemDiskList[0] = &newDisk

	ports, diagErr := createPorts(d, manager)
	if diagErr != nil {
		return diagErr
	}

	var floatingIp *string = nil
	if d.Get("floating").(bool) {
		floatingIpStr := "RANDOM_FIP"
		floatingIp = &floatingIpStr
	}

	newVm := rustack.NewVm(vmName, cpu, ram, template, nil, &userData, ports,
		systemDiskList, floatingIp)

	err = targetVdc.CreateVm(&newVm)
	if err != nil {
		return diag.Errorf("Error creating vm: %s", err)
	}

	systemDisk := make([]interface{}, 1)
	systemDisk[0] = map[string]interface{}{
		"id":                 newVm.Disks[0].ID,
		"name":               "Основной диск",
		"size":               newVm.Disks[0].Size,
		"storage_profile_id": newVm.Disks[0].StorageProfile.ID,
	}

	syncDisks(d, manager, targetVdc, &newVm)

	d.Set("system_disk", systemDisk)
	d.SetId(newVm.ID)

	log.Printf("[INFO] VM created, ID: %s", d.Id())

	return resourceRustackVmRead(ctx, d, meta)
}

func createPorts(d *schema.ResourceData, manager *rustack.Manager) ([]*rustack.Port, diag.Diagnostics) {

	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return nil, diag.Errorf("Error getting VDC: %s", err)
	}
	portCount := d.Get("port.#").(int)
	ports := make([]*rustack.Port, portCount)

	for i := 0; i < portCount; i++ {
		portPrefix := fmt.Sprint("port.", i)

		newPort, err := createPort(d, manager, &portPrefix)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		if newPort.Network.Vdc.Id != targetVdc.ID {
			return nil, diag.Errorf("ERROR: Network %s not in target's vdc", newPort.Network.ID)
		}

		ports[i] = newPort

		log.Printf("Create port with network '%s'", newPort.Network.ID)
	}

	return ports, nil
}

func resourceRustackVmRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diagErr diag.Diagnostics) {
	manager := meta.(*CombinedConfig).rustackManager()
	vm, err := manager.GetVm(d.Id())
	if err != nil {
		diagErr = diag.Errorf("Error getting vm: %s", err)
		return
	}

	d.SetId(vm.ID)
	d.Set("name", vm.Name)
	d.Set("cpu", vm.Cpu)
	d.Set("ram", vm.Ram)
	d.Set("template_id", vm.Template.ID)

	flattenDisks := make([]string, len(vm.Disks)-1)
	for i, disk := range vm.Disks {
		if i == 0 {
			systemDisk := make([]interface{}, 1)
			systemDisk[0] = map[string]interface{}{
				"id":                 disk.ID,
				"name":               "Основной диск",
				"size":               disk.Size,
				"storage_profile_id": disk.StorageProfile.ID,
			}

			d.Set("system_disk", systemDisk)
			continue
		}
		flattenDisks[i-1] = disk.ID
	}
	d.Set("disks", flattenDisks)

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

	return
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
			vm.Floating = &rustack.Port{IpAddress: nil}
		} else {
			vm.Floating = &rustack.Port{ID: "RANDOM_FIP"}
		}
		d.Set("floating", vm.Floating != nil)
	}

	if needUpdate {
		if err := repeatOnError(vm.Update, vm); err != nil {
			return diag.Errorf("Error updating vm: %s", err)
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

	disksIds := d.Get("disks").(*schema.Set).List()
	for _, diskId := range disksIds {
		disk, err := manager.GetDisk(diskId.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		err = vm.DetachDisk(disk)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = vm.Delete()
	if err != nil {
		return diag.Errorf("Error deleting vm: %s", err)
	}

	return nil
}

func createPort(d *schema.ResourceData, manager *rustack.Manager, portPrefix *string) (*rustack.Port, error) {
	portNetwork, err := GetNetworkById(d, manager, portPrefix)
	if err != nil {
		return nil, err
	}

	firewallsCount := d.Get(MakePrefix(portPrefix, "firewall_templates.#")).(int)
	firewalls := make([]*rustack.FirewallTemplate, firewallsCount)
	firewallsResourceData := d.Get(MakePrefix(portPrefix, "firewall_templates")).(*schema.Set).List()
	for j, firewallId := range firewallsResourceData {
		portFirewall, err := manager.GetFirewallTemplate(firewallId.(string))
		if err != nil {
			return nil, err
		}
		firewalls[j] = portFirewall
	}
	ipAddressStr := d.Get(MakePrefix(portPrefix, "ip_address")).(string)
	ipAddress := &ipAddressStr
	if ipAddressStr == "" {
		ipAddress = nil
	}

	newPort := rustack.NewPort(portNetwork, firewalls, ipAddress)
	return &newPort, nil
}

func syncDisks(d *schema.ResourceData, manager *rustack.Manager, vdc *rustack.Vdc, vm *rustack.Vm) (diagErr diag.Diagnostics) {
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("ERROR. Something wrong with Vdc: %s", err)
	}

	// List disks to join
	diagErr = attachNewDisk(d, manager, vm)
	if diagErr != nil {
		return
	}

	// Which disks are present on vm and not mentioned in the state?
	// Detach disks
	diagErr = detachOldDisk(d, manager, vm)
	if diagErr != nil {
		return
	}

	// System disk resize
	if d.HasChange("system_disk") {
		systemDiskArgs := d.Get("system_disk.0").(map[string]interface{})
		systemDiskId := systemDiskArgs["id"].(string)
		diskSize := systemDiskArgs["size"].(int)
		systemDisk, err := manager.GetDisk(systemDiskId)
		if err != nil {
			return diag.Errorf("ERROR. Wrong system disk id: %s", err)
		}

		if err = systemDisk.Resize(diskSize); err != nil {
			return diag.Errorf("Error resizing disk: %s", err)
		}

		if !d.HasChange("system_disk.0.storage_profile_id") {
			return
		}
		storageProfileId := d.Get("system_disk.0.storage_profile_id").(string)
		storageProfile, err := targetVdc.GetStorageProfile(storageProfileId)
		if err != nil {
			return diag.Errorf("Error getting storage profile: %s", err)
		}

		err = systemDisk.UpdateStorageProfile(*storageProfile)
		if err != nil {
			return diag.Errorf("Error updating storage: %s", err)
		}
	}

	return
}

func syncPorts(d *schema.ResourceData, manager *rustack.Manager, vdc *rustack.Vdc, vm *rustack.Vm) diag.Diagnostics {
	// Delete disconnected ports and create a new if connected
	portList := d.Get("port").([]interface{})
	if portList != nil {
		return nil
	}
	if err := manageVmPorts(d, manager); err != nil {
		return diag.FromErr(err)
	}

	// Detect port changes for found ports with the same id
	if diagErr := updateVmPorts(d, manager); diagErr != nil {
		return diagErr
	}

	return nil
}

func manageVmPorts(d *schema.ResourceData, manager *rustack.Manager) (err error) {
	needReload := false
	ports, err := connectVmPorts(d, manager)
	if err != nil {
		return err
	}

	vm, err := manager.GetVm(d.Id())
	if err != nil {
		return err
	}
	// disconnect ports
	oldPortList := vm.Ports
	newPortList := ports
	for _, old_port := range oldPortList {
		found := false
		for _, portNew := range newPortList {
			portId := portNew.ID
			if portId == old_port.ID {
				found = true
				break
			}
		}

		if !found {
			log.Printf("Port %s found on vm and not mentioned in the state. Port will be deleted", old_port.ID)
			if err := old_port.ForceDelete(); err != nil {
				return err
			}

			needReload = true
		}
	}
	if needReload {
		if err := vm.Reload(); err != nil {
			return err
		}
	}

	return nil
}

func connectVmPorts(d *schema.ResourceData, manager *rustack.Manager) (ports []*rustack.Port, err error) {
	portList := d.Get("port").([]interface{})
	vm, err := manager.GetVm(d.Id())
	if err != nil {
		return nil, err
	}

	for i, portNew := range portList {
		portId := portNew.(map[string]interface{})["id"].(string)
		portPrefix := fmt.Sprint("port.", i)

		found := false
		for _, port := range vm.Ports {
			if port.ID == portId {
				found = true
				ports = append(ports, port)
				break
			}
		}

		if !found {
			log.Printf("Port %s found in the state and is not present on vm. Port will be created", portPrefix)

			newPort, err := createPort(d, manager, &portPrefix)
			if err != nil {
				return nil, err
			}

			f := func() error { return vm.AddPort(newPort) }
			if err := repeatOnError(f, vm); err != nil {
				return nil, err
			}

			if err := repeatOnError(vm.Reload, vm); err != nil {
				return nil, err
			}
			ports = append(ports, newPort)
		}
	}

	return
}

func updateVmPorts(d *schema.ResourceData, manager *rustack.Manager) diag.Diagnostics {
	vm, err := manager.GetVm(d.Id())
	if err != nil {
		return diag.Errorf("Error getting vm: %s", err)
	}

	for _, port := range vm.Ports {
		var portExists bool
		var portId string
		var portPrefix string
		portList := d.Get("port").([]interface{})
		for i, p := range portList {
			portPrefix = fmt.Sprint("port.", i)
			portId = p.(map[string]interface{})["id"].(string)
			if portId == port.ID {
				portExists = true
				break
			}
		}
		if !portExists {
			// That case has been resolved above
			continue
		}

		// Compare port
		pseudoPort, err := createPort(d, manager, &portPrefix)
		if err != nil {
			return diag.FromErr(err)
		}

		if err = port.UpdateFirewall(pseudoPort.FirewallTemplates); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func attachNewDisk(d *schema.ResourceData, manager *rustack.Manager, vm *rustack.Vm) (diagErr diag.Diagnostics) {
	disksIds := d.Get("disks").(*schema.Set).List()
	systemDiskResource := d.Get("system_disk.0")
	systemDisk := systemDiskResource.(map[string]interface{})["id"].(string)
	var needReload bool
	disksIds = append(disksIds, systemDisk)

	for _, diskId := range disksIds {
		found := false
		for _, disk := range vm.Disks {
			if diskId == disk.ID {
				found = true
				break
			}
		}

		if !found {
			disk, err := manager.GetDisk(diskId.(string))
			if err != nil {
				diagErr = diag.FromErr(err)
				return
			}
			log.Printf("Disk `%s` will be Attached", disk.ID)
			if err = vm.AttachDisk(disk); err != nil {
				diagErr = diag.Errorf("ERROR. Cannot attach disk `%s`: %s", disk.ID, err)
				return
			}
			needReload = true
		}
	}

	if needReload {
		if err := vm.Reload(); err != nil {
			return diag.FromErr(err)
		}
	}

	return
}

func detachOldDisk(d *schema.ResourceData, manager *rustack.Manager, vm *rustack.Vm) (diagErr diag.Diagnostics) {
	disksIds := d.Get("disks").(*schema.Set).List()
	systemDiskResource := d.Get("system_disk.0")
	systemDisk := systemDiskResource.(map[string]interface{})["id"].(string)
	var needReload bool
	disksIds = append(disksIds, systemDisk)

	for _, disk := range vm.Disks {
		found := false
		for _, diskId := range disksIds {
			if diskId == disk.ID {
				found = true
				break
			}
		}

		if !found {
			log.Printf("Disk %s found on vm and not mentioned in the state."+
				" Disk will be detached", disk.ID)
			vm.DetachDisk(disk)
			needReload = true
		}
	}

	if needReload {
		if err := vm.Reload(); err != nil {
			return diag.FromErr(err)
		}
	}

	return
}

func resourceRustackVmPortRead(vm *rustack.Vm, d *schema.ResourceData) (diagErr diag.Diagnostics) {
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

	return
}
