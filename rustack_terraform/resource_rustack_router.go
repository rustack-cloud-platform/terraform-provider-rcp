package rustack_terraform

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pilat/rustack-go/rustack"
)

func resourceRustackRouter() *schema.Resource {
	args := Defaults()
	args.injectContextVdcById()
	args.injectCreateRouter()

	return &schema.Resource{
		CreateContext: resourceRustackRouterCreate,
		ReadContext:   resourceRustackRouterRead,
		UpdateContext: resourceRustackRouterUpdate,
		DeleteContext: resourceRustackRouterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: args,
	}
}

func resourceRustackRouterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	vdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting Ports from vdc: %s", err)
	}
	router := rustack.NewRouter(d.Get("name").(string))

	ports, err := preparePortsToConnect(manager, d)
	if err != nil {
		return diag.Errorf("Error getting errors: %s", err)
	}

	router.Vdc.Id = vdc.ID
	if ipAddress, ok := d.GetOk("floating"); ok {
		d.Set("floating", ipAddress.(string))
		router.Floating, err = vdc.GetFloatingByAddress(ipAddress.(string))
		if err != nil {
			return diag.Errorf("Error floating set up: %s", err)
		}
	}

	log.Printf("[DEBUG] Router create request: %#v", router)
	vdc.WaitLock()

	// Wait networks and routers of each ports
	for _, port := range ports {
		port.Network.WaitLock()
		for {
			networkCheck, err := manager.GetNetwork(port.Network.ID)
			if err != nil {
				return diag.Errorf("Error creating Router: %s", err)
			}
			if len(networkCheck.Subnets) != 0 {
				networkCheck.Subnets[0].WaitLock()
				break
			}
			time.Sleep(time.Second)
		}
		portRouter, err := getRouterByNetwork(*manager, *port.Network)
		if err != nil {
			return diag.Errorf("Error creating Router: %s", err)
		}
		if portRouter != nil {
			portRouter.WaitLock()
		}
	}
	err = vdc.CreateRouter(&router, ports...)
	if err != nil {
		return diag.Errorf("Error creating Router: %s", err)
	}

	d.SetId(router.ID)
	d.Set("floating_id", router.Floating.ID)
	log.Printf("[INFO] Router created, ID: %s", d.Id())

	return resourceRustackRouterRead(ctx, d, meta)
}

func resourceRustackRouterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	Router, err := manager.GetRouter(d.Id())
	if err != nil {
		return diag.Errorf("Error getting Router: %s", err)
	}

	d.SetId(Router.ID)
	d.Set("name", Router.Name)

	return nil
}

func resourceRustackRouterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()

	vdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting Vdc: %s", err)
	}
	router, err := manager.GetRouter(d.Id())
	if err != nil {
		return diag.Errorf("Error getting Router: %s", err)
	}
	router.Name = d.Get("name").(string)
	if d.HasChange("floating") {
		router.Floating = nil
		if sourceFip, ok := d.GetOk("floating"); ok {
			fip, err := vdc.GetFloatingByAddress(sourceFip.(string))
			if err != nil {
				return diag.Errorf("Error getting fip address: %s", err)
			}
			router.Floating = fip
		}
	}

	// Delete disconnected ports
	oldPortList, newPortList := d.GetChange("networks")
	for _, s1 := range oldPortList.(*schema.Set).List() {
		found := false
		for _, s2 := range newPortList.(*schema.Set).List() {
			if s1 == s2 {
				found = true
				break
			}
		}
		if !found {
			vdcPorts, err := vdc.GetPorts()
			if err != nil {
				return diag.Errorf("Error getting Ports from vdc: %s", err)
			}
			var portId string
			for _, p := range vdcPorts {
				if p.Connected != nil && p.Connected.ID == router.ID && p.Network.ID == s1 {
					portId = p.ID
					break
				}
			}
			if portId == "" {
				return diag.Errorf("Port with current network=%s not found: %s", s1.(string), err)
			}
			portToDelete, err := manager.GetPort(portId)
			if err != nil {
				return diag.Errorf("Port not found: %s", err)
			}
			portToDelete.WaitLock()
			if err := portToDelete.Delete(); err != nil {
				return diag.Errorf("One of the ports cannot be deleted: %s", err)
			}
		}
	}

	// Connect ports
	ports, err := preparePortsToConnect(manager, d)
	if err != nil {
		return diag.Errorf("ERROR. getting ports: %s", err)
	}
	router.Ports = ports
	router.WaitLock()
	err = router.Update()
	log.Printf("[INFO] Updated Router, ID: %v", router)

	if err != nil {
		return diag.Errorf("Error updating Router: %s", err)
	}
	log.Printf("[INFO] Updated Router, ID: %v", router)

	return resourceRustackRouterRead(ctx, d, meta)
}

func resourceRustackRouterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()

	routerId := d.Id()

	router, err := manager.GetRouter(routerId)
	if err != nil {
		return diag.Errorf("Error getting Router: %s", err)
	}

	router.WaitLock()
	err = router.Delete()
	if err != nil {
		return diag.Errorf("Error deleting Router: %s", err)
	}

	d.SetId("")
	log.Printf("[INFO] Router deleted, ID: %s", routerId)

	return nil
}

func preparePortsToConnect(manager *rustack.Manager, d *schema.ResourceData) (ports []*rustack.Port, err error) {

	newIp := "0.0.0.0"
	netArray := d.Get("networks").(*schema.Set).List()
	vdc, err := GetVdcById(d, manager)
	if err != nil {
		return nil, fmt.Errorf("ERROR getting Vdc: %s", err)
	}
	vdcPorts, err := vdc.GetPorts()
	if err != nil {
		return nil, fmt.Errorf("ERROR getting Ports from vdc")
	}

	router, err := manager.GetRouter(d.Id())
	var noRouter bool
	if err != nil {
		err = nil
		noRouter = true
	}

	for _, networkId := range netArray {
		network, err := manager.GetNetwork(networkId.(string))
		if err != nil {
			return nil, err
		}
		if network.Vdc.Id != vdc.ID {
			return nil, errors.New("ERROR: Network should belong to routers vdc")
		}
		var newPort rustack.Port
		newPort.Network = network
		newPort.IpAddress = &newIp
		var found bool
		if noRouter {
			ports = append(ports, &newPort)
			continue
		}
		for _, port := range vdcPorts {
			if port.Connected != nil && port.Connected.ID == router.ID && port.Network.ID == network.ID {
				port.Network.WaitLock()
				ports = append(ports, port)
				found = true
				break
			}
		}
		if found {
			continue
		}
		router.WaitLock()
		if err = router.CreatePort(&newPort, router); err != nil {
			return nil, err
		}
		ports = append(ports, &newPort)
	}
	return
}
