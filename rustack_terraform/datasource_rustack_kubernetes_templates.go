package rustack_terraform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/hashstructure/v2"
)

func dataSourceRustackKubernetesTemplates() *schema.Resource {
	args := Defaults()
	args.injectContextVdcById()
	args.injectResultListKubernetesTemplate()

	return &schema.Resource{
		ReadContext: dataSourceRustackKubernetesTemplateReadRead,
		Schema:      args,
	}
}

func dataSourceRustackKubernetesTemplateReadRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting vdc: %s", err)
	}

	allKubernetesTemplateRead, err := targetVdc.GetKubernetesTemplates()
	if err != nil {
		return diag.Errorf("Error retrieving KubernetesTemplateRead: %s", err)
	}

	flattenedRecords := make([]map[string]interface{}, len(allKubernetesTemplateRead))
	for i, KubernetesTemplateRead := range allKubernetesTemplateRead {
		flattenedRecords[i] = map[string]interface{}{
			"id":           KubernetesTemplateRead.ID,
			"name":         KubernetesTemplateRead.Name,
			"min_node_cpu": KubernetesTemplateRead.MinNodeCpu,
			"min_node_ram": KubernetesTemplateRead.MinNodeRam,
			"min_node_hdd": KubernetesTemplateRead.MinNodeHdd,
		}
	}

	hash, err := hashstructure.Hash(allKubernetesTemplateRead, hashstructure.FormatV2, nil)
	if err != nil {
		diag.Errorf("unable to set `kubernetes_templates` attribute: %s", err)
	}

	d.SetId(fmt.Sprintf("kubernetes_templates/%d", hash))

	if err := d.Set("kubernetes_templates", flattenedRecords); err != nil {
		return diag.Errorf("unable to set `kubernetes_templates` attribute: %s", err)
	}

	return nil
}
