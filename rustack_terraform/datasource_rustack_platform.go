package rustack_terraform

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rustack-cloud-platform/rcp-go/rustack"
)

func dataSourceRustackPlatform() *schema.Resource {
	args := Defaults()
	args.injectResultPlatform()
	args.injectContextVdcById()
	args.injectContextGetPlatform() // override name

	return &schema.Resource{
		ReadContext: dataSourceRustackPlatformRead,
		Schema:      args,
	}
}

func dataSourceRustackPlatformRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()

	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting vdc: %s", err)
	}

	target, err := checkDatasourceNameOrId(d)
	if err != nil {
		return diag.Errorf("Error getting Platform: %s", err)
	}
	var targetPlatform *rustack.Platform
	if target == "id" {
		targetPlatform, err = manager.GetPlatform(d.Get("id").(string))
		if err != nil {
			return diag.Errorf("Error getting Platform: %s", err)
		}
	} else {
		targetPlatform, err = GetPlatformByName(d, manager, targetVdc)
		if err != nil {
			return diag.Errorf("Error getting Platform: %s", err)
		}
	}

	flatten := map[string]interface{}{
		"id":   targetPlatform.ID,
		"name": targetPlatform.Name,
	}

	if err := setResourceDataFromMap(d, flatten); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(targetPlatform.ID)
	return nil
}
