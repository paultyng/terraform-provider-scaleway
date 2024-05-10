---
subcategory: "Account"
page_title: "Scaleway: scaleway_availability_zones"
---

# scaleway_availability_zones

Use this data source to get the available zones information based on its Region.

For technical and legal reasons, some products are split by Region or by Availability Zones. When using such product,
you can choose the location that better fits your need (country, latency, …).

## Example Usage

```hcl
# Get info by Region key
data scaleway_availability_zones main {
  region = "nl-ams"
}
```

## Argument Reference

- `region` - Region is represented as a Geographical area such as France. Defaults: `fr-par`.

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `id` - The Region ID
- `zones` - List of availability zones by regions
