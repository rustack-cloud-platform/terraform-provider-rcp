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

resource "rustack_project" "demo_project" {
    name = "Terraform Project"
}
