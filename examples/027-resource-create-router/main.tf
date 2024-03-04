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
  name       = "Terraform VDC"
}

data "rustack_network" "network1" {
    vdc_id = data.rustack_vdc.single_vdc.id
    name = "Network 1"
}

resource "rustack_port" "router_port" {
    vdc_id = resource.rustack_vdc.single_vdc.id
    network_id = resource.rustack_network.network1.id
}

resource "rustack_router" "new_router" {
  vdc_id = data.rustack_vdc.single_vdc.id
  name = "New router3"
  ports = [
    resource.rustack_port.router_port.id
  ]
  # floating = false

  # System router creating automatically with vdc, to manage it there is special flag "system". 
  # Terraform wont create new router, but read existing one so you can manage it like resource
  # On delete terraform will disconnect all networks except the default and return floating to default "true" value.

  # system = true 
}
