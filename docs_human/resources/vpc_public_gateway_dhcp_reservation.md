---
subcategory: "VPC"
page_title: "Scaleway: scaleway_vpc_public_gateway_dhcp_reservation"
---

# Resource: scaleway_vpc_public_gateway_dhcp_reservation

Creates and manages the [Scaleway DHCP Reservations](https://www.scaleway.com/en/docs/network/vpc/concepts/#dhcp).

The static associations are used to assign IP addresses based on the MAC addresses of the Instance.

Statically assigned IP addresses should fall within the configured subnet, but be outside of the dynamic range.

For more information, see [the documentation](https://developers.scaleway.com/en/products/vpc-gw/api/v1/#dhcp-c05544) and [configuration guide](https://www.scaleway.com/en/docs/network/vpc/how-to/configure-a-public-gateway/#how-to-review-and-configure-dhcp).

[DHCP reservations](https://developers.scaleway.com/en/products/vpc-gw/api/v1/#dhcp-entries-e40fb6) hold both dynamic DHCP leases (IP addresses dynamically assigned by the gateway to instances) and static user-created DHCP reservations.

## Example Usage

```terraform
resource scaleway_vpc_private_network main {
    name = "your_private_network"
}

resource "scaleway_instance_server" "main" {
    image = "ubuntu_jammy"
    type  = "DEV1-S"
    zone = "fr-par-1"

    private_network {
        pn_id = scaleway_vpc_private_network.main.id
    }
}

resource scaleway_vpc_public_gateway_ip main {
}

resource scaleway_vpc_public_gateway_dhcp main {
    subnet = "192.168.1.0/24"
}

resource scaleway_vpc_public_gateway main {
    name = "foobar"
    type = "VPC-GW-S"
    ip_id = scaleway_vpc_public_gateway_ip.main.id
}

resource scaleway_vpc_gateway_network main {
    gateway_id = scaleway_vpc_public_gateway.main.id
    private_network_id = scaleway_vpc_private_network.main.id
    dhcp_id = scaleway_vpc_public_gateway_dhcp.main.id
    cleanup_dhcp = true
    enable_masquerade = true
    depends_on = [scaleway_vpc_public_gateway_ip.main, scaleway_vpc_private_network.main]
}

resource scaleway_vpc_public_gateway_dhcp_reservation main {
    gateway_network_id = scaleway_vpc_gateway_network.main.id
    mac_address = scaleway_instance_server.main.private_network.0.mac_address
    ip_address = "192.168.1.1"
}
```

## Argument Reference

The following arguments are supported:

- `gateway_network_id` - (Required) The ID of the owning GatewayNetwork.
- `ip_address` - (Required) The IP address to give to the machine (IP address).
- `mac_address` - (Required) The MAC address to give a static entry to.
- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the public gateway DHCP config should be created.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the public gateway DHCP Reservation config.

~> **Important:** Public gateway DHCP reservations configurations' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

- `hostname` - The Hostname of the client machine.
- `type` - The reservation type, either static (DHCP reservation) or dynamic (DHCP lease). Possible values are reservation and lease.
- `created_at` - The date and time of the creation of the public gateway DHCP config.
- `updated_at` - The date and time of the last update of the public gateway DHCP config.

## Import

Public gateway DHCP Reservation config can be imported using the `{zone}/{id}`, e.g.

```bash
$ terraform import scaleway_vpc_public_gateway_dhcp_reservation.main fr-par-1/11111111-1111-1111-1111-111111111111
```
