---
subcategory: "Instances"
page_title: "Scaleway: scaleway_marketplace_image"
---

# scaleway_marketplace_image

Gets local image ID of an image from its label name.

## Example Usage

```hcl
data "scaleway_marketplace_image" "my_image" {
  label  = "ubuntu_jammy"
}
```

## Argument Reference

- `label` - (Required) Exact label of the desired image. You can use [this endpoint](https://api-marketplace.scaleway.com/images?page=1&per_page=100)
to find the right `label`.

- `instance_type` - (Optional, default `DEV1-S`) The instance type the image is compatible with.
You find all the available types on the [pricing page](https://www.scaleway.com/en/pricing/).

- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the image exists.

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `id` - The ID of the local image.

- ~> **Important:** Instance local images' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`
