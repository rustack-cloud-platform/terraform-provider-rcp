package rustack_terraform

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRustackVm() *schema.Resource {
	args := Defaults()
	args.injectResultVm()
	args.injectContextVdcById()
	args.injectContextVmByName() // override "name"

	return &schema.Resource{
		ReadContext: dataSourceRustackVmRead,
		Schema:      args,
	}
}

func dataSourceRustackVmRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting vdc: %s", err)
	}

	targetVm, err := GetVmByName(d, manager, targetVdc)
	if err != nil {
		return diag.Errorf("Error getting vm: %s", err)
	}

	flatten := map[string]interface{}{
		"id":            targetVm.ID,
		"name":          targetVm.Name,
		"cpu":           targetVm.Cpu,
		"ram":           targetVm.Ram,
		"template_id":   targetVm.Template.ID,
		"template_name": targetVm.Template.Name,
		"power":         targetVm.Power,
		"floating":      nil,
		"floating_ip":   nil,
	}

	if targetVm.Floating != nil {
		flatten["floating"] = true
		flatten["floating_ip"] = targetVm.Floating.IpAddress
	}

	if err := setResourceDataFromMap(d, flatten); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(targetVm.ID)
	return nil
}
