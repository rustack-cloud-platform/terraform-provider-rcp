package rustack_terraform

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rustack-cloud-platform/rcp-go/rustack"
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
		diagErr = setServiceRouter(d, manager)
	} else {
		diagErr = createRouter(d, manager)
	}
	if diagErr != nil {
		return diagErr
	}

	return resourceRustackRouterRead(ctx, d, meta)
}

func resourceRustackRouterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diagErr diag.Diagnostics) {
	manager := meta.(*CombinedConfig).rustackManager()
	router, err := manager.GetRouter(d.Id())
	if err != nil {
		if err.(*rustack.RustackApiError).Code() == 404 {
			d.SetId("")
			return nil
		} else {
			return diag.Errorf("id: Error getting Router: %s", err)
		}
	}

	d.SetId(router.ID)
	d.Set("name", router.Name)

	d.Set("floating", router.Floating != nil)
	d.Set("floating_ip", "")
	if router.Floating != nil {
		d.Set("floating_ip", router.Floating.IpAddress)
	}

	ports := make([]*string, len(router.Ports))
	for i, port := range router.Ports {
		ports[i] = &port.ID
	}

	d.Set("ports", ports)
	d.Set("vdc_id", router.Vdc.Id)
	d.Set("tags", marshalTagNames(router.Tags))

	return
}

func resourceRustackRouterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	router, err := manager.GetRouter(d.Id())
	if err != nil {
		return diag.Errorf("id: Error getting Router: %s", err)
	}
	shouldUpdate := false
	if d.HasChange("name") {
		router.Name = d.Get("name").(string)
		shouldUpdate = true
	}
	if d.HasChange("tags") {
		router.Tags = unmarshalTagNames(d.Get("tags"))
		shouldUpdate = true
	}
	if shouldUpdate {
		if err := router.Update(); err != nil {
			return diag.Errorf("error on router's update %s", err)
		}
	}

	if err := syncFloating(d, router); err != nil {
		return diag.FromErr(err)
	}

	// Disconnect ports and connect new
	err = syncRouterPorts(d, manager, router)
	if err != nil {
		return diag.FromErr(err)
	}
	router.WaitLock()

	return resourceRustackRouterRead(ctx, d, meta)
}

func resourceRustackRouterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	portsIds := d.Get("ports").(*schema.Set).List()
	routerId := d.Id()
	router, err := manager.GetRouter(routerId)
	if err != nil {
		return diag.Errorf("id: Error getting Router: %s", err)
	}

	// Disconnect custom ports from system router
	if d.Get("system").(bool) {
		if err != nil {
			return diag.Errorf("Error getting service Network: %s", err)
		}

		for _, port := range router.Ports {
			network, err := manager.GetNetwork(port.Network.ID)
			if err != nil {
				return diag.FromErr(err)
			}
			if !network.IsDefault {
				err = router.DisconnectPort(port)
				if err != nil {
					return diag.FromErr(err)
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

	// Detach ports and delete custom router
	for _, portId := range portsIds {
		port, err := manager.GetPort(portId.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		err = router.DisconnectPort(port)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if err = repeatOnError(router.Delete, router); err != nil {
		return diag.Errorf("Error deleting Router: %s", err)
	}
	router.WaitLock()

	d.SetId("")
	log.Printf("[INFO] Router deleted, ID: %s", routerId)

	return nil
}

func setServiceRouter(d *schema.ResourceData, manager *rustack.Manager) diag.Diagnostics {
	router, err := getSystemRouter(d, manager)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(router.ID)
	if router.Floating != nil {
		d.Set("floating_id", router.Floating.ID)
	}

	portsIds := d.Get("ports").(*schema.Set).List()
	ports := make([]*rustack.Port, len(portsIds))

	for i, portId := range portsIds {
		port, err := manager.GetPort(portId.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		ports[i] = port
	}

	d.Set("ports", ports)

	if err := syncFloating(d, router); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Router, ID: %s", d.Id())

	return nil
}

func createRouter(d *schema.ResourceData, manager *rustack.Manager) (diagErr diag.Diagnostics) {
	vdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("ports: Error getting Ports from vdc: %s", err)
	}
	if _, ok := d.GetOk("ports"); !ok {
		return diag.Errorf("ports: Error You should setup a port for non default routers")
	}

	var floatingIp *string = nil
	if d.Get("floating").(bool) {
		floatingIpStr := "RANDOM_FIP"
		floatingIp = &floatingIpStr
	}

	router := rustack.NewRouter(d.Get("name").(string), floatingIp)
	router.Tags = unmarshalTagNames(d.Get("tags"))
	portsIds := d.Get("ports").(*schema.Set).List()
	ports := make([]*rustack.Port, len(portsIds))

	for i, portId := range portsIds {
		port, err := manager.GetPort(portId.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		ports[i] = port
	}

	router.Vdc.Id = vdc.ID

	log.Printf("[DEBUG] Router create request: %#v", router)
	vdc.WaitLock()

	err = vdc.CreateRouter(&router, ports...)
	if err != nil {
		return diag.Errorf("Error creating Router: %s", err)
	}
	router.WaitLock()

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
	if router == nil {
		return nil, fmt.Errorf("ERROR: Default router not found in vdc %s", vdc.ID)
	}
	d.SetId(router.ID)
	tags := unmarshalTagNames(d.Get("tags"))
	shouldUpdate := false
	if len(tags) != len(router.Tags) {
		router.Tags = tags
		shouldUpdate = true
	} else {
		sort.Slice(tags, func(i, j int) bool { return tags[i].Name < tags[j].Name })
		sort.Slice(router.Tags, func(i, j int) bool { return tags[i].Name < tags[j].Name })
		for i := 0; i < len(tags); i++ {
			if tags[i].Name != router.Tags[i].Name {
				router.Tags = tags
				shouldUpdate = true
				break
			}
		}
	}
	if shouldUpdate {
		if err := router.Update(); err != nil {
			return nil, err
		}
	}

	err = syncRouterPorts(d, manager, router)
	if err != nil {
		return nil, err
	}

	// Connect ports

	portsIds := d.Get("ports").(*schema.Set).List()
	ports := make([]*rustack.Port, 0, len(portsIds))

	for _, portId := range portsIds {
		port, err := manager.GetPort(portId.(string))
		if err != nil {
			return nil, err
		}
		ports = append(ports, port)
	}

	router.Ports = ports

	return
}

func syncRouterPorts(d *schema.ResourceData, manager *rustack.Manager, router *rustack.Router) (err error) {
	portsIds := d.Get("ports").(*schema.Set).List()
	router_id := d.Id()

	for _, port := range router.Ports {
		found := false
		for _, portId := range portsIds {
			if portId == port.ID {
				found = true
				break
			}
		}

		if !found {
			if port.Connected != nil && port.Connected.ID == router_id {
				log.Printf("Port %s found on vm and not mentioned in the state."+
					" Port will be detached", port.ID)
				router.DisconnectPort(port)
				port.WaitLock()
			}
		}
	}

	for _, portId := range portsIds {
		found := false
		for _, port := range router.Ports {
			if port.ID == portId {
				found = true
				break
			}
		}

		if !found {
			port, err := manager.GetPort(portId.(string))
			if err != nil {
				return fmt.Errorf("ports: getting Port from vdc")
			}
			if port.Connected != nil && port.Connected.Type == "vm_int" {
				return fmt.Errorf("ports: Unable to bind a port that is already connected to the server")
			}
			if port.Connected != nil && port.Connected.ID != router_id {
				router.DisconnectPort(port)
				port.WaitLock()
			}
			port, err = manager.GetPort(portId.(string))
			if err != nil {
				return fmt.Errorf("ERROR: Cannot get port `%s`: %s", portId, err)
			}
			log.Printf("Port `%s` will be Attached", port.ID)
			if err := router.ConnectPort(port, true); err != nil {
				return fmt.Errorf("ERROR: Cannot attach port `%s`: %s", port.ID, err)
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
			return fmt.Errorf("ERROR: Can't update Router: %s", err)
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
			return nil, fmt.Errorf("ports: Error getting Ports from vdc")
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
			return nil, errors.New("ports: Error Ports should belong to routers vdc")
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
