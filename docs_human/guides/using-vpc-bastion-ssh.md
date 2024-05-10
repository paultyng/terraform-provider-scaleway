---
page_title: "Using Scaleway SSH Bastion"
---

# How to use Scaleway VPC SSH Bastion config

In this guide you'll learn how to deploy Scaleway SSH bastion to your Scaleway Private Network using the Scaleway Terraform provider.
After Bastion is deployed, you can connect (SSH) to virtual machines in the virtual network via Bastion using the private IP address of the VM.
When you connect to a VM, it doesn't need a public IP address, client software, agent, or a special configuration.

## Prerequisites

*	You have created a virtual machine (Instance) in a VPC Private Network. Check our example below.

1. When you deploy Bastion, the values are pulled from the Private Network in which your VM resides.
   1. This VM doesn't become a part of the Bastion deployment itself, but you do connect to it later in the exercise.

2. If you don't have any VMs connected to the Private Network, use the `scaleway_instance_private_nic` or the attribute `private_network` on `scaleway_instance_server` to connect.

3. Detach any VMs that are attached to a `scaleway_instance_ip`.

  **Note**: Your VMs and Private Network should be in the same Availability Zone. e.g. `fr-par-1`

```hcl
provider "scaleway" {
	zone = "fr-par-1"
}
```

```hcl
variable "machine_count" {
  description = "Number of virtual machines in private network"
  default = 3
}

# SCALEWAY VPC PRIVATE NETWORK
resource scaleway_vpc_private_network "pn" {
  name = "myprivatenetwork"
  zone = "fr-par-1"
}

# SCALEWAY VPC VIRTUAL MACHINES
resource scaleway_instance_server "servers" {
  count	= var.machine_count
  name 	= "machine${count.index}"
  image = "ubuntu_focal"
  type  = "DEV1-S"
}

# SCALEWAY INSTANCES PRIVATE NETWORK CONNECTION
resource scaleway_instance_private_nic "nic" {
  count              = length(scaleway_instance_server.servers)
  private_network_id = scaleway_vpc_private_network.pn.id
  server_id          = scaleway_instance_server.servers[count.index].id
}
```

## Reserve a public gateway IP

Reserve your public IP, allowing it to reach the public Internet, as well as to forward (masquerade) traffic from member Instances of attached Private Networks.

This IP is a static IPv4 address designed for dynamic cloud computing.

```hcl
# SCALEWAY PUBLIC GATEWAY IP
resource scaleway_vpc_public_gateway_ip "pgw_ip" {
}
```

## Set up your Public Gateway

Public Gateways sit at the border of Private Networks and allow you to enable the bastion.
You can also choose your port of preference on `bastion_port` option. The default port is `61000`

You can check the types of gateways currently supported via our CLI.

```shell
scw vpc-gw gateway-type list
```

Example:

```hcl
resource scaleway_vpc_public_gateway "pgw" {
  type = "VPC-GW-S"
  bastion_enabled = true
  ip_id = scaleway_vpc_public_gateway_ip.pgw_ip.id
}
```

## Configure your DHCP on your subnet

The [DHCP](https://fr.wikipedia.org/wiki/Dynamic_Host_Configuration_Protocol) server sets the IPv4 address dynamically,
which is required to communicate over the private network.

The `dns_local_name` is the [TLD](https://en.wikipedia.org/wiki/Top-level_domain), the value by default is `priv`.
This is used to resolve your Instance on a Private Network.

In order to resolve the Instances using your Bastion you should set the `dns_local_name` with `scaleway_vpc_private_network.pn.name`.

Please check our API [documentation](https://www.scaleway.com/en/developers/api/public-gateway/#path-dhcp-create-a-dhcp-configuration) for more details.

```hcl
resource scaleway_vpc_public_gateway_dhcp "dhcp" {
  subnet = "192.168.1.0/24"
  dns_local_name = scaleway_vpc_private_network.pn.name
}
```

## Attach your VPC Gateway Network to a Private Network

To enable DHCP on this Private Network you must set `enable_dhcp` and `dhcp_id`.
Do not set the `address` attribute.

```hcl
resource scaleway_vpc_gateway_network "gn" {
  gateway_id          = scaleway_vpc_public_gateway.pgw.id
  private_network_id  = scaleway_vpc_private_network.pn.id
  dhcp_id             = scaleway_vpc_public_gateway_dhcp.dhcp.id
  enable_dhcp         = true
}
```

## Config my Bastion config

You should add your config on your local config file e.g: `~/.ssh/config`

```
Host *.myprivatenetwork
ProxyJump bastion@<your-public-ip>:<bastion_port>
```

Then try to connect to it:

```shell
ssh root@<vm-name>.myprivatenetwork
```

For further information using our console please check [our dedicated documentation](https://www.scaleway.com/en/docs/network/vpc/how-to/use-ssh-bastion/).