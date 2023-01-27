---
page_title: "rustack_router Data Source - terraform-provider-rustack"
---
# rustack_router (Data Source)

Get information about a Routers for use in other resources. 

## Example Usage

```hcl

data "rustack_project" "single_project" {
    name = "Terraform Project"
}

data "rustack_vdc" "single_vdc" {
    project_id = data.rustack_project.single_project.id"
    name = "Terraform VDC"
}

data "rustack_routers" "vdc_routers" {
    vdc_id = data.rustack_vdc.single_vdc.id
}

```
## Schema

### Required

- **vdc_id** (String) id of the VDC

### Read-Only

- **routers** (List of Object) (see [below for nested schema](#nestedatt--router))

<a id="nestedatt--router"></a>
### Nested Schema for `router`

Read-Only:

- **id** (String)
- **name** (String)
