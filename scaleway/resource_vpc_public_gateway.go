package scaleway

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpcgw "github.com/scaleway/scaleway-sdk-go/api/vpcgw/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	retryIntervalVPCPublicGatewayNetwork = 30 * time.Second
)

func resourceScalewayVPCPublicGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScalewayVPCPublicGatewayCreate,
		ReadContext:   resourceScalewayVPCPublicGatewayRead,
		UpdateContext: resourceScalewayVPCPublicGatewayUpdate,
		DeleteContext: resourceScalewayVPCPublicGatewayDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 0,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "name of the gateway",
			},
			"type": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "gateway type",
				DiffSuppressFunc: diffSuppressFuncIgnoreCase,
			},
			"upstream_dns_servers": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "override the gateway's default recursive DNS servers, if DNS features are enabled",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ip_id": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				Description:      "attach an existing IP to the gateway",
				DiffSuppressFunc: diffSuppressFuncLocality,
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The tags associated with public gateway",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"project_id": projectIDSchema(),
			"zone":       zoneSchema(),
			// Computed elements
			"organization_id": organizationIDSchema(),
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of the creation of the public gateway",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of the last update of the public gateway",
			},
		},
	}
}

func resourceScalewayVPCPublicGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcgwAPI, zone, err := vpcgwAPIWithZone(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	req := &vpcgw.CreateGatewayRequest{
		Name:               expandOrGenerateString(d.Get("name"), "pn"),
		Type:               d.Get("type").(string),
		Tags:               expandStrings(d.Get("tags")),
		UpstreamDNSServers: expandStrings(d.Get("upstream_dns_servers")),
		ProjectID:          d.Get("project_id").(string),
		Zone:               zone,
	}

	if ipID, ok := d.GetOk("ip_id"); ok {
		req.IPID = expandStringPtr(expandZonedID(ipID).ID)
	}

	res, err := vpcgwAPI.CreateGateway(req, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newZonedIDString(zone, res.ID))

	retryInterval := retryGWTimeout
	_, err = vpcgwAPI.WaitForGateway(&vpcgw.WaitForGatewayRequest{
		GatewayID:     res.ID,
		Timeout:       scw.TimeDurationPtr(defaultVPCGatewayTimeout),
		RetryInterval: &retryInterval,
		Zone:          zone,
	}, scw.WithContext(ctx))
	// check err waiting process
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceScalewayVPCPublicGatewayRead(ctx, d, meta)
}

func resourceScalewayVPCPublicGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcgwAPI, zone, ID, err := vpcgwAPIWithZoneAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	retryInterval := retryGWTimeout
	gateway, err := vpcgwAPI.WaitForGateway(&vpcgw.WaitForGatewayRequest{
		GatewayID:     ID,
		Timeout:       scw.TimeDurationPtr(defaultVPCGatewayTimeout),
		RetryInterval: &retryInterval,
		Zone:          zone,
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("name", gateway.Name)
	_ = d.Set("organization_id", gateway.OrganizationID)
	_ = d.Set("project_id", gateway.ProjectID)
	_ = d.Set("created_at", gateway.CreatedAt.Format(time.RFC3339))
	_ = d.Set("updated_at", gateway.UpdatedAt.Format(time.RFC3339))
	_ = d.Set("zone", zone)
	_ = d.Set("tags", gateway.Tags)
	_ = d.Set("upstream_dns_servers", gateway.UpstreamDNSServers)
	_ = d.Set("ip_id", newZonedID(zone, gateway.IP.ID).String())

	return nil
}

func resourceScalewayVPCPublicGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcgwAPI, zone, ID, err := vpcgwAPIWithZoneAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	retryInterval := retryGWTimeout
	gateway, err := vpcgwAPI.WaitForGateway(&vpcgw.WaitForGatewayRequest{
		GatewayID:     ID,
		Timeout:       scw.TimeDurationPtr(defaultVPCGatewayTimeout),
		RetryInterval: &retryInterval,
		Zone:          zone,
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges("name", "tags", "upstream_dns_servers") {
		updateRequest := &vpcgw.UpdateGatewayRequest{
			GatewayID:          gateway.ID,
			Zone:               zone,
			Name:               scw.StringPtr(d.Get("name").(string)),
			Tags:               scw.StringsPtr(expandStrings(d.Get("tags"))),
			UpstreamDNSServers: scw.StringsPtr(expandStrings(d.Get("upstream_dns_servers"))),
		}

		_, err = vpcgwAPI.UpdateGateway(updateRequest, scw.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	_, err = vpcgwAPI.WaitForGateway(&vpcgw.WaitForGatewayRequest{
		GatewayID:     ID,
		Timeout:       scw.TimeDurationPtr(defaultVPCGatewayTimeout),
		RetryInterval: &retryInterval,
		Zone:          zone,
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceScalewayVPCPublicGatewayRead(ctx, d, meta)
}

func resourceScalewayVPCPublicGatewayDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcgwAPI, zone, ID, err := vpcgwAPIWithZoneAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	retryInterval := retryGWTimeout
	_, err = vpcgwAPI.WaitForGateway(&vpcgw.WaitForGatewayRequest{
		GatewayID:     ID,
		Timeout:       scw.TimeDurationPtr(defaultVPCGatewayTimeout),
		RetryInterval: &retryInterval,
		Zone:          zone,
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	err = vpcgwAPI.DeleteGateway(&vpcgw.DeleteGatewayRequest{
		GatewayID: ID,
		Zone:      zone,
	}, scw.WithContext(ctx))
	if err != nil && !is404Error(err) {
		return diag.FromErr(err)
	}

	return nil
}
