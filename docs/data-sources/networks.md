---
page_title: "rustack_networks Data Source - terraform-provider-rustack"
---
# rustack_networks (Data Source)

Get information on Networks in the Vdc for use in other resources.

Note: You can use the [`rustack_network`](Network) data source to obtain metadata
about a single Network if you already know the `name` and `vdc_id` to retrieve.

## Example Usage

```hcl

data "rustack_project" "single_project" {
    name = "Terraform Project"
}

data "rustack_vdc" "single_vdc" {
    project_id = "${data.rustack_project.single_project.id}"
    name = "Terraform VDC"
}

data "rustack_networks" "all_networks" {
    vdc_id = data.rustack_vdc.single_vdc.id
}

```

## Schema

### Required

- **vdc_id** (String) id of the VDC

### Read-Only

- **networks** (List of Object) (see [below for nested schema](#nestedatt--networks))

<a id="nestedatt--networks"></a>
### Nested Schema for `networks`

Read-Only:

- **id** (String)
- **name** (String)
- **subnets** (List of Object) (see [below for nested schema](#nestedobjatt--networks--subnets))

<a id="nestedobjatt--networks--subnets"></a>
### Nested Schema for `networks.subnets`

Read-Only:

- **cidr** (String)
- **dhcp** (Boolean)
- **dns** (List of String)
- **end_ip** (String)
- **gateway** (String)
- **id** (String)
- **start_ip** (String)
