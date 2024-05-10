---
subcategory: "Load Balancers"
page_title: "Scaleway: scaleway_lb_certificate"
---

# Resource: scaleway_lb_certificate

Creates and manages Scaleway Load-Balancer Certificates.
For more information, see [the documentation](https://www.scaleway.com/en/developers/api/load-balancer/zoned-api/#path-certificate).

## Example Usage

### Let's Encrypt

```terraform
resource "scaleway_lb_certificate" "cert01" {
  lb_id = scaleway_lb.lb01.id
  name  = "cert1"

  letsencrypt {
    common_name = "example.org"
    subject_alternative_name = [
      "sub1.example.com",
      "sub2.example.com"
    ]
  }
  # Make sure the new certificate is created before the old one can be replaced
  lifecycle {
      create_before_destroy = true
  }
}
```

### Custom Certificate

```terraform
resource "scaleway_lb_certificate" "cert01" {
  lb_id = scaleway_lb.lb01.id
  name  = "custom-cert"
  custom_certificate {
    certificate_chain = <<EOF
CERTIFICATE_CHAIN_CONTENTS
EOF
  }
}
```

## Argument Reference

The following arguments are supported:

### Basic arguments

- `lb_id` - (Required) The load-balancer ID this certificate is attached to.

~> **Important:** Updates to `lb_id` will recreate the load-balancer certificate.

- `name` - (Optional) The name of the certificate backend.

- `letsencrypt` - (Optional) Configuration block for Let's Encrypt configuration. Only one of `letsencrypt` and `custom_certificate` should be specified.

    - `common_name` - (Required) Main domain of the certificate. A new certificate will be created if this field is changed.

    - `subject_alternative_name` - (Optional) Array of alternative domain names.  A new certificate will be created if this field is changed.

~> **Important:** Updates to `letsencrypt` will recreate the load-balancer certificate.

- `custom_certificate` - (Optional) Configuration block for custom certificate chain. Only one of `letsencrypt` and `custom_certificate` should be specified.

    - `certificate_chain` - (Required) Full PEM-formatted certificate chain.

~> **Important:** Updates to `custom_certificate` will recreate the load-balancer certificate.

- `zone` - (Defaults to [provider](../index.md#zone) `zone`) The [zone](../guides/regions_and_zones.md#zones) of the certificate.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the load-balancer certificate.

~> **Important:** Load-Balancers certificates' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

- `common_name` - Main domain of the certificate
- `subject_alternative_name` - The alternative domain names of the certificate
- `fingerprint` - The identifier (SHA-1) of the certificate
- `not_valid_before` - The not valid before validity bound timestamp
- `not_valid_after` - The not valid after validity bound timestamp
- `status` - Certificate status

## Additional notes

* Ensure that all domain names used in configuration are pointing to the load balancer IP.
  You can achieve this by creating a DNS record through terraform pointing to  `ip_address` property of `lb_beta` entity.
* In case there are any issues with the certificate, you will receive a `400` error from the `apply` operation.
  Use `export TF_LOG=DEBUG` to view exact problem returned by the api.
* Wildcards are not supported with Let's Encrypt yet.
* Use `lifecycle` instruction with `create_before_destroy = true` to permit correct certificate replacement and prevent a `400` error from the `apply` operation.
