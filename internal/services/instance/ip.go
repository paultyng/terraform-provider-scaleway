package instance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	instanceSDK "github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/zonal"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/account"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/types"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/verify"
)

func ResourceIP() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceInstanceIPCreate,
		ReadContext:   ResourceInstanceIPRead,
		UpdateContext: ResourceInstanceIPUpdate,
		DeleteContext: ResourceInstanceIPDelete,
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
			"prefix": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP prefix",
			},
			"type": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				Description:      "The type of instance IP",
				ValidateDiagFunc: verify.ValidateEnum[instanceSDK.IPType](),
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
			"tags": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "The tags associated with the ip",
			},
			"zone":            zonal.Schema(),
			"organization_id": account.OrganizationIDSchema(),
			"project_id":      account.ProjectIDSchema(),
		},
	}
}

func ResourceInstanceIPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	instanceAPI, zone, err := newAPIWithZone(d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	req := &instanceSDK.CreateIPRequest{
		Zone:    zone,
		Project: types.ExpandStringPtr(d.Get("project_id")),
		Type:    instanceSDK.IPType(d.Get("type").(string)),
	}
	tags := types.ExpandStrings(d.Get("tags"))
	if len(tags) > 0 {
		req.Tags = tags
	}
	res, err := instanceAPI.CreateIP(req, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	reverseRaw, ok := d.GetOk("reverse")
	if ok {
		reverseStrPtr := types.ExpandStringPtr(reverseRaw)
		req := &instanceSDK.UpdateIPRequest{
			IP:      res.IP.ID,
			Reverse: &instanceSDK.NullableStringValue{Value: *reverseStrPtr},
			Zone:    zone,
		}
		_, err = instanceAPI.UpdateIP(req, scw.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(zonal.NewIDString(zone, res.IP.ID))
	return ResourceInstanceIPRead(ctx, d, m)
}

func ResourceInstanceIPUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	instanceAPI, zone, ID, err := NewAPIWithZoneAndID(m, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	req := &instanceSDK.UpdateIPRequest{
		IP:   ID,
		Zone: zone,
	}

	if d.HasChange("tags") {
		req.Tags = types.ExpandUpdatedStringsPtr(d.Get("tags"))
	}

	_, err = instanceAPI.UpdateIP(req, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return ResourceInstanceIPRead(ctx, d, m)
}

func ResourceInstanceIPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	instanceAPI, zone, ID, err := NewAPIWithZoneAndID(m, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := instanceAPI.GetIP(&instanceSDK.GetIPRequest{
		IP:   ID,
		Zone: zone,
	}, scw.WithContext(ctx))
	if err != nil {
		// We check for 403 because instanceSDK API returns 403 for a deleted IP
		if httperrors.Is404(err) || httperrors.Is403(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	address := res.IP.Address.String()
	prefix := res.IP.Prefix.String()
	if prefix == types.NetIPNil {
		ipnet := scw.IPNet{}
		_ = (&ipnet).UnmarshalJSON([]byte("\"" + res.IP.Address.String() + "\""))
		prefix = ipnet.String()
	}
	if address == types.NetIPNil {
		address = res.IP.Prefix.IP.String()
	}

	_ = d.Set("address", address)
	_ = d.Set("prefix", prefix)
	_ = d.Set("zone", zone)
	_ = d.Set("organization_id", res.IP.Organization)
	_ = d.Set("project_id", res.IP.Project)
	_ = d.Set("reverse", res.IP.Reverse)
	_ = d.Set("type", res.IP.Type)

	if len(res.IP.Tags) > 0 {
		_ = d.Set("tags", types.FlattenSliceString(res.IP.Tags))
	}

	if res.IP.Server != nil {
		_ = d.Set("server_id", zonal.NewIDString(res.IP.Zone, res.IP.Server.ID))
	} else {
		_ = d.Set("server_id", "")
	}

	return nil
}

func ResourceInstanceIPDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	instanceAPI, zone, ID, err := NewAPIWithZoneAndID(m, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = instanceAPI.DeleteIP(&instanceSDK.DeleteIPRequest{
		IP:   ID,
		Zone: zone,
	}, scw.WithContext(ctx))
	if err != nil {
		// We check for 403 because instanceSDK API returns 403 for a deleted IP
		if httperrors.Is404(err) || httperrors.Is403(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}
