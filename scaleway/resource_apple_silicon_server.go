package scaleway

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	applesilicon "github.com/scaleway/scaleway-sdk-go/api/applesilicon/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func resourceScalewayAppleSiliconServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScalewayAppleSiliconServerCreate,
		ReadContext:   resourceScalewayAppleSiliconServerRead,
		UpdateContext: resourceScalewayAppleSiliconServerUpdate,
		DeleteContext: resourceScalewayAppleSiliconServerDelete,
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(defaultAppleSiliconServerTimeout),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 0,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the server",
				Computed:    true,
				Optional:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "Type of the server",
				Required:    true,
				ForceNew:    true,
			},

			// Computed
			"ip": {
				Type:        schema.TypeString,
				Description: "IPv4 address of the server",
				Computed:    true,
			},
			"vnc_url": {
				Type:        schema.TypeString,
				Description: "VNC url use to connect remotely to the desktop GUI",
				Computed:    true,
			},

			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of the creation of the server",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of the last update of the server",
			},
			"deletable_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The minimal date and time on which you can delete this server due to Apple licence",
			},

			// Common
			"zone":            zoneSchema(),
			"organization_id": organizationIDSchema(),
			"project_id":      projectIDSchema(),
		},
	}
}

func resourceScalewayAppleSiliconServerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	asAPI, zone, err := asAPIWithZone(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	createReq := &applesilicon.CreateServerRequest{
		Zone:      zone,
		Name:      expandOrGenerateString(d.Get("name"), "m1"),
		Type:      d.Get("type").(string),
		ProjectID: d.Get("project_id").(string),
	}

	if definedZone, ok := d.GetOk("zone"); ok {
		zone = scw.Zone(definedZone.(string))
		createReq.Zone = zone
	}

	res, err := asAPI.CreateServer(createReq, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newZonedIDString(zone, res.ID))

	_, err = asAPI.WaitForServer(&applesilicon.WaitForServerRequest{
		ServerID:      res.ID,
		Timeout:       scw.TimeDurationPtr(defaultAppleSiliconServerTimeout),
		RetryInterval: nil,
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceScalewayRdbInstanceRead(ctx, d, meta)
}

func resourceScalewayAppleSiliconServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	asAPI, zone, ID, err := asAPIWithZoneAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := asAPI.GetServer(&applesilicon.GetServerRequest{
		Zone:     zone,
		ServerID: ID,
	}, scw.WithContext(ctx))
	if err != nil {
		if is404Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	_ = d.Set("name", res.Name)
	_ = d.Set("type", res.Type)
	_ = d.Set("created_at", res.CreatedAt.Format(time.RFC3339))
	_ = d.Set("updated_at", res.UpdatedAt.Format(time.RFC3339))
	_ = d.Set("deletable_at", res.DeletableAt.Format(time.RFC3339))
	_ = d.Set("ip", res.IP.String())
	_ = d.Set("vnc_url", res.VncURL)

	_ = d.Set("zone", zone.String())
	_ = d.Set("organization_id", res.OrganizationID)
	_ = d.Set("project_id", res.ProjectID)

	return nil
}

func resourceScalewayAppleSiliconServerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	asAPI, zone, ID, err := asAPIWithZoneAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	req := &applesilicon.UpdateServerRequest{
		Zone:     zone,
		ServerID: ID,
	}

	if d.HasChange("name") {
		req.Name = expandStringPtr(d.Get("name"))
	}

	_, err = asAPI.UpdateServer(req, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceScalewayAppleSiliconServerRead(ctx, d, meta)
}

func resourceScalewayAppleSiliconServerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	asAPI, zone, ID, err := asAPIWithZoneAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = asAPI.DeleteServer(&applesilicon.DeleteServerRequest{
		Zone:     zone,
		ServerID: ID,
	}, scw.WithContext(ctx))

	if err != nil && !is404Error(err) {
		return diag.FromErr(err)
	}

	return nil
}
