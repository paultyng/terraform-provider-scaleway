---
subcategory: "Containers"
page_title: "Scaleway: scaleway_container"
---

# Resource: scaleway_container

Creates and manages Scaleway Container.

For more information consult the [documentation](https://www.scaleway.com/en/docs/faq/serverless-containers/).

For more details about the limitation check [containers-limitations](https://www.scaleway.com/en/docs/compute/containers/reference-content/containers-limitations/).

You can check also our [containers guide](https://www.scaleway.com/en/docs/compute/containers/concepts/).

## Example Usage

```terraform
resource scaleway_container_namespace main {
    name = "my-ns-test"
    description = "test container"
}

resource scaleway_container main {
    name = "my-container-02"
    description = "environment variables test"
    namespace_id = scaleway_container_namespace.main.id
    registry_image = "${scaleway_container_namespace.main.registry_endpoint}/alpine:test"
    port = 9997
    cpu_limit = 140
    memory_limit = 256
    min_scale = 3
    max_scale = 5
    timeout = 600
    max_concurrency = 80
    privacy = "private"
    protocol = "http1"
    deploy = true

    environment_variables = {
        "foo" = "var"
    }
    secret_environment_variables = {
      "key" = "secret"
    }
}
```

## Argument Reference

The following arguments are required:

- `name` - (Required) The unique name of the container name.

- `namespace_id` - (Required) The container namespace ID of the container.

~> **Important** Updates to `name` will recreate the container.

The following arguments are optional:

- `description` (Optional) The description of the container.

- `environment_variables` - (Optional) The [environment](https://www.scaleway.com/en/docs/compute/containers/concepts/#environment-variables) variables of the container.

- `secret_environment_variables` - (Optional) The [secret environment](https://www.scaleway.com/en/docs/compute/containers/concepts/#secrets) variables of the container.

- `min_scale` - (Optional) The minimum of running container instances continuously. Defaults to 0.

- `max_scale` - (Optional) The maximum of number of instances this container can scale to. Default to 20.

- `memory_limit` - (Optional) The memory computing resources in MB to allocate to each container. Defaults to 256.

- `cpu_limit` - (Optional) The amount of vCPU computing resources to allocate to each container. Defaults to 140.

- `timeout` - (Optional) The maximum amount of time in seconds during which your container can process a request before we stop it. Defaults to 300s.

- `privacy` - (Optional) The privacy type define the way to authenticate to your container. Please check our dedicated [section](https://developers.scaleway.com/en/products/containers/api/#protocol-9dd4c8).

- `registry_image` - (Optional) The registry image address. e.g: **"rg.fr-par.scw.cloud/$NAMESPACE/$IMAGE"**.

- `registry_sha256` - (Optional) The sha256 of your source registry image, changing it will re-apply the deployment. Can be any string.

- `max_concurrency` - (Optional) The maximum number of simultaneous requests your container can handle at the same time. Defaults to 50.

- `protocol` - (Optional) The communication [protocol](https://developers.scaleway.com/en/products/containers/api/#protocol-9dd4c8) http1 or h2c. Defaults to http1.

- `port` - (Optional) The port to expose the container. Defaults to 8080.

- `deploy` - (Optional) Boolean controlling whether the container is on a production environment.

Note that if you want to use your own configuration, you must consult our configuration [restrictions](https://www.scaleway.com/en/docs/compute/containers/reference-content/containers-limitations/#configuration-restrictions) section.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The container's ID.

~> **Important:** Containers' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`

- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#regions) in which the container was created.
- `status` - The container status.
- `cron_status` - The cron status of the container.
- `error_message` - The error message of the container.
- `domain_name` - The native domain name of the container

## Import

Container can be imported using the `{region}/{id}`, e.g.

```bash
$ terraform import scaleway_container.main fr-par/11111111-1111-1111-1111-111111111111
```

## Protocols

The supported protocols are:

* `h2c`: HTTP/2 over TCP.
* `http1`: Hypertext Transfer Protocol.

**Important:** For details about the protocols check [this](https://httpd.apache.org/docs/2.4/howto/http2.html)

## Privacy

By default, creating a container will make it `public`, meaning that anybody knowing the endpoint could execute it.
A container can be made `private` with the privacy parameter.

Please check our [authentication](https://developers.scaleway.com/en/products/containers/api/#protocol-9dd4c8) section

## Memory and vCPUs configuration

The vCPU represents a portion or share of the underlying, physical CPU that is assigned to a particular virtual machine (VM).

You may decide how much computing resources to allocate to each container.
The `memory_limit` (in MB) must correspond with the right amount of vCPU.

**Important:** The right choice for your container's resources is very important, as you will be billed based on compute usage over time and the number of Containers executions.

Please check our [price](https://www.scaleway.com/en/docs/faq/serverless-containers/#prices) section for more details.

| Memory (in MB) | vCPU |
|----------------|------|
| 128            | 70m  |
| 256            | 140m |
| 512            | 280m |
| 1024           | 560m |

**Note:** 560mCPU accounts roughly for half of one CPU power of a Scaleway General Purpose instance