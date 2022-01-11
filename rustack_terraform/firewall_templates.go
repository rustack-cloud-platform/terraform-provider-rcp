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
			Description: "name of the Firewall Template",
		},
	})
}

func (args *Arguments) injectResultFirewallTemplate() {
	ingressRule := Defaults()
	ingressRule.injectCreateFirewallRule()
	egressRule := Defaults()
	egressRule.injectCreateFirewallRule()

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
		"ingress_rule": {
			Type:     schema.TypeList,
			Optional: true,
			// TODO: setup limits
			MinItems:    0,
			MaxItems:    20,
			Description: "list of ingress rules",
			Elem: &schema.Resource{
				Schema: ingressRule,
			},
		},
		"egress_rule": {
			Type:     schema.TypeList,
			Optional: true,
			// TODO: setup limits
			MinItems:    0,
			MaxItems:    20,
			Description: "list of egress rules",
			Elem: &schema.Resource{
				Schema: egressRule,
			},
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
	ingressRule := Defaults()
	ingressRule.injectCreateFirewallRule()
	egressRule := Defaults()
	egressRule.injectCreateFirewallRule()

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
				validation.StringLenBetween(2, 100),
			),
			Description: "name of the firewall template",
		},
		"ingress_rule": {
			Type:     schema.TypeList,
			Optional: true,
			// TODO: setup limits
			MinItems:    0,
			MaxItems:    20,
			Description: "list of ingress rules",
			Elem: &schema.Resource{
				Schema: ingressRule,
			},
		},
		"egress_rule": {
			Type:     schema.TypeList,
			Optional: true,
			// TODO: setup limits
			MinItems:    0,
			MaxItems:    20,
			Description: "list of egress rules",
			Elem: &schema.Resource{
				Schema: egressRule,
			},
		},
	})
}
