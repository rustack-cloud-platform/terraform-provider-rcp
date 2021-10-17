package rustack_terraform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/hashstructure/v2"
)

func dataSourceRustackNetworks() *schema.Resource {
	args := Defaults()
	args.injectContextVdcById()
	args.injectResultListNetwork()

	return &schema.Resource{
		ReadContext: dataSourceRustackNetworksRead,
		Schema:      args,
	}
}

func dataSourceRustackNetworksRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting vdc: %s", err)
	}

	allNetworks, err := targetVdc.GetNetworks()
	if err != nil {
		return diag.Errorf("Error retrieving networks: %s", err)
	}

	flattenedRecords := make([]map[string]interface{}, len(allNetworks))
	for i, network := range allNetworks {
		allSubnets, err := network.GetSubnets()
		if err != nil {
			return diag.Errorf("Error getting subnets")
		}

		flattenRecords2 := make([]map[string]interface{}, len(allSubnets))
		for i2, subnet := range allSubnets {
			dnsStrings := make([]string, len(subnet.DnsServers))
			for i3, dns := range subnet.DnsServers {
				dnsStrings[i3] = dns.DNSServer
			}

			flattenRecords2[i2] = map[string]interface{}{
				"id":       subnet.ID,
				"cidr":     subnet.CIDR,
				"dhcp":     subnet.IsDHCP,
				"gateway":  subnet.Gateway,
				"start_ip": subnet.StartIp,
				"end_ip":   subnet.EndIp,
				"dns":      dnsStrings,
			}
		}

		flattenedRecords[i] = map[string]interface{}{
			"id":      network.ID,
			"name":    network.Name,
			"subnets": flattenRecords2,
		}
	}

	hash, err := hashstructure.Hash(allNetworks, hashstructure.FormatV2, nil)
	if err != nil {
		diag.Errorf("unable to set `networks` attribute: %s", err)
	}

	d.SetId(fmt.Sprintf("networks/%d", hash))
	// d.Set("vdc_id", nil)
	// d.Set("vdc_name", nil)

	if err := d.Set("networks", flattenedRecords); err != nil {
		return diag.Errorf("unable to set `networks` attribute: %s", err)
	}

	return nil
}
