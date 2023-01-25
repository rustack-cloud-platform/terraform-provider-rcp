package rustack_terraform

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pilat/rustack-go/rustack"
)

func resourceRustackS3Storage() *schema.Resource {
	args := Defaults()
	args.injectContextProjectById()
	args.injectCreateS3Storage()

	return &schema.Resource{
		CreateContext: resourceRustackS3StorageCreate,
		ReadContext:   resourceRustackS3StorageRead,
		UpdateContext: resourceRustackS3StorageUpdate,
		DeleteContext: resourceRustackS3StorageDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: args,
	}
}

func resourceRustackS3StorageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	project, err := GetProjectById(d, manager)
	if err != nil {
		return diag.Errorf("project_id: Error getting Project: %s", err)
	}
	name := d.Get("name").(string)
	newS3Storage := rustack.NewS3Storage(name)

	err = project.CreateS3Storage(&newS3Storage)
	if err != nil {
		return diag.Errorf("Error creating S3Storage: %s", err)
	}

	newS3Storage.WaitLock()
	d.SetId(newS3Storage.ID)
	log.Printf("[INFO] S3Storage created, ID: %s", d.Id())

	return resourceRustackS3StorageRead(ctx, d, meta)
}

func resourceRustackS3StorageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()

	s3, err := manager.GetS3Storage(d.Id())
	if d.HasChange("name") {
		s3.Name = d.Get("name").(string)
	}

	err = s3.Update()
	if err != nil {
		return diag.Errorf("Error updating S3Storage: %s", err)
	}
	s3.WaitLock()
	log.Printf("[INFO] S3Storage updated, ID: %s", d.Id())

	return resourceRustackS3StorageRead(ctx, d, meta)
}

func resourceRustackS3StorageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	S3Storage, err := manager.GetS3Storage(d.Id())
	if err != nil {
		return diag.Errorf("id: Error getting S3Storage: %s", err)
	}

	d.SetId(S3Storage.ID)
	d.Set("name", S3Storage.Name)
	d.Set("project", S3Storage.Project.ID)
	d.Set("client_endpoint", S3Storage.ClientEndpoint)
	d.Set("secret_key", S3Storage.SecretKey)
	d.Set("access_key", S3Storage.AccessKey)

	return nil
}

func resourceRustackS3StorageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	s3_id := d.Id()
	s3, err := manager.GetS3Storage(d.Id())
	if err != nil {
		return diag.Errorf("id: Error getting S3Storage: %s", err)
	}

	err = s3.Delete()
	if err != nil {
		return diag.Errorf("Error deleting S3Storage: %s", err)
	}
	s3.WaitLock()

	d.SetId("")
	log.Printf("[INFO] S3Storage deleted, ID: %s", s3_id)

	return nil
}
