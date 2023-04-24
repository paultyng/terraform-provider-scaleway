package scaleway

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/scaleway-sdk-go/validation"
)

func DataSourceAvailabilityZone() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceAvailabilityZoneRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Region is represented as a Geographical area such as France",
				Default:     scw.RegionFrPar,
			},
			"zones": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "Availability Zones (AZ)",
			},
		},
	}
}

func dataSourceAvailabilityZoneRead(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	regionStr := d.Get("region").(string)

	if !validation.IsRegion(regionStr) {
		return diag.FromErr(SingularDataSourceFindError("Availability Zone", fmt.Errorf("not a supported region %s", regionStr)))
	}

	region := scw.Region(regionStr)
	d.SetId(regionStr)
	_ = d.Set("zones", region.GetZones())

	return nil
}
