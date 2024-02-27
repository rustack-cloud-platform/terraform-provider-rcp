package rustack_terraform

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rustack-cloud-platform/rcp-go/rustack"
)

func dataSourceRustackLbaas() *schema.Resource {
	args := Defaults()
	args.injectContextVdcById()
	args.injectResultLbaas()
	args.injectContextGetLbaas()

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

	target, err := checkDatasourceNameOrId(d)
	if err != nil {
		return diag.Errorf("Error getting Lbaas: %s", err)
	}
	var targetLbaas *rustack.LoadBalancer
	if target == "id" {
		targetLbaas, err = manager.GetLoadBalancer(d.Get("id").(string))
		if err != nil {
			return diag.Errorf("Error getting Lbaas: %s", err)
		}
	} else {
		targetLbaas, err = GetLbaasByName(d, manager, targetVdc)
		if err != nil {
			return diag.Errorf("Error getting Lbaas: %s", err)
		}
	}

	flatten := map[string]interface{}{
		"id":   targetLbaas.ID,
		"name": targetLbaas.Name,
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
