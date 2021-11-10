---
page_title: "rustack_vdc Data Source - terraform-provider-rustack"
---
# rustack_vdc (Data Source)

Get information about a Vdc for use in other resources. 
This is useful if you need to utilize any of the Vdc's data and Vdc not managed by Terraform.

## Example Usage

```hcl

data "rustack_project" "single_project" {
    name = "Terraform Project"
}

data "rustack_vdc" "single_vdc" {
    project_id = "${data.rustack_project.single_project.id}"
    name = "Terraform VDC"
}

data "rustack_vdc" "single_vdc2" {
    name = "Terraform VDC"
}

```

## Schema

### Required

- **name** (String) name of the vdc

### Optional

- **project_id** (String) id of the Project

### Read-Only

- **hypervisor** (String) name of the Hypervisor
- **hypervisor_type** (String) type of the Hypervisor
- **id** (String) id of the VDC
