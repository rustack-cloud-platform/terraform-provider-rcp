---
page_title: "rustack_vm Resource - terraform-provider-rustack"
---
# rustack_vm (Resource)

This data source provides creating and deleting vms. You should have a vdc to create a vm.

## Example Usage

```hcl 
data "rustack_project" "single_project" {
    name = "Terraform Project"
}

data "rustack_vdc" "single_vdc" {
    project_id = data.rustack_project.single_project.id
    name = "Terraform VDC"
}

data "rustack_network" "service_network" {
    vdc_id = data.rustack_vdc.single_vdc.id
    name = "Сеть"
}

data "rustack_storage_profile" "ssd" {
    vdc_id = data.rustack_vdc.single_vdc.id
    name = "ssd"
}

data "rustack_storage_profile" "sas" {
    vdc_id = data.rustack_vdc.single_vdc.id
    name = "sas"
}

data "rustack_template" "debian10" {
    vdc_id = data.rustack_vdc.single_vdc.id
    name = "Debian 10"
}

data "rustack_firewall_template" "allow_default" {
    vdc_id = data.rustack_vdc.single_vdc.id
    name = "По-умолчанию"
}

data "rustack_firewall_template" "allow_web" {
    vdc_id = data.rustack_vdc.single_vdc.id
    name = "Разрешить WEB"
}

data "rustack_firewall_template" "allow_ssh" {
    vdc_id = data.rustack_vdc.single_vdc.id
    name = "Разрешить SSH"
}

data "rustack_port" "vm_port" {
    vdc_id = resource.rustack_vdc.single_vdc.id

    network_id = resource.rustack_network.network.id
    firewall_templates = [data.rustack_firewall_template.allow_default.id]
}

resource "rustack_vm" "vm1" {
    vdc_id = data.rustack_vdc.single_vdc.id

    name = "Server 1"
    cpu = 2
    ram = 4

    template_id = data.rustack_template.debian10.id

    user_data = file("user_data.yaml")

    system_disk {
        size = 10
        storage_profile_id = data.rustack_storage_profile.ssd.id
    }
    
    disks = [
        data.rustack_disk.new_disk1,
        data.rustack_disk.new_disk2,
    ]

    ports {
        data.rustack_port.vm_port
    }

    floating = false
}
```

## Schema

### Required

- **cpu** (Number) the number of virtual cpus
- **system_disk** System disk (Min: 1, Max: 1). 
- **name** (String) name of the Vm
- **ports** (List of String) list of Ports id attached to the Vm. 
- **ram** (Number) memory of the Vm in gigabytes
- **template_id** (String) id of the Template
- **user_data** (String) script for cloud-init
- **vdc_id** (String) id of the VDC

### Optional

- **floating** (Boolean) enable floating ip for the Vm
- **id** (String) The ID of this resource.
- **timeouts** (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- **disks** (Toset, String) list of Disks id attached to the Vm.


### Read-Only

- **floating_ip** (String) floating ip for the Vm. May be omitted

<a id="nestedblock--disk"></a>
### Nested Schema for `disk`

Required:

- **name** (String) name of the Disk
- **size** (Number) the size of the Disk in gigabytes
- **storage_profile_id** (String) Id of the storage profile
- **vdc_id** (String) id of the VDC

Optional:

- **timeouts** (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

Read-Only:

- **id** (String) id of the Disk
