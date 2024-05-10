---
subcategory: "IAM"
page_title: "Scaleway: scaleway_iam_group"
---

# Resource: scaleway_iam_group

Creates and manages Scaleway IAM Groups.
For more information, see [the documentation](https://developers.scaleway.com/en/products/iam/api/v1alpha1/#groups-f592eb).

## Example Usage

### Basic

```terraform
resource "scaleway_iam_group" "basic" {
  name            = "iam_group_basic"
  description     = "basic description"
  application_ids = []
  user_ids        = []
}
```

### With applications

```terraform
resource "scaleway_iam_application" "app" {}

resource "scaleway_iam_group" "with_app" {
  name = "iam_group_with_app"
  application_ids = [
    scaleway_iam_application.app.id,
  ]
  user_ids = []
}
```

### With users

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

resource "scaleway_iam_group" "with_users" {
  name            = "iam_group_with_app"
  application_ids = []
  user_ids        = [for user in data.scaleway_iam_user.users : user.id]
}
```

## Argument Reference

- `name` - (Optional) The name of the IAM group.

- `description` - (Optional) The description of the IAM group.

- `tags` - (Optional) The tags associated with the group.

- `application_ids` - (Optional) The list of IDs of the applications attached to the group.

- `user_ids` - (Optional) The list of IDs of the users attached to the group.

- `external_membership` - (Optional) Manage membership externally. This make the resource ignore user_ids and application_ids. Should be used when using [iam_group_membership](iam_group_membership.md)

- `organization_id` - (Defaults to [provider](../index.md#organization_d) `organization_id`) The ID of the organization the group is associated with.

## Attributes Reference

No additional attributes are exported.

## Import

IAM groups can be imported using the `{id}`, e.g.

```bash
$ terraform import scaleway_iam_group.basic 11111111-1111-1111-1111-111111111111
```
