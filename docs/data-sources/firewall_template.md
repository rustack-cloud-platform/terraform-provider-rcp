---
page_title: "rustack_firewall_template Data Source - terraform-provider-rustack"
---
# rustack_firewall_template (Data Source)

Get information about a Firewall Template for use in other resources. 

## Example Usage

```hcl

data "rustack_project" "single_project" {
    name = "Terraform Project"
}

data "rustack_vdc" "single_vdc" {
    project_id = data.rustack_project.single_project.id
    name = "Terraform VDC"
}

data "rustack_firewall_template" "single_template" {
    vdc_id = data.rustack_vdc.single_vdc.id
    name = "Разрешить Web"
}

```
## Schema

### Required

- **name** (String) name of the Template
- **vdc_id** (String) id of the VDC

### Read-Only

- **id** (String) id of the Template
