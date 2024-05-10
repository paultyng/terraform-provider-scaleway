---
subcategory: "Databases"
page_title: "Scaleway: scaleway_rdb_database_backup"
---

# scaleway_rdb_database_backup

Gets information about an RDB backup.

## Example Usage

```hcl
data scaleway_rdb_database_backup find_by_name {
  name        = "mybackup"
}

data scaleway_rdb_database_backup find_by_name_and_instance {
  name        = "mybackup"
  instance_id = "11111111-1111-1111-1111-111111111111"
}

data scaleway_rdb_database_backup find_by_id {
  backup_id = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

- `backup_id` - (Optional) The RDB backup ID.
  Only one of the `name` and `backup_id` should be specified.

- `instance_id` - (Optional) The RDB instance ID.

- `name` - (Optional) The name of the RDB instance.
  Only one of the `name` and `backup_id` should be specified.

- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#regions) in which the Database Backup is associated with.

- `project_id` - (Optional) The ID of the project the Database Backup is associated with.

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `id` - The ID of the backup.

~> **Important:** RDB databases backups' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`
