package rustack_terraform

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RUSTACK_TOKEN", nil),
				Description: "The token key for API operations.",
			},
			"api_endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("RUSTACK_API_URL", "https://cp.sbcloud.ru"),
				Description: "The URL to use for the Rustack API.",
			},
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RUSTACK_CLIENT_ID", nil),
				Description: "The client id to use for managing instances.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"rustack_account": dataSourceRustackAccount(),

			"rustack_project":            dataSourceRustackProject(),           // 002-data-get-project +
			"rustack_projects":           dataSourceRustackProjects(),          // 003-data-get-projects +
			"rustack_hypervisor":         dataSourceRustackHypervisor(),        // 004-data-get-hypervisor +
			"rustack_hypervisors":        dataSourceRustackHypervisors(),       // 005-data-get-hypervisors +
			"rustack_vdc":                dataSourceRustackVdc(),               // 007-data-get-vdc +
			"rustack_vdcs":               dataSourceRustackVdcs(),              // 008-data-get-vdcs +
			"rustack_network":            dataSourceRustackNetwork(),           // 010-data-get-network +
			"rustack_networks":           dataSourceRustackNetworks(),          // 011-data-get-networks +
			"rustack_storage_profile":    dataSourceRustackStorageProfile(),    // 012-data-get-storage-profile +
			"rustack_storage_profiles":   dataSourceRustackStorageProfiles(),   // 013-data-get-storage-profiles +
			"rustack_disk":               dataSourceRustackDisk(),              // 015-data-get-disk +
			"rustack_disks":              dataSourceRustackDisks(),             // 016-data-get-disks +
			"rustack_template":           dataSourceRustackTemplate(),          // 017-data-get-template +
			"rustack_templates":          dataSourceRustackTemplates(),         // 018-data-get-templates +
			"rustack_firewall_template":  dataSourceRustackFirewallTemplate(),  // 019-data-get-template +
			"rustack_firewall_templates": dataSourceRustackFirewallTemplates(), // 020-data-get-templates +
			"rustack_vm":                 dataSourceRustackVm(),                // 022-data-get-vm
			"rustack_vms":                dataSourceRustackVms(),               // 023-data-get-vms
			"rustack_router":             dataSourceRustackRouter(),            // 025-data-get-router +
			"rustack_routers":            dataSourceRustackRouters(),           // 026-data-get-routers +
			"rustack_port":               dataSourceRustackPort(),              // 027-data-get-port +
			"rustack_ports":              dataSourceRustackPorts(),             // 027-data-get-ports +
			"rustack_dns":                dataSourceRustackDns(),               // 028-data-get-dns +
			"rustack_dnss":               dataSourceRustackDnss(),              // 028-data-get-dnss +
			"rustack_lbaas":              dataSourceRustackLbaas(),             // 028-data-get-lbaas +
			"rustack_lbaass":             dataSourceRustackLoadBalancers(),     // 028-data-get-lbaass +
			"rustack_s3":                 dataSourceRustackLoadBalancers(),     // 028-data-get-lbaass +
			"rustack_s3_storage":         dataSourceRustackS3Storage(),     	// 028-data-get-s3-storage +
			"rustack_s3_storages":        dataSourceRustackS3Storages(),     	// 028-data-get-s3-storages +
		},

		ResourcesMap: map[string]*schema.Resource{
			"rustack_project":                resourceRustackProject(),          // 001-resource-create-project +
			"rustack_vdc":                    resourceRustackVdc(),              // 006-resource-create-vdc +
			"rustack_network":                resourceRustackNetwork(),          // 009-resource-create-network +
			"rustack_disk":                   resourceRustackDisk(),             // 014-resource-create-disk +
			"rustack_vm":                     resourceRustackVm(),               // 021-resource-create-vm +
			"rustack_firewall_template":      resourceRustackFirewallTemplate(), // 024-resource-create-firewall-template +
			"rustack_router":                 resourceRustackRouter(),           // 027-resource-create-router +
			"rustack_port":                   resourceRustackPort(),             // 027-resource-create-port +
			"rustack_dns":                    resourceRustackDns(),              // 028-resource-create-dns +
			"rustack_dns_record":             resourceRustackDnsRecord(),        // 028-resource-create-dns-record +
			"rustack_firewall_template_rule": resourceRustackFirewallRule(),     // 029-resource-create-firewall-rule +
			"rustack_lbaas":                  resourceRustackLbaas(),            // 029-resource-create-lbaas +
			"rustack_lbaas_pool":             resourceRustackLbaasPool(),        // 029-resource-create-lbaas-pool +
			"rustack_s3_storage":             resourceRustackS3Storage(),        // 029-resource-create-s3-storage +
			"rustack_s3_storage_bucket":      resourceRustackS3StorageBucket(),  // 029-resource-create-s3-storage-bucket +
		},
	}

	p.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		terraformVersion := p.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(d, terraformVersion)
	}

	return p
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, diag.Diagnostics) {
	config := Config{
		Token:            d.Get("token").(string),
		APIEndpoint:      d.Get("api_endpoint").(string),
		ClientID:         d.Get("client_id").(string),
		TerraformVersion: terraformVersion,
	}

	return config.Client()
}
