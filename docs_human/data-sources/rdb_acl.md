---
subcategory: "Databases"
page_title: "Scaleway: scaleway_rdb_acl"
---

# scaleway_rdb_acl

Gets information about the RDB instance network Access Control List.

## Example Usage

```hcl
# Get the database ACL for the instance id 11111111-1111-1111-1111-111111111111 located in the default region e.g: fr-par
data "scaleway_rdb_acl" "my_acl" {
  instance_id = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

- `instance_id` - (Required) The RDB instance ID.
- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#regions) in which the Database Instance should be created.

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `id` - The ID of the ACL.

~> **Important:** RDB instances ACLs' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`

- `acl_rules` - A list of ACLs rules (structure is described below)

The `acl_rules` block supports:

- `ip` - The ip range to whitelist in [CIDR notation](https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing#CIDR_notation)
- `description` - A simple text describing this rule
