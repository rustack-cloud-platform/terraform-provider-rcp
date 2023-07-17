package rustack_terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func (args *Arguments) injectContextGetS3Storage() {
	args.merge(Arguments{
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "name of the S3Storage",
		},
		"id": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "id of the S3Storage",
		},
	})
}

func (args *Arguments) injectContextS3StorageById() {
	args.merge(Arguments{
		"s3_storage_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "id of the S3Storage",
		},
	})
}

func (args *Arguments) injectCreateS3Storage() {
	args.merge(Arguments{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.All(
				validation.NoZeroValues,
				validation.StringLenBetween(1, 255),
			),
			Description: "name of the S3Storage",
		},
		"backend": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "backend for s3",
		},
		"client_endpoint": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "url for connecting to s3",
		},
		"access_key": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "access_key for access to s3",
		},
		"secret_key": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "secret_key for access to s3",
		},
	})
}

func (args *Arguments) injectResultS3Storage() {
	args.merge(Arguments{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "id of the S3Storage",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "name of the S3Storage",
		},
		"backend": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "backend for s3",
		},
		"client_endpoint": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "url for connecting to s3",
		},
		"access_key": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "access_key for access to s3",
		},
		"secret_key": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "secret_key for access to s3",
		},
	})
}

func (args *Arguments) injectResultListS3Storage() {
	s := Defaults()
	s.injectResultS3Storage()

	args.merge(Arguments{
		"s3_storages": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: s,
			},
		},
	})
}
