---
subcategory: "Instances"
page_title: "Scaleway: scaleway_instance_security_group"
---

# scaleway_instance_security_group

Gets information about a Security Group.

## Example Usage

```hcl
# Get info by security group name
data "scaleway_instance_security_group" "my_key" {
  name  = "my-security-group-name"
}

# Get info by security group id
data "scaleway_instance_security_group" "my_key" {
  security_group_id = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

- `name` - (Optional) The security group name. Only one of `name` and `security_group_id` should be specified.

- `security_group_id` - (Optional) The security group id. Only one of `name` and `security_group_id` should be specified.

- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the security group exists.

- `project_id` - (Optional) The ID of the project the security group is associated with.

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `id` - The ID of the security group.

~> **Important:** Instance security groups' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

- `organization_id` - The ID of the organization the security group is associated with.

- `inbound_default_policy` - The default policy on incoming traffic. Possible values are: `accept` or `drop`.

- `outbound_default_policy` - The default policy on outgoing traffic. Possible values are: `accept` or `drop`.

- `inbound_rule` - A list of inbound rule to add to the security group. (Structure is documented below.)

- `outbound_rule` - A list of outbound rule to add to the security group. (Structure is documented below.)

The `inbound_rule` and `outbound_rule` block supports:

- `action` - The action to take when rule match. Possible values are: `accept` or `drop`.

- `protocol`- The protocol this rule apply to. Possible values are: `TCP`, `UDP`, `ICMP` or `ANY`.

- `port`- The port this rule apply to. If no port is specified, rule will apply to all port.

- `ip`- The ip this rule apply to.

- `ip_range`- The ip range (e.g `192.168.1.0/24`) this rule apply to.
