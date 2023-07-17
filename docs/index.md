---
page_title: "Rustack Provider"
---
# rustack Provider

The Rustack provider is used to interact with the Rustack cloud. 
The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
terraform {
  required_providers {
    rustack = {
      source  = "pilat/rustack"
    }
  }
}

# Set the variable value in *.tfvars file
# or using -var="rustack_token=..." CLI option
variable "rustack_token" {}

# Configure the Rustack Provider
provider "rustack" {
    api_endpoint = "https://cp.iteco.cloud"
    token = var.rustack_token
}

```

-> **Note for Module Developers** Although provider configurations are shared between modules, each module must
declare its own [provider requirements](https://www.terraform.io/docs/language/providers/requirements.html). See the [module development documentation](https://www.terraform.io/docs/language/modules/develop/providers.html) for additional information.

## Schema

### Optional

- **api_endpoint** (String) The URL to use for the Rustack API.
- **token** (String) The token key for API operations.
