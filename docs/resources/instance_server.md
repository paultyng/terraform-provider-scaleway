---
subcategory: "Instances"
page_title: "Scaleway: scaleway_instance_server"
---

# Resource: scaleway_instance_server

Creates and manages Scaleway compute Instances. For more information, see [the documentation](https://www.scaleway.com/en/developers/api/instance/#path-instances-list-all-instances).

Please check our [FAQ - Instances](https://www.scaleway.com/en/docs/faq/instances).

## Example Usage

### Basic

```terraform
resource "scaleway_instance_ip" "public_ip" {}

resource "scaleway_instance_server" "web" {
  type = "DEV1-S"
  image = "ubuntu_jammy"
  ip_id = scaleway_instance_ip.public_ip.id
}
```

### With additional volumes and tags

```terraform
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

```terraform
resource "scaleway_instance_ip" "ip" {}

resource "scaleway_instance_server" "web" {
  type = "DEV1-S"
  image = "f974feac-abae-4365-b988-8ec7d1cec10d"

  tags = [ "hello", "public" ]

  ip_id = scaleway_instance_ip.ip.id
}
```

### With security group

```terraform
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

```terraform
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

```terraform
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

```terraform
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

```terraform
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

#### Using Scaleway Block Storage (SBS) volume

```terraform
resource "scaleway_instance_server" "server" {
  type = "PLAY2-MICRO"
  image = "ubuntu_jammy"
  root_volume {
    volume_type = "sbs_volume"
    sbs_iops = 15000
    size_in_gb = 50
  }
}
```

## Argument Reference

The following arguments are supported:

- `type` - (Required) The commercial type of the server.
You find all the available types on the [pricing page](https://www.scaleway.com/en/pricing/).
Updates to this field will migrate the server, local storage constraint must be respected. [More info](https://www.scaleway.com/en/docs/compute/instances/api-cli/migrating-instances/).
Use `replace_on_type_change` to trigger replacement instead of migration.

~> **Important:** If `type` change and migration occurs, the server will be stopped and changed backed to its original state. It will be started again if it was running.

- `image` - (Optional) The UUID or the label of the base image used by the server. You can use [this endpoint](https://www.scaleway.com/en/developers/api/marketplace/#path-marketplace-images-list-marketplace-images)
to find either the right `label` or the right local image `ID` for a given `type`. Optional when creating an instance with an existing root volume.

You can check the available labels with our [CLI](https://www.scaleway.com/en/docs/compute/instances/api-cli/creating-managing-instances-with-cliv2/). ```scw marketplace image list```

To retrieve more information by label please use: ```scw marketplace image get label=<LABEL>```

- `name` - (Optional) The name of the server.

- `tags` - (Optional) The tags associated with the server.

- `security_group_id` - (Optional) The [security group](https://www.scaleway.com/en/developers/api/instance/#path-security-groups-update-a-security-group9) the server is attached to.

- `placement_group_id` - (Optional) The [placement group](https://www.scaleway.com/en/developers/api/instance/#path-security-groups-update-a-security-group the server is attached to.


~> **Important:** When updating `placement_group_id` the `state` must be set to `stopped`, otherwise it will fail.

- `root_volume` - (Optional) Root [volume](https://www.scaleway.com/en/developers/api/instance/#path-volume-types-list-volume-types) attached to the server on creation.
    - `volume_id` - (Optional) The volume ID of the root volume of the server, allows you to create server with an existing volume. If empty, will be computed to a created volume ID.
    - `size_in_gb` - (Required) Size of the root volume in gigabytes.
      To find the right size use [this endpoint](https://www.scaleway.com/en/developers/api/instance/#path-instances-list-all-instances) and
      check the `volumes_constraint.{min|max}_size` (in bytes) for your `commercial_type`.
      Updates to this field will recreate a new resource.
    - `volume_type` - (Optional) Volume type of root volume, can be `b_ssd`, `l_ssd` or `sbs_volume`, default value depends on server type
    - `delete_on_termination` - (Defaults to `true`) Forces deletion of the root volume on instance termination.
    - `sbs_iops` - (Optional) Choose IOPS of your sbs volume, has to be used with `sbs_volume` for root volume type.

~> **Important:** Updates to `root_volume.size_in_gb` will be ignored after the creation of the server.

- `additional_volume_ids` - (Optional) The [additional volumes](https://www.scaleway.com/en/developers/api/instance/#path-volume-types-list-volume-types)
attached to the server. Updates to this field will trigger a stop/start of the server.

~> **Important:** If this field contains local volumes, the `state` must be set to `stopped`, otherwise it will fail.

~> **Important:** If this field contains local volumes, you have to first detach them, in one apply, and then delete the volume in another apply.

- `enable_ipv6` - (Defaults to `false`) Determines if IPv6 is enabled for the server. Useful only with `routed_ip_enabled` as false, otherwise ipv6 is always supported.
  Deprecated: Please use a scaleway_instance_ip with a `routed_ipv6` type.

- `ip_id` - (Optional) The ID of the reserved IP that is attached to the server.

- `ip_ids` - (Optional) List of ID of reserved IPs that are attached to the server. Cannot be used with `ip_id`.

~> `ip_id` to `ip_ids` migration: if moving the ip from the old `ip_id` field to the new `ip_ids`, it should not detach the ip.

- `enable_dynamic_ip` - (Defaults to `false`) If true a dynamic IP will be attached to the server.

- `routed_ip_enabled` - (Defaults to `true`) If true, the server will support routed ips only. Changing it to true will migrate the server and its IP to routed type.

~> **Important:** Enabling routed ip will restart the server

- `state` - (Defaults to `started`) The state of the server. Possible values are: `started`, `stopped` or `standby`.

- `user_data` - (Optional) The user data associated with the server.
  Use the `cloud-init` key to use [cloud-init](https://cloudinit.readthedocs.io/en/latest/) on your instance.
  You can define values using:
    - string
    - UTF-8 encoded file content using [file](https://www.terraform.io/language/functions/file)
    - Binary files using [filebase64](https://www.terraform.io/language/functions/filebase64).

- `private_network` - (Optional) The private network associated with the server.
   Use the `pn_id` key to attach a [private_network](https://www.scaleway.com/en/developers/api/instance/#path-private-nics-list-all-private-nics) on your instance.

- `private_ip` - The list of private IP addresses associated with the resource.
    - `id` - The ID of the IP address resource.
    - `address` - The private IP address.

- `boot_type` - The boot Type of the server. Possible values are: `local`, `bootscript` or `rescue`.

- `replace_on_type_change` - (Defaults to false) If true, the server will be replaced if `type` is changed. Otherwise, the server will migrate.

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

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the server.

~> **Important:** Instance servers' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

- `placement_group_policy_respected` - True when the placement group policy is respected.
- `root_volume`
    - `volume_id` - The volume ID of the root volume of the server.
- `private_ip` - The Scaleway internal IP address of the server.
- `public_ip` - The public IP address of the server (Deprecated use `public_ips` instead).
- `public_ips` - The list of public IPs of the server.
    - `id` - The ID of the IP
    - `address` - The address of the IP
- `ipv6_address` - The default ipv6 address routed to the server. ( Only set when enable_ipv6 is set to true )
  Deprecated: Please use a scaleway_instance_ip with a `routed_ipv6` type.
- `ipv6_gateway` - The ipv6 gateway address. ( Only set when enable_ipv6 is set to true )
  Deprecated: Please use a scaleway_instance_ip with a `routed_ipv6` type.
- `ipv6_prefix_length` - The prefix length of the ipv6 subnet routed to the server. ( Only set when enable_ipv6 is set to true )
  Deprecated: Please use a scaleway_instance_ip with a `routed_ipv6` type.
- `boot_type` - The boot Type of the server. Possible values are: `local`, `bootscript` or `rescue`.
- `organization_id` - The organization ID the server is associated with.

## Import

Instance servers can be imported using the `{zone}/{id}`, e.g.

```bash
terraform import scaleway_instance_server.web fr-par-1/11111111-1111-1111-1111-111111111111
```
