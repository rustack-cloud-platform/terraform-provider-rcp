---
page_title: "rustack_kubernetes Data Source - terraform-provider-rustack"
---
# rustack_kubernetes (Data Source)

Get information about a Kubernetes for use in other resources. 
This is useful if you need to utilize any of the Kubernetes's data and Kubernetes not managed by Terraform.

**Note:** This data source returns a single Kubernetes. When specifying a `name`, an
error is triggered if more than one Kubernetes is found.

## Example Usage

```hcl

data "rustack_project" "single_project" {
    name = "Terraform Project"
}

data "rustack_vdc" "single_vdc" {
    project_id = data.rustack_project.single_project.id
    name = "Terraform VDC"
}

data "rustack_kubernetes" "single_k8s" {
    vdc_id = data.rustack_vdc.single_vdc.id
    
    name = "Server 1"
    # or
    id = "id"
}

```

## Schema

### Required

- **name** (String) name of the Kubernetes `or` **id** (String) id of the Kubernetes
- **vdc_id** (String) id of the VDC

### Read-Only

- **node_cpu** (Integer) the number of virtual cpus
- **floating** (Boolean) enable floating ip for the Kubernetes
- **floating_ip** (String) floating_ip of the Kubernetes. May be omitted
- **nodes_count** (String) vms count of the Kubernetes
- **node_disk_size** (String) disk size of the attached vms in Kubernetes
- **user_public_key_id** (String) public key of the Kubernetes
- **node_storage_profile_id** (String) storage profile id of the attached vms in Kubernetes
- **dashboard_url** (String) dashboard url of the Kubernetes
- **node_ram** (Integer) the number of ram of the attached vms in Kubernetes
- **template_id** (String) id of the Template
- **vms** (List) List of vms attached to Kubernetes

## Getting information about kubernetes

### Get dashboard url
- **This block will print dashboard_url in console**
```
    output "dashboard_url" {
        value = data.rustack_kubernetes.single_k8s.dashboard_url
    }
```
### Get kubectl config
- **When kubernetes is received, the kubectl configuration will appear in the workdir folder.**
