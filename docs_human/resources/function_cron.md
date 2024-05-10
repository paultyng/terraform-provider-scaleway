---
subcategory: "Functions"
page_title: "Scaleway: scaleway_function_cron"
---

# Resource: scaleway_function_cron

Creates and manages Scaleway Function Triggers. For the moment, the feature is limited to CRON Schedule (time-based).

For more details about the limitation
check [functions-limitations](https://www.scaleway.com/en/docs/compute/functions/reference-content/functions-limitations/).

You can check also
our [functions cron api documentation](https://developers.scaleway.com/en/products/functions/api/#crons-942bf4).

## Example Usage

```terraform
resource scaleway_function_namespace main {
    name = "test-cron"
}

resource scaleway_function main {
    name = "test-cron"
    namespace_id = scaleway_function_namespace.main.id
    runtime = "node14"
    privacy = "private"
    handler = "handler.handle"
}

resource scaleway_function_cron main {
    name = "test-cron"
    function_id = scaleway_function.main.id
    schedule = "0 0 * * *"
    args = jsonencode({test = "scw"})
}

resource scaleway_function_cron func {
    function_id = scaleway_function.main.id
    schedule = "0 1 * * *"
    args = jsonencode({my_var = "terraform"})
}
```

## Argument Reference

The following arguments are required:

- `schedule` - (Required) Cron format string, e.g. @hourly, as schedule time of its jobs to be created and
  executed.
- `function_id` - (Required) The function ID to link with your cron.
- `args`   - (Required) The key-value mapping to define arguments that will be passed to your function’s event object
  during
- `name` - (Optional) The name of the cron. If not provided, the name is generated.
- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#regions)
  in where the job was created.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The function CRON's ID.

~> **Important:** Function CRONs' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`

- `status` - The cron status.

## Import

Container Cron can be imported using the `{region}/{id}`, e.g.

```bash
$ terraform import scaleway_function_cron.main fr-par/11111111-1111-1111-1111-111111111111
```
