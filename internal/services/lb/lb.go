package lb

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ipamAPI "github.com/scaleway/scaleway-sdk-go/api/ipam/v1"
	lbSDK "github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/cdf"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/dsf"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/regional"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/zonal"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/account"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/ipam"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/types"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/verify"
)

const (
	DefaultWaitLBRetryInterval = 30 * time.Second
)

func ResourceLb() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLbCreate,
		ReadContext:   resourceLbRead,
		UpdateContext: resourceLbUpdate,
		DeleteContext: resourceLbDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create:  schema.DefaultTimeout(defaultLbLbTimeout),
			Read:    schema.DefaultTimeout(defaultLbLbTimeout),
			Update:  schema.DefaultTimeout(defaultLbLbTimeout),
			Delete:  schema.DefaultTimeout(defaultLbLbTimeout),
			Default: schema.DefaultTimeout(defaultLbLbTimeout),
		},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{Version: 0, Type: lbUpgradeV1SchemaType(), Upgrade: UpgradeStateV1Func},
		},
		CustomizeDiff: customdiff.All(
			cdf.LocalityCheck("ip_id", "private_network.#.private_network_id"),
			customizeDiffLBIPIDs,
			customizeDiffAssignFlexibleIPv6,
		),
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the lb",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the lb",
			},
			"type": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: dsf.IgnoreCase,
				Description:      "The type of load-balancer you want to create",
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Array of tags to associate with the load-balancer",
			},
			"ip_id": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "The load-balance public IP ID",
				DiffSuppressFunc: dsf.Locality,
				ValidateDiagFunc: verify.IsUUIDorUUIDWithLocality(),
				Deprecated:       "Please use ip_ids",
			},
			"ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The load-balance public IPv4 address",
			},
			"ipv6_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The load-balance public IPv6 address",
			},
			"release_ip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Release the IPs related to this load-balancer",
				Deprecated:  "The resource ip will be destroyed by it's own resource. Please set this to `false`",
			},
			"private_network": {
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    8,
				Set:         lbPrivateNetworkSetHash,
				Description: "List of private network to connect with your load balancer",
				DiffSuppressFunc: func(k, oldValue, newValue string, _ *schema.ResourceData) bool {
					// Check if the key is for the 'private_network_id' attribute
					if strings.HasSuffix(k, "private_network_id") {
						return locality.ExpandID(oldValue) == locality.ExpandID(newValue)
					}
					// For all other attributes, don't suppress the diff
					return false
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"private_network_id": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: verify.IsUUIDorUUIDWithLocality(),
							Description:      "The Private Network ID",
						},
						"static_config": {
							Description: "Define an IP address in the subnet of your private network that will be assigned to your load balancer instance",
							Type:        schema.TypeList,
							Optional:    true,
							Deprecated:  "static_config field is deprecated, please use `private_network_id` or `ipam_ids` instead",
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: verify.IsStandaloneIPorCIDR(),
							},
							MaxItems: 1,
						},
						"dhcp_config": {
							Description: "Set to true if you want to let DHCP assign IP addresses",
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Deprecated:  "dhcp_config field is deprecated, please use `private_network_id` or `ipam_ids` instead",
						},
						"ipam_ids": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "IPAM ID of a pre-reserved IP address to assign to the Load Balancer on this Private Network",
						},
						// Readonly attributes
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of private network connection",
						},
						"zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"ssl_compatibility_level": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Enforces minimal SSL version (in SSL/TLS offloading context)",
				Default:          lbSDK.SSLCompatibilityLevelSslCompatibilityLevelIntermediate.String(),
				ValidateDiagFunc: verify.ValidateEnum[lbSDK.SSLCompatibilityLevel](),
			},
			"assign_flexible_ip": {
				Type:          schema.TypeBool,
				Optional:      true,
				ForceNew:      true,
				Description:   "Defines whether to automatically assign a flexible public IP to the load balancer",
				ConflictsWith: []string{"ip_ids"},
			},
			"assign_flexible_ipv6": {
				Type:          schema.TypeBool,
				Optional:      true,
				Description:   "Defines whether to automatically assign a flexible public IPv6 to the load balancer",
				ConflictsWith: []string{"ip_ids"},
			},
			"ip_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: verify.IsUUIDorUUIDWithLocality(),
				},
				Description:      "List of IP IDs to attach to the Load Balancer",
				DiffSuppressFunc: dsf.OrderDiff,
				ConflictsWith:    []string{"assign_flexible_ip", "assign_flexible_ipv6"},
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
			"region":          regional.ComputedSchema(),
			"zone":            zonal.Schema(),
			"organization_id": account.OrganizationIDSchema(),
			"project_id":      account.ProjectIDSchema(),
		},
	}
}

func resourceLbCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	lbAPI, zone, err := lbAPIWithZone(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	createReq := &lbSDK.ZonedAPICreateLBRequest{
		Zone:                  zone,
		IPIDs:                 types.ExpandSliceIDs(d.Get("ip_ids")),
		IPID:                  types.ExpandStringPtr(locality.ExpandID(d.Get("ip_id"))),
		ProjectID:             types.ExpandStringPtr(d.Get("project_id")),
		Name:                  types.ExpandOrGenerateString(d.Get("name"), "lb"),
		Description:           d.Get("description").(string),
		Type:                  d.Get("type").(string),
		SslCompatibilityLevel: lbSDK.SSLCompatibilityLevel(*types.ExpandStringPtr(d.Get("ssl_compatibility_level"))),
		AssignFlexibleIP:      types.ExpandBoolPtr(types.GetBool(d, "assign_flexible_ip")),
		AssignFlexibleIPv6:    types.ExpandBoolPtr(types.GetBool(d, "assign_flexible_ipv6")),
	}

	if tags, ok := d.GetOk("tags"); ok {
		for _, tag := range tags.([]interface{}) {
			createReq.Tags = append(createReq.Tags, tag.(string))
		}
	}

	lb, err := lbAPI.CreateLB(createReq, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(zonal.NewIDString(zone, lb.ID))

	// check err waiting process
	_, err = waitForLB(ctx, lbAPI, zone, lb.ID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	// attach private network
	pnConfigs, pnExist := d.GetOk("private_network")
	if pnExist {
		pnConfigs, err := expandPrivateNetworks(pnConfigs)
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = attachLBPrivateNetworks(ctx, lbAPI, zone, pnConfigs, lb.ID, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = waitForLB(ctx, lbAPI, zone, lb.ID, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceLbRead(ctx, d, m)
}

func resourceLbRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	lbAPI, zone, ID, err := NewAPIWithZoneAndID(m, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	lb, err := waitForInstances(ctx, lbAPI, zone, ID, d.Timeout(schema.TimeoutRead))
	if err != nil {
		if httperrors.Is404(err) || httperrors.Is403(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	// set the region from zone
	region, err := zone.Region()
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("release_ip", false)
	_ = d.Set("name", lb.Name)
	_ = d.Set("description", lb.Description)
	_ = d.Set("zone", lb.Zone.String())
	_ = d.Set("region", region.String())
	_ = d.Set("organization_id", lb.OrganizationID)
	_ = d.Set("project_id", lb.ProjectID)
	_ = d.Set("tags", lb.Tags)
	// For now API return lowercase lb type. This should be fixed in a near future on the API side
	_ = d.Set("type", strings.ToUpper(lb.Type))
	_ = d.Set("ssl_compatibility_level", lb.SslCompatibilityLevel.String())
	if len(lb.IP) > 0 {
		_ = d.Set("ip_id", zonal.NewIDString(zone, lb.IP[0].ID))
		_ = d.Set("ip_ids", flattenLBIPIDs(zone, lb.IP))
		var ipv4Address, ipv6Address string
		for _, ip := range lb.IP {
			parsedIP := net.ParseIP(ip.IPAddress)
			if parsedIP != nil {
				if parsedIP.To4() != nil {
					ipv4Address = ip.IPAddress
				} else {
					ipv6Address = ip.IPAddress
				}
			}
		}
		_ = d.Set("ip_address", ipv4Address)
		_ = d.Set("ipv6_address", ipv6Address)
	}

	// retrieve attached private networks
	privateNetworks, err := waitForPrivateNetworks(ctx, lbAPI, zone, ID, d.Timeout(schema.TimeoutRead))
	if err != nil {
		if httperrors.Is404(err) {
			return nil
		}
		return diag.FromErr(err)
	}
	_ = d.Set("private_network", flattenPrivateNetworkConfigs(privateNetworks))

	privateNetworkIDs := make([]string, 0, len(privateNetworks))
	for _, pn := range privateNetworks {
		privateNetworkIDs = append(privateNetworkIDs, pn.PrivateNetworkID)
	}

	var allPrivateIPs []map[string]interface{}
	resourceType := ipamAPI.ResourceTypeLBServer
	for _, privateNetworkID := range privateNetworkIDs {
		opts := &ipam.GetResourcePrivateIPsOptions{
			ResourceType:     &resourceType,
			PrivateNetworkID: &privateNetworkID,
		}
		privateIPs, err := ipam.GetResourcePrivateIPs(ctx, m, region, opts)
		if err != nil {
			return diag.FromErr(err)
		}
		if privateIPs != nil {
			allPrivateIPs = append(allPrivateIPs, privateIPs...)
		}
	}

	_ = d.Set("private_ip", allPrivateIPs)

	return nil
}

//gocyclo:ignore
func resourceLbUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	lbAPI, zone, ID, err := NewAPIWithZoneAndID(m, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	req := &lbSDK.ZonedAPIUpdateLBRequest{
		Zone:                  zone,
		LBID:                  ID,
		Name:                  d.Get("name").(string),
		Tags:                  types.ExpandStrings(d.Get("tags")),
		Description:           d.Get("description").(string),
		SslCompatibilityLevel: lbSDK.SSLCompatibilityLevel(*types.ExpandStringPtr(d.Get("ssl_compatibility_level"))),
	}

	_, err = lbAPI.UpdateLB(req, scw.WithContext(ctx))
	if err != nil && !httperrors.Is404(err) {
		return diag.FromErr(err)
	}

	_, err = waitForLB(ctx, lbAPI, zone, ID, d.Timeout(schema.TimeoutUpdate))
	if err != nil && !httperrors.Is404(err) {
		return diag.FromErr(err)
	}

	if d.HasChange("type") {
		_, err = waitForInstances(ctx, lbAPI, zone, ID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}

		lbType := d.Get("type").(string)
		migrateReq := &lbSDK.ZonedAPIMigrateLBRequest{
			Zone: zone,
			LBID: ID,
			Type: lbType,
		}

		lb, err := lbAPI.MigrateLB(migrateReq, scw.WithContext(ctx))
		if err != nil {
			return diag.FromErr(fmt.Errorf("couldn't migrate load balancer on type %s: %w", migrateReq.Type, err))
		}

		_, err = waitForLB(ctx, lbAPI, zone, lb.ID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("ip_ids") {
		oldIPIDs, newIPIDs := d.GetChange("ip_ids")
		oldIPIDsSet := make(map[string]struct{})
		newIPIDsSet := make(map[string]struct{})

		for _, id := range oldIPIDs.([]interface{}) {
			oldIPIDsSet[id.(string)] = struct{}{}
		}

		for _, id := range newIPIDs.([]interface{}) {
			newIPIDsSet[id.(string)] = struct{}{}
		}

		// Check if an IPv6 address is being added to an existing LB with an IPv4 address
		if len(oldIPIDsSet) == 1 && len(newIPIDsSet) == 2 {
			var ipv4ID, ipv6ID string
			for id := range oldIPIDsSet {
				ipv4ID = id
			}
			for id := range newIPIDsSet {
				if id != ipv4ID {
					ipv6ID = id
					break
				}
			}

			res, err := lbAPI.GetIP(&lbSDK.ZonedAPIGetIPRequest{
				Zone: zone,
				IPID: locality.ExpandID(ipv6ID),
			}, scw.WithContext(ctx))
			if err != nil {
				return diag.FromErr(err)
			}

			_, err = lbAPI.UpdateIP(&lbSDK.ZonedAPIUpdateIPRequest{
				Zone: zone,
				IPID: res.ID,
				LBID: types.ExpandStringPtr(ID),
			}, scw.WithContext(ctx))
			if err != nil {
				return diag.FromErr(err)
			}

			_, err = waitForLB(ctx, lbAPI, zone, ID, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("assign_flexible_ipv6") {
		assignFlexibleIPv6 := d.Get("assign_flexible_ipv6").(bool)
		if assignFlexibleIPv6 {
			createReq := &lbSDK.ZonedAPICreateIPRequest{
				Zone:      zone,
				ProjectID: types.ExpandStringPtr(d.Get("project_id")),
				IsIPv6:    true,
			}

			res, err := lbAPI.CreateIP(createReq, scw.WithContext(ctx))
			if err != nil {
				return diag.FromErr(err)
			}

			updateReq := &lbSDK.ZonedAPIUpdateIPRequest{
				Zone: zone,
				IPID: res.ID,
				LBID: types.ExpandStringPtr(ID),
			}

			_, err = lbAPI.UpdateIP(updateReq, scw.WithContext(ctx))
			if err != nil {
				return diag.FromErr(err)
			}

			_, err = waitForLB(ctx, lbAPI, zone, ID, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	////
	// Attach / Detach Private Networks
	////
	if d.HasChange("private_network") {
		// check current lb stability state
		_, err = waitForLB(ctx, lbAPI, zone, ID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}

		// check that pns are in a stable state
		_, err := waitForPrivateNetworks(ctx, lbAPI, zone, ID, d.Timeout(schema.TimeoutUpdate))
		if err != nil && !httperrors.Is404(err) {
			return diag.FromErr(err)
		}

		oldPNs, newPNs := d.GetChange("private_network")
		oldPNConfigs, err := expandPrivateNetworks(oldPNs)
		if err != nil {
			return diag.FromErr(err)
		}
		newPNConfigs, err := expandPrivateNetworks(newPNs)
		if err != nil {
			return diag.FromErr(err)
		}

		toDetach, toAttach := PrivateNetworksCompare(oldPNConfigs, newPNConfigs)

		// detach private networks
		for _, pn := range toDetach {
			err = lbAPI.DetachPrivateNetwork(&lbSDK.ZonedAPIDetachPrivateNetworkRequest{
				Zone:             zone,
				LBID:             ID,
				PrivateNetworkID: pn.PrivateNetworkID,
			}, scw.WithContext(ctx))
			if err != nil && !httperrors.Is404(err) {
				return diag.FromErr(err)
			}
		}

		_, err = waitForLB(ctx, lbAPI, zone, ID, d.Timeout(schema.TimeoutUpdate))
		if err != nil && !httperrors.Is404(err) {
			return diag.FromErr(err)
		}

		_, err = waitForPrivateNetworks(ctx, lbAPI, zone, ID, d.Timeout(schema.TimeoutUpdate))
		if err != nil && !httperrors.Is404(err) {
			return diag.FromErr(err)
		}

		// attach private networks
		_, err = attachLBPrivateNetworks(ctx, lbAPI, zone, toAttach, ID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}

		privateNetworks, err := waitForPrivateNetworks(ctx, lbAPI, zone, ID, d.Timeout(schema.TimeoutUpdate))
		if err != nil && !httperrors.Is404(err) {
			return diag.FromErr(err)
		}

		for _, pn := range privateNetworks {
			tflog.Debug(ctx, fmt.Sprintf("PrivateNetwork ID %s state: %v", pn.PrivateNetworkID, pn.Status))
			if pn.Status == lbSDK.PrivateNetworkStatusError {
				err = lbAPI.DetachPrivateNetwork(&lbSDK.ZonedAPIDetachPrivateNetworkRequest{
					Zone:             zone,
					LBID:             ID,
					PrivateNetworkID: pn.PrivateNetworkID,
				}, scw.WithContext(ctx))
				if err != nil && !httperrors.Is404(err) {
					return diag.FromErr(err)
				}
				return diag.Errorf("attaching private network with id: %s on error state. please check your config", pn.PrivateNetworkID)
			}
		}
	}

	return resourceLbRead(ctx, d, m)
}

func resourceLbDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	lbAPI, zone, ID, err := NewAPIWithZoneAndID(m, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// check if current lb is on stable state
	currentLB, err := waitForInstances(ctx, lbAPI, zone, ID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	if currentLB.PrivateNetworkCount != 0 {
		lbPNs, err := lbAPI.ListLBPrivateNetworks(&lbSDK.ZonedAPIListLBPrivateNetworksRequest{
			Zone: zone,
			LBID: ID,
		}, scw.WithContext(ctx))
		if err != nil && !httperrors.Is404(err) {
			return diag.FromErr(err)
		}

		// detach private networks
		for _, pn := range lbPNs.PrivateNetwork {
			err = lbAPI.DetachPrivateNetwork(&lbSDK.ZonedAPIDetachPrivateNetworkRequest{
				Zone:             zone,
				LBID:             ID,
				PrivateNetworkID: pn.PrivateNetworkID,
			}, scw.WithContext(ctx))
			if err != nil && !httperrors.Is404(err) {
				return diag.FromErr(err)
			}
		}

		_, err = waitForInstances(ctx, lbAPI, zone, ID, d.Timeout(schema.TimeoutDelete))
		if err != nil && !httperrors.Is404(err) {
			return diag.FromErr(err)
		}
	}

	err = lbAPI.DeleteLB(&lbSDK.ZonedAPIDeleteLBRequest{
		Zone:      zone,
		LBID:      ID,
		ReleaseIP: false,
	}, scw.WithContext(ctx))
	if err != nil && !httperrors.Is404(err) {
		return diag.FromErr(err)
	}

	_, err = waitForLB(ctx, lbAPI, zone, ID, d.Timeout(schema.TimeoutDelete))
	if err != nil && !httperrors.Is404(err) {
		return diag.FromErr(err)
	}

	_, err = waitForInstances(ctx, lbAPI, zone, ID, d.Timeout(schema.TimeoutDelete))
	if err != nil && !httperrors.Is404(err) {
		return diag.FromErr(err)
	}

	return nil
}
