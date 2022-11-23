package rustack_terraform

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pilat/rustack-go/rustack"
)

func resourceRustackVdc() *schema.Resource {
	args := Defaults()
	args.injectCreateVdc()

	return &schema.Resource{
		CreateContext: resourceRustackVdcCreate,
		ReadContext:   resourceRustackVdcRead,
		UpdateContext: resourceRustackVdcUpdate,
		DeleteContext: resourceRustackVdcDelete,
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

func resourceRustackVdcCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetProject, err := manager.GetProject(d.Get("project_id").(string))
	if err != nil {
		return diag.Errorf("project_id: Error getting project: %s", err)
	}

	targetHypervisor, err := GetHypervisorById(d, manager, targetProject)
	if err != nil {
		return diag.Errorf("hypervisor_idЖ Error getting Hypervisor: %s", err)
	}

	vdc := rustack.NewVdc(d.Get("name").(string), targetHypervisor)

	// if we creating multiple vdc at once, there are need some time to get new vnid
	f := func() error { return targetProject.CreateVdc(&vdc) }
	err = repeatOnError(f, targetProject)

	if err != nil {
		return diag.Errorf("Error creating vdc: %s", err)
	}

	d.SetId(vdc.ID)
	log.Printf("[INFO] VDC created, ID: %s", d.Id())

	return resourceRustackVdcRead(ctx, d, meta)
}

func resourceRustackVdcRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	vdc, err := manager.GetVdc(d.Id())
	if err != nil {
		return diag.Errorf("id: Error getting vdc: %s", err)
	}

	flattenedProject := map[string]interface{}{
		"name":          vdc.Name,
		"project_id":    vdc.Project.ID,
		"hypervisor_id": vdc.Hypervisor.ID,
	}

	if err := setResourceDataFromMap(d, flattenedProject); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(vdc.ID)
	return nil
}

func resourceRustackVdcUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()

	vdc, err := manager.GetVdc(d.Id())
	if err != nil {
		return diag.Errorf("id: Error getting vdc: %s", err)
	}

	err = vdc.Rename(d.Get("name").(string))
	if err != nil {
		return diag.Errorf("name: Error rename vdc: %s", err)
	}

	return resourceRustackVdcRead(ctx, d, meta)
}

func resourceRustackVdcDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	vdc, err := manager.GetVdc(d.Id())
	if err != nil {
		return diag.Errorf("id: Error getting vdc: %s", err)
	}

	err = vdc.Delete()
	if err != nil {
		return diag.Errorf("Error deleting vdc: %s", err)
	}

	return nil
}
