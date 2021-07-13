---
layout: "scaleway"
page_title: "Scaleway: scaleway_iot_device"
description: |-
  Manages Scaleway IoT Hub device.
---

# scaleway_iot_device

-> **Note:** This terraform resource is currently in beta and might include breaking change in future releases.

Creates and manages Scaleway IoT Hub Instances. For more information, see [the documentation](https://developers.scaleway.com/en/products/iot/api).

## Examples

### Basic

```hcl
resource scaleway_iot_hub main {
    name         = "test-iot"
    product_plan = "plan_shared"
}

resource scaleway_iot_device main {
    hub_id = scaleway_iot_hub.main.id
    name   = "test-iot"
}
```

## Arguments Reference

The following arguments are supported:

- `hub_id` - (Required) The ID of the hub on which this device will be created.

- `name` - (Required) The name of the IoT device you want to create (e.g. `my-device`).

~> **Important:** Updates to `name` will destroy and recreate a new resource.

- `description` - (Optional) The description of the IoT device (e.g. `living room`).

- `allow_insecure` - (Optional) Allow plain and server-authenticated TLS connections in addition to mutually-authenticated ones.

~> **Important:** Updates to `allow_insecure` can disconnect eventually connected devices.

- `allow_multiple_connections` - (Optional) Allow more than one simultaneous connection using the same device credentials.

~> **Important:** Updates to `allow_multiple_connections` can disconnect eventually connected devices.

- `message_filters` - (Optional) Rules that define which messages are authorized or denied based on their topic.
    - `publish` - (Optional) Rules used to restrict topics the device can publish to.
        - `policy` (Optional) Filtering policy (eg `accept` or `reject`)
        - `topics` (Optional) List of topics to match (eg `foo/bar/+/baz/#`)
    - `subscribe` - (Optional) Rules used to restrict topics the device can subscribe to.
        - `policy` (Optional) Same as publish rules.
        - `topics` (Optional) Same as publish rules.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the device.
- `created_at` - The date and time the device was created.
- `updated_at` - The date and time the device resource was updated.
- `certificate` - The certificate bundle of the device.
    - `crt` - The certificate of the device.
    - `key` - The private key of the device.
- `status` - The current status of the device.
- `last_activity_at` - The last MQTT activity of the device.
- `is_connected` - The current connection status of the device.


## Import

IoT devices can be imported using the `{region}/{id}`, e.g.

```bash
$ terraform import scaleway_iot_device.device01 fr-par/11111111-1111-1111-1111-111111111111
```
