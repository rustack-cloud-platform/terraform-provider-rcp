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

	var diagErr diag.Diagnostics
	if d.Get("system").(bool) {
		diagErr = setSetviceRouter(d, manager)
	} else {
		diagErr = createRouter(d, manager)
	}
	if diagErr != nil {
		return diagErr
	}

	return resourceRustackRouterRead(ctx, d, meta)
}

func resourceRustackRouterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	router, err := manager.GetRouter(d.Id())
	if err != nil {
		return diag.Errorf("Error getting Router: %s", err)
	}

	d.SetId(router.ID)
	d.Set("name", router.Name)
	d.Set("floating", router.Floating)
	if router.Floating != nil {
		d.Set("floating_id", router.Floating.ID)
	}
	networks := []string{}
	for _, port := range router.Ports {
		networks = append(networks, port.Network.ID)
	}
	d.Set("networks", networks)
	d.Set("vdc_id", router.Vdc.Id)

	return nil
}

func resourceRustackRouterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	router, err := manager.GetRouter(d.Id())
	if err != nil {
		return diag.Errorf("Error getting Router: %s", err)
	}
	router.Name = d.Get("name").(string)
	if err := syncFloating(d, router); err != nil {
		return diag.FromErr(err)
	}

	// Delete disconnected ports and create a new if connected
	err = syncRouterPorts(d, manager, router)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceRustackRouterRead(ctx, d, meta)
}

func resourceRustackRouterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	routerId := d.Id()
	router, err := manager.GetRouter(routerId)
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting Router: %s", err)
	}

	if d.Get("system").(bool) {
		network := GetServiseNetworkByVdc(targetVdc)
		var newPort rustack.Port
		newPort.Network = network
		router.CreatePort(&newPort, router)
		router.WaitLock()

		for _, port := range router.Ports {
			network, err := manager.GetNetwork(port.Network.ID)
			if err != nil {
				return diag.FromErr(err)
			}
			if !network.IsDefault {
				if err = repeatOnError(port.Delete, port); err != nil {
					return diag.Errorf("Error deleting Router: %s", err)
				}
			}
			if router.Floating == nil {
				router.Floating = &rustack.Port{ID: "RANDOM_FIP"}
				if err = repeatOnError(router.Update, router); err != nil {
					return diag.Errorf("ERROR: Can't return router to default state: %s", err)
				}
			}
		}

		return nil
	}

	if err = repeatOnError(router.Delete, router); err != nil {
		return diag.Errorf("Error deleting Router: %s", err)
	}
	
	d.SetId("")
	log.Printf("[INFO] Router deleted, ID: %s", routerId)

	return nil
}

func setSetviceRouter(d *schema.ResourceData, manager *rustack.Manager) diag.Diagnostics {
	router, err := getSystemRouter(d, manager)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(router.ID)
	if router.Floating != nil {
		d.Set("floating_id", router.Floating.ID)
	}

	flattenedRecords := make([]string, len(router.Ports))
	for i, port := range router.Ports {
		flattenedRecords[i] = port.Network.ID
	}

	if err := d.Set("networks", flattenedRecords); err != nil {
		return diag.Errorf("unable to set `networks` attribute: %s", err)
	}

	if err := syncFloating(d, router); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Router, ID: %s", d.Id())

	return nil
}

func createRouter(d *schema.ResourceData, manager *rustack.Manager) (diagErr diag.Diagnostics) {
	vdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting Ports from vdc: %s", err)
	}
	if _, ok := d.GetOk("networks"); !ok {
		return diag.Errorf("ERROR: You should setup a network for non default routers")
	}

	var floatingIp *string = nil
	if d.Get("floating").(bool) {
		floatingIpStr := "RANDOM_FIP"
		floatingIp = &floatingIpStr
	}

	router := rustack.NewRouter(d.Get("name").(string), floatingIp)

	ports, err := preparePortsToConnect(manager, d)
	if err != nil {
		return diag.Errorf("Error getting errors: %s", err)
	}

	router.Vdc.Id = vdc.ID

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
	d.Set("floating", router.Floating)
	if router.Floating != nil {
		d.Set("floating_id", router.Floating.ID)
	}
	log.Printf("[INFO] Router created, ID: %s", router.ID)

	return
}

