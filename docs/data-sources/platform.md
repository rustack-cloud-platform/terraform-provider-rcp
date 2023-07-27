---
page_title: "rustack_platform Data Source - terraform-provider-rustack"
---
# rustack_platform (Data Source)
### `Only for Vmware Hypervisor`
Get information about a Platform for use in other resources. 

## Example Usage

```hcl
data "rustack_project" "single_project" {
    name = "Terraform Project"
}
data "rustack_vdc" "single_vdc" {
    project_id = data.rustack_project.single_project.id
    name = "Terraform VDC"
}
data "rustack_platform" "platform"{
    vdc_id = resource.rustack_vdc.single_vdc.id
    name = "Intel Cascade Lake"
    # or
    id = ""
}
```
## Schema

### Required

- **name** (String) name of the Platform `or` **id** (String) id of the Platform
