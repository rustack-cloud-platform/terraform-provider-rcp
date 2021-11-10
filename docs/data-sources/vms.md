---
page_title: "rustack_vms Data Source - terraform-provider-rustack"
---
# rustack_vms (Data Source)

Returns a list of Rustack vms.

Get information about Vms in the Vdc for use in other resources.

Note: You can use the [`rustack_vm`](Vm) data source to obtain metadata
about a single Vm if you already know the `name` and `vdc_id` to retrieve.

## Example Usage

```hcl

data "rustack_project" "single_project" {
    name = "Terraform Project"
}

data "rustack_vdc" "single_vdc" {
    project_id = "${data.rustack_project.single_project.id}"
    name = "Terraform VDC"
}

data "rustack_vms" "all_vms" {
    vdc_id = data.rustack_vdc.single_vdc.id
}

```

## Schema

### Required

- **vdc_id** (String) id of the VDC

### Read-Only

- **vms** (List of Object) (see [below for nested schema](#nestedatt--vms))

<a id="nestedatt--vms"></a>
### Nested Schema for `vms`

Read-Only:

- **cpu** (Number)
- **floating** (Boolean)
- **floating_ip** (String)
- **id** (String)
- **name** (String)
- **ram** (Number)
- **template_id** (String)
- **template_name** (String)
