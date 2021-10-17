terraform {
  required_version = ">= 1.0.0"

  required_providers {
    rustack = {
      source  = "rustack/rustack"
      version = "~> 0.1"
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
    project_id = "${data.rustack_project.single_project.id}"
    name = "Terraform VDC"
}

data "rustack_storage_profiles" "all_storage_profiles" {
    vdc_id = data.rustack_vdc.single_vdc.id
    # vdc_name = "Terraform VDC"
    # vdc_id = "e76abe25-2e02-4652-b8b8-39531be04c63"
}

