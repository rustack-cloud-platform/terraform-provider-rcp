package rustack_terraform

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rustack-cloud-platform/rcp-go/rustack"
)

func dataSourceRustackDns() *schema.Resource {
	args := Defaults()
	args.injectContextProjectById()
	args.injectResultDns()
	args.injectContextGetDns() // override name

	return &schema.Resource{
		ReadContext: dataSourceRustackDnsRead,
		Schema:      args,
	}
}

func dataSourceRustackDnsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()

	target, err := checkDatasourceNameOrId(d)
	if err != nil {
		return diag.Errorf("Error getting dns: %s", err)
	}
	var targetDns *rustack.Dns
	if target == "id" {
		targetDns, err = manager.GetDns(d.Get("id").(string))
		if err != nil {
			return diag.Errorf("Error getting dns: %s", err)
		}
	} else {
		targetDns, err = GetDnsByName(d, manager)
		if err != nil {
			return diag.Errorf("Error getting dns: %s", err)
		}
	}

	flatten := map[string]interface{}{
		"id":         targetDns.ID,
		"name":       targetDns.Name,
		"project_id": targetDns.Project.ID,
	}

	if err := setResourceDataFromMap(d, flatten); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(targetDns.ID)
	return nil
}
