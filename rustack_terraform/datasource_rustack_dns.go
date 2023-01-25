package rustack_terraform

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRustackDns() *schema.Resource {
	args := Defaults()
	args.injectContextProjectById()
	args.injectResultDns()
	args.injectContextDnsByName() // override name

	return &schema.Resource{
		ReadContext: dataSourceRustackDnsRead,
		Schema:      args,
	}
}

func dataSourceRustackDnsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()

	targetDns, err := GetDnsByName(d, manager)
	if err != nil {
		return diag.Errorf("Error getting dns: %s", err)
	}

	flatten := map[string]interface{}{
		"id":      targetDns.ID,
		"name":    targetDns.Name,
		"project_id": targetDns.Project.ID,
	}

	if err := setResourceDataFromMap(d, flatten); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(targetDns.ID)
	return nil
}
