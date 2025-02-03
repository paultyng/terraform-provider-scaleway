---
subcategory: "Edge Services"
page_title: "Scaleway: scaleway_edge_services_backend_stage"
---

# Resource: scaleway_edge_services_backend_stage

Creates and manages Scaleway Edge Services Backend Stages.

## Example Usage

### Basic

```terraform
resource "scaleway_object_bucket" "main" {
    name = "my-bucket-name"
    tags = {
        foo = "bar"
    }
}

resource "scaleway_edge_services_backend_stage" "main" {
  s3_backend_config {
    bucket_name   = scaleway_object_bucket.main.name
    bucket_region = "fr-par"
  }
}
```

### Custom Certificate

```terraform
```

## Argument Reference

- `s3_backend_config` - (Required) The Scaleway Object Storage origin bucket (S3) linked to the backend stage.
    - `bucket_name` - The name of the Bucket.
    - `bucket_region` - The region of the Bucket.
    - `is_website` - Defines whether the bucket website feature is enabled.
- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the backend stage is associated with.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the backend stage (UUID format).
- `created_at` - The date and time of the creation of the backend stage.
- `updated_at` - The date and time of the last update of the backend stage.
- `pipeline_id` - The pipeline ID the backend stage belongs to.

## Import

Backend stages can be imported using the `{id}`, e.g.

```bash
$ terraform import scaleway_edge_services_backend_stage.basic 11111111-1111-1111-1111-111111111111
```
