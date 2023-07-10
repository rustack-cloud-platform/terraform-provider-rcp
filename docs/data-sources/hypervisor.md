---
page_title: "rustack_hypervisor Data Source - terraform-provider-rustack"
---
# rustack_hypervisor (Data Source)

Get information about a Hypervisor for use in other resources. 

## Example Usage

```hcl

data "rustack_project" "single_project" {
    name = "Terraform Project"
}

data "rustack_hypervisor" "single_hypervisor" {
    project_id = data.rustack_project.single_project.id
    
    name = "VMWARE"
    # or
    id ="id"
}

```

## Schema

### Required

- **name** (String) name of the Hypervisor `or` **id** (String) id of the Hypervisor
- **project_id** (String) id of the Project

### Read-Only

- **type** (String) type of the Hypervisor
