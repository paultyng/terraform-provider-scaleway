---
subcategory: "Load Balancers"
page_title: "Scaleway: scaleway_lb_acl"
---

# Resource: scaleway_lb_acl

Creates and manages Scaleway Load-Balancer ACLs. For more information, see [the documentation](https://www.scaleway.com/en/developers/api/load-balancer/zoned-api/#path-acls).

## Example Usage

### Basic

```terraform
resource "scaleway_lb_acl" "acl01" {
  frontend_id  = scaleway_lb_frontend.frt01.id
  name         = "acl01"
  description  = "Exclude well-known IPs"
  index        = 0
  # Allow downstream requests from: 192.168.0.1, 192.168.0.2 or 192.168.10.0/24
  action {
    type = "allow"
  }
  match {
    ip_subnet = ["192.168.0.1", "192.168.0.2", "192.168.10.0/24"]
  }
}
```

## Argument Reference

The following arguments are supported:

- `frontend_id` - (Required) The load-balancer Frontend ID to attach the ACL to.

- `name` - (Optional) The ACL name. If not provided it will be randomly generated.

- `description` - (Optional) The ACL description.

- `index` - (Required) The Priority of this ACL (ACLs are applied in ascending order, 0 is the first ACL executed).

- `action` - (Required) Action to undertake when an ACL filter matches.

    - `type` - (Required) The action type. Possible values are: `allow` or `deny` or `redirect`.

    - `redirect` - (Optional) Redirect parameters when using an ACL with `redirect` action.

        - `type`  - (Optional) The redirect type. Possible values are: `location` or `scheme`.

        - `target`  - (Optional) An URL can be used in case of a location redirect (e.g. `https://scaleway.com` will redirect to this same URL). A scheme name (e.g. `https`, `http`, `ftp`, `git`) will replace the request's original scheme.

        - `code`  - (Optional) The HTTP redirect code to use. Valid values are `301`, `302`, `303`, `307` and `308`.

- `match` - (Required) The ACL match rule. At least `ip_subnet` or `http_filter` and `http_filter_value` are required.

    - `ip_subnet` - (Optional) A list of IPs or CIDR v4/v6 addresses of the client of the session to match.

    - `http_filter` - (Optional) The HTTP filter to match. This filter is supported only if your backend protocol has an HTTP forward protocol.
      It extracts the request's URL path, which starts at the first slash and ends before the question mark (without the host part).
      Possible values are: `acl_http_filter_none`, `path_begin`, `path_end`, `http_header_match` or `regex`.

    - `http_filter_value` - (Optional) A list of possible values to match for the given HTTP filter.
      Keep in mind that in the case of `http_header_match` the HTTP header field name is case-insensitive.

    - `http_filter_option` - (Optional) If you have `http_filter` at `http_header_match`, you can use this field to filter on the HTTP header's value.

    - `invert` - (Optional) If set to `true`, the condition will be of type "unless".

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the load-balancer ACL.

~> **Important:** Load-Balancers ACLs' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`


## Import

Load-Balancer ACL can be imported using the `{zone}/{id}`, e.g.

```bash
$ terraform import scaleway_lb_acl.acl01 fr-par-1/11111111-1111-1111-1111-111111111111
```
