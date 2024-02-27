package rustack_terraform

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rustack-cloud-platform/rcp-go/rustack"
)

func resourceRustackDns() *schema.Resource {
	args := Defaults()
	args.injectCreateDns()

	return &schema.Resource{
		CreateContext: resourceRustackDnsCreate,
		ReadContext:   resourceRustackDnsRead,
		DeleteContext: resourceRustackDnsDelete,
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

func resourceRustackDnsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	project, err := GetProjectById(d, manager)
	if err != nil {
		return diag.Errorf("project_id: Error getting Project: %s", err)
	}
	name := d.Get("name").(string)
	newDns := rustack.NewDns(name)
	if !strings.HasSuffix(name, ".") {
		return diag.Errorf("name: must be ending by '.'")
	}
	newDns.Tags = unmarshalTagNames(d.Get("tags"))

	err = project.CreateDns(&newDns)
	if err != nil {
		return diag.Errorf("Error creating Dns: %s", err)
	}

	d.SetId(newDns.ID)
	log.Printf("[INFO] Dns created, ID: %s", d.Id())

	return resourceRustackDnsRead(ctx, d, meta)
}

func resourceRustackDnsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	Dns, err := manager.GetDns(d.Id())
	if err != nil {
		if err.(*rustack.RustackApiError).Code() == 404 {
			d.SetId("")
			return nil
		} else {
			return diag.Errorf("id: Error getting Dns: %s", err)
		}
	}

	d.SetId(Dns.ID)
	d.Set("name", Dns.Name)
	d.Set("project", Dns.Project.ID)
	d.Set("tags", marshalTagNames(Dns.Tags))

	return nil
}

func resourceRustackDnsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	dns, err := manager.GetDns(d.Id())
	if err != nil {
		return diag.Errorf("id: Error getting Dns: %s", err)
	}

	err = dns.Delete()
	if err != nil {
		return diag.Errorf("Error deleting Dns: %s", err)
	}

	return nil
}
