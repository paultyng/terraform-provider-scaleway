package instance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	ipamAPI "github.com/scaleway/scaleway-sdk-go/api/ipam/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/cdf"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/dsf"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/regional"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/zonal"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/ipam"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/types"
)

func ResourcePrivateNIC() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceInstancePrivateNICCreate,
		ReadContext:   ResourceInstancePrivateNICRead,
		UpdateContext: ResourceInstancePrivateNICUpdate,
		DeleteContext: ResourceInstancePrivateNICDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create:  schema.DefaultTimeout(defaultInstancePrivateNICWaitTimeout),
			Read:    schema.DefaultTimeout(defaultInstancePrivateNICWaitTimeout),
			Update:  schema.DefaultTimeout(defaultInstancePrivateNICWaitTimeout),
			Delete:  schema.DefaultTimeout(defaultInstancePrivateNICWaitTimeout),
			Default: schema.DefaultTimeout(defaultInstancePrivateNICWaitTimeout),
		},
		Schema: map[string]*schema.Schema{
			"server_id": {
				Type:        schema.TypeString,
				Description: "The server ID",
				Required:    true,
				ForceNew:    true,
			},
			"private_network_id": {
				Type:             schema.TypeString,
				Description:      "The private network ID",
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: dsf.Locality,
			},
			"mac_address": {
				Type:        schema.TypeString,
				Description: "MAC address of the NIC",
				Computed:    true,
			},
			"tags": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "The tags associated with the private-nic",
			},
			"ip_ids": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "IPAM ip list, should be for internal use only",
				ForceNew:    true,
			},
			"private_ip": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of private IP addresses associated with the resource",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the IP address resource",
						},
						"address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The private IP address",
						},
					},
				},
			},
			"ipam_ip_ids": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				ForceNew:    true,
				Description: "IPAM IDs of a pre-reserved IP addresses to assign to the Instance in the requested private network",
			},
			"zone": zonal.Schema(),
		},
		CustomizeDiff: cdf.LocalityCheck("server_id", "private_network_id"),
	}
}

func ResourceInstancePrivateNICCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	instanceAPI, zone, err := newAPIWithZone(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = waitForServer(ctx, instanceAPI, zone, locality.ExpandID(d.Get("server_id")), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	createPrivateNICRequest := &instance.CreatePrivateNICRequest{
		Zone:             zone,
		ServerID:         zonal.ExpandID(d.Get("server_id").(string)).ID,
		PrivateNetworkID: regional.ExpandID(d.Get("private_network_id").(string)).ID,
		Tags:             types.ExpandStrings(d.Get("tags")),
		IPIDs:            types.ExpandStringsPtr(d.Get("ip_ids")),
		IpamIPIDs:        locality.ExpandIDs(d.Get("ipam_ip_ids")),
	}

	privateNIC, err := instanceAPI.CreatePrivateNIC(
		createPrivateNICRequest,
		scw.WithContext(ctx),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = waitForPrivateNIC(ctx, instanceAPI, zone, privateNIC.PrivateNic.ServerID, privateNIC.PrivateNic.ID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(
		zonal.NewNestedIDString(
			zone,
			privateNIC.PrivateNic.ServerID,
			privateNIC.PrivateNic.ID,
		),
	)

	return ResourceInstancePrivateNICRead(ctx, d, m)
}

func ResourceInstancePrivateNICRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	instanceAPI, _, err := newAPIWithZone(d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	zone, privateNICID, serverID, err := zonal.ParseNestedID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	privateNIC, err := waitForPrivateNIC(ctx, instanceAPI, zone, serverID, privateNICID, d.Timeout(schema.TimeoutRead))
	if err != nil {
		if httperrors.Is404(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	fetchRegion, err := zone.Region()
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("zone", zone)
	_ = d.Set("server_id", zonal.NewID(zone, privateNIC.ServerID).String())
	_ = d.Set("private_network_id", regional.NewIDString(fetchRegion, privateNIC.PrivateNetworkID))
	_ = d.Set("mac_address", privateNIC.MacAddress)

	if len(privateNIC.Tags) > 0 {
		_ = d.Set("tags", privateNIC.Tags)
	}

	region, err := zone.Region()
	if err != nil {
		return diag.FromErr(err)
	}

	resourceType := ipamAPI.ResourceTypeInstancePrivateNic
	opts := &ipam.GetResourcePrivateIPsOptions{
		ResourceID:       &privateNIC.ID,
		ResourceType:     &resourceType,
		PrivateNetworkID: &privateNIC.PrivateNetworkID,
	}
	privateIP, err := ipam.GetResourcePrivateIPs(ctx, m, region, opts)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("private_ip", privateIP)

	return nil
}

func ResourceInstancePrivateNICUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	instanceAPI, _, err := newAPIWithZone(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	zone, privateNICID, serverID, err := zonal.ParseNestedID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("tags") {
		_, err := instanceAPI.UpdatePrivateNIC(
			&instance.UpdatePrivateNICRequest{
				Zone:         zone,
				ServerID:     serverID,
				PrivateNicID: privateNICID,
				Tags:         types.ExpandUpdatedStringsPtr(d.Get("tags")),
			},
			scw.WithContext(ctx),
		)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return ResourceInstancePrivateNICRead(ctx, d, m)
}

func ResourceInstancePrivateNICDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	instanceAPI, _, err := newAPIWithZone(d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	zone, privateNICID, serverID, err := zonal.ParseNestedID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = waitForPrivateNIC(ctx, instanceAPI, zone, serverID, privateNICID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		if httperrors.Is404(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	err = instanceAPI.DeletePrivateNIC(&instance.DeletePrivateNICRequest{
		ServerID:     serverID,
		PrivateNicID: privateNICID,
		Zone:         zone,
	}, scw.WithContext(ctx))
	if err != nil {
		if httperrors.Is404(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	_, err = waitForPrivateNIC(ctx, instanceAPI, zone, serverID, privateNICID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		if httperrors.Is404(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}
