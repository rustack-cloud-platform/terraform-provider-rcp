package rustack_terraform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/hashstructure/v2"
)

func dataSourceRustackProjects() *schema.Resource {
	args := Defaults()
	args.injectResultListProject()

	return &schema.Resource{
		ReadContext: dataSourceRustackProjectsRead,
		Schema:      args,
	}
}

func dataSourceRustackProjectsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()

	allProjects, err := manager.GetProjects()
	if err != nil {
		return diag.Errorf("Error getting projects: %s", err)
	}

	flattenedRecords := make([]map[string]interface{}, len(allProjects))
	for i, project := range allProjects {
		flattenedRecords[i] = map[string]interface{}{
			"id":   project.ID,
			"name": project.Name,
		}
	}

	hash, err := hashstructure.Hash(allProjects, hashstructure.FormatV2, nil)
	if err != nil {
		diag.Errorf("unable to set `projects` attribute: %s", err)
	}

	d.SetId(fmt.Sprintf("projects/%d", hash))

	if err := d.Set("projects", flattenedRecords); err != nil {
		return diag.Errorf("unable to set `projects` attribute: %s", err)
	}

	return nil
}
