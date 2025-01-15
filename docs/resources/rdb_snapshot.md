---
subcategory: "Databases"
page_title: "Scaleway: scaleway_rdb_snapshot"
---

# Resource: scaleway_rdb_snapshot

Creates and manages Scaleway RDB (Relational Database) Snapshots.
Snapshots are point-in-time backups of a database instance that can be used for recovery or duplication.
For more information, refer to [the API documentation](https://www.scaleway.com/en/developers/api/managed-database-postgre-mysql/).

## Example Usage

### Example Basic Snapshot

```terraform
resource "scaleway_rdb_instance" "main" {
  name              = "test-rdb-instance"
  node_type         = "db-dev-s"
  engine            = "PostgreSQL-15"
  is_ha_cluster     = false
  disable_backup    = true
  user_name         = "my_initial_user"
  password          = "thiZ_is_v&ry_s3cret"
  tags              = ["terraform-test", "scaleway_rdb_instance", "minimal"]
  volume_type       = "bssd"
  volume_size_in_gb = 10
}

resource "scaleway_rdb_snapshot" "test" {
  name        = "initial-snapshot"
  instance_id = scaleway_rdb_instance.main.id
  depends_on  = [scaleway_rdb_instance.main]
}
```

### Example with Expiration

```terraform
resource "scaleway_rdb_snapshot" "snapshot_with_expiration" {
  name        = "snapshot-with-expiration"
  instance_id = scaleway_rdb_instance.main.id
  expires_at  = "2025-01-31T00:00:00Z"
}
```

### Example with Multiple Snapshots

```terraform
resource "scaleway_rdb_snapshot" "daily_snapshot" {
  name        = "daily-backup"
  instance_id = scaleway_rdb_instance.main.id
  depends_on  = [scaleway_rdb_instance.main]
}

resource "scaleway_rdb_snapshot" "weekly_snapshot" {
  name        = "weekly-backup"
  instance_id = scaleway_rdb_instance.main.id
  expires_at  = "2025-02-07T00:00:00Z"
  depends_on  = [scaleway_rdb_instance.main]
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The name of the snapshot.
- `instance_id` - (Required) The UUID of the database instance for which the snapshot is created.
- `snapshot_id` - (Optional, ForceNew) The ID of an existing snapshot. This allows creating an instance from a specific snapshot ID. Conflicts with `engine`.
- `expires_at` - (Optional) Expiration date of the snapshot in ISO 8601 format (e.g., `2025-01-31T00:00:00Z`). If not set, the snapshot will not expire automatically.

### Additional Computed Attributes

In addition to the arguments above, the following attributes are exported:

- `id` - The unique ID of the snapshot.
- `created_at` - The timestamp when the snapshot was created, in ISO 8601 format.
- `updated_at` - The timestamp when the snapshot was last updated, in ISO 8601 format.
- `status` - The current status of the snapshot (e.g., `ready`, `creating`, `error`).
- `size` - The size of the snapshot in bytes.
- `node_type` - The type of the database instance for which the snapshot was created.
- `volume_type` - The type of volume used by the snapshot.

## Attributes Reference

- `region` - The region where the snapshot is stored. Defaults to the region set in the provider configuration.

## Import

RDB Snapshots can be imported using the `{region}/{snapshot_id}` format.

### Example:

```bash
terraform import scaleway_rdb_snapshot.example fr-par/11111111-1111-1111-1111-111111111111
```

## Limitations

- Snapshots are tied to the database instance and region where they are created.
- Expired snapshots are automatically deleted and cannot be restored.

## Notes

- Ensure the `instance_id` corresponds to an existing database instance.
- Use the `depends_on` argument when creating snapshots right after creating an instance to ensure proper dependency management.
