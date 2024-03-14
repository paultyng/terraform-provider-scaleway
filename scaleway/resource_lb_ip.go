package scaleway

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lbSDK "github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/regional"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/zonal"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/types"
)

func ResourceScalewayLbIP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScalewayLbIPCreate,
		ReadContext:   resourceScalewayLbIPRead,
		UpdateContext: resourceScalewayLbIPUpdate,
		DeleteContext: resourceScalewayLbIPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Read:    schema.DefaultTimeout(defaultLbLbTimeout),
			Update:  schema.DefaultTimeout(defaultLbLbTimeout),
			Delete:  schema.DefaultTimeout(defaultLbLbTimeout),
			Default: schema.DefaultTimeout(defaultLbLbTimeout),
		},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{Version: 0, Type: lbUpgradeV1SchemaType(), Upgrade: LbUpgradeV1SchemaUpgradeFunc},
		},
		Schema: map[string]*schema.Schema{
			"reverse": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The reverse domain name for this IP",
			},
			"zone": zonal.Schema(),
			// Computed
			"organization_id": organizationIDSchema(),
			"project_id":      projectIDSchema(),
			"ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The load-balancer public IP address",
			},
			"lb_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the load balancer attached to this IP, if any",
			},
			"region": regional.ComputedSchema(),
		},
	}
}

func resourceScalewayLbIPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	lbAPI, zone, err := lbAPIWithZone(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	zoneAttribute, ok := d.GetOk("zone")
	if ok {
		zone = scw.Zone(zoneAttribute.(string))
	}

	createReq := &lbSDK.ZonedAPICreateIPRequest{
		Zone:      zone,
		ProjectID: types.ExpandStringPtr(d.Get("project_id")),
		Reverse:   types.ExpandStringPtr(d.Get("reverse")),
	}

	res, err := lbAPI.CreateIP(createReq, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(zonal.NewIDString(zone, res.ID))

	return resourceScalewayLbIPRead(ctx, d, m)
}

func resourceScalewayLbIPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	lbAPI, zone, ID, err := LbAPIWithZoneAndID(m, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	ip, err := lbAPI.GetIP(&lbSDK.ZonedAPIGetIPRequest{
		Zone: zone,
		IPID: ID,
	}, scw.WithContext(ctx))
	if err != nil {
		if httperrors.Is404(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	// check lb state if it is attached
	if ip.LBID != nil {
		_, err = waitForLB(ctx, lbAPI, zone, *ip.LBID, d.Timeout(schema.TimeoutRead))
		if err != nil {
			if httperrors.Is403(err) {
				d.SetId("")
				return nil
			}
			return diag.FromErr(err)
		}
	}

	// set the region from zone
	region, err := zone.Region()
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("region", string(region))
	_ = d.Set("zone", ip.Zone.String())
	_ = d.Set("organization_id", ip.OrganizationID)
	_ = d.Set("project_id", ip.ProjectID)
	_ = d.Set("ip_address", ip.IPAddress)
	_ = d.Set("reverse", ip.Reverse)
	_ = d.Set("lb_id", types.FlattenStringPtr(ip.LBID))

	return nil
}

func resourceScalewayLbIPUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	lbAPI, zone, ID, err := LbAPIWithZoneAndID(m, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	var ip *lbSDK.IP
	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		res, errGet := lbAPI.GetIP(&lbSDK.ZonedAPIGetIPRequest{
			Zone: zone,
			IPID: ID,
		}, scw.WithContext(ctx))
		if err != nil {
			if httperrors.Is403(errGet) {
				return resource.RetryableError(errGet)
			}
			return resource.NonRetryableError(errGet)
		}

		ip = res
		return nil
	})
	if err != nil {
		if httperrors.Is404(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if ip.LBID != nil {
		_, err = waitForLB(ctx, lbAPI, zone, *ip.LBID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			if httperrors.Is403(err) {
				d.SetId("")
				return nil
			}
			return diag.FromErr(err)
		}
	}

	if d.HasChange("reverse") {
		req := &lbSDK.ZonedAPIUpdateIPRequest{
			Zone:    zone,
			IPID:    ID,
			Reverse: types.ExpandStringPtr(d.Get("reverse")),
		}

		_, err = lbAPI.UpdateIP(req, scw.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if ip.LBID != nil {
		_, err = waitForLB(ctx, lbAPI, zone, *ip.LBID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			if httperrors.Is403(err) {
				d.SetId("")
				return nil
			}
			return diag.FromErr(err)
		}
	}

	return resourceScalewayLbIPRead(ctx, d, m)
}

//gocyclo:ignore
func resourceScalewayLbIPDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	lbAPI, zone, ID, err := LbAPIWithZoneAndID(m, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	var ip *lbSDK.IP
	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		res, errGet := lbAPI.GetIP(&lbSDK.ZonedAPIGetIPRequest{
			Zone: zone,
			IPID: ID,
		}, scw.WithContext(ctx))
		if err != nil {
			if httperrors.Is403(errGet) {
				return resource.RetryableError(errGet)
			}
			return resource.NonRetryableError(errGet)
		}

		ip = res
		return nil
	})
	if err != nil {
		return diag.FromErr(err)
	}

	// check lb state
	if ip != nil && ip.LBID != nil {
		_, err = waitForLB(ctx, lbAPI, zone, *ip.LBID, d.Timeout(schema.TimeoutDelete))
		if err != nil {
			if httperrors.Is403(err) {
				d.SetId("")
				return nil
			}
			return diag.FromErr(err)
		}
	}

	err = lbAPI.ReleaseIP(&lbSDK.ZonedAPIReleaseIPRequest{
		Zone: zone,
		IPID: ID,
	}, scw.WithContext(ctx))

	if err != nil && !httperrors.Is404(err) {
		return diag.FromErr(err)
	}

	// check lb state
	if ip != nil && ip.LBID != nil {
		_, err = waitForLB(ctx, lbAPI, zone, *ip.LBID, d.Timeout(schema.TimeoutDelete))
		if err != nil {
			if httperrors.Is404(err) || httperrors.Is403(err) {
				d.SetId("")
				return nil
			}
			return diag.FromErr(err)
		}
	}

	return nil
}
