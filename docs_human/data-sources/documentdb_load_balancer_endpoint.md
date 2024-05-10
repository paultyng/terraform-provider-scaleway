---
subcategory: "Databases"
page_title: "Scaleway: scaleway_documentdb_load_balancer_endpoint"
---

# scaleway_documentdb_load_balancer_endpoint

Gets information about an DocumentDB load balancer endpoint.

## Example Usage

```terraform
# Get info by instance name
data "scaleway_documentdb_load_balancer_endpoint" "my_endpoint" {
  instance_name = "foobar"
}

# Get info by instance ID
data "scaleway_documentdb_load_balancer_endpoint" "my_endpoint" {
  instance_id = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

- `instance_name` - (Optional) The DocumentDB Instance Name on which the endpoint is attached. Only one of `instance_name` and `instance_id` should be specified.
- `instance_id` - (Optional) The DocumentDB Instance on which the endpoint is attached. Only one of `instance_name` and `instance_id` should be specified.
- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#zones) in which the DocumentDB endpoint exists.
- `project_id` - (Optional) The ID of the project the DocumentDB endpoint is associated with.

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `id` - The ID of the DocumentDB endpoint.

~> **Important:** DocumentDB endpoints' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`

- `ip` - The IP of your load balancer service.
- `port` - The port of your load balancer service.
- `name` - The name of your load balancer service.
- `hostname` - The hostname of your endpoint.
