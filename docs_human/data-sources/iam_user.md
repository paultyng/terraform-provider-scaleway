---
subcategory: "IAM"
page_title: "Scaleway: scaleway_iam_user"
---

# scaleway_iam_user

Use this data source to get information on an existing IAM user based on its ID or email address.
For more information,
see [the documentation](https://developers.scaleway.com/en/products/iam/api/v1alpha1/#users-06bdcf).

## Example Usage

```hcl
# Get info by user id
data "scaleway_iam_user" "find_by_id" {
  user_id = "11111111-1111-1111-1111-111111111111"
}
# Get info by email address
data "scaleway_iam_user" "find_by_email" {
  email = "foo@bar.com"
}
```

## Argument Reference

- `email` - (Optional) The email address of the IAM user. Only one of the `email` and `user_id` should be specified.
- `user_id` - (Optional) The ID of the IAM user. Only one of the `email` and `user_id` should be specified.
- `organization_id` - (Optional. Defaults to [provider](../index.md#organization_d) `organization_id`) The ID of the
  organization the user is associated with.

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `id` - The ID of the IAM user.
