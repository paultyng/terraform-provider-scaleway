package scaleway

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func resourceScalewayLb() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScalewayLbCreate,
		ReadContext:   resourceScalewayLbRead,
		UpdateContext: resourceScalewayLbUpdate,
		DeleteContext: resourceScalewayLbDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(defaultLbLbTimeout),
		},
		SchemaVersion: 0,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the lb",
			},
			"type": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: diffSuppressFuncIgnoreCase,
				Description:      "The type of load-balancer you want to create",
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Array of tags to associate with the load-balancer",
			},
			"ip_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "The load-balance public IP ID",
				DiffSuppressFunc: diffSuppressFuncLocality,
			},
			"ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The load-balance public IP address",
			},
			"region":          regionSchema(),
			"organization_id": organizationIDSchema(),
			"project_id":      projectIDSchema(),
		},
	}
}

func resourceScalewayLbCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	lbAPI, region, err := lbAPIWithRegion(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	createReq := &lb.CreateLBRequest{
		Region:    region,
		IPID:      expandStringPtr(expandID(d.Get("ip_id"))),
		ProjectID: expandStringPtr(d.Get("project_id")),
		Name:      expandOrGenerateString(d.Get("name"), "lb"),
		Type:      d.Get("type").(string),
	}

	if definedRegion, ok := d.GetOk("region"); ok {
		region = scw.Region(definedRegion.(string))
		createReq.Region = region
	}

	if raw, ok := d.GetOk("tags"); ok {
		for _, tag := range raw.([]interface{}) {
			createReq.Tags = append(createReq.Tags, tag.(string))
		}
	}
	res, err := lbAPI.CreateLB(createReq, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newRegionalIDString(region, res.ID))

	_, err = lbAPI.WaitForLb(&lb.WaitForLBRequest{
		Region:  region,
		LBID:    res.ID,
		Timeout: scw.TimeDurationPtr(defaultInstanceServerWaitTimeout),
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceScalewayLbRead(ctx, d, meta)
}

func resourceScalewayLbRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	lbAPI, region, ID, err := lbAPIWithRegionAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := lbAPI.GetLB(&lb.GetLBRequest{
		Region: region,
		LBID:   ID,
	}, scw.WithContext(ctx))

	if err != nil {
		if is404Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	_ = d.Set("name", res.Name)
	_ = d.Set("region", string(region))
	_ = d.Set("organization_id", res.OrganizationID)
	_ = d.Set("project_id", res.ProjectID)
	_ = d.Set("tags", res.Tags)
	// For now API return lowercase lb type. This should be fix in a near future on the API side
	_ = d.Set("type", strings.ToUpper(res.Type))
	_ = d.Set("ip_id", newRegionalIDString(region, res.IP[0].ID))
	_ = d.Set("ip_address", res.IP[0].IPAddress)

	return nil
}

func resourceScalewayLbUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	lbAPI, region, ID, err := lbAPIWithRegionAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges("name", "tags") {
		req := &lb.UpdateLBRequest{
			Region: region,
			LBID:   ID,
			Name:   d.Get("name").(string),
			Tags:   expandStrings(d.Get("tags")),
		}

		_, err = lbAPI.UpdateLB(req, scw.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceScalewayLbRead(ctx, d, meta)
}

func resourceScalewayLbDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	lbAPI, region, ID, err := lbAPIWithRegionAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = lbAPI.DeleteLB(&lb.DeleteLBRequest{
		Region:    region,
		LBID:      ID,
		ReleaseIP: false,
	}, scw.WithContext(ctx))

	if err != nil && !is404Error(err) {
		return diag.FromErr(err)
	}

	_, err = lbAPI.WaitForLb(&lb.WaitForLBRequest{
		LBID:    ID,
		Region:  region,
		Timeout: scw.TimeDurationPtr(LbWaitForTimeout),
	}, scw.WithContext(ctx))

	if err != nil && !is404Error(err) {
		return diag.FromErr(err)
	}

	return nil
}
