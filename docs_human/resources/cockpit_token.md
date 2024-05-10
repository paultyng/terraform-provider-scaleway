---
subcategory: "Cockpit"
page_title: "Scaleway: scaleway_cockpit_token"
---

# Resource: scaleway_cockpit_token

Creates and manages Scaleway Cockpit Tokens.

For more information consult the [documentation](https://www.scaleway.com/en/docs/observability/cockpit/concepts/#tokens).

## Example Usage

```terraform
// Get the cockpit of the default project
data "scaleway_cockpit" "main" {}

// Create a token for the cockpit that can write metrics and logs
resource "scaleway_cockpit_token" "main" {
  project_id = data.scaleway_cockpit.main.project_id
  
  name = "my-awesome-token"
}
```

```terraform
// Get the cockpit of the default project
data "scaleway_cockpit" "main" {}

// Create a token for the cockpit that can read metrics and logs but not write
resource "scaleway_cockpit_token" "main" {
  project_id = data.scaleway_cockpit.main.project_id
  
  name = "my-awesome-token"
  scopes {
    query_metrics = true
    write_metrics = false

    query_logs = true
    write_logs = false
  }
}
```

## Argument Reference

- `name` - (Required) The name of the token.
- `scopes` - (Optional) Allowed scopes.
    - `query_metrics` - (Defaults to `false`) Query metrics.
    - `write_metrics` - (Defaults to `true`) Write metrics.
    - `setup_metrics_rules` - (Defaults to `false`) Setup metrics rules.
    - `query_logs` - (Defaults to `false`) Query logs.
    - `write_logs` - (Defaults to `true`) Write logs.
    - `setup_logs_rules` - (Defaults to `false`) Setup logs rules.
    - `setup_alerts` - (Defaults to `false`) Setup alerts.
    - `query_traces` - (Defaults to `false`) Query traces.
    - `write_traces` - (Defaults to `false`) Write traces.
- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the cockpit is associated with.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `secret_key` - The secret key of the token.

## Import

Cockpits can be imported using the token ID, e.g.

```bash
$ terraform import scaleway_cockpit_token.main 11111111-1111-1111-1111-111111111111
```
