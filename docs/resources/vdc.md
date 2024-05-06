---
page_title: "rustack_vdc Resource - terraform-provider-rustack"
---
# rustack_vdc (Resource)

Provides a Rustack VDC resource to determinate hipervisor to use.

## Example Usage

```hcl
data "rustack_project" "single_project" {
    name = "Terraform Project"
}

data "rustack_hypervisor" "single_hypervisor" {
    project_id = data.rustack_project.single_project.id
    name = "VMWARE"
}

resource "rustack_vdc" "vdc1" {
    name = "Terraform VDC"
    project_id = data.rustack_project.single_project.id
    hypervisor_id = data.rustack_hypervisor.single_hypervisor.id
}
```

## Schema

### Required

- **hypervisor_id** (String) id of the Hypervisor
- **name** (String) name of the VDC
- **project_id** (String) id of the Project

### Optional

- **id** (String) The ID of this resource.
- **timeouts** (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- **tags** (Toset, String) list of Tags added to the VDC.
- **default_network_mtu** (Integer) maximum transmission unit for the default network of the vdc

### Read-only

- **default_network_id** (String) id of the default network of the vdc
- **default_network_name** (String) name of the default network of the vdc
- **default_network_subnets** (Block List) (see [below for nested schema](#nestedblock--subnets))

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- **create** (String)
- **delete** (String)

<a id="nestedblock--subnets"></a>
### Nested Schema for `subnets`

Read-Only:

- **cidr** (String)
- **dhcp** (Boolean)
- **dns** (List of String)
- **end_ip** (String)
- **gateway** (String)
- **id** (String)
- **start_ip** (String)
