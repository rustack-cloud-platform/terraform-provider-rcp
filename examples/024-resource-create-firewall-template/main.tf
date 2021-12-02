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
  name       = "Terraform VDC"
}

resource "rustack_firewall_template" "single_template" {
  vdc_id = data.rustack_vdc.single_vdc.id
  name   = "New custom template"
}

resource "rustack_firewall_template" "single_template" {
  vdc_id = rustack_vdc.single_vdc.id
  name   = "New custom template"

  ingress_rule {
    protocol       = "tcp"
    port_range     = "80"
    destination_ip = "2.0.0.0"
  }

  ingress_rule {
    protocol       = "icmp"
    destination_ip = "1.0.0.0/24"
  }

  egress_rule {
    protocol       = "udp"
    port_range     = "53"
    destination_ip = "5.0.0.0/24"
  }
}
