---
subcategory: "Secrets"
page_title: "Scaleway: scaleway_secret_version"
---

# Resource: scaleway_secret_version

Creates and manages Scaleway Secret Versions.
For more information, see [the documentation](https://developers.scaleway.com/en/products/secret_manager/api/v1alpha1/#secret-versions-079501).

## Example Usage

### Basic

```terraform
resource "scaleway_secret" "main" {
  name        = "foo"
  description = "barr"
  tags        = ["foo", "terraform"]
}

resource "scaleway_secret_version" "v1" {
  description = "version1"
  secret_id   = scaleway_secret.main.id
  data        = "my_new_secret"
}
```

## Argument Reference

The following arguments are supported:

- `secret_id` - (Required) The Secret ID associated wit the secret version.
- `data` - (Required) The data payload of the secret version. Must be no larger than 64KiB. (e.g. `my-secret-version-payload`). more on the [data section](#data)
- `description` - (Optional) Description of the secret version (e.g. `my-new-description`).
- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#regions)
  in which the resource exists.

### Data

Note: The `data` should be a base64 encoded string when sent from the API. **It is already handled by the provider so you don't need to encode it yourself.**

Updating `data` will force creating a new the secret version.

Be aware that this is a sensitive attribute. For more information, see [Sensitive Data in State](https://developer.hashicorp.com/terraform/language/state/sensitive-data).

~> **Important:**  This property is sensitive and will not be displayed in the plan.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `revision` - The revision for this Secret Version.
- `status` - The status of the Secret Version.
- `created_at` - Date and time of secret version's creation (RFC 3339 format).
- `updated_at` - Date and time of secret version's last update (RFC 3339 format).

## Import

The Secret Version can be imported using the `{region}/{id}/{revision}`, e.g.

~> **Important:** Be aware if you import with revision `latest` you will overwrite the version you used before.

```bash
$ terraform import scaleway_secret_version.main fr-par/11111111-1111-1111-1111-111111111111/2
```
