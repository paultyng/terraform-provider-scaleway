---
page_title: "Provider: Scaleway"
description: |-
  The Scaleway provider is used to manage Scaleway resources. The provider needs to be configured with the proper credentials before it can be used.
---

# Scaleway Provider

The Scaleway provider is used to manage Scaleway resources.
The provider needs to be configured with the proper credentials before it can be used.

**This is the documentation for the version `>= 1.11.0` of the provider. If you come from `< v1.11.0`, checkout to [migration guide](./guides/migration_guide_v2.md).**

Use the navigation to the left to read about the available resources.

## Terraform 0.13 and later

For Terraform 0.13 and later, please also include this:

```hcl
terraform {
  required_providers {
    scaleway = {
      source = "scaleway/scaleway"
    }
  }
  required_version = ">= 0.13"
}
```

## Example

Here is an example that will set up a web server with an additional volume, a public IP and a security group.

You can test this config by creating a `test.tf` and run terraform commands from this directory:

- Get your [Scaleway credentials](https://console.scaleway.com/account/credentials)
- Initialize a Terraform working directory: `terraform init`
- Generate and show the execution plan: `terraform plan`
- Build the infrastructure: `terraform apply`

```hcl
variable "project_id" {
  type        = string
  description = "Your project ID."
}

variable "dns_zone" {
  type        = string
  description = "The DNS Zone used for testing records."
}

terraform {
  required_providers {
    scaleway = {
      source = "scaleway/scaleway"
    }
  }
  required_version = ">= 0.13"
}

provider "scaleway" {
  zone   = "fr-par-1"
  region = "fr-par"
}

resource "scaleway_instance_ip" "public_ip" {
  project_id = var.project_id
}
resource "scaleway_instance_ip" "public_ip_backup" {
  project_id = var.project_id
}

resource "scaleway_instance_volume" "data" {
  project_id = var.project_id
  size_in_gb = 30
  type       = "l_ssd"
}

resource "scaleway_instance_volume" "data_backup" {
  project_id = var.project_id
  size_in_gb = 10
  type       = "l_ssd"
}

resource "scaleway_instance_security_group" "www" {
  project_id              = var.project_id
  inbound_default_policy  = "drop"
  outbound_default_policy = "accept"

  inbound_rule {
    action = "accept"
    port   = "22"
    ip     = "212.47.225.64"
  }

  inbound_rule {
    action = "accept"
    port   = "80"
  }

  inbound_rule {
    action = "accept"
    port   = "443"
  }
}

resource "scaleway_instance_server" "web" {
  project_id = var.project_id
  type       = "DEV1-L"
  image      = "ubuntu_focal"

  tags = ["front", "web"]

  ip_id = scaleway_instance_ip.public_ip.id

  additional_volume_ids = [scaleway_instance_volume.data.id]

  root_volume {
    # The local storage of a DEV1-L instance is 80 GB, subtract 30 GB from the additional l_ssd volume, then the root volume needs to be 50 GB.
    size_in_gb = 50
  }

  security_group_id = scaleway_instance_security_group.www.id
}
```

## Authentication

The Scaleway authentication is based on an **access key**, and a **secret key**.
Since secret keys are only revealed one time (when it is first created) you might
need to create a new one in the section "API Keys" of the [Scaleway console](https://console.scaleway.com/project/credentials).
Click on the "Generate new API key" button to create them.
Giving it a friendly-name is recommended.

The Scaleway provider offers three ways of providing these credentials.
The following methods are supported, in this priority order:

1. [Environment variables](#environment-variables)
1. [Static credentials](#static-credentials)
1. [Shared configuration file](#shared-configuration-file)

### Environment variables

You can provide your credentials via the `SCW_ACCESS_KEY`, `SCW_SECRET_KEY` environment variables.

Example:

```hcl
provider "scaleway" {}
```

Usage:

```bash
$ export SCW_ACCESS_KEY="my-access-key"
$ export SCW_SECRET_KEY="my-secret-key"
$ terraform plan
```

### Static credentials

!> **Warning**: Hard-coding credentials into any Terraform configuration is not recommended, and risks secret leakage should this file ever be committed to a public version control system.

Static credentials can be provided by adding `access_key` and `secret_key` attributes in-line in the Scaleway provider block:

Example:

```hcl
provider "scaleway" {
  access_key = "my-access-key"
  secret_key = "my-secret-key"
}
```

### Shared configuration file

It is a YAML configuration file shared between the majority of the
[Scaleway developer tools](https://developers.scaleway.com/en/community-tools/#official-repos).
Its default location is `$HOME/.config/scw/config.yaml` (`%USERPROFILE%/.config/scw/config.yaml` on Windows).
If it fails to detect credentials inline, or in the environment, Terraform will check this file.

You can optionally specify a different location with `SCW_CONFIG_PATH` environment variable.
You can find more information about this configuration [in the documentation](https://github.com/scaleway/scaleway-sdk-go/blob/master/scw/README.md#scaleway-config).

## Arguments Reference

In addition to [generic provider arguments](https://www.terraform.io/docs/configuration/providers.html) (e.g. `alias` and `version`), the following arguments are supported in the Scaleway provider block:

| Provider Argument | [Environment Variables](#environment-variables) | Description                                                                                                                             | Mandatory |
|-------------------|-------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------|-----------|
| `access_key`      | `SCW_ACCESS_KEY`                                | [Scaleway access key](https://console.scaleway.com/project/credentials)                                                                 | ✅        |
| `secret_key`      | `SCW_SECRET_KEY`                                | [Scaleway secret key](https://console.scaleway.com/project/credentials)                                                                 | ✅        |
| `project_id`      | `SCW_DEFAULT_PROJECT_ID`                        | The [project ID](https://console.scaleway.com/project/settings) that will be used as default value for all resources.                   | ✅        |
| `region`          | `SCW_DEFAULT_REGION`                            | The [region](./guides/regions_and_zones.md#regions)  that will be used as default value for all resources. (`fr-par` if none specified) |           |
| `zone`            | `SCW_DEFAULT_ZONE`                              | The [zone](./guides/regions_and_zones.md#zones) that will be used as default value for all resources. (`fr-par-1` if none specified)    |           |

## Store terraform state on Scaleway S3-compatible object storage

[Scaleway object storage](https://www.scaleway.com/en/object-storage/) can be used to store your Terraform state.
Configure your backend as:

```
terraform {
  backend "s3" {
    bucket                      = "terraform_state"
    key                         = "my_state.tfstate"
    region                      = "fr-par"
    endpoint                    = "https://s3.fr-par.scw.cloud"
    access_key                  = "my-access-key"
    secret_key                  = "my-secret-key"
    skip_credentials_validation = true
    skip_region_validation      = true
  }
}
```

Be careful as no locking mechanism are yet supported.
Using scaleway object storage as terraform backend is not suitable if you work in a team with a risk of simultaneous access to the same plan.

Note: For security reason it's not recommended to store secrets in terraform files.
If you want to configure the backend with environment var, you need to use `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` [source](https://www.terraform.io/docs/backends/types/s3.html#access_key).

```bash
export AWS_ACCESS_KEY_ID=$SCW_ACCESS_KEY
export AWS_SECRET_ACCESS_KEY=$SCW_SECRET_KEY
```
