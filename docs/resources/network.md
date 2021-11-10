---
page_title: "rustack_network Resource - terraform-provider-rustack"
---
# rustack_network (Resource)

Provides a Rustack network to provide connections of two or more computers that are linked in order to share resources.

## Example Usage

```hcl
data "rustack_project" "single_project" {
    name = "Terraform Project"
}

data "rustack_vdc" "single_vdc" {
    project_id = "${data.rustack_project.single_project.id}"
    name = "Terraform VDC"
}

resource "rustack_network" "network2" {
    count = 1

    vdc_id = data.rustack_vdc.single_vdc.id

    name = format("Сеть %s", count.index + 1)

    subnets {
        cidr = format("10.20.%s.0/24", (count.index + 3) * 10)
        dhcp = true
        gateway = format("10.20.%s.1", (count.index + 3) * 10)
        start_ip = format("10.20.%s.2", (count.index + 3) * 10)
        end_ip = format("10.20.%s.254", (count.index + 3) * 10)
        dns = ["8.8.8.8", "8.8.4.4", "1.1.1.1"]
    }
}
```

## Schema

### Required

- **name** (String) name of the Network
- **subnets** (Block List, Min: 1, Max: 1) (see [below for nested schema](#nestedblock--subnets))
- **vdc_id** (String) id of the VDC

### Optional

- **id** (String) The ID of this resource.

<a id="nestedblock--subnets"></a>
### Nested Schema for `subnets`

Required:

- **cidr** (String) cidr of the Subnet
- **dhcp** (Boolean) enable dhcp service of the Subnet
- **dns** (List of String) dns servers list
- **end_ip** (String) pool end ip of the Subnet
- **gateway** (String) gateway of the Subnet
- **start_ip** (String) pool start ip of the Subnet

Read-Only:

- **id** (String) id of the Subnet
