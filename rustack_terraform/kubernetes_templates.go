package rustack_terraform

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func (args *Arguments) injectContextKubernetesTemplateById() {
	args.merge(Arguments{
		"template_id": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Required:    true,
			Description: "id of the Kubernetes Template",
		},
	})
}

func (args *Arguments) injectContextGetKubernetesTemplate() {
	args.merge(Arguments{
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "name of the Kubernetes Template",
		},
		"id": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "id of the Kubernetes Template",
		},
	})
}

func (args *Arguments) injectResultKubernetesTemplate() {
	args.merge(Arguments{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "id of the Kubernetes Template",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "name of the Kubernetes Template",
		},
		"min_node_cpu": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "minimum cpu required for node by the Kubernetes Template",
		},
		"min_node_ram": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "minimum ram in GB required for node by the Kubernetes Template",
		},
		"min_node_hdd": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "minimum disk size in GB required for node by the Kubernetes Template",
		},
	})
}

func (args *Arguments) injectResultListKubernetesTemplate() {
	s := Defaults()
	s.injectResultKubernetesTemplate()

	args.merge(Arguments{
		"kubernetes_templates": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: s,
			},
		},
	})
}
