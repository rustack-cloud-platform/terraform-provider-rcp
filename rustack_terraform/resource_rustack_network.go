package rustack_terraform

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pilat/rustack-go/rustack"
)

func resourceRustackNetwork() *schema.Resource {
	args := Defaults()
	args.injectContextVdcById()
	args.injectCreateNetwork()

	return &schema.Resource{
		CreateContext: resourceRustackNetworkCreate,
		ReadContext:   resourceRustackNetworkRead,
		UpdateContext: resourceRustackNetworkUpdate,
		DeleteContext: resourceRustackNetworkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: args,
	}
}

func resourceRustackNetworkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("vdc_id: Error getting VDC: %s", err)
	}

	log.Printf("[DEBUG] subnetInfo: %#v", targetVdc)
	network := rustack.NewNetwork(d.Get("name").(string))

	targetVdc.WaitLock()
	if err = targetVdc.CreateNetwork(&network); err != nil {
		return diag.Errorf("Error creating network: %s", err)
	}
	d.SetId(network.ID)

	diag := createSubnet(d, manager)
	if diag != nil {
		return diag
	}
	network.WaitLock()

	log.Printf("[INFO] Network created, ID: %s", d.Id())

	return resourceRustackNetworkRead(ctx, d, meta)
}

func resourceRustackNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	network, err := manager.GetNetwork(d.Id())
	if err != nil {
		if err.(*rustack.RustackApiError).Code() == 404 {
			d.SetId("")
			return nil
		} else {
			return diag.Errorf("id: Error getting network: %s", err)
		}
	}

	d.Set("name", network.Name)

	subnets, err := network.GetSubnets()
	if err != nil {
		return diag.Errorf("subnets: Error getting subnets: %s", err)
	}

	flattenedRecords := make([]map[string]interface{}, len(subnets))
	for i, subnet := range subnets {
		dnsStrings := make([]string, len(subnet.DnsServers))
		for i2, dns := range subnet.DnsServers {
			dnsStrings[i2] = dns.DNSServer
		}
		flattenedRecords[i] = map[string]interface{}{
			"id":       subnet.ID,
			"cidr":     subnet.CIDR,
			"dhcp":     subnet.IsDHCP,
			"gateway":  subnet.Gateway,
			"start_ip": subnet.StartIp,
			"end_ip":   subnet.EndIp,
			"dns":      dnsStrings,
		}
	}

	if err := d.Set("subnets", flattenedRecords); err != nil {
		return diag.Errorf("subnets: unable to set `subnet` attribute: %s", err)
	}

	return nil
}

func resourceRustackNetworkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()

	network, err := manager.GetNetwork(d.Id())
	if err != nil {
		return diag.Errorf("id: Error getting network: %s", err)
	}

	if d.HasChange("name") {
		err = network.Rename(d.Get("name").(string))
		if err != nil {
			return diag.Errorf("name: Error rename network: %s", err)
		}
	}

	if d.HasChange("subnets") {
		diagErr := updateSubnet(d, manager)
		if diagErr != nil {
			return diagErr
		}
	}
	network.WaitLock()

	return resourceRustackNetworkRead(ctx, d, meta)
}

func resourceRustackNetworkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	network, err := manager.GetNetwork(d.Id())
	if err != nil {
		return diag.Errorf("id: Error getting network: %s", err)

	}

	if err = repeatOnError(network.Delete, network); err != nil {
		return diag.Errorf("Error deleting network: %s", err)
	}
	network.WaitLock()

	return nil
}

func createSubnet(d *schema.ResourceData, manager *rustack.Manager) (diagErr diag.Diagnostics) {
	subnets := d.Get("subnets").([]interface{})
	log.Printf("[DEBUG] subnets: %#v", subnets)
	network, err := manager.GetNetwork(d.Id())
	if err != nil {
		return diag.Errorf("id: Unable to get network: %s", err)
	}

	for _, subnetInfo := range subnets {
		log.Printf("[DEBUG] subnetInfo: %#v", subnetInfo)
		subnetInfo2 := subnetInfo.(map[string]interface{})

		// Create subnet
		subnet := rustack.NewSubnet(subnetInfo2["cidr"].(string), subnetInfo2["gateway"].(string), subnetInfo2["start_ip"].(string), subnetInfo2["end_ip"].(string), subnetInfo2["dhcp"].(bool))

		if err := network.CreateSubnet(&subnet); err != nil {
			return diag.Errorf("subnets: Error creating subnet: %s", err)
		}

		dnsServersList := subnetInfo2["dns"].([]interface{})
		dnsServers := make([]*rustack.SubnetDNSServer, len(dnsServersList))
		for i, dns := range dnsServersList {
			s1 := rustack.NewSubnetDNSServer(dns.(string))
			dnsServers[i] = &s1
		}

		if err := subnet.UpdateDNSServers(dnsServers); err != nil {
			return diag.Errorf("dns: Error Update DNS Servers: %s", err)
		}

	}

	return
}

func updateSubnet(d *schema.ResourceData, manager *rustack.Manager) (diagErr diag.Diagnostics) {

	subnets := d.Get("subnets").([]interface{})
	network, err := manager.GetNetwork(d.Id())
	if err != nil {
		return diag.Errorf("id: Unable to get network: %s", err)
	}
	subnetsRaw, err := network.GetSubnets()
	if err != nil {
		return diag.Errorf("subnets: Unable to get subnets: %s", err)
	}

	for _, subnetInfo := range subnets {
		subnetInfo2 := subnetInfo.(map[string]interface{})
		for _, subnet := range subnetsRaw {
			if subnet.ID == subnetInfo2["id"] {
				if subnet.Gateway != subnetInfo2["gateway"] || subnet.StartIp != subnetInfo2["start_ip"] || subnet.EndIp != subnetInfo2["end_ip"] {
					return diag.Errorf("You cannot change params (gateway, start_ip, end_ip)")
				}
				newDHCPValue := subnetInfo2["dhcp"].(bool)
				if subnet.IsDHCP != newDHCPValue {
					if newDHCPValue {
						err = subnet.EnableDHCP()
						if err != nil {
							return diag.Errorf("dhcp: Unable to toggle DHCP: %s", err)
						}
					} else {
						err = subnet.DisableDHCP()
						if err != nil {
							return diag.Errorf("dhcp: Unable to toggle DHCP: %s", err)
						}
					}
				}

				// Set DNS again
				dnsServersList := subnetInfo2["dns"].([]interface{})
				dnsServers := make([]*rustack.SubnetDNSServer, len(dnsServersList))
				for i, dns := range dnsServersList {
					s1 := rustack.NewSubnetDNSServer(dns.(string))
					dnsServers[i] = &s1
				}

				subnet.UpdateDNSServers(dnsServers)
			}
		}
	}

	return
}
