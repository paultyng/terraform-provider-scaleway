package redis

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ipamAPI "github.com/scaleway/scaleway-sdk-go/api/ipam/v1"
	"github.com/scaleway/scaleway-sdk-go/api/redis/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/cdf"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/dsf"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/zonal"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/account"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/ipam"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/types"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/verify"
)

func ResourceCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceClusterCreate,
		ReadContext:   ResourceClusterRead,
		UpdateContext: ResourceClusterUpdate,
		DeleteContext: ResourceClusterDelete,
		Timeouts: &schema.ResourceTimeout{
			Create:  schema.DefaultTimeout(defaultRedisClusterTimeout),
			Update:  schema.DefaultTimeout(defaultRedisClusterTimeout),
			Delete:  schema.DefaultTimeout(defaultRedisClusterTimeout),
			Default: schema.DefaultTimeout(defaultRedisClusterTimeout),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 0,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the redis cluster",
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Redis version of the cluster",
			},
			"node_type": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Type of node to use for the cluster",
				DiffSuppressFunc: dsf.IgnoreCase,
			},
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the user created when the cluster is created",
			},
			"password": {
				Type:        schema.TypeString,
				Sensitive:   true,
				Required:    true,
				Description: "Password of the user",
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of tags [\"tag1\", \"tag2\", ...] attached to a redis cluster",
			},
			"cluster_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Number of nodes for the cluster.",
			},
			"tls_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether or not TLS is enabled.",
				ForceNew:    true,
			},
			"acl": {
				Type:          schema.TypeSet,
				Description:   "List of acl rules.",
				Optional:      true,
				ConflictsWith: []string{"private_network"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "ID of the rule (UUID format).",
							Computed:    true,
						},
						"ip": {
							Type:         schema.TypeString,
							Description:  "IPv4 network address of the rule (IP network in a CIDR format).",
							Required:     true,
							ValidateFunc: validation.IsCIDR,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "Description of the rule.",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"settings": {
				Type:        schema.TypeMap,
				Description: "Map of settings to define for the cluster.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"private_network": {
				Type:          schema.TypeSet,
				Optional:      true,
				Description:   "Private network specs details",
				ConflictsWith: []string{"acl"},
				Set:           privateNetworkSetHash,
				DiffSuppressFunc: func(k, oldValue, newValue string, _ *schema.ResourceData) bool {
					// Check if the key is for the 'id' attribute
					if strings.HasSuffix(k, "id") {
						return locality.ExpandID(oldValue) == locality.ExpandID(newValue)
					}
					// For all other attributes, don't suppress the diff
					return false
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: verify.IsUUIDorUUIDWithLocality(),
							Description:      "UUID of the private network to be connected to the cluster",
						},
						"service_ips": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.IsCIDR,
							},
							Description: "List of IPv4 addresses of the private network with a CIDR notation",
						},
						// computed
						"endpoint_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "UUID of the endpoint to be connected to the cluster",
						},
						"zone": zonal.ComputedSchema(),
					},
				},
			},
			// Computed
			"public_network": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Public network specs details",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "TCP port of the endpoint",
						},
						"ips": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed: true,
						},
					},
				},
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
			"certificate": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "public TLS certificate used by redis cluster, empty if tls is disabled",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of the creation of the Redis cluster",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of the last update of the Redis cluster",
			},
			// Common
			"zone":       zonal.Schema(),
			"project_id": account.ProjectIDSchema(),
		},
		CustomizeDiff: customdiff.All(
			cdf.LocalityCheck("private_network.#.id"),
			customizeDiffMigrateClusterSize(),
		),
	}
}

func customizeDiffMigrateClusterSize() schema.CustomizeDiffFunc {
	return func(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
		oldSize, newSize := diff.GetChange("cluster_size")
		if newSize == 2 {
			return errors.New("cluster_size can be either 1 (standalone) ou >3 (cluster mode), not 2")
		}
		if oldSize == 1 && newSize != 1 || newSize.(int) < oldSize.(int) {
			return diff.ForceNew("cluster_size")
		}
		return nil
	}
}

func ResourceClusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	redisAPI, zone, err := newAPIWithZone(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	createReq := &redis.CreateClusterRequest{
		Zone:      zone,
		ProjectID: d.Get("project_id").(string),
		Name:      types.ExpandOrGenerateString(d.Get("name"), "redis"),
		Version:   d.Get("version").(string),
		NodeType:  d.Get("node_type").(string),
		UserName:  d.Get("user_name").(string),
		Password:  d.Get("password").(string),
	}

	tags, tagsExist := d.GetOk("tags")
	if tagsExist {
		createReq.Tags = types.ExpandStrings(tags)
	}
	clusterSize, clusterSizeExist := d.GetOk("cluster_size")
	if clusterSizeExist {
		createReq.ClusterSize = scw.Int32Ptr(int32(clusterSize.(int)))
	}
	tlsEnabled, tlsEnabledExist := d.GetOk("tls_enabled")
	if tlsEnabledExist {
		createReq.TLSEnabled = tlsEnabled.(bool)
	}
	aclRules, aclRulesExist := d.GetOk("acl")
	if aclRulesExist {
		rules, err := expandACLSpecs(aclRules)
		if err != nil {
			return diag.FromErr(err)
		}
		createReq.ACLRules = rules
	}
	settings, settingsExist := d.GetOk("settings")
	if settingsExist {
		createReq.ClusterSettings = expandSettings(settings)
	}

	pn, pnExists := d.GetOk("private_network")
	if pnExists {
		pnSpecs, err := expandPrivateNetwork(pn.(*schema.Set).List())
		if err != nil {
			return diag.FromErr(err)
		}
		createReq.Endpoints = pnSpecs
	}

	res, err := redisAPI.CreateCluster(createReq, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(zonal.NewIDString(zone, res.ID))

	_, err = waitForCluster(ctx, redisAPI, zone, res.ID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return ResourceClusterRead(ctx, d, m)
}

func ResourceClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	redisAPI, zone, ID, err := NewAPIWithZoneAndID(m, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	getReq := &redis.GetClusterRequest{
		Zone:      zone,
		ClusterID: ID,
	}
	cluster, err := redisAPI.GetCluster(getReq, scw.WithContext(ctx))
	if err != nil {
		if httperrors.Is404(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	_ = d.Set("name", cluster.Name)
	_ = d.Set("node_type", cluster.NodeType)
	_ = d.Set("user_name", d.Get("user_name").(string))
	_ = d.Set("password", d.Get("password").(string))
	_ = d.Set("zone", cluster.Zone.String())
	_ = d.Set("project_id", cluster.ProjectID)
	_ = d.Set("version", cluster.Version)
	_ = d.Set("cluster_size", int(cluster.ClusterSize))
	_ = d.Set("created_at", cluster.CreatedAt.Format(time.RFC3339))
	_ = d.Set("updated_at", cluster.UpdatedAt.Format(time.RFC3339))
	_ = d.Set("acl", flattenACLs(cluster.ACLRules))
	_ = d.Set("settings", flattenSettings(cluster.ClusterSettings))

	if len(cluster.Tags) > 0 {
		_ = d.Set("tags", cluster.Tags)
	}

	// set endpoints
	var allPrivateIPs []map[string]interface{}
	pnI, pnExists := flattenPrivateNetwork(cluster.Endpoints)
	if pnExists {
		_ = d.Set("private_network", pnI)

		var privateNetworkIDs []string
		for _, endpoint := range cluster.Endpoints {
			if endpoint.PrivateNetwork != nil {
				privateNetworkIDs = append(privateNetworkIDs, endpoint.PrivateNetwork.ID)
			}
		}

		resourceType := ipamAPI.ResourceTypeRedisCluster
		region, err := zone.Region()
		if err != nil {
			return diag.FromErr(err)
		}
		for _, privateNetworkID := range privateNetworkIDs {
			opts := &ipam.GetResourcePrivateIPsOptions{
				ResourceType:     &resourceType,
				PrivateNetworkID: &privateNetworkID,
				ResourceID:       &cluster.ID,
			}
			privateIPs, err := ipam.GetResourcePrivateIPs(ctx, m, region, opts)
			if err != nil {
				return diag.FromErr(err)
			}
			if privateIPs != nil {
				allPrivateIPs = append(allPrivateIPs, privateIPs...)
			}
		}
	}
	_ = d.Set("private_ip", allPrivateIPs)
	_ = d.Set("public_network", flattenPublicNetwork(cluster.Endpoints))

	if cluster.TLSEnabled {
		certificate, err := redisAPI.GetClusterCertificate(&redis.GetClusterCertificateRequest{
			Zone:      zone,
			ClusterID: cluster.ID,
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch cluster certificate: %w", err))
		}

		certificateContent, err := io.ReadAll(certificate.Content)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to read cluster certificate: %w", err))
		}

		_ = d.Set("certificate", string(certificateContent))
	} else {
		_ = d.Set("certificate", "")
	}

	return nil
}

func ResourceClusterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	redisAPI, zone, ID, err := NewAPIWithZoneAndID(m, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	req := &redis.UpdateClusterRequest{
		Zone:      zone,
		ClusterID: ID,
	}

	if d.HasChange("name") {
		req.Name = types.ExpandStringPtr(d.Get("name"))
	}
	if d.HasChange("user_name") {
		req.UserName = types.ExpandStringPtr(d.Get("user_name"))
	}
	if d.HasChange("password") {
		req.Password = types.ExpandStringPtr(d.Get("password"))
	}
	if d.HasChange("tags") {
		req.Tags = types.ExpandUpdatedStringsPtr(d.Get("tags"))
	}
	if d.HasChange("acl") {
		diagnostics := updateACL(ctx, d, redisAPI, zone, ID)
		if diagnostics != nil {
			return diagnostics
		}
	}
	if d.HasChange("settings") {
		diagnostics := updateSettings(ctx, d, redisAPI, zone, ID)
		if diagnostics != nil {
			return diagnostics
		}
	}

	_, err = waitForCluster(ctx, redisAPI, zone, ID, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = redisAPI.UpdateCluster(req, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	migrateClusterRequests := []redis.MigrateClusterRequest(nil)
	if d.HasChange("cluster_size") {
		migrateClusterRequests = append(migrateClusterRequests, redis.MigrateClusterRequest{
			Zone:        zone,
			ClusterID:   ID,
			ClusterSize: scw.Uint32Ptr(uint32(d.Get("cluster_size").(int))),
		})
	}
	if d.HasChange("version") {
		migrateClusterRequests = append(migrateClusterRequests, redis.MigrateClusterRequest{
			Zone:      zone,
			ClusterID: ID,
			Version:   types.ExpandStringPtr(d.Get("version")),
		})
	}
	if d.HasChange("node_type") {
		migrateClusterRequests = append(migrateClusterRequests, redis.MigrateClusterRequest{
			Zone:      zone,
			ClusterID: ID,
			NodeType:  types.ExpandStringPtr(d.Get("node_type")),
		})
	}
	for i := range migrateClusterRequests {
		_, err = waitForCluster(ctx, redisAPI, zone, ID, d.Timeout(schema.TimeoutUpdate))
		if err != nil && !httperrors.Is404(err) {
			return diag.FromErr(err)
		}
		_, err = redisAPI.MigrateCluster(&migrateClusterRequests[i], scw.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = waitForCluster(ctx, redisAPI, zone, ID, d.Timeout(schema.TimeoutUpdate))
		if err != nil && !httperrors.Is404(err) {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("private_network") {
		diagnostics := ResourceClusterUpdateEndpoints(ctx, d, redisAPI, zone, ID)
		if diagnostics != nil {
			return diagnostics
		}
	}

	_, err = waitForCluster(ctx, redisAPI, zone, ID, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.FromErr(err)
	}

	return ResourceClusterRead(ctx, d, m)
}

func updateACL(ctx context.Context, d *schema.ResourceData, redisAPI *redis.API, zone scw.Zone, clusterID string) diag.Diagnostics {
	rules, err := expandACLSpecs(d.Get("acl"))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = redisAPI.SetACLRules(&redis.SetACLRulesRequest{
		Zone:      zone,
		ClusterID: clusterID,
		ACLRules:  rules,
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func updateSettings(ctx context.Context, d *schema.ResourceData, redisAPI *redis.API, zone scw.Zone, clusterID string) diag.Diagnostics {
	settings := expandSettings(d.Get("settings"))

	_, err := redisAPI.SetClusterSettings(&redis.SetClusterSettingsRequest{
		Zone:      zone,
		ClusterID: clusterID,
		Settings:  settings,
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ResourceClusterUpdateEndpoints(ctx context.Context, d *schema.ResourceData, redisAPI *redis.API, zone scw.Zone, clusterID string) diag.Diagnostics {
	// retrieve state
	cluster, err := waitForCluster(ctx, redisAPI, zone, clusterID, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.FromErr(err)
	}

	// get new desired state of endpoints
	rawNewEndpoints := d.Get("private_network")
	newEndpoints, err := expandPrivateNetwork(rawNewEndpoints.(*schema.Set).List())
	if err != nil {
		return diag.FromErr(err)
	}
	if len(newEndpoints) == 0 {
		newEndpoints = append(newEndpoints, &redis.EndpointSpec{
			PublicNetwork: &redis.EndpointSpecPublicNetworkSpec{},
		})
	}
	// send request
	_, err = redisAPI.SetEndpoints(&redis.SetEndpointsRequest{
		Zone:      cluster.Zone,
		ClusterID: cluster.ID,
		Endpoints: newEndpoints,
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = waitForCluster(ctx, redisAPI, zone, clusterID, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ResourceClusterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	redisAPI, zone, ID, err := NewAPIWithZoneAndID(m, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = waitForCluster(ctx, redisAPI, zone, ID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = redisAPI.DeleteCluster(&redis.DeleteClusterRequest{
		Zone:      zone,
		ClusterID: ID,
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = waitForCluster(ctx, redisAPI, zone, ID, d.Timeout(schema.TimeoutDelete))
	if err != nil && !httperrors.Is404(err) {
		return diag.FromErr(err)
	}

	return nil
}
