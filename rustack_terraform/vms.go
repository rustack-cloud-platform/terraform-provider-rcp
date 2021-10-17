package rustack_terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func (args *Arguments) injectContextVmByName() {
	args.merge(Arguments{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "name of the Vm",
		},
	})
}

func (args *Arguments) injectCreateVm() {
	diskCreation := Defaults()
	diskCreation.injectCreateDisk()

	portCreation := Defaults()
	portCreation.injectCreatePort()

	args.merge(Arguments{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.All(
				validation.NoZeroValues,
				validation.StringLenBetween(2, 100),
			),
			Description: "name of the Vm",
		},
		"cpu": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 128),
			Description:  "the number of virtual cpus",
		},
		"ram": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 256),
			Description:  "memory of the Vm in gigabytes",
		},
		"template_id": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "id of the Template",
		},
		"user_data": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "script for cloud-init",
		},
		"disk": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			MaxItems:    20,
			Description: "list of Disks attached to the Vm",
			Elem: &schema.Resource{
				Schema: diskCreation,
			},
		},
		"port": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			MaxItems:    10,
			Description: "list of Ports attached to the Vm",
			Elem: &schema.Resource{
				Schema: portCreation,
			},
		},
		"floating_ip": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "enable floating ip for the Vm",
		},
	})
}

func (args *Arguments) injectResultVm() {
	args.merge(Arguments{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "id of the VDC",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "name of the Vm",
		},
		"cpu": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "the number of virtual cpus",
		},
		"ram": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "memory of the Vm in gigabytes",
		},
		"template_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "id of the Template",
		},
		"template_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "name of the Template",
		},
		"floating_ip": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "enable floating ip for the Vm",
		},
	})
}

func (args *Arguments) injectResultListVm() {
	s := Defaults()
	s.injectResultVm()

	args.merge(Arguments{
		"vms": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: s,
			},
		},
	})
}
