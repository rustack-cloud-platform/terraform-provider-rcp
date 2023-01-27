---
page_title: "rustack_lbaas Resource - terraform-provider-rustack"
---
# rustack_lbaas (Resource)

Provides a Rustack DNS record resource.

## Example Usage

```hcl
data "rustack_project" "single_project" {
    name = "Terraform Project"
}

data "rustack_vdc" "single_vdc" {
    project_id = data.rustack_project.single_project.id
    name = "Terraform VDC"
}

data "rustack_network" "new_network" {
    vdc_id =  data.rustack_vdc.single_vdc.id
    name = "New network"
}

resource "rustack_lbaas" "lbaas" {
    vdc_id = data.rustack_project.single_vdc.id
    name = "lbaas"
    port{
        network_id = data.rustack_network.new_network.id
    }
}

```

## Schema

### Required

- **vdc_id** (String) id of Vdc
- **name** (String) name of LoadBalancer
- **Port** (String) parameter that specifies which network will be connected to LoadBalancer  (see [below for nested schema](#nestedblock--port))


### Optional

- **floating** (Boolean) enable floating ip for the LoadBalancer.
- **timeouts** (Block, Optional)

<a id="nestedblock--port"></a>
### Nested Schema for `port`

Required:

- **network_id** (String) id of the Network

Optional:

- **ip_address** (String) ip address of port
