---
subcategory: "Cockpit"
page_title: "Scaleway: scaleway_cockpit_plan"
---
# scaleway_cockpit_plan

Gets information about a Scaleway Cockpit plan.

## Example Usage

```hcl
data "scaleway_cockpit_plan" "premium" {
  name = "premium"
}

resource "scaleway_cockpit" "main" {
  plan = data.scaleway_cockpit_plan.premium.id
}
```

## Arguments Reference

- `name` - (Required) The name of the plan.
