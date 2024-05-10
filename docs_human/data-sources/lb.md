---
subcategory: "Load Balancers"
page_title: "Scaleway: scaleway_lb"
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

- `name` - (Optional) The load balancer name.

- `ip_id` - (Optional) The load balancer IP ID.

- `project_id` - (Optional) The ID of the project the LB is associated with.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the load-balancer.

~> **Important:** Load-Balancers' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

- `ip_address` - The load-balancer public IP Address.

- `type` - The type of the load-balancer.

- `tags` - The tags associated with the load-balancer.

- `zone` -  (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the LB exists.