---
page_title: "rustack_vm Data Source - terraform-provider-rustack"
---
# rustack_vm (Data Source)

Get information about a Vm for use in other resources. 
This is useful if you need to utilize any of the Vm's data and Vm not managed by Terraform.

**Note:** This data source returns a single Vm. When specifying a `name`, an
error is triggered if more than one Vm is found.

## Example Usage

```hcl

data "rustack_project" "single_project" {
    name = "Terraform Project"
}

data "rustack_vdc" "single_vdc" {
    project_id = data.rustack_project.single_project.id
    name = "Terraform VDC"
}

data "rustack_vm" "single_vm" {
    vdc_id = data.rustack_vdc.single_vdc.id
    
    name = "Server 1"
    # or
    id = "id"
}

```

## Schema

### Required

- **name** (String) name of the Vm `or` **id** (String) id of the Vm
- **vdc_id** (String) id of the VDC

### Read-Only

- **cpu** (Integer) the number of virtual cpus
- **floating** (Boolean) enable floating ip for the Vm
- **floating_ip** (String) floating_ip of the Vm. May be omitted
- **ram** (Float) memory of the Vm in gigabytes
- **template_id** (String) id of the Template
- **template_name** (String) name of the Template
- **power** (Boolean) the vm state
- **ports** (Block List)    (see [below for nested schema](#nestedblock--port))

<a id="nestedblock--port"></a>
### Nested Schema for `port`

Required:

- **id** (String) id of the Port

Read-Only:

- **ip_address** (String) IP of the Port
