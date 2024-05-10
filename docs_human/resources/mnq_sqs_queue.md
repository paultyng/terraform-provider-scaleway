---
subcategory: "Messaging and Queuing"
page_title: "Scaleway: scaleway_mnq_sqs_queue"
---

# Resource: scaleway_mnq_sqs_queue

Creates and manages Scaleway Messaging and queuing SQS Queues.
For further information please check
our [documentation](https://www.scaleway.com/en/docs/serverless/messaging/how-to/create-manage-queues/)

## Example Usage

### Basic

```terraform
resource "scaleway_mnq_sqs" "main" {}

resource scaleway_mnq_sqs_credentials main {
  project_id = scaleway_mnq_sqs.main.project_id
  name = "sqs-credentials"

  permissions {
    can_manage = true
    can_receive = false
    can_publish = false
  }
}

resource scaleway_mnq_sqs_queue main {
  project_id = scaleway_mnq_sqs.main.project_id
  name = "my-queue"
  sqs_endpoint = scaleway_mnq_sqs.main.endpoint
  access_key = scaleway_mnq_sqs_credentials.main.access_key
  secret_key = scaleway_mnq_sqs_credentials.main.secret_key
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The unique name of the sqs queue. Either `name` or `name_prefix` is required. Conflicts with `name_prefix`.

- `name_prefix` - (Optional) Creates a unique name beginning with the specified prefix. Conflicts with `name`.

- `sqs_endpoint` - (Optional) The endpoint of the SQS queue. Can contain a {region} placeholder. Defaults to `https://sqs.mnq.{region}.scaleway.com`.

- `access_key` - (Required) The access key of the SQS queue.

- `secret_key` - (Required) The secret key of the SQS queue.

- `fifo_queue` - (Optional) Whether the queue is a FIFO queue. If true, the queue name must end with .fifo. Defaults to `false`.

- `content_based_deduplication` - (Optional) Specifies whether to enable content-based deduplication. Defaults to `false`.

- `receive_wait_time_seconds` - (Optional) The number of seconds to wait for a message to arrive in the queue before returning. Must be between 0 and 20. Defaults to 0.

- `visibility_timeout_seconds` - (Optional) The number of seconds a message is hidden from other consumers. Must be between 0 and 43_200. Defaults to 30.

- `message_max_age` - (Optional) The number of seconds the queue retains a message. Must be between 60 and 1_209_600. Defaults to 345_600.

- `message_max_size` - (Optional) The maximum size of a message. Should be in bytes. Must be between 1024 and 262_144. Defaults to 262_144.

- `region` - (Defaults to [provider](../index.md#region) `region`). The [region](../guides/regions_and_zones.md#regions) in which sqs is enabled.

- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the sqs is enabled for.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the queue with format `{region/{project-id}/{queue-name}`

- `url` - The URL of the queue.
