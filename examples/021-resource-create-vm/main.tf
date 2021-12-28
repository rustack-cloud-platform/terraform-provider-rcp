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
    token = "[PLACE_YOUR_TOKEN_HERE]"
}

data "rustack_project" "single_project" {
    name = "Terraform Project"
}

data "rustack_vdc" "single_vdc" {
    project_id = data.rustack_project.single_project.id
    name = "Terraform VDC"
}

data "rustack_network" "service_network" {
    vdc_id = data.rustack_vdc.single_vdc.id
    name = "Сеть"
}

data "rustack_storage_profile" "ssd" {
    vdc_id = data.rustack_vdc.single_vdc.id
    name = "ssd"
}

data "rustack_storage_profile" "sas" {
    vdc_id = data.rustack_vdc.single_vdc.id
    name = "sas"
}

data "rustack_template" "debian10" {
    vdc_id = data.rustack_vdc.single_vdc.id
    name = "Debian 10"
}

data "rustack_firewall_template" "allow_default" {
    vdc_id = data.rustack_vdc.single_vdc.id
    name = "По-умолчанию"
}

data "rustack_firewall_template" "allow_web" {
    vdc_id = data.rustack_vdc.single_vdc.id
    name = "Разрешить WEB"
}

data "rustack_firewall_template" "allow_ssh" {
    vdc_id = data.rustack_vdc.single_vdc.id
    name = "Разрешить SSH"
}

data "rustack_disk" "new_disk1" {
    vdc_id = data.rustack_vdc.single_vdc.id
    name = "Disk 1"
}

data "rustack_disk" "new_disk2" {
    vdc_id = data.rustack_vdc.single_vdc.id
    name = "Disk 2"
}

resource "rustack_vm" "vm1" {
    vdc_id = data.rustack_vdc.single_vdc.id

    name = "Сервер 1"
    cpu = 2
    ram = 4

    template_id = data.rustack_template.debian10.id

    user_data = "${file("user_data.yaml")}"

    system_disk {
        size = 10
        storage_profile_id = data.rustack_storage_profile.ssd.id
    }
    
    disks = [
        data.rustack_disk.new_disk1.id,
        data.rustack_disk.new_disk2.id,
    ]

    port {
        network_id = data.rustack_network.service_network.id
        firewall_templates = [data.rustack_firewall_template.allow_default.id,
            data.rustack_firewall_template.allow_web.id,
            data.rustack_firewall_template.allow_ssh.id
        ]
    }

    floating = false
}
