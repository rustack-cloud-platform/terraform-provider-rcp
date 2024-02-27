package rustack_terraform

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rustack-cloud-platform/rcp-go/rustack"
)

func dataSourceRustackRouter() *schema.Resource {
	args := Defaults()
	args.injectContextVdcById()
	args.injectResultRouter()
	args.injectContextGetRouter()

	return &schema.Resource{
		ReadContext: dataSourceRustackRouterRead,
		Schema:      args,
	}
}

func dataSourceRustackRouterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	target, err := checkDatasourceNameOrId(d)
	if err != nil {
		return diag.Errorf("Error getting router: %s", err)
	}
	var router *rustack.Router
	if target == "id" {
		router, err = manager.GetRouter(d.Get("id").(string))
		if err != nil {
			return diag.Errorf("Error getting router: %s", err)
		}
	} else {
		router, err = GetRouterByName(d, manager)
		if err != nil {
			return diag.Errorf("Error getting router: %s", err)
		}
	}

	routerMap := map[string]interface{}{
		"id":   router.ID,
		"name": router.Name,
	}

	if err := setResourceDataFromMap(d, routerMap); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(router.ID)
	return nil
}
