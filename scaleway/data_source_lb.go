package scaleway

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func dataSourceScalewayLb() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := datasourceSchemaFromResourceSchema(resourceScalewayLb().Schema)

	// Set 'Optional' schema elements
	addOptionalFieldsToSchema(dsSchema, "name", "region")

	dsSchema["name"].ConflictsWith = []string{"lb_id"}
	dsSchema["lb_id"] = &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		Description:   "The ID of the load-balancer",
		ValidateFunc:  validationUUIDorUUIDWithLocality(),
		ConflictsWith: []string{"name"},
	}

	return &schema.Resource{
		ReadContext: dataSourceScalewayLbRead,
		Schema:      dsSchema,
	}
}

func dataSourceScalewayLbRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api, region, err := lbAPIWithRegion(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	lbID, ok := d.GetOk("lb_id")
	if !ok { // Get LB by name.
		res, err := api.ListLBs(&lb.ListLBsRequest{
			Region:    region,
			Name:      expandStringPtr(d.Get("name")),
			ProjectID: expandStringPtr(d.Get("project_id")),
		}, scw.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
		if len(res.LBs) == 0 {
			return diag.FromErr(fmt.Errorf("no lbs found with the name %s", d.Get("name")))
		}
		if len(res.LBs) > 1 {
			return diag.FromErr(fmt.Errorf("%d lbs found with the same name %s", len(res.LBs), d.Get("name")))
		}
		lbID = res.LBs[0].ID
	}

	regionalID := datasourceNewRegionalizedID(lbID, region)
	d.SetId(regionalID)
	err = d.Set("lb_id", regionalID)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceScalewayLbRead(ctx, d, meta)
}
