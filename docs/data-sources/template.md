---
page_title: "rustack_template Data Source - terraform-provider-rustack"
---
# rustack_template (Data Source)

Get information on a Template for use in other resources. 

## Example Usage

```hcl

data "rustack_project" "single_project" {
    name = "Terraform Project"
}

data "rustack_vdc" "single_vdc" {
    project_id = "${data.rustack_project.single_project.id}"
    name = "Terraform VDC"
}

data "rustack_template" "single_template" {
    vdc_id = data.rustack_vdc.single_vdc.id
    name = "Debian 10"
}

```

## Schema

### Required

- **name** (String) name of the Template
- **vdc_id** (String) id of the VDC

### Read-Only

- **id** (String) id of the Template
- **min_cpu** (Number) minimum cpu required by the Template
- **min_disk** (Number) minimum disk size in GB required by the Template
- **min_ram** (Number) minimum ram in GB required by the Template