func getSystemRouter(d *schema.ResourceData, manager *rustack.Manager) (router *rustack.Router, err error) {
	vdc, err := GetVdcById(d, manager)
	if err != nil {
		return nil, fmt.Errorf("ERROR: Can't get Ports from vdc: %s", err)
	}
	routerList, err := vdc.GetRouters()
	if err != nil {
		return nil, fmt.Errorf("ERROR: Can't get routers from vdc: %s", err)
	}
	for _, router = range routerList {
		if router.IsDefault {
			break
		}
	}
	d.SetId(router.ID)

	err = syncRouterPorts(d, manager, router)
	if err != nil {
		return nil, err
	}

	// Connect ports
	ports, err := preparePortsToConnect(manager, d)
	if err != nil {
		return nil, fmt.Errorf("ERROR. getting ports: %s", err)
	}
	router.Ports = ports

	return
}

func syncRouterPorts(d *schema.ResourceData, manager *rustack.Manager, router *rustack.Router) (err error) {
	vdc, err := GetVdcById(d, manager)
	if err != nil {
		return fmt.Errorf("ERROR: Can't get Ports from vdc: %s", err)
	}
	// Connect ports
	ports, err := preparePortsToConnect(manager, d)
	if err != nil {
		return fmt.Errorf("ERROR: Can't get ports: %s", err)
	}
	router.Ports = ports
	if err = repeatOnError(router.Update, router); err != nil {
		return fmt.Errorf("ERROR: Can't update Router: %s", err)
	}
	log.Printf("[INFO] Updated Router, ID: %v", router)

	// disconnect ports
	oldPortList := router.Ports
	newPortList := d.Get("networks")
	for _, s1 := range oldPortList {
		found := false
		for _, s2 := range newPortList.(*schema.Set).List() {
			if s1.Network.ID == s2 {
				found = true
				break
			}
		}
		if !found {
			vdcPorts, err := vdc.GetPorts()
			if err != nil {
				return fmt.Errorf("ERROR: Can't get Ports from vdc: %s", err)
			}
			var portId string
			for _, p := range vdcPorts {
				if p.Connected != nil && p.Connected.ID == router.ID && p.Network.ID == s1.Network.ID {
					portId = p.ID
					break
				}
			}
			if portId == "" {
				return fmt.Errorf("ERROR: Port with current network=%s not found: %s", s1.Network.ID, err)
			}
			portToDelete, err := manager.GetPort(portId)
			if err != nil {
				return fmt.Errorf("ERROR: Port not found: %s", err)
			}

			if err = repeatOnError(portToDelete.ForceDelete, portToDelete); err != nil {
				return fmt.Errorf("ERROR: One of the ports cannot be deleted: %s", err)
			}
		}
	}

	return
}

func syncFloating(d *schema.ResourceData, router *rustack.Router) (err error) {
	floating := d.Get("floating")
	if floating.(bool) && (router.Floating == nil) {
		// add floating if it was removed
		router.Floating = &rustack.Port{ID: "RANDOM_FIP"}
		if err = repeatOnError(router.Update, router); err != nil {
			return fmt.Errorf("ERROR: Can't update Router: %s", err)
		}
		d.Set("floating", true)
		d.Set("floating_id", router.Floating.ID)
	} else if !floating.(bool) && (router.Floating != nil) {
		// remove floating if needed
		router.Floating = nil

		if err = repeatOnError(router.Update, router); err != nil {
			return
		}
	} else if floating.(bool) && (router.Floating != nil) {
		d.Set("floating", true)
		d.Set("floating_id", router.Floating.ID)
	}
	return
}

func preparePortsToConnect(manager *rustack.Manager, d *schema.ResourceData) (ports []*rustack.Port, err error) {
	netArray := d.Get("networks").(*schema.Set).List()
	vdc, err := GetVdcById(d, manager)
	if err != nil {
		return nil, fmt.Errorf("ERROR getting Vdc: %s", err)
	}
	for _, networkId := range netArray {
		newIp := "0.0.0.0"

		vdcPorts, err := vdc.GetPorts()
		if err != nil {
			return nil, fmt.Errorf("ERROR getting Ports from vdc")
		}

		noRouter := false
		router, err := manager.GetRouter(d.Id())
		if err != nil {
			err = nil
			noRouter = true
		}

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
			if port.Network != nil && port.Network.ID == network.ID {
				port.Network.WaitLock()
			}
			if port.Connected != nil && port.Connected.ID == router.ID && port.Network.ID == network.ID {
				ports = append(ports, port)
				found = true
				break
			}
		}
		if found {
			continue
		}
		if err = router.CreatePort(&newPort, router); err != nil {
			return nil, err
		}
		ports = append(ports, &newPort)
	}
	return
}
