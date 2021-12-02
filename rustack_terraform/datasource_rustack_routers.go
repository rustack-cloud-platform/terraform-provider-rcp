package rustack_terraform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/hashstructure/v2"
)

func dataSourceRustackRouters() *schema.Resource {
	args := Defaults()
	args.injectContextVdcById()
	args.injectResultListRouter()

	return &schema.Resource{
		ReadContext: dataSourceRustackRoutersRead,
		Schema:      args,
	}
}

func dataSourceRustackRoutersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	vdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Unable to get vdc: %s", err)
	}

	allRouters, err := vdc.GetRouters()
	if err != nil {
		return diag.Errorf("Error getting routers: %s", err)
	}

	routersMap := make([]map[string]interface{}, len(allRouters))
	for i, project := range allRouters {
		routersMap[i] = map[string]interface{}{
			"id":   project.ID,
			"name": project.Name,
		}
	}

	hash, err := hashstructure.Hash(allRouters, hashstructure.FormatV2, nil)
	if err != nil {
		return diag.Errorf("unable to set `routers` attribute: %s", err)
	}

	d.SetId(fmt.Sprintf("routers/%d", hash))

	if err := d.Set("routers", routersMap); err != nil {
		return diag.Errorf("unable to set `routers` attribute: %s", err)
	}

	return nil
}
