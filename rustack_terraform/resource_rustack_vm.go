package rustack_terraform

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rustack-cloud-platform/rcp-go/rustack"
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

func getVmPortsIds(d *schema.ResourceData) (portsIds []string) {
	if d.HasChange("ports") {
		portsIdsValue := d.Get("ports").([]interface{})
		portsIds = make([]string, 0, len(portsIdsValue))
		for _, portIdValue := range portsIdsValue {
			portsIds = append(portsIds, portIdValue.(string))
		}
	} else {
		networks := d.Get("networks").([]interface{})
		portsIds = make([]string, 0, len(networks))
		for _, network := range networks {
			portMap := network.(map[string]interface{})
			portsIds = append(portsIds, portMap["id"].(string))
		}
	}
	return
}

func resourceRustackVmCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("vdc_id: Error getting VDC: %s", err)
	}

	template, err := GetTemplateById(d, manager)
	if err != nil {
		return diag.Errorf("template_id: Error getting template: %s", err)
	}

	vmName := d.Get("name").(string)
	cpu := d.Get("cpu").(int)
	ram := d.Get("ram").(float64)
	userData := d.Get("user_data").(string)
	log.Printf(vmName, cpu, ram, userData, template.Name)

	// System disk creation
	systemDiskArgs := d.Get("system_disk.0").(map[string]interface{})
	diskSize := systemDiskArgs["size"].(int)
	storageProfileId := systemDiskArgs["storage_profile_id"].(string)

	storageProfile, err := targetVdc.GetStorageProfile(storageProfileId)
	if err != nil {
		return diag.Errorf("storage_profile_id: Error storage profile %s not found", storageProfileId)
	}

	systemDiskList := make([]*rustack.Disk, 1)
	newDisk := rustack.NewDisk("Основной диск", diskSize, storageProfile)
	systemDiskList[0] = &newDisk

	portsIds := getVmPortsIds(d)
	ports := make([]*rustack.Port, len(portsIds))
	for i, portId := range portsIds {
		port, err := manager.GetPort(portId)
		if err != nil {
			return diag.FromErr(err)
		}
		ports[i] = port
	}
	var floatingIp *string = nil
	if d.Get("floating").(bool) {
		floatingIpStr := "RANDOM_FIP"
		floatingIp = &floatingIpStr
	}

	newVm := rustack.NewVm(vmName, cpu, ram, template, nil, &userData, ports,
		systemDiskList, floatingIp)
	newVm.Tags = unmarshalTagNames(d.Get("tags"))

	err = targetVdc.CreateVm(&newVm)
	if err != nil {
		return diag.Errorf("Error creating vm: %s", err)
	}

	newVm.WaitLock()
	vm_power := d.Get("power").(bool)
	if !vm_power {
		newVm.PowerOff()
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

func resourceRustackVmRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	vm, err := manager.GetVm(d.Id())
	if err != nil {
		if err.(*rustack.RustackApiError).Code() == 404 {
			d.SetId("")
			return nil
		} else {
			return diag.Errorf("id: Error getting vm: %s", err)
		}
	}

	d.SetId(vm.ID)
	d.Set("name", vm.Name)
	d.Set("cpu", vm.Cpu)
	d.Set("ram", vm.Ram)
	d.Set("template_id", vm.Template.ID)
	d.Set("power", vm.Power)

	flattenDisks := make([]string, len(vm.Disks)-1)
	for i, disk := range vm.Disks {
		if i == 0 {
			systemDisk := make([]interface{}, 1)
			systemDisk[0] = map[string]interface{}{
				"id":                 disk.ID,
				"name":               "Основной диск",
				"size":               disk.Size,
				"storage_profile_id": disk.StorageProfile.ID,
				"external_id":        disk.ExternalID,
			}

			d.Set("system_disk", systemDisk)
			continue
		}
		flattenDisks[i-1] = disk.ID
	}
	d.Set("disks", flattenDisks)

	flattenPorts := make([]string, len(vm.Ports))
	flattenNetworks := make([]map[string]interface{}, 0, len(vm.Ports))
	for i, port := range vm.Ports {
		flattenPorts[i] = port.ID
		flattenNetworks = append(flattenNetworks, map[string]interface{}{
			"id":         port.ID,
			"ip_address": port.IpAddress,
		})
	}
	d.Set("ports", flattenPorts)
	d.Set("networks", flattenNetworks)

	d.Set("floating", vm.Floating != nil)
	d.Set("floating_ip", "")
	if vm.Floating != nil {
		d.Set("floating_ip", vm.Floating.IpAddress)
	}
	d.Set("tags", marshalTagNames(vm.Tags))

	return nil
}

func resourceRustackVmUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("vdc_id: Error getting VDC: %s", err)
	}

	hasFlavorChanged := false
	needUpdate := false

	vm, err := manager.GetVm(d.Id())
	if err != nil {
		return diag.Errorf("id: Error getting vm: %s", err)
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
		vm.Ram = d.Get("ram").(float64)
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
	if d.HasChange("tags") {
		needUpdate = true
		vm.Tags = unmarshalTagNames(d.Get("tags"))
	}

	if needUpdate {
		if err := repeatOnError(vm.Update, vm); err != nil {
			return diag.Errorf("Error updating vm: %s", err)
		}
	}

	if needPowerOn {
		vm.PowerOn()
	}
	a := d.Get("power").(bool)
	if a {
		vm.PowerOn()
	} else {
		vm.PowerOff()
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
		return diag.Errorf("id: Error getting vm: %s", err)
	}

	vm.Floating = &rustack.Port{IpAddress: nil}
	if err := repeatOnError(vm.Update, vm); err != nil {
		return diag.Errorf("Error updating vm: %s", err)
	}
	vm.WaitLock()

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

	portsIds := getVmPortsIds(d)
	for _, portId := range portsIds {
		port, err := manager.GetPort(portId)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := vm.DisconnectPort(port); err != nil {
			return diag.FromErr(err)
		}
	}
	vm.WaitLock()

	err = vm.Delete()
	if err != nil {
		return diag.Errorf("Error deleting vm: %s", err)
	}
	vm.WaitLock()

	return nil
}

