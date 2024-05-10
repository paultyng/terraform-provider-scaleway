---
subcategory: "Load Balancers"
page_title: "Scaleway: scaleway_lb_route"
---

# Resource: scaleway_lb_route

Creates and manages Scaleway Load-Balancer Routes.
For more information, see [the documentation](https://www.scaleway.com/en/developers/api/load-balancer/zoned-api/#path-route).

## Example Usage

### With SNI for direction to TCP backends

```terraform
resource "scaleway_lb_ip" "ip01" {}

resource "scaleway_lb" "lb01" {
  ip_id = scaleway_lb_ip.ip01.id
  name  = "test-lb"
  type  = "lb-s"
}

resource "scaleway_lb_backend" "bkd01" {
  lb_id            = scaleway_lb.lb01.id
  forward_protocol = "tcp"
  forward_port     = 80
  proxy_protocol   = "none"
}

resource "scaleway_lb_frontend" "frt01" {
  lb_id        = scaleway_lb.lb01.id
  backend_id   = scaleway_lb_backend.bkd01.id
  inbound_port = 80
}

resource "scaleway_lb_route" "rt01" {
  frontend_id = scaleway_lb_frontend.frt01.id
  backend_id  = scaleway_lb_backend.bkd01.id
  match_sni   = "sni.scaleway.com"
}
```

### With host-header for direction to HTTP backends

```terraform
resource "scaleway_lb_ip" "ip01" {}

resource "scaleway_lb" "lb01" {
  ip_id = scaleway_lb_ip.ip01.id
  name  = "test-lb"
  type  = "lb-s"
}

resource "scaleway_lb_backend" "bkd01" {
  lb_id            = scaleway_lb.lb01.id
  forward_protocol = "http"
  forward_port     = 80
  proxy_protocol   = "none"
}

resource "scaleway_lb_frontend" "frt01" {
  lb_id        = scaleway_lb.lb01.id
  backend_id   = scaleway_lb_backend.bkd01.id
  inbound_port = 80
}

resource "scaleway_lb_route" "rt01" {
  frontend_id       = scaleway_lb_frontend.frt01.id
  backend_id        = scaleway_lb_backend.bkd01.id
  match_host_header = "host.scaleway.com"
}
```

## Argument Reference

The following arguments are supported:

- `backend_id` - (Required) The ID of the backend to which the route is associated.
- `frontend_id` - (Required) The ID of the frontend to which the route is associated.
- `match_sni` - The Server Name Indication (SNI) value to match. Value to match in the Server Name Indication TLS extension (SNI) field from an incoming connection made via an SSL/TLS transport layer.
  Only one of `match_sni` and `match_host_header` should be specified.

~> **Important:** This field should be set for routes on TCP Load Balancers.

- `match_host_header` - The HTTP host header to match. Value to match in the HTTP Host request header from an incoming connection.
  Only one of `match_sni` and `match_host_header` should be specified.  

~> **Important:** This field should be set for routes on HTTP Load Balancers.

- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the load-balancer was created.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the route

~> **Important:** Load-balancers routes' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

- `created_at` - The date at which the route was created.
- `updated_at` - The date at which the route was last updated.

## Import

Load-Balancer frontend can be imported using the `{zone}/{id}`, e.g.

```bash
$ terraform import scaleway_lb_route.main fr-par-1/11111111-1111-1111-1111-111111111111
```
