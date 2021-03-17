package scaleway

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func resourceScalewayInstanceIP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScalewayInstanceIPCreate,
		ReadContext:   resourceScalewayInstanceIPRead,
		DeleteContext: resourceScalewayInstanceIPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(defaultInstanceIPTimeout),
		},
		SchemaVersion: 0,
		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP address",
			},
			"reverse": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The reverse DNS for this IP",
			},
			"server_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The server associated with this IP",
			},
			"zone":            zoneSchema(),
			"organization_id": organizationIDSchema(),
			"project_id":      projectIDSchema(),
		},
	}
}

func resourceScalewayInstanceIPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	instanceAPI, zone, err := instanceAPIWithZone(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	req := &instance.CreateIPRequest{
		Zone:    zone,
		Project: expandStringPtr(d.Get("project_id")),
	}

	if definedZone, ok := d.GetOk("zone"); ok {
		zone = scw.Zone(definedZone.(string))
		req.Zone = zone
	}

	res, err := instanceAPI.CreateIP(req, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newZonedIDString(zone, res.IP.ID))
	return resourceScalewayInstanceIPRead(ctx, d, meta)
}

func resourceScalewayInstanceIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	instanceAPI, zone, ID, err := instanceAPIWithZoneAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := instanceAPI.GetIP(&instance.GetIPRequest{
		IP:   ID,
		Zone: zone,
	}, scw.WithContext(ctx))

	if err != nil {
		// We check for 403 because instance API returns 403 for a deleted IP
		if is404Error(err) || is403Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	_ = d.Set("address", res.IP.Address.String())
	_ = d.Set("zone", zone)
	_ = d.Set("organization_id", res.IP.Organization)
	_ = d.Set("project_id", res.IP.Project)
	_ = d.Set("reverse", res.IP.Reverse)

	if res.IP.Server != nil {
		_ = d.Set("server_id", newZonedIDString(res.IP.Zone, res.IP.Server.ID))
	} else {
		_ = d.Set("server_id", "")
	}

	return nil
}

func resourceScalewayInstanceIPDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	instanceAPI, zone, ID, err := instanceAPIWithZoneAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = instanceAPI.DeleteIP(&instance.DeleteIPRequest{
		IP:   ID,
		Zone: zone,
	}, scw.WithContext(ctx))

	if err != nil && !is404Error(err) && !is403Error(err) {
		return diag.FromErr(err)
	}

	return nil
}
