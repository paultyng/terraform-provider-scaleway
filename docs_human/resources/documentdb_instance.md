---
subcategory: "Databases"
page_title: "Scaleway: scaleway_documentdb_instance"
---

# Resource: scaleway_documentdb_instance

Creates and manages Scaleway Database Instances.
For more information, see [the documentation](https://www.scaleway.com/en/developers/api/document_db/).

## Example Usage

### Example Basic

```terraform
resource "scaleway_documentdb_instance" "main" {
  name              = "test-documentdb-instance-basic"
  node_type         = "docdb-play2-pico"
  engine            = "FerretDB-1"
  user_name         = "my_initial_user"
  password          = "thiZ_is_v&ry_s3cret"
  tags              = ["terraform-test", "scaleway_documentdb_instance", "minimal"]
  volume_size_in_gb = 20
}
```

## Argument Reference

The following arguments are supported:

- `node_type` - (Required) The type of database instance you want to create (e.g. `docdb-play2-pico`).

~> **Important:** Updates to `node_type` will upgrade the Database Instance to the desired `node_type` without any
interruption. Keep in mind that you cannot downgrade a Database Instance.

- `engine` - (Required) Database Instance's engine version (e.g. `FerretDB-1`).

~> **Important:** Updates to `engine` will recreate the Database Instance.

- `volume_type` - (Optional, default to `lssd`) Type of volume where data are stored (`bssd` or `lssd`).

- `volume_size_in_gb` - (Optional) Volume size (in GB) when `volume_type` is set to `bssd`.

- `user_name` - (Optional) Identifier for the first user of the database instance.

~> **Important:** Updates to `user_name` will recreate the Database Instance.

- `password` - (Optional) Password for the first user of the database instance.

- `is_ha_cluster` - (Optional) Enable or disable high availability for the database instance.

- `telemetry_enabled` - (Optional) Enable telemetry to collects basic anonymous usage data and sends them to FerretDB telemetry service. More about the telemetry [here](https://docs.ferretdb.io/telemetry/#configure-telemetry).

~> **Important:** Updates to `is_ha_cluster` will recreate the Database Instance.

- `name` - (Optional) The name of the Database Instance.

- `tags` - (Optional) The tags associated with the Database Instance.

- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#regions)
  in which the Database Instance should be created.

- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the Database
  Instance is associated with.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the Database Instance.

~> **Important:** Database instances' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they
are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`

- `organization_id` - The organization ID the Database Instance is associated with.

## Import

Database Instance can be imported using the `{region}/{id}`, e.g.

```bash
$ terraform import scaleway_documentdb_instance.db fr-par/11111111-1111-1111-1111-111111111111
```
