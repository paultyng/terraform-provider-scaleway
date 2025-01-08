---
page_title: "Migrating from Scaleway Cockpit to the New Infrastructure"
---

# How to Migrate from Deprecated Resource `scaleway_cockpit` to `scaleway_cockpit_source`

## Overview

This guide provides a step-by-step process to remove the deprecated `scaleway_cockpit` resource from your Terraform configurations and transition to the new `scaleway_cockpit_source` resource. Note that this migration involves breaking down the functionalities of `scaleway_cockpit` into multiple specialized resources to manage endpoints effectively.

> **Note:**
> Scaleway Cockpit plans are scheduled for deprecation on **January 1st, 2025**. While the retention period for your logs and metrics will remain unchanged, you will be able to edit the retention period for metrics, logs, and traces for free during the Beta period.

## Prerequisites

### Ensure the Latest Provider Version

Ensure your Scaleway provider is updated to at least version `2.49.0`.

```hcl
terraform {
  required_providers {
    scaleway = {
      source  = "scaleway/scaleway"
      version = "~> 2.49.0"
    }
  }
}

provider "scaleway" {
  # Configuration details
}
```

Run the following command to initialize the updated provider:

```bash
terraform init
```

## Migrating Resources

### Transitioning from `scaleway_cockpit`

The `scaleway_cockpit` resource is deprecated. Its functionalities, including endpoint management, are now divided across multiple specialized resources. Below are the steps to migrate:

#### Deprecated Resource: `scaleway_cockpit`

The following resource will no longer be supported after January 1st, 2025:

```hcl
resource "scaleway_cockpit" "main" {
  project_id = "11111111-1111-1111-1111-111111111111"
  plan       = "premium"
}
```

#### New Resources

To handle specific functionalities previously managed by `scaleway_cockpit`, you need to use the following resources:

**Data Source Management:**

In the deprecated `scaleway_cockpit` resource, the `plan` argument determined the retention period for logs, metrics, and traces. Now, retention periods are set individually for each data source using the `retention_days` argument in `scaleway_cockpit_source` resources.

```hcl
resource "scaleway_account_project" "project" {
  name = "test project data source"
}

resource "scaleway_cockpit_source" "metrics" {
  project_id     = scaleway_account_project.project.id
  name           = "metrics-source"
  type           = "metrics"
  retention_days = 6 # Customize retention period (1-365 days)
}

resource "scaleway_cockpit_source" "logs" {
  project_id     = scaleway_account_project.project.id
  name           = "logs-source"
  type           = "logs"
  retention_days = 30
}

resource "scaleway_cockpit_source" "traces" {
  project_id     = scaleway_account_project.project.id
  name           = "traces-source"
  type           = "traces"
  retention_days = 15
}
```

**Alert Manager:**

To retrieve the deprecated `alertmanager_url`, you must now explicitly create an Alert Manager using the `scaleway_cockpit_alert_manager` resource:

```hcl
resource "scaleway_cockpit_alert_manager" "alert_manager" {
  project_id            = scaleway_account_project.project.id
  enable_managed_alerts = true

  contact_points {
    email = "alert1@example.com"
  }

  contact_points {
    email = "alert2@example.com"
  }
}
```

**Grafana User:**

To retrieve the deprecated `grafana_url`, you must create a Grafana user. Creating the user will trigger the creation of the Grafana instance:

```hcl
resource "scaleway_cockpit_grafana_user" "main" {
  project_id = scaleway_account_project.project.id
  login      = "my-awesome-user"
  role       = "editor"
}
```

### Notes on Regionalization

- As of September 2024, Cockpit resources are regionalized for improved flexibility and resilience. Update your queries in Grafana to use the new regionalized data sources.
- Metrics, logs, and traces now have dedicated resources that allow granular control over retention policies.

### Before and After Example

#### Before: Using `scaleway_cockpit` to Retrieve Endpoints

```hcl
resource "scaleway_cockpit" "main" {
  project_id = "11111111-1111-1111-1111-111111111111"
  plan       = "premium"
}

output "endpoints" {
  value = scaleway_cockpit.main.endpoints
}
```

#### After: Using Specialized Resources

To retrieve all endpoints (metrics, logs, traces, alert manager, and Grafana):

```hcl
resource "scaleway_cockpit_source" "metrics" {
  project_id     = scaleway_account_project.project.id
  name           = "metrics-source"
  type           = "metrics"
  retention_days = 6
}

resource "scaleway_cockpit_source" "logs" {
  project_id     = scaleway_account_project.project.id
  name           = "logs-source"
  type           = "logs"
  retention_days = 30
}

resource "scaleway_cockpit_source" "traces" {
  project_id     = scaleway_account_project.project.id
  name           = "traces-source"
  type           = "traces"
  retention_days = 15
}

resource "scaleway_cockpit_alert_manager" "alert_manager" {
  project_id = scaleway_account_project.project.id
  enable_managed_alerts = true
}

resource "scaleway_cockpit_grafana_user" "main" {
  project_id = scaleway_account_project.project.id
  login      = "my-awesome-user"
  role       = "editor"
}

output "endpoints" {
  value = {
    metrics        = scaleway_cockpit_source.metrics.url
    logs           = scaleway_cockpit_source.logs.url
    traces         = scaleway_cockpit_source.traces.url
    alert_manager  = scaleway_cockpit_alert_manager.alert_manager.alert_manager_url
    grafana        = scaleway_cockpit_grafana_user.main.grafana_url
  }
}
```

## Importing Resources

### Import a Cockpit Source

To import an existing `scaleway_cockpit_source` resource:

```bash
terraform import scaleway_cockpit_source.main fr-par/11111111-1111-1111-1111-111111111111
```

### Import a Grafana User

To import an existing Grafana user:

```bash
terraform import scaleway_cockpit_grafana_user.main 11111111-1111-1111-1111-111111111111
```

## Conclusion

By following this guide, you can successfully transition from the deprecated `scaleway_cockpit` resource to the new set of specialized resources. This ensures compatibility with the latest Terraform provider and Scaleway's updated infrastructure.

