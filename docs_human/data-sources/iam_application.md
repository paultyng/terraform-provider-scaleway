---
subcategory: "IAM"
page_title: "Scaleway: scaleway_iam_application"
---

# scaleway_iam_application

Gets information about an existing IAM application.

## Example Usage

```hcl
# Get info by name
data "scaleway_iam_application" "find_by_name" {
  name = "foobar"
}
# Get info by application ID
data "scaleway_iam_application" "find_by_id" {
  application_id = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

- `name` - (Optional) The name of the IAM application.
  Only one of the `name` and `application_id` should be specified.

- `application_id` - (Optional) The ID of the IAM application.
  Only one of the `name` and `application_id` should be specified.

- `organization_id` - (Optional. Defaults to [provider](../index.md#organization_d) `organization_id`) The ID of the
  organization the application is associated with.

## Attribute Reference

Exported attributes are the ones from `iam_application` [resource](../resources/iam_application.md)
