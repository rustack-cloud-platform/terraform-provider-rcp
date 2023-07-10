package rustack_terraform

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func (args *Arguments) injectContextPlatformById() {
	args.merge(Arguments{
		"id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "id of the Platform",
		},
	})
}

func (args *Arguments) injectContextGetPlatform() {
	args.merge(Arguments{
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "name of the Platform",
		},
		"id": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "id of the Platform",
		},
	})
}

func (args *Arguments) injectResultPlatform() {
	args.merge(Arguments{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "id of the Platform",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "name of the Platform",
		},
	})
}

func (args *Arguments) injectResultListPlatforms() {
	s := Defaults()
	s.injectResultPlatform()

	args.merge(Arguments{
		"platforms": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: s,
			},
		},
	})
}
