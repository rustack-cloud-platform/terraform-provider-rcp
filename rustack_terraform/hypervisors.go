package rustack_terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (args *Arguments) injectContextHypervisorByName() {
	args.merge(Arguments{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "name of the Hypervisor",
		},
	})
}

func (args *Arguments) injectContextHypervisorById() {
	args.merge(Arguments{
		"hypervisor_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "id of the Hypervisor",
		},
	})
}

func (args *Arguments) injectResultHypervisor() {
	args.merge(Arguments{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "id of the Hypervisor",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "name of the Hypervisor",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "type of the Hypervisor",
		},
	})
}

func (args *Arguments) injectResultListHypervisor() {
	s := Defaults()
	s.injectResultHypervisor()

	args.merge(Arguments{
		"hypervisors": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: s,
			},
		},
	})
}
