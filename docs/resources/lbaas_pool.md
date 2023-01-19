---
page_title: "rustack_dns_record Resource - terraform-provider-rustack"
---
# rustack_dns_record (Resource)

Provides a Rustack DNS record resource.

## Example Usage

```hcl
data "rustack_project" "single_project" {
    name = "Terraform Project"
}

data "rustack_vdc" "single_vdc" {
    project_id = data.rustack_project.single_project.id
    name = "Terraform VDC"
}

data "rustack_network" "new_network" {
    vdc_id =  data.rustack_vdc.single_vdc.id
    name = "New network"
}

data "rustack_template" "debian10" {
    vdc_id = data.rustack_vdc.single_vdc.id
    name = "Debian 10"
}

data "rustack_firewall_template" "allow_default" {
    vdc_id = data.rustack_vdc.single_vdc.id
    name = "По-умолчанию"
}

data "rustack_lbaas" "lbaas" {
    vdc_id = data.rustack_project.single_vdc.id
    name = "lbaas"
}

data "rustack_port" "vm_port" {
    vdc_id = resource.rustack_vdc.single_vdc.id

    network_id = resource.rustack_network.new_network.id
    firewall_templates = [data.rustack_firewall_template.allow_default.id]
}

data "rustack_vm" "vm" {
    vdc_id = data.rustack_vdc.single_vdc.id
    name = "Server 1"
}

resource "rustack_lbaas_pool" "pool" {
    lbaas_id = data.rustack_lbaas.lbaas.id
    connlimit = 65536
    method = "ROUND_ROBIN"
    port = 2
    protocol = "TCP"
    member {
        port = 2
        weight = 1
        vm_id = data.rustack_vm.vm.id
    }
}

```

## Schema

### Required

- **lbaas_id** (String) id of LoadBalancer
- **port** (Integer) port of LoadBalancerPool
- **member** (String) parameter that specifies which network will be connected to LoadBalancer  (see [below for nested schema](#nestedblock--member))


### Optional

- **method** (String) method of LoadBalancerPool 
> Can be chosen ROUND_ROBIN, LEAST_CONNECTIONS, SOURCE_IP
- **protocol** (String) method of LoadBalancerPool
> Can be chosen TCP, HTTP, HTTPS
- **connlimit** (Integer) connlimit of LoadBalancerPool
- **timeouts** (Block, Optional)

<a id="nestedblock--member"></a>
### Nested Schema for `member`

Required:

- **port** (Integer) id of the Network
- **vm_id** (String) id of the Network

Optional:

- **weight** (Integer) id of the Network
