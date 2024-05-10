---
subcategory: "Elastic Metal"
page_title: "Scaleway: scaleway_baremetal_os"
---

# scaleway_baremetal_os

Gets information about a baremetal operating system.
For more information, see [the documentation](https://developers.scaleway.com/en/products/baremetal/api).

You can also use the [scaleway-cli](https://github.com/scaleway/scaleway-cli) with `scw baremetal os list` to list all available operating systems.

## Example Usage

```hcl
# Get info by os name and version
data "scaleway_baremetal_os" "by_name" {
  name = "Ubuntu"
  version = "20.04 LTS (Focal Fossa)"
}

# Get info by os id
data "scaleway_baremetal_os" "by_id" {
  os_id = "03b7f4ba-a6a1-4305-984e-b54fafbf1681"
}
```

## Argument Reference

- `name` - (Optional) The os name. Only one of `name` and `os_id` should be specified.
- `version` - (Optional) The os version.
- `os_id` - (Optional) The operating system id. Only one of `name` and `os_id` should be specified.
- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the os exists.

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `id` - The resource's ID

~> **Important:** Baremetal operating systems' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`
