terraform {
  required_version = ">= 1.0.0"

  required_providers {
    rustack = {
      source  = "rustack-cloud-platform/rcp"
    }
  }
}

provider "rustack" {
  token = "[PLACE_YOUR_TOKEN_HERE]"
}
data "rustack_project" "single_project" {
  name = "Terraform Project"
}

data "rustack_vdc" "single_vdc" {
  project_id = data.rustack_project.single_project.id
  name       = "Terraform VDC"
}

data "rustack_storage_profile" "single_storage_profile" {
  vdc_id = data.rustack_vdc.single_vdc.id
  name   = "sas"
}

resource "rustack_disk" "disk2" {
  vdc_id = data.rustack_vdc.single_vdc.id

  name               = "Disk 1"
  storage_profile_id = data.rustack_storage_profile.single_storage_profile.id
  size               = 1
}
