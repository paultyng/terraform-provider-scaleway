---
subcategory: "Load Balancers"
page_title: "Scaleway: scaleway_lbs"
---

# scaleway_lbs

Gets information about multiple Load Balancers.

## Example Usage

```hcl
# Find LBs by name
data "scaleway_lbs" "my_key" {
  name = "foobar"
}

# Find LBs by name and zone
data "scaleway_lbs" "my_key" {
  name = "foobar"
  zone = "fr-par-2"
}
```

## Argument Reference

- `name` - (Optional) The load balancer name used as a filter. LBs with a name like it are listed.

- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which LBs exist.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `lbs` - List of found LBs
    - `id` - The ID of the load-balancer.
    - `tags` - The tags associated with the load-balancer.
    - `description` - The description of the load-balancer.
    - `status` - The state of the LB's instance. Possible values are: `unknown`, `ready`, `pending`, `stopped`, `error`, `locked` and `migrating`.
    - `zone` - The [zone](../guides/regions_and_zones.md#zones) in which the load-balancer is.
    - `name` - The name of the load-balancer.
    - `type` - The offer type of the load-balancer.
    - `instances` - List of underlying instances.
    - `ips` - List of IPs attached to the Load balancer.
    - `frontend_count` - Number of frontends the Load balancer has.
    - `backend_count` - Number of backends the Load balancer has.
    - `private_network_count` - Number of private networks attached to the Load balancer.
    - `route_count` - Number of routes the Load balancer has.
    - `subscriber` - The subscriber information.
    - `ssl_compatibility_level` - Determines the minimal SSL version which needs to be supported on client side.
    - `created_at` - Date at which the Load balancer was created.
    - `updated_at` - Date at which the Load balancer was updated.
    - `organization_id` - The organization ID the load-balancer is associated with.
    - `project_id` - The ID of the project the load-balancer is associated with.
