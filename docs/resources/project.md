---
page_title: "rustack_project Resource - terraform-provider-rustack"
---
# rustack_project (Resource)

Projects allow you to organize your resources into groups that fit the way you work.

The Vdcs can be associated with a project:

## Example Usage

```hcl
resource "rustack_project" "demo_project" {
    name = "Terraform Project"
}
```

## Schema

### Required

- **name** (String) name of the Project

### Optional

- **id** (String) The ID of this resource.
