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

resource "rustack_project" "project" {
    name = "Terraform GitLab Demo"
}

data "rustack_hypervisor" "vmware" {
    project_id = rustack_project.project.id
    name = "vmware"
}

resource "rustack_vdc" "vdc" {
    name = "Gitlab"
    project_id = "${rustack_project.project.id}"
    hypervisor_id = data.rustack_hypervisor.vmware.id
}


data "rustack_firewall_template" "allow_default" {
    vdc_id = rustack_vdc.vdc.id
    name = "По-умолчанию"
}

data "rustack_firewall_template" "allow_web" {
    vdc_id = rustack_vdc.vdc.id
    name = "Разрешить WEB"
}

data "rustack_firewall_template" "allow_ssh" {
    vdc_id = rustack_vdc.vdc.id
    name = "Разрешить SSH"
}

data "rustack_storage_profile" "ssd" {
    vdc_id = rustack_vdc.vdc.id
    name = "ssd"
}

data "rustack_network" "service_network" {
    vdc_id = rustack_vdc.vdc.id
    name = "Сеть"
}

data "rustack_template" "ubuntu20" {
    vdc_id = rustack_vdc.vdc.id
    name = "Ubuntu 20.04"
}

resource "random_password" "password" {
    length           = 16
    special          = true
    override_special = "_-#"
}

data "template_file" "cloud_init" {
    template = file("cloud-config.tpl")
    vars = {
        user_login        = var.user_login
        public_key        = file(var.public_key)
        hostname          = "gitlab"
        gitlab_password   = random_password.password.result
    }
}

resource "rustack_vm" "gitlab" {
    vdc_id = rustack_vdc.vdc.id

    name = "GitLab"
    cpu = 8
    ram = 16

    template_id = data.rustack_template.ubuntu20.id

    user_data = data.template_file.cloud_init.rendered

    disk {
        name = "Root"
        size = 50
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

output "gitlab_ip" {
  value = rustack_vm.gitlab.floating_ip
}

output "gitlab_user" {
  value = "root"
}

output "gitlab_password" {
  value = nonsensitive(random_password.password.result)
}
