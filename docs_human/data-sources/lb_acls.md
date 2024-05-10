---
subcategory: "Load Balancers"
page_title: "Scaleway: scaleway_lb_acls"
---

# scaleway_lb_acls

Gets information about multiple Load Balancer ACLs.

## Example Usage

```hcl
# Find acls that share the same frontend ID
data "scaleway_lb_acls" "byFrontID" {
  frontend_id = "${scaleway_lb_frontend.frt01.id}"
}
# Find acls by frontend ID and name
data "scaleway_lb_acls" "byFrontID_and_name" {
  frontend_id = "${scaleway_lb_frontend.frt01.id}"
  name        = "tf-acls-datasource"
}
```

## Argument Reference

- `frontend_id` - (Required) The frontend ID this ACL is attached to. ACLs with a frontend ID like it are listed.
  ~> **Important:** LB Frontends' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

- `name` - (Optional) The ACL name used as filter. ACLs with a name like it are listed.

- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) in which ACLs exist.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `acls` - List of found ACLs
    - `id` - The associated ACL ID.
      ~> **Important:** LB ACLs' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`
    - `created_at` - The date at which the ACL was created (RFC 3339 format).
    - `update_at` - The date at which the ACL was last updated (RFC 3339 format).
    - `index` - The order between the ACLs.
    - `description` - The description of the ACL resource.
    - `action` - The action that has been undertaken when an ACL filter had matched.
        - `type` - The action type.
        - `redirect` - Redirect parameters when using an ACL with `redirect` action.
            - `type`  - The redirect type.
            - `target`  - The URL used in case of a location redirect or the scheme name that replaces the request's original scheme.
            - `code`  - The HTTP redirect code used.
    - `match` - The ACL match rule.
        - `ip_subnet` - A list of matched IPs or CIDR v4/v6 addresses of the client of the session.
        - `http_filter` - The matched HTTP filter.
        - `http_filter_value` - The possible values matched for a given HTTP filter.
        - `http_filter_option` - A list of possible values for the HTTP filter based on the HTTP header.
        - `invert` -  The condition will be of type "unless" if invert is set to `true`