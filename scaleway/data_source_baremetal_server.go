package scaleway

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func dataSourceScalewayBaremetalServer() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := datasourceSchemaFromResourceSchema(resourceScalewayBaremetalServer().Schema)

	// Set 'Optional' schema elements
	addOptionalFieldsToSchema(dsSchema, "name")

	dsSchema["name"].ConflictsWith = []string{"server_id"}
	dsSchema["server_id"] = &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		Description:   "The ID of the server",
		ValidateFunc:  validationUUIDorUUIDWithLocality(),
		ConflictsWith: []string{"name"},
	}

	return &schema.Resource{
		ReadContext: dataSourceScalewayBaremetalServerRead,
		Schema:      dsSchema,
	}
}

func dataSourceScalewayBaremetalServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api, zone, err := baremetalAPIWithZone(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	serverID, ok := d.GetOk("server_id")
	if !ok { // Get server by zone and name.
		res, err := api.ListServers(&baremetal.ListServersRequest{
			Zone: zone,
			Name: scw.StringPtr(d.Get("name").(string)),
		}, scw.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
		if len(res.Servers) == 0 {
			return diag.FromErr(fmt.Errorf("no servers found with the name %s", d.Get("name")))
		}
		if len(res.Servers) > 1 {
			return diag.FromErr(fmt.Errorf("%d servers found with the same name %s", len(res.Servers), d.Get("name")))
		}
		serverID = res.Servers[0].ID
	}

	zoneID := datasourceNewZonedID(serverID, zone)
	d.SetId(zoneID)
	err = d.Set("server_id", zoneID)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceScalewayBaremetalServerRead(ctx, d, meta)
}
