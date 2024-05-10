---
subcategory: "Secrets"
page_title: "Scaleway: scaleway_secret_version"
---

# scaleway_secret_version

Gets information about Scaleway a Secret Version.
For more information, see [the documentation](https://developers.scaleway.com/en/products/secret_manager/api/v1alpha1/#secret-versions-079501).

## Examples

### Basic

```hcl
resource "scaleway_secret" "main" {
  name        = "fooii"
  description = "barr"
}

resource "scaleway_secret_version" "main" {
  description = "your description"
  secret_id   = scaleway_secret.main.id
  data        = "your_secret"
}

data "scaleway_secret_version" "data_by_secret_id" {
  secret_id  = scaleway_secret.main.id
  revision   = "1"
  depends_on = [scaleway_secret_version.main]
}

data "scaleway_secret_version" "data_by_secret_name" {
  secret_name = scaleway_secret.main.name
  revision    = "1"
  depends_on  = [scaleway_secret_version.main]
}

#Output Sensitive data
output "scaleway_secret_access_payload" {
  value = data.scaleway_secret_version.data_by_secret_name.data
}

#Output Sensitive data
output "scaleway_secret_access_payload_by_id" {
  value = data.scaleway_secret_version.data_by_secret_id.data
}
```

## Arguments Reference

The following arguments are supported:

- `secret_id` - (Optional) The Secret ID associated wit the secret version.
  Only one of `secret_id` and `secret_name` should be specified.

- `secret_name` - (Optional) The Name of Secret associated wit the secret version.
  Only one of `secret_id` and `secret_name` should be specified.

- `revision` - The revision for this Secret Version.

- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#regions)
  in which the resource exists.

- `project_id` - (Optional) The ID of the project the Secret version is associated with.

## Data

Note: This Data Source give you **access** to the secret payload encoded en base64.

Be aware that this is a sensitive attribute. For more information,
see [Sensitive Data in State](https://developer.hashicorp.com/terraform/language/state/sensitive-data).

~> **Important:**  This property is sensitive and will not be displayed in the plan.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `description` - (Optional) Description of the secret version (e.g. `my-new-description`).
- `data` - The data payload of the secret version. more on the [data section](#data)
- `status` - The status of the Secret Version.
- `created_at` - Date and time of secret version's creation (RFC 3339 format).
- `updated_at` - Date and time of secret version's last update (RFC 3339 format).

Exported attributes are the ones from `scaleway_secret_version` [resource](../resources/secret_version.md)
