---
subcategory: "Load Balancers"
page_title: "Scaleway: scaleway_lb_backend"
---

# scaleway_lb_backend

Get information about Scaleway Load-Balancer Backends.
For more information, see [the documentation](https://www.scaleway.com/en/developers/api/load-balancer/zoned-api/#path-backends).

## Example Usage

```hcl
resource "scaleway_lb_ip" "main" {
}

resource "scaleway_lb" "main" {
  ip_id = scaleway_lb_ip.main.id
  name  = "data-test-lb-backend"
  type  = "LB-S"
}

resource "scaleway_lb_backend" "main" {
  lb_id            = scaleway_lb.main.id
  name             = "backend01"
  forward_protocol = "http"
  forward_port     = "80"
}

data "scaleway_lb_backend" "byID" {
  backend_id = scaleway_lb_backend.main.id
}

data "scaleway_lb_backend" "byName" {
  name  = scaleway_lb_backend.main.name
  lb_id = scaleway_lb.main.id
}
```

## Arguments Reference

The following arguments are supported:

- `backend_id` - (Optional) The backend id.
    - Only one of `name` and `backend_id` should be specified.

- `name` - (Optional) The name of the backend.
    - When using the `name` you should specify the `lb-id`

- `lb_id` - (Required) The load-balancer ID this backend is attached to.

## Attributes Reference

See the [LB Backend Resource](../resources/lb_backend.md) for details on the returned attributes - they are identical.