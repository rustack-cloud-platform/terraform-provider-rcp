package rustack_terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)


func (args *Arguments) injectContextFirewallTemplateByName() {
	args.merge(Arguments{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "name of the Firewall Template",
		},
	})
}


func (args *Arguments) injectContextFirewallTemplateById() {
	args.merge(Arguments{
		"firewall_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "id of the Firewall Template",
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
	firewallTemplate := Defaults()
	firewallTemplate.injectResultFirewallTemplate()

	args.merge(Arguments{
		"firewall_templates": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: firewallTemplate,
			},
		},
	})
}

func (args *Arguments) injectCreateFirewallTemplate() {
	args.merge(Arguments{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "id of the firewall template",
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.All(
				validation.NoZeroValues,
				validation.StringLenBetween(1, 100),
			),
			Description: "name of the firewall template",
		},
	})
}
