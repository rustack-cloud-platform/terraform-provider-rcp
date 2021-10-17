package rustack_terraform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/hashstructure/v2"
)

func dataSourceRustackStorageProfiles() *schema.Resource {
	args := Defaults()
	args.injectContextVdcById()
	args.injectResultListStorageProfile()

	return &schema.Resource{
		ReadContext: dataSourceRustackStorageProfilesRead,
		Schema:      args,
	}
}

func dataSourceRustackStorageProfilesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting VDC: %s", err)
	}

	storageProfiles, err := targetVdc.GetStorageProfiles()
	if err != nil {
		return diag.Errorf("Error getting storage profiles")
	}

	flattenedStorageProfiles := make([]map[string]interface{}, len(storageProfiles))
	for i, storageProfile := range storageProfiles {
		flattenedStorageProfiles[i] = map[string]interface{}{
			"id":   storageProfile.ID,
			"name": storageProfile.Name,
		}
	}

	hash, err := hashstructure.Hash(storageProfiles, hashstructure.FormatV2, nil)
	if err != nil {
		diag.Errorf("unable to set `networks` attribute: %s", err)
	}

	d.SetId(fmt.Sprintf("storage_profiles/%d", hash))
	// d.Set("vdc_id", nil)
	// d.Set("vdc_name", nil)

	if err := d.Set("storage_profiles", flattenedStorageProfiles); err != nil {
		return diag.Errorf("unable to set `storage_profiles` attribute: %s", err)
	}

	return nil
}
