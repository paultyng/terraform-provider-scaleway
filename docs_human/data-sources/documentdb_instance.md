---
subcategory: "Databases"
page_title: "Scaleway: scaleway_documentdb_instance"
---

# scaleway_documentdb_instance

Gets information about an DocumentDB instance. For further information see our [developers website](https://www.scaleway.com/en/developers/api/document_db/)

## Example Usage

```hcl
# Get info by name
data "scaleway_documentdb_instance" "db" {
  name = "foobar"
}

# Get info by instance ID
data "scaleway_documentdb_instance" "db" {
  instance_id = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

- `name` - (Optional) The name of the DocumentDB instance.
  Only one of `name` and `instance_id` should be specified.

- `instance_id` - (Optional) The DocumentDB instance ID.
  Only one of `name` and `instance_id` should be specified.

- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#zones) in which the DocumentDB instance exists.

- `organization_id` - (Defaults to [provider](../index.md#organization_id) `organization_id`) The ID of the organization the DocumentDB instance is in.

- `project_id` - (Optional) The ID of the project the DocumentDB instance is associated with.


## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `id` - The ID of the DocumentDB instance.

~> **Important:** DocumentDB instances' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`

Exported attributes are the ones from `scaleway_documentdb_instance` [resource](../resources/documentdb_instance.md)
