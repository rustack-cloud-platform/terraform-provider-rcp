package rustack_terraform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/hashstructure/v2"
)

func dataSourceRustackTemplates() *schema.Resource {
	args := Defaults()
	args.injectContextVdcById()
	args.injectResultListTemplate()

	return &schema.Resource{
		ReadContext: dataSourceRustackTemplatesRead,
		Schema:      args,
	}
}

func dataSourceRustackTemplatesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting vdc: %s", err)
	}

	allTemplates, err := targetVdc.GetTemplates()
	if err != nil {
		return diag.Errorf("Error retrieving templates: %s", err)
	}

	flattenedRecords := make([]map[string]interface{}, len(allTemplates))
	for i, templates := range allTemplates {
		flattenedRecords[i] = map[string]interface{}{
			"id":       templates.ID,
			"name":     templates.Name,
			"min_cpu":  templates.MinCpu,
			"min_ram":  templates.MinRam,
			"min_disk": templates.MinHdd,
		}
	}

	hash, err := hashstructure.Hash(allTemplates, hashstructure.FormatV2, nil)
	if err != nil {
		diag.Errorf("unable to set `templates` attribute: %s", err)
	}

	d.SetId(fmt.Sprintf("templates/%d", hash))

	if err := d.Set("templates", flattenedRecords); err != nil {
		return diag.Errorf("unable to set `templates` attribute: %s", err)
	}

	return nil
}
