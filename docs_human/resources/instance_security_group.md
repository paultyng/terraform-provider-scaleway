---
subcategory: "Instances"
page_title: "Scaleway: scaleway_instance_security_group"
---

# Resource: scaleway_instance_security_group

Creates and manages Scaleway Compute Instance security groups. For more information, see [the documentation](https://developers.scaleway.com/en/products/instance/api/#security-groups-8d7f89).

## Example Usage

### Basic

```terraform
resource "scaleway_instance_security_group" "allow_all" {
}

resource "scaleway_instance_security_group" "web" {
  inbound_default_policy = "drop" # By default we drop incoming traffic that do not match any inbound_rule

  inbound_rule {
    action = "accept"
    port   = 22
    ip     = "212.47.225.64"
  }

  inbound_rule {
    action = "accept"
    port   = 80
  }

  inbound_rule {
    action     = "accept"
    protocol   = "UDP"
    port_range = "22-23"
  }
}
```

### Web server with banned IP and restricted internet access

```terraform
resource "scaleway_instance_security_group" "web" {
  inbound_default_policy  = "drop" # By default we drop incoming traffic that do not match any inbound_rule.
  outbound_default_policy = "drop" # By default we drop outgoing traffic that do not match any outbound_rule.

  inbound_rule {
    action = "drop"
    ip     = "1.1.1.1" # Banned IP
  }

  inbound_rule {
    action = "accept"
    port   = 22
    ip     = "212.47.225.64"
  }

  inbound_rule {
    action = "accept"
    port   = 443
  }

  outbound_rule {
    action = "accept"
    ip     = "8.8.8.8" # Only allow outgoing connection to this IP.
  }
}
```

### Trusted IP for SSH access (using for_each)

If you use terraform >= 0.12.6, you can leverage the [`for_each`](https://www.terraform.io/docs/configuration/resources.html#for_each-multiple-resource-instances-defined-by-a-map-or-set-of-strings) feature with this resource.

```terraform
locals {
  trusted = ["192.168.0.1", "192.168.0.2", "192.168.0.3"]
}

resource "scaleway_instance_security_group" "dummy" {
  inbound_default_policy  = "drop"
  outbound_default_policy = "accept"

  dynamic "inbound_rule" {
    for_each = local.trusted

    content {
      action = "accept"
      port   = 22
      ip     = inbound_rule.value
    }
  }
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The name of the security group.

- `description` - (Optional) The description of the security group.

- `stateful` - (Defaults to `true`) A boolean to specify whether the security group should be stateful or not.

- `inbound_default_policy` - (Defaults to `accept`) The default policy on incoming traffic. Possible values are: `accept` or `drop`.

- `outbound_default_policy` - (Defaults to `accept`) The default policy on outgoing traffic. Possible values are: `accept` or `drop`.

- `inbound_rule` - (Optional) A list of inbound rule to add to the security group. (Structure is documented below.)

- `outbound_rule` - (Optional) A list of outbound rule to add to the security group. (Structure is documented below.)

- `external_rules` - (Defaults to `false`) A boolean to specify whether to use [instance_security_group_rules](../resources/instance_security_group_rules.md).
  If `external_rules` is set to `true`, `inbound_rule` and `outbound_rule` can not be set directly in the security group.

- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the security group should be created.


- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the security group is associated with.

- `enable_default_security` - Whether to block SMTP on IPv4/IPv6 (Port 25, 465, 587). Set to false will unblock SMTP if your account is authorized to. If your organization is not yet authorized to send SMTP traffic, [open a support ticket](https://console.scaleway.com/support/tickets).

The `inbound_rule` and `outbound_rule` block supports:

- `action` - (Required) The action to take when rule match. Possible values are: `accept` or `drop`.

- `protocol`- (Defaults to `TCP`) The protocol this rule apply to. Possible values are: `TCP`, `UDP`, `ICMP` or `ANY`.

- `port`- (Optional) The port this rule applies to. If no `port` nor `port_range` are specified, the rule will apply to all port. Only one of `port` and `port_range` should be specified.

- `port_range`- Need terraform >= 0.13.0 (Optional) The port range (e.g `22-23`) this rule applies to.
  Port range MUST comply the Scaleway-notation: interval between ports must be a power of 2 `2^X-1` number (e.g 2^13-1=8191 in port_range = "10000-18191").
  If no `port` nor `port_range` are specified, rule will apply to all port.
  Only one of `port` and `port_range` should be specified.

- `ip`- (Optional) The ip this rule apply to. If no `ip` nor `ip_range` are specified, rule will apply to all ip. Only one of `ip` and `ip_range` should be specified.

- `ip_range`- (Optional) The ip range (e.g `192.168.1.0/24`) this rule applies to. If no `ip` nor `ip_range` are specified, rule will apply to all ip. Only one of `ip` and `ip_range` should be specified.

- `tags`- (Optional) The tags of the security group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the security group.

~> **Important:** Instance security groups' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

- `organization_id` - The organization ID the security group is associated with.

## Import

Instance security group can be imported using the `{zone}/{id}`, e.g.

```bash
$ terraform import scaleway_instance_security_group.web fr-par-1/11111111-1111-1111-1111-111111111111
```
