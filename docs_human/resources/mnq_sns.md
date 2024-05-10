---
subcategory: "Messaging and Queuing"
page_title: "Scaleway: scaleway_mnq_sns"
---

# Resource: scaleway_mnq_sns

Activate Scaleway Messaging and queuing SNS for a project.
For further information please check
our [documentation](https://www.scaleway.com/en/docs/serverless/messaging/reference-content/sns-overview/)

## Example Usage

### Basic

Activate SNS for default project

```terraform
resource scaleway_mnq_sns "main" {}
```

Activate SNS for a specific project

```terraform
data scaleway_account_project project {
  name = "default"
}

// For specific project in default region
resource scaleway_mnq_sns "for_project" {
  project_id = data.scaleway_account_project.project.id
}
```

## Argument Reference

The following arguments are supported:


- `region` - (Defaults to [provider](../index.md#region) `region`). The [region](../guides/regions_and_zones.md#regions)
  in which sns will be enabled.

- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the sns will be enabled for.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the project

~> **Important:** Messaging and Queueing sns' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`

- `endpoint` - The endpoint of the SNS service for this project.

## Import

SNS status can be imported using the `{region}/{project_id}`, e.g.

```bash
$ terraform import scaleway_mnq_sns.main fr-par/11111111111111111111111111111111
```
