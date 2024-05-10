---
subcategory: "Elastic Metal"
page_title: "Scaleway: scaleway_baremetal_offer"
---

# scaleway_baremetal_offer

Gets information about a baremetal offer. For more information, see [the documentation](https://developers.scaleway.com/en/products/baremetal/api).

## Example Usage

```hcl
# Get info by offer name
data "scaleway_baremetal_offer" "my_offer" {
  zone = "fr-par-2"
  name = "EM-A210R-SATA"
}

# Get info by offer id
data "scaleway_baremetal_offer" "my_offer" {
  zone     = "fr-par-2"
  offer_id = "25dcf38b-c90c-4b18-97a2-6956e9d1e113"
}
```

## Argument Reference

- `name` - (Optional) The offer name. Only one of `name` and `offer_id` should be specified.

- `subscription_period` - (Optional) Period of subscription the desired offer. Should be `hourly` or `monthly`.

- `offer_id` - (Optional) The offer id. Only one of `name` and `offer_id` should be specified.

- `allow_disabled` - (Optional, default `false`) Include disabled offers.

- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the offer should be created.

## Attributes Reference

In addition to all above arguments, the following attributes are exported:

- `id` - The ID of the offer.

~> **Important:** Baremetal offers' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

- `bandwidth` - Available Bandwidth with the offer.

- `commercial_range` - Commercial range of the offer.

- `cpu` - A list of cpu specifications. (Structure is documented below.)

- `disk` - A list of disk specifications. (Structure is documented below.)

- `memory` - A list of memory specifications. (Structure is documented below.)

- `stock` - Stock status for this offer. Possible values are: `empty`, `low` or `available`.

The `cpu` block supports:

- `name` - Name of the CPU.

- `core_count`- Number of core on this CPU.

- `frequency`- Frequency of the CPU in MHz.

- `thread_count`- Number of thread on this CPU.

The `disk` block supports:

- `type` - Type of disk.

- `capacity`- Capacity of the disk in GB.

The `memory` block supports:

- `type` - Type of memory.

- `capacity`- Capacity of the memory in GB.

- `frequency` - Frequency of the memory in MHz.

- `is_ecc`- True if error-correcting code is available on this memory.
