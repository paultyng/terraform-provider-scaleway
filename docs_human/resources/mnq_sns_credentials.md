---
subcategory: "Messaging and Queuing"
page_title: "Scaleway: scaleway_mnq_sns_credentials"
---

# Resource: scaleway_mnq_sns_credentials

Creates and manages Scaleway Messaging and queuing SNS Credentials.
For further information please check
our [documentation](https://www.scaleway.com/en/docs/serverless/messaging/reference-content/sns-overview/)

## Example Usage

### Basic

```terraform
resource "scaleway_mnq_sns" "main" {}

resource scaleway_mnq_sns_credentials main {
  project_id = scaleway_mnq_sns.main.project_id
  name = "sns-credentials"

  permissions {
    can_manage = false
    can_receive = true
    can_publish = false
  }
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The unique name of the sns credentials.

- `permissions` - (Optional). List of permissions associated to these credentials. Only one of permissions may be set.
    - `can_publish` - (Optional). Defines if user can publish messages to the service.
    - `can_receive` - (Optional). Defines if user can receive messages from the service.
    - `can_manage` - (Optional). Defines if user can manage the associated resource(s).


- `region` - (Defaults to [provider](../index.md#region) `region`). The [region](../guides/regions_and_zones.md#regions) in which sns is enabled.

- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the sns is enabled for.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the credentials

~> **Important:** Messaging and Queueing sns credentials' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`

- `access_key` - The ID of the key.
- `secret_key` - The secret value of the key.

## Import

SNS credentials can be imported using the `{region}/{id}`, e.g.

```bash
$ terraform import scaleway_mnq_sns_credentials.main fr-par/11111111111111111111111111111111
```
