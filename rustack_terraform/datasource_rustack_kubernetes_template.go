package rustack_terraform

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pilat/rustack-go/rustack"
)

func dataSourceRustackKubernetesTemplate() *schema.Resource {
	args := Defaults()
	args.injectResultKubernetesTemplate()
	args.injectContextVdcById()
	args.injectContextGetKubernetesTemplate() // override name

	return &schema.Resource{
		ReadContext: dataSourceRustackKubernetesTemplateRead,
		Schema:      args,
	}
}

func dataSourceRustackKubernetesTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting vdc: %s", err)
	}

	target, err := checkDatasourceNameOrId(d)
	if err != nil {
		return diag.Errorf("Error getting KubernetesTemplate: %s", err)
	}
	var targetKubernetesTemplate *rustack.KubernetesTemplate
	if target == "id" {
		targetKubernetesTemplate, err = manager.GetKubernetesTemplate(d.Get("id").(string))
		if err != nil {
			return diag.Errorf("Error getting KubernetesTemplate: %s", err)
		}
	} else {
		targetKubernetesTemplate, err = GetKubernetesTemplateByName(d, manager, targetVdc)
		if err != nil {
			return diag.Errorf("Error getting KubernetesTemplate: %s", err)
		}
	}

	flatten := map[string]interface{}{
		"id":           targetKubernetesTemplate.ID,
		"name":         targetKubernetesTemplate.Name,
		"min_node_cpu": targetKubernetesTemplate.MinNodeCpu,
		"min_node_ram": targetKubernetesTemplate.MinNodeRam,
		"min_node_hdd": targetKubernetesTemplate.MinNodeHdd,
	}

	if err := setResourceDataFromMap(d, flatten); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(targetKubernetesTemplate.ID)
	return nil
}
