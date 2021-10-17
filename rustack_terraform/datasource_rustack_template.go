package rustack_terraform

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRustackTemplate() *schema.Resource {
	args := Defaults()
	args.injectResultTemplate()
	args.injectContextVdcById()
	args.injectContextTemplateByName() // override name

	return &schema.Resource{
		ReadContext: dataSourceRustackTemplateRead,
		Schema:      args,
	}
}

func dataSourceRustackTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting vdc: %s", err)
	}

	targetTemplate, err := GetTemplateByName(d, manager, targetVdc)
	if err != nil {
		return diag.Errorf("Error getting template: %s", err)
	}

	flatten := map[string]interface{}{
		"id":       targetTemplate.ID,
		"name":     targetTemplate.Name,
		"min_cpu":  targetTemplate.MinCpu,
		"min_ram":  targetTemplate.MinRam,
		"min_disk": targetTemplate.MinHdd,
	}

	if err := setResourceDataFromMap(d, flatten); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(targetTemplate.ID)
	return nil
}
