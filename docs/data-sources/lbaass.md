---
page_title: "rustack_lbaass Data Source - terraform-provider-rustack"
---
# rustack_lbaass (Data Source)

Returns a list of Rustack lbaass.

Get information about Lbaass in the Vdc for use in other resources.

Note: You can use the [`rustack_lbaas`](Lbaas) data source to obtain metadata
about a single lbaas if you already know the `name` and `vdc_id` to retrieve.

## Example Usage

```hcl

data "rustack_project" "single_project" {
    name = "Terraform Project"
}

data "rustack_vdc" "single_vdc" {
    project_id = data.rustack_project.single_project.id
    name = "Terraform VDC"
}

data "rustack_lbaass" "all_lbaass" {
    vdc_id = data.rustack_vdc.single_vdc.id
}

```

## Schema

### Required

- **vdc_id** (String) id of the VDC

### Read-Only

- **lbaas** (List of Object) (see [below for nested schema](#nestedatt--lbaas))

<a id="nestedatt--lbaas"></a>
### Nested Schema for `lbaas`

Read-Only:

- **floating** (Boolean)
- **id** (String)
- **name** (String)
