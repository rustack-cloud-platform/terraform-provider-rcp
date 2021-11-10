---
page_title: "rustack_account Data Source - terraform-provider-rustack"
---
# rustack_account (Data Source)

Get information on a Acconut for use in other resources. 

## Example Usage

```hcl

data "rustack_account" "account" { }

```
## Schema

### Read-Only

- **email** (String) The email address of current user
- **id** (String) The identifier for the current user
- **username** (String) The username of current user
