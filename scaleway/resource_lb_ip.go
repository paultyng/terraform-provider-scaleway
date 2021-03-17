package scaleway

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func resourceScalewayLbIP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScalewayLbIPCreate,
		ReadContext:   resourceScalewayLbIPRead,
		UpdateContext: resourceScalewayLbIPUpdate,
		DeleteContext: resourceScalewayLbIPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(defaultLbLbTimeout),
		},
		SchemaVersion: 0,
		Schema: map[string]*schema.Schema{
			"reverse": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The reverse domain name for this IP",
			},
			"region":          regionSchema(),
			"organization_id": organizationIDSchema(),
			"project_id":      projectIDSchema(),
			// Computed
			"ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The load-balancer public IP address",
			},
			"lb_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the loadbalancer attached to this IP, if any",
			},
		},
	}
}

func resourceScalewayLbIPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	lbAPI, region, err := lbAPIWithRegion(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	createReq := &lb.CreateIPRequest{
		Region:    region,
		ProjectID: expandStringPtr(d.Get("project_id")),
		Reverse:   expandStringPtr(d.Get("reverse")),
	}

	if definedRegion, ok := d.GetOk("region"); ok {
		region = scw.Region(definedRegion.(string))
		createReq.Region = region
	}

	res, err := lbAPI.CreateIP(createReq, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newRegionalIDString(region, res.ID))

	return resourceScalewayLbIPRead(ctx, d, meta)
}

func resourceScalewayLbIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	lbAPI, region, ID, err := lbAPIWithRegionAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := lbAPI.GetIP(&lb.GetIPRequest{
		Region: region,
		IPID:   ID,
	}, scw.WithContext(ctx))

	if err != nil {
		if is404Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	_ = d.Set("region", string(region))
	_ = d.Set("organization_id", res.OrganizationID)
	_ = d.Set("project_id", res.ProjectID)
	_ = d.Set("ip_address", res.IPAddress)
	_ = d.Set("reverse", res.Reverse)
	_ = d.Set("lb_id", flattenStringPtr(res.LBID))

	return nil
}

func resourceScalewayLbIPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	lbAPI, region, ID, err := lbAPIWithRegionAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("reverse") {
		req := &lb.UpdateIPRequest{
			Region:  region,
			IPID:    ID,
			Reverse: expandStringPtr(d.Get("reverse")),
		}

		_, err = lbAPI.UpdateIP(req, scw.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceScalewayLbIPRead(ctx, d, meta)
}

func resourceScalewayLbIPDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	lbAPI, region, ID, err := lbAPIWithRegionAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = lbAPI.ReleaseIP(&lb.ReleaseIPRequest{
		Region: region,
		IPID:   ID,
	}, scw.WithContext(ctx))

	if err != nil && !is404Error(err) {
		return diag.FromErr(err)
	}

	return nil
}
