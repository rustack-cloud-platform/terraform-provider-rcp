package rustack_terraform

import (
	"context"
	"log"
	"time"

	"github.com/pilat/rustack-go/rustack"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRustackFirewallTemplate() *schema.Resource {
	args := Defaults()
	args.injectContextVdcById()
	args.injectCreateFirewallTemplate()

	return &schema.Resource{
		CreateContext: resourceRustackFirewallTemplateCreate,
		ReadContext:   resourceRustackFirewallTemplateRead,
		UpdateContext: resourceRustackFirewallTemplateUpdate,
		DeleteContext: resourceRustackFirewallTemplateDelete,
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

func resourceRustackFirewallTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("vdc_id: Error getting VDC: %s", err)
	}

	newFirewallTemplate := rustack.NewFirewallTemplate(d.Get("name").(string))
	err = targetVdc.CreateFirewallTemplate(&newFirewallTemplate)
	if err != nil {
		return diag.Errorf("Error creating Firewall Template: %s", err)
	}

	d.SetId(newFirewallTemplate.ID)
	log.Printf("[INFO] FirewallTemplate created, ID: %s", d.Id())

	return resourceRustackFirewallTemplateRead(ctx, d, meta)
}

func resourceRustackFirewallTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	firewallTemplate, err := manager.GetFirewallTemplate(d.Id())
	if err != nil {
		if err.(*rustack.RustackApiError).Code() == 404 {
			d.SetId("")
			return nil
		} else {
			return diag.Errorf("id: Error getting Firewall Template: %s", err)
		}
	}

	d.SetId(firewallTemplate.ID)
	d.Set("name", firewallTemplate.Name)

	return nil
}

func resourceRustackFirewallTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()

	firewallTemplate, err := manager.GetFirewallTemplate(d.Id())
	if err != nil {
		return diag.Errorf("id: Error getting FirewallTemplate: %s", err)
	}

	if d.HasChange("name") {
		if err = firewallTemplate.Rename(d.Get("name").(string)); err != nil {
			return diag.Errorf("name: Error rename Firewall Template: %s", err)
		}
	}

	return resourceRustackFirewallTemplateRead(ctx, d, meta)
}

func resourceRustackFirewallTemplateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	FirewallTemplate, err := manager.GetFirewallTemplate(d.Id())
	if err != nil {
		return diag.Errorf("id: Error getting FirewallTemplate: %s", err)
	}

	err = FirewallTemplate.Delete()
	if err != nil {
		return diag.Errorf("Error deleting FirewallTemplate: %s", err)
	}

	return nil
}
