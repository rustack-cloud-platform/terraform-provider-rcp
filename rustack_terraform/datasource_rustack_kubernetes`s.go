package rustack_terraform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/hashstructure/v2"
)

func dataSourceRustackKubernetess() *schema.Resource {
	args := Defaults()
	args.injectContextVdcById()
	args.injectResultListKubernetes()

	return &schema.Resource{
		ReadContext: dataSourceRustackKubernetessRead,
		Schema:      args,
	}
}

func dataSourceRustackKubernetessRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting vdc: %s", err)
	}

	allk8s, err := targetVdc.GetKubernetes()
	if err != nil {
		return diag.Errorf("Error retrieving Kubernetess: %s", err)
	}

	flattenedRecords := make([]map[string]interface{}, len(allk8s))
	for i, targetKubernetes := range allk8s {

		dashboard, err := targetKubernetes.GetKubernetesDashBoardUrl()
		if err != nil {
			return diag.Errorf("id: Error getting Kubernetes dashboard url: %s", err)
		}
		dashboard_url := fmt.Sprint(manager.BaseURL, *dashboard.DashBoardUrl)
		err = targetKubernetes.GetKubernetesConfigUrl()
		if err != nil {
			return diag.Errorf("id: Error creating Kubernetes config file url: %s", err)
		}

		err = targetKubernetes.GetKubernetesConfigUrl()
		if err != nil {
			return diag.Errorf("id: Error creating Kubernetes config file url: %s", err)
		}

		vms := make([]*string, len(targetKubernetes.Vms))
		for i, vm := range targetKubernetes.Vms {
			vms[i] = &vm.ID
		}
		flattenedRecords[i] = map[string]interface{}{
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
			flattenedRecords[i]["floating"] = true
			flattenedRecords[i]["floating_ip"] = targetKubernetes.Floating.IpAddress
		}
	}

	hash, err := hashstructure.Hash(allk8s, hashstructure.FormatV2, nil)
	if err != nil {
		diag.Errorf("unable to set `kubernetess` attribute: %s", err)
	}

	d.SetId(fmt.Sprintf("kubernetess/%d", hash))

	if err := d.Set("kubernetess", flattenedRecords); err != nil {
		return diag.Errorf("unable to set `kubernetess` attribute: %s", err)
	}

	return nil
}
