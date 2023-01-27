---
page_title: "rustack_s3_storage Data Source - terraform-provider-rustack"
---
# rustack_s3_storage (Data Source)

Get information about a S3Storage for use in other resources. 

## Example Usage

```hcl

data "rustack_project" "single_project" {
    name = "Terraform Project"
}

data "rustack_s3_storage" "s3_storage" {
    project_id = resource.rustack_project.single_project.id
    name = "s3_storage"
}

```

## Schema

### Required

- **project_id** (String) id of the project
- **name** (String) name of the S3Storage

### Read-Only

- **id** (String) id of the S3Storage
- **client_endpoint** (String) url for connecting to s3"
- **access_key** (String) access_key for access to s3
- **secret_key** (String) secret_key for access to s3