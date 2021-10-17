package rustack_terraform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/hashstructure/v2"
)

func dataSourceRustackVms() *schema.Resource {
	args := Defaults()
	args.injectContextVdcById()
	args.injectResultListVm()

	return &schema.Resource{
		ReadContext: dataSourceRustackVmsRead,
		Schema:      args,
	}
}

func dataSourceRustackVmsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting vdc: %s", err)
	}

	allVms, err := targetVdc.GetVms()
	if err != nil {
		return diag.Errorf("Error retrieving vms: %s", err)
	}

	flattenedRecords := make([]map[string]interface{}, len(allVms))
	for i, vm := range allVms {
		flattenedRecords[i] = map[string]interface{}{
			"id":            vm.ID,
			"name":          vm.Name,
			"cpu":           vm.Cpu,
			"ram":           vm.Ram,
			"template_id":   vm.Template.ID,
			"template_name": vm.Template.Name,
			"floating_ip":   nil,
		}

		if vm.Floating != nil {
			flattenedRecords[i]["floating_ip"] = vm.Floating.IpAddress
		}
	}

	hash, err := hashstructure.Hash(allVms, hashstructure.FormatV2, nil)
	if err != nil {
		diag.Errorf("unable to set `vms` attribute: %s", err)
	}

	d.SetId(fmt.Sprintf("vms/%d", hash))

	if err := d.Set("vms", flattenedRecords); err != nil {
		return diag.Errorf("unable to set `vms` attribute: %s", err)
	}

	return nil
}
