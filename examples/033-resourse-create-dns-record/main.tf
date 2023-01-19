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

data "rustack_dns" "dns" {
    name = "test.test."
    project_id = data.rustack_project.single_project.id
}

resource "rustack_dns_record" "dns_record1" {
    dns_id = data.rustack_dns.dns.id
    type = "A"
    host = "test.test.test."
    data = "8.8.8.8"
}