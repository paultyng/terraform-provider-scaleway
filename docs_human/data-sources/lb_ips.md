---
subcategory: "Load Balancers"
page_title: "Scaleway: scaleway_lb_ips"
---

# scaleway_lb_ips

Gets information about multiple Load Balancer IPs.

## Example Usage

```hcl
# Find multiple IPs that share the same CIDR block
data "scaleway_lb_ips" "my_key" {
  ip_cidr_range = "0.0.0.0/0"
}
# Find IPs by CIDR block and zone
data "scaleway_lb_ips" "my_key" {
  ip_cidr_range = "0.0.0.0/0"
  zone       = "fr-par-2"
}
```

## Argument Reference

- `ip_cidr_range` - (Optional) The IP CIDR range used as a filter. IPs within a CIDR block like it are listed.

- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which IPs exist.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `ips` - List of found IPs
    - `id` - The associated IP ID.
    - `lb_id` - The associated load-balancer ID if any
    - `ip_address` - The IP Address
    - `zone` - The [zone](../guides/regions_and_zones.md#zones) in which the load-balancer is.
    - `reverse` - The reverse domain associated with this IP.
    - `organization_id` - The organization ID the load-balancer is associated with.
    - `project_id` - The ID of the project the load-balancer is associated with.