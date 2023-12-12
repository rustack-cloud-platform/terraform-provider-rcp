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

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- **create** (String)
- **delete** (String)
