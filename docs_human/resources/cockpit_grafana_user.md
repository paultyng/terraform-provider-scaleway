---
subcategory: "Cockpit"
page_title: "Scaleway: scaleway_cockpit_grafana_user"
---

# Resource: scaleway_cockpit_grafana_user

Creates and manages Scaleway Cockpit Grafana Users.

For more information consult the [documentation](https://www.scaleway.com/en/docs/observability/cockpit/concepts/#grafana-users).

## Example Usage

```terraform
// Get the cockpit of the default project
data "scaleway_cockpit" "main" {}

// Create an editor grafana user for the cockpit
resource "scaleway_cockpit_grafana_user" "main" {
  project_id = data.scaleway_cockpit.main.project_id
  
  login = "my-awesome-user"
  role = "editor"
}
```


## Argument Reference

- `login` - (Required) The login of the grafana user.
- `role` - (Required) The role of the grafana user. Must be `editor` or `viewer`.
- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the cockpit is associated with.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `password` - The password of the grafana user

## Import

Cockpits Grafana Users can be imported using the project ID and the grafana user ID formatted `{project_id}/{grafana_user_id}`, e.g.

```bash
$ terraform import scaleway_cockpit_grafana_user.main 11111111-1111-1111-1111-111111111111/2
```
