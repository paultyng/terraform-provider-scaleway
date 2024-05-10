---
subcategory: "IAM"
page_title: "Scaleway: scaleway_iam_ssh_key"
---

# scaleway_iam_ssh_key

Use this data source to get SSH key information based on its ID or name.

## Example Usage

```hcl
# Get info by SSH key name
data "scaleway_iam_ssh_key" "my_key" {
  name  = "my-key-name"
}

# Get info by SSH key id
data "scaleway_iam_ssh_key" "my_key" {
  ssh_key_id = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

- `name` - The SSH key name. Only one of `name` and `ssh_key_id` should be specified.
- `ssh_key_id` - The SSH key id. Only one of `name` and `ssh_key_id` should be specified.
- `project_id` (Optional. Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the SSH
  key is associated with.

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `id` - The ID of the SSH public key.
- `public_key` - The SSH public key string
- `organization_id` - The ID of the organization the SSH key is associated with.
- `created_at` - The date and time of the creation of the SSH key.
- `updated_at` - The date and time of the last update of the SSH key.
- `disabled` - The SSH key status.
