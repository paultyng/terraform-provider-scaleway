---
subcategory: "Containers"
page_title: "Scaleway: scaleway_container"
---
# scaleway_container

The `scaleway_container` data source is used to retrieve information about a Serverless Container.

Refer to the Serverless Containers [product documentation](https://www.scaleway.com/en/docs/serverless/containers/) and [API documentation](https://www.scaleway.com/en/developers/api/serverless-containers/) for more information.

For more information on the limitations of Serverless Containers, refer to the [dedicated documentation](https://www.scaleway.com/en/docs/serverless-containers/reference-content/containers-limitations/).

## Retrieve a Serverless Container

The following commands allow you to:

- retrieve a container by its name
- retrieve a container by its ID

```hcl
resource scaleway_container_namespace main {
}

resource scaleway_container main {
    name = "test-container-data"
    namespace_id = scaleway_container_namespace.main.id
}

// Get info by container name
data "scaleway_container" "by_name" {
    namespace_id = scaleway_container_namespace.main.id
    name = scaleway_container.main.name
}

// Get info by container ID
data "scaleway_container" "by_id" {
    namespace_id = scaleway_container_namespace.main.id
    container_id = scaleway_container.main.id
}
```

## Arguments reference

This section lists the arguments that you can provide to the `scaleway_container` data source to filter and retrieve the desired namespace. Each argument has a specific purpose:

- `name` - (Required) The unique name of the container.

- `namespace_id` - (Required) The container namespace ID of the container.

- `project_id` - (Optional) The unique identifier of the project with which the container is associated.

~> **Important** Updating the `name` argument will recreate the container.

## Attributes Reference

The `scaleway_container` data source exports certain attributes once the container information is retrieved. These attributes can be referenced in other parts of your Terraform configuration.

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the container

~> **Important:** Containers' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`

- `description` The description of the container.

- `environment_variables` - The [environment](https://www.scaleway.com/en/docs/serverless-containers/concepts/#environment-variables) variables of the container.

- `min_scale` - The minimum number of container instances running continuously.

- `max_scale` - The maximum number of instances the container can scale to.

- `memory_limit` - The memory resources in MB to allocate to each container.

- `cpu_limit` - The amount of vCPU computing resources to allocate to each container.

- `timeout` - The maximum amount of time your container can spend processing a request before being stopped.

- `privacy` - The privacy type define the way to authenticate to your container. Refer to the [dedicated documentation](https://www.scaleway.com/en/developers/api/serverless-containers/#path-containers-update-an-existing-container) for more information.

- `registry_image` - The registry image address (e.g. `rg.fr-par.scw.cloud/$NAMESPACE/$IMAGE`).

- `registry_sha256` - The sha256 of your source registry image, changing it will re-apply the deployment. Can be any string.

- `max_concurrency` - The maximum number of simultaneous requests your container can handle at the same time.

- `domain_name` - The container domain name.

- `protocol` - The communication [protocol](https://www.scaleway.com/en/developers/api/serverless-containers/#path-containers-update-an-existing-container) `http1` or `h2c`. Defaults to `http1`.

- `port` - The port to expose the container.

- `deploy` - Boolean indicating whether the container is on a production environment.

- `sandbox` - Execution environment of the container.

- `heath_check` - Health check configuration block of the container.
    - `http` - HTTP health check configuration.
        - `path` - Path to use for the HTTP health check.
    - `failure_threshold` - Number of consecutive health check failures before considering the container unhealthy.
    - `interval`- Period between health checks (in seconds).

- `status` - The container status.

- `cron_status` - The cron status of the container.

- `error_message` - The error message of the container.

- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#regions) in which the container was created.
