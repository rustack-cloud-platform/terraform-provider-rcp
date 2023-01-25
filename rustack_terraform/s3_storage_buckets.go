package rustack_terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func (args *Arguments) injectContextS3StorageBucketByName() {
	args.merge(Arguments{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "name of the S3StorageBucket",
		},
	})
}

func (args *Arguments) injectContextS3StorageBucketById() {
	args.merge(Arguments{
		"s3_bucket_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "id of the S3StorageBucket",
		},
	})
}

func (args *Arguments) injectCreateS3StorageBucket() {
	args.merge(Arguments{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "id of the S3Storage",
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.All(
				validation.NoZeroValues,
				validation.StringLenBetween(1, 255),
			),
			Description: "name of the S3Storage",
		},
		"external_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "url for connecting to s3",
		},
	})
}

func (args *Arguments) injectResultS3StorageBucket() {
	args.merge(Arguments{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "id of the S3StorageBucket",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "name of the S3StorageBucket",
		},
		"external_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "external_name of the S3StorageBucket",
		},
	})
}

func (args *Arguments) injectResultListS3StorageBucket() {
	s := Defaults()
	s.injectResultS3StorageBucket()

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
