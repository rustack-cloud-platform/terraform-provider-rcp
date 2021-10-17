package rustack_terraform

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRustackFirewallTemplate() *schema.Resource {
	args := Defaults()
	args.injectResultTemplate()
	args.injectContextVdcById()
	args.injectContextTemplateByName() // override name

	return &schema.Resource{
		ReadContext: dataSourceRustackFirewallTemplateRead,
		Schema:      args,
	}
}

func dataSourceRustackFirewallTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting vdc: %s", err)
	}

	targetTemplate, err := GetFirewallTemplateByName(d, manager, targetVdc)
	if err != nil {
		return diag.Errorf("Error getting template: %s", err)
	}

	flatten := map[string]interface{}{
		"id":   targetTemplate.ID,
		"name": targetTemplate.Name,
	}

	if err := setResourceDataFromMap(d, flatten); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(targetTemplate.ID)
	return nil
}
