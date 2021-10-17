package rustack_terraform

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pilat/rustack-go/rustack"
)

func resourceRustackDisk() *schema.Resource {
	args := Defaults()
	args.injectCreateDisk()
	args.injectContextVdcById()
	args.injectContextStorageProfileById() // override storage_profile_id

	return &schema.Resource{
		CreateContext: resourceRustackDiskCreate,
		ReadContext:   resourceRustackDiskRead,
		UpdateContext: resourceRustackDiskUpdate,
		DeleteContext: resourceRustackDiskDelete,
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

func resourceRustackDiskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting VDC: %s", err)
	}

	targetStorageProfile, err := GetStorageProfileById(d, manager, targetVdc, nil)
	if err != nil {
		return diag.Errorf("Error getting storage profile: %s", err)
	}

	newDisk := rustack.NewDisk(d.Get("name").(string), d.Get("size").(int), targetStorageProfile)
	err = targetVdc.CreateDisk(&newDisk)
	if err != nil {
		return diag.Errorf("Error creating disk: %s", err)
	}

	d.SetId(newDisk.ID)
	log.Printf("[INFO] Disk created, ID: %s", d.Id())

	return resourceRustackDiskRead(ctx, d, meta)
}

func resourceRustackDiskRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	disk, err := manager.GetDisk(d.Id())
	if err != nil {
		return diag.Errorf("Error getting disk: %s", err)
	}

	d.SetId(disk.ID)
	d.Set("name", disk.Name)
	d.Set("size", disk.Size)

	return nil
}

func resourceRustackDiskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()

	disk, err := manager.GetDisk(d.Id())
	if err != nil {
		return diag.Errorf("Error getting disk: %s", err)
	}

	if d.HasChange("name") {
		err = disk.Rename(d.Get("name").(string))
		if err != nil {
			return diag.Errorf("Error rename disk: %s", err)
		}
	}

	if d.HasChange("size") {
		err = disk.Resize(d.Get("size").(int))
		if err != nil {
			return diag.Errorf("Error resizing disk: %s", err)
		}
	}

	if d.HasChange("storage_profile_id") {
		targetVdc, err := GetVdcById(d, manager)
		if err != nil {
			return diag.Errorf("Error getting VDC: %s", err)
		}

		targetStorageProfile, err := GetStorageProfileById(d, manager, targetVdc, nil)
		if err != nil {
			return diag.Errorf("Error getting storage profile: %s", err)
		}

		err = disk.UpdateStorageProfile(*targetStorageProfile)
		if err != nil {
			return diag.Errorf("Error updating storage: %s", err)
		}
	}

	return resourceRustackDiskRead(ctx, d, meta)
}

func resourceRustackDiskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	disk, err := manager.GetDisk(d.Id())
	if err != nil {
		return diag.Errorf("Error getting disk: %s", err)
	}

	err = disk.Delete()
	if err != nil {
		return diag.Errorf("Error deleting disk: %s", err)
	}

	return nil
}
