package rustack_terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func (args *Arguments) injectContextDiskByName() {
	args.merge(Arguments{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "name of the disk",
		},
	})
}

func (args *Arguments) injectCreateDisk() {
	args.merge(Arguments{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "id of the Disk",
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.All(
				validation.NoZeroValues,
				validation.StringLenBetween(2, 100),
			),
			Description: "name of the Disk",
		},
		"size": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 10000),
			Description:  "the size of the Disk in gigabytes",
		},
		"storage_profile_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "the id of the StorageProfile",
		},
	})
}

func (args *Arguments) injectResultDisk() {
	args.merge(Arguments{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "id of the Disk",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "name of the Disk",
		},
		"size": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "the size of the Disk in gigabytes",
		},
		"storage_profile_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "the id of the StorageProfile",
		},
		"storage_profile_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "the name of the StorageProfile",
		},
	})
}

func (args *Arguments) injectResultListDisk() {
	s := Defaults()
	s.injectResultDisk()

	args.merge(Arguments{
		"disks": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: s,
			},
		},
	})
}
