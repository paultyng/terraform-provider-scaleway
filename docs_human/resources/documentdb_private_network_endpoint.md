---
subcategory: "Databases"
page_title: "Scaleway: scaleway_documentdb_private_network_endpoint"
---

# Resource: scaleway_documentdb_private_network_endpoint

Creates and manages Scaleway Database Private Network Endpoint.
For more information, see [the documentation](https://www.scaleway.com/en/developers/api/document_db/).

## Example Usage

### Example Basic

```terraform
resource "scaleway_vpc_private_network" "pn" {
  name = "my_private_network"
}

resource "scaleway_documentdb_private_network_endpoint" "main" {
  instance_id    = "11111111-1111-1111-1111-111111111111"
  private_network {
    ip_net = "172.16.32.3/22"
    id     = scaleway_vpc_private_network.pn.id
  }
  depends_on = [scaleway_vpc_private_network.pn]
}
```

## Argument Reference

The following arguments are supported:

- `instance_id` - (Required) UUID of the documentdb instance.

- `ip_net` - (Optional) The IP network address within the private subnet. This must be an IPv4 address with a
  CIDR notation. The IP network address within the private subnet is determined by the IP Address Management (IPAM)
  service if not set.

- `private_network_id` - (Required) The ID of the private network.

## Private Network

~> **Important:** Updates to `private_network_id` will recreate the attachment Instance.

~> **NOTE:** Please calculate your host IP.
using [cirhost](https://developer.hashicorp.com/terraform/language/functions/cidrhost). Otherwise, lets IPAM service
handle the host IP on the network.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `ip` - IPv4 address on the network.
- `port` - Port in the Private Network.
- `name` - Name of the endpoint.
- `hostname` - Hostname of the endpoint.


~> **Important:** Database instances' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they
are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`

## Import

Database Instance Endpoint can be imported using the `{region}/{endpoint_id}`, e.g.

```bash
$ terraform import scaleway_documentdb_private_network_endpoint.end fr-par/11111111-1111-1111-1111-111111111111
```
