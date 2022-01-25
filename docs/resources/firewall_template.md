---
page_title: "rustack_firewall_template Resource - terraform-provider-rustack"
---
# rustack_firewall_template (Resource)

Firewall allow you to manage your network traffic.

## Example Usage

```hcl

resource "rustack_firewall_template" "single_template" {
  vdc_id = data.rustack_vdc.single_vdc.id
  name   = "New custom template"

  ingress_rule {
    protocol       = "tcp"
    port_range     = "80"
    destination_ip = "2.0.0.0"
  }

  ingress_rule {
    protocol       = "icmp"
    destination_ip = "1.0.0.0/24"
  }

  egress_rule {
    protocol       = "udp"
    port_range     = "53"
    destination_ip = "5.0.0.0/24"
  }
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
