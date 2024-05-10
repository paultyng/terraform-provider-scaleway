---
subcategory: "Load Balancers"
page_title: "Scaleway: scaleway_lb"
---

# Resource: scaleway_lb

Creates and manages Scaleway Load-Balancers.
For more information, see [the documentation](https://www.scaleway.com/en/developers/api/load-balancer/zoned-api).

## Example Usage

### Basic

```terraform
resource "scaleway_lb_ip" "main" {
  zone = "fr-par-1"
}

resource "scaleway_lb" "base" {
  ip_id  = scaleway_lb_ip.main.id
  zone   = scaleway_lb_ip.main.zone
  type   = "LB-S"
}
```

### Private LB

```terraform

resource "scaleway_lb" "base" {
  ip_id  = scaleway_lb_ip.main.id
  zone   = scaleway_lb_ip.main.zone
  type   = "LB-S"
  assign_flexible_ip = false
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
    name = "private network with static config"
}

### Scaleway Load Balancer
resource scaleway_lb main {
    ip_id = scaleway_lb_ip.main.id
    name = "MyTest"
    type = "LB-S"

    private_network {
        private_network_id = scaleway_vpc_private_network.main.id
        static_config = ["172.16.0.100"]
    }

    private_network {
        private_network_id = scaleway_vpc_private_network.main.id
        dhcp_config = true
    }

    depends_on = [scaleway_vpc_public_gateway.main]
}
```

## Argument Reference

The following arguments are supported:

- `ip_id` - (Optional) The ID of the associated LB IP. See below.

~> **Important:** Updates to `ip_id` will recreate the load-balancer.

- `type` - (Required) The type of the load-balancer. Please check the [migration section](#migration) to upgrade the type.

- `assign_flexible_ip` - (Optional) Defines whether to automatically assign a flexible public IP to the load-balancer.

- `name` - (Optional) The name of the load-balancer.

- `description` - (Optional) The description of the load-balancer.

- `tags` - (Optional) The tags associated with the load-balancers.

- `release_ip` - (Defaults to false) The release_ip allow release the ip address associated with the load-balancers.

- `ssl_compatibility_level` - (Optional) Enforces minimal SSL version (in SSL/TLS offloading context). Please check [possible values](https://www.scaleway.com/en/developers/api/load-balancer/zoned-api/#path-load-balancer-create-a-load-balancer).

- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) of the load-balancer.

- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the load-balancer is associated with.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the load-balancer.

~> **Important:** Load-Balancers' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

- `ip_address` -  The load-balance public IP Address
- `organization_id` - The organization ID the load-balancer is associated with.

~> **Important:** `release_ip` will not be supported. This prevents the destruction of the IP from releasing a LBs.
The `resource_lb_ip` will be the only resource that handles those IPs.

## Migration

In order to migrate to other types you can check the migration up or down via our CLI `scw lb lb-types list`.
this change will not recreate your Load Balancer.

Please check our [documentation](https://www.scaleway.com/en/developers/api/load-balancer/zoned-api/#path-load-balancer-migrate-a-load-balancer) for further details

## IP ID

Since v1.15.0, `ip_id` is a required field. This means that now a separate `scaleway_lb_ip` is required.
When importing, the IP needs to be imported as well as the LB.
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

## Private Network with static config

```terraform
resource scaleway_lb_ip main {
}

resource scaleway_vpc_private_network main {
    name = "MyTest"
}

resource scaleway_lb main {
    ip_id = scaleway_lb_ip.main.id
    name = "MyTest"
    type = "LB-S"
    release_ip = false
    private_network {
        private_network_id = scaleway_vpc_private_network.main.id
        static_config = ["172.16.0.100"]
    }
}
```

## Attributes Reference

- `private_network_id` - (Required) The ID of the Private Network to associate.

- ~> **Important:** Updates to `private_network` will recreate the attachment.

- `static_config` - (Optional) Define a local ip address of your choice for the load balancer instance. See below.

- `dhcp_config` - (Optional) Set to true if you want to let DHCP assign IP addresses. See below.

~> **Important:**  Only one of static_config and dhcp_config may be set.

- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the private network was created.


## Import

Load-Balancer can be imported using the `{zone}/{id}`, e.g.

```bash
$ terraform import scaleway_lb.main fr-par-1/11111111-1111-1111-1111-111111111111
```

Be aware that you will also need to import the `scaleway_lb_ip` resource.
