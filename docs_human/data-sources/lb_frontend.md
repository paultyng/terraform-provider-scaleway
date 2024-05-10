---
subcategory: "Load Balancers"
page_title: "Scaleway: scaleway_lb_frontend"
---

# scaleway_lb_frontend

Get information about Scaleway Load-Balancer Frontends.
For more information, see [the documentation](https://www.scaleway.com/en/developers/api/load-balancer/zoned-api/#path-frontends).

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

data "scaleway_lb_frontend" "byID" {
  frontend_id = scaleway_lb_frontend.frt01.id
}

data "scaleway_lb_frontend" "byName" {
  name  = scaleway_lb_frontend.frt01.name
  lb_id = scaleway_lb.lb01.id
}
```

## Arguments Reference

The following arguments are supported:

- `frontend_id` - (Optional) The frontend id.
    - Only one of `name` and `frontend_id` should be specified.

- `name` - (Optional) The name of the frontend.
    - When using the `name` you should specify the `lb-id`

- `lb_id` - (Required) The load-balancer ID this frontend is attached to.

## Attributes Reference

See the [LB Frontend Resource](../resources/lb_frontend.md) for details on the returned attributes - they are identical.