---
layout: "scaleway"
page_title: "Scaleway: scaleway_lb_ip_beta"
description: |-
  Gets information about a Load Balancer IP.
---

# scaleway_lb_ip_beta

Gets information about a Load Balancer IP.

## Example Usage

```hcl
# Get info by IP address
data "scaleway_lb_ip_beta" "my_ip" {
  ip_address = "0.0.0.0"
}

# Get info by IP ID
data "scaleway_lb_ip_beta" "my_ip" {
  ip_id = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

- `ip_address` - (Optional) The IP address.
  Only one of `ip_address` and `lb_id` should be specified.

- `lb_id` - (Optional) The IP ID.
  Only one of `ip_address` and `ip_id` should be specified.

- `region` - (Defaults to [provider](../index.html#region) `region`) The [region](../guides/regions_and_zones.html#zones) in which the LB IP exists.

- `organization_id` - (Defaults to [provider](../index.html#organization_id) `organization_id`) The ID of the organization the LB IP is associated with.

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `reverse` - The reverse domain associated with this IP.

- `lb_id` - The associated load-balance ID if any
