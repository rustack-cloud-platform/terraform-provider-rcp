---
page_title: "rustack_port Resource - terraform-provider-rustack"
---
# rustack_port (Resource)

Provides a Rustack port which can be attached to a VM and Router in order to provide connection with network.

## Example Usage

```hcl

data "rustack_project" "single_project" {
    name = "Terraform Project"
}

data "rustack_vdc" "single_vdc" {
    project_id = data.rustack_project.single_project.id
    name = "Terraform VDC"
}

data "rustack_firewall_template" "allow_default" {
    vdc_id = resource.rustack_vdc.vdc1.id
    name = "Разрешить входящие"
}


resource "rustack_network" "network" {
    vdc_id = resource.rustack_vdc.single_vdc.id
    name = "network"

    subnets {
        cidr = "10.20.3.0/24"
        dhcp = true
        gateway = "10.20.3.1"
        start_ip = "10.20.3.2"
        end_ip = "10.20.3.254"
        dns = ["8.8.8.8", "8.8.4.4", "1.1.1.1"]
    }
}

resource "rustack_port" "router_port" {
    vdc_id = resource.rustack_vdc.single_vdc.id

    network_id = resource.rustack_network.network.id
    ip_address = "199.199.199.199"
    firewall_templates = [data.rustack_firewall_template.allow_default.id]
}
```

## Schema

### Required

- **network_id** String) id of the Network
- **vdc_id** (String) id of the VDC

### Optional

- **firewall_templates** (List of String) list of firewall rule ids of the Port
- **ip_address** (String) must be accurate

### Read-Only

- **id** (String) id of the Port
