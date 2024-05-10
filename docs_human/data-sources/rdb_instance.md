---
subcategory: "Databases"
page_title: "Scaleway: scaleway_rdb_instance"
---

# scaleway_rdb_instance

Gets information about an RDB instance. For further information see our [developers website](https://developers.scaleway.com/en/products/rdb/api/#database-instance)

## Example Usage

```hcl
# Get info by name
data "scaleway_rdb_instance" "my_instance" {
  name = "foobar"
}

# Get info by instance ID
data "scaleway_rdb_instance" "my_instance" {
  instance_id = "11111111-1111-1111-1111-111111111111"
}

# Get other attributes
output "load_balancer_ip_addr" {
  description = "IP address of load balancer"
  value = data.scaleway_rdb_instance.my_instance.load_balancer.0.ip
}
```

## Argument Reference

- `name` - (Optional) The name of the RDB instance.
  Only one of `name` and `instance_id` should be specified.

- `instance_id` - (Optional) The RDB instance ID.
  Only one of `name` and `instance_id` should be specified.

- `project_id` - (Optional) The ID of the project the RDB instance is in. Can be used to filter instances when using `name`.

- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#zones) in which the RDB instance exists.

- `organization_id` - (Defaults to [provider](../index.md#organization_id) `organization_id`) The ID of the organization the RDB instance is in.

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `id` - The ID of the RDB instance.

~> **Important:** RDB instances' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`

Exported attributes are the ones from `scaleway_rdb_instance` [resource](../resources/rdb_instance.md)
