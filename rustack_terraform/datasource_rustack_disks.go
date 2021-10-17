package rustack_terraform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/hashstructure/v2"
)

func dataSourceRustackDisks() *schema.Resource {
	args := Defaults()
	args.injectContextVdcById()
	args.injectResultListDisk()

	return &schema.Resource{
		ReadContext: dataSourceRustackDisksRead,
		Schema:      args,
	}
}

func dataSourceRustackDisksRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting vdc: %s", err)
	}

	allDisks, err := targetVdc.GetDisks()
	if err != nil {
		return diag.Errorf("Error retrieving disks: %s", err)
	}

	flattenedRecords := make([]map[string]interface{}, len(allDisks))
	for i, disk := range allDisks {
		flattenedRecords[i] = map[string]interface{}{
			"id":                   disk.ID,
			"name":                 disk.Name,
			"size":                 disk.Size,
			"storage_profile_id":   disk.StorageProfile.ID,
			"storage_profile_name": disk.StorageProfile.Name,
		}
	}

	hash, err := hashstructure.Hash(allDisks, hashstructure.FormatV2, nil)
	if err != nil {
		diag.Errorf("unable to set `disks` attribute: %s", err)
	}

	d.SetId(fmt.Sprintf("disks/%d", hash))
	// d.Set("vdc_id", nil)
	// d.Set("vdc_name", nil)

	if err := d.Set("disks", flattenedRecords); err != nil {
		return diag.Errorf("unable to set `disks` attribute: %s", err)
	}

	return nil
}
