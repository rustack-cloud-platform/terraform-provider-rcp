package rustack_terraform

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	// "github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"

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
		// https://www.terraform.io/docs/extend/resources/customizing-differences.html
		// CustomizeDiff: customdiff.Sequence(
		// 	// Clear the diff if the old and new allocation_pools are the same.
		// 	func(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
		// 		if diff.Id() != "" {
		// 			// hyp_id, proj_id
		// 			// diff.HasChange("project_id")
		// 			o, n := diff.GetChange("project_id")
		// 			log.Printf("[DEBUG] Project id. Old: %s, New: %s", o, n)

		// 			o, n = diff.GetChange("project_name")
		// 			log.Printf("[DEBUG] Project name. Old: %s, New: %s", o, n)
		// 		}
		// 		// return networkingSubnetV2AllocationPoolsCustomizeDiff(diff)
		// 		// https://github.com/terraform-provider-openstack/terraform-provider-openstack/blob/be9221d28ec309b57d3bba067abb26ddad35ecf8/openstack/networking_subnet_v2.go#L142
		// 		// if diff.Id() != "" && diff.HasChange("allocation_pools") {
		// 		// 	// o, n := diff.GetChange("allocation_pools")
		// 		// 	// oldPools := o.([]interface{})
		// 		// 	// newPools := n.([]interface{})

		// 		// 	samePools := true // networkingSubnetV2AllocationPoolsMatch(oldPools, newPools)

		// 		// 	if samePools {
		// 		// 		log.Printf("[DEBUG] allocation_pools have not changed. clearing diff")
		// 		// 		return diff.Clear("allocation_pools")
		// 		// 	}
		// 		// }
		// 		return nil
		// 	},
		// ),
	}
}

func resourceRustackVdcCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetProject, err := manager.GetProject(d.Get("project_id").(string))
	if err != nil {
		return diag.Errorf("Error getting project: %s", err)
	}

	targetHypervisor, err := GetHypervisorById(d, manager, targetProject)
	if err != nil {
		return diag.Errorf("Error getting Hypervisor: %s", err)
	}

	vdc := rustack.NewVdc(d.Get("name").(string), targetHypervisor)

	err = targetProject.CreateVdc(&vdc)
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
		return diag.Errorf("Error getting vdc: %s", err)
	}

	// d.SetId(vdc.ID)
	// d.Set("name", vdc.Name)

	// return nil
	flattenedProject := map[string]interface{}{
		// "id":         vdc.ID,
		"name":       vdc.Name,
		"project_id": vdc.Project.ID,
		// "project_name":    vdc.Project.Name,
		"hypervisor_id": vdc.Hypervisor.ID,
		// "hypervisor_name": vdc.Hypervisor.Name,
		// "hypervisor": strings.ToLower(vdc.Hypervisor.Type),
		// "vdc_id":          "",
		// "vdc_name":        "",
		// "project_id":      "",
		// "project_name":    "",
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
		return diag.Errorf("Error getting vdc: %s", err)
	}

	err = vdc.Rename(d.Get("name").(string))
	if err != nil {
		return diag.Errorf("Error rename vdc: %s", err)
	}

	return resourceRustackVdcRead(ctx, d, meta)
}

func resourceRustackVdcDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	vdc, err := manager.GetVdc(d.Id())
	if err != nil {
		return diag.Errorf("Error getting vdc: %s", err)
	}

	err = vdc.Delete()
	if err != nil {
		return diag.Errorf("Error deleting vdc: %s", err)
	}

	return nil
}
