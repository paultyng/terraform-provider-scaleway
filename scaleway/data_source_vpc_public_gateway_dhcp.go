package scaleway

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceScalewayVPCPublicGatewayDHCP() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := datasourceSchemaFromResourceSchema(resourceScalewayVPCPublicGatewayDHCP().Schema)

	dsSchema["dhcp_id"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		Description:  "The ID of the public gateway DHCP configuration",
		ValidateFunc: validationUUIDorUUIDWithLocality(),
	}

	return &schema.Resource{
		Schema:      dsSchema,
		ReadContext: dataSourceScalewayVPCPublicGatewayDHCPRead,
	}
}

func dataSourceScalewayVPCPublicGatewayDHCPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, zone, err := vpcgwAPIWithZone(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	dhcpID, _ := d.GetOk("dhcp_id")

	zonedID := datasourceNewZonedID(dhcpID, zone)
	d.SetId(zonedID)
	_ = d.Set("dhcp_id", zonedID)
	return resourceScalewayVPCPublicGatewayDHCPRead(ctx, d, meta)
}
