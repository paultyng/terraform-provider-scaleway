---
page_title: "Scaleway: scaleway_lb_frontend"
description: |-
  Manages Scaleway Load-Balancer Frontends.
---

# scaleway_lb_frontend

Creates and manages Scaleway Load-Balancer Frontends. For more information, see [the documentation](https://developers.scaleway.com/en/products/lb/zoned_api).

## Examples Usage

### Basic

```hcl
resource "scaleway_lb_frontend" "frontend01" {
  lb_id        = scaleway_lb.lb01.id
  backend_id   = scaleway_lb_backend.backend01.id
  name         = "frontend01"
  inbound_port = "80"
}
```

## With Certificate

```hcl
resource scaleway_lb_ip ip01 {}

resource scaleway_lb lb01 {
    ip_id = scaleway_lb_ip.ip01.id
    name = "test-lb"
    type = "lb-s"
}

resource scaleway_lb_backend bkd01 {
    lb_id = scaleway_lb.lb01.id
    forward_protocol = "tcp"
    forward_port = 443
    proxy_protocol = "none"
}

resource scaleway_lb_certificate cert01 {
    lb_id = scaleway_lb.lb01.id
    name = "test-cert-front-end"
    letsencrypt {
        common_name = "${replace(scaleway_lb_ip.ip01.ip_address,".", "-")}.lb.${scaleway_lb.lb01.region}.scw.cloud"
    }
    # Make sure the new certificate is created before the old one can be replaced
    lifecycle {
        create_before_destroy = true
    }
}

resource scaleway_lb_frontend frt01 {
    lb_id = scaleway_lb.lb01.id
    backend_id = scaleway_lb_backend.bkd01.id
    inbound_port = 443
    certificate_ids = [scaleway_lb_certificate.cert01.id]
}
```

## With ACLs

```hcl
resource "scaleway_lb_frontend" "frontend01" {
  lb_id        = scaleway_lb.lb01.id
  backend_id   = scaleway_lb_backend.backend01.id
  name         = "frontend01"
  inbound_port = "80"

  # Allow downstream requests from: 192.168.0.1, 192.168.0.2 or 192.168.10.0/24
  acl {
    name = "blacklist wellknwon IPs"
    action {
      type = "allow"
    }
    match {
      ip_subnet = ["192.168.0.1", "192.168.0.2", "192.168.10.0/24"]
    }
  }

  # Deny downstream requests from: 51.51.51.51 that match "^foo*bar$"
  acl {
    action {
      type = "deny"
    }
    match {
      ip_subnet         = ["51.51.51.51"]
      http_filter       = "regex"
      http_filter_value = ["^foo*bar$"]
    }
  }

  # Allow downstream http requests that begins with "/foo" or "/bar"
  acl {
    action {
      type = "allow"
    }
    match {
      http_filter       = "path_begin"
      http_filter_value = ["foo", "bar"]
    }
  }

  # Allow upstream http requests that DO NOT begins with "/hi"
  acl {
    action {
      type = "allow"
    }
    match {
      http_filter       = "path_begin"
      http_filter_value = ["hi"]
      invert            = "true"
    }
  }

  # Allow upstream http requests that have an HTTP header "foo" that matches "bar"
  acl {
    action {
      type = "allow"
    }

    match {
      http_filter       = "http_header_match"
      http_filter_value = "foo"
      http_value_option = "bar"
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

- `lb_id` - (Required) The load-balancer ID this frontend is attached to.

- `backend_id` - (Required) The load-balancer backend ID this frontend is attached to.

~> **Important:** Updates to `lb_id` or `backend_id` will recreate the frontend.

- `inbound_port` - (Required) TCP port to listen on the front side.

- `name` - (Optional) The name of the load-balancer frontend.

- `timeout_client` - (Optional) Maximum inactivity time on the client side. (e.g.: `1s`)

- `certificate_ids` - (Optional) List of Certificate IDs that should be used by the frontend.

~> **Important:** Certificates are not allowed on port 80.

- `acl` - (Optional) A list of ACL rules to apply to the load-balancer frontend.  Defined below.

## acl

- `name` - (Optional) The ACL name. If not provided it will be randomly generated.
  
- `action` - (Required) Action to undertake when an ACL filter matches.
  
    - `type` - (Required) The action type. Possible values are: `allow` or `deny`.
  
- `match` - (Required) The ACL match rule. At least `ip_subnet` or `http_filter` and `http_filter_value` are required.

    - `ip_subnet` - (Optional) A list of IPs or CIDR v4/v6 addresses of the client of the session to match.

    - `http_filter` - (Optional) The HTTP filter to match. This filter is supported only if your backend protocol has an HTTP forward protocol.
       It extracts the request's URL path, which starts at the first slash and ends before the question mark (without the host part).
       Possible values are: `acl_http_filter_none`, `path_begin`, `path_end`, `http_header_match` or `regex`.

    - `http_filter_value` - (Optional) A list of possible values to match for the given HTTP filter.
      Keep in mind that in the case of `http_header_match` the HTTP header field name is case-insensitive.

    - `http_value_option` - (Optional) If you have `http_filter` at `http_header_match`, you can use this field to filter on the HTTP header's value.

    - `invert` - (Optional) If set to `true`, the condition will be of type "unless".

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the load-balancer frontend.
- `certificate_id` - (Deprecated) first certificate ID used by the frontend.


## Import

Load-Balancer frontend can be imported using the `{zone}/{id}`, e.g.

```bash
$ terraform import scaleway_lb_frontend.frontend01 fr-par-1/11111111-1111-1111-1111-111111111111
```
