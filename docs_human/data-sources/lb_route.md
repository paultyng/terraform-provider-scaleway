---
subcategory: "Load Balancers"
page_title: "Scaleway: scaleway_lb_route"
---

# scaleway_lb_route

Get information about Scaleway Load-Balancer Routes.
For more information, see [the documentation](https://www.scaleway.com/en/developers/api/load-balancer/zoned-api/#path-route).

## Example Usage

```hcl
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

data "scaleway_lb_route" "byID" {
  route_id = scaleway_lb_route.rt01.id
}
```

## Argument Reference

The following argument is supported:

- `route_id` - (Required) The route id.

## Attributes Reference

See the [LB Route Resource](../resources/lb_route.md) for details on the returned attributes - they are identical.