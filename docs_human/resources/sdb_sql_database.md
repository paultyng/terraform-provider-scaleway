---
subcategory: "Databases"
page_title: "Scaleway: scaleway_sdb_sql_database"
---

# Resource: scaleway_sdb_sql_database

Creates and manages Scaleway Serverless SQL Databases. For more information, see [the documentation](https://www.scaleway.com/en/developers/api/serverless-databases/).

## Example Usage

### Basic

```hcl
resource scaleway_sdb_sql_database "database" {
  name = "my-database"
  min_cpu = 0
  max_cpu = 8
}
```

### With IAM Application

This example creates an [IAM application](https://www.scaleway.com/en/docs/identity-and-access-management/iam/concepts/#application) and a secret key used to connect to the DATABASE.
For more information, see [How to connect to a Serverless SQL Database](https://www.scaleway.com/en/docs/serverless/sql-databases/how-to/connect-to-a-database/)

```hcl
data scaleway_account_project "default" {
  name = "default"
}

resource scaleway_iam_application "app" {
  name = "my app"
}

resource scaleway_iam_policy "db_access" {
  name = "my policy"
  description = "gives app access to serverless database in project"
  application_id = scaleway_iam_application.app.id
  rule {
    project_ids = [data.scaleway_account_project.default.id]
    permission_set_names = ["ServerlessSQLDatabaseReadWrite"]
  }
}

resource scaleway_iam_api_key "api_key" {
  application_id = scaleway_iam_application.app.id
}

resource scaleway_sdb_sql_database "database" {
  name = "my-database"
  min_cpu = 0
  max_cpu = 8
}

output "database_connection_string" {
  // Output as an example, you can give this string to your application
  value = format("postgres://%s:%s@%s",
    scaleway_iam_application.app.id,
    scaleway_iam_api_key.api_key.secret_key,
    trimprefix(scaleway_sdb_sql_database.database.endpoint, "postgres://"),
  )
  sensitive = true
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) Name of the database (e.g. `my-new-database`).

    ~> **Important:** Updates to `name` will recreate the database.

- `min_cpu` - (Optional) The minimum number of CPU units for your database. Defaults to 0.
- `max_cpu` - (Optional) The maximum number of CPU units for your database. Defaults to 15.

- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#regions) in which the resource exists.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the database, which is of the form `{region}/{id}` e.g. `fr-par/11111111-1111-1111-1111-111111111111`
- `endpoint` - Endpoint of the database

## Import

Serverless SQL Database can be imported using the `{region}/{id}`, e.g.

```bash
$ terraform import scaleway_sdb_sql_database.database fr-par/11111111-1111-1111-1111-111111111111
```
