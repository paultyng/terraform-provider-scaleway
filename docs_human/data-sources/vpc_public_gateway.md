---
subcategory: "VPC"
page_title: "Scaleway: scaleway_vpc_public_gateway"
---

# scaleway_vpc_public_gateway

Gets information about a public gateway.

## Example Usage

```hcl
resource "scaleway_vpc_public_gateway" "main" {
  name = "demo"
  type = "VPC-GW-S"
  zone = "nl-ams-1"
}

data "scaleway_vpc_public_gateway" "pg_test_by_name" {
  name = "${scaleway_vpc_public_gateway.main.name}"
  zone = "nl-ams-1"
}

data "scaleway_vpc_public_gateway" "pg_test_by_id" {
  public_gateway_id = "${scaleway_vpc_public_gateway.main.id}"
}
```

## Argument Reference

- `name` - (Required) Exact name of the public gateway.
- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which
  the public gateway should be created.
- `project_id` - (Optional) The ID of the project the public gateway is associated with.

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `id` - The ID of the public gateway.

~> **Important:** Public gateways' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`
