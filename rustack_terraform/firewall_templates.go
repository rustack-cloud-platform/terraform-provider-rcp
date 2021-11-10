package rustack_terraform

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func (args *Arguments) injectContextFirewallTemplateById() {
	args.merge(Arguments{
		"firewall_template_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "id of the Firewall Template",
		},
	})
}

func (args *Arguments) injectContextFirewallTemplateByName() {
	args.merge(Arguments{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "name of the Firewall Template",
		},
	})
}

func (args *Arguments) injectResultFirewallTemplate() {
	args.merge(Arguments{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "id of the Firewall Template",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "name of the Firewall Template",
		},
	})
}

func (args *Arguments) injectResultListFirewallTemplate() {
	s := Defaults()
	s.injectResultFirewallTemplate()

	args.merge(Arguments{
		"firewall_templates": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: s,
			},
		},
	})
}
