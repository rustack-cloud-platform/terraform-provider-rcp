package rustack_terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (args *Arguments) injectCreateDnsRecord() {
	args.merge(Arguments{
		"data": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "data of dns record",
		},
		"flag": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "flag of dns record",
		},
		"host": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "host of dns record",
		},
		"port": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "port of dns record",
		},
		"priority": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "priority of dns record",
		},
		"tag": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "tag of dns record",
		},
		"ttl": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     "86400",
			Description: "ttl of dns record",
		},
		"type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "type of dns record",
		},
		"weight": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "weight of dns record",
		},
	})
}
