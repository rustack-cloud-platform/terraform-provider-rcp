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

data "rustack_hypervisors" "all_hypervisors" {
    project_id = data.rustack_project.single_project.id
}
