package rustack_terraform

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pilat/rustack-go/rustack"
)

func dataSourceRustackVdc() *schema.Resource {
	args := Defaults()
	args.injectResultVdc()
	args.injectContextProjectByIdOptional()
	args.injectContextGetVdc() // override name

	return &schema.Resource{
		ReadContext: dataSourceRustackVdcRead,
		Schema:      args,
	}
}

func dataSourceRustackVdcRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()

	var targetProject *rustack.Project

	if _, exists := d.GetOk("project_id"); exists {
		project, err := GetProjectById(d, manager)
		if err != nil {
			return diag.Errorf("Error getting project: %s", err)
		}

		targetProject = project
	}
	
	target, err := checkDatasourceNameOrId(d)
	if err != nil {
		return diag.Errorf("Error getting VDC: %s", err)
	}
	var targetVdc *rustack.Vdc
	if target == "id" {
		targetVdc, err = manager.GetVdc(d.Get("id").(string))
		if err != nil {
			return diag.Errorf("Error getting VDC: %s", err)
		}
	} else {
		targetVdc, err = GetVdcByName(d, manager, targetProject)
		if err != nil {
			return diag.Errorf("Error getting VDC: %s", err)
		}
	}

	flattenedVdc := map[string]interface{}{
		"id":              targetVdc.ID,
		"name":            targetVdc.Name,
		"hypervisor":      targetVdc.Hypervisor.Name,
		"hypervisor_type": targetVdc.Hypervisor.Type,
	}

	if err := setResourceDataFromMap(d, flattenedVdc); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(targetVdc.ID)
	return nil
}
