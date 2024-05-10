---
subcategory: "Containers"
page_title: "Scaleway: scaleway_container_domain"
---

# Resource: scaleway_container_domain

Creates and manages Scaleway Container domain name bindings.
You can check our [containers guide](https://www.scaleway.com/en/docs/compute/containers/how-to/add-a-custom-domain-to-a-container/) for further information.

## Example Usage

### Simple

```terraform
resource scaleway_container app {}

resource scaleway_container_domain "app" {
  container_id = scaleway_container.app.id
  hostname = "container.domain.tld"
}
```

### Complete example with domain

```terraform
resource scaleway_container_namespace main {
    name = "my-ns-test"
    description = "test container"
}

resource scaleway_container app {
    name = "app"
    namespace_id = scaleway_container_namespace.main.id
    registry_image = "${scaleway_container_namespace.main.registry_endpoint}/nginx:alpine"
    port = 80
    cpu_limit = 140
    memory_limit = 256
    min_scale = 1
    max_scale = 1
    timeout = 600
    max_concurrency = 80
    privacy = "public"
    protocol = "http1"
    deploy = true
}

resource scaleway_domain_record "app" {
  dns_zone = "domain.tld"
  name     = "subdomain"
  type     = "CNAME"
  data     = "${scaleway_container.app.domain_name}." // Trailing dot is important in CNAME
  ttl      = 3600
}

resource scaleway_container_domain "app" {
  container_id = scaleway_container.app.id
  hostname = "${scaleway_domain_record.app.name}.${scaleway_domain_record.app.dns_zone}"
}
```

## Argument Reference

The following arguments are required:

- `hostname` - (Required) The hostname with a CNAME record.

- `container_id` - (Required) The ID of the container.

- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#regions) in which the container exists

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The container domain's ID.

~> **Important:** Container domains' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`

- `url` - The URL used to query the container


## Import

Container domain binding can be imported using the `{region}/{id}`, e.g.

```bash
$ terraform import scaleway_container_domain.main fr-par/11111111-1111-1111-1111-111111111111
```
