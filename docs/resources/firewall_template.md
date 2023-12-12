---
page_title: "rustack_firewall_template Resource - terraform-provider-rustack"
---
# rustack_firewall_template (Resource)

Firewall allow you to manage your network traffic.

## Example Usage

```hcl

data "rustack_project" "single_project" {
    name = "Terraform Project"
}

data "rustack_vdc" "single_vdc" {
    project_id = data.rustack_project.single_project.id
    name = "Terraform VDC"
}

resource "rustack_firewall_template" "single_template" {
  vdc_id = data.rustack_vdc.single_vdc.id
  name   = "New custom template"
  tags = ["created_by:terraform"]
}

```

## Schema

### Required

- **name** (String) name of the FirewallTemplate
- **vdc_id** (String) id of the VDC

### Optional

- **tags** (Toset, String) list of Tags added to the FirewallTemplate
