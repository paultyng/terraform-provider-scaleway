---
subcategory: "IoT Hub"
page_title: "Scaleway: scaleway_iot_network"
---

# Resource: scaleway_iot_network

-> **Note:** This terraform resource is currently in beta and might include breaking change in future releases.

Creates and manages Scaleway IoT Networks. For more information, see the following:

- [API documentation](https://developers.scaleway.com/en/products/iot/api).
- [Product documentation](https://www.scaleway.com/en/docs/scaleway-iothub-networks/)

For more step-by-step instructions on how to setup the networks on the external providers backends, you can follow these guides:

- [Configuring the Sigfox backend](https://www.scaleway.com/en/docs/scaleway-iothub-networks/#-Configuring-the-Sigfox-backend)
- [Using the Rest Network](https://www.scaleway.com/en/docs/scaleway-iothub-networks/#-Using-the-Rest-Network)

## Example Usage

```terraform
resource "scaleway_iot_network" "main" {
	name   = "main"
	hub_id = scaleway_iot_hub.main.id
	type   = "sigfox"
}
resource "scaleway_iot_hub" "main" {
	name         = "main"
	product_plan = "plan_shared"
}
```

## Argument Reference

~> **Important:** Updates to any value will recreate the IoT Route.

The following arguments are supported:

- `name` - (Required) The name of the IoT Network you want to create (e.g. `my-net`).

- `hub_id` - (Required) The hub ID to which the Network will be attached to.

- `type` - (Required) The network type to create (e.g. `sigfox`).

- `topic_prefix` - (Optional) The prefix that will be prepended to all topics for this Network.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the Network.

~> **Important:** IoT networks' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`

- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#regions) in which the Network is attached to.
- `created_at` - The date and time the Network was created.
- `endpoint` - The endpoint to use when interacting with the network.
- `secret` - The endpoint key to keep secret.

## Import

IoT Networks can be imported using the `{region}/{id}`, e.g.

```bash
$ terraform import scaleway_iot_network.net01 fr-par/11111111-1111-1111-1111-111111111111
```

