package scaleway

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	sdkValidation "github.com/scaleway/scaleway-sdk-go/validation"
)

func resourceScalewayBaremetalServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScalewayBaremetalServerCreate,
		ReadContext:   resourceScalewayBaremetalServerRead,
		UpdateContext: resourceScalewayBaremetalServerUpdate,
		DeleteContext: resourceScalewayBaremetalServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 0,
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(defaultBaremetalServerTimeout),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the server",
			},
			"hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Hostname of the server",
			},
			"offer": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID or name of the server offer",
			},
			"offer_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the server offer",
			},
			"os": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The base image of the server",
				ValidateFunc: validationUUID(),
			},
			"os_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The base image ID of the server",
			},
			"ssh_key_ids": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validationUUID(),
				},
				Required:    true,
				Description: "Array of SSH key IDs allowed to SSH to the server",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
				Description:  "Some description to associate to the server, max 255 characters",
			},
			"tags": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Array of tags to associate with the server",
			},
			"zone":            zoneSchema(),
			"organization_id": organizationIDSchema(),
			"project_id":      projectIDSchema(),
			"ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the IP",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The version of the IP",
						},
						"address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address of the IP",
						},
						"reverse": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Reverse of the IP",
						},
					},
				},
			},
			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceScalewayBaremetalServerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	baremetalAPI, zone, err := baremetalAPIWithZone(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	offerID := expandZonedID(d.Get("offer"))
	if !sdkValidation.IsUUID(offerID.ID) {
		o, err := baremetalAPI.GetOfferByName(&baremetal.GetOfferByNameRequest{
			OfferName: offerID.ID,
			Zone:      zone,
		})
		if err != nil {
			return diag.FromErr(err)
		}
		offerID = newZonedID(zone, o.ID)
	}

	req := &baremetal.CreateServerRequest{
		Zone:        zone,
		Name:        expandOrGenerateString(d.Get("name"), "bm"),
		ProjectID:   expandStringPtr(d.Get("project_id")),
		Description: d.Get("description").(string),
		OfferID:     offerID.ID,
		Tags:        expandStrings(d.Get("tags")),
	}

	if definedZone, ok := d.GetOk("zone"); ok {
		zone = scw.Zone(definedZone.(string))
		req.Zone = zone
	}

	server, err := baremetalAPI.CreateServer(req, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newZonedID(server.Zone, server.ID).String())

	_, err = baremetalAPI.WaitForServer(&baremetal.WaitForServerRequest{
		Zone:     server.Zone,
		ServerID: server.ID,
		Timeout:  scw.TimeDurationPtr(baremetalServerWaitForTimeout),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = baremetalAPI.InstallServer(&baremetal.InstallServerRequest{
		Zone:      server.Zone,
		ServerID:  server.ID,
		OsID:      expandZonedID(d.Get("os")).ID,
		Hostname:  expandStringWithDefault(d.Get("hostname"), server.Name),
		SSHKeyIDs: expandStrings(d.Get("ssh_key_ids")),
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = baremetalAPI.WaitForServerInstall(&baremetal.WaitForServerInstallRequest{
		Zone:     server.Zone,
		ServerID: server.ID,
		Timeout:  scw.TimeDurationPtr(baremetalServerWaitForTimeout),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceScalewayBaremetalServerRead(ctx, d, meta)
}

func resourceScalewayBaremetalServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	baremetalAPI, zonedID, err := baremetalAPIWithZoneAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	server, err := baremetalAPI.GetServer(&baremetal.GetServerRequest{
		Zone:     zonedID.Zone,
		ServerID: zonedID.ID,
	}, scw.WithContext(ctx))
	if err != nil {
		if is404Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	offer, err := baremetalAPI.GetOffer(&baremetal.GetOfferRequest{
		Zone:    server.Zone,
		OfferID: server.OfferID,
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("name", server.Name)
	_ = d.Set("zone", server.Zone.String())
	_ = d.Set("organization_id", server.OrganizationID)
	_ = d.Set("project_id", server.ProjectID)
	_ = d.Set("offer_id", newZonedID(server.Zone, offer.ID).String())
	_ = d.Set("tags", server.Tags)
	_ = d.Set("domain", server.Domain)
	_ = d.Set("ips", flattenBaremetalIPs(server.IPs))
	if server.Install != nil {
		_ = d.Set("os_id", newZonedID(server.Zone, server.Install.OsID).String())
		_ = d.Set("ssh_key_ids", server.Install.SSHKeyIDs)
	}
	_ = d.Set("description", server.Description)

	return nil
}

func resourceScalewayBaremetalServerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	baremetalAPI, zonedID, err := baremetalAPIWithZoneAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = baremetalAPI.UpdateServer(&baremetal.UpdateServerRequest{
		Zone:        zonedID.Zone,
		ServerID:    zonedID.ID,
		Name:        expandStringPtr(d.Get("name")),
		Description: expandStringPtr(d.Get("description")),
		Tags:        scw.StringsPtr(expandStrings(d.Get("tags"))),
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges("os", "ssh_key_ids") {
		installReq := &baremetal.InstallServerRequest{
			Zone:      zonedID.Zone,
			ServerID:  zonedID.ID,
			OsID:      expandZonedID(d.Get("os")).ID,
			Hostname:  expandStringWithDefault(d.Get("hostname"), d.Get("name").(string)),
			SSHKeyIDs: expandStrings(d.Get("ssh_key_ids")),
		}

		server, err := baremetalAPI.InstallServer(installReq, scw.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = baremetalAPI.WaitForServerInstall(&baremetal.WaitForServerInstallRequest{
			Zone:     server.Zone,
			ServerID: server.ID,
			Timeout:  scw.TimeDurationPtr(baremetalServerWaitForTimeout),
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceScalewayBaremetalServerRead(ctx, d, meta)
}

func resourceScalewayBaremetalServerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	baremetalAPI, zonedID, err := baremetalAPIWithZoneAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	server, err := baremetalAPI.DeleteServer(&baremetal.DeleteServerRequest{
		Zone:     zonedID.Zone,
		ServerID: zonedID.ID,
	}, scw.WithContext(ctx))

	if err != nil {
		if is404Error(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	_, err = baremetalAPI.WaitForServer(&baremetal.WaitForServerRequest{
		Zone:     server.Zone,
		ServerID: server.ID,
		Timeout:  scw.TimeDurationPtr(baremetalServerWaitForTimeout),
	})

	if err != nil && !is404Error(err) {
		return diag.FromErr(err)
	}

	return nil
}
