---
subcategory: "Jobs"
page_title: "Scaleway: scaleway_job_definition"
---

# Resource: scaleway_job_definition

Creates and manages a Scaleway Serverless Job Definition. For more information, see [the documentation](https://pkg.go.dev/github.com/scaleway/scaleway-sdk-go@master/api/jobs/v1alpha1).

## Example Usage

### Basic

```terraform
resource scaleway_job_definition main {
  name = "testjob"
  cpu_limit = 140
  memory_limit = 256
  image_uri = "docker.io/alpine:latest"
  command = "ls"
  timeout = "10m"

  env = {
    foo: "bar"
  }

  cron {
    schedule = "5 4 1 * *" # cron at 04:05 on day-of-month 1
    timezone = "Europe/Paris"
  }
}
```

## Argument Reference

The following arguments are supported:

- `cpu_limit` - (Required) The amount of vCPU computing resources to allocate to each container running the job.
- `memory_limit` - (Required) The memory computing resources in MB to allocate to each container running the job.
- `image_uri` - (Required) The uri of the container image that will be used for the job run.
- `name` - (Optional) The name of the job.
- `description` - (Optional) The description of the job
- `command` - (Optional) The command that will be run in the container if specified.
- `timeout` - (Optional) The job run timeout, in Go Time format (ex: `2h30m25s`)
- `env` - (Optional) The environment variables of the container.
- `cron` - (Optional) The cron configuration
    - `schedule` - Cron format string.
    - `timezone` - The timezone, must be a canonical TZ identifier as found in this [list](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones).
- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#regions) of the Job.
- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the Job is associated with.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the Job Definition.

~> **Important:** Serverless Jobs Definition's IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`

## Import

Serverless Jobs can be imported using the `{region}/{id}`, e.g.

```bash
$ terraform import scaleway_job_definition.job fr-par/11111111-1111-1111-1111-111111111111
```
