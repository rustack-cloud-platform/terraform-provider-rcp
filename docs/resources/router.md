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

data "rustack_port" "vm_port" {
    vdc_id = resource.rustack_vdc.single_vdc.id

    network_id = resource.rustack_network.default_network.id
}


resource "rustack_router" "new_router" {
  vdc_id =  data.rustack_vdc.single_vdc.id
  name = "New router"
  ports = [
    data.rustack_port.vm_port.id,
  ]
  floating = false
  tags = ["created_by:terraform"]
}

```

## Schema

### Required

- **name** (String) name of the Network
- **ports** (Toset, String) list of Ports id attached to the Router.
- **vdc_id** (String) id of the VDC

### Optional

- **system** (Bool) let terraform treat system router properly. False by default. There can be only 1 router with the system = ture
- **floating** (Bool) enable floating ip for the Router. True by default.
- **is_default** (Bool) Set up this option to set router by default.
- **tags** (Toset, String) list of Tags added to the Router

Read-Only:

- **id** (String) id of the Subnet
- **floating_id** (String) id of the Floating address
