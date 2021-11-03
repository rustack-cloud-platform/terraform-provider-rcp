package rustack_terraform

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func (args *Arguments) injectContextTemplateById() {
	args.merge(Arguments{
		"template_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "id of the Template",
		},
	})
}

func (args *Arguments) injectContextTemplateByName() {
	args.merge(Arguments{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "name of the Template",
		},
	})
}

func (args *Arguments) injectResultTemplate() {
	args.merge(Arguments{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "id of the Template",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "name of the Template",
		},
		"min_cpu": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "minimum cpu required by the Template",
		},
		"min_ram": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "minimum ram in GB required by the Template",
		},
		"min_disk": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "minimum disk size in GB required by the Template",
		},
	})
}

func (args *Arguments) injectResultListTemplate() {
	s := Defaults()
	s.injectResultTemplate()

	args.merge(Arguments{
		"templates": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: s,
			},
		},
	})
}
