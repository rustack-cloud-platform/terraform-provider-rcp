package rustack_terraform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRustackPaasTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRustackPaasTemplateRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "id of Paas Template",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "id of Project",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "name of Paas Template",
			},
		},
	}
}

func dataSourceRustackPaasTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	template, err := manager.GetPaasTemplate(d.Get("id").(int), d.Get("project_id").(string))
	if err != nil {
		return diag.Errorf("Error getting paas template: %s", err)
	}
	flatten := map[string]interface{}{
		"id":   template.ID,
		"name": template.Name,
	}
	if err := setResourceDataFromMap(d, flatten); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprint(template.ID))
	return nil
}
