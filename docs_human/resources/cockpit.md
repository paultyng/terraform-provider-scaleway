---
subcategory: "Cockpit"
page_title: "Scaleway: scaleway_cockpit"
---

# Resource: scaleway_cockpit

Creates and manages Scaleway Cockpit.

For more information consult the [documentation](https://www.scaleway.com/en/docs/observability/cockpit/concepts/).

## Example Usage

### Activate Cockpit in the default project

```terraform
resource "scaleway_cockpit" "main" {}
```

### Activate Cockpit in a specific project

```terraform
resource "scaleway_cockpit" "main" {
  project_id = "11111111-1111-1111-1111-111111111111"
}
```

### Use the Grafana Terraform provider

```terraform
resource "scaleway_cockpit" "main" {}

resource "scaleway_cockpit_grafana_user" "main" {
  project_id = scaleway_cockpit.main.project_id
  login      = "example"
  role       = "editor"
}

provider "grafana" {
  url  = scaleway_cockpit.main.endpoints.0.grafana_url
  auth = "${scaleway_cockpit_grafana_user.main.login}:${scaleway_cockpit_grafana_user.main.password}"
}

resource "grafana_folder" "test_folder" {
  title = "Test Folder"
}
```

## Argument Reference

- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the cockpit is associated with.
- `plan` - (Optional) Name or ID of the plan to use.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `plan_id` - The ID of the current plan.
- `endpoints` - Endpoints.
    - `metrics_url` - The metrics URL.
    - `logs_url` - The logs URL.
    - `alertmanager_url` - The alertmanager URL.
    - `grafana_url` - The grafana URL.
    - `traces_url` - The traces URL.

## Import

Cockpits can be imported using the `{project_id}`, e.g.

```bash
$ terraform import scaleway_cockpit.main 11111111-1111-1111-1111-111111111111
```
