---
subcategory: "Databases"
page_title: "Scaleway: scaleway_documentdb_privilege"
---

# Resource: scaleway_documentdb_privilege

Create and manage Scaleway DocumentDB database privilege.
For more information, see [the documentation](https://www.scaleway.com/en/developers/api/document_db/).

## Example Usage

```terraform
resource "scaleway_documentdb_privilege" "main" {
  instance_id   = "11111111-1111-1111-1111-111111111111"
  user_name     = "my-db-user"
  database_name = "my-db-name"
  permission    = "all"
}

```

## Argument Reference

The following arguments are supported:

- `instance_id` - (Required) UUID of the rdb instance.

- `user_name` - (Required) Name of the user (e.g. `my-db-user`).

- `database_name` - (Required) Name of the database (e.g. `my-db-name`).

- `permission` - (Required) Permission to set. Valid values are `readonly`, `readwrite`, `all`, `custom` and `none`.

- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#regions) in which the resource exists.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the user privileges, which is of the form `{region}/{instance_id}/{database_name}/{user_name}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111/database_name/foo`

## Import

The user privileges can be imported using the `{region}/{instance_id}/{database_name}/{user_name}`, e.g.

```bash
$ terraform import scaleway_documentdb_privilege.o fr-par/11111111-1111-1111-1111-111111111111/database_name/foo
```
