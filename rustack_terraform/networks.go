package rustack_terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func (args *Arguments) injectContextNetworkById() {
	args.merge(Arguments{
		"network_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "id of the Network",
		},
	})
}

func (args *Arguments) injectContextNetworkByName() {
	args.merge(Arguments{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "name of the Network",
		},
	})
}

func (args *Arguments) injectCreateNetwork() {
	args.merge(Arguments{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.All(
				validation.NoZeroValues,
				validation.StringLenBetween(2, 100),
			),
			Description: "name of the Network",
		},
		"subnets": {
			Type:     schema.TypeList,
			Required: true,
			MinItems: 1,
			MaxItems: 1, // Rustack doesn't support several subnets
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "id of the Subnet",
					},
					"cidr": {
						Type:        schema.TypeString,
						ForceNew:    true,
						Required:    true,
						Description: "cidr of the Subnet",
					},
					"gateway": {
						Type:        schema.TypeString,
						ForceNew:    true,
						Required:    true,
						Description: "gateway of the Subnet",
					},
					"start_ip": {
						Type:        schema.TypeString,
						ForceNew:    true,
						Required:    true,
						Description: "pool start ip of the Subnet",
					},
					"end_ip": {
						Type:        schema.TypeString,
						ForceNew:    true,
						Required:    true,
						Description: "pool end ip of the Subnet",
					},
					"dhcp": {
						Type:        schema.TypeBool,
						Required:    true,
						Description: "enable dhcp service of the Subnet",
					},
					"dns": {
						Type:        schema.TypeList,
						Required:    true,
						Elem:        &schema.Schema{Type: schema.TypeString},
						Description: "dns servers list",
					},
				},
			},
		},
	})
}

func (args *Arguments) injectResultNetwork() {
	args.merge(Arguments{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "id of the Network",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "name of the Network",
		},
		"subnets": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "list of subnets",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "id of the Subnet",
					},
					"cidr": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "cidr of the Subnet",
					},
					"gateway": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "gateway of the Subnet",
					},
					"start_ip": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "pool start ip of the Subnet",
					},
					"end_ip": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "pool end ip of the Subnet",
					},
					"dhcp": {
						Type:        schema.TypeBool,
						Computed:    true,
						Description: "enable dhcp service of the Subnet",
					},
					"dns": {
						Type:        schema.TypeList,
						Computed:    true,
						Elem:        &schema.Schema{Type: schema.TypeString},
						Description: "dns servers list",
					},
				},
			},
		},
	})
}

func (args *Arguments) injectResultListNetwork() {
	s := Defaults()
	s.injectResultNetwork()

	args.merge(Arguments{
		"networks": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: s,
			},
		},
	})
}
