package rustack_terraform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/hashstructure/v2"
)

func dataSourceRustackDnss() *schema.Resource {
	args := Defaults()
	args.injectContextProjectById()
	args.injectResultListDns()

	return &schema.Resource{
		ReadContext: dataSourceRustackDnssRead,
		Schema:      args,
	}
}

func dataSourceRustackDnssRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	project, err := GetProjectById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting project: %s", err)
	}

	allDns, err := project.GetDnss()
	if err != nil {
		return diag.Errorf("Error retrieving dnss: %s", err)
	}

	flattenedRecords := make([]map[string]interface{}, len(allDns))
	for i, dns := range allDns {
		flattenedRecords[i] = map[string]interface{}{
			"id":   dns.ID,
			"name": dns.Name,
		}
	}

	hash, err := hashstructure.Hash(allDns, hashstructure.FormatV2, nil)
	if err != nil {
		diag.Errorf("unable to set `dnss` attribute: %s", err)
	}

	d.SetId(fmt.Sprintf("dnss/%d", hash))

	if err := d.Set("dnss", flattenedRecords); err != nil {
		return diag.Errorf("unable to set `dnss` attribute: %s", err)
	}

	return nil
}
