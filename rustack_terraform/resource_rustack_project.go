package rustack_terraform

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rustack-cloud-platform/rcp-go/rustack"
)

func resourceRustackProject() *schema.Resource {
	args := Defaults()
	args.injectCreateProject()

	return &schema.Resource{
		CreateContext: resourceRustackProjectCreate,
		ReadContext:   resourceRustackProjectRead,
		UpdateContext: resourceRustackProjectUpdate,
		DeleteContext: resourceRustackProjectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: args,
	}
}

func resourceRustackProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	client_id := manager.ClientID
	var client *rustack.Client
	var err error

	if client_id != "" {
		client, err = manager.GetClient(client_id)
		if err != nil {
			return diag.Errorf("Error getting client: %s", err)
		}
	} else {
		allClients, err := manager.GetClients()
		if err != nil {
			return diag.Errorf("Error there are no clients available for management: %s", err)
		}
		if len(allClients) == 0 {
			return diag.Errorf("There are no available clients")
		}
		if len(allClients) > 1 {
			return diag.Errorf("More than one client available for you")
		}

		client = allClients[0]
	}

	project := rustack.NewProject(
		d.Get("name").(string),
	)
	project.Tags = unmarshalTagNames(d.Get("tags"))
	log.Printf("[DEBUG] Project create request: %#v", project)
	err = client.CreateProject(&project)
	if err != nil {
		return diag.Errorf("id: Error creating project: %s", err)
	}
	project.WaitLock()

	d.SetId(project.ID)
	log.Printf("[INFO] Project created, ID: %s", d.Id())

	return resourceRustackProjectRead(ctx, d, meta)
}

func resourceRustackProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	project, err := manager.GetProject(d.Id())
	if err != nil {
		if err.(*rustack.RustackApiError).Code() == 404 {
			d.SetId("")
			return nil
		} else {
			return diag.Errorf("id: Error getting project: %s", err)
		}
	}

	d.SetId(project.ID)
	d.Set("name", project.Name)
	d.Set("tags", marshalTagNames(project.Tags))

	return nil
}

func resourceRustackProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()

	project, err := manager.GetProject(d.Id())
	if err != nil {
		return diag.Errorf("id: Error getting project: %s", err)
	}

	if d.HasChange("name") {
		project.Name = d.Get("name").(string)
	}
	if d.HasChange("tags") {
		project.Tags = unmarshalTagNames(d.Get("tags"))
	}
	err = project.Update()
	if err != nil {
		return diag.Errorf("name: Error rename project: %s", err)
	}
	project.WaitLock()

	log.Printf("[INFO] Updated Project, ID: %#v", project)

	return resourceRustackProjectRead(ctx, d, meta)
}

func resourceRustackProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()

	projectId := d.Id()

	project, err := manager.GetProject(projectId)
	if err != nil {
		return diag.Errorf("id: Error getting project: %s", err)
	}

	err = project.Delete()
	if err != nil {
		return diag.Errorf("Error deleting project: %s", err)
	}
	project.WaitLock()

	d.SetId("")
	log.Printf("[INFO] Project deleted, ID: %s", projectId)

	return nil
}
