---
page_title: "Scaleway: scaleway_domain_record"
description: |-
  Manages Scaleway Domain records.
---

# scaleway_domain_record

Creates and manages Scaleway Domain record.  
For more information, see [the documentation](https://www.scaleway.com/en/docs/scaleway-dns/).

## Examples

### Basic

```hcl
resource "scaleway_domain_record" "www" {
  dns_zone = "domain.tld"
  name     = "www"
  type     = "A"
  data     = "1.2.3.4"
  ttl      = 3600
}

resource "scaleway_domain_record" "www2" {
  dns_zone = "domain.tld"
  name     = "www"
  type     = "A"
  data     = "1.2.3.5"
  ttl      = 3600
}

resource "scaleway_domain_record" "mx" {
  dns_zone = "domain.tld"
  name     = ""
  type     = "MX"
  data     = "mx.online.net."
  ttl      = 3600
  priority = 10
}

resource "scaleway_domain_record" "mx2" {
  dns_zone = "domain.tld"
  name     = ""
  type     = "MX"
  data     = "mx-cache.online.net."
  ttl      = 3600
  priority = 20
}
```

### With dynamic records

```hcl
resource "scaleway_domain_record" "geo_ip" {
  dns_zone = "domain.tld"
  name     = "images"
  type     = "A"
  data     = "1.2.3.4"
  ttl      = 3600

  geo_ip {
    matches {
      continents = ["EU"]
      countries  = ["FR"]
      data       = "1.2.3.5"
    }
    
    matches {
      continents = ["NA"]
      data       = "4.3.2.1"
    }
  }
}

resource "scaleway_domain_record" "http_service" {
  dns_zone = "domain.tld"
  name     = "app"
  type     = "A"
  data     = "1.2.3.4"
  ttl      = 3600

  http_service {
    ips          = ["1.2.3.5", "1.2.3.6"]
    must_contain = "up"
    url          = "http://mywebsite.com/health"
    user_agent   = "scw_service_up"
    strategy     = "hashed"
  }
}

resource "scaleway_domain_record" "view" {
  dns_zone = "domain.tld"
  name     = "db"
  type     = "A"
  data     = "1.2.3.4"
  ttl      = 3600

  view {
    subnet = "100.0.0.0/16"
    data   = "1.2.3.5"
  }

  view {
    subnet = "100.1.0.0/16"
    data   = "1.2.3.6"
  }
}

resource "scaleway_domain_record" "weighted" {
  dns_zone = "domain.tld"
  name     = "web"
  type     = "A"
  data     = "1.2.3.4"
  ttl      = 3600

  weighted {
    ip     = "1.2.3.5"
    weight = 1
  }

  weighted {
    ip     = "1.2.3.6"
    weight = 2
  }
}
```

### Create an instance and add records with the new instance IP

```hcl
variable "project_id" {
  type        = string
  description = "Your project ID."
}

variable "dns_zone" {
  type        = string
  description = "The DNS Zone used for testing records."
}

resource "scaleway_instance_ip" "public_ip" {
  project_id = var.project_id
}

resource "scaleway_instance_server" "web" {
  project_id = var.project_id
  type       = "DEV1-S"
  image      = "ubuntu_focal"
  tags       = ["front", "web"]
  ip_id      = scaleway_instance_ip.public_ip.id

  root_volume {
    size_in_gb = 20
  }
}

resource "scaleway_domain_record" "web_A" {
  dns_zone = var.dns_zone
  name     = "web"
  type     = "A"
  data     = scaleway_instance_server.web.public_ip
  ttl      = 3600
}

resource "scaleway_domain_record" "web_cname" {
  dns_zone = var.dns_zone
  name     = "www"
  type     = "CNAME"
  data     = "web.${var.dns_zone}."
  ttl      = 3600
}

resource "scaleway_domain_record" "web_alias" {
  dns_zone = var.dns_zone
  name     = ""
  type     = "ALIAS"
  data     = "web.${var.dns_zone}."
  ttl      = 3600
}
```

## Arguments Reference

The following arguments are supported:

- `dns_zone` - (Required) The DNS Zone of the domain. If the DNS zone doesn't exist, it will be automatically created.

- `keep_empty_zone` - (Optional, default: `false`) When destroying a resource, if only NS records remain and this is set to `false`, the zone will be deleted. Please note, each zone not deleted will [cost you money](https://www.scaleway.com/en/dns/)

- `name` - (Required) The name of the record (can be an empty string for a root record).

- `type` - (Required) The type of the record (`A`, `AAAA`, `MX`, `CNAME`, `ALIAS`, `NS`, `PTR`, `SRV`, `TXT`, `TLSA`, or `CAA`).

- `data` - (Required) The content of the record (an IPv4 for an `A`, a string for a `TXT`...).

- `ttl` - (Optional, default: `3600`) Time To Tive of the record in seconds.
  
- `priority` - (Optional, default: `0`) The priority of the record (mostly used with an `MX` record)

**Dynamic records:**

- `geo_ip` - (Optional) The Geo IP feature provides DNS resolution, based on the user’s geographical location. You can define a default IP that resolves if no Geo IP rule matches, and specify IPs for each geographical zone. [Documentation and usage example](https://www.scaleway.com/en/docs/scaleway-dns/#-Geo-IP-Records)
    - `matches` - (Required) The list of matches. *(Can be more than 1)*
        - `countries` - (Optional) List of countries (eg: `FR` for France, `US` for the United States, `GB` for Great Britain...). [List of all countries code](https://api.scaleway.com/domain-private/v2beta1/countries)
        - `continents` - (Optional) List of continents (eg: `EU` for Europe, `NA` for North America, `AS` for Asia...). [List of all continents code](https://api.scaleway.com/domain-private/v2beta1/continents)
        - `data` (Required) The data of the match result


- `http_service` - (Optional) The DNS service checks the provided URL on the configured IPs and resolves the request to one of the IPs by excluding the ones not responding to the given string to check. [Documentation and usage example](https://www.scaleway.com/en/docs/scaleway-dns/#-Healthcheck-records)
    - `ips` - (Required) List of IPs to check
    - `must_contain` - (Required) Text to search
    - `url` - (Required) URL to match the `must_contain` text to validate an IP
    - `user_agent` - (Optional) User-agent used when checking the URL
    - `strategy` - (Required) Strategy to return an IP from the IPs list. Can be `random` or `hashed`


- `view` - (Optional) The answer to a DNS request is based on the client’s (resolver) subnet. *(Can be more than 1)* [Documentation and usage example](https://www.scaleway.com/en/docs/scaleway-dns/#-Views-records)
    - `subnet` - (Required) The subnet of the view
    - `data` - (Required) The data of the view record


- `weighted` - (Optional) You provide a list of IPs with their corresponding weights. These weights are used to proportionally direct requests to each IP. Depending on the weight of a record more or fewer requests are answered with its related IP compared to the others in the list. *(Can be more than 1)* [Documentation and usage example](https://www.scaleway.com/en/docs/scaleway-dns/#-Weight-Records)
    - `ip` - (Required) The weighted IP
    - `weight` - (Required) The weight of the IP as an integer UInt32.

## Multiple records

Some record types can have multiple `data` with the same `name` (eg: `A`, `AAAA`, `MX`, `NS`...).  
You can duplicate a resource `scaleway_domain_record` with the same `name`, the records will be added.

Please note, some record (eg: `CNAME`, Multiple dynamic records of different types...) has to be unique.

## Import

Record can be imported using the `{dns_zone}/{id}`, e.g.

```bash
$ terraform import scaleway_domain_record.www subdomain.domain.tld/11111111-1111-1111-1111-111111111111
```
