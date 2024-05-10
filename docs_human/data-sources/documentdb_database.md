---
subcategory: "Databases"
page_title: "Scaleway: scaleway_documentdb_database"
---

# scaleway_documentdb_database

Gets information about DocumentDB database. More on our official [site](https://www.scaleway.com/en/developers/api/document_db/)

## Example Usage

```hcl
# Get the database foobar hosted on instance id 11111111-1111-1111-1111-111111111111
data scaleway_documentdb_database main {
  instance_id = "11111111-1111-1111-1111-111111111111"
  name        = "foobar"
}
```

## Argument Reference

- `instance_id` - (Required) The DocumentDB instance ID.

- `name` - (Required) The name of the DocumentDB instance.

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `id` - The ID of the database.

~> **Important:** DocumentDB databases' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{instance-id}/{database-name}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111/database`

- `owner` - The name of the owner of the database.
- `managed` - Whether the database is managed or not.
- `size` - Size of the database (in bytes).
