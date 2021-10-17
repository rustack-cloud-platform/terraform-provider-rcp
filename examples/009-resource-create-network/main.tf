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

resource "rustack_network" "network2" {
    count = 1

    vdc_id = data.rustack_vdc.single_vdc.id

    name = format("Сеть %s", count.index + 1)

    subnets {
        cidr = format("10.20.%s.0/24", (count.index + 3) * 10)
        dhcp = true
        gateway = format("10.20.%s.1", (count.index + 3) * 10)
        start_ip = format("10.20.%s.2", (count.index + 3) * 10)
        end_ip = format("10.20.%s.254", (count.index + 3) * 10)
        dns = ["8.8.8.8", "8.8.4.4", "1.1.1.1"]
    }
}
