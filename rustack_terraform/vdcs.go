package rustack_terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func (args *Arguments) injectCreateVdc() {
	args.merge(Arguments{
		"project_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "id of the Project",
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.All(
				validation.NoZeroValues,
				validation.StringLenBetween(1, 100),
			),
			Description: "name of the VDC",
		},
		"default_network_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "id of the Network",
		},
		"default_network_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "name of the Network",
		},
		"default_network_mtu": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "mtu of the Network",
		},
		"default_network_subnets": {
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
		"tags": newTagNamesResourceSchema("tags of the VDC"),
	})
}

func (args *Arguments) injectContextVdcById() {
	args.merge(Arguments{
		"vdc_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "id of the VDC",
		},
	})
}

func (args *Arguments) injectContextVdcByIdForData() {
	args.merge(Arguments{
		"vdc_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "id of the VDC",
		},
	})
}

func (args *Arguments) injectContextGetVdc() {
	args.merge(Arguments{
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "name of the vdc",
		},
		"id": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "id of the VDC",
		},
	})
}

func (args *Arguments) injectResultVdc() {
	args.merge(Arguments{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "id of the VDC",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "name of the VDC",
		},
		"hypervisor": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "name of the Hypervisor",
		},
		"hypervisor_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "type of the Hypervisor",
		},
	})
}

func (args *Arguments) injectResultListVdc() {
	s := Defaults()
	s.injectResultVdc()

	args.merge(Arguments{
		"vdcs": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: s,
			},
		},
	})
}
