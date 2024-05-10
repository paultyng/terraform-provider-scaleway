---
page_title: "Provider: Scaleway"
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

- Get your [Scaleway credentials](https://console.scaleway.com/iam/api-keys)
- Get your [project ID](https://console.scaleway.com/project/settings)
- Initialize a Terraform working directory: `terraform init`
- Generate and show the execution plan: `terraform plan`
- Build the infrastructure: `terraform apply`

```hcl
variable "project_id" {
  type        = string
  description = "Your project ID."
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
    action   = "accept"
    port     = "22"
    ip_range = "0.0.0.0/0"
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
  image      = "ubuntu_jammy"

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
2. [Static credentials](#static-credentials)
3. [Shared configuration file](#shared-configuration-file)

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

~> **Warning**: Hard-coding credentials into any Terraform configuration is not recommended, and risks secret leakage should this file ever be committed to a public version control system.

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

This method also supports a `profile` configuration:

Example:

If your shared configuration file contains:

```yaml
profiles:
  myProfile:
    access_key: xxxxxxxxxxxxxxxxxxxx
    secret_key: xxxxxxxx-xxx-xxxx-xxxx-xxxxxxxxxxx
    default_organization_id: xxxxxxxx-xxx-xxxx-xxxx-xxxxxxxxxxx 
    default_project_id: xxxxxxxx-xxx-xxxx-xxxx-xxxxxxxxxxx
    default_zone: fr-par-2
    default_region: fr-par
    api_url: https://api.scaleway.com
    insecure: false
```

You can invoke and use this profile in the provider declaration:

```hcl
provider "scaleway" {
  alias   = "p2"
  profile = "myProfile"
}

resource "scaleway_instance_ip" "server_ip" {
  provider = scaleway.p2
}
```

## Arguments Reference

In addition to [generic provider arguments](https://www.terraform.io/docs/configuration/providers.html) (e.g. `alias` and `version`), the following arguments are supported in the Scaleway provider block:

| Provider Argument | [Environment Variables](#environment-variables) | Description                                                                                                                                      | Mandatory |
| ----------------- | ----------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------ | --------- |
| `access_key`      | `SCW_ACCESS_KEY`                                | [Scaleway access key](https://console.scaleway.com/project/credentials)                                                                          | ✅         |
| `secret_key`      | `SCW_SECRET_KEY`                                | [Scaleway secret key](https://console.scaleway.com/project/credentials)                                                                          | ✅         |
| `project_id`      | `SCW_DEFAULT_PROJECT_ID`                        | The [project ID](https://console.scaleway.com/project/settings) that will be used as default value for project-scoped resources.                | ✅         |
| `organization_id` | `SCW_DEFAULT_ORGANIZATION_ID`                   | The [organization ID](https://console.scaleway.com/organization/settings) that will be used as default value for organization-scoped resources. |           |
| `region`          | `SCW_DEFAULT_REGION`                            | The [region](./guides/regions_and_zones.md#regions)  that will be used as default value for all resources. (`fr-par` if none specified)          |           |
| `zone`            | `SCW_DEFAULT_ZONE`                              | The [zone](./guides/regions_and_zones.md#zones) that will be used as default value for all resources. (`fr-par-1` if none specified)             |           |

## Store terraform state on Scaleway S3-compatible object storage

[Scaleway object storage](https://www.scaleway.com/en/object-storage/) can be used to store your Terraform state.
Configure your backend as:

```
terraform {
  backend "s3" {
    bucket                      = "terraform-state"
    key                         = "my_state.tfstate"
    region                      = "fr-par"
    endpoint                    = "https://s3.fr-par.scw.cloud"
    access_key                  = "my-access-key"
    secret_key                  = "my-secret-key"
    skip_credentials_validation = true
    skip_region_validation      = true
    # Need terraform>=1.6.1
    skip_requesting_account_id  = true
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

## Custom User-Agent Information

The Scaleway Terraform Provider allows you to append custom information to the User-Agent header of HTTP requests made to the Scaleway API. This can be useful for tracking requests for auditing, logging, or analytics purposes.

To append custom information to the User-Agent header, you can use the `TF_APPEND_USER_AGENT` environment variable. The value you set for this variable will be appended to the User-Agent header of all HTTP requests made by the provider.

For example, to add custom information indicating the request is coming from a specific CI/CD job or system, you could set the environment variable as follows:

```bash
$ export TF_APPEND_USER_AGENT="CI/CD System XYZ Job #1234"
```

## Debugging a deployment

In case you want to [debug a deployment](https://www.terraform.io/internals/debugging), you can use the following command to increase the level of verbosity.

`SCW_DEBUG=1 TF_LOG=WARN TF_LOG_PROVIDER=DEBUG terraform apply`

- `SCW_DEBUG`: set the debug level of the scaleway SDK.
- `TF_LOG`: set the level of the Terraform logging.
- `TF_LOG_PROVIDER`: set the level of the Scaleway Terraform provider logging.

### Submitting a bug report or a feature request

In case you find something wrong with the scaleway provider, please submit a bug report on the [Terraform provider repository](https://github.com/scaleway/terraform-provider-scaleway/issues/new/choose).
If it is a bug report, please include a **minimal** snippet of the Terraform configuration that triggered the error.
This helps a lot to debug the issue.
