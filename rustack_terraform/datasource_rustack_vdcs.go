package rustack_terraform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/hashstructure/v2"
	"github.com/pilat/rustack-go/rustack"
)

func dataSourceRustackVdcs() *schema.Resource {
	args := Defaults()
	args.injectContextProjectById()
	args.injectResultListVdc()

	return &schema.Resource{
		ReadContext: dataSourceRustackVdcsRead,
		Schema:      args,
	}
}

func dataSourceRustackVdcsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetProject, err := GetProjectById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting project: %s", err)
	}

	allVdcs, err := manager.GetVdcs(rustack.Arguments{"project": targetProject.ID})
	if err != nil {
		return diag.Errorf("Error retrieving vdcs: %s", err)
	}

	flattenedRecords := make([]map[string]interface{}, len(allVdcs))
	for i, vdc := range allVdcs {
		flattenedRecords[i] = map[string]interface{}{
			"id":              vdc.ID,
			"name":            vdc.Name,
			"hypervisor":      vdc.Hypervisor.Name,
			"hypervisor_type": vdc.Hypervisor.Type,
			// "project":         vdc.Project.Name,
		}
	}

	hash, err := hashstructure.Hash(allVdcs, hashstructure.FormatV2, nil)
	if err != nil {
		diag.Errorf("unable to set `vdcs` attribute: %s", err)
	}

	d.SetId(fmt.Sprintf("vdcs/%d", hash))

	if err := d.Set("vdcs", flattenedRecords); err != nil {
		return diag.Errorf("unable to set `vdcs` attribute: %s", err)
	}

	return nil
}
