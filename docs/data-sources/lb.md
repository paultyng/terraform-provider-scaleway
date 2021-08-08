---
page_title: "Scaleway: scaleway_lb"
description: |-
  Gets information about a Load Balancer.
---

# scaleway_lb

Gets information about a Load Balancer.

## Example Usage

```hcl
# Get info by name
data "scaleway_lb" "by_name" {
  name = "foobar"
}

# Get info by ID
data "scaleway_lb" "by_id" {
  lb_id = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

- `name` - (Optional) The IP address.
  Only one of `name` and `lb_id` should be specified.

- `lb_id` - (Optional) The ID.
  Only one of `ip_address` and `lb_id` should be specified.

- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#zones) in which the LB exists.

- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the LB is associated with.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the load-balancer.
- `ip_address` -  The load-balancer public IP Address
- `organization_id` - The organization ID the load-balancer is associated with.
- `project_id` - The project ID the load-balancer is associated with.
- `tags` - (Optional) The tags associated with the load-balancers.
- `type` - (Required) The type of the load-balancer.
