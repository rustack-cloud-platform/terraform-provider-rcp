package rustack_terraform

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pilat/rustack-go/rustack"
)

func resourceRustackLbaas() *schema.Resource {
	args := Defaults()
	args.injectContextVdcById()
	args.injectCreateLbaas()

	return &schema.Resource{
		CreateContext: resourceRustackLbaasCreate,
		ReadContext:   resourceRustackLbaasRead,
		UpdateContext: resourceRustackLbaasUpdate,
		DeleteContext: resourceRustackLbaasDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: args,
	}
}

func resourceRustackLbaasCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()

	vdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("vdc_id: Error getting vdc : %s", err)
	}

	// create port
	var floatingIp *rustack.Port = nil
	if d.Get("floating").(bool) {
		floatingIp = &rustack.Port{ID: "RANDOM_FIP"}
	}
	portPrefix := "port.0"
	lbaasPort := d.Get("port.0").(map[string]interface{})

	network, err := manager.GetNetwork(lbaasPort["network_id"].(string))
	if err != nil {
		return diag.Errorf("network_id: Error getting network by id: %s", err)
	}
	network.WaitLock()
	firewalls := make([]*rustack.FirewallTemplate, 0)
	ipAddressStr := d.Get(MakePrefix(&portPrefix, "ip_address")).(string)
	if ipAddressStr == "" {
		ipAddressStr = "0.0.0.0"
	}
	port := rustack.NewPort(network, firewalls, ipAddressStr)

	newLbaas := rustack.NewLoadBalancer(d.Get("name").(string), vdc, &port, floatingIp)

	err = vdc.Create(&newLbaas)
	if err != nil {
		return diag.Errorf("Error creating Lbaas: %s", err)
	}
	newLbaas.WaitLock()
	d.SetId(newLbaas.ID)
	return resourceRustackLbaasRead(ctx, d, meta)
}

func resourceRustackLbaasRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diagErr diag.Diagnostics) {
	manager := meta.(*CombinedConfig).rustackManager()
	lbaas, err := manager.GetLoadBalancer(d.Id())
	if err != nil {
		if err.(*rustack.RustackApiError).Code() == 404 {
			d.SetId("")
			return nil
		} else {
			return diag.Errorf("id: Error getting Lbaas: %s", err)
		}
	}
	d.SetId(lbaas.ID)
	d.Set("name", lbaas.Name)
	d.Set("floating", lbaas.Floating != nil)
	d.Set("floating_ip", "")
	if lbaas.Floating != nil {
		d.Set("floating_ip", lbaas.Floating.IpAddress)
	}
	lbaasPort := make([]interface{}, 1)
	lbaasPort[0] = map[string]interface{}{
		"ip_address": lbaas.Port.IpAddress,
		"network_id": lbaas.Port.Network.ID,
	}
	d.Set("port", lbaasPort)
	d.Set("vdc_id", lbaas.Vdc.ID)

	return
}

func resourceRustackLbaasUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	lbaas, err := manager.GetLoadBalancer(d.Id())
	if err != nil {
		return diag.Errorf("id: Error getting Lbaas: %s", err)
	}
	if d.HasChange("name") {
		lbaas.Name = d.Get("name").(string)
	}
	if d.HasChange("floating") {
		if !d.Get("floating").(bool) {
			lbaas.Floating = &rustack.Port{IpAddress: nil}
		} else {
			lbaas.Floating = &rustack.Port{ID: "RANDOM_FIP"}
		}
		d.Set("floating", lbaas.Floating != nil)
	}
	lbaasPort := d.Get("port.0").(map[string]interface{})
	ip_address := lbaasPort["ip_address"].(string)
	if ip_address != *lbaas.Port.IpAddress {
		lbaas.Port.IpAddress = &ip_address
	}
	if err := repeatOnError(lbaas.Update, lbaas); err != nil {
		return diag.Errorf("Error updating lbaas: %s", err)
	}
	lbaas.WaitLock()

	return resourceRustackLbaasRead(ctx, d, meta)
}

func resourceRustackLbaasDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	lbaasId := d.Id()

	lbaas, err := manager.GetLoadBalancer(lbaasId)
	if err != nil {
		return diag.Errorf("id: Error getting Lbaas: %s", err)
	}

	lbaas.Delete()
	if err != nil {
		return diag.Errorf("Error deleting Lbaas: %s", err)
	}
	lbaas.WaitLock()

	d.SetId("")
	log.Printf("[INFO] Lbaas deleted, ID: %s", lbaasId)

	return nil
}
