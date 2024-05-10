---
subcategory: "Databases"
page_title: "Scaleway: scaleway_rdb_instance"
---

# Resource: scaleway_rdb_instance

Creates and manages Scaleway Database Instances.
For more information, see [the documentation](https://developers.scaleway.com/en/products/rdb/api).

## Example Usage

### Example Basic

```terraform
resource "scaleway_rdb_instance" "main" {
  name           = "test-rdb"
  node_type      = "DB-DEV-S"
  engine         = "PostgreSQL-11"
  is_ha_cluster  = true
  disable_backup = true
  user_name      = "my_initial_user"
  password       = "thiZ_is_v&ry_s3cret"
}
```

### Example with Settings

```terraform
resource "scaleway_rdb_instance" "main" {
  name           = "test-rdb"
  node_type      = "db-dev-s"
  disable_backup = true
  engine         = "MySQL-8"
  user_name      = "my_initial_user"
  password       = "thiZ_is_v&ry_s3cret"
  init_settings  = {
    "lower_case_table_names" = 1
  }
  settings = {
    "max_connections" = "350"
  }
}
```

### Example with backup schedule

```terraform
resource "scaleway_rdb_instance" "main" {
  name          = "test-rdb"
  node_type     = "DB-DEV-S"
  engine        = "PostgreSQL-11"
  is_ha_cluster = true
  user_name     = "my_initial_user"
  password      = "thiZ_is_v&ry_s3cret"

  disable_backup            = false
  backup_schedule_frequency = 24 # every day
  backup_schedule_retention = 7  # keep it one week
}
```

### Examples of endpoints configuration

RDB Instances can have a maximum of 1 public endpoint and 1 private endpoint. It can have both, or none.

#### 1 static private network endpoint

```terraform
resource "scaleway_vpc_private_network" "pn" {
  ipv4_subnet {
    subnet = "172.16.20.0/22"
  }
}

resource "scaleway_rdb_instance" "main" {
  node_type         = "db-dev-s"
  engine            = "PostgreSQL-11"
  private_network {
    pn_id  = scaleway_vpc_private_network.pn.id
    ip_net = "172.16.20.4/22"   # IP address within a given IP network
    # enable_ipam = false
  }
}
```

#### 1 IPAM private network endpoint + 1 public endpoint

```terraform
resource "scaleway_vpc_private_network" "pn" {}

resource "scaleway_rdb_instance" "main" {
  node_type      = "DB-DEV-S"
  engine         = "PostgreSQL-11"
  private_network {
    pn_id = scaleway_vpc_private_network.pn.id
    enable_ipam = true
  }
  load_balancer {}
}
```

#### Default: 1 public endpoint

```terraform
resource "scaleway_rdb_instance" "main" {
  node_type         = "db-dev-s"
  engine            = "PostgreSQL-11"
}  
```

-> If nothing is defined, your instance will have a default public load-balancer endpoint

## Argument Reference

The following arguments are supported:

- `node_type` - (Required) The type of database instance you want to create (e.g. `db-dev-s`).

~> **Important:** Updates to `node_type` will upgrade the Database Instance to the desired `node_type` without any
interruption. Keep in mind that you cannot downgrade a Database Instance.

~> **Important:** Once your instance reaches `disk_full` status, if you are using `lssd` storage, you should upgrade the node_type,
and if you are using `bssd` storage, you should increase the volume size before making any other change to your instance.

- `engine` - (Required) Database Instance's engine version (e.g. `PostgreSQL-11`).

~> **Important:** Updates to `engine` will recreate the Database Instance.

- `volume_type` - (Optional, default to `lssd`) Type of volume where data are stored (`bssd`, `lssd` or `sbs_5k`).

- `volume_size_in_gb` - (Optional) Volume size (in GB). Cannot be used when `volume_type` is set to `lssd`.

~> **Important:** Once your instance reaches `disk_full` status, you should increase the volume size before making any other change to your instance.

- `user_name` - (Optional) Identifier for the first user of the database instance.

~> **Important:** Updates to `user_name` will recreate the Database Instance.

- `password` - (Optional) Password for the first user of the database instance.

- `is_ha_cluster` - (Optional) Enable or disable high availability for the database instance.

~> **Important:** Updates to `is_ha_cluster` will recreate the Database Instance.

- `name` - (Optional) The name of the Database Instance.

- `tags` - (Optional) The tags associated with the Database Instance.

- `region` - (Defaults to [provider](../index.md#arguments-reference) `region`) The [region](../guides/regions_and_zones.md#regions)
  in which the Database Instance should be created.

- `project_id` - (Defaults to [provider](../index.md#arguments-reference) `project_id`) The ID of the project the Database
  Instance is associated with.

### Backups

- `disable_backup` - (Optional) Disable automated backup for the database instance.

- `backup_schedule_frequency` - (Optional) Backup schedule frequency in hours.

- `backup_schedule_retention` - (Optional) Backup schedule retention in days.

- `backup_same_region` - (Optional) Boolean to store logical backups in the same region as the database instance.

### Settings

- `settings` - (Optional) Map of engine settings to be set. Using this option will override default config.

- `init_settings` - (Optional) Map of engine settings to be set at database initialisation.

~> **Important:** Updates to `init_settings` will recreate the Database Instance.

Please consult the [GoDoc](https://pkg.go.dev/github.com/scaleway/scaleway-sdk-go@v1.0.0-beta.9/api/rdb/v1#EngineVersion) to list all available `settings` and `init_settings` for the `node_type` of your convenience.

### Endpoints

- `private_network` - List of private networks endpoints of the database instance.

    - `pn_id` - (Required) The ID of the private network.
    - `ip_net` - (Optional) The IP network address within the private subnet. This must be an IPv4 address with a CIDR notation. If not set, The IP network address within the private subnet is determined by the IP Address Management (IPAM) service.
    - `enable_ipam` - (Optional) If true, the IP network address within the private subnet is determined by the IP Address Management (IPAM) service.
  
~> **NOTE:** Please calculate your host IP using [cidrhost](https://developer.hashicorp.com/terraform/language/functions/cidrhost). Otherwise, let IPAM service
handle the host IP on the network.

~> **Important:** Updates to `private_network` will recreate the Instance's endpoint

- `load_balancer` - (Optional) List of load balancer endpoints of the database instance. A load-balancer endpoint will be set by default if no private network is.
This block must be defined if you want a public endpoint in addition to your private endpoint.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the Database Instance.

~> **Important:** Database instances' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they
are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`

- `endpoint_ip` - (Deprecated) The IP of the Database Instance.
- `endpoint_port` - (Deprecated) The port of the Database Instance.
- `read_replicas` - List of read replicas of the database instance.
    - `ip` - IP of the replica.
    - `port` - Port of the replica.
    - `name` - Name of the replica.
- `load_balancer` - List of load balancer endpoints of the database instance.
    - `endpoint_id` - The ID of the endpoint of the load balancer.
    - `ip` - IP of the endpoint.
    - `port` - Port of the endpoint.
    - `name` - Name of the endpoint.
    - `hostname` - Name of the endpoint.
- `private_network` - List of private networks endpoints of the database instance.
    - `endpoint_id` - The ID of the endpoint.
    - `ip` - IPv4 address on the network.
    - `port` - Port in the Private Network.
    - `name` - Name of the endpoint.
    - `hostname` - Hostname of the endpoint.
- `certificate` - Certificate of the database instance.
- `organization_id` - The organization ID the Database Instance is associated with.

## Limitations

The Managed Database product is only compliant with the private network in the default availability zone (AZ).
i.e. `fr-par-1`, `nl-ams-1`, `pl-waw-1`. To learn more, read our
section [How to connect a PostgreSQL and MySQL Database Instance to a Private Network](https://www.scaleway.com/en/docs/managed-databases/postgresql-and-mysql/how-to/connect-database-private-network/)

## Import

Database Instance can be imported using the `{region}/{id}`, e.g.

```bash
$ terraform import scaleway_rdb_instance.rdb01 fr-par/11111111-1111-1111-1111-111111111111
```
