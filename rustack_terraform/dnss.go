package rustack_terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func (args *Arguments) injectContextGetDns() {
	args.merge(Arguments{
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "name of the Dns",
		},
		"id": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "id of the Dns",
		},
	})
}

func (args *Arguments) injectContextDnsById() {
	args.merge(Arguments{
		"dns_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "id of the Dns",
		},
	})
}

func (args *Arguments) injectCreateDns() {
	args.merge(Arguments{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.NoZeroValues,
				validation.StringLenBetween(1, 255),
			),
			Description: "name of the Dns",
		},
		"project_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "id of the Project",
		},
	})
}

func (args *Arguments) injectResultDns() {
	args.merge(Arguments{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "id of the Dns",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "name of the Dns",
		},
	})
}

func (args *Arguments) injectResultListDns() {
	s := Defaults()
	s.injectResultDns()

	args.merge(Arguments{
		"dnss": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: s,
			},
		},
	})
}
