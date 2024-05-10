---
subcategory: "IPAM"
page_title: "Scaleway: scaleway_ipam_ip_reverse_dns"
---

# Resource: scaleway_ipam_ip_reverse_dns

Manages Scaleway IPAM IP Reverse DNS.

## Example Usage

```terraform
resource "scaleway_instance_ip" "ip01" {
  type = "routed_ipv6"
}

resource "scaleway_instance_server" "srv01" {
  name   = "tf-tests-instance-server-ips"
  ip_ids = [scaleway_instance_ip.ip01.id]
  image  = "ubuntu_jammy"
  type   = "PRO2-XXS"
  state  = "stopped"
}

data "scaleway_ipam_ip" "ipam01" {
  resource {
    id   = scaleway_instance_server.srv01.id
    type = "instance_server"
  }
  type = "ipv6"
}

resource "scaleway_domain_record" "tf_AAAA" {
  dns_zone = "example.com"
  name     = ""
  type     = "AAAA"
  data     = cidrhost(data.scaleway_ipam_ip.ipam01.address_cidr, 42)
  ttl      = 3600
  priority = 1
}

resource "scaleway_ipam_ip_reverse_dns" "base" {
  ipam_ip_id = data.scaleway_ipam_ip.ipam01.id

  hostname   = "example.com"
  address    = cidrhost(data.scaleway_ipam_ip.ipam01.address_cidr, 42)
}
```

## Argument Reference

The following arguments are supported:

- `ipam_ip_id` - (Required) The IPAM IP ID.
- `hostname` - (Required) The reverse domain name.
- `address` - (Required) The IP corresponding to the hostname.
- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#regions) of the IP reverse DNS.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the IPAM IP for which the DNS reverse is configured.

~> **Important:** IPAM IPs' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111

## Import

IPAM IP reverse DNS can be imported using the `{region}/{id}`, e.g.

```bash
$ terraform import scaleway_ipam_ip_reverse_dns.main fr-par/11111111-1111-1111-1111-111111111111
```
