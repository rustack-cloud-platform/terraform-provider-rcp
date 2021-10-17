package rustack_terraform

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRustackStorageProfile() *schema.Resource {
	args := Defaults()
	args.injectResultStorageProfile()
	args.injectContextVdcById()
	args.injectContextStorageProfileByName() // override name

	return &schema.Resource{
		ReadContext: dataSourceRustackStorageProfileRead,
		Schema:      args,
	}
}

func dataSourceRustackStorageProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting vdc: %s", err)
	}

	targetStorageProfile, err := GetStorageProfileByName(d, manager, targetVdc)
	if err != nil {
		return diag.Errorf("Error getting storage profile: %s", err)
	}

	flatten := map[string]interface{}{
		"id":   targetStorageProfile.ID,
		"name": targetStorageProfile.Name,
	}

	if err := setResourceDataFromMap(d, flatten); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(targetStorageProfile.ID)
	return nil
}
