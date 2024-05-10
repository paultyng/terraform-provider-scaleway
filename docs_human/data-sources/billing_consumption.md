---
subcategory: "Billing"
page_title: "Scaleway: scaleway_billing_consumptions"
---

# scaleway_billing_consumptions

Gets information about your Consumptions.

## Example Usage

```hcl
# Find your detailed monthly consumption list
data "scaleway_billing_consumptions" "my-consumption" {
  organization_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Argument Reference

- `organization_id` - (Defaults to [provider](../index.md#organization_d) `organization_id`) The ID of the organization the consumption list is associated with.
- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the consumption list is associated with.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `consumptions` - List of found consumptions
    - `value` - The monetary value of the consumption.
    - `product_name` - The product name.
    - `category_name` - The name of the consumption category.
    - `sku` - The unique identifier of the product.
    - `unit` - The unit of consumed quantity.
    - `billed_quantity` - The consumed quantity.
    - `project_id` - The project ID of the consumption.
- `updated_at` - The last consumption update date.