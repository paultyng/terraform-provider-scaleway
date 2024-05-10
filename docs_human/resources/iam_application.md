---
subcategory: "IAM"
page_title: "Scaleway: scaleway_iam_application"
---

# Resource: scaleway_iam_application

Creates and manages Scaleway IAM Applications. For more information, see [the documentation](https://developers.scaleway.com/en/products/iam/api/v1alpha1/#applications-83ce5e).

## Example Usage

```terraform
resource "scaleway_iam_application" "main" {
  name        = "My application"
  description = "a description"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The name of the iam application.
- `description` - (Optional) The description of the iam application.
- `tags` - (Optional) The tags associated with the application.
- `organization_id` - (Defaults to [provider](../index.md#organization_d) `organization_id`) The ID of the organization the application is associated with.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the application.
- `created_at` - The date and time of the creation of the application.
- `updated_at` - The date and time of the last update of the application.
- `editable` - Whether the application is editable.

## Import

Applications can be imported using the `{id}`, e.g.

```bash
$ terraform import scaleway_iam_application.main 11111111-1111-1111-1111-111111111111
```
