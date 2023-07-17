package rustack_terraform

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pilat/rustack-go/rustack"
)

func dataSourceRustackProject() *schema.Resource {
	args := Defaults()
	args.injectResultProject()
	args.injectContextGetProject()

	return &schema.Resource{
		ReadContext: dataSourceRustackProjectRead,
		Schema:      args,
	}
}

func dataSourceRustackProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()

	target, err := checkDatasourceNameOrId(d)
	if err != nil {
		return diag.Errorf("Error getting project: %s", err)
	}
	var targetProject *rustack.Project
	if target == "id" {
		targetProject, err = manager.GetProject(d.Get("id").(string))
		if err != nil {
			return diag.Errorf("Error getting project: %s", err)
		}
	} else {
		targetProject, err = GetProjectByName(d, manager)
		if err != nil {
			return diag.Errorf("Error getting project: %s", err)
		}
	}

	flattenedProject := map[string]interface{}{
		"id":   targetProject.ID,
		"name": targetProject.Name,
		// "project_id":   nil,
		// "project_name": nil,
	}

	if err := setResourceDataFromMap(d, flattenedProject); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(targetProject.ID)
	return nil
}
