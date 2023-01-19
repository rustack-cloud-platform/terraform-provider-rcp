terraform {
  required_version = ">= 1.0.0"

  required_providers {
    rustack = {
      source  = "pilat/rustack"
    }
  }
}

provider "rustack" {
  token = "[PLACE_YOUR_TOKEN_HERE]"
}

data "rustack_project" "single_project" {
  name = "Terraform Project"
}

resource "rustack_s3_storage" "s3_storage" {
    project_id = resource.rustack_project.single_project.id
    name = "s3_storage"
}