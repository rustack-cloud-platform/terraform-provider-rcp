---
page_title: "rustack_ports Data Source - terraform-provider-rustack"
---
# rustack_ports (Data Source)

Get information about list of Ports in the Vdc for use in other resources.

## Example Usage

```hcl

data "rustack_project" "single_project" {
    name = "Terraform Project"
}

data "rustack_vdc" "single_vdc" {
    project_id = data.rustack_project.single_project.id
    name = "Terraform VDC"
}

data "rustack_port" "all_port" {
    vdc_id = data.rustack_vdc.single_vdc.id
}

```

## Schema

### Required

- **vdc_id** (String) id of the VDC

### Read-Only

- **ports** (List of Object) (see [below for nested schema](#nestedatt--ports))

<a id="nestedatt--ports"></a>
### Nested Schema for `ports`

Read-Only:

- **id** (String)
- **network_id** (String)
- **ip_address** (String)
- **firewall_templates** (List of String)
