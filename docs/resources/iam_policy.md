---
page_title: "Scaleway: scaleway_iam_policy"
description: |-
Manages Scaleway IAM Policies.
---

# scaleway_iam_policy

| WARNING: This resource is in beta version. If your are in the beta group, please set the variable `SCW_ENABLE_BETA=true` in your `env` in order to use this resource. |
|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------|

Creates and manages Scaleway IAM Policies. For more information, see [the documentation](https://developers.scaleway.com/en/products/iam/api/v1alpha1/#policies-54b8a7).

## Example Usage

### Create a policy for an organization's project

```hcl
provider scaleway {
  organization_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

data scaleway_account_project "default" {
  name = "default"
}

resource scaleway_iam_application "app" {
  name = "my app"
}

resource scaleway_iam_policy "object_read_only" {
  name = "my policy"
  description = "gives app readonly access to object storage in project"
  application_id = scaleway_iam_application.app.id
  rule {
    project_ids = [data.scaleway_account_project.default.id]
    permission_set_names = ["ObjectStorageReadOnly"]
  }
}
```

### Create a permission for multiple users using a group

```hcl
locals {
  users = [
  "user1@mail.com",
  "user2@mail.com",
  ]
  project_name = "default"
}

data scaleway_account_project project {
  name = local.project_name
}

data "scaleway_iam_user" "users" {
  for_each = toset(local.users)
  email = each.value
}

resource "scaleway_iam_group" "with_users" {
  name = "developers"
  user_ids = [for user in data.scaleway_iam_user.users : user.id]
}

resource scaleway_iam_policy "iam_tf_storage_policy" {
  name = "developers permissions"
  group_id = scaleway_iam_group.with_users.id
  rule {
    project_ids = [data.scaleway_account_project.project.id]
    permission_set_names = ["InstancesReadOnly"]
  }
}

```

## Arguments Reference

The following arguments are supported:

- `name` - .The name of the iam policy.
- `description` - The description of the iam policy.
- `organization_id` - (Defaults to [provider](../index.md#organization_d) `organization_id`) The ID of the organization the policy is associated with.
- `user_id` - ID of the User the policy will be linked to
- `group_id` - ID of the Group the policy will be linked to
- `application_id` - ID of the Application the policy will be linked to
- `no_principal` - If the policy doesn't apply to a principal.

~> **Important** Only one of `user_id`, `group_id`, `application_id` and `no_principal`  may be set.

- `rule` - List of rules in the policy.
    - `organization_id` - ID of organization scoped to the rule.
    - `project_ids` - List of project IDs scoped to the rule.

  ~> **Important** One of `organization_id` or `project_ids`  must be set per rule.

    - `permission_set_names` - Names of permission sets bound to the rule.

  **_TIP:_**  You can use the Scaleway CLI to list the permissions details. e.g:

```shell
  $ SCW_ENABLE_BETA=1 scw iam permission-set list
```

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `created_at` - The date and time of the creation of the policy.
- `updated_at` - The date and time of the last update of the policy.
- `editable` - Whether the policy is editable.

## Import

Policies can be imported using the `{id}`, e.g.

```bash
$ terraform import scaleway_iam_policy.main 11111111-1111-1111-1111-111111111111
```
