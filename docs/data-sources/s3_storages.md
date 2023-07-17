---
page_title: "rustack_s3_storages Data Source - terraform-provider-rustack"
---
# rustack_s3_storages (Data Source)

Returns a list of Rustack S3Storage.

Get information about S3Storage in the Vdc for use in other resources.

Note: You can use the [`rustack_s3_storage`](S3Storage) data source to obtain metadata
about a single s3 storage if you already know the `name` and `project_id` to retrieve.

## Example Usage

```hcl

data "rustack_project" "single_project" {
    name = "Terraform Project"
}


data "rustack_s3_storages" "s3_storages" {
    project_id = resource.rustack_project.single_project.id
}

```

## Schema

### Required

- **project_id** (String) id of the VDC

### Read-Only

- **s3_storages** (List of Object) (see [below for nested schema](#nestedatt--s3_storage))

<a id="nestedatt--s3_storage"></a>
### Nested Schema for `s3_storage`

Read-Only:

- **id** (String)
- **client_endpoint** (String)
- **access_key** (String)
- **secret_key** (String)
- **name** (String)
- **backend** (String)
