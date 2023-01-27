package rustack_terraform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/hashstructure/v2"
)

func dataSourceRustackS3Storages() *schema.Resource {
	args := Defaults()
	args.injectContextProjectById()
	args.injectResultListS3Storage()

	return &schema.Resource{
		ReadContext: dataSourceRustackS3Read,
		Schema:      args,
	}
}

func dataSourceRustackS3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	project, err := GetProjectById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting project: %s", err)
	}

	s3_storages, err := project.GetS3Storages()
	if err != nil {
		return diag.Errorf("Error retrieving storages: %s", err)
	}

	flattenedRecords := make([]map[string]interface{}, len(s3_storages))
	for i, s3 := range s3_storages {
		flattenedRecords[i] = map[string]interface{}{
			"id":              s3.ID,
			"name":            s3.Name,
			"client_endpoint": s3.ClientEndpoint,
			"access_key":      s3.AccessKey,
			"secret_key":      s3.SecretKey,
		}

	}

	hash, err := hashstructure.Hash(s3_storages, hashstructure.FormatV2, nil)
	if err != nil {
		diag.Errorf("unable to set `s3_storages` attribute: %s", err)
	}

	d.SetId(fmt.Sprintf("s3_storages/%d", hash))

	if err := d.Set("s3_storages", flattenedRecords); err != nil {
		return diag.Errorf("unable to set `s3_storages` attribute: %s", err)
	}

	return nil
}
