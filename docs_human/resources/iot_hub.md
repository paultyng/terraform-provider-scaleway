---
subcategory: "IoT Hub"
page_title: "Scaleway: scaleway_iot_hub"
---

# Resource: scaleway_iot_hub

-> **Note:** This terraform resource is currently in beta and might include breaking change in future releases.

Creates and manages Scaleway IoT Hub Instances. For more information, see [the documentation](https://developers.scaleway.com/en/products/iot/api).

## Example Usage

### Basic

```terraform
resource "scaleway_iot_hub" "main" {
    name = "test-iot"
    product_plan = "plan_shared"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The name of the IoT Hub instance you want to create (e.g. `my-hub`).

- `product_plan` - (Required) Product plan to create the hub, see documentation for available product plans (e.g. `plan_shared`)

~> **Important:** Updates to `product_plan` will recreate the IoT Hub Instance.

- `enabled` - (Optional) Wether the IoT Hub instance should be enabled or not.

~> **Important:** Updates to `enabled` will disconnect eventually connected devices.

- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#regions) in which the Database Instance should be created.

- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the IoT Hub Instance is associated with.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the Hub.

~> **Important:** IoT Hub instances' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`

- `created_at` - The date and time the Hub was created.
- `updated_at` - The date and time the Hub resource was updated.
- `status` - The current status of the Hub.
- `endpoint` - The MQTT network endpoint to connect MQTT devices to.
- `device_count` - The number of registered devices in the Hub.
- `connected_device_count` - The current number of connected devices in the Hub.
- `mqtt_ca_url` - The MQTT ca url
- `mqtt_ca` - The MQTT certificat content


## Import

IoT Hubs can be imported using the `{region}/{id}`, e.g.

```bash
$ terraform import scaleway_iot_hub.hub01 fr-par/11111111-1111-1111-1111-111111111111
```
