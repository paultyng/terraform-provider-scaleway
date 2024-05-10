---
subcategory: "VPC"
page_title: "Scaleway: scaleway_vpc_public_gateway_ip"
---

# Resource: scaleway_vpc_public_gateway_ip

Creates and manages Scaleway VPC Public Gateway IP.
For more information, see [the documentation](https://developers.scaleway.com/en/products/vpc-gw/api/v1/#ips-268151).

## Example Usage

```terraform
resource "scaleway_domain_record" "tf_A" {
    dns_zone = "example.com"
    name     = "tf"
    type     = "A"
    data     = "${scaleway_vpc_public_gateway_ip.main.address}"
    ttl      = 3600
    priority = 1
}

resource scaleway_vpc_public_gateway_ip main {
	reverse = "tf.example.com"
}
```

## Argument Reference

The following arguments are supported:

- `reverse` - (Optional) The reverse domain name for the IP address
- `tags` - (Optional) The tags associated with the public gateway IP.
- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the public gateway ip should be created.
- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the public gateway ip is associated with.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the public gateway IP.

~> **Important:** Public gateway IPs' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

- `address` - The IP address itself.
- `organization_id` - The organization ID the public gateway ip is associated with.
- `created_at` - The date and time of the creation of the public gateway ip.
- `updated_at` - The date and time of the last update of the public gateway ip.

## Import

Public gateway can be imported using the `{zone}/{id}`, e.g.

```bash
$ terraform import scaleway_vpc_public_gateway_ip.main fr-par-1/11111111-1111-1111-1111-111111111111
```
