---
page_title: "rustack_paas_service Resource - terraform-provider-rustack"
---
# rustack_paas_service (Resource)

Provides a Rustack PaaS Service resource.

## Example Usage

```hcl
data "rustack_project" "single_project" {
    name = "Terraform Project"
}

resource "rustack_paas_service" "service" {
  name = "terraform paas service"
  project_id = data.rustack_project.single_project.id
  paas_service_id = 1
  paas_service_inputs = jsonencode({
    "change_password": false,
    "enable_ssh_password":true,
    "enable_sudo":true,
    "passwordless_sudo":false,
    "cpu_num":1,
    "ram_size":1,
    "volume_size":10,
    "network_id":"c7ec518c-a62c-42cc-adee-9cf731d066e4",
    "vdcs_id":"e9da805f-04c3-4a93-9542-9724486bebe0",
    "vm_name":"vm name",
    "template_id":"4c8b06e6-7909-4fb9-9097-6ff900a848d0",
    "storage_profile":"0e1aead5-ef4b-46e2-b988-64dea6d146f8",
    "firewall_profiles":["85929e4e-d12d-411a-9763-b3f8c3d279a0","dc4203e4-d7fe-45f0-91b8-a33128e1a089"],
    "user_name":"ubuntu",
    "user_password":"ubuntu"
  })
}

```

## Schema

### Required

- **project_id** (String) id of Project
- **name** (String) name of PaaS Service
- **paas_service_id** (String) id of PaaS Service Template
- **paas_service_id** (String) id of PaaS Service Template


### Read-Only

- **id** (Boolean) id of PaaS Service
