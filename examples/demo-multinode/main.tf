terraform {
    required_version = ">= 1.0.0"

    required_providers {
        rustack = {
            source  = "pilat/rustack"
            version = "~> 0.1"
        }
    }
}

provider "rustack" {
    api_endpoint = var.rustack_endpoint
    token = var.rustack_token
}

provider "random" {}

provider "template" {}

resource "rustack_project" "demo_project" {
    name = "Terraform Demo"
}

data "rustack_hypervisor" "single_hypervisor" {
    project_id = rustack_project.demo_project.id
    name = "VMWARE"
}

resource "rustack_vdc" "vdc1" {
    name = "Terraform Demo VDC"
    project_id = "${rustack_project.demo_project.id}"
    hypervisor_id = data.rustack_hypervisor.single_hypervisor.id
}


data "rustack_firewall_template" "allow_default" {
    vdc_id = rustack_vdc.vdc1.id
    name = "По-умолчанию"
}

data "rustack_firewall_template" "allow_web" {
    vdc_id = rustack_vdc.vdc1.id
    name = "Разрешить WEB"
}

data "rustack_firewall_template" "allow_ssh" {
    vdc_id = rustack_vdc.vdc1.id
    name = "Разрешить SSH"
}


///
data "rustack_storage_profile" "ssd" {
    vdc_id = rustack_vdc.vdc1.id
    name = "ssd"
}

data "rustack_storage_profile" "sas" {
    vdc_id = rustack_vdc.vdc1.id
    name = "sas"
}


///////
data "rustack_network" "service_network" {
    vdc_id = rustack_vdc.vdc1.id
    name = "Сеть"
}


////
data "rustack_template" "debian10" {
    vdc_id = rustack_vdc.vdc1.id
    name = "Debian 10"
}


resource "random_password" "password" {
    length           = 16
    special          = true
    override_special = "_-#"
}

# locals {
#   expanded_names = [
#       for i in range(var.nodes_count) : format("host-%s", i)
#   ]
# }


data "template_file" "cloud_init_node" {
    count = var.nodes_count

    template = file("cloud-config-node.tpl")
    vars = {
        vm_password   = random_password.password.result
        hostname      = format("host-%s", count.index)
    }
}


resource "rustack_vm" "vm_node" {
    vdc_id = rustack_vdc.vdc1.id

    name = "Host 1"
    cpu = 2
    ram = 4

    template_id = data.rustack_template.debian10.id

    user_data = "${element(data.template_file.cloud_init_node.*.rendered, 1)}" 

    disk {
        name = "Root disk"
        size = 10
        storage_profile_id = data.rustack_storage_profile.ssd.id
    }

    port {
        network_id = data.rustack_network.service_network.id
        firewall_templates = [
            "${data.rustack_firewall_template.allow_default.id}"
        ]
    }

    floating = false
}


data "template_file" "cloud_init_master" {
    template = file("cloud-config-balancer.tpl")
    vars = {
        vm_password   = random_password.password.result
        hostname      = "balancer"
        balancer_upstream = <<-EOT
            %{ for k, v in rustack_vm.vm_node ~}
            server ${v.port.0.ip_address}:80;
            %{ endfor ~}
        EOT
    }
}

resource "rustack_vm" "vm_master" {
    vdc_id = rustack_vdc.vdc1.id

    name = "Master"
    cpu = 2
    ram = 4

    template_id = data.rustack_template.debian10.id

    user_data = data.template_file.cloud_init_master.rendered

    disk {
        name = "Root disk"
        size = 10
        storage_profile_id = data.rustack_storage_profile.ssd.id
    }

    port {
        network_id = data.rustack_network.service_network.id
        firewall_templates = [
            "${data.rustack_firewall_template.allow_default.id}",
            "${data.rustack_firewall_template.allow_web.id}",
            "${data.rustack_firewall_template.allow_ssh.id}"
        ]
    }

    floating = true
}

output "instance_ip" {
  value = rustack_vm.vm_master.floating_ip
}

output "instance_password" {
  value = nonsensitive(random_password.password.result)
}


