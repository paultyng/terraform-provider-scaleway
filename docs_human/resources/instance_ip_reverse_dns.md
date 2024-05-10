---
subcategory: "Instances"
page_title: "Scaleway: scaleway_instance_ip_reverse_dns"
---

# Resource: scaleway_instance_ip_reverse_dns

Manages Scaleway Compute Instance IPs Reverse DNS.

Please check our [guide](https://www.scaleway.com/en/docs/compute/instances/how-to/configure-reverse-dns/) for more details

## Example Usage

```terraform
resource "scaleway_instance_ip" "server_ip" {}

resource "scaleway_domain_record" "tf_A" {
  dns_zone = "scaleway.com"
  name     = "www"
  type     = "A"
  data     = "${scaleway_instance_ip.server_ip.address}"
  ttl      = 3600
  priority = 1
}

resource "scaleway_instance_ip_reverse_dns" "reverse" {
  ip_id   = scaleway_instance_ip.server_ip.id
  reverse = "www.scaleway.com"
}
```

## Argument Reference

The following arguments are supported:

- `ip_id` - (Required) The IP ID
- `reverse` - (Required) The reverse DNS for this IP.
- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the IP should be reserved.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the IP for which the reverse DNS is configured.

~> **Important:** Instance IPs' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

## Import

IPs reverse DNS can be imported using the `{zone}/{id}`, e.g.

```bash
$ terraform import scaleway_instance_ip_reverse_dns.reverse fr-par-1/11111111-1111-1111-1111-111111111111
```
