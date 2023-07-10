---
page_title: "rustack_s3_storage Resource - terraform-provider-rustack"
---
# rustack_s3_storage (Resource)

This data source provides creating and deleting s3 storage. You should have a project to create a s3_storage.

## Example Usage

```hcl 
data "rustack_project" "single_project" {
    name = "Terraform Project"
}
resource "rustack_s3_storage" "s3_storage" {
    project_id = resource.rustack_project.single_project.id
    name = "s3_storage"
    backend = "minio" # or "netapp"
}
```

## Schema

### Required

- **name** (String) name of the s3_storage
- **project_id** (String) id of the project
- **backend** (String) backend of the s3_storage (`minio` or `netapp`)

### Optional

- **id** (String) The ID of this resource.
- **client_endpoint** (Boolean) url for connecting to s3
- **access_key** (String) access_key for connecting to s3
- **secret_key** (String) secret_key for connecting to s3
