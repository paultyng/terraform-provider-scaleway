---
subcategory: "Account"
page_title: "Scaleway: scaleway_account_project"
---

# Resource: scaleway_account_project

Manages organization's projects on Scaleway.

## Example Usage

### Basic

```terraform
resource "scaleway_account_project" "project" {
  name = "project"
}
```

### Use project in provider configuration

If you want to use as default a project created in terraform you can use a temporary provider alias.
This project can then be used to configure your default provider.

```terraform
provider "scaleway" {
  alias = "tmp"
}

resource scaleway_account_project "project" {
  provider = scaleway.tmp
  name = "my_project"
}

provider "scaleway" {
  project_id = scaleway_account_project.project.id
}

resource "scaleway_instance_server" "server" { // Will use scaleway_account_project.project
  image = "ubuntu_jammy"
  type  = "PRO2-XXS"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The name of the Project.
- `description` - (Optional) The description of the Project.
- `organization_id` - (Optional. Defaults to [provider](../index.md#organization_id) `organization_id`)The organization ID the Project is associated with. Please note that any change in `organization_id` will recreate the resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the project (UUID format).
- `created_at` - The Project creation time.
- `updated_at` - The Project last update time.

## Import

Projects can be imported using the `id`, e.g.

```bash
$ terraform import scaleway_account_project.project 11111111-1111-1111-1111-111111111111
```
