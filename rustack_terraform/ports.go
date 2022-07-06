package rustack_terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (args *Arguments) injectCreatePort() {
	args.injectContextNetworkById()

	args.merge(Arguments{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "id of the Port",
		},
		"ip_address": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "ip_address of the Port",
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				return new == ""
			},
		},
		"firewall_templates": {
			Type:        schema.TypeSet,
			Required:    true,
			Description: "list of firewall rule ids of the Port",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	})
}
