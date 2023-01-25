terraform {
  required_version = ">= 1.0.0"

  required_providers {
    rustack = {
      source  = "pilat/rustack"
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

resource "rustack_firewall_template_rule" "rule_1" {
    firewall_id = resource.rustack_firewall_template.single_template.id
    name = "test"
    direction = "ingress"
    protocol = "tcp"
    port_range = "80"
    destination_ip = "0.0.0.0/0"
}

resource "rustack_firewall_template_rule" "rule_2" {
    firewall_id = resource.rustack_firewall_template.single_template.id
    name = "test"
    direction = "egress"
    protocol = "tcp"
    destination_ip = "0.0.0.0/0"
}
