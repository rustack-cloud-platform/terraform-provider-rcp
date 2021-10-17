package rustack_terraform

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRustackProject() *schema.Resource {
	args := Defaults()
	args.injectResultProject()
	args.injectContextProjectName()

	return &schema.Resource{
		ReadContext: dataSourceRustackProjectRead,
		Schema:      args,
	}
}

func dataSourceRustackProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetProject, err := GetProjectByName(d, manager)
	if err != nil {
		return diag.Errorf("Error getting project: %s", err)
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
