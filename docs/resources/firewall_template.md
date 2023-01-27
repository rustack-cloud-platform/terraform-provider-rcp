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

data "rustack_firewall_template" "single_template" {
  vdc_id = data.rustack_vdc.single_vdc.id
  name   = "New custom template"
}

```

## Schema

### Required

- **name** (String) name of the Firewall
- **vdc_id** (String) id of the VDC

### Optional

- **ingress_rule** (Schema) Schema for ingress template rule.

    **protocol** (String) udp/tcp/icmp protocols
    **port_range** (String) You can set only one number or range like `80:90`
    **destination_ip** (String) Destination Ip address 

- **egress_rule** (Schema) Schema for egress template rule.

    **protocol** (String) udp/tcp/icmp protocols
    **port_range** (String) You can set only one number or range like `80:90`
    **destination_ip** (String) Destination Ip address 
