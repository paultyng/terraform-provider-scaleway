---
subcategory: "VPC"
page_title: "Scaleway: scaleway_vpc"
---

# Resource: scaleway_vpc

Creates and manages Scaleway Virtual Private Clouds.
For more information, see [the documentation](https://www.scaleway.com/en/docs/network/vpc/concepts/).

## Example Usage

```terraform
resource "scaleway_vpc" "vpc01" {
    name = "my-vpc"
    tags = ["demo", "terraform"]
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The name of the VPC. If not provided it will be randomly generated.
- `tags` - (Optional) The tags associated with the VPC.
- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#regions) of the VPC.
- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the VPC is associated with.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the VPC.
- `is_default` - Defines whether the VPC is the default one for its Project.
- `created_at` - Date and time of VPC's creation (RFC 3339 format).
- `updated_at` - Date and time of VPC's last update (RFC 3339 format).

~> **Important:** VPCs' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111

- `organization_id` - The organization ID the VPC is associated with.

## Import

VPCs can be imported using the `{region}/{id}`, e.g.

```bash
$ terraform import scaleway_vpc.vpc_demo fr-par/11111111-1111-1111-1111-111111111111
```
