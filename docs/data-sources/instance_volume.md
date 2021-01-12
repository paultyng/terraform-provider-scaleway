---
page_title: "scaleway_instance_volume Data Source - terraform-provider-scaleway"
subcategory: ""
description: |-
  
---

# Data Source `scaleway_instance_volume`





## Schema

### Optional

- **id** (String) The ID of this resource.
- **name** (String) The name of the volume
- **volume_id** (String) The ID of the volume
- **zone** (String) The zone you want to attach the resource to

### Read-only

- **from_snapshot_id** (String) Create a volume based on a image
- **from_volume_id** (String) Create a copy of an existing volume
- **organization_id** (String) The organization_id you want to attach the resource to
- **project_id** (String) The project_id you want to attach the resource to
- **server_id** (String) The server associated with this volume
- **size_in_gb** (Number) The size of the volume in gigabyte
- **type** (String) The volume type


