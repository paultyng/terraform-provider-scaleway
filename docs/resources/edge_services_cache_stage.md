---
subcategory: "Edge Services"
page_title: "Scaleway: scaleway_edge_services_cache_stage"
---

# Resource: scaleway_edge_services_cache_stage

Creates and manages Scaleway Edge Services Cache Stages.

## Example Usage

### Basic

```terraform
resource "scaleway_edge_services_cache_stage" "main" {
  backend_stage_id = scaleway_edge_services_backend_stage.main.id
}
```

### Purge request

```terraform
resource "scaleway_edge_services_cache_stage" "main" {
  backend_stage_id = scaleway_edge_services_backend_stage.main.id

  purge {
    pipeline_id = scaleway_edge_services_pipeline.main.id
    all         = true
  } 
}
```

## Argument Reference

- `backend_stage_id` - (Optional) The backend stage ID the cache stage will be linked to.
- `fallback_ttl` - (Optional) The Time To Live (TTL) in seconds. Defines how long content is cached.
- `refresh_cache` - (Optional) Trigger a refresh of the cache by changing this field's value.
- `purge_requests` - (Optional) The Scaleway Object Storage origin bucket (S3) linked to the backend stage.
    - `pipeline_id` - The pipeline ID in which the purge request will be created.
    - `assets` - The list of asserts to purge.
    - `all` - Defines whether to purge all content.
- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the cache stage is associated with.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the cache stage (UUID format).
- `created_at` - The date and time of the creation of the cache stage.
- `updated_at` - The date and time of the last update of the cache stage.
- `pipeline_id` - The pipeline ID the cache stage belongs to.

## Import

Cache stages can be imported using the `{id}`, e.g.

```bash
$ terraform import scaleway_edge_services_cache_stage.basic 11111111-1111-1111-1111-111111111111
```
