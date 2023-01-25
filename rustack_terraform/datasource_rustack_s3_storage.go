package rustack_terraform

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRustackS3Storage() *schema.Resource {
	args := Defaults()
	args.injectContextProjectById()
	args.injectResultS3Storage()
	args.injectContextS3StorageByName() // override name

	return &schema.Resource{
		ReadContext: dataSourceRustackS3StorageRead,
		Schema:      args,
	}
}

func dataSourceRustackS3StorageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()

	s3_storage, err := GetS3ByName(d, manager)
	if err != nil {
		return diag.Errorf("Error getting s3: %s", err)
	}

	flatten := map[string]interface{}{
		"id":              s3_storage.ID,
		"name":            s3_storage.Name,
		"client_endpoint": s3_storage.ClientEndpoint,
		"access_key":      s3_storage.AccessKey,
		"secret_key":      s3_storage.SecretKey,
	}

	if err := setResourceDataFromMap(d, flatten); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(s3_storage.ID)
	return nil
}
