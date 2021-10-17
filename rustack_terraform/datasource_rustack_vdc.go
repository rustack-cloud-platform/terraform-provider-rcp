package rustack_terraform

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRustackVdc() *schema.Resource {
	args := Defaults()
	args.injectResultVdc()
	args.injectContextProjectById()
	args.injectContextVdcByName() // override name

	return &schema.Resource{
		ReadContext: dataSourceRustackVdcRead,
		Schema:      args,
	}
}

func dataSourceRustackVdcRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetProject, err := GetProjectById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting project: %s", err)
	}

	targetVdc, err := GetVdcByName(d, manager, targetProject)
	if err != nil {
		return diag.Errorf("Error getting VDC: %s", err)
	}

	flattenedVdc := map[string]interface{}{
		"id": targetVdc.ID,
		// "vdc":             targetVdc.Name,  // ??
		// "vdc":             nil,
		"name":            targetVdc.Name,
		"hypervisor":      targetVdc.Hypervisor.Name,
		"hypervisor_type": targetVdc.Hypervisor.Type,
		// "project":         targetVdc.Project.Name,
	}

	if err := setResourceDataFromMap(d, flattenedVdc); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(targetVdc.ID)
	return nil
}
