package rustack_terraform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/hashstructure/v2"
)

func dataSourceRustackHypervisors() *schema.Resource {
	args := Defaults()
	args.injectContextProjectById()
	args.injectResultListHypervisor()

	return &schema.Resource{
		ReadContext: dataSourceRustackHypervisorsRead,
		Schema:      args,
	}
}

func dataSourceRustackHypervisorsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetProject, err := GetProjectById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting project: %s", err)
	}

	hypervisors, err := targetProject.GetAvailableHypervisors()
	if err != nil {
		return diag.Errorf("Error getting available hypervisors")
	}

	flattenedHypervisors := make([]map[string]interface{}, len(hypervisors))
	for i, hypervisor := range hypervisors {
		flattenedHypervisors[i] = map[string]interface{}{
			"id":   hypervisor.ID,
			"name": hypervisor.Name,
			"type": hypervisor.Type,
		}
	}

	hash, err := hashstructure.Hash(hypervisors, hashstructure.FormatV2, nil)
	if err != nil {
		diag.Errorf("unable to set `hypervisors` attribute: %s", err)
	}

	d.SetId(fmt.Sprintf("hypervisors/%d", hash))
	d.Set("project_id", nil)
	// d.Set("project_name", nil)

	if err := d.Set("hypervisors", flattenedHypervisors); err != nil {
		return diag.Errorf("unable to set `hypervisors` attribute: %s", err)
	}

	return nil
}
