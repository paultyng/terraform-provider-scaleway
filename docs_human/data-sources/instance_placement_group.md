---
subcategory: "Instances"
page_title: "Scaleway: scaleway_instance_placement_group"
---

# scaleway_instance_placement_group

Gets information about a Security Group.

## Example Usage

```hcl
# Get info by placement group name
data "scaleway_instance_placement_group" "my_key" {
  name  = "my-placement-group-name"
}

# Get info by placement group id
data "scaleway_instance_placement_group" "my_key" {
  placement_group_id = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

- `name` - (Optional) The placement group name. Only one of `name` and `placement_group_id` should be specified.

- `placement_group_id` - (Optional) The placement group id. Only one of `name` and `placement_group_id` should be specified.

- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the placement group is associated with.

- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the placement group exists.

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `id` - The ID of the placement group.

- `policy_type` - The [policy type](https://developers.scaleway.com/en/products/instance/api/#placement-groups-d8f653) of the placement group.
- `policy_mode` -The [policy mode](https://developers.scaleway.com/en/products/instance/api/#placement-groups-d8f653) of the placement group.
- `tags` - A list of tags to apply to the placement group.
- `policy_respected` - Is true when the policy is respected.
- `organization_id` - The organization ID the placement group is associated with.