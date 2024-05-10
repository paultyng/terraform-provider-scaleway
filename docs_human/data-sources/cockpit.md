---
subcategory: "Cockpit"
page_title: "Scaleway: scaleway_cockpit"
---
# scaleway_cockpit

Gets information about the Scaleway Cockpit.

For more information consult the [documentation](https://www.scaleway.com/en/docs/observability/cockpit/concepts/).

## Example Usage

```hcl
// Get default project's cockpit
data "scaleway_cockpit" "main" {}
```

```hcl
// Get a specific project's cockpit
data "scaleway_cockpit" "main" {
  project_id = "11111111-1111-1111-1111-111111111111"
}
```

## Arguments Reference

- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the cockpit is associated with.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `plan_id` - The ID of the current plan
- `endpoints` - Endpoints
    - `metrics_url` - The metrics URL
    - `logs_url` - The logs URL
    - `alertmanager_url` - The alertmanager URL
    - `grafana_url` - The grafana URL
