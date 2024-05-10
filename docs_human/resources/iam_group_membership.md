---
subcategory: "IAM"
page_title: "Scaleway: scaleway_iam_group_membership"
---

# Resource: scaleway_iam_group_membership

Add members to an IAM group.
For more information, see [the documentation](https://developers.scaleway.com/en/products/iam/api/v1alpha1/#groups-f592eb).

## Example Usage

### Application Membership

```terraform
resource "scaleway_iam_group" "group" {
  name                = "my_group"
  external_membership = true
}

resource "scaleway_iam_application" "app" {
  name = "my_app"
}

resource "scaleway_iam_group_membership" "member" {
  group_id       = scaleway_iam_group.group.id
  application_id = scaleway_iam_application.app.id
}
```

### Users membership

```terraform
locals {
  users = toset([
    "user1@mail.com",
    "user2@mail.com"
  ])
}

data "scaleway_iam_user" "users" {
  for_each = local.users
  email    = each.value
}

resource "scaleway_iam_group" "group" {
  name                = "my_group"
  external_membership = true
}

resource "scaleway_iam_group_membership" "members" {
  for_each = data.scaleway_iam_user.users
  group_id = scaleway_iam_group.group.id
  user_id = each.value.id
}
```

## Argument Reference

- `group_id` - (Required) ID of the group to add members to.

- `application_id` - (Optional) The ID of the application that will be added to the group.

- `user_id` - (Optional) The ID of the user that will be added to the group

- ~> Only one of `application_id` or `user_id` must be specified

## Attributes Reference

No additional attributes are exported.

## Import

IAM group memberships can be imported using two format:

- For user: `{group_id}/user/{user_id}`
- For application: `{group_id}/app/{application_id}`

```bash
$ terraform import scaleway_iam_group_membership.app 11111111-1111-1111-1111-111111111111/app/11111111-1111-1111-1111-111111111111
```
