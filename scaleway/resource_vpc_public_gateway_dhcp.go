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

func resourceScalewayVPCPublicGatewayDHCP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScalewayVPCPublicGatewayDHCPCreate,
		ReadContext:   resourceScalewayVPCPublicGatewayDHCPRead,
		UpdateContext: resourceScalewayVPCPublicGatewayDHCPUpdate,
		DeleteContext: resourceScalewayVPCPublicGatewayDHCPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 0,
		Schema: map[string]*schema.Schema{
			"project_id": projectIDSchema(),
			"zone":       zoneSchema(),
			"subnet": {
				Type:         schema.TypeString,
				Description:  "subnet for the DHCP server",
				ValidateFunc: validation.IsCIDR,
				Required:     true,
			},
			"address": {
				Type:        schema.TypeString,
				Description: "Address: address of the DHCP server. This will be the gateway's address in the private network. Defaults to the first address of the subnet",
				Optional:    true,
				Computed:    true,
			},
			"pool_low": {
				Type:         schema.TypeString,
				Description:  "low IP (included) of the dynamic address pool. Defaults to the second address of the subnet.",
				ValidateFunc: validation.IsIPAddress,
				Computed:     true,
				Optional:     true,
			},
			"pool_high": {
				Type:         schema.TypeString,
				Description:  "High IP (included) of the dynamic address pool. Defaults to the last address of the subnet.",
				ValidateFunc: validation.IsIPAddress,
				Computed:     true,
				Optional:     true,
			},
			"enable_dynamic": {
				Type:        schema.TypeBool,
				Description: "Whether to enable dynamic pooling of IPs. By turning the dynamic pool off, only pre-existing DHCP reservations will be handed out. Defaults to true.",
				Computed:    true,
				Optional:    true,
			},
			"valid_lifetime": {
				Type:        schema.TypeInt,
				Description: "For how long, in seconds, will DHCP entries will be valid. Defaults to 1h (3600s).",
				Computed:    true,
				Optional:    true,
			},
			"renew_timer": {
				Type:        schema.TypeInt,
				Description: "After how long, in seconds, a renew will be attempted. Must be 30s lower than `rebind_timer`. Defaults to 50m (3000s).",
				Computed:    true,
				Optional:    true,
			},
			"rebind_timer": {
				Type:        schema.TypeInt,
				Description: "After how long, in seconds, a DHCP client will query for a new lease if previous renews fail. Must be 30s lower than `valid_lifetime`. Defaults to 51m (3060s).",
				Computed:    true,
				Optional:    true,
			},
			"push_default_route": {
				Type:        schema.TypeBool,
				Description: "Whether the gateway should push a default route to DHCP clients or only hand out IPs. Defaults to true",
				Computed:    true,
				Optional:    true,
			},
			"push_dns_server": {
				Type:        schema.TypeBool,
				Description: "Whether the gateway should push custom DNS servers to clients. This allows for instance hostname -> IP resolution. Defaults to true.",
				Computed:    true,
				Optional:    true,
			},
			"dns_server_override": {
				Type:        schema.TypeList,
				Description: "override the DNS server list pushed to DHCP clients, instead of the gateway itself",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"dns_search": {
				Type:        schema.TypeList,
				Description: "additional DNS search paths",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"dns_local_name": {
				Type:        schema.TypeString,
				Description: "TLD given to hostnames in the Private Network. Allowed characters are `a-z0-9-.`. Defaults to the slugified Private Network name if created along a GatewayNetwork, or else to `priv`.",
				Optional:    true,
				Computed:    true,
			},
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

func resourceScalewayVPCPublicGatewayDHCPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcgwAPI, zone, err := vpcgwAPIWithZone(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	req := &vpcgw.CreateDHCPRequest{
		Zone:               zone,
		ProjectID:          d.Get("project_id").(string),
		Subnet:             expandIPNet(d.Get("subnet").(string)),
		EnableDynamic:      expandBoolPtr(d.Get("enable_dynamic")),
		PushDefaultRoute:   expandBoolPtr(d.Get("push_default_route")),
		PushDNSServer:      expandBoolPtr(d.Get("push_dns_servers")),
		DNSServersOverride: expandStringsPtr(d.Get("dns_servers_override")),
		DNSSearch:          expandStringsPtr(d.Get("dns_search")),
		DNSLocalName:       expandStringPtr(d.Get("dns_local_name")),
	}

	if address, ok := d.GetOk("address"); ok {
		req.Address = scw.IPPtr(net.ParseIP(address.(string)))
	}

	if renewTimer, ok := d.GetOk("renew_timer"); ok {
		req.RenewTimer = &scw.Duration{Seconds: renewTimer.(int64)}
	}

	if validLifetime, ok := d.GetOk("valid_lifetime"); ok {
		req.ValidLifetime = &scw.Duration{Seconds: validLifetime.(int64)}
	}

	if rebindTimer, ok := d.GetOk("rebind_timer"); ok {
		req.RebindTimer = &scw.Duration{Seconds: rebindTimer.(int64)}
	}

	if poolLow, ok := d.GetOk("pool_low"); ok {
		req.PoolLow = scw.IPPtr(net.ParseIP(poolLow.(string)))
	}

	if poolHigh, ok := d.GetOk("pool_high"); ok {
		req.PoolLow = scw.IPPtr(net.ParseIP(poolHigh.(string)))
	}

	res, err := vpcgwAPI.CreateDHCP(req, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newZonedIDString(zone, res.ID))

	return resourceScalewayVPCPublicGatewayDHCPRead(ctx, d, meta)
}

func resourceScalewayVPCPublicGatewayDHCPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcgwAPI, zone, ID, err := vpcgwAPIWithZoneAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	dhcp, err := vpcgwAPI.GetDHCP(&vpcgw.GetDHCPRequest{
		DHCPID: ID,
		Zone:   zone,
	}, scw.WithContext(ctx))
	if err != nil {
		if is404Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	_ = d.Set("address", dhcp.Address.String())
	_ = d.Set("created_at", dhcp.CreatedAt.Format(time.RFC3339))
	_ = d.Set("dns_local_name", dhcp.DNSLocalName)
	_ = d.Set("dns_search", dhcp.DNSSearch)
	_ = d.Set("dns_server_override", dhcp.DNSServersOverride)
	_ = d.Set("enable_dynamic", dhcp.EnableDynamic)
	_ = d.Set("organization_id", dhcp.OrganizationID)
	_ = d.Set("pool_high", dhcp.PoolLow.String())
	_ = d.Set("pool_low", dhcp.PoolLow.String())
	_ = d.Set("project_id", dhcp.ProjectID)
	_ = d.Set("push_default_route", dhcp.PushDefaultRoute)
	_ = d.Set("push_dns_server", dhcp.PushDNSServer)
	_ = d.Set("rebind_timer", dhcp.RebindTimer.Seconds)
	_ = d.Set("renew_timer", dhcp.RenewTimer.Seconds)
	_ = d.Set("subnet", dhcp.Subnet.String())
	_ = d.Set("updated_at", dhcp.UpdatedAt.Format(time.RFC3339))
	_ = d.Set("valid_lifetime", dhcp.ValidLifetime.Seconds)
	_ = d.Set("zone", zone)

	return nil
}

func resourceScalewayVPCPublicGatewayDHCPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcgwAPI, zone, ID, err := vpcgwAPIWithZoneAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges("enable_dynamic", "push_default_route", "push_dns_servers", "dns_servers_override", "dns_search", "dns_local_name") {
		req := &vpcgw.UpdateDHCPRequest{
			DHCPID:             ID,
			Zone:               zone,
			EnableDynamic:      expandBoolPtr(d.Get("enable_dynamic")),
			PushDefaultRoute:   expandBoolPtr(d.Get("push_default_route")),
			PushDNSServer:      expandBoolPtr(d.Get("push_dns_servers")),
			DNSServersOverride: expandStringsPtr(d.Get("dns_servers_override")),
			DNSSearch:          expandStringsPtr(d.Get("dns_search")),
			DNSLocalName:       expandStringPtr(d.Get("dns_local_name")),
		}

		_, err = vpcgwAPI.UpdateDHCP(req, scw.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceScalewayVPCPublicGatewayRead(ctx, d, meta)
}

func resourceScalewayVPCPublicGatewayDHCPDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcgwAPI, zone, ID, err := vpcgwAPIWithZoneAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = vpcgwAPI.DeleteDHCP(&vpcgw.DeleteDHCPRequest{
		DHCPID: ID,
		Zone:   zone,
	}, scw.WithContext(ctx))

	if err != nil && !is404Error(err) {
		return diag.FromErr(err)
	}

	return nil
}
