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

data "rustack_network" "network1" {
    vdc_id = data.rustack_vdc.single_vdc.id
    name = "Network 1"
}

data "rustack_network" "network2" {
    vdc_id = data.rustack_vdc.single_vdc.id
    name = "Network 2"
}

resource "rustack_router" "new_router" {
  vdc_id = data.rustack_vdc.single_vdc.id
  name = "New router3"
  networks = [
    data.rustack_network.network1.id,
    data.rustack_network.network2.id
  ]
  # floating = false

  # System router creating automatically with vdc, to manage it there is special flag "system". 
  # Terraform wont create new router, but read existing one so you can manage it like resource
  # On delete terraform will disconnect all networks except the default and return floating to default "true" value.

  # system = true 
}
