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

resource "rustack_vm" "vm1" {
    vdc_id = data.rustack_vdc.single_vdc.id

    name = "Server 1"
    cpu = 2
    ram = 4

    template_id = data.rustack_template.debian10.id

    user_data = file("user_data.yaml")

    system_disk = "1-ssd" 
    
    disks = [
        data.rustack_disk.new_disk1,
        data.rustack_disk.new_disk2,
    ]

    port {
        network_id = data.rustack_network.service_network.id
        firewall_templates = [data.rustack_firewall_template.allow_default.id,
            data.rustack_firewall_template.allow_web.id,
            data.rustack_firewall_template.allow_ssh.id
        ]
    }

    floating = false
}
```

## Schema

### Required

- **cpu** (Number) the number of virtual cpus
- **system_disk** (String) System disk. Format `1-ssd` where 1 is size in Gb and `ssd` is storage profile.
- **name** (String) name of the Vm
- **port** (Block List, Min: 1, Max: 10) list of Ports attached to the Vm (see [below for nested schema](#nestedblock--port))
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
- **storage_profile_id** (String) the id of the StorageProfile

Read-Only:

- **id** (String) id of the Disk


<a id="nestedblock--port"></a>
### Nested Schema for `port`

Required:

- **firewall_templates** (List of String) list of firewall rule ids of the Port
- **network_id** (String) id of the Network

Read-Only:

- **id** (String) id of the Port
- **ip_address** (String) ip_address of the Port

Optional:

- **create** (String)
- **delete** (String)
