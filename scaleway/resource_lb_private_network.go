package scaleway

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func resourceScalewayLbPrivateNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScalewayLbPrivateNetworkCreate,
		ReadContext:   resourceScalewayLbPrivateNetworkRead,
		DeleteContext: resourceScalewayLbPrivateNetworkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(defaultLbLbTimeout),
		},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{Version: 0, Type: lbUpgradeV1SchemaType(), Upgrade: lbUpgradeV1SchemaUpgradeFunc},
		},
		Schema: map[string]*schema.Schema{
			"lb_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validationUUIDorUUIDWithLocality(),
				Description:  "The load-balancer ID",
			},
			"private_network_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validationUUIDorUUIDWithLocality(),
				Description:  "The Private Network ID",
			},
			"static_config": {
				ConflictsWith: []string{"dhcp_config"},
				Description:   "Define two local IP addresses of your choice for each load balancer instance",
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsIPAddress,
				},
			},
			"dhcp_config": {
				ConflictsWith: []string{"static_config"},
				Description:   "Set to true if you want to let DHCP assign IP addresses",
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{},
				},
			},

			// Readonly attributes
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of private network connection",
			},
		},
	}
}

func resourceScalewayLbPrivateNetworkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	lbAPI, _, err := lbAPIWithZone(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	zone, lbID, err := parseZonedID(d.Get("lb_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	zone_pn, pnID, err := parseZonedID(d.Get("private_network_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if zone_pn != zone {
		return diag.Errorf("LB and Private Network must be in the same zone (got %s and %s)", zone, zone_pn)
	}

	createReq := &lb.ZonedAPIAttachPrivateNetworkRequest{
		Zone:             zone,
		LBID:             lbID,
		PrivateNetworkID: pnID,
		StaticConfig:     expandLbPrivateNetworkStaticConfig(d.Get("static_config")),
		DHCPConfig:       expandLbPrivateNetworkDHCPConfig(d.Get("dhcp_config")),
	}

	res, err := lbAPI.AttachPrivateNetwork(createReq, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newZonedIDString(zone, res.PrivateNetworkID))
	d.Set("status", res.Status)

	return resourceScalewayLbPrivateNetworkRead(ctx, d, meta)
}

func resourceScalewayLbPrivateNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	lbAPI, zone, pnID, err := lbAPIWithZoneAndID(meta, d.Get("private_network_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, lbID, err := parseZonedID(d.Get("lb_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := lbAPI.ListLBPrivateNetworks(&lb.ZonedAPIListLBPrivateNetworksRequest{
		Zone: zone,
		LBID: lbID,
	}, scw.WithContext(ctx))

	if err != nil {
		if is404Error(err) {
			d.Set("lb_id", "")
			return nil
		}
		return diag.FromErr(err)
	}

	pn := findPn(res.PrivateNetwork, pnID)
	if pn == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("lb_id", newZonedIDString(zone, pn.LB.ID))
	_ = d.Set("private_network_id", newZonedIDString(zone, pn.PrivateNetworkID))
	_ = d.Set("static_config", flattenLbPrivateNetworkStaticConfig(pn.StaticConfig))
	_ = d.Set("dhcp_config", flattenLbPrivateNetworkDHCPConfig(pn.DHCPConfig))
	_ = d.Set("status", pn.Status)

	return nil
}

func findPn(private_networks []*lb.PrivateNetwork, id string) *lb.PrivateNetwork {
	for _, pn := range private_networks {
		if pn.PrivateNetworkID == id {
			return pn
		}
	}
	return nil
}

func resourceScalewayLbPrivateNetworkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	lbAPI, zone, pnID, err := lbAPIWithZoneAndID(meta, d.Get("private_network_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, lbID, err := parseZonedID(d.Get("lb_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	err = lbAPI.DetachPrivateNetwork(&lb.ZonedAPIDetachPrivateNetworkRequest{
		Zone:             zone,
		LBID:             lbID,
		PrivateNetworkID: pnID,
	})
	if err != nil && !is404Error(err) {
		return diag.FromErr(err)
	}

	return nil
}
