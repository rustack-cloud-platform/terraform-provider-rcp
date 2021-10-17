package rustack_terraform

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pilat/rustack-go/rustack"
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
	allClients, err := manager.GetClients()
	if err != nil {
		return diag.Errorf("Error getting list of clients: %s", err)
	}
	if len(allClients) == 0 {
		return diag.Errorf("There are no available clients")
	}
	if len(allClients) > 1 {
		return diag.Errorf("More than one client available for you") // TODO: use provider's variable
	}

	client := allClients[0]

	project := rustack.NewProject(
		d.Get("name").(string),
	)

	log.Printf("[DEBUG] Project create request: %#v", project)
	err = client.CreateProject(&project)
	if err != nil {
		return diag.Errorf("Error creating project: %s", err)
	}

	d.SetId(project.ID)
	log.Printf("[INFO] Project created, ID: %s", d.Id())

	return resourceRustackProjectRead(ctx, d, meta)
}

func resourceRustackProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	project, err := manager.GetProject(d.Id())
	if err != nil {
		return diag.Errorf("Error getting project: %s", err)
	}

	d.SetId(project.ID)
	d.Set("name", project.Name)

	return nil
}

func resourceRustackProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()

	project, err := manager.GetProject(d.Id())
	if err != nil {
		return diag.Errorf("Error getting project: %s", err)
	}

	err = project.Rename(d.Get("name").(string))
	if err != nil {
		return diag.Errorf("Error rename project: %s", err)
	}

	log.Printf("[INFO] Updated Project, ID: %#v", project)

	return resourceRustackProjectRead(ctx, d, meta)
}

func resourceRustackProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()

	projectId := d.Id()

	project, err := manager.GetProject(projectId)
	if err != nil {
		return diag.Errorf("Error getting project: %s", err)
	}

	err = project.Delete()
	if err != nil {
		return diag.Errorf("Error deleting project: %s", err)
	}

	d.SetId("")
	log.Printf("[INFO] Project deleted, ID: %s", projectId)

	return nil
}
