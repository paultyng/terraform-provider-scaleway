---
subcategory: "Functions"
page_title: "Scaleway: scaleway_function"
---

# scaleway_function

Gets information about a function.

## Example Usage

```terraform
// Get info by function name
data "scaleway_function" "my_function" {
  name = "my-namespace-name"
  namespace_id = "11111111-1111-1111-1111-111111111111"
}

// Get info by function ID
data "scaleway_function" "my_function" {
  function_id = "11111111-1111-1111-1111-111111111111"
  namespace_id = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

- `namespace_id` - (Required) The namespace id associated with this function.
- `name` - (Optional) The function name. Only one of `name` and `namespace_id` should be specified.
- `function_id` - (Optional) The function id. Only one of `name` and `function_id` should be specified.
- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#regions) in which the function exists.
- `project_id` - (Optional) The ID of the project the function is associated with.

## Attributes Reference

Exported attributes are the ones from `scaleway_function` [resource](../resources/function.md)
