package vpcgw

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ipamAPI "github.com/scaleway/scaleway-sdk-go/api/ipam/v1"
	"github.com/scaleway/scaleway-sdk-go/api/vpcgw/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/cdf"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/dsf"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/regional"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/zonal"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/ipam"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/transport"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/types"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/verify"
)

func ResourceNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceVPCGatewayNetworkCreate,
		ReadContext:   ResourceVPCGatewayNetworkRead,
		UpdateContext: ResourceVPCGatewayNetworkUpdate,
		DeleteContext: ResourceVPCGatewayNetworkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create:  schema.DefaultTimeout(defaultTimeout),
			Read:    schema.DefaultTimeout(defaultTimeout),
			Update:  schema.DefaultTimeout(defaultTimeout),
			Delete:  schema.DefaultTimeout(defaultTimeout),
			Default: schema.DefaultTimeout(defaultTimeout),
		},
		SchemaVersion: 0,
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: verify.IsUUIDorUUIDWithLocality(),
				Description:      "The ID of the public gateway where connect to",
			},
			"private_network_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: verify.IsUUIDorUUIDWithLocality(),
				DiffSuppressFunc: dsf.Locality,
				Description:      "The ID of the private network where connect to",
			},
			"dhcp_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: verify.IsUUIDorUUIDWithLocality(),
				Description:      "The ID of the public gateway DHCP config",
				ConflictsWith:    []string{"static_address", "ipam_config"},
			},
			"enable_masquerade": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable masquerade on this network",
			},
			"enable_dhcp": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable DHCP config on this network",
			},
			"cleanup_dhcp": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Remove DHCP config on this network on destroy",
			},
			"static_address": {
				Type:          schema.TypeString,
				Description:   "The static IP address in CIDR on this network",
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.IsCIDR,
				ConflictsWith: []string{"dhcp_id", "ipam_config"},
			},
			"ipam_config": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				Description:   "Auto-configure the Gateway Network using IPAM (IP address management service)",
				ConflictsWith: []string{"dhcp_id", "static_address"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"push_default_route": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Defines whether the default route is enabled on that Gateway Network",
						},
						"ipam_ip_id": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							Description:      "Use this IPAM-booked IP ID as the Gateway's IP in this Private Network",
							ValidateDiagFunc: verify.IsUUIDorUUIDWithLocality(),
							DiffSuppressFunc: dsf.Locality,
						},
					},
				},
			},
			// Computed elements
			"mac_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The mac address on this network",
			},
			"private_ip": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of private IP addresses associated with the resource.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the IP address resource.",
						},
						"address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The private IP address.",
						},
					},
				},
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of the creation of the gateway network",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of the last update of the gateway network",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the Public Gateway's connection to the Private Network",
			},
			"zone": zonal.Schema(),
		},
		CustomizeDiff: cdf.LocalityCheck("gateway_id", "private_network_id", "dhcp_id"),
	}
}

func ResourceVPCGatewayNetworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api, zone, err := newAPIWithZone(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	gatewayID := zonal.ExpandID(d.Get("gateway_id").(string)).ID

	gateway, err := waitForVPCPublicGateway(ctx, api, zone, gatewayID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	req := &vpcgw.CreateGatewayNetworkRequest{
		Zone:             zone,
		GatewayID:        gateway.ID,
		PrivateNetworkID: regional.ExpandID(d.Get("private_network_id").(string)).ID,
		EnableMasquerade: *types.ExpandBoolPtr(d.Get("enable_masquerade")),
		EnableDHCP:       types.ExpandBoolPtr(d.Get("enable_dhcp")),
		IpamConfig:       expandIpamConfig(d.Get("ipam_config")),
	}
	staticAddress, staticAddressExist := d.GetOk("static_address")
	if staticAddressExist {
		address, err := types.ExpandIPNet(staticAddress.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		req.Address = &address
	}

	dhcpID, dhcpExist := d.GetOk("dhcp_id")
	if dhcpExist {
		dhcpZoned := zonal.ExpandID(dhcpID.(string))
		req.DHCPID = &dhcpZoned.ID
	}

	gatewayNetwork, err := transport.RetryOnTransientStateError(func() (*vpcgw.GatewayNetwork, error) {
		return api.CreateGatewayNetwork(req, scw.WithContext(ctx))
	}, func() (*vpcgw.Gateway, error) {
		tflog.Warn(ctx, "Public gateway is in transient state after waiting, retrying...")
		return waitForVPCPublicGateway(ctx, api, zone, gatewayID, d.Timeout(schema.TimeoutCreate))
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(zonal.NewIDString(zone, gatewayNetwork.ID))

	_, err = waitForVPCPublicGateway(ctx, api, zone, gatewayNetwork.GatewayID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = waitForVPCGatewayNetwork(ctx, api, zone, gatewayNetwork.ID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return ResourceVPCGatewayNetworkRead(ctx, d, m)
}

func ResourceVPCGatewayNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api, zone, ID, err := NewAPIWithZoneAndID(m, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	gatewayNetwork, err := waitForVPCGatewayNetwork(ctx, api, zone, ID, d.Timeout(schema.TimeoutRead))
	if err != nil {
		if httperrors.Is404(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	_, err = waitForVPCPublicGateway(ctx, api, zone, gatewayNetwork.GatewayID, d.Timeout(schema.TimeoutRead))
	if err != nil {
		return diag.FromErr(err)
	}

	if dhcp := gatewayNetwork.DHCP; dhcp != nil {
		_ = d.Set("dhcp_id", zonal.NewID(zone, dhcp.ID).String())
	}

	if staticAddress := gatewayNetwork.Address; staticAddress != nil {
		staticAddressValue, err := types.FlattenIPNet(*staticAddress)
		if err != nil {
			return diag.FromErr(err)
		}
		_ = d.Set("static_address", staticAddressValue)
	}

	if macAddress := gatewayNetwork.MacAddress; macAddress != nil {
		_ = d.Set("mac_address", types.FlattenStringPtr(macAddress).(string))
	}

	if enableDHCP := gatewayNetwork.EnableDHCP; enableDHCP {
		_ = d.Set("enable_dhcp", enableDHCP)
	}

	if ipamConfig := gatewayNetwork.IpamConfig; ipamConfig != nil {
		_ = d.Set("ipam_config", flattenIpamConfig(ipamConfig))
	}

	var cleanUpDHCPValue bool
	cleanUpDHCP, cleanUpDHCPExist := d.GetOk("cleanup_dhcp")
	if cleanUpDHCPExist {
		cleanUpDHCPValue = *types.ExpandBoolPtr(cleanUpDHCP)
	}

	gatewayNetwork, err = waitForVPCGatewayNetwork(ctx, api, zone, gatewayNetwork.ID, d.Timeout(schema.TimeoutRead))
	if err != nil {
		return diag.FromErr(err)
	}

	fetchRegion, err := zone.Region()
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("gateway_id", zonal.NewID(zone, gatewayNetwork.GatewayID).String())
	_ = d.Set("private_network_id", regional.NewIDString(fetchRegion, gatewayNetwork.PrivateNetworkID))
	_ = d.Set("enable_masquerade", gatewayNetwork.EnableMasquerade)
	_ = d.Set("cleanup_dhcp", cleanUpDHCPValue)
	_ = d.Set("created_at", gatewayNetwork.CreatedAt.Format(time.RFC3339))
	_ = d.Set("updated_at", gatewayNetwork.UpdatedAt.Format(time.RFC3339))
	_ = d.Set("zone", zone.String())
	_ = d.Set("status", gatewayNetwork.Status.String())

	var privateIP []map[string]interface{}
	if gatewayNetwork.PrivateNetworkID != "" {
		resourceID := gatewayNetwork.ID
		region, err := zone.Region()
		if err != nil {
			return diag.FromErr(err)
		}

		resourceType := ipamAPI.ResourceTypeVpcGatewayNetwork
		opts := &ipam.GetResourcePrivateIPsOptions{
			ResourceID:       &resourceID,
			ResourceType:     &resourceType,
			PrivateNetworkID: &gatewayNetwork.PrivateNetworkID,
		}
		privateIP, err = ipam.GetResourcePrivateIPs(ctx, m, region, opts)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	_ = d.Set("private_ip", privateIP)

	return nil
}

func ResourceVPCGatewayNetworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api, zone, ID, err := NewAPIWithZoneAndID(m, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = waitForVPCGatewayNetwork(ctx, api, zone, ID, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.FromErr(err)
	}

	updateRequest := &vpcgw.UpdateGatewayNetworkRequest{
		GatewayNetworkID: ID,
		Zone:             zone,
	}

	if d.HasChange("enable_masquerade") {
		updateRequest.EnableMasquerade = types.ExpandBoolPtr(d.Get("enable_masquerade"))
	}
	if d.HasChange("enable_dhcp") {
		updateRequest.EnableDHCP = types.ExpandBoolPtr(d.Get("enable_dhcp"))
	}
	if d.HasChange("dhcp_id") {
		dhcpID := zonal.ExpandID(d.Get("dhcp_id").(string)).ID
		updateRequest.DHCPID = &dhcpID
	}
	if d.HasChange("ipam_config") {
		updateRequest.IpamConfig = expandUpdateIpamConfig(d.Get("ipam_config"))
	}
	if d.HasChange("static_address") {
		staticAddress, staticAddressExist := d.GetOk("static_address")
		if staticAddressExist {
			address, err := types.ExpandIPNet(staticAddress.(string))
			if err != nil {
				return diag.FromErr(err)
			}
			updateRequest.Address = &address
		}
	}

	_, err = api.UpdateGatewayNetwork(updateRequest, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = waitForVPCGatewayNetwork(ctx, api, zone, ID, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.FromErr(err)
	}

	return ResourceVPCGatewayNetworkRead(ctx, d, m)
}

func ResourceVPCGatewayNetworkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api, zone, id, err := NewAPIWithZoneAndID(m, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	gwNetwork, err := waitForVPCGatewayNetwork(ctx, api, zone, id, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	req := &vpcgw.DeleteGatewayNetworkRequest{
		GatewayNetworkID: gwNetwork.ID,
		Zone:             gwNetwork.Zone,
		CleanupDHCP:      *types.ExpandBoolPtr(d.Get("cleanup_dhcp")),
	}
	err = api.DeleteGatewayNetwork(req, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = waitForVPCGatewayNetwork(ctx, api, zone, id, d.Timeout(schema.TimeoutDelete))
	if err != nil && !httperrors.Is404(err) {
		return diag.FromErr(err)
	}

	_, err = waitForVPCPublicGateway(ctx, api, zone, gwNetwork.GatewayID, d.Timeout(schema.TimeoutDelete))
	if err != nil && !httperrors.Is404(err) {
		return diag.FromErr(err)
	}

	return nil
}
