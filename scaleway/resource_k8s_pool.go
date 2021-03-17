package scaleway

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func resourceScalewayK8SPool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScalewayK8SPoolCreate,
		ReadContext:   resourceScalewayK8SPoolRead,
		UpdateContext: resourceScalewayK8SPoolUpdate,
		DeleteContext: resourceScalewayK8SPoolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(defaultK8SPoolTimeout),
		},
		SchemaVersion: 0,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the cluster on which this pool will be created",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the cluster",
			},
			"node_type": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "Server type of the pool servers",
				DiffSuppressFunc: diffSuppressFuncIgnoreCaseAndHyphen,
			},
			"autoscaling": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable the autoscaling on the pool",
			},
			"autohealing": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable the autohealing on the pool",
			},
			"size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Size of the pool",
			},
			"min_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Minimun size of the pool",
			},
			"max_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Maximum size of the pool",
			},
			"tags": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "The tags associated with the pool",
			},
			"container_runtime": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     k8s.RuntimeDocker.String(),
				ForceNew:    true,
				Description: "Container runtime for the pool",
				ValidateFunc: validation.StringInSlice([]string{
					k8s.RuntimeDocker.String(),
					k8s.RuntimeContainerd.String(),
					k8s.RuntimeCrio.String(),
				}, false),
			},
			"wait_for_pool_ready": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to wait for the pool to be ready",
			},
			"placement_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     nil,
				Description: "ID of the placement group",
			},
			"kubelet_args": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "The Kubelet arguments to be used by this pool",
			},
			"upgrade_policy": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "The Pool upgrade policy",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_unavailable": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1,
							Description: "The maximum number of nodes that can be not ready at the same time",
						},
						"max_surge": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: "The maximum number of nodes to be created during the upgrade",
						},
					},
				},
			},
			"zone":   zoneSchema(),
			"region": regionSchema(),
			// Computed elements
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of the creation of the pool",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of the last update of the pool",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Kubernetes version of the pool",
			},
			"current_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The actual size of the pool",
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the node",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the node",
						},
						"public_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The public IPv4 address of the node",
						},
						"public_ip_v6": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The public IPv6 address of the node",
						},
					},
				},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the pool",
			},
		},
	}
}

