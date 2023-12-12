package rustack_terraform

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pilat/rustack-go/rustack"
)

func resourceRustackKubernetes() *schema.Resource {
	args := Defaults()
	args.injectCreateKubernetes()
	args.injectContextVdcById()
	args.injectContextKubernetesTemplateById() // override template_id

	return &schema.Resource{
		CreateContext: resourceRustackKubernetesCreate,
		ReadContext:   resourceRustackKubernetesRead,
		UpdateContext: resourceRustackKubernetesUpdate,
		DeleteContext: resourceRustackKubernetesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: args,
	}
}

func resourceRustackKubernetesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("vdc_id: Error getting VDC: %s", err)
	}

	platform_id := d.Get("platform").(string)
	if targetVdc.Hypervisor.Type == "Vmware" && platform_id == "" {
		return diag.Errorf("platform: This field is required for %s Hypervisor", targetVdc.Hypervisor.Type)
	}
	platform, err := manager.GetPlatform(platform_id)
	if err != nil {
		return diag.Errorf("template_id: Error getting template: %s", err)
	}
	template, err := GetKubernetesTemplateById(d, manager, targetVdc)
	if err != nil {
		return diag.Errorf("template_id: Error getting template: %s", err)
	}

	sp_id := d.Get("node_storage_profile_id").(string)
	storage_profile, err := targetVdc.GetStorageProfile(sp_id)
	if err != nil {
		return diag.Errorf("storage_profile_id: Error storage profile %s not found", sp_id)
	}

	userPublicKey := d.Get("user_public_key_id").(string)
	pub_key, err := manager.GetPublicKey(userPublicKey)
	if err != nil {
		return diag.Errorf("storage_profile_id: Error storage profile %s not found", userPublicKey)
	}
	name := d.Get("name").(string)
	cpu := d.Get("node_cpu").(int)
	ram := d.Get("node_ram").(int)
	nodesCount := d.Get("nodes_count").(int)
	nodeDiskSize := d.Get("node_disk_size").(int)
	log.Printf(name, cpu, ram, template.Name)

	var floatingIp *string = nil
	if d.Get("floating").(bool) {
		floatingIpStr := "RANDOM_FIP"
		floatingIp = &floatingIpStr
	}

	newKubernetes := rustack.NewKubernetes(name, cpu, ram, nodesCount, nodeDiskSize, floatingIp, template, storage_profile, pub_key.ID, platform)
	newKubernetes.Tags = unmarshalTagNames(d.Get("tags"))

	err = targetVdc.CreateKubernetes(&newKubernetes)
	if err != nil {
		return diag.Errorf("Error creating Kubernetes: %s", err)
	}

	newKubernetes.WaitLock()

	d.SetId(newKubernetes.ID)

	log.Printf("[INFO] Kubernetes created, ID: %s", d.Id())

	return resourceRustackKubernetesRead(ctx, d, meta)
}

func resourceRustackKubernetesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diagErr diag.Diagnostics) {
	manager := meta.(*CombinedConfig).rustackManager()
	Kubernetes, err := manager.GetKubernetes(d.Id())
	if err != nil {
		if err.(*rustack.RustackApiError).Code() == 404 {
			d.SetId("")
			return nil
		} else {
			return diag.Errorf("id: Error getting Kubernetes: %s", err)
		}
	}

	d.SetId(Kubernetes.ID)
	d.Set("name", Kubernetes.Name)
	d.Set("node_cpu", Kubernetes.NodeCpu)
	d.Set("node_ram", Kubernetes.NodeRam)
	d.Set("nodes_count", Kubernetes.NodesCount)
	d.Set("node_disk_size", Kubernetes.NodeDiskSize)
	d.Set("platform", Kubernetes.NodePlatform.ID)
	d.Set("template_id", Kubernetes.Template.ID)
	d.Set("tags", marshalTagNames(Kubernetes.Tags))

	vms := make([]*string, len(Kubernetes.Vms))
	for i, vm := range Kubernetes.Vms {
		vms[i] = &vm.ID
	}
	d.Set("vms", vms)

	d.Set("floating", Kubernetes.Floating != nil)
	d.Set("floating_ip", "")
	if Kubernetes.Floating != nil {
		d.Set("floating_ip", Kubernetes.Floating.IpAddress)
	}

	err = Kubernetes.GetKubernetesConfigUrl()
	if err != nil {
		diagErr = diag.Errorf("config: Error getting Kubernetes config: %s", err)
		return
	}

	dashboard, err := Kubernetes.GetKubernetesDashBoardUrl()
	if err != nil {
		diagErr = diag.Errorf("dashboard_url: Error getting Kubernetes dashboard url: %s", err)
		return
	}
	dashboard_url := fmt.Sprint(manager.BaseURL, *dashboard.DashBoardUrl)
	d.Set("dashboard_url", dashboard_url)

	return
}

func resourceRustackKubernetesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("vdc_id: Error getting VDC: %s", err)
	}

	needUpdate := false

	kubernetes, err := manager.GetKubernetes(d.Id())
	if err != nil {
		return diag.Errorf("id: Error getting Kubernetes: %s", err)
	}

	// Detect Kubernetes changes
	if d.HasChange("name") {
		needUpdate = true
		kubernetes.Name = d.Get("name").(string)
	}
	if d.HasChange("tags") {
		needUpdate = true
		kubernetes.Tags = unmarshalTagNames(d.Get("tags"))
	}
	needUpdate = true
	sp_id := d.Get("node_storage_profile_id").(string)
	storage_profile, err := targetVdc.GetStorageProfile(sp_id)
	if err != nil {
		return diag.Errorf("storage_profile_id: Error storage profile %s not found", sp_id)
	}

	userPublicKey := d.Get("user_public_key_id").(string)
	pub_key, err := manager.GetPublicKey(userPublicKey)
	if err != nil {
		return diag.Errorf("storage_profile_id: Error storage profile %s not found", userPublicKey)
	}
	kubernetes.NodeRam = d.Get("node_cpu").(int)
	kubernetes.NodeCpu = d.Get("node_ram").(int)
	kubernetes.UserPublicKey = pub_key.ID
	kubernetes.NodeStorageProfile = storage_profile
	kubernetes.NodeDiskSize = d.Get("node_disk_size").(int)
	kubernetes.NodesCount = d.Get("nodes_count").(int)

	if d.HasChange("floating") {
		needUpdate = true
		if !d.Get("floating").(bool) {
			kubernetes.Floating = &rustack.Port{IpAddress: nil}
		} else {
			kubernetes.Floating = &rustack.Port{ID: "RANDOM_FIP"}
		}
		d.Set("floating", kubernetes.Floating != nil)
	}

	if needUpdate {
		if err := repeatOnError(kubernetes.Update, kubernetes); err != nil {
			return diag.Errorf("Error updating Kubernetes: %s", err)
		}
	}

	return resourceRustackKubernetesRead(ctx, d, meta)
}

func resourceRustackKubernetesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	kubernetes, err := manager.GetKubernetes(d.Id())
	if err != nil {
		return diag.Errorf("id: Error getting Kubernetes: %s", err)
	}

	err = kubernetes.Delete()
	if err != nil {
		return diag.Errorf("Error deleting Kubernetes: %s", err)
	}
	kubernetes.WaitLock()

	return nil
}
