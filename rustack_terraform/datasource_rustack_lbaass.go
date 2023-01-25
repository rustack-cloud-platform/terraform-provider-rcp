package rustack_terraform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/hashstructure/v2"
)

func dataSourceRustackLoadBalancers() *schema.Resource {
	args := Defaults()
	args.injectContextVdcById()
	args.injectResultListLbaas()

	return &schema.Resource{
		ReadContext: dataSourceRustackLoadBalancersRead,
		Schema:      args,
	}
}

func dataSourceRustackLoadBalancersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting vdc: %s", err)
	}

	allLoadBalancers, err := targetVdc.GetLoadBalancers()
	if err != nil {
		return diag.Errorf("Error retrieving lbs: %s", err)
	}

	flattenedRecords := make([]map[string]interface{}, len(allLoadBalancers))
	for i, lb := range allLoadBalancers {
		flattenedRecords[i] = map[string]interface{}{
			"id":            lb.ID,
			"name":          lb.Name,
		}

		if lb.Floating != nil {
			flattenedRecords[i]["floating"] = true
			flattenedRecords[i]["floating_ip"] = lb.Floating.IpAddress
		}
	}

	hash, err := hashstructure.Hash(allLoadBalancers, hashstructure.FormatV2, nil)
	if err != nil {
		diag.Errorf("unable to set `lbs` attribute: %s", err)
	}

	d.SetId(fmt.Sprintf("lbs/%d", hash))

	if err := d.Set("lbaass", flattenedRecords); err != nil {
		return diag.Errorf("unable to set `lbs` attribute: %s", err)
	}

	return nil
}
