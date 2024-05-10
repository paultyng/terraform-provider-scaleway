---
subcategory: "Instances"
page_title: "Scaleway: scaleway_instance_private_nic"
---

# scaleway_instance_private_nic

Gets information about an instance private NIC.

## Example Usage

```hcl
data "scaleway_instance_private_nic" "by_nic_id" {
  server_id = "11111111-1111-1111-1111-111111111111"
  private_nic_id = "11111111-1111-1111-1111-111111111111"
}

data "scaleway_instance_private_nic" "by_pn_id" {
  server_id = "11111111-1111-1111-1111-111111111111"
  private_network_id = "11111111-1111-1111-1111-111111111111"
}

data "scaleway_instance_private_nic" "by_tags" {
  server_id = "11111111-1111-1111-1111-111111111111"
  tags = ["mytag"]
}
```

## Argument Reference

- `server_id` - (Required) The server's id

- `tags` (Optional) The tags associated with the private NIC.
  As datasource only returns one private NIC, the search with given tags must return only one result

- `private_nic_id` - (Optional) The ID of the instance server private nic
  Only one of `private_nic_id` and `private_network_id` should be specified.

- `private_network_id` - (Optional) The ID of the private network
  Only one of `private_nic_id` and `private_network_id` should be specified.

- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the private nic exists.

## Attributes Reference

Exported attributes are the ones from `instance_private_nic` [resource](../resources/instance_private_nic.md)
