---
page_title: "rustack_dns_record Resource - terraform-provider-rustack"
---
# rustack_dns_record (Resource)

Provides a Rustack DNS record resource.

## Example Usage

```hcl
data "rustack_project" "single_project" {
    name = "Terraform Project"
}

data "rustack_dns" "dns" {
    name="dns.teraform."
    project_id = data.rustack_project.single_project.id
}

resource "rustack_dns_record" "dns_record1" {
    dns_id = data.rustack_dns.dns.id
    type = "A"
    host = "test2.dns.teraform."
    data = "8.8.8.8"
}

```

## Schema

### Required

> required for all types

- **dns_id** (String) name of the Dns
- **type** (String) type of Dns record
- **host** (String) host of Dns record
- **data** (String) data of Dns record

> for type CAA parameters are required to

- **tag** (String) tag of Dns record
- **flag** (String) flag of Dns record. Can be chosen
    **0 (not critical)**, **128 (critical)**

> for type MX parameters are required to

- **Priority** (String) Priority of Dns record

> for type SRV parameters are required to

- **Priority** (String) Priority of Dns record
- **Weight** (String) Weight of Dns record
- **Port** (String) Port of Dns record

### Optional

- **id** (String) The ID of this resource.
- **timeouts** (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- **create** (String)
- **delete** (String)
