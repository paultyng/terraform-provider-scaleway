---
subcategory: "Transactional Email"
page_title: "Scaleway: scaleway_tem_domain"
---

# Resource: scaleway_tem_domain

Creates and manages Scaleway Transactional Email Domains.
For more information see [the documentation](https://developers.scaleway.com/en/products/transactional_email/api/).

## Example Usage

### Basic

```terraform
resource "scaleway_tem_domain" "main" {
  accept_tos = true
  name       = "example.com"
}
```

### Add the required records to your DNS zone

```terraform
variable "domain_name" {
  type    = string
}

resource "scaleway_tem_domain" "main" {
  name       = var.domain_name
  accept_tos = true
}

resource "scaleway_domain_record" "spf" {
  dns_zone = var.domain_name
  type     = "TXT"
  data     = "v=spf1 ${scaleway_tem_domain.main.spf_config} -all"
}

resource "scaleway_domain_record" "dkim" {
  dns_zone = var.domain_name
  name     = "${scaleway_tem_domain.main.project_id}._domainkey"
  type     = "TXT"
  data     = scaleway_tem_domain.main.dkim_config
}

resource "scaleway_domain_record" "mx" {
  dns_zone = var.domain_name
  type     = "MX"
  data     = "."
}
```


### Configuring GitLab Project Variables

```terraform

variable "domain_name" {
  type    = string
}

data "scaleway_tem_domain" "my_domain" {
  name       = var.domain_name
}

resource "gitlab_project_variable" "smtp_host" {
  key   = "SMTP_AUTH_USER"
  value = data.scaleway_tem_domain.my_domain.smtps_auth_user
}

resource "gitlab_project_variable" "smtp_port" {
  key   = "SMTP_PORT"
  value = data.scaleway_tem_domain.my_domain.smtps_port
}

resource "gitlab_project_variable" "smtp_host" {
  key   = "SMTP_HOST"
  value = data.scaleway_tem_domain.my_domain.smtps_host
}

```


## Argument Reference

The following arguments are supported:

- `name` - (Required) The domain name, must not be used in another Transactional Email Domain.
~> **Important:** Updates to `name` will recreate the domain.

- `accept_tos` - (Required) Acceptation of the [Term of Service](https://tem.s3.fr-par.scw.cloud/antispam_policy.pdf).
~> **Important:**  This attribute must be set to `true`.

- `region` - (Defaults to [provider](../index.md#region) `region`). The [region](../guides/regions_and_zones.md#regions) in which the domain should be created.

- `project_id` - (Defaults to [provider](../index.md#project_id) `project_id`) The ID of the project the domain is associated with.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the Transaction Email Domain.

~> **Important:** Transaction Email Domains' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`

- `status` - The status of the Transaction Email Domain.

- `created_at` - The date and time of the Transaction Email Domain's creation (RFC 3339 format).

- `next_check_at` - The date and time of the next scheduled check (RFC 3339 format).

- `last_valid_at` - The date and time the domain was last found to be valid (RFC 3339 format).

- `revoked_at` - The date and time of the revocation of the domain (RFC 3339 format).

- `last_error` - The error message if the last check failed.

- `spf_config` - The snippet of the SPF record that should be registered in the DNS zone.

- `dkim_config` - The DKIM public key, as should be recorded in the DNS zone.

- `smtp_host` - The SMTP host to use to send emails.

- `smtp_port_unsecure` - The SMTP port to use to send emails.

- `smtp_port` - The SMTP port to use to send emails over TLS.

- `smtp_port_alternative` - The SMTP port to use to send emails over TLS.

- `smtps_port` - The SMTPS port to use to send emails over TLS Wrapper.

- `smtps_port_alternative` - The SMTPS port to use to send emails over TLS Wrapper.

- `smtps_auth_user` - SMTPS auth user refers to the identifier for a user authorized to send emails via SMTPS, ensuring secure email transmission.

- `mx_blackhole` - The Scaleway's blackhole MX server to use if you do not have one.

- `reputation` - The domain's reputation.
    - `status` - The status of the domain's reputation.
    - `score` - A range from 0 to 100 that determines your domain's reputation score.
    - `scored_at` - The time and date the score was calculated.
    - `previous_score` - The previously-calculated domain's reputation score.
    - `previous_scored_at` - The time and date the previous reputation score was calculated.

## Import

Domains can be imported using the `{region}/{id}`, e.g.

```bash
$ terraform import scaleway_tem_domain.main fr-par/11111111-1111-1111-1111-111111111111
```
