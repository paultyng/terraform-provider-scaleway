---
subcategory: "Databases"
page_title: "Scaleway: scaleway_documentdb_database"
---

# Resource: scaleway_documentdb_database

Creates and manages Scaleway DocumentDB database.
For more information, see [the documentation](https://www.scaleway.com/en/developers/api/document_db).

## Example Usage

### Basic

```terraform
resource "scaleway_documentdb_database" "main" {
  instance_id    = "11111111-1111-1111-1111-111111111111"
  name           = "my-new-database"
}
```

## Argument Reference

The following arguments are supported:

- `instance_id` - (Required) UUID of the documentdb instance.

~> **Important:** Updates to `instance_id` will recreate the Database.

- `name` - (Required) Name of the database (e.g. `my-new-database`).

- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#regions) in which the resource exists.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the database, which is of the form `{region}/{id}/{DBNAME}` e.g. `fr-par/11111111-1111-1111-1111-111111111111/mydb`
- `owner` - The name of the owner of the database.
- `managed` - Whether the database is managed or not.
- `size` - Size in gigabytes of the database.

## Import

DocumentDB Database can be imported using the `{region}/{id}/{DBNAME}`, e.g.

```bash
$ terraform import scaleway_documentdb_database.mydb fr-par/11111111-1111-1111-1111-111111111111/mydb
```
