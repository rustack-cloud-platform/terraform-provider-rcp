package rustack_terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (args *Arguments) injectContextGetStorageProfile() {
	args.merge(Arguments{
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Name of the storage profile",
		},
		"id": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "id of the storage profile",
		},
	})
}

func (args *Arguments) injectContextStorageProfileById() {
	args.merge(Arguments{
		"storage_profile_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Id of the storage profile",
		},
	})
}

func (args *Arguments) injectResultStorageProfile() {
	args.merge(Arguments{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The identifier for the storage profile",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Name of the storage profile",
		},
	})
}

func (args *Arguments) injectResultListStorageProfile() {
	s := Defaults()
	s.injectResultStorageProfile()

	args.merge(Arguments{
		"storage_profiles": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: s,
			},
		},
	})
}
