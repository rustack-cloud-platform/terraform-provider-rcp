---
page_title: "rustack_platforms Data Source - terraform-provider-rustack"
---
# rustack_platforms (Data Source)
### `Only for Vmware Hypervisor`

Get information about Platforms in the Vdc for use in other resources.

Note: You can use the [`rustack_platforms`](Platforms) data source to obtain metadata
about a single Platforms if you already know the `name` to retrieve.

## Example Usage

```hcl
data "rustack_project" "single_project" {
    name = "Terraform Project"
}
data "rustack_vdc" "single_vdc" {
    project_id = data.rustack_project.single_project.id
    name = "Terraform VDC"
}
data "rustack_platforms" "platforms"{
    vdc_id = resource.rustack_vdc.single_vdc.id
}
```

## Schema

### Read-Only

- **platforms** (List of Object) (see [below for nested schema](#nestedatt--projects))

<a id="nestedatt--platforms"></a>
### Nested Schema for `platforms`

Read-Only:

- **id** (String)
- **name** (String)