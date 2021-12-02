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

data "rustack_network" "network1" {
    vdc_id = rustack_vdc.single_vdc.id
    name = "Network 1"
}

data "rustack_network" "network2" {
    vdc_id = rustack_vdc.single_vdc.id
    name = "Network 2"
}

resource "rustack_router" "new_router" {
  vdc_id = rustack_vdc.single_vdc.id
  name = "New router3"
  networks = [
    data.rustack_network.network1.id,
    data.rustack_network.network2.id
  ]
  # floating = "10.11.133.8"
}
