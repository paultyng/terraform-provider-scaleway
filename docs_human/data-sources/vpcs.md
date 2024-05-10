---
subcategory: "VPC"
page_title: "Scaleway: scaleway_vpcs"
---

# scaleway_vpcs

Gets information about multiple Virtual Private Clouds.

## Example Usage

```hcl
# Find VPCs that share the same tags
data "scaleway_vpcs" "my_key" {
  tags = ["tag"]
}

# Find VPCs by name and region
data "scaleway_vpcs" "my_key" {
  name   = "tf-vpc-datasource"
  region = "nl-ams"
}
```

## Argument Reference

- `name` - (Optional) The VPC name used as filter. VPCs with a name like it are listed.

- `tags` - (Optional) List of tags used as filter. VPCs with these exact tags are listed.

- `region` - (Defaults to [provider](../index.md#region) `region`). The [region](../guides/regions_and_zones.md#regions) in which vpcs exist.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `vpcs` - List of found vpcs
    - `id` - The associated VPC ID.
      ~> **Important:** VPCs' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111
    - `is_default` - Defines whether the VPC is the default one for its Project.
    - `created_at` - Date and time of VPC's creation (RFC 3339 format).
    - `updated_at` - Date and time of VPC's last update (RFC 3339 format).
    - `organization_id` - The organization ID the VPC is associated with.
    - `project_id` - The ID of the project the VPC is associated with.