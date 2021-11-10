---
page_title: "rustack_disks Data Source - terraform-provider-rustack"
---
# rustack_disks (Data Source)

Get information about list of Disks in the Vdc for use in other resources.

Note: You can use the [`rustack_storage_profile`](Disk) data source to obtain metadata
about a single Disk if you already know the `name` and `vdc_id` to retrieve.

## Example Usage

```hcl

data "rustack_project" "single_project" {
    name = "Terraform Project"
}

data "rustack_vdc" "single_vdc" {
    project_id = "${data.rustack_project.single_project.id}"
    name = "Terraform VDC"
}

data "rustack_disks" "all_disks" {
    vdc_id = data.rustack_vdc.single_vdc.id
}

```

## Schema

### Required

- **vdc_id** (String) id of the VDC

### Read-Only

- **disks** (List of Object) (see [below for nested schema](#nestedatt--disks))

<a id="nestedatt--disks"></a>
### Nested Schema for `disks`

Read-Only:

- **id** (String)
- **name** (String)
- **size** (Number)
- **storage_profile_id** (String)
- **storage_profile_name** (String)
