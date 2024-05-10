---
page_title: "Using Scaleway Messaging and Queuing service with NATS Terraform provider"
---

# How to use Scaleway Messaging and Queuing config

This guide shows the combination of Scaleway Messaging and Queuing configuration with to terraform NATS Jetstream
provider. Il will allow you to provision and manage NATS Jetstream resources.

## Examples

```hcl
terraform {
  required_providers {
    scaleway = {
      source  = "scaleway/scaleway"
    }
    jetstream = {
      source  = "nats-io/jetstream"
      version = "~> 0.0.34"
    }
  }
}

resource "scaleway_mnq_nats_account" "account" {}

resource "scaleway_mnq_nats_credentials" "creds" {
  account_id = scaleway_mnq_nats_account.account.id
}

provider "jetstream" {
  servers = scaleway_mnq_nats_account.account.endpoint
  credential_data = scaleway_mnq_nats_credentials.creds.file
}

// Use any jetstream resources
resource "jetstream_stream" "ORDERS" {
  name     = "ORDERS"
  subjects = ["ORDERS.*"]
  storage  = "file"
  max_age  = 60 * 60 * 24 * 365
}
```
