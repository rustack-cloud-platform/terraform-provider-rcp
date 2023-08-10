---
page_title: "rustack_kubernetes_template Data Source - terraform-provider-rustack"
---
# rustack_kubernetes_template (Data Source)

Get information about a kubernetes template for use in other resources. 

## Example Usage

```hcl

data "rustack_project" "single_project" {
    name = "Terraform Project"
}

data "rustack_vdc" "single_vdc" {
    project_id = data.rustack_project.single_project.id
    name = "Terraform VDC"
}

data "rustack_kubernetes_template" "k8s_template" {
    vdc_id = data.rustack_vdc.single_vdc.id
    
    name = "Kubernetes 1.22.1"
    # or
    id = "id"
}

```

## Schema

### Required

- **name** (String) name of the kubernetes template `or` **id** (String) id of the kubernetes template
- **vdc_id** (String) id of the VDC

### Read-Only

- **min_node_cpu** (Integer) minimum cpu required by the kubernetes template
- **min_node_hdd** (Integer) minimum disk size in GB required by the kubernetes template
- **min_node_ram** (Integer) minimum ram in GB required by the kubernetes template

