---
subcategory: "Databases"
page_title: "Scaleway: scaleway_rdb_privilege"
---

# scaleway_rdb_privilege

Gets information about the privilege on RDB database.

## Example Usage

```hcl
# Get the database privilege for the user "my-user" on the database "my-database" hosted on instance id 11111111-1111-1111-1111-111111111111 and on the default region. e.g: fr-par
data "scaleway_rdb_privilege" "main" {
  instance_id   = "11111111-1111-111111111111"
  user_name     = "my-user"
  database_name = "my-database"
}
```

## Argument Reference

- `instance_id` - (Required) The RDB instance ID.

- `user_name` - (Required) The user name.

- `database_name` - (Required) The database name.

- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#regions) in which the resource exists.

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `id` - The ID of the user privileges.

~> **Important:** RDB databases user privileges' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{instance-id}/{database}/{user-name}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111/database/user`

- `permission` - The permission for this user on the database. Possible values are `readonly`, `readwrite`, `all`
  , `custom` and `none`.
