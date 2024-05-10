---
subcategory: "Instances"
page_title: "Scaleway: scaleway_instance_user_data"
---

# Resource: scaleway_instance_user_data

Creates and manages Scaleway Compute Instance User Data values.

User data is a key value store API you can use to provide data from and to your server without authentication. It is the mechanism by which a user can pass information contained in a local file to an Instance at launch time.

The typical use case is to pass something like a shell script or a configuration file as user data.

For more information about [user_data](https://developers.scaleway.com/en/products/instance/api/#patch-9ef3ec)  check our documentation guide [here](https://www.scaleway.com/en/docs/compute/instances/how-to/use-boot-modes/#how-to-use-cloud-init).

About cloud-init documentation please check this [link](https://cloudinit.readthedocs.io/en/latest/).

## Example Usage

### Basic

```terraform
variable user_data {
  type = map
  default = {
    "cloud-init" = <<-EOF
    #cloud-config
    apt-update: true
    apt-upgrade: true
    EOF
    "foo" = "bar"
  }
}

# User data with a single value
resource "scaleway_instance_user_data" "main" {
  server_id = scaleway_instance_server.main.id
  key = "foo"
  value = "bar"
}

# User Data with many keys.
resource scaleway_instance_user_data data {
  server_id = scaleway_instance_server.main.id
  for_each = var.user_data
  key = each.key
  value = each.value
}

resource "scaleway_instance_server" "main" {
  image = "ubuntu_focal"
  type  = "DEV1-S"
}
```

## Argument Reference

The following arguments are required:

- `server_id` - (Required) The ID of the server associated with.
- `key` - (Required) Key of the user data.
- `value` - (Required) Value associated with your key
- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the server should be created.

~> **Important:**   Use the `cloud-init` key to use [cloud-init](https://cloudinit.readthedocs.io/en/latest/) on your instance.
  You can define values using:
    - string
    - UTF-8 encoded file content using [file](https://www.terraform.io/language/functions/file)

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the instance's user data.

~> **Important:** Instances user data's IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{key}/{server_id}`, e.g. `fr-par-1/cloud-init/11111111-1111-1111-1111-111111111111`

## Import

User data can be imported using the `{zone}/{key}/{server_id}`, e.g.

```bash
$ terraform import scaleway_instance_user_data.main fr-par-1/cloud-init/11111111-1111-1111-1111-111111111111
```
