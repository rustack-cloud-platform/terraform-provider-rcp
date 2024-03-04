terraform {
  required_version = ">= 1.0.0"

  required_providers {
    rustack = {
      source  = "rustack-cloud-platform/rcp"
    }
  }
}

provider "rustack" {
  token = "[PLACE_YOUR_TOKEN_HERE]"
}

data "rustack_project" "single_project" {
  name = "Terraform Project"
}

data "rustack_account" "me"{}

data "rustack_hypervisor" "vmware" {
    project_id = resource.rustack_project.project1.id
    name = "VMware"
}

resource "rustack_vdc" "vdc" {
    name = "Vmware Terraform"
    project_id = resource.rustack_project.project1.id
    hypervisor_id = data.rustack_hypervisor.vmware.id
}

data "rustack_kubernetes_template" "kubernetes_template"{
    name = "Kubernetes 1.22.1"
    vdc_id = resource.rustack_vdc.vdc1.id
    
}

data "rustack_storage_profile" "ssd" {
    vdc_id = resource.rustack_vdc.vdc.id
    name = "ssd"
}

data "rustack_pub_key" "key"{
    name = "test"
    account_id = data.rustack_account.me.id
}

data "rustack_platform" "platform"{
    vdc_id = resource.rustack_vdc.vdc1.id
    name = "Intel Cascade Lake"
    
}

resource "rustack_kubernetes" "k8s" {
    vdc_id = resource.rustack_vdc.vdc1.id
    name = "kubernetes"
    node_ram = 3
    node_cpu = 3
    platform = data.rustack_platform.platform.id # vmware hypervosor only
    template_id = data.rustack_kubernetes_template.kuber.id
    nodes_count = 2
    node_disk_size = 10
    node_storage_profile_id = data.rustack_storage_profile.ssd.id
    user_public_key_id = data.rustack_pub_key.key.id
    floating = true
}

# For get dashboard url
output "dashboard_k8s" {
    value = resource.rustack_kubernetes.k8s.dashboard_url
}