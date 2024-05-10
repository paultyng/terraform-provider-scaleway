---
subcategory: "Account"
page_title: "Scaleway: scaleway_account_ssh_key"
---

# Resource: scaleway_account_ssh_key

Manages user SSH keys to access servers provisioned on Scaleway.

~> **Important:**  The resource `scaleway_account_ssh_key` has been deprecated and will no longer be supported. Instead, use `scaleway_iam_ssh_key`.

## Example Usage

```terraform
resource "scaleway_account_ssh_key" "main" {
    name 	   = "main"
    public_key = "<YOUR-PUBLIC-SSH-KEY>"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The name of the SSH key.
- `public_key` - (Required) The public SSH key to be added.
- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the SSH key is associated with.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the SSH key (UUID format).
- `organization_id` - The organization ID the SSH key is associated with.

## Import

SSH keys can be imported using the `id`, e.g.

```bash
$ terraform import scaleway_account_ssh_key.main 11111111-1111-1111-1111-111111111111
```
