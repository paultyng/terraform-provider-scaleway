package applesilicon

import (
	"context"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/meta"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/verify"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	applesilicon "github.com/scaleway/scaleway-sdk-go/api/applesilicon/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/zonal"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/account"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/types"
)

func ResourceServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceAppleSiliconServerCreate,
		ReadContext:   ResourceAppleSiliconServerRead,
		UpdateContext: ResourceAppleSiliconServerUpdate,
		DeleteContext: ResourceAppleSiliconServerDelete,
		Timeouts: &schema.ResourceTimeout{
			Create:  schema.DefaultTimeout(defaultAppleSiliconServerTimeout),
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
			"enable_vpc": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether or not to enable VPC access",
			},
			"private_network": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The private networks to attach to the server",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:             schema.TypeString,
							Description:      "The private network ID",
							Required:         true,
							ValidateDiagFunc: verify.IsUUIDorUUIDWithLocality(),
							StateFunc: func(i interface{}) string {
								return locality.ExpandID(i.(string))
							},
						},
						"ipam_ip_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: verify.IsUUIDorUUIDWithLocality(),
							},
							Description: "List of IPAM IP IDs to attach to the server",
						},
						// computed
						"vlan": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The VLAN ID associated to the private network",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The private network status",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time of the creation of the private network",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time of the last update of the private network",
						},
					},
				},
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
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The state of the server",
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
			"vpc_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The VPC status of the server",
			},

			// Common
			"zone":            zonal.Schema(),
			"organization_id": account.OrganizationIDSchema(),
			"project_id":      account.ProjectIDSchema(),
		},
	}
}

func ResourceAppleSiliconServerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	asAPI, zone, err := newAPIWithZone(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	createReq := &applesilicon.CreateServerRequest{
		Name:      types.ExpandOrGenerateString(d.Get("name"), "m1"),
		Type:      d.Get("type").(string),
		ProjectID: d.Get("project_id").(string),
		EnableVpc: d.Get("enable_vpc").(bool),
	}

	res, err := asAPI.CreateServer(createReq, scw.WithContext(ctx))

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(zonal.NewIDString(zone, res.ID))

	_, err = waitForAppleSiliconServer(ctx, asAPI, zone, res.ID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	if pn, ok := d.GetOk("private_network"); ok {
		privateNetworkAPI := applesilicon.NewPrivateNetworkAPI(meta.ExtractScwClient(m))
		req := &applesilicon.PrivateNetworkAPISetServerPrivateNetworksRequest{
			Zone:                       zone,
			ServerID:                   res.ID,
			PerPrivateNetworkIpamIPIDs: expandPrivateNetworks(pn),
		}
		_, err := privateNetworkAPI.SetServerPrivateNetworks(req, scw.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = waitForAppleSiliconPrivateNetworkServer(ctx, privateNetworkAPI, zone, res.ID, d.Timeout(schema.TimeoutCreate))
	}

	return ResourceAppleSiliconServerRead(ctx, d, m)
}

func ResourceAppleSiliconServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	asAPI, zone, ID, err := NewAPIWithZoneAndID(m, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	privateNetworkAPI := applesilicon.NewPrivateNetworkAPI(meta.ExtractScwClient(m))

	res, err := asAPI.GetServer(&applesilicon.GetServerRequest{
		Zone:     zone,
		ServerID: ID,
	}, scw.WithContext(ctx))
	if err != nil {
		if httperrors.Is404(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	_ = d.Set("name", res.Name)
	_ = d.Set("type", res.Type)
	_ = d.Set("state", res.Status.String())
	_ = d.Set("created_at", res.CreatedAt.Format(time.RFC3339))
	_ = d.Set("updated_at", res.UpdatedAt.Format(time.RFC3339))
	_ = d.Set("deletable_at", res.DeletableAt.Format(time.RFC3339))
	_ = d.Set("ip", res.IP.String())
	_ = d.Set("vnc_url", res.VncURL)
	_ = d.Set("vpc_status", res.VpcStatus)
	_ = d.Set("zone", res.Zone.String())
	_ = d.Set("organization_id", res.OrganizationID)
	_ = d.Set("project_id", res.ProjectID)
	listPrivateNetworks, err := privateNetworkAPI.ListServerPrivateNetworks(&applesilicon.PrivateNetworkAPIListServerPrivateNetworksRequest{
		Zone:     res.Zone,
		ServerID: &res.ID,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	pnRegion, err := res.Zone.Region()
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("private_network", flattenPrivateNetworks(pnRegion, listPrivateNetworks.ServerPrivateNetworks))

	return nil
}

func ResourceAppleSiliconServerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	asAPI, zone, ID, err := NewAPIWithZoneAndID(m, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	req := &applesilicon.UpdateServerRequest{
		Zone:     zone,
		ServerID: ID,
	}

	if d.HasChange("name") {
		req.Name = types.ExpandStringPtr(d.Get("name"))
	}

	if d.HasChange("enable_vpc") {
		enableVpc := d.Get("enable_vpc").(bool)
		req.EnableVpc = &enableVpc
	}

	_, err = asAPI.UpdateServer(req, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return ResourceAppleSiliconServerRead(ctx, d, m)
}

func ResourceAppleSiliconServerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	asAPI, zone, ID, err := NewAPIWithZoneAndID(m, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = asAPI.DeleteServer(&applesilicon.DeleteServerRequest{
		Zone:     zone,
		ServerID: ID,
	}, scw.WithContext(ctx))

	if err != nil && !httperrors.Is404(err) {
		return diag.FromErr(err)
	}

	return nil
}
