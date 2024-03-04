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
    name = "Terraform VDC"
}

resource "rustack_network" "network2" {
    vdc_id = data.rustack_vdc.single_vdc.id

    name = "Сеть 1"

    subnets {
        cidr = "10.20.40.0/24"
        dhcp = true
        gateway = "10.20.40.1"
        start_ip = "10.20.40.2"
        end_ip = "10.20.40.254"
        dns = ["8.8.8.8", "8.8.4.4", "1.1.1.1"]
    }
}
