---
page_title: "rustack_pub_key Data Source - terraform-provider-rustack"
---
# rustack_pub_key (Data Source)

Get information about a public key for use in other resources. 

## Example Usage

```hcl

data "rustack_account" "me"{}

data "rustack_pub_key" "key" {
    account_id = data.rustack_account.me.id
    
    name = "Debian 10"
    # or
    id = "id"
}

```

## Schema

### Required

- **name** (String) name of the public key `or` **id** (String) id of the public key
- **account_id** (String) id of the account

### Read-Only

- **fingerprint** (Number) fingerprint of public key
- **public_key** (Number) public_key value of public key data source
