package rustack_terraform

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pilat/rustack-go/rustack"
)

func dataSourceRustackDisk() *schema.Resource {
	args := Defaults()
	args.injectContextVdcById()
	args.injectResultDisk()
	args.injectContextGetDisk() // override name

	return &schema.Resource{
		ReadContext: dataSourceRustackDiskRead,
		Schema:      args,
	}
}

func dataSourceRustackDiskRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting vdc: %s", err)
	}

	target, err := checkDatasourceNameOrId(d)
	if err != nil {
		return diag.Errorf("Error getting disk: %s", err)
	}
	var targetDisk *rustack.Disk
	if target == "id" {
		targetDisk, err = manager.GetDisk(d.Get("id").(string))
		if err != nil {
			return diag.Errorf("Error getting disk: %s", err)
		}
	} else {
		targetDisk, err = GetDiskByName(d, manager, targetVdc)
		if err != nil {
			return diag.Errorf("Error getting disk: %s", err)
		}
	}

	flatten := map[string]interface{}{
		"id":                   targetDisk.ID,
		"name":                 targetDisk.Name,
		"size":                 targetDisk.Size,
		"storage_profile_id":   targetDisk.StorageProfile.ID,
		"storage_profile_name": targetDisk.StorageProfile.Name,
	}

	if err := setResourceDataFromMap(d, flatten); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(targetDisk.ID)
	return nil
}
