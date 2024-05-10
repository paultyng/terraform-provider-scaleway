---
subcategory: "IoT Hub"
page_title: "Scaleway: scaleway_iot_hub"
---

# scaleway_iot_hub

Gets information about an IOT Hub.

## Example Usage

```hcl
# Get info by name
data "scaleway_iot_hub" "my_hub" {
  name = "foobar"
}

# Get info by hub ID
data "scaleway_iot_hub" "my_hub" {
  hub_id = "11111111-1111-1111-1111-111111111111"
}

```

## Argument Reference

- `name` - (Optional) The name of the Hub.
  Only one of the `name` and `hub_id` should be specified.

- `hub_id` - (Optional) The Hub ID.
  Only one of the `name` and `hub_id` should be specified.

- `region` - (Default to [provider](../index.md) `region`) The [region](../guides/regions_and_zones.md#zones) in which the hub exists.

- `project_id` - (Optional) The ID of the project the hub is associated with.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the Hub.

~> **Important:** IoT Hub instances' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`
