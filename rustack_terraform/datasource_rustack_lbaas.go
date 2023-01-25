package rustack_terraform

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRustackLbaas() *schema.Resource {
	args := Defaults()
	args.injectContextVdcById()
	args.injectResultLbaas()
	args.injectContextLbaasByName()

	return &schema.Resource{
		ReadContext: dataSourceRustackLbaasRead,
		Schema:      args,
	}
}

func dataSourceRustackLbaasRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting vdc: %s", err)
	}
	targetLbaas, err := GetLbaasByName(d, manager, targetVdc)
	if err != nil {
		return diag.Errorf("Error getting Lbaas: %s", err)
	}

	flatten := map[string]interface{}{
		"id":            targetLbaas.ID,
		"name":          targetLbaas.Name,
	}

	if targetLbaas.Floating != nil {
		flatten["floating"] = true
		flatten["floating_ip"] = targetLbaas.Floating.IpAddress
	}

	if err := setResourceDataFromMap(d, flatten); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(targetLbaas.ID)
	return nil
}
