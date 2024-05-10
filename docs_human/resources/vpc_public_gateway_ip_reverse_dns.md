---
subcategory: "VPC"
page_title: "Scaleway: scaleway_vpc_public_gateway_ip_reverse_dns"
---

# Resource: scaleway_vpc_public_gateway_ip_reverse_dns

Manages Scaleway VPC Public Gateways IPs reverse DNS.
For more information, see [the documentation](https://developers.scaleway.com/en/products/vpc-gw/api/v1/#ips-268151).

## Example Usage

```terraform
resource "scaleway_vpc_public_gateway_ip" "main" {}

resource "scaleway_domain_record" "tf_A" {
    dns_zone = "example.com"
    name     = "tf"
    type     = "A"
    data     = "${scaleway_vpc_public_gateway_ip.main.address}"
    ttl      = 3600
    priority = 1
}

resource "scaleway_vpc_public_gateway_ip_reverse_dns" "main" {
    gateway_ip_id   = scaleway_vpc_public_gateway_ip.main.id
    reverse         = "tf.example.com"
}
```

## Argument Reference

The following arguments are supported:

- `gateway_ip_id` - (Required) The public gateway IP ID
- `reverse` - (Required) The reverse domain name for this IP address
- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the IP should be reserved.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the public gateway IP for which the reverse DNS is configured.

~> **Important:** Public gateway IPs' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`


## Import

Public gateway IPs reverse DNS can be imported using the `{zone}/{id}`, e.g.

```bash
$ terraform import scaleway_vpc_public_gateway_ip_reverse_dns.reverse fr-par-1/11111111-1111-1111-1111-111111111111
```
