---
subcategory: "Databases"
page_title: "Scaleway: scaleway_rdb_user"
---

# Resource: scaleway_rdb_user

Creates and manages Scaleway Database Users.
For more information, see [the documentation](https://developers.scaleway.com/en/products/rdb/api).

## Example Usage

### Basic

```terraform
resource "random_password" "db_password" {
  length  = 16
  special = true
}

resource "scaleway_rdb_user" "db_admin" {
  instance_id = scaleway_rdb_instance.main.id
  name        = "devtools"
  password    = random_password.db_password.result
  is_admin    = true
}
```

## Argument Reference

The following arguments are supported:

- `instance_id` - (Required) UUID of the rdb instance.

~> **Important:** Updates to `instance_id` will recreate the Database User.

- `name` - (Required) Database User name.

~> **Important:** Updates to `name` will recreate the Database User.

- `password` - (Required) Database User password.

- `is_admin` - (Optional) Grant admin permissions to the Database User.

- `region` - The Scaleway region this resource resides in.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the user, which is of the form `{region}/{instance_id}/{user_name}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111/admin`

## Import

Database User can be imported using `{region}/{instance_id}/{user_name}`, e.g.

```bash
$ terraform import scaleway_rdb_user.admin fr-par/11111111-1111-1111-1111-111111111111/admin
```
