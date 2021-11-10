---
page_title: "rustack_project Data Source - terraform-provider-rustack"
---
# rustack_project (Data Source)

Get information about a Project for use in other resources. 

## Example Usage

```hcl

data "rustack_project" "single_project" {
    name = "Terraform Project"
}

```
## Schema

### Required

- **name** (String) name of the Project

### Read-Only

- **id** (String) id of the Project
