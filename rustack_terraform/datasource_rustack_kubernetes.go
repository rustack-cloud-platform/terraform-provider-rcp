package rustack_terraform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pilat/rustack-go/rustack"
)

func dataSourceRustackKubernetes() *schema.Resource {
	args := Defaults()
	args.injectResultKubernetes()
	args.injectContextVdcById()
	args.injectContextGetKubernetes() // override "name"

	return &schema.Resource{
		ReadContext: dataSourceRustackKubernetesRead,
		Schema:      args,
	}
}

func dataSourceRustackKubernetesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting vdc: %s", err)
	}

	target, err := checkDatasourceNameOrId(d)
	if err != nil {
		return diag.Errorf("Error getting Kubernetes: %s", err)
	}
	var targetKubernetes *rustack.Kubernetes
	if target == "id" {
		targetKubernetes, err = manager.GetKubernetes(d.Get("id").(string))
		if err != nil {
			return diag.Errorf("Error getting Kubernetes: %s", err)
		}
	} else {
		targetKubernetes, err = GetKubernetesByName(d, manager, targetVdc)
		if err != nil {
			return diag.Errorf("Error getting Kubernetes: %s", err)
		}
	}

	vms := make([]*string, len(targetKubernetes.Vms))
	for i, vm := range targetKubernetes.Vms {
		vms[i] = &vm.ID
	}

	dashboard, err := targetKubernetes.GetKubernetesDashBoardUrl()
	if err != nil {
		return diag.Errorf("id: Error getting Kubernetes dashboard url: %s", err)
	}
	dashboard_url := fmt.Sprint(manager.BaseURL, *dashboard.DashBoardUrl)

	flatten := map[string]interface{}{
		"id":                      targetKubernetes.ID,
		"name":                    targetKubernetes.Name,
		"node_cpu":                targetKubernetes.NodeCpu,
		"node_ram":                targetKubernetes.NodeRam,
		"template_id":             targetKubernetes.Template.ID,
		"node_disk_size":          targetKubernetes.NodeDiskSize,
		"nodes_count":             targetKubernetes.NodesCount,
		"user_public_key_id":      targetKubernetes.UserPublicKey,
		"node_storage_profile_id": targetKubernetes.NodeStorageProfile.ID,
		"floating":                nil,
		"floating_ip":             nil,
		"vms":                     vms,
		"dashboard_url":           dashboard_url,
	}

	if targetKubernetes.Floating != nil {
		flatten["floating"] = true
		flatten["floating_ip"] = targetKubernetes.Floating.IpAddress
	}

	if err := setResourceDataFromMap(d, flatten); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(targetKubernetes.ID)
	return nil
}
