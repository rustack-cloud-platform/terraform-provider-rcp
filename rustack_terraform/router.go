package rustack_terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func (args *Arguments) injectContextGetRouter() {
	args.merge(Arguments{
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "name of the Router",
		},
		"id": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "id of the Router",
		},
	})
}

func (args *Arguments) injectResultRouter() {

	args.merge(Arguments{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "id of the Router",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "name of the Router",
		},
	})
}

func (args *Arguments) injectResultListRouter() {
	Router := Defaults()
	Router.injectResultRouter()

	args.merge(Arguments{
		"routers": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: Router,
			},
		},
	})
}

func (args *Arguments) injectCreateRouter() {
	args.merge(Arguments{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Id of the Router",
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ValidateFunc: validation.All(
				validation.NoZeroValues,
				validation.StringLenBetween(1, 100),
			),
			Description: "Name of the Router",
		},
		"is_default": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Set if this is default router",
		},
		"floating": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Enable floating ip for the Vm",
		},
		"floating_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Floating id address.",
		},
		"ports": {
			Type:        schema.TypeSet,
			Optional:    true,
			Computed:    true,
			MinItems:    1,
			MaxItems:    10,
			Description: "List of Ports connected to the router",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"system": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Determinate if router is system.",
		},
		"tags": newTagNamesResourceSchema("tags of the router"),
	})
}