func resourceScalewayK8SPoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	k8sAPI, region, err := k8sAPIWithRegion(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	////
	// Create pool
	////

	req := &k8s.CreatePoolRequest{
		Region:      region,
		ClusterID:   expandID(d.Get("cluster_id")),
		Name:        expandOrGenerateString(d.Get("name"), "pool"),
		NodeType:    d.Get("node_type").(string),
		Autoscaling: d.Get("autoscaling").(bool),
		Autohealing: d.Get("autohealing").(bool),
		Size:        uint32(d.Get("size").(int)),
		Tags:        expandStrings(d.Get("tags")),
		Zone:        scw.Zone(d.Get("zone").(string)),
		KubeletArgs: expandKubeletArgs(d.Get("kubelet_args")),
	}

	if definedRegion, ok := d.GetOk("region"); ok {
		region = scw.Region(definedRegion.(string))
		req.Region = region
	}

	if placementGroupID, ok := d.GetOk("placement_group_id"); ok {
		req.PlacementGroupID = expandStringPtr(expandID(placementGroupID))
	}

	if minSize, ok := d.GetOk("min_size"); ok {
		req.MinSize = scw.Uint32Ptr(uint32(minSize.(int)))
	}

	if maxSize, ok := d.GetOk("max_size"); ok {
		req.MaxSize = scw.Uint32Ptr(uint32(maxSize.(int)))
	} else {
		req.MaxSize = scw.Uint32Ptr(req.Size)
	}

	if containerRuntime, ok := d.GetOk("container_runtime"); ok {
		req.ContainerRuntime = k8s.Runtime(containerRuntime.(string))
	}

	upgradePolicyReq := &k8s.CreatePoolRequestUpgradePolicy{}

	if maxSurge, ok := d.GetOk("upgrade_policy.0.max_surge"); ok {
		req.UpgradePolicy = upgradePolicyReq
		upgradePolicyReq.MaxSurge = scw.Uint32Ptr(uint32(maxSurge.(int)))
	}

	if maxUnavailable, ok := d.GetOk("upgrade_policy.0.max_unavailable"); ok {
		req.UpgradePolicy = upgradePolicyReq
		upgradePolicyReq.MaxUnavailable = scw.Uint32Ptr(uint32(maxUnavailable.(int)))
	}

	// check if the cluster is waiting for a pool
	cluster, err := k8sAPI.GetCluster(&k8s.GetClusterRequest{
		ClusterID: expandID(d.Get("cluster_id")),
		Region:    region,
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	waitForCluster := false

	if cluster.Status == k8s.ClusterStatusPoolRequired {
		waitForCluster = true
	} else if cluster.Status == k8s.ClusterStatusCreating {
		err = waitK8SCluster(ctx, k8sAPI, region, cluster.ID, k8s.ClusterStatusReady)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	res, err := k8sAPI.CreatePool(req, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newRegionalIDString(region, res.ID))

	if waitForCluster {
		err = waitK8SCluster(ctx, k8sAPI, region, cluster.ID, k8s.ClusterStatusReady)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.Get("wait_for_pool_ready").(bool) { // wait for the pool to be ready if specified (including all its nodes)
		err = waitK8SPoolReady(ctx, k8sAPI, region, res.ID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceScalewayK8SPoolRead(ctx, d, meta)
}

func resourceScalewayK8SPoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	k8sAPI, region, poolID, err := k8sAPIWithRegionAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	////
	// Read Pool
	////
	pool, err := k8sAPI.GetPool(&k8s.GetPoolRequest{
		Region: region,
		PoolID: poolID,
	}, scw.WithContext(ctx))
	if err != nil {
		if is404Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	nodes, err := getNodes(ctx, k8sAPI, pool)
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("cluster_id", newRegionalIDString(region, pool.ClusterID))
	_ = d.Set("name", pool.Name)
	_ = d.Set("node_type", pool.NodeType)
	_ = d.Set("autoscaling", pool.Autoscaling)
	_ = d.Set("autohealing", pool.Autohealing)
	_ = d.Set("current_size", int(pool.Size))
	if !pool.Autoscaling {
		_ = d.Set("size", int(pool.Size))
	}
	_ = d.Set("version", pool.Version)
	_ = d.Set("min_size", int(pool.MinSize))
	_ = d.Set("max_size", int(pool.MaxSize))
	_ = d.Set("tags", pool.Tags)
	_ = d.Set("container_runtime", pool.ContainerRuntime)
	_ = d.Set("created_at", pool.CreatedAt.Format(time.RFC3339))
	_ = d.Set("updated_at", pool.UpdatedAt.Format(time.RFC3339))
	_ = d.Set("nodes", nodes)
	_ = d.Set("status", pool.Status)
	_ = d.Set("kubelet_args", flattenKubeletArgs(pool.KubeletArgs))
	_ = d.Set("zone", pool.Zone)
	_ = d.Set("upgrade_policy", poolUpgradePolicyFlatten(pool))

	if pool.PlacementGroupID != nil {
		zone := scw.Zone(region + "-1") // Placement groups are zoned resources.
		_ = d.Set("placement_group_id", newZonedID(zone, *pool.PlacementGroupID).String())
	}

	return nil
}

func resourceScalewayK8SPoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	k8sAPI, region, poolID, err := k8sAPIWithRegionAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	////
	// Update Pool
	////
	updateRequest := &k8s.UpdatePoolRequest{
		Region: region,
		PoolID: poolID,
	}

	if d.HasChange("autoscaling") {
		updateRequest.Autoscaling = scw.BoolPtr(d.Get("autoscaling").(bool))
	}

	if d.HasChange("autohealing") {
		updateRequest.Autohealing = scw.BoolPtr(d.Get("autohealing").(bool))
	}

	if d.HasChange("min_size") {
		updateRequest.MinSize = scw.Uint32Ptr(uint32(d.Get("min_size").(int)))
	}

	if d.HasChange("max_size") {
		updateRequest.MaxSize = scw.Uint32Ptr(uint32(d.Get("max_size").(int)))
	}

	if !d.Get("autoscaling").(bool) && d.HasChange("size") {
		updateRequest.Size = scw.Uint32Ptr(uint32(d.Get("size").(int)))
	}

	if d.HasChange("tags") {
		updateRequest.Tags = scw.StringsPtr(expandStrings(d.Get("tags")))
	}

	if d.HasChange("kubelet_args") {
		kubeletArgs := expandKubeletArgs(d.Get("kubelet_args"))
		updateRequest.KubeletArgs = &kubeletArgs
	}

	upgradePolicyReq := &k8s.UpdatePoolRequestUpgradePolicy{}

	if d.HasChange("upgrade_policy.0.max_surge") {
		upgradePolicyReq.MaxSurge = scw.Uint32Ptr(uint32(d.Get("upgrade_policy.0.max_surge").(int)))
	}

	if d.HasChange("upgrade_policy.0.max_unavailable") {
		upgradePolicyReq.MaxUnavailable = scw.Uint32Ptr(uint32(d.Get("upgrade_policy.0.max_unavailable").(int)))
	}

	updateRequest.UpgradePolicy = upgradePolicyReq

	res, err := k8sAPI.UpdatePool(updateRequest, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	if d.Get("wait_for_pool_ready").(bool) { // wait for the pool to be ready if specified (including all its nodes)
		err = waitK8SPoolReady(ctx, k8sAPI, region, res.ID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceScalewayK8SPoolRead(ctx, d, meta)
}

func resourceScalewayK8SPoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	k8sAPI, region, poolID, err := k8sAPIWithRegionAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	////
	// Delete Pool
	////
	_, err = k8sAPI.DeletePool(&k8s.DeletePoolRequest{
		Region: region,
		PoolID: poolID,
	}, scw.WithContext(ctx))
	if err != nil {
		if !is404Error(err) {
			return diag.FromErr(err)
		}
	}

	return nil
}
