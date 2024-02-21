---
page_title: "rustack_paas_template Data Source - terraform-provider-rustack"
---
# rustack_paas_template (Data Source)

Get information about a PaaS Service Template for use in other resources. 

## Example Usage

```hcl

data "rustack_project" "single_project" {
    name = "Terraform Project"
}

data "rustack_paas_template" "paas_template" {
  id = 1
  project_id = data.rustack_project.single_project.id
}
```
## Schema

### Required

- **project_id** (String) id of Project
- **id** (String) id of PaaS Service Template

### Read-Only

- **name** (String) name of PaaS Service Template
