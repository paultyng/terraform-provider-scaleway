---
subcategory: "IoT Hub"
page_title: "Scaleway: scaleway_iot_route"
---

# Resource: scaleway_iot_route

-> **Note:** This terraform resource is currently in beta and might include breaking change in future releases.

Creates and manages Scaleway IoT Routes. For more information, see the following:

- [API documentation](https://developers.scaleway.com/en/products/iot/api).
- [Product documentation](https://www.scaleway.com/en/docs/scaleway-iothub-route/)

## Example Usage

### Database Route

```terraform
resource "scaleway_iot_route" "main" {
	name   = "default"
	hub_id = scaleway_iot_hub.main.id
	topic  = "#"
	database {
		query  = <<-EOT
			INSERT INTO measurements(
				push_time,
				report_time,
				station_id,
				temperature,
				humidity
			) VALUES (
				NOW(),
				TIMESTAMP 'epoch' + (($PAYLOAD::jsonb->'last_reported')::integer * INTERVAL '1 second'),
				($PAYLOAD::jsonb->'station_id')::uuid,
				($PAYLOAD::jsonb->'temperature')::decimal,
				($PAYLOAD::jsonb->'humidity'):decimal:
			);
			EOT
		host   = scaleway_rdb_instance.iot.endpoint_ip
		port   = scaleway_rdb_instance.iot.endpoint_port
		dbname = "rdb"
		username = scaleway_rdb_instance.iot.user_name
		password = scaleway_rdb_instance.iot.password
	}
}

resource "scaleway_iot_hub" "main" {
	name         = "main"
	product_plan = "plan_shared"
}

resource "scaleway_rdb_instance" "iot" {
	name           = "iot"
	node_type      = "db-dev-s"
	engine         = "PostgreSQL-12"
	user_name      = "root"
	password       = "T3stP4ssw0rdD0N0tUs3!"
}
```

### S3 Route

```terraform
resource "scaleway_iot_route" "main" {
	name   = "main"
	hub_id = scaleway_iot_hub.main.id
	topic  = "#"
	s3 {
		bucket_region = scaleway_object_bucket.main.region
		bucket_name   = scaleway_object_bucket.main.name
		object_prefix = "foo"
		strategy      = "per_topic"
	}
}

resource "scaleway_iot_hub" "main" {
	name         = "main"
	product_plan = "plan_shared"
}

resource "scaleway_object_bucket" "main" {
	region = "fr-par"
	name = "my_awesome-bucket"
}
```

### Rest Route

```terraform
resource "scaleway_iot_route" "main" {
	name   = "main"
	hub_id = scaleway_iot_hub.main.id
	topic  = "#"
	rest {
		verb = "get"
		uri  = "http://scaleway.com"
		headers = {
			X-awesome-header = "my-awesome-value"
		}
	}
}

resource "scaleway_iot_hub" "main" {
	name         = "main"
	product_plan = "plan_shared"
}
```

## Argument Reference

~> **Important:** Updates to any value will recreate the IoT Route.

The following arguments are supported:

- `name` - (Required) The name of the IoT Route you want to create (e.g. `my-route`).

- `hub_id` - (Required) The hub ID to which the Route will be attached to.

- `topic` - (Required) The topic the Route subscribes to, wildcards allowed (e.g. `thelab/+/temperature/#`).

- `database` - (Optional) Configuration block for the database routes. See  [product documentation](https://www.scaleway.com/en/docs/scaleway-iothub-route/#-Database-Route) for a better understanding of the parameters.
    - `query` - (Required) The SQL query that will be executed when receiving a message ($TOPIC and $PAYLOAD variables are available, see documentation, e.g. `INSERT INTO mytable(date, topic, value) VALUES (NOW(), $TOPIC, $PAYLOAD)`).
    - `host` - (Required) The database hostname. Can be an IP or a FQDN.
    - `port` - (Required) The database port (e.g. `5432`)
    - `dbname` - (Required) The database name (e.g. `measurements`).
    - `username` - (Required) The database username.
    - `password` - (Required) The database password.

- `rest` (Optional) - Configuration block for the rest routes. See [product documentation](https://www.scaleway.com/en/docs/scaleway-iothub-route/#-REST-Route) for a better understanding of the parameters.
    - `verb` - (Required) The HTTP Verb used to call Rest URI (e.g. `post`).
    - `uri` - (Required) The URI of the Rest endpoint (e.g. `https://internal.mycompany.com/ingest/mqttdata`).
    - `headers` - (Required) a map of the extra headers to send with the HTTP call (e.g. `X-Header = Value`).

- `s3` (Optional) - Configuration block for the S3 routes. See [product documentation](https://www.scaleway.com/en/docs/scaleway-iothub-route/#-Scaleway-Object-Storage-Route) for a better understanding of the parameters.
    - `bucket_region` (Required) - The region of the S3 route's destination bucket (e.g. `fr-par`).
    - `bucket_name` (Required) - The name of the S3 route's destination bucket (e.g. `my-object-storage`).
    - `object_prefix` (Required) - The string to prefix object names with (e.g. `mykeyprefix-`).
    - `strategy` (Required) - How the S3 route's objects will be created (e.g. `per_topic`). See [documentation](https://www.scaleway.com/en/docs/scaleway-iothub-route/#-Messages-Store-Strategies) for behaviour details.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the Route.

~> **Important:** IoT routes' IDs are [regional](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{region}/{id}`, e.g. `fr-par/11111111-1111-1111-1111-111111111111`

- `region` - (Defaults to [provider](../index.md#region) `region`) The [region](../guides/regions_and_zones.md#regions) in which the Route is attached to.
- `created_at` - The date and time the Route was created.


## Import

IoT Routes can be imported using the `{region}/{id}`, e.g.

```bash
$ terraform import scaleway_iot_route.route01 fr-par/11111111-1111-1111-1111-111111111111
```

