---
page_title: "Scaleway: scaleway_instance_server"
description: |-
  Manages Scaleway Compute Instance servers.
---

# scaleway_instance_server

Creates and manages Scaleway Compute Instance servers. For more information, see [the documentation](https://developers.scaleway.com/en/products/instance/api/#servers-8bf7d7).

Please check our [FAQ - Instances](https://www.scaleway.com/en/docs/faq/instances).

## Examples

### Basic

```hcl
resource "scaleway_instance_ip" "public_ip" {}

resource "scaleway_instance_server" "web" {
  type = "DEV1-S"
  image = "ubuntu_jammy"
  ip_id = scaleway_instance_ip.public_ip.id
}
```

### With additional volumes and tags

```hcl
resource "scaleway_instance_volume" "data" {
  size_in_gb = 100
  type = "b_ssd"
}

resource "scaleway_instance_server" "web" {
  type = "DEV1-S"
  image = "ubuntu_jammy"

  tags = [ "hello", "public" ]

  root_volume {
    delete_on_termination = false
  }

  additional_volume_ids = [ scaleway_instance_volume.data.id ]
}
```

### With a reserved IP

```hcl
resource "scaleway_instance_ip" "ip" {}

resource "scaleway_instance_server" "web" {
  type = "DEV1-S"
  image = "f974feac-abae-4365-b988-8ec7d1cec10d"

  tags = [ "hello", "public" ]

  ip_id = scaleway_instance_ip.ip.id
}
```

### With security group

```hcl
resource "scaleway_instance_security_group" "www" {
  inbound_default_policy = "drop"
  outbound_default_policy = "accept"

  inbound_rule {
    action = "accept"
    port = "22"
    ip = "212.47.225.64"
  }

  inbound_rule {
    action = "accept"
    port = "80"
  }

  inbound_rule {
    action = "accept"
    port = "443"
  }

  outbound_rule {
    action = "drop"
    ip_range = "10.20.0.0/24"
  }
}

resource "scaleway_instance_server" "web" {
  type = "DEV1-S"
  image = "ubuntu_jammy"

  security_group_id= scaleway_instance_security_group.www.id
}
```

### With user data and cloud-init

```hcl
resource "scaleway_instance_server" "web" {
  type  = "DEV1-S"
  image = "ubuntu_jammy"

  user_data = {
    foo        = "bar"
    cloud-init = file("${path.module}/cloud-init.yml")
  }
}
```

### With private network

```hcl
resource scaleway_vpc_private_network pn01 {
    name = "private_network_instance"
}

resource "scaleway_instance_server" "base" {
  image = "ubuntu_jammy"
  type  = "DEV1-S"

  private_network {
    pn_id = scaleway_vpc_private_network.pn01.id
  }
}
```

### Root volume configuration

#### Resized block volume with installed image

```hcl
resource "scaleway_instance_server" "image" {
  type = "PRO2-XXS"
  image = "ubuntu_jammy"
  root_volume {
    volume_type = "b_ssd"
    size_in_gb = 100
  }
}
```

#### From snapshot

```hcl
data "scaleway_instance_snapshot" "snapshot" {
  name = "my_snapshot"
}

resource "scaleway_instance_volume" "from_snapshot" {
  from_snapshot_id = data.scaleway_instance_snapshot.snapshot.id
  type = "b_ssd"
}

resource "scaleway_instance_server" "from_snapshot" {
  type = "PRO2-XXS"
  root_volume {
    volume_id = scaleway_instance_volume.from_snapshot.id
  }
}
```

## Arguments Reference

The following arguments are supported:

- `type` - (Required) The commercial type of the server.
You find all the available types on the [pricing page](https://www.scaleway.com/en/pricing/).
Updates to this field will recreate a new resource.

- `image` - (Optional) The UUID or the label of the base image used by the server. You can use [this endpoint](https://api-marketplace.scaleway.com/images?page=1&per_page=100)
to find either the right `label` or the right local image `ID` for a given `type`. Optional when creating an instance with an existing root volume.

You can check the available labels with our [CLI](https://www.scaleway.com/en/docs/compute/instances/api-cli/creating-managing-instances-with-cliv2/). ```scw marketplace image list```

To retrieve more information by label please use: ```scw marketplace image get label=<LABEL>```

- `name` - (Optional) The name of the server.

- `tags` - (Optional) The tags associated with the server.

- `security_group_id` - (Optional) The [security group](https://developers.scaleway.com/en/products/instance/api/#security-groups-8d7f89) the server is attached to.

- `placement_group_id` - (Optional) The [placement group](https://developers.scaleway.com/en/products/instance/api/#placement-groups-d8f653) the server is attached to.


~> **Important:** When updating `placement_group_id` the `state` must be set to `stopped`, otherwise it will fail.

- `root_volume` - (Optional) Root [volume](https://developers.scaleway.com/en/products/instance/api/#volumes-7e8a39) attached to the server on creation.
    - `volume_id` - (Optional) The volume ID of the root volume of the server, allows you to create server with an existing volume. If empty, will be computed to a created volume ID.
    - `size_in_gb` - (Required) Size of the root volume in gigabytes.
      To find the right size use [this endpoint](https://api.scaleway.com/instance/v1/zones/fr-par-1/products/servers) and
      check the `volumes_constraint.{min|max}_size` (in bytes) for your `commercial_type`.
      Updates to this field will recreate a new resource.
    - `volume_type` - (Optional) Volume type of root volume, can be `b_ssd` or `l_ssd`, default value depends on server type
    - `delete_on_termination` - (Defaults to `true`) Forces deletion of the root volume on instance termination.

~> **Important:** Updates to `root_volume.size_in_gb` will be ignored after the creation of the server.

- `additional_volume_ids` - (Optional) The [additional volumes](https://developers.scaleway.com/en/products/instance/api/#volumes-7e8a39)
attached to the server. Updates to this field will trigger a stop/start of the server.

~> **Important:** If this field contains local volumes, the `state` must be set to `stopped`, otherwise it will fail.

~> **Important:** If this field contains local volumes, you have to first detach them, in one apply, and then delete the volume in another apply.

- `enable_ipv6` - (Defaults to `false`) Determines if IPv6 is enabled for the server.

- `ip_id` = (Optional) The ID of the reserved IP that is attached to the server.

- `enable_dynamic_ip` - (Defaults to `false`) If true a dynamic IP will be attached to the server.

- `state` - (Defaults to `started`) The state of the server. Possible values are: `started`, `stopped` or `standby`.

- `user_data` - (Optional) The user data associated with the server.
  Use the `cloud-init` key to use [cloud-init](https://cloudinit.readthedocs.io/en/latest/) on your instance.
  You can define values using:
    - string
    - UTF-8 encoded file content using [file](https://www.terraform.io/language/functions/file)
    - Binary files using [filebase64](https://www.terraform.io/language/functions/filebase64).

- `private_network` - (Optional) The private network associated with the server.
   Use the `pn_id` key to attach a [private_network](https://developers.scaleway.com/en/products/instance/api/#private-nics-a42eea) on your instance.

- `boot_type` - The boot Type of the server. Possible values are: `local`, `bootscript` or `rescue`.

- `bootscript_id` - The ID of the bootscript to use  (set boot_type to `bootscript`).

- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the server should be created.

- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the server is associated with.


## Private Network

~> **Important:** Updates to `private_network` will recreate a new private network interface.

- `pn_id` - (Required) The private network ID where to connect.
- `mac_address` The private NIC MAC address.
- `status` The private NIC state.
- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the server must be created.

~> **Important:**

- You can only attach an instance in the same [zone](../guides/regions_and_zones.md#zones) as a private network.
- Instance supports maximum 8 different private networks.

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `id` - The ID of the server.
- `placement_group_policy_respected` - True when the placement group policy is respected.
- `root_volume`
    - `volume_id` - The volume ID of the root volume of the server.
- `private_ip` - The Scaleway internal IP address of the server.
- `public_ip` - The public IPv4 address of the server.
- `ipv6_address` - The default ipv6 address routed to the server. ( Only set when enable_ipv6 is set to true )
- `ipv6_gateway` - The ipv6 gateway address. ( Only set when enable_ipv6 is set to true )
- `ipv6_prefix_length` - The prefix length of the ipv6 subnet routed to the server. ( Only set when enable_ipv6 is set to true )
- `boot_type` - The boot Type of the server. Possible values are: `local`, `bootscript` or `rescue`.
- `organization_id` - The organization ID the server is associated with.

## Import

Instance servers can be imported using the `{zone}/{id}`, e.g.

```bash
$ terraform import scaleway_instance_server.web fr-par-1/11111111-1111-1111-1111-111111111111
```
