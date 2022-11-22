package rustack_terraform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/hashstructure/v2"
)

func dataSourceRustackPorts() *schema.Resource {
	args := Defaults()
	args.injectContextVdcById()
	args.injectResultListPort()

	return &schema.Resource{
		ReadContext: dataSourceRustackPortsRead,
		Schema:      args,
	}
}

func dataSourceRustackPortsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting vdc: %s", err)
	}

	allPorts, err := targetVdc.GetPorts()
	if err != nil {
		return diag.Errorf("Error retrieving ports: %s", err)
	}

	flattenedRecords := make([]map[string]interface{}, len(allPorts))
	for i, port := range allPorts {
		flattenedRecords[i] = map[string]interface{}{
			"id":         port.ID,
			"ip_address": port.IpAddress,
			"network":    port.Network.ID,
		}
	}

	hash, err := hashstructure.Hash(allPorts, hashstructure.FormatV2, nil)
	if err != nil {
		diag.Errorf("unable to set `ports` attribute: %s", err)
	}

	d.SetId(fmt.Sprintf("port/%d", hash))

	if err := d.Set("ports", flattenedRecords); err != nil {
		return diag.Errorf("unable to set `ports` attribute: %s", err)
	}

	return nil
}