func syncDisks(d *schema.ResourceData, manager *rustack.Manager, vdc *rustack.Vdc, vm *rustack.Vm) (diagErr diag.Diagnostics) {
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("vdc_id: Error getting VDC: %s", err)
	}

	// Which disks are present on vm and not mentioned in the state?
	// Detach disks
	diagErr = detachOldDisk(d, manager, vm)
	if diagErr != nil {
		return
	}

	// List disks to join
	diagErr = attachNewDisk(d, manager, vm)
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
			return diag.Errorf("system_disk: Error getting system disk id: %s", err)
		}

		if err = systemDisk.Resize(diskSize); err != nil {
			return diag.Errorf("size: Error resizing disk: %s", err)
		}

		if !d.HasChange("system_disk.0.storage_profile_id") {
			return
		}

		storageProfileId := d.Get("system_disk.0.storage_profile_id").(string)
		storageProfile, err := targetVdc.GetStorageProfile(storageProfileId)
		if err != nil {
			return diag.Errorf("storage_profile_id: Error getting storage profile: %s", err)
		}

		err = systemDisk.UpdateStorageProfile(*storageProfile)
		if err != nil {
			return diag.Errorf("Error updating storage: %s", err)
		}
	}

	return
}

func syncPorts(d *schema.ResourceData, manager *rustack.Manager, vdc *rustack.Vdc, vm *rustack.Vm) (diagErr diag.Diagnostics) {

	// Delete ConnectNewPort ports and create a new if connected
	diagErr = DisconnectOldPort(d, manager, vm)
	if diagErr != nil {
		return

	}

	diagErr = ConnectNewPort(d, manager, vm)
	if diagErr != nil {
		return
	}

	return
}

func ConnectNewPort(d *schema.ResourceData, manager *rustack.Manager, vm *rustack.Vm) (diagErr diag.Diagnostics) {
	portsIds := getVmPortsIds(d)
	for _, portId := range portsIds {
		found := false
		for _, port := range vm.Ports {
			if port.ID == portId {
				found = true
				break
			}
		}

		if !found {
			port, err := manager.GetPort(portId)

			if err != nil {
				diagErr = diag.FromErr(err)
				return
			}
			if port.Connected != nil && port.Connected.ID != vm.ID {

				if err := vm.DisconnectPort(port); err != nil {
					return diag.FromErr(err)
				}
				vm.WaitLock()
			}
			log.Printf("Port `%s` will be Attached", port.ID)

			if err := vm.ConnectPort(port, true); err != nil {
				diagErr = diag.Errorf("Ports: Error Cannot attach port `%s`: %s", port.ID, err)
				return
			}
		}
	}
	return
}

func DisconnectOldPort(d *schema.ResourceData, manager *rustack.Manager, vm *rustack.Vm) diag.Diagnostics {
	portsIds := getVmPortsIds(d)
	for _, port := range vm.Ports {
		found := false
		for _, portId := range portsIds {
			if portId == port.ID {
				found = true
				break
			}
		}
		if !found {
			if port.Connected != nil && port.Connected.ID == vm.ID {
				log.Printf("Port %s found on vm and not mentioned in the state."+
					" Port will be detached", port.ID)

				if err := vm.DisconnectPort(port); err != nil {
					return diag.FromErr(err)
				}
				vm.WaitLock()
			}
		}
	}

	return nil
}

func attachNewDisk(d *schema.ResourceData, manager *rustack.Manager, vm *rustack.Vm) (diagErr diag.Diagnostics) {
	disksIds := d.Get("disks").(*schema.Set).List()
	// Save system_disk
	systemDiskResource := d.Get("system_disk.0")
	systemDisk := systemDiskResource.(map[string]interface{})["id"].(string)
	var needReload bool
	disksIds = append(disksIds, systemDisk)
	vm_id := vm.ID

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
			if disk.Vm != nil && disk.Vm.ID != vm_id {
				log.Printf("Disk %s found on other vm and will be detached for attached to vm.", disk.ID)
				vm.DetachDisk(disk)
				if err := vm.Reload(); err != nil {
					return diag.FromErr(err)
				}
				vm.WaitLock()
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
	vm_id := vm.ID

	for _, disk := range vm.Disks {
		found := false
		for _, diskId := range disksIds {
			if diskId == disk.ID {
				found = true
				break
			}
		}

		if !found {
			disk, err := manager.GetDisk(disk.ID)
			if err != nil {
				diagErr = diag.FromErr(err)
				return
			}
			if disk.Vm != nil && disk.Vm.ID == vm_id {
				log.Printf("Disk %s found on vm and not mentioned in the state."+
					" Disk will be detached", disk.ID)
				vm.DetachDisk(disk)
				needReload = true
			}
		}
	}

	if needReload {
		if err := vm.Reload(); err != nil {
			return diag.FromErr(err)
		}
	}

	return
}
