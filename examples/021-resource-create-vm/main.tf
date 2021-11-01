terraform {
  required_version = ">= 1.0.0"

  required_providers {
    rustack = {
      source  = "rustack/rustack"
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
    project_id = "${data.rustack_project.single_project.id}"
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

resource "rustack_vm" "vm1" {
    count = 2

    vdc_id = data.rustack_vdc.single_vdc.id

    name = format("Сервер %s", count.index + 1)
    cpu = 2
    ram = 4

    template_id = data.rustack_template.debian10.id

    user_data = "${file("user_data.yaml")}"

    disk {
        name = "Загрузочный диск"
        size = 10
        storage_profile_id = data.rustack_storage_profile.ssd.id
    }

    # disk {
    #     name = "Диск 2"
    #     size = 1
    #     storage_profile_id = data.rustack_storage_profile.sas.id
    # }

    port {
        network_id = data.rustack_network.service_network.id
        firewall_templates = ["${data.rustack_firewall_template.allow_default.id}",
            "${data.rustack_firewall_template.allow_web.id}",
            "${data.rustack_firewall_template.allow_ssh.id}"
        ]
    }

    floating = false
}
