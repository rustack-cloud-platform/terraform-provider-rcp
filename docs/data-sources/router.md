---
page_title: "rustack_router Data Source - terraform-provider-rustack"
---
# rustack_router (Data Source)

Get information about a Router for use in other resources. 

## Example Usage

```hcl

data "rustack_project" "single_project" {
    name = "Terraform Project"
}

data "rustack_vdc" "single_vdc" {
    project_id = data.rustack_project.single_project.id"
    name = "Terraform VDC"
}

data "rustack_router" "single_Router" {
    vdc_id = data.rustack_vdc.single_vdc.id
    
    name = "Terraform Router"
    # or
    id = "id"
}

```
## Schema

### Required

- **vdc_id** (String) id of the VDC
- **name** (String) name of the Router `or` **id** (String) id of the Router

