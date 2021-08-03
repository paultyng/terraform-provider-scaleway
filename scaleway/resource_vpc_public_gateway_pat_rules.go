package scaleway

import (
	"context"
	"net"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	vpcgw "github.com/scaleway/scaleway-sdk-go/api/vpcgw/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func resourceScalewayVPCPublicGatewayPATRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScalewayVPCPublicGatewayPATRulesCreate,
		ReadContext:   resourceScalewayVPCPublicGatewayPATRulesRead,
		UpdateContext: resourceScalewayVPCPublicGatewayPATRulesUpdate,
		DeleteContext: resourceScalewayVPCPublicGatewayPATRulesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 0,
		Schema: map[string]*schema.Schema{
			//_ = d.Set("private_ip", patRules.PrivateIP)
			"gateway_id": {
				Type:         schema.TypeString,
				Description:  "The ID of the gateway this PAT rule is applied to",
				Required:     true,
				ValidateFunc: validationUUIDorUUIDWithLocality(),
			},
			"private_ip": {
				Type:         schema.TypeString,
				Description:  "The private IP used in the PAT rule",
				Required:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"public_port": {
				Type:         schema.TypeInt,
				Description:  "The public port used in the PAT rule",
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 65535),
			},
			"private_port": {
				Type:         schema.TypeInt,
				Description:  "The private port used in the PAT rule",
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 65535),
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					vpcgw.PATRuleProtocolTCP.String(),
					vpcgw.PATRuleProtocolUDP.String(),
					vpcgw.PATRuleProtocolBoth.String(),
				}, true),
			},
			"zone": zoneSchema(),
			// Computed elements
			"organization_id": organizationIDSchema(),
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of the creation of the PAT rule",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of the last update of the PAT rule",
			},
		},
	}
}

func resourceScalewayVPCPublicGatewayPATRulesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcgwAPI, zone, err := vpcgwAPIWithZone(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	req := &vpcgw.CreatePATRuleRequest{
		Zone:        zone,
		GatewayID:   d.Get("gateway_id").(string),
		PublicPort:  uint32(d.Get("public_port").(int)),
		PrivateIP:   net.ParseIP(d.Get("private_ip").(string)),
		PrivatePort: uint32(d.Get("private_port").(int)),
		Protocol:    vpcgw.PATRuleProtocol(d.Get("protocol").(string)),
	}

	res, err := vpcgwAPI.CreatePATRule(req, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newZonedIDString(zone, res.ID))

	return resourceScalewayVPCPublicGatewayPATRulesRead(ctx, d, meta)
}

func resourceScalewayVPCPublicGatewayPATRulesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcgwAPI, zone, ID, err := vpcgwAPIWithZoneAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	patRules, err := vpcgwAPI.GetPATRule(&vpcgw.GetPATRuleRequest{
		PatRuleID: ID,
		Zone:      zone,
	}, scw.WithContext(ctx))
	if err != nil {
		if is404Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	_ = d.Set("created_at", patRules.CreatedAt.Format(time.RFC3339))
	_ = d.Set("updated_at", patRules.UpdatedAt.Format(time.RFC3339))
	_ = d.Set("gateway_id", patRules.GatewayID)
	_ = d.Set("private_ip", patRules.PrivateIP)
	_ = d.Set("private_port", patRules.PrivatePort)
	_ = d.Set("public_port", patRules.PublicPort)
	_ = d.Set("protocol", patRules.Protocol.String())
	_ = d.Set("zone", zone)

	return nil
}

func resourceScalewayVPCPublicGatewayPATRulesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, _, _, err := vpcgwAPIWithZoneAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceScalewayVPCPublicGatewayPATRulesRead(ctx, d, meta)
}

func resourceScalewayVPCPublicGatewayPATRulesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcgwAPI, zone, ID, err := vpcgwAPIWithZoneAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = vpcgwAPI.DeletePATRule(&vpcgw.DeletePATRuleRequest{
		PatRuleID: ID,
		Zone:      zone,
	}, scw.WithContext(ctx))

	if err != nil && !is404Error(err) {
		return diag.FromErr(err)
	}

	return nil
}
