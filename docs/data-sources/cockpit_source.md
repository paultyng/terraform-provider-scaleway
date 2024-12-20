---
subcategory: "Cockpit"
page_title: "Scaleway: scaleway_cockpit_source"
---

# Data Source: scaleway_cockpit_source

The `scaleway_cockpit_source` data source allows you to retrieve information about a specific [data source](https://www.scaleway.com/en/docs/observability/cockpit/concepts/#data-sources) in Scaleway's Cockpit.

Refer to Cockpit's [product documentation](https://www.scaleway.com/en/docs/observability/cockpit/concepts/) and [API documentation](https://www.scaleway.com/en/developers/api/cockpit/regional-api) for more information.

## Example Usage

### Retrieve a specific data source by ID

The following example retrieves a Cockpit data source by its unique ID.

```terraform
data "scaleway_cockpit_source" "example" {
    id = "fr-par/11111111-1111-1111-1111-111111111111"
}
```

### Retrieve a data source by filters

You can also retrieve a data source by specifying filtering criteria such as `name`, `type`, and `origin`.

```terraform
data "scaleway_cockpit_source" "filtered" {
    project_id = "11111111-1111-1111-1111-111111111111"
    region     = "fr-par"
    name       = "my-data-source"
}
```

## Argument Reference

This section lists the arguments that are supported:

- `id` - (Optional) The unique identifier of the Cockpit data source in the `{region}/{id}` format. If specified, other filters are ignored.

- `region` - (Optional) The [region](../guides/regions_and_zones.md#regions) where the data source is located. Defaults to the region specified in the [provider configuration](../index.md#region).

- `project_id` - (Required unless `id` is specified) The ID of the Project the data source is associated with. Defaults to the Project ID specified in the [provider configuration](../index.md#project_id).

- `name` - (Optional) The name of the data source.

- `type` - (Optional) The [type](https://www.scaleway.com/en/docs/observability/cockpit/concepts/#data-types) of data source. Possible values are: `metrics`, `logs`, or `traces`.

- `origin` - (Optional) The origin of the data source. Possible values are:
    - `scaleway` - Data source managed by Scaleway.
    - `external` - Data source created by the user.
    - `custom` - User-defined custom data source.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier of the data source in the `{region}/{id}` format.

- `url` - The URL of the Cockpit data source.

- `created_at` - The date and time the data source was created (in RFC 3339 format).

- `updated_at` - The date and time the data source was last updated (in RFC 3339 format).

- `origin` - The origin of the data source.

- `synchronized_with_grafana` - Indicates whether the data source is synchronized with Grafana.

- `retention_days` - The number of days the data is retained in the data source.

## Import

You can import a Cockpit data source using its unique ID in the `{region}/{id}` format.

```bash
terraform import scaleway_cockpit_source.example fr-par/11111111-1111-1111-1111-111111111111
```
