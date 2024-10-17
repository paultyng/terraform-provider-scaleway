---
subcategory: "Load Balancers"
page_title: "Scaleway: scaleway_lb"
---

# Resource: scaleway_lb

Creates and manages Scaleway Load Balancers.

For more information, see the [main documentation](https://www.scaleway.com/en/docs/network/load-balancer/concepts/#load-balancers) or [API documentation](https://www.scaleway.com/en/developers/api/load-balancer/zoned-api/#path-load-balancer-list-load-balancers).

## Example Usage

### Basic

```terraform
resource "scaleway_lb_ip" "main" {
  zone = "fr-par-1"
}

resource "scaleway_lb" "base" {
  ip_ids  = [scaleway_lb_ip.main.id]
  zone    = scaleway_lb_ip.main.zone
  type    = "LB-S"
}
```

### Private LB

```terraform

resource "scaleway_lb" "base" {
  name               = "private-lb"
  type               = "LB-S"
  assign_flexible_ip = false
}
```

### With IPv6

```terraform
resource "scaleway_lb_ip" "v4" {
}
resource "scaleway_lb_ip" "v6" {
  is_ipv6 = true
}
resource scaleway_lb main {
  ip_ids = [scaleway_lb_ip.v4.id, scaleway_lb_ip.v6.id]
  name   = "ipv6-lb"
  type   = "LB-S"
}
```

### With IPAM IDs

```terraform
resource "scaleway_vpc" "vpc01" {
  name = "my vpc"
}

resource "scaleway_vpc_private_network" "pn01" {
  vpc_id = scaleway_vpc.vpc01.id
  ipv4_subnet {
    subnet = "172.16.32.0/22"
  }
}

resource "scaleway_ipam_ip" "ip01" {
  address = "172.16.32.7"
  source {
    private_network_id = scaleway_vpc_private_network.pn01.id
  }
}

resource "scaleway_lb_ip" "v4" {
}

resource scaleway_lb main {
  ip_ids = [scaleway_lb_ip.v4.id]
  name   = "my-lb"
  type   = "LB-S"

  private_network {
    private_network_id = scaleway_vpc_private_network.pn01.id
    ipam_ids = [scaleway_ipam_ip.ip01.id]
  }
}
```

### Multiple configurations

```terraform
### IP for Public Gateway
resource "scaleway_vpc_public_gateway_ip" "main" {
}

### Scaleway Private Network
resource scaleway_vpc_private_network main {
}

### VPC Public Gateway Network
resource "scaleway_vpc_public_gateway" "main" {
    name  = "tf-test-public-gw"
    type  = "VPC-GW-S"
    ip_id = scaleway_vpc_public_gateway_ip.main.id
}

### VPC Public Gateway Network DHCP config
resource "scaleway_vpc_public_gateway_dhcp" "main" {
    subnet = "10.0.0.0/24"
}

### VPC Gateway Network
resource "scaleway_vpc_gateway_network" "main" {
    gateway_id         = scaleway_vpc_public_gateway.main.id
    private_network_id = scaleway_vpc_private_network.main.id
    dhcp_id            = scaleway_vpc_public_gateway_dhcp.main.id
    cleanup_dhcp       = true
    enable_masquerade  = true
}

### Scaleway Instance
resource "scaleway_instance_server" "main" {
    name        = "Scaleway Terraform Provider"
    type        = "DEV1-S"
    image       = "debian_bullseye"
    enable_ipv6 = false

    private_network {
        pn_id = scaleway_vpc_private_network.main.id
    }
}

### IP for LB IP
resource scaleway_lb_ip main {
}

### Scaleway Private Network
resource scaleway_vpc_private_network "main" {
    name = "private network with DHCP config"
}

### Scaleway Load Balancer
resource scaleway_lb main {
    ip_id = scaleway_lb_ip.main.id
    name = "MyTest"
    type = "LB-S"

    private_network {
        private_network_id = scaleway_vpc_private_network.main.id
        dhcp_config = true
    }

    depends_on = [scaleway_vpc_public_gateway.main]
}
```

## Argument Reference

The following arguments are supported:

- `ip_ids` - (Optional) The List of IP IDs to attach to the Load Balancer.

- `ip_id` - (Deprecated) The ID of the associated Load Balancer IP. See below.

~> **Important:** Updates to `ip_id` will recreate the Load Balancer.

- `type` - (Required) The type of the Load Balancer. Please check the [migration section](#migration) to upgrade the type.

- `assign_flexible_ip` - (Optional) Defines whether to automatically assign a flexible public IPv4 to the Load Balancer.

- `assign_flexible_ipv6` - (Optional) Defines whether to automatically assign a flexible public IPv6 to the Load Balancer.

- `name` - (Optional) The name of the Load Balancer.

- `description` - (Optional) The description of the Load Balancer.

- `tags` - (Optional) The tags associated with the Load Balancer.

- `release_ip` - (Defaults to false) The `release_ip` allow the release of the IP address associated with the Load Balancer.

- `ssl_compatibility_level` - (Optional) Enforces minimal SSL version (in SSL/TLS offloading context). Please check [possible values](https://www.scaleway.com/en/developers/api/load-balancer/zoned-api/#path-load-balancer-create-a-load-balancer).

- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) of the Load Balancer.

- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the Project the Load Balancer is associated with.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the Load Balancer.

~> **Important:** Load Balancers IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

- `ip_address` -  The Load Balancer public IPv4 address.
- `ipv6_address` -  The Load Balancer public IPv6 address.
- `organization_id` - The ID of the Organization ID the Load Balancer is associated with.

~> **Important:** `release_ip` will not be supported. This prevents the destruction of the IP from releasing a Load Balancer.
The `resource_lb_ip` will be the only resource that handles those IPs.

## Migration

In order to migrate to other Load Balancer types, you can check upwards or downwards migration via our CLI `scw lb lb-types list`.
This change will not recreate your Load Balancer.

Please check our [documentation](https://www.scaleway.com/en/developers/api/load-balancer/zoned-api/#path-load-balancer-migrate-a-load-balancer) for further details

## IP ID

Since v1.15.0, `ip_id` is a required field. This means that now a separate `scaleway_lb_ip` is required.
When importing, the IP needs to be imported as well as the Load Balancer.
When upgrading to v1.15.0, you will need to create a new `scaleway_lb_ip` resource and import it.

For instance, if you had the following:

```terraform
resource "scaleway_lb" "main" {
  zone = "fr-par-1"
  type   = "LB-S"
}
```

You will need to update it to:

```terraform
resource "scaleway_lb_ip" "main" {
}

resource "scaleway_lb" "main" {
  ip_id  = scaleway_lb_ip.main.id
  zone = "fr-par-1"
  type   = "LB-S"
  release_ip = false
}
```

## Attributes Reference

- `private_network_id` - (Required) The ID of the Private Network to attach to.

- ~> **Important:** Updates to `private_network` will recreate the attachment.

- `ipam_ids` - (Optional) IPAM ID of a pre-reserved IP address to assign to the Load Balancer on this Private Network.

- `dhcp_config` - (Deprecated) Please use `ipam_ids`. Set to `true` if you want to let DHCP assign IP addresses. See below.

- `static_config` - (Deprecated) Please use `ipam_ids`. Define a local ip address of your choice for the load balancer instance.

- `private_ip` - The list of private IP addresses associated with the resource.
    - `id` - The ID of the IP address resource.
    - `address` - The private IP address.

- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the Private Network was created.


## Import

Load Balancers can be imported using `{zone}/{id}`, e.g.

```bash
terraform import scaleway_lb.main fr-par-1/11111111-1111-1111-1111-111111111111
```

Be aware that you will also need to import the `scaleway_lb_ip` resource.
