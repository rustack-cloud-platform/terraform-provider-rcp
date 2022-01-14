package rustack_terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func (args *Arguments) injectContextRouterByName() {
	args.merge(Arguments{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "name of the Router",
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
				validation.StringLenBetween(2, 100),
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
		"networks": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			// TODO: setup limits
			// MinItems:    1,
			// MaxItems:    20,
			Description: "list of networks",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"system": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Determinate if router is system.",
		},
	})
}
