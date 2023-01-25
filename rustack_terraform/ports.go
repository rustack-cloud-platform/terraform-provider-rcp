package rustack_terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (args *Arguments) injectContextPortById() {
	args.merge(Arguments{
		"id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "id of the port",
		},
	})
}

func (args *Arguments) injectContextPortByIp() {
	args.merge(Arguments{
		"ip_address": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "ip_address of the Port",
		},
	})
}

func (args *Arguments) injectCreatePort() {
	args.injectContextNetworkById()

	args.merge(Arguments{
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
			Optional:    true,
			Description: "list of firewall templates ids of the Port",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	})
}

func (args *Arguments) injectResultPort() {
	args.merge(Arguments{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "id of the Port",
		},
		"network": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "id of the Network",
		},
		"ip_address": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "ip_address of the Port",
		},
	})
}

func (args *Arguments) injectResultListPort() {
	portSchema := Defaults()
	portSchema.injectResultPort()

	args.merge(Arguments{
		"ports": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: portSchema,
			},
		},
	})
}
