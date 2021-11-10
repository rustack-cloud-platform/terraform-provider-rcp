---
page_title: "rustack_storage_profile Data Source - terraform-provider-rustack"
---
# rustack_storage_profile (Data Source)

Get information about a Storage Profile for use in other resources. 

## Example Usage

```hcl

data "rustack_project" "single_project" {
    name = "Terraform Project"
}

data "rustack_vdc" "single_vdc" {
    project_id = "${data.rustack_project.single_project.id}"
    name = "Terraform VDC"
}

data "rustack_storage_profile" "single_storage_profile" {
    vdc_id = data.rustack_vdc.single_vdc.id
    name = "ssd"
}

```
## Schema

### Required

- **name** (String) Name of the storage profile
- **vdc_id** (String) id of the VDC

### Read-Only

- **id** (String) The identifier for the storage profile
