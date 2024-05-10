---
subcategory: "IAM"
page_title: "Scaleway: scaleway_iam_api_key"
---

# Resource: scaleway_iam_api_key

Creates and manages Scaleway IAM API Keys. For more information, please
check [the documentation](https://developers.scaleway.com/en/products/iam/api/v1alpha1/#api-keys-3665ae)

## Example Usage

```terraform
resource "scaleway_iam_application" "ci_cd" {
  name = "My application"
}

resource "scaleway_iam_api_key" "main" {
  application_id = scaleway_iam_application.main.id
  description    = "a description"
}
```

## Argument Reference

The following arguments are supported:

- `description`: (Optional) The description of the iam api key.
- `application_id`: (Optional) ID of the application attached to the api key.
  Only one of the `application_id` and `user_id` should be specified.
- `user_id` - (Optional) ID of the user attached to the api key.
  Only one of the `application_id` and `user_id` should be specified.
- `expires_at` - (Optional) The date and time of the expiration of the iam api key. Please note that in case of change,
  the resource will be recreated.
- `default_project_id` - (Optional) The default project ID to use with object storage.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the API key, which is the access key.
- `created_at` - The date and time of the creation of the iam api key.
- `updated_at` - The date and time of the last update of the iam api key.
- `editable` - Whether the iam api key is editable.
- `access_key` - The access key of the iam api key.
- `secret_key`: The secret Key of the iam api key.
- `creation_ip` - The IP Address of the device which created the API key.

## Import

Api keys can be imported using the `{id}`, e.g.

```bash
$ terraform import scaleway_iam_api_key.main 11111111111111111111
```
