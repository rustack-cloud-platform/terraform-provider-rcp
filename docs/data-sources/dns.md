---
page_title: "rustack_dns Data Source - terraform-provider-rustack"
---
# rustack_dns (Data Source)

Get information about a Dns for use in other resources. 
This is useful if you need to utilize any of the Dns's data and dns not managed by Terraform.

## Example Usage

```hcl

data "rustack_project" "single_project" {
    name = "Terraform Project"
}

data "rustack_dns" "dns" {
    name="dns.teraform."
    project_id = data.rustack_project.single_project.id
}

```

## Schema

### Required

- **project_id** (String) id of the Project
- **name** (String) name of the dns zone

### Read-Only

- **id** (String) id of the VDC
