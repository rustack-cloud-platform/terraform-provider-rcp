package rustack_terraform

import (
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
		},

		ResourcesMap: map[string]*schema.Resource{
			"rustack_project":           resourceRustackProject(),          // 001-resource-create-project +
			"rustack_vdc":               resourceRustackVdc(),              // 006-resource-create-vdc +
			"rustack_network":           resourceRustackNetwork(),          // 009-resource-create-network +
			"rustack_disk":              resourceRustackDisk(),             // 014-resource-create-disk +
			"rustack_vm":                resourceRustackVm(),               // 021-resource-create-vm +
			"rustack_firewall_template": resourceRustackFirewallTemplate(), // 024-resource-create-firewall-template +
		},
	}

	p.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
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

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	config := Config{
		Token:            d.Get("token").(string),
		APIEndpoint:      d.Get("api_endpoint").(string),
		TerraformVersion: terraformVersion,
	}

	return config.Client()
}
