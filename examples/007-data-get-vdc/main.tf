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
    name = "Terraform VDC"
    # or
    id = "id"
}

data "rustack_vdc" "single_vdc2" {
    name = "Terraform VDC"
    # or
    id = "id"
}
