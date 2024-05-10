---
subcategory: "Messaging and Queuing"
page_title: "Scaleway: scaleway_mnq_sns_topic"
---

# Resource: scaleway_mnq_sns_topic

Manage Scaleway Messaging and queuing SNS Topics.
For further information please check
our [documentation](https://www.scaleway.com/en/docs/serverless/messaging/reference-content/sns-overview/)

## Example Usage

### Basic

```terraform
resource "scaleway_mnq_sns" "main" {}

resource scaleway_mnq_sns_credentials main {
  project_id = scaleway_mnq_sns.main.project_id
  permissions {
    can_manage = true
  }
}

resource "scaleway_mnq_sns_topic" "topic" {
  project_id = scaleway_mnq_sns.main.project_id
  name = "my-topic"
  access_key = scaleway_mnq_sns_credentials.main.access_key
  secret_key = scaleway_mnq_sns_credentials.main.secret_key
}
```

## Argument Reference

The following arguments are supported:


- `name` - (Optional) The unique name of the sns topic. Either `name` or `name_prefix` is required. Conflicts with `name_prefix`.

- `name_prefix` - (Optional) Creates a unique name beginning with the specified prefix. Conflicts with `name`.

- `access_key` - (Optional) The access key of the SNS credentials.

- `secret_key` - (Optional) The secret key of the SNS credentials.

- `content_based_deduplication` - (Optional) Specifies whether to enable content-based deduplication.

- `fifo_topic` - (Optional) Whether the topic is a FIFO. If true, the topic name must end with .fifo.

- `sns_endpoint` - (Optional) The endpoint of the SNS service. Can contain a {region} placeholder. Defaults to `https://sns.mnq.{region}.scaleway.com`.

- `region` - (Defaults to [provider](../index.md#region) `region`). The [region](../guides/regions_and_zones.md#regions)
  in which sns is enabled.

- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the sns is enabled for.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the topic with format `{region}/{project-id}/{topic-name}`

- `arn` - The ARN of the topic

## Import

SNS topic can be imported using the `{region}/{project-id}/{topic-name}`, e.g.

```bash
$ terraform import scaleway_mnq_sns_topic.main fr-par/11111111111111111111111111111111/my-topic
```
