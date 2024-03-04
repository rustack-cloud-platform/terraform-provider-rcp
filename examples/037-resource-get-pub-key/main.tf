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

data "rustack_account" "me"{}

data "rustack_pub_key" "key"{
    name = "name"
    # or
    or = "id"
    account_id = data.rustack_account.me.id
}