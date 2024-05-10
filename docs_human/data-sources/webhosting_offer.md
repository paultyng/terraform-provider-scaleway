---
subcategory: "Web Hosting"
page_title: "Scaleway: scaleway_webhosting_offer"
---

# scaleway_webhosting_offer

Gets information about a webhosting offer.

## Example Usage

```hcl
# Get info by offer name
data "scaleway_webhosting_offer" "by_name" { 
  name = "performance"
}

# Get info by offer id
data "scaleway_webhosting_offer" "by_id" {
  offer_id = "de2426b4-a9e9-11ec-b909-0242ac120002"
}
```

## Argument Reference

- `name` - (Optional) The offer name. Only one of `name` and `offer_id` should be specified.

- `offer_id` - (Optional) The offer id. Only one of `name` and `offer_id` should be specified.

- `region` - (Defaults to [provider](../index.md#zone) `region`) The [region](../guides/regions_and_zones.md#zones) in which offer exists.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `billing_operation_path` - The unique identifier used for billing.
- `product` - The offer product.
    - `option` - The product option.
    - `email_accounts_quota` - The quota of email accounts.
    - `email_storage_quota` - The quota of email storage.
    - `databases_quota` - The quota of databases.
    - `hosting_storage_quota` - The quota of hosting storage.
    - `support_included` - If support is included.
    - `v_cpu` - The number of cores.
    - `ram` - The capacity of the memory in GB.
- `price` - The offer price.

