---
subcategory: "Containers"
page_title: "Scaleway: scaleway_container"
---
# scaleway_container

Gets information about the Scaleway Container.

For more information consult the [documentation](https://www.scaleway.com/en/docs/faq/serverless-containers/).

For more details about the limitation check [containers-limitations](https://www.scaleway.com/en/docs/compute/containers/reference-content/containers-limitations/).

You can check also our [containers guide](https://www.scaleway.com/en/docs/compute/containers/concepts/).

## Example Usage

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

## Arguments Reference

The following arguments are required:

- `name` - (Required) The unique name of the container name.

- `namespace_id` - (Required) The container namespace ID of the container.

- `project_id` - (Optional) The ID of the project the container is associated with.

~> **Important** Updates to `name` will recreate the container.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the container

~> **Important:** Containers' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`

- `description` The description of the container.

- `environment_variables` - The [environment](https://www.scaleway.com/en/docs/compute/containers/concepts/#environment-variables) variables of the container.

- `min_scale` - The minimum of running container instances continuously. Defaults to 0.

- `max_scale` - The maximum of number of instances this container can scale to. Default to 20.

- `memory_limit` - The memory computing resources in MB to allocate to each container. Defaults to 128.

- `cpu_limit` - The amount of vCPU computing resources to allocate to each container. Defaults  to 70.

- `timeout` - The maximum amount of time in seconds during which your container can process a request before we stop it. Defaults to 300s.

- `privacy` - The privacy type define the way to authenticate to your container. Please check our dedicated [section](https://developers.scaleway.com/en/products/containers/api/#protocol-9dd4c8).

- `registry_image` - The registry image address. e.g: **"rg.fr-par.scw.cloud/$NAMESPACE/$IMAGE"**.

- `registry_sha256` - The sha256 of your source registry image, changing it will re-apply the deployment. Can be any string.

- `max_concurrency` - The maximum number of simultaneous requests your container can handle at the same time. Defaults to 50.

- `domain_name` - The container domain name.

- `protocol` - The communication [protocol](https://developers.scaleway.com/en/products/containers/api/#protocol-9dd4c8) http1 or h2c. Defaults to http1.

- `port` - The port to expose the container. Defaults to 8080.

- `deploy` - Boolean indicating whether the container is on a production environment.

- `status` - The container status.

- `cron_status` - The cron status of the container.

- `error_message` - The error message of the container.

- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#regions) in which the container was created.
