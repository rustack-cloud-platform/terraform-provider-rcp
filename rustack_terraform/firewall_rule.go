package rustack_terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (args *Arguments) injectCreateFirewallRule() {
	args.merge(Arguments{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "id of the firewall rule",
		},
		"direction": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "direction of the firewall rule (ingress, egress)",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "name of the firewall rule",
		},
		"destination_ip": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "destination ip address",
		},
		"port_range": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "max range of port",
		},
		"protocol": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "protocol tcp/upd/icmp",
		},
	})
}
