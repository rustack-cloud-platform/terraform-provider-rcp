package rustack_terraform

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRustackHypervisor() *schema.Resource {
	args := Defaults()
	args.injectResultHypervisor()
	args.injectContextProjectById()
	args.injectContextHypervisorByName()

	return &schema.Resource{
		ReadContext: dataSourceRustackHypervisorRead,
		Schema:      args,
	}
}

func dataSourceRustackHypervisorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetProject, err := GetProjectById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting project: %s", err)
	}

	targetHypervisor, err := GetHypervisorByName(d, manager, targetProject)
	if err != nil {
		return diag.Errorf("Error getting hypervisor: %s", err)
	}

	flatten := map[string]interface{}{
		"id":         targetHypervisor.ID,
		"name":       targetHypervisor.Name,
		"type":       targetHypervisor.Type,
		"project_id": targetProject.ID,
	}

	if err := setResourceDataFromMap(d, flatten); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(targetHypervisor.ID)
	return nil
}
