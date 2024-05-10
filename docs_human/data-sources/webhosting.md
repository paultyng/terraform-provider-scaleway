---
subcategory: "Web Hosting"
page_title: "Scaleway: scaleway_webhosting"
---

# scaleway_webhosting

Gets information about a webhosting.

## Example Usage

```hcl
# Get info by offer domain
data "scaleway_webhosting" "by_domain" {
  domain = "foobar.com"
}

# Get info by id
data "scaleway_webhosting" "by_id" {
  webhosting_id = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

* `domain` - (Optional) The hosting domain name. Only one of `domain` and `webhosting_id` should be specified.
* `webhosting_id` - (Optional) The hosting id. Only one of `domain` and `webhosting_id` should be specified.
* `organization_id` - The ID of the organization the hosting is associated with.
* `project_id` - (Optional. Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the hosting is associated with.
* `region` - (Defaults to [provider](../index.md#zone) `region`) The [region](../guides/regions_and_zones.md#zones) in which hosting exists.

## Attributes Reference

Exported attributes are the ones from `scaleway_webhosting` [resource](../resources/webhosting.md)
