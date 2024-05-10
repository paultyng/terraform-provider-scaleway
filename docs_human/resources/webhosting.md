---
subcategory: "Web Hosting"
page_title: "Scaleway: scaleway_webhosting"
---

# Resource: scaleway_webhosting

Creates and manages Scaleway Web Hostings.
For more information, see [the documentation](https://www.scaleway.com/en/developers/api/webhosting/).

## Example Usage

```terraform
data "scaleway_webhosting_offer" "by_name" {
  name = "lite"
}

resource "scaleway_webhosting" "main" {
  offer_id = data.scaleway_webhosting_offer.by_name.offer_id
  email    = "your@email.com"
  domain   = "yourdomain.com"
  tags     = ["webhosting", "provider", "terraform"]
}
```

## Argument Reference

The following arguments are supported:

- `offer_id` - (Required) The ID of the selected offer for the hosting.
- `email` - (Required) The contact email of the client for the hosting.
- `domain` - (Required) The domain name of the hosting.
- `option_ids` - (Optional) The IDs of the selected options for the hosting.
- `tags` - (Optional) The tags associated with the hosting.
- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#regions) of the Hosting.
- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the VPC is associated with.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the hosting.

~> **Important:** Hostings' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111

- `status` - The hosting status.
- `created_at` - Date and time of hosting's creation (RFC 3339 format).
- `updated_at` - Date and time of hosting's last update (RFC 3339 format).
- `platform_hostname` - The hostname of the host platform.
- `platform_number` - The number of the host platform.
- `offer_name` - The name of the active offer.
- `options` - The active options of the hosting.
    - `id` - The option ID.
    - `name` - The option name.
- `dns_status` - The DNS status of the hosting.
- `cpanel_urls` - The URL to connect to cPanel Dashboard and to Webmail interface.
    - `dashboard` - The URL of the Dashboard.
    - `webmail` - The URL of the Webmail interface.
- `username` - The main hosting cPanel username.
- `organization_id` - The organization ID the hosting is associated with.

## Import

Hostings can be imported using the `{region}/{id}`, e.g.

```bash
$ terraform import scaleway_webhosting.hosting01 fr-par/11111111-1111-1111-1111-111111111111
```