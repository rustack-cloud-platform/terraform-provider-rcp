---
page_title: "rustack_projects Data Source - terraform-provider-rustack"
---
# rustack_projects (Data Source)

Get information about Projects in the Vdc for use in other resources.

Note: You can use the [`rustack_project`](Project) data source to obtain metadata
about a single Project if you already know the `name` to retrieve.

## Example Usage

```hcl

data "rustack_projects" "all_projects" { }

```

## Schema

### Read-Only

- **projects** (List of Object) (see [below for nested schema](#nestedatt--projects))

<a id="nestedatt--projects"></a>
### Nested Schema for `projects`

Read-Only:

- **id** (String)
- **name** (String)


