---
subcategory: "Transactional Email"
page_title: "Scaleway: scaleway_tem_domain"
---

# scaleway_tem_domain

Gets information about a transactional email domain.

## Example Usage

```hcl
// Get info by domain name
data "scaleway_tem_domain" "my_domain" {
  name = "example.com"
}

// Get info by domain ID
data "scaleway_tem_domain" "my_domain" {
  domain_id = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

- `name` - (Optional) The domain name.
  Only one of `name` and `domain_id` should be specified.

- `domain_id` - (Optional) The domain id.
  Only one of `name` and `domain_id` should be specified.

- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#regions) in which the domain exists.

- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the domain is associated with.

## Attributes Reference

Exported attributes are the ones from `scaleway_tem_domain` [resource](../resources/tem_domain.md)
