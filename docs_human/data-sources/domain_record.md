---
subcategory: "Domains and DNS"
page_title: "Scaleway: scaleway_domain_record"
---

# scaleway_domain_record

Gets information about a domain record.

## Example Usage

```hcl
# Get record by name, type and data
data "scaleway_domain_record" "by_content" {
  dns_zone = "domain.tld"
  name     = "www"
  type     = "A"
  data     = "1.2.3.4"
}

# Get info by ID
data "scaleway_domain_record" "by_id" {
  dns_zone  = "domain.tld"
  record_id = "11111111-1111-1111-1111-111111111111"
}
```

## Argument Reference

- `dns_zone` - (Optional) The IP address.

- `name` - (Required) The name of the record (can be an empty string for a root record).
  Cannot be used with `record_id`.

- `type` - (Required) The type of the record (`A`, `AAAA`, `MX`, `CNAME`, `DNAME`, `ALIAS`, `NS`, `PTR`, `SRV`, `TXT`, `TLSA`, or `CAA`).
  Cannot be used with `record_id`.

- `data` - (Required) The content of the record (an IPv4 for an `A`, a string for a `TXT`...).
  Cannot be used with `record_id`.

- `record_id` - (Optional) The record ID.
  Cannot be used with `name`, `type` and `data`.

- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the domain is associated with.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the record.

~> **Important:** Domain records' IDs are of the form `{dns_zone}/{id}`, e.g. `subdomain.domain.tld/11111111-1111-1111-1111-111111111111`

- `ttl` - Time To Live of the record in seconds.
- `priority` - The priority of the record (mostly used with an `MX` record)
- `geo_ip` - Dynamic record base on user geolocalisation ([More information about dynamic records](../resources/domain_record.md#dynamic-records))
- `http_service` - Dynamic record base on URL resolve ([More information about dynamic records](../resources/domain_record.md#dynamic-records))
- `weighted` - Dynamic record base on IP weights ([More information about dynamic records](../resources/domain_record.md#dynamic-records))
- `view` - Dynamic record based on the client’s (resolver) subnet ([More information about dynamic records](../resources/domain_record.md#dynamic-records))
