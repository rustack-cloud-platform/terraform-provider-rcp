package rustack_terraform

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rustack-cloud-platform/rcp-go/rustack"
)

func resourceRustackVdc() *schema.Resource {
	args := Defaults()
	args.injectCreateVdc()
	args.injectContextHypervisorById()

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
		CustomizeDiff: func(ctx context.Context, rd *schema.ResourceDiff, i interface{}) error {
			if rd.Id() != "" && !rd.HasChange("project_id") {
				rd.Clear("id")
				rd.Clear("default_network_id")
			}
			return nil
		},
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
		return diag.Errorf("hypervisor_id: Error getting Hypervisor: %s", err)
	}

	vdc := rustack.NewVdc(d.Get("name").(string), targetHypervisor)
	vdc.Tags = unmarshalTagNames(d.Get("tags"))
	// if we creating multiple vdc at once, there are need some time to get new vnid
	f := func() error { return targetProject.CreateVdc(&vdc) }
	err = repeatOnError(f, targetProject)

	if err != nil {
		return diag.Errorf("Error creating vdc: %s", err)
	}

	vdc.WaitLock()
	if mtu, ok := d.GetOk("default_network_mtu"); ok {
		networks, err := vdc.GetNetworks(rustack.Arguments{"defaults_only": "true"})
		if err != nil {
			return diag.Errorf("Error getting vdc networks: %s", err)
		}
		if len(networks) != 1 {
			return diag.Errorf("Expected 1 network, got %d networks", len(networks))
		}
		network := networks[0]
		mtuValue := mtu.(int)
		network.Mtu = &mtuValue
		err = network.Update()
		if err != nil {
			return diag.Errorf("Error updating vdc default network: %s", err)
		}
	}
	vdc.GetNetworks()
	d.SetId(vdc.ID)
	log.Printf("[INFO] VDC created, ID: %s", d.Id())

	return resourceRustackVdcRead(ctx, d, meta)
}

func resourceRustackVdcRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	vdc, err := manager.GetVdc(d.Id())
	if err != nil {
		if err.(*rustack.RustackApiError).Code() == 404 {
			d.SetId("")
			return nil
		} else {
			return diag.Errorf("id: Error getting vdc: %s", err)
		}
	}
	networks, err := vdc.GetNetworks(rustack.Arguments{"defaults_only": "true"})
	if err != nil {
		return diag.Errorf("error getting default network: %s", err)
	}
	if len(networks) != 1 {
		return diag.Errorf("expected 1 default network, receive %d default networks", len(networks))
	}
	network := networks[0]
	subnets, err := network.GetSubnets()
	if err != nil {
		return diag.Errorf("subnets: Error getting subnets: %s", err)
	}

	flattenedSubnets := make([]map[string]interface{}, len(subnets))
	for i, subnet := range subnets {
		dnsStrings := make([]string, len(subnet.DnsServers))
		for i2, dns := range subnet.DnsServers {
			dnsStrings[i2] = dns.DNSServer
		}
		flattenedSubnets[i] = map[string]interface{}{
			"id":       subnet.ID,
			"cidr":     subnet.CIDR,
			"dhcp":     subnet.IsDHCP,
			"gateway":  subnet.Gateway,
			"start_ip": subnet.StartIp,
			"end_ip":   subnet.EndIp,
			"dns":      dnsStrings,
		}
	}
	flattenedVdc := map[string]interface{}{
		"name":                    vdc.Name,
		"project_id":              vdc.Project.ID,
		"default_network_id":      network.ID,
		"default_network_name":    network.Name,
		"default_network_subnets": flattenedSubnets,
		"default_network_mtu":     network.Mtu,
		"tags":                    marshalTagNames(vdc.Tags),
	}

	if err := setResourceDataFromMap(d, flattenedVdc); err != nil {
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
	if d.HasChange("hypervisor_id") {
		return diag.Errorf("hypervisor_id: you can`t change hypervisor type on created vdc")
	}
	if d.HasChange("name") {
		vdc.Name = d.Get("name").(string)
	}
	if d.HasChange("tags") {
		vdc.Tags = unmarshalTagNames(d.Get("tags"))
	}
	err = vdc.Update()
	if err != nil {
		return diag.Errorf("name: Error rename vdc: %s", err)
	}
	if d.HasChange("default_network_mtu") {
		networks, err := vdc.GetNetworks(rustack.Arguments{"defaults_only": "true"})
		if err != nil {
			return diag.Errorf("Error getting vdc networks: %s", err)
		}
		if len(networks) != 1 {
			return diag.Errorf("Expected 1 network, got %d networks", len(networks))
		}
		network := networks[0]

		if mtu, ok := d.GetOk("default_network_mtu"); ok {
			mtuValue := mtu.(int)
			network.Mtu = &mtuValue
		} else {
			network.Mtu = nil
		}
		err = network.Update()
		if err != nil {
			return diag.Errorf("Error updating vdc default network: %s", err)
		}
	}

	vdc.WaitLock()

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
	vdc.WaitLock()

	return nil
}
