---
page_title: "rustack_router Resource - terraform-provider-rustack"
---
# rustack_router (Resource)

Provides a Rustack network to provide connections of two or more computers that are linked in order to share resources.

## Example Usage

```hcl
data "rustack_project" "single_project" {
    name = "Terraform Project"
}

data "rustack_vdc" "single_vdc" {
    project_id = data.rustack_project.single_project.id
    name = "Terraform VDC"
}

data "rustack_network" "default_network" {
  vdc_id =  data.rustack_vdc.single_vdc.id
  name = "Network"
}

data "rustack_network" "new_network" {
  vdc_id =  data.rustack_vdc.single_vdc.id
  name = "New network"
}

resource "rustack_router" "new_router" {
  vdc_id =  data.rustack_vdc.single_vdc.id
  name = "New router"
  networks = [
    data.rustack_network.new_network.id,
    data.rustack_network.default_network.id
  ]
  floating = false
}

```

## Schema

### Required

- **name** (String) name of the Network
- **networks** (Toset, String, Min: 1, Max: 20) List of network id.
- **vdc_id** (String) id of the VDC

### Optional

- **system** (Bool) let terraform treat system router properly. False by default.
- **floating** (Bool) enable floating ip for the Router. True by default.
- **is_default** (Bool) Set up this option to set router by default.

Read-Only:

- **id** (String) id of the Subnet
- **floating_id** (String) id of the Floating address
