package rustack_terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func (args *Arguments) injectContextLbaasPoolByName() {
	args.merge(Arguments{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "name of the Lbaas",
		},
	})
}

func (args *Arguments) injectCreateLbaasPool() {
	poolMembers := Defaults()
	poolMembers.injectLbaasPoolMembers()
	

	args.merge(Arguments{
		"connlimit": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default: "65536",
			Description: "id of the Template",
		},
		"cookie_name": {
			Type:     schema.TypeString,
			Optional: true,
			Default: "",
			Description: "name of the Lbaas",
		},
		"method": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "enable floating ip for the Lbaas",
		},
		"port": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "floating ip for the Lbaas. May be omitted",
		},
		"protocol": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "floating ip for the Lbaas. May be omitted",
		},
		"session_persistence": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "floating ip for the Lbaas. May be omitted",
		},
		"member": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Resource{
				Schema: poolMembers,
			},
			Description: "Lbaas members.",
		},
	})
}

func (args *Arguments) injectResultLbaasPool() {
	args.merge(Arguments{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "id of the Lbaas",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "name of the Lbaas",
		},
		"floating": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "enable floating ip for the Lbaas",
		},
		"floating_ip": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "floating_ip of the Lbaas. May be omitted",
		},
	})
}

func (args *Arguments) injectLbaasPoolMembers() {
	args.injectContextVmById()

	args.merge(Arguments{
		"port": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "list of firewall templates ids of the Port",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"weight": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "ip_address of the Port",
			ValidateFunc: validation.All(
				validation.IntBetween(0, 256),
			),
		},
	},
	)
}
