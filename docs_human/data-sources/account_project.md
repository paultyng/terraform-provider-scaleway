---
subcategory: "Account"
page_title: "Scaleway: scaleway_account_project"
---

# scaleway_account_project

Gets information about an existing Project.

## Example Usage

```hcl
# Get info by name
data scaleway_account_project "by_name" {
  name            = "myproject"
  organization_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
# Get default project
data scaleway_account_project "by_name" {
  name = "default"
}
# Get info by ID
data scaleway_account_project "by_id" {
  project_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Argument Reference

- `name` - (Optional) The name of the Project.
  Only one of the `name` and `project_id` should be specified.

- `project_id` - (Optional) The ID of the Project.
  Only one of the `name` and `project_id` should be specified.

- `organization_id` - (Optional) The organization ID the Project is associated with.
  If no default organization_id is set, one must be set explicitly in this datasource

## Attribute Reference

Exported attributes are the ones from `account_project` [resource](../resources/account_project.md)
