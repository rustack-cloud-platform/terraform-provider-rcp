---
page_title: "rustack_dnss Data Source - terraform-provider-rustack"
---
# rustack_dnss (Data Source)

Get information about Dnss in the Project for use in other resources.

Note: You can use the [`rustack_dns`](Dns) data source to obtain metadata
about a single Dns if you already know the `name` and unique `project_id` to retrieve.

## Example Usage

```hcl

data "rustack_project" "single_project" {
    name = "Terraform Project"
}

data "rustack_dnss" "dns" {
    project_id = data.rustack_project.single_project.id
}


```

## Schema

### Required

- **project_id** (String) id of the Project

### Read-Only

- **dnss** (List of Object) (see [below for nested schema](#nestedatt--dnss))

<a id="nestedatt--dnss"></a>
### Nested Schema for `dnss`

Read-Only:

- **id** (String)
- **name** (String)
