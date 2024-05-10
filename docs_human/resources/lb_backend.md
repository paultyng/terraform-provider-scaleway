---
subcategory: "Load Balancers"
page_title: "Scaleway: scaleway_lb_backend"
---

# Resource: scaleway_lb_backend

Creates and manages Scaleway Load-Balancer Backends.
For more information, see [the documentation](https://www.scaleway.com/en/developers/api/load-balancer/zoned-api/#path-backends).

## Example Usage

### Basic

```terraform
resource "scaleway_lb_backend" "backend01" {
  lb_id            = scaleway_lb.lb01.id
  name             = "backend01"
  forward_protocol = "http"
  forward_port     = "80"
}
```

### With HTTP Health Check

```terraform
resource "scaleway_lb_backend" "backend01" {
  lb_id            = scaleway_lb.lb01.id
  name             = "backend01"
  forward_protocol = "http"
  forward_port     = "80"

  health_check_http {
    uri = "www.test.com/health"
  }
}
```

## Argument Reference

The following arguments are supported:

### Basic arguments

- `lb_id`                       - (Required) The load-balancer ID this backend is attached to.
~> **Important:** Updates to `lb_id` will recreate the backend.
- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which the load-balancer was created.
- `forward_protocol`            - (Required) Backend protocol. Possible values are: `tcp` or `http`.
- `name`                        - (Optional) The name of the load-balancer backend.
- `forward_port`                - (Required) User sessions will be forwarded to this port of backend servers.
- `forward_port_algorithm`      - (Default: `roundrobin`) Load balancing algorithm. Possible values are: `roundrobin`, `leastconn` and `first`.
- `sticky_sessions`             - (Default: `none`) The type of sticky sessions. The only current possible values are: `none`, `cookie` and `table`.
- `sticky_sessions_cookie_name` - (Optional) Cookie name for sticky sessions. Only applicable when sticky_sessions is set to `cookie`.
- `server_ips`                  - (Optional) List of backend server IP addresses. Addresses can be either IPv4 or IPv6.
- `send_proxy_v2`               - DEPRECATED please use `proxy_protocol` instead - (Default: `false`) Enables PROXY protocol version 2.
- `proxy_protocol`              - (Default: `none`) Choose the type of PROXY protocol to enable (`none`, `v1`, `v2`, `v2_ssl`, `v2_ssl_cn`)
- `timeout_server`              - (Optional) Maximum server connection inactivity time. (e.g.: `1s`)
- `timeout_connect`             - (Optional) Maximum initial server connection establishment time. (e.g.: `1s`)
- `timeout_tunnel`              - (Optional) Maximum tunnel inactivity time. (e.g.: `1s`)
- `failover_host`               - (Optional) Scaleway S3 bucket website to be served in case all backend servers are down.
~> **Note:** Only the host part of the Scaleway S3 bucket website is expected:
e.g. 'failover-website.s3-website.fr-par.scw.cloud' if your bucket website URL is 'https://failover-website.s3-website.fr-par.scw.cloud/'.
- `ssl_bridging`                - (Default: `false`) Enables SSL between load balancer and backend servers.
- `ignore_ssl_server_verify`    - (Default: `false`) Specifies whether the Load Balancer should check the backend server’s certificate before initiating a connection.
- `max_connections`             - (Optional) Maximum number of connections allowed per backend server.
- `timeout_queue`               - (Optional) Maximum time for a request to be left pending in queue when `max_connections` is reached. (e.g.: `1s`)
- `redispatch_attempt_count`    - (Optional) Whether to use another backend server on each attempt.
- `max_retries`                 - (Optional) Number of retries when a backend server connection failed.

### Health Check arguments

Backends use Health Check to test if a backend server is ready to receive requests.
You may use one of the following health check types: `TCP`, `HTTP` or `HTTPS`. (Default: `TCP`)

- `health_check_timeout`          - (Default: `30s`) Timeout before we consider a HC request failed.
- `health_check_delay`            - (Default: `60s`) Interval between two HC requests.
- `health_check_port`             - (Default: `forward_port`) Port the HC requests will be send to.
- `health_check_max_retries`      - (Default: `2`) Number of allowed failed HC requests before the backend server is marked down.
- `health_check_tcp`              - (Optional) This block enable TCP health check. Only one of `health_check_tcp`, `health_check_http` and `health_check_https` should be specified.
- `health_check_http`             - (Optional) This block enable HTTP health check. Only one of `health_check_tcp`, `health_check_http` and `health_check_https` should be specified.
    - `uri`                         - (Required) The HTTP endpoint URL to call for HC requests.
    - `method`                      - (Default: `GET`) The HTTP method to use for HC requests.
    - `code`                        - (Default: `200`) The expected HTTP status code.
    - `host_header`                 - (Optional) The HTTP host header to use for HC requests.
- `health_check_https`            - (Optional) This block enable HTTPS health check. Only one of `health_check_tcp`, `health_check_http` and `health_check_https` should be specified.
    - `uri`                         - (Required) The HTTPS endpoint URL to call for HC requests.
    - `method`                      - (Default: `GET`) The HTTP method to use for HC requests.
    - `code`                        - (Default: `200`) The expected HTTP status code.
    - `host_header`                 - (Optional) The HTTP host header to use for HC requests.
    - `sni`                         - (Optional) The SNI to use for HC requests over SSL.
- `on_marked_down_action`         - (Default: `none`) Modify what occurs when a backend server is marked down. Possible values are: `none` and `shutdown_sessions`.
- `health_check_transient_delay`  - (Default: `0.5s`) The time to wait between two consecutive health checks when a backend server is in a transient state (going UP or DOWN).
- `health_check_send_proxy`       - (Default: `false`) Defines whether proxy protocol should be activated for the health check.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the loadbalancer backend.

~> **Important:** Load-Balancers backends' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

## Import

Load-Balancer backend can be imported using the `{zone}/{id}`, e.g.

```bash
$ terraform import scaleway_lb_backend.backend01 fr-par-1/11111111-1111-1111-1111-111111111111
```
