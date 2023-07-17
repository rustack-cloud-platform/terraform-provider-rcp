package rustack_terraform

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pilat/rustack-go/rustack"
)

func resourceRustackLbaasPool() *schema.Resource {
	args := Defaults()
	args.injectContextLbaasByID()
	args.injectCreateLbaasPool()

	return &schema.Resource{
		CreateContext: resourceRustackLbaasPoolCreate,
		ReadContext:   resourceRustackLbaasPoolRead,
		UpdateContext: resourceRustackLbaasPoolUpdate,
		DeleteContext: resourceRustackLbaasPoolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: args,
	}
}

func resourceRustackLbaasPoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()

	lbaasId := d.Get("lbaas_id").(string)

	lbaas, err := manager.GetLoadBalancer(lbaasId)
	if err != nil {
		return diag.Errorf("id: Error getting Lbaas: %s", err)
	}
	// Get members
	membersCount := d.Get("member.#").(int)
	members := make([]*rustack.PoolMember, membersCount)

	for i := 0; i < membersCount; i++ {
		memberPrefix := fmt.Sprint("member.", i)
		member := d.Get(memberPrefix).(map[string]interface{})
		vm_id := member["vm_id"].(string)
		port := member["port"].(int)
		weight := member["weight"].(int)

		vm, err := manager.GetVm(vm_id)
		if err != nil {
			return diag.Errorf("vm_id: Error getting vm: %s", err)
		}

		newMember := rustack.NewLoadBalancerPoolMember(port, weight, vm)
		if err != nil {
			return diag.FromErr(err)
		}
		members[i] = &newMember
	}

	newPool := rustack.NewLoadBalancerPool(
		*lbaas,
		d.Get("port").(int),
		d.Get("connlimit").(int),
		members,
		d.Get("method").(string),
		d.Get("protocol").(string),
		d.Get("session_persistence").(string),
	)
	err = lbaas.CreatePool(&newPool)
	if err != nil {
		return diag.Errorf("id: Error creating Lbaas pool: %s", err)
	}
	lbaas.WaitLock()
	d.SetId(newPool.ID)
	return resourceRustackLbaasPoolRead(ctx, d, meta)
}

func resourceRustackLbaasPoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diagErr diag.Diagnostics) {
	manager := meta.(*CombinedConfig).rustackManager()
	lbaasPoolId := d.Id()
	lbaasId := d.Get("lbaas_id").(string)

	lbaas, err := manager.GetLoadBalancer(lbaasId)
	if err != nil {
		return diag.Errorf("id: Error getting Lbaas: %s", err)
	}

	pool, err := lbaas.GetLoadBalancerPool(lbaasPoolId)
	if err != nil {
		return diag.Errorf("Error getting LbaasPool: %s", err)
	}

	d.SetId(pool.ID)
	d.Set("port", pool.Port)
	d.Set("connlimit", pool.Connlimit)
	d.Set("method", pool.Method)
	d.Set("protocol", pool.Protocol)
	d.Set("session_persistence", pool.SessionPersistence)

	flattenedPools := make([]map[string]interface{}, len(pool.Members))
	for i, member := range pool.Members {
		flattenedPools[i] = map[string]interface{}{
			"id":     member.ID,
			"port":   member.Port,
			"weight": member.Weight,
			"vm":     member.Vm.ID,
		}
	}

	return
}

func resourceRustackLbaasPoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	lbaasPoolId := d.Id()
	lbaasId := d.Get("lbaas_id").(string)

	lbaas, err := manager.GetLoadBalancer(lbaasId)
	if err != nil {
		return diag.Errorf("id: Error getting Lbaas: %s", err)
	}

	pool, err := lbaas.GetLoadBalancerPool(lbaasPoolId)
	if err != nil {
		return diag.Errorf("Error getting LbaasPool: %s", err)
	}

	if d.HasChange("port") {
		pool.Port = d.Get("port").(int)
	}
	if d.HasChange("connlimit") {
		pool.Connlimit = d.Get("connlimit").(int)
	}
	if d.HasChange("method") {
		pool.Method = d.Get("method").(string)
	}
	if d.HasChange("protocol") {
		pool.Protocol = d.Get("protocol").(string)
	}
	if d.HasChange("session_persistence") {
		pool.SessionPersistence = d.Get("session_persistence").(*string)
	}
	if d.HasChange("member") {
		membersCount := d.Get("member.#").(int)
		members := make([]*rustack.PoolMember, membersCount)
		for i := 0; i < membersCount; i++ {
			memberPrefix := fmt.Sprint("member.", i)
			member := d.Get(memberPrefix).(map[string]interface{})
			vm_id := member["vm_id"].(string)
			port := member["port"].(int)
			weight := member["weight"].(int)

			vm, err := manager.GetVm(vm_id)
			if err != nil {
				return diag.Errorf("vm_id: Error getting vm: %s", err)
			}

			newMember := rustack.NewLoadBalancerPoolMember(port, weight, vm)
			if err != nil {
				return diag.FromErr(err)
			}
			members[i] = &newMember
		}
		pool.Members = members
	}
	err = lbaas.UpdatePool(&pool)
	if err != nil {
		return diag.Errorf("Error updating Lbaas pool: %s", err)
	}
	lbaas.WaitLock()

	return resourceRustackLbaasPoolRead(ctx, d, meta)
}

func resourceRustackLbaasPoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	lbaasPoolId := d.Id()
	lbaasId := d.Get("lbaas_id").(string)

	lbaas, err := manager.GetLoadBalancer(lbaasId)
	if err != nil {
		return diag.Errorf("id: Error getting Lbaas: %s", err)
	}

	_, err = lbaas.GetLoadBalancerPool(lbaasPoolId)
	if err != nil {
		return diag.Errorf("Error getting LbaasPool: %s", err)
	}

	lbaas.DeletePool(lbaasPoolId)
	if err != nil {
		return diag.Errorf("Error deleting LbaasPool: %s", err)
	}
	lbaas.WaitLock()

	d.SetId("")
	log.Printf("[INFO] LbaasPool deleted, ID: %s", lbaasId)

	return nil
}
