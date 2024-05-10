---
subcategory: "Functions"
page_title: "Scaleway: scaleway_function"
---

# Resource: scaleway_function

Creates and manages Scaleway Functions.
For more information see [the documentation](https://developers.scaleway.com/en/products/functions/api/).

## Example Usage

### Basic

```terraform
resource "scaleway_function_namespace" "main" {
  name        = "main-function-namespace"
  description = "Main function namespace"
}

resource scaleway_function main {
  namespace_id = scaleway_function_namespace.main.id
  runtime      = "go118"
  handler      = "Handle"
  privacy      = "private"
}
```

### With sources and deploy

You create a zip of your function (ex: `zip function.zip -r go.mod go.sum handler.go`)

```terraform
resource "scaleway_function_namespace" "main" {
  name        = "main-function-namespace"
  description = "Main function namespace"
}

resource scaleway_function main {
  namespace_id = scaleway_function_namespace.main.id
  runtime      = "go118"
  handler      = "Handle"
  privacy      = "private"
  timeout      = 10
  zip_file     = "function.zip"
  zip_hash     = filesha256("function.zip")
  deploy       = true
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The unique name of the function.

~> **Important** Updates to `name` will recreate the function.

- `description` (Optional) The description of the function.

- `environment_variables` - The environment variables of the function.

- `secret_environment_variables` - (Optional) The [secret environment](https://www.scaleway.com/en/docs/compute/functions/concepts/#secrets) variables of the function.

- `privacy` - Privacy of the function. Can be either `private` or `public`. Read more on [authentication](https://developers.scaleway.com/en/products/functions/api/#authentication)

- `runtime` - Runtime of the function. Runtimes can be fetched using [specific route](https://developers.scaleway.com/en/products/functions/api/#get-f7de6a)

- `min_scale` - Minimum replicas for your function, defaults to 0, Note that a function is billed when it gets executed, and using a min_scale greater than 0 will cause your function container to run constantly.

- `max_scale` - Maximum replicas for your function (defaults to 20), our system will scale your functions automatically based on incoming workload, but will never scale the number of replicas above the configured max_scale.

- `memory_limit` - Memory limit in MB for your function, defaults to 128MB

- `handler` - Handler of the function. Depends on the runtime ([function guide](https://developers.scaleway.com/en/products/functions/api/#create-a-function))

- `timeout` - Holds the max duration (in seconds) the function is allowed for responding to a request

- `zip_file` - Location of the zip file to upload containing your function sources

- `zip_hash` - The hash of your source zip file, changing it will re-apply function. Can be any string, changing it will just trigger state change. You can use any terraform hash function to trigger a change on your zip change (see examples)

- `deploy` - Define if the function should be deployed, terraform will wait for function to be deployed. Function will get deployed if you change source zip

- `region` - (Defaults to [provider](../index.md#region) `region`). The [region](../guides/regions_and_zones.md#regions) in which the namespace should be created.

- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the namespace is associated with.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the function

~> **Important:** Functions' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`

- `namespace_id` - The namespace ID the function is associated with.
- `domain_name` - The native domain name of the function
- `organization_id` - The organization ID the function is associated with.
- `cpu_limit` - The CPU limit in mCPU for your function. More infos on resources [here](https://developers.scaleway.com/en/products/functions/api/#functions)

## Import

Functions can be imported using the `{region}/{id}`, e.g.

```bash
$ terraform import scaleway_function.main fr-par/11111111-1111-1111-1111-111111111111
```
