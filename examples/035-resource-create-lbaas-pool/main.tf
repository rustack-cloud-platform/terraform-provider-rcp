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

data "rustack_vdc" "single_vdc" {
  project_id = data.rustack_project.single_project.id
  name       = "Terraform VDC"
}

data "rustack_network" "service_network" {
    vdc_id = data.rustack_vdc.single_vdc.id
    name = "Сеть"
}

data "rustack_firewall_template" "allow_default" {
    vdc_id = resource.rustack_vdc.vdc1.id
    name = "Разрешить входящие"
}

data "rustack_storage_profile" "ssd" {
    vdc_id = resource.rustack_vdc.vdc1.id
    name = "ssd"
}


resource "rustack_port" "vm_port" {
    vdc_id = resource.rustack_vdc.vdc1.id

    network_id = resource.rustack_network.service_network.id
    firewall_templates = [data.rustack_firewall_template.allow_default.id]
}

resource "rustack_vm" "vm" {
    vdc_id = resource.rustack_vdc.vdc1.id
    name = "Server 1"
    cpu = 3
    ram = 3
    power = true

    template_id = data.rustack_template.ubuntu20.id

    user_data = "${file("user_data.yaml")}"

    system_disk {
        size = 10
        storage_profile_id = data.rustack_storage_profile.ssd.id
    }

    ports = [resource.rustack_port.vm_port.id]

    floating = true
}


resource "rustack_lbaas" "lbaas" {
    vdc_id = resource.rustack_vdc.single_vdc.id
    name = "lbaas"
    floating = true
    port {
        network_id = data.rustack_network.service_network.id
    }
}

resource "rustack_lbaas_pool" "pool" {
    lbaas_id = resource.rustack_lbaas.lbaas.id
    connlimit = 34
    method = "SOURCE_IP"
    protocol = "TCP"
    port = 80

     member {
        port = 80
        weight = 50
        vm_id = resource.rustack_vm.vm.id
    }

    depends_on = [rustack_vm.vm1]
}
