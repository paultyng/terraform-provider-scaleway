---
subcategory: "Edge Services"
page_title: "Scaleway: scaleway_edge_services_dns_stage"
---

# Resource: scaleway_edge_services_dns_stage

Creates and manages Scaleway Edge Services DNS Stages.

## Example Usage

### Basic

```terraform
resource "scaleway_edge_services_dns_stage" "main" {
  fqdns = ["subdomain.example.com"]
}
```

## Argument Reference

- `backend_stage_id` - (Optional) The backend stage ID the DNS stage will be linked to.
- `tls_stage_id` - (Optional) The TLS stage ID the DNS stage will be linked to.
- `cache_stage_id` - (Optional) The cache stage ID the DNS stage will be linked to.
- `fqdns` - (Optional) Fully Qualified Domain Name (in the format subdomain.example.com) to attach to the stage.
- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the DNS stage is associated with.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the DNS stage (UUID format).
- `type` - The type of the stage.
- `created_at` - The date and time of the creation of the DNS stage.
- `updated_at` - The date and time of the last update of the DNS stage.
- `pipeline_id` - The pipeline ID the DNS stage belongs to.

## Import

DNS stages can be imported using the `{id}`, e.g.

```bash
$ terraform import scaleway_edge_services_dns_stage.basic 11111111-1111-1111-1111-111111111111
```
