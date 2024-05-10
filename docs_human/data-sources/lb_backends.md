---
subcategory: "Load Balancers"
page_title: "Scaleway: scaleway_lb_backends"
---

# scaleway_lb_backends

Gets information about multiple Load Balancer Backends.

## Example Usage

```hcl
# Find backends that share the same LB ID
data "scaleway_lb_backends" "byLBID" {
  lb_id = "${scaleway_lb.lb01.id}"
}
# Find backends by LB ID and name
data "scaleway_lb_backends" "byLBID_and_name" {
  lb_id = "${scaleway_lb.lb01.id}"
  name  = "tf-backend-datasource"
}
```

## Argument Reference

- `lb_id` - (Required) The load-balancer ID this backend is attached to. backends with a LB ID like it are listed.

- `name` - (Optional) The backend name used as filter. Backends with a name like it are listed.

- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which backends exist.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `backends` - List of found backends
    - `id` - The associated backend ID.
    - `forward_protocol` - Backend protocol.
    - `created_at` - The date at which the backend was created (RFC 3339 format).
    - `update_at` - The date at which the backend was last updated (RFC 3339 format).
    - `forward_port` - User sessions will be forwarded to this port of backend servers.
    - `forward_port_algorithm` - Load balancing algorithm.
    - `sticky_sessions` - Enables cookie-based session persistence.
    - `sticky_sessions_cookie_name` - Cookie name for sticky sessions.
    - `server_ips` - List of backend server IP addresses.
    - `proxy_protocol` - The type of PROXY protocol.
    - `timeout_server` - Maximum server connection inactivity time.
    - `timeout_connect` - Maximum initial server connection establishment time.
    - `timeout_tunnel` - Maximum tunnel inactivity time.
    - `failover_host` - Scaleway S3 bucket website to be served in case all backend servers are down.
    - `ssl_bridging` - Enables SSL between load balancer and backend servers.
    - `ignore_ssl_server_verify` - Specifies whether the Load Balancer should check the backend server’s certificate before initiating a connection.
    - `health_check_timeout` - Timeout before we consider a HC request failed.
    - `health_check_delay` - Interval between two HC requests.
    - `health_check_port` - Port the HC requests will be sent to.
    - `health_check_max_retries` - Number of allowed failed HC requests before the backend server is marked down.
    - `health_check_tcp` - This block enable TCP health check.
    - `health_check_http` - This block enable HTTP health check.
        - `uri` - The HTTP endpoint URL to call for HC requests.
        - `method` - The HTTP method to use for HC requests.
        - `code` - The expected HTTP status code.
        - `host_header` -  The HTTP host header to use for HC requests.
    - `health_check_https` - This block enable HTTPS health check.
        - `uri` - The HTTPS endpoint URL to call for HC requests.
        - `method` - The HTTP method to use for HC requests.
        - `code` - The expected HTTP status code.
        - `host_header` - The HTTP host header to use for HC requests.
        - `sni` - The SNI to use for HC requests over SSL.
    - `on_marked_down_action` - Modify what occurs when a backend server is marked down.
