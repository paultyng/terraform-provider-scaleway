---
page_title: "Scaleway: scaleway_lb_ip"
description: |-
  Gets information about a Load Balancer IP.
---

# scaleway_lb_ip

Gets information about a Load Balancer IP.

## Example Usage

```hcl
# Get info by IP address
data "scaleway_lb_ip" "my_ip" {
  ip_address = "0.0.0.0"
}

# Get info by IP ID
data "scaleway_lb_ip" "my_ip" {
  ip_id = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

The following arguments are supported:

- `ip_address` - (Optional) The IP address.
  Only one of `ip_address` and `ip_id` should be specified.

- `ip_id` - (Optional) The IP ID.
  Only one of `ip_address` and `ip_id` should be specified.

- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the IP should be reserved.

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `reverse` - The reverse domain associated with this IP.

- `lb_id` - The associated load-balancer ID if any

- `organization_id` - (Defaults to [provider](../index.md#organization_id) `organization_id`) The ID of the organization the LB IP is associated with.
