---
subcategory: "IAM"
page_title: "Scaleway: scaleway_iam_ssh_key"
---

# Resource: scaleway_iam_ssh_key

Creates and manages Scaleway IAM SSH Keys.
For more information,
see [the documentation](https://developers.scaleway.com/en/products/iam/api/v1alpha1/#ssh-keys-d8ccd4).

## Example Usage

```terraform
resource "scaleway_iam_ssh_key" "main" {
  name       = "main"
  public_key = "<YOUR-PUBLIC-SSH-KEY>"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The name of the SSH key.
- `public_key` - (Required) The public SSH key to be added.
- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the SSH key is
  associated with.
- `disabled` - (Optional) The SSH key status.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the SSH public key.
- `fingerprint` - The fingerprint of the iam SSH key.
- `organization_id` - The ID of the organization the SSH key is associated with.
- `created_at` - The date and time of the creation of the SSH key.
- `updated_at` - The date and time of the last update of the SSH key.

## Import

SSH keys can be imported using the `id`, e.g.

```bash
$ terraform import scaleway_iam_ssh_key.main 11111111-1111-1111-1111-111111111111
```
