---
subcategory: "Functions"
page_title: "Scaleway: scaleway_function_domain"
---

# Resource: scaleway_function_domain

Creates and manages Scaleway Function Domain bindings.
For more information see [the documentation](https://developers.scaleway.com/en/products/functions/api/).

## Example Usage

### Basic

```terraform
resource "scaleway_function_domain" "main" {
  function_id = scaleway_function.main.id
  hostname    = "example.com"

  depends_on = [
    scaleway_function.main,
  ]
}

resource scaleway_function_namespace main {}

resource scaleway_function main {
  namespace_id = scaleway_function_namespace.main.id
  runtime = "go118"
  privacy = "private"
  handler = "Handle"
  zip_file = "testfixture/gofunction.zip"
  deploy = true
}
```

## Argument Reference

The following arguments are supported:

- `function_id` - (Required) The ID of the function you want to create a domain with.
- `hostname` - (Required) The hostname that should resolve to your function id native domain.
  You should use a CNAME domain record that point to your native function `domain_name` for it.

~> **Important** Updates to `function_id` or `hostname` will recreate the domain.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The function domain's ID.

~> **Important:** Function domains' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`

- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#regions) in where the domain was created.
- `url` - The URL that triggers the function

## Import

Domain can be imported using the `{region}/{id}`, e.g.

```bash
$ terraform import scaleway_function_domain.main fr-par/11111111-1111-1111-1111-111111111111
```
