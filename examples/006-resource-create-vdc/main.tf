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

data "rustack_hypervisor" "single_hypervisor" {
    project_id = data.rustack_project.single_project.id
    name = "VMWARE"
}

resource "rustack_vdc" "vdc1" {
    name = "Terraform VDC"
    project_id = "${data.rustack_project.single_project.id}"
    hypervisor_id = data.rustack_hypervisor.single_hypervisor.id
}
