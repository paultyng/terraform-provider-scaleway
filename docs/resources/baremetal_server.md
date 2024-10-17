---
subcategory: "Elastic Metal"
page_title: "Scaleway: scaleway_baremetal_server"
---

# Resource: scaleway_baremetal_server

Creates and manages Scaleway Compute Baremetal servers. For more information, see [the documentation](https://www.scaleway.com/en/developers/api/elastic-metal/).

## Example Usage

### Basic

```terraform
data "scaleway_account_ssh_key" "main" {
  name = "main"
}

resource "scaleway_baremetal_server" "base" {
  zone		  = "fr-par-2"
  offer       = "GP-BM1-S"
  os          = "d17d6872-0412-45d9-a198-af82c34d3c5c"
  ssh_key_ids = [data.scaleway_account_ssh_key.main.id]
}
```

### With option

```terraform
data "scaleway_account_ssh_key" "main" {
  name = "main"
}

data "scaleway_baremetal_os" "my_os" {
  zone    = "fr-par-2"
  name    = "Ubuntu"
  version = "22.04 LTS (Jammy Jellyfish)"
}

data "scaleway_baremetal_offer" "my_offer" {
  zone = "fr-par-2"
  name = "EM-B112X-SSD"
}

data "scaleway_baremetal_option" "private_network" {
  zone = "fr-par-2"
  name = "Private Network"
}

data "scaleway_baremetal_option" "remote_access" {
  zone = "fr-par-2"
  name = "Remote Access"
}

resource "scaleway_baremetal_server" "base" {
  zone        = "fr-par-2"
  offer       = data.scaleway_baremetal_offer.my_offer.offer_id
  os          = data.scaleway_baremetal_os.my_os.os_id
  ssh_key_ids = [data.scaleway_account_ssh_key.main.id]

  options {
    id = data.scaleway_baremetal_option.private_network.option_id
  }

  options {
    id = data.scaleway_baremetal_option.remote_access.option_id
  }
}
```

### With private network

```terraform
data "scaleway_account_ssh_key" "main" {
  name = "main"
}

data "scaleway_baremetal_os" "my_os" {
  zone    = "fr-par-2"
  name    = "Ubuntu"
  version = "22.04 LTS (Jammy Jellyfish)"
}

data "scaleway_baremetal_offer" "my_offer" {
  zone = "fr-par-2"
  name = "EM-B112X-SSD"
}

data "scaleway_baremetal_option" "private_network" {
  zone = "fr-par-2"
  name = "Private Network"
}

resource "scaleway_vpc_private_network" "pn" {
  region = "fr-par"
  name   = "baremetal_private_network"
}

resource "scaleway_baremetal_server" "base" {
  zone        = "fr-par-2"
  offer       = data.scaleway_baremetal_offer.my_offer.offer_id
  os          = data.scaleway_baremetal_os.my_os.os_id
  ssh_key_ids = [data.scaleway_account_ssh_key.main.id]

  options {
    id = data.scaleway_baremetal_option.private_network.option_id
  }
  private_network {
    id = scaleway_vpc_private_network.pn.id
  }
}
```

### Without install config

```terraform
data "scaleway_baremetal_offer" "my_offer" {
  zone = "fr-par-2"
  name = "EM-B112X-SSD"
}

resource "scaleway_baremetal_server" "base" {
  zone	                   = "fr-par-2"
  offer                    = data.scaleway_baremetal_offer.my_offer.offer_id
  install_config_afterward = true
}
```

## Argument Reference

The following arguments are supported:

- `offer` - (Required) The offer name or UUID of the baremetal server.
  Use [this endpoint](https://www.scaleway.com/en/developers/api/elastic-metal/#path-servers-get-a-specific-elastic-metal-server) to find the right offer.

~> **Important:** Updates to `offer` will recreate the server.

- `os` - (Required) The UUID of the os to install on the server.
  Use [this endpoint](https://www.scaleway.com/en/developers/api/elastic-metal/#path-os-list-available-oses) to find the right OS ID.
  ~> **Important:** Updates to `os` will reinstall the server.
- `ssh_key_ids` - (Required) List of SSH keys allowed to connect to the server.
- `user` - (Optional) User used for the installation.
- `password` - (Optional) Password used for the installation. May be required depending on used os.
- `service_user` - (Optional) User used for the service to install.
- `service_password` - (Optional) Password used for the service to install. May be required depending on used os.
- `reinstall_on_config_changes` - (Optional) If True, this boolean allows to reinstall the server on install config changes.
  ~> **Important:** Updates to `ssh_key_ids`, `user`, `password`, `service_user` or `service_password` will not take effect on the server, it requires to reinstall it. To do so please set 'reinstall_on_config_changes' argument to true.
- `install_config_afterward` - (Optional) If True, this boolean allows to create a server without the install config if you want to provide it later.
- `name` - (Optional) The name of the server.
- `hostname` - (Optional) The hostname of the server.
- `description` - (Optional) A description for the server.
- `tags` - (Optional) The tags associated with the server.
- `options` - (Optional) The options to enable on the server.
  ~> The `options` block supports:
    - `id` - (Required) The id of the option to enable. Use [this endpoint](https://www.scaleway.com/en/developers/api/elastic-metal/#path-options-list-options) to find the available options IDs.
    - `expires_at` - (Optional) The auto expiration date for compatible options
- `private_network` - (Required) The private networks to attach to the server. For more information, see [the documentation](https://www.scaleway.com/en/docs/compute/elastic-metal/how-to/use-private-networks/)
    - `id` - (Required) The id of the private network to attach.
- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the server should be created.
- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the server is associated with.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the server.

~> **Important:** Baremetal servers' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

- `offer_id` - The ID of the offer.
- `offer_name` - The name of the offer.
- `os_name` - The name of the os.
- `private_network` - The private networks attached to the server.
    - `id` - The ID of the private network.
    - `vlan` - The VLAN ID associated to the private network.
    - `status` - The private network status.
    - `created_at` - The date and time of the creation of the private network.
    - `updated_at` - The date and time of the last update of the private network.
- `private_ip` - The list of private IP addresses associated with the resource.
    - `id` - The ID of the IP address resource.
    - `address` - The private IP address.
- `ips` - (List of) The IPs of the server.
    - `id` - The ID of the IP.
    - `address` - The address of the IP.
    - `reverse` - The reverse of the IP.
    - `version` - The type of the IP.
- `ipv4` - (List of) The IPv4 addresses of the server.
    - `id` - The ID of the IPv4.
    - `address` - The address of the IPv4.
    - `reverse` - The reverse of the IPv4.
    - `version` - The type of the IPv4.
- `ipv6` - (List of) The IPv6 addresses of the server.
    - `id` - The ID of the IPv6.
    - `address` - The address of the IPv6.
    - `reverse` - The reverse of the IPv6.
    - `version` - The type of the IPv6.
- `domain` - The domain of the server.
- `organization_id` - The organization ID the server is associated with.

## Import

Baremetal servers can be imported using the `{zone}/{id}`, e.g.

```bash
terraform import scaleway_baremetal_server.web fr-par-2/11111111-1111-1111-1111-111111111111
```
