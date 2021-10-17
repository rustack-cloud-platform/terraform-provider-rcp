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

data "rustack_disk" "single_disk" {
    vdc_id = data.rustack_vdc.single_vdc.id
    name = "Disk 2"
}
