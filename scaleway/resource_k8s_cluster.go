package scaleway

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func resourceScalewayK8SCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScalewayK8SClusterCreate,
		ReadContext:   resourceScalewayK8SClusterRead,
		UpdateContext: resourceScalewayK8SClusterUpdate,
		DeleteContext: resourceScalewayK8SClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(defaultK8SClusterTimeout),
		},
		SchemaVersion: 0,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the cluster",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The description of the cluster",
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The version of the cluster",
			},
			"cni": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The CNI plugin of the cluster",
				ValidateFunc: validation.StringInSlice([]string{
					k8s.CNICilium.String(),
					k8s.CNICalico.String(),
					k8s.CNIFlannel.String(),
					k8s.CNIWeave.String(),
				}, false),
			},
			"enable_dashboard": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable the dashboard on the cluster",
			},
			"ingress": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     k8s.IngressNone.String(),
				Description: "The ingress to be deployed on the cluster",
				ValidateFunc: validation.StringInSlice([]string{
					k8s.IngressNone.String(),
					k8s.IngressTraefik.String(),
					k8s.IngressTraefik2.String(),
					k8s.IngressNginx.String(),
				}, false),
			},
			"tags": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "The tags associated with the cluster",
			},
			"autoscaler_config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "The autoscaler configuration for the cluster",
				Elem:        autoscalerConfigSchema(),
			},
			"auto_upgrade": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "The auto upgrade configuration for the cluster",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Enables the Kubernetes patch version auto upgrade",
						},
						"maintenance_window_start_hour": {
							Type:         schema.TypeInt,
							Required:     true,
							Description:  "Start hour of the 2-hour maintenance window",
							ValidateFunc: validateHour(),
						},
						"maintenance_window_day": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Day of the maintenance window",
							ValidateFunc: validation.StringInSlice([]string{
								k8s.MaintenanceWindowDayOfTheWeekAny.String(),
								k8s.MaintenanceWindowDayOfTheWeekMonday.String(),
								k8s.MaintenanceWindowDayOfTheWeekTuesday.String(),
								k8s.MaintenanceWindowDayOfTheWeekWednesday.String(),
								k8s.MaintenanceWindowDayOfTheWeekThursday.String(),
								k8s.MaintenanceWindowDayOfTheWeekFriday.String(),
								k8s.MaintenanceWindowDayOfTheWeekSaturday.String(),
								k8s.MaintenanceWindowDayOfTheWeekSunday.String(),
							}, false),
						},
					},
				},
			},
			"feature_gates": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "The list of feature gates to enable on the cluster",
			},
			"admission_plugins": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "The list of admission plugins to enable on the cluster",
			},
			"open_id_connect_config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "The OpenID Connect configuration of the cluster",
				Elem:        openIDConnectConfigSchema(),
			},
			"apiserver_cert_sans": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Additional Subject Alternative Names for the Kubernetes API server certificate",
			},
			"delete_additional_resources": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Delete additional resources like block volumes and loadbalancers on cluster deletion",
			},
			"region":          regionSchema(),
			"organization_id": organizationIDSchema(),
			"project_id":      projectIDSchema(),
			// Computed elements
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of the creation of the Kubernetes cluster",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of the last update of the Kubernetes cluster",
			},
			"apiserver_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Kubernetes API server URL",
			},
			"wildcard_dns": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Wildcard DNS pointing to all the ready nodes",
			},
			"kubeconfig": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The kubeconfig configuration file of the Kubernetes cluster",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_file": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The whole kubeconfig file",
						},
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The kubernetes master URL",
						},
						"cluster_ca_certificate": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The kubernetes cluster CA certificate",
						},
						"token": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The kubernetes cluster admin token",
						},
					},
				},
			},
			"upgrade_available": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "True if an upgrade is available",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the cluster",
			},
		},
	}
}

func resourceScalewayK8SClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	k8sAPI, region, err := k8sAPIWithRegion(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	////
	// Create cluster
	////

	description, ok := d.GetOk("description")
	if !ok {
		description = ""
	}

	req := &k8s.CreateClusterRequest{
		Region:            region,
		ProjectID:         expandStringPtr(d.Get("project_id")),
		Name:              expandOrGenerateString(d.Get("name"), "cluster"),
		Description:       description.(string),
		Cni:               k8s.CNI(d.Get("cni").(string)),
		Tags:              expandStrings(d.Get("tags")),
		FeatureGates:      expandStrings(d.Get("feature_gates")),
		AdmissionPlugins:  expandStrings(d.Get("admission_plugins")),
		ApiserverCertSans: expandStrings(d.Get("apiserver_cert_sans")),
	}

	if definedRegion, ok := d.GetOk("region"); ok {
		region = scw.Region(definedRegion.(string))
		req.Region = region
	}

	if dashboard, ok := d.GetOk("enable_dashboard"); ok {
		req.EnableDashboard = dashboard.(bool)
	}

	if ingress, ok := d.GetOk("ingress"); ok {
		req.Ingress = k8s.Ingress(ingress.(string))
	}

	autoscalerReq := &k8s.CreateClusterRequestAutoscalerConfig{}

	if scaleDownDisabled, ok := d.GetOk("autoscaler_config.0.disable_scale_down"); ok {
		autoscalerReq.ScaleDownDisabled = scw.BoolPtr(scaleDownDisabled.(bool))
	}

	if scaleDownDelayAfterAdd, ok := d.GetOk("autoscaler_config.0.scale_down_delay_after_add"); ok {
		autoscalerReq.ScaleDownDelayAfterAdd = expandStringPtr(scaleDownDelayAfterAdd)
	}

	if scaleDownUneededTime, ok := d.GetOk("autoscaler_config.0.scale_down_unneeded_time"); ok {
		autoscalerReq.ScaleDownUnneededTime = expandStringPtr(scaleDownUneededTime)
	}

	if estimator, ok := d.GetOk("autoscaler_config.0.estimator"); ok {
		autoscalerReq.Estimator = k8s.AutoscalerEstimator(estimator.(string))
	}

	if expander, ok := d.GetOk("autoscaler_config.0.expander"); ok {
		autoscalerReq.Expander = k8s.AutoscalerExpander(expander.(string))
	}

	if ignoreDaemonsetsUtilization, ok := d.GetOk("autoscaler_config.0.ignore_daemonsets_utilization"); ok {
		autoscalerReq.IgnoreDaemonsetsUtilization = scw.BoolPtr(ignoreDaemonsetsUtilization.(bool))
	}

	if balanceSimilarNodeGroups, ok := d.GetOk("autoscaler_config.0.balance_similar_node_groups"); ok {
		autoscalerReq.BalanceSimilarNodeGroups = scw.BoolPtr(balanceSimilarNodeGroups.(bool))
	}

	if balanceSimilarNodeGroups, ok := d.GetOk("autoscaler_config.0.balance_similar_node_groups"); ok {
		autoscalerReq.BalanceSimilarNodeGroups = scw.BoolPtr(balanceSimilarNodeGroups.(bool))
	}

	autoscalerReq.ExpendablePodsPriorityCutoff = scw.Int32Ptr(int32(d.Get("autoscaler_config.0.expendable_pods_priority_cutoff").(int)))

	if utilizationThreshold, ok := d.GetOk("autoscaler_config.0.scale_down_utilization_threshold"); ok {
		autoscalerReq.ScaleDownUtilizationThreshold = scw.Float32Ptr(float32(utilizationThreshold.(float64)))
	}

	autoscalerReq.MaxGracefulTerminationSec = scw.Uint32Ptr(uint32(d.Get("autoscaler_config.0.max_graceful_termination_sec").(int)))

	req.AutoscalerConfig = autoscalerReq

	createClusterRequestOpenIDConnectConfig := &k8s.CreateClusterRequestOpenIDConnectConfig{}

	if issuerURL, ok := d.GetOk("open_id_connect_config.0.issuer_url"); ok {
		req.OpenIDConnectConfig = createClusterRequestOpenIDConnectConfig
		createClusterRequestOpenIDConnectConfig.IssuerURL = issuerURL.(string)
	}

	if clientID, ok := d.GetOk("open_id_connect_config.0.client_id"); ok {
		req.OpenIDConnectConfig = createClusterRequestOpenIDConnectConfig
		createClusterRequestOpenIDConnectConfig.ClientID = clientID.(string)
	}

	// createClusterRequestOpenIDConnectConfig is always defined here

	if usernameClaim, ok := d.GetOk("open_id_connect_config.0.username_claim"); ok {
		createClusterRequestOpenIDConnectConfig.UsernameClaim = scw.StringPtr(usernameClaim.(string))
	}

	if usernamePrefix, ok := d.GetOk("open_id_connect_config.0.username_prefix"); ok {
		createClusterRequestOpenIDConnectConfig.UsernamePrefix = scw.StringPtr(usernamePrefix.(string))
	}

	if groupsClaim, ok := d.GetOk("open_id_connect_config.0.groups_claim"); ok {
		createClusterRequestOpenIDConnectConfig.GroupsClaim = scw.StringsPtr(expandStrings(groupsClaim))
	}

	if groupsPrefix, ok := d.GetOk("open_id_connect_config.0.groups_prefix"); ok {
		createClusterRequestOpenIDConnectConfig.GroupsPrefix = scw.StringPtr(groupsPrefix.(string))
	}

	if requiredClaim, ok := d.GetOk("open_id_connect_config.0.required_claim"); ok {
		createClusterRequestOpenIDConnectConfig.RequiredClaim = scw.StringsPtr(expandStrings(requiredClaim))
	}

	autoUpgradeEnable, okAutoUpgradeEnable := d.GetOkExists("auto_upgrade.0.enable")
	autoUpgradeStartHour, okAutoUpgradeStartHour := d.GetOkExists("auto_upgrade.0.maintenance_window_start_hour")
	autoUpgradeDay, okAutoUpgradeDay := d.GetOk("auto_upgrade.0.maintenance_window_day")

	if okAutoUpgradeEnable {
		// check if either all or none of the auto upgrade attribute are set.
		// if one auto upgrade attribute is set, they all must be set.
		// if none is set, auto upgrade attributes will be computed.
		if !(okAutoUpgradeDay && okAutoUpgradeStartHour) {
			return diag.FromErr(fmt.Errorf("all field or zero field of auto_upgrade must be set"))
		}
	}

	clusterAutoUpgradeEnabled := false

	if okAutoUpgradeDay && okAutoUpgradeEnable && okAutoUpgradeStartHour {
		clusterAutoUpgradeEnabled = autoUpgradeEnable.(bool)
		req.AutoUpgrade = &k8s.CreateClusterRequestAutoUpgrade{
			Enable: clusterAutoUpgradeEnabled,
			MaintenanceWindow: &k8s.MaintenanceWindow{
				StartHour: uint32(autoUpgradeStartHour.(int)),
				Day:       k8s.MaintenanceWindowDayOfTheWeek(autoUpgradeDay.(string)),
			},
		}
	}

	version := d.Get("version").(string)
	versionIsOnlyMinor := len(strings.Split(version, ".")) == 2

	if versionIsOnlyMinor != clusterAutoUpgradeEnabled {
		return diag.FromErr(fmt.Errorf("minor version x.y must be used with auto upgrade enabled"))
	}

	if versionIsOnlyMinor {
		version, err = k8sGetLatestVersionFromMinor(ctx, k8sAPI, region, version)
		if err != nil {
			return diag.FromErr(fmt.Errorf("minor version x.y must be used with auto upgrade enabled"))
		}
	}

	req.Version = version

	res, err := k8sAPI.CreateCluster(req, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	err = waitK8SCluster(ctx, k8sAPI, region, res.ID, k8s.ClusterStatusPoolRequired)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newRegionalIDString(region, res.ID))

	return resourceScalewayK8SClusterRead(ctx, d, meta)
}

func resourceScalewayK8SClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	k8sAPI, region, clusterID, err := k8sAPIWithRegionAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	////
	// Read Cluster
	////
	response, err := k8sAPI.GetCluster(&k8s.GetClusterRequest{
		Region:    region,
		ClusterID: clusterID,
	}, scw.WithContext(ctx))
	if err != nil {
		if is404Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	_ = d.Set("region", string(region))
	_ = d.Set("name", response.Name)
	_ = d.Set("organization_id", response.OrganizationID)
	_ = d.Set("project_id", response.ProjectID)
	_ = d.Set("description", response.Description)
	_ = d.Set("cni", response.Cni)
	_ = d.Set("tags", response.Tags)
	_ = d.Set("apiserver_cert_sans", response.ApiserverCertSans)
	_ = d.Set("created_at", response.CreatedAt.Format(time.RFC3339))
	_ = d.Set("updated_at", response.UpdatedAt.Format(time.RFC3339))
	_ = d.Set("apiserver_url", response.ClusterURL)
	_ = d.Set("wildcard_dns", response.DNSWildcard)
	_ = d.Set("status", response.Status.String())
	_ = d.Set("upgrade_available", response.UpgradeAvailable)
	_ = d.Set("feature_gates", response.FeatureGates)
	_ = d.Set("admission_plugins", response.AdmissionPlugins)

	// if autoupgrade is enabled, we only set the minor k8s version (x.y)
	version := response.Version
	if response.AutoUpgrade != nil && response.AutoUpgrade.Enabled {
		version, err = k8sGetMinorVersionFromFull(version)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	_ = d.Set("version", version)

	// autoscaler_config
	_ = d.Set("autoscaler_config", clusterAutoscalerConfigFlatten(response))
	_ = d.Set("open_id_connect_config", clusterOpenIDConnectConfigFlatten(response))
	_ = d.Set("auto_upgrade", clusterAutoUpgradeFlatten(response))

	////
	// Read kubeconfig
	////
	kubeconfig, err := k8sAPI.GetClusterKubeConfig(&k8s.GetClusterKubeConfigRequest{
		Region:    region,
		ClusterID: clusterID,
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	kubeconfigServer, err := kubeconfig.GetServer()
	if err != nil {
		return diag.FromErr(err)
	}

	kubeconfigCa, err := kubeconfig.GetCertificateAuthorityData()
	if err != nil {
		return diag.FromErr(err)
	}

	kubeconfigToken, err := kubeconfig.GetToken()
	if err != nil {
		return diag.FromErr(err)
	}

	kubeconf := map[string]interface{}{}
	kubeconf["config_file"] = string(kubeconfig.GetRaw())
	kubeconf["host"] = kubeconfigServer
	kubeconf["cluster_ca_certificate"] = kubeconfigCa
	kubeconf["token"] = kubeconfigToken

	_ = d.Set("kubeconfig", []map[string]interface{}{kubeconf})

	return nil
}

func resourceScalewayK8SClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	k8sAPI, region, clusterID, err := k8sAPIWithRegionAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	canUpgrade := false

	////
	// Construct UpdateClusterRequest
	////
	updateRequest := &k8s.UpdateClusterRequest{
		Region:    region,
		ClusterID: clusterID,
	}

	if d.HasChange("name") {
		updateRequest.Name = expandStringPtr(d.Get("name"))
	}

	if d.HasChange("description") {
		updateRequest.Description = expandStringPtr(d.Get("description"))
	}

	if d.HasChange("tags") {
		tags := expandStrings(d.Get("tags"))
		updateRequest.Tags = scw.StringsPtr(tags)
	}

	if d.HasChange("apiserver_cert_sans") {
		updateRequest.ApiserverCertSans = scw.StringsPtr(expandStrings(d.Get("apiserver_cert_sans")))
	}

	if d.HasChange("feature_gates") {
		updateRequest.FeatureGates = scw.StringsPtr(expandStrings(d.Get("feature_gates")))
	}

	if d.HasChange("admission_plugins") {
		updateRequest.AdmissionPlugins = scw.StringsPtr(expandStrings(d.Get("admission_plugins")))
	}

	if d.HasChange("ingress") {
		updateRequest.Ingress = k8s.Ingress(d.Get("ingress").(string))
	}

	if d.HasChange("enable_dashboard") {
		updateRequest.EnableDashboard = scw.BoolPtr(d.Get("enable_dashboard").(bool))
	}

	updateRequest.AutoUpgrade = &k8s.UpdateClusterRequestAutoUpgrade{}
	autoupgradeEnabled := d.Get("auto_upgrade.0.enable").(bool)

	if d.HasChange("auto_upgrade.0.enable") {
		updateRequest.AutoUpgrade.Enable = scw.BoolPtr(d.Get("auto_upgrade.0.enable").(bool))
	}

	if d.HasChanges("auto_upgrade.0.maintenance_window_start_hour", "auto_upgrade.0.maintenance_window_day") {
		updateRequest.AutoUpgrade.MaintenanceWindow = &k8s.MaintenanceWindow{}
		updateRequest.AutoUpgrade.MaintenanceWindow.StartHour = uint32(d.Get("auto_upgrade.0.maintenance_window_start_hour").(int))
		updateRequest.AutoUpgrade.MaintenanceWindow.Day = k8s.MaintenanceWindowDayOfTheWeek(d.Get("auto_upgrade.0.maintenance_window_day").(string))
	}

	version := d.Get("version").(string)
	versionIsOnlyMinor := len(strings.Split(version, ".")) == 2

	if versionIsOnlyMinor != autoupgradeEnabled {
		return diag.FromErr(fmt.Errorf("minor version x.y must be used with auto upgrades enabled"))
	}

	if versionIsOnlyMinor {
		version, err = k8sGetLatestVersionFromMinor(ctx, k8sAPI, region, version)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("version") {
		// maybe it's a change from minor to patch or patch to minor
		// we need to check the current version

		clusterResp, err := k8sAPI.GetCluster(&k8s.GetClusterRequest{
			ClusterID: clusterID,
			Region:    region,
		}, scw.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}

		if clusterResp.Version == version {
			// no upgrades if same version
			canUpgrade = false
		} else {
			// we let the API decide if we can upgrade
			canUpgrade = true
		}
	}

	autoscalerReq := &k8s.UpdateClusterRequestAutoscalerConfig{}

	if d.HasChange("autoscaler_config.0.disable_scale_down") {
		autoscalerReq.ScaleDownDisabled = scw.BoolPtr(d.Get("autoscaler_config.0.disable_scale_down").(bool))
	}

	if d.HasChange("autoscaler_config.0.scale_down_delay_after_add") {
		autoscalerReq.ScaleDownDelayAfterAdd = expandStringPtr(d.Get("autoscaler_config.0.scale_down_delay_after_add"))
	}

	if d.HasChange("autoscaler_config.0.scale_down_unneeded_time") {
		autoscalerReq.ScaleDownUnneededTime = expandStringPtr(d.Get("autoscaler_config.0.scale_down_unneeded_time"))
	}

	if d.HasChange("autoscaler_config.0.estimator") {
		autoscalerReq.Estimator = k8s.AutoscalerEstimator(d.Get("autoscaler_config.0.estimator").(string))
	}

	if d.HasChange("autoscaler_config.0.expander") {
		autoscalerReq.Expander = k8s.AutoscalerExpander(d.Get("autoscaler_config.0.expander").(string))
	}

	if d.HasChange("autoscaler_config.0.ignore_daemonsets_utilization") {
		autoscalerReq.IgnoreDaemonsetsUtilization = scw.BoolPtr(d.Get("autoscaler_config.0.ignore_daemonsets_utilization").(bool))
	}

	if d.HasChange("autoscaler_config.0.balance_similar_node_groups") {
		autoscalerReq.BalanceSimilarNodeGroups = scw.BoolPtr(d.Get("autoscaler_config.0.balance_similar_node_groups").(bool))
	}

	if d.HasChange("autoscaler_config.0.expendable_pods_priority_cutoff") {
		autoscalerReq.ExpendablePodsPriorityCutoff = scw.Int32Ptr(int32(d.Get("autoscaler_config.0.expendable_pods_priority_cutoff").(int)))
	}

	if d.HasChange("autoscaler_config.0.scale_down_utilization_threshold") {
		autoscalerReq.ScaleDownUtilizationThreshold = scw.Float32Ptr(float32(d.Get("autoscaler_config.0.scale_down_utilization_threshold").(float64)))
	}

	if d.HasChange("autoscaler_config.0.max_graceful_termination_sec") {
		autoscalerReq.MaxGracefulTerminationSec = scw.Uint32Ptr(uint32(d.Get("autoscaler_config.0.max_graceful_termination_sec").(int)))
	}

	updateRequest.AutoscalerConfig = autoscalerReq

	updateClusterRequestOpenIDConnectConfig := &k8s.UpdateClusterRequestOpenIDConnectConfig{}

	if d.HasChange("open_id_connect_config.0.issuer_url") {
		updateClusterRequestOpenIDConnectConfig.IssuerURL = scw.StringPtr(d.Get("open_id_connect_config.0.issuer_url").(string))
	}

	if d.HasChange("open_id_connect_config.0.client_id") {
		updateClusterRequestOpenIDConnectConfig.ClientID = scw.StringPtr(d.Get("open_id_connect_config.0.client_id").(string))
	}

	if d.HasChange("open_id_connect_config.0.username_claim") {
		updateClusterRequestOpenIDConnectConfig.UsernameClaim = scw.StringPtr(d.Get("open_id_connect_config.0.username_claim").(string))
	}

	if d.HasChange("open_id_connect_config.0.username_prefix") {
		updateClusterRequestOpenIDConnectConfig.UsernamePrefix = scw.StringPtr(d.Get("open_id_connect_config.0.username_prefix").(string))
	}

	if d.HasChange("open_id_connect_config.0.groups_claim") {
		updateClusterRequestOpenIDConnectConfig.GroupsClaim = scw.StringsPtr(expandStrings(d.Get("open_id_connect_config.0.groups_claim")))
	}

	if d.HasChange("open_id_connect_config.0.groups_prefix") {
		updateClusterRequestOpenIDConnectConfig.GroupsPrefix = scw.StringPtr(d.Get("open_id_connect_config.0.groups_prefix").(string))
	}

	if d.HasChange("open_id_connect_config.0.required_claim") {
		updateClusterRequestOpenIDConnectConfig.RequiredClaim = scw.StringsPtr(expandStrings(d.Get("open_id_connect_config.0.required_claim")))
	}

	updateRequest.OpenIDConnectConfig = updateClusterRequestOpenIDConnectConfig

	////
	// Apply Update
	////
	_, err = k8sAPI.UpdateCluster(updateRequest, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	err = waitK8SCluster(ctx, k8sAPI, region, clusterID, k8s.ClusterStatusReady, k8s.ClusterStatusPoolRequired)
	if err != nil {
		return diag.FromErr(err)
	}

	////
	// Upgrade if needed
	////
	if canUpgrade {
		upgradeRequest := &k8s.UpgradeClusterRequest{
			Region:       region,
			ClusterID:    clusterID,
			Version:      version,
			UpgradePools: true,
		}
		_, err = k8sAPI.UpgradeCluster(upgradeRequest)
		if err != nil {
			return diag.FromErr(err)
		}

		err = waitK8SCluster(ctx, k8sAPI, region, clusterID, k8s.ClusterStatusReady, k8s.ClusterStatusPoolRequired)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceScalewayK8SClusterRead(ctx, d, meta)
}

func resourceScalewayK8SClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	k8sAPI, region, clusterID, err := k8sAPIWithRegionAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	deleteAdditionalResources := d.Get("delete_additional_resources").(bool)

	////
	// Delete Cluster
	////
	_, err = k8sAPI.DeleteCluster(&k8s.DeleteClusterRequest{
		Region:                  region,
		ClusterID:               clusterID,
		WithAdditionalResources: deleteAdditionalResources,
	}, scw.WithContext(ctx))
	if err != nil {
		if is404Error(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	err = waitK8SClusterDeleted(ctx, k8sAPI, region, clusterID)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func autoscalerConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"disable_scale_down": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Disable the scale down feature of the autoscaler",
			},
			"scale_down_delay_after_add": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "10m",
				Description: "How long after scale up that scale down evaluation resumes",
			},
			"scale_down_unneeded_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "10m",
				Description: "How long a node should be unneeded before it is eligible for scale down",
			},
			"estimator": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     k8s.AutoscalerEstimatorBinpacking.String(),
				Description: "Type of resource estimator to be used in scale up",
				ValidateFunc: validation.StringInSlice([]string{
					k8s.AutoscalerEstimatorBinpacking.String(),
				}, false),
			},
			"expander": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     k8s.AutoscalerExpanderRandom.String(),
				Description: "Type of node group expander to be used in scale up",
				ValidateFunc: validation.StringInSlice([]string{
					k8s.AutoscalerExpanderRandom.String(),
					k8s.AutoscalerExpanderLeastWaste.String(),
					k8s.AutoscalerExpanderMostPods.String(),
					k8s.AutoscalerExpanderPriority.String(),
					k8s.AutoscalerExpanderPrice.String(),
				}, false),
			},
			"ignore_daemonsets_utilization": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Ignore DaemonSet pods when calculating resource utilization for scaling down",
			},
			"balance_similar_node_groups": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Detect similar node groups and balance the number of nodes between them",
			},
			"expendable_pods_priority_cutoff": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     -10,
				Description: "Pods with priority below cutoff will be expendable. They can be killed without any consideration during scale down and they don't cause scale up. Pods with null priority (PodPriority disabled) are non expendable",
			},
			"scale_down_utilization_threshold": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Default:     0.5,
				Description: "Node utilization level, defined as sum of requested resources divided by capacity, below which a node can be considered for scale down",
			},
			"max_graceful_termination_sec": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     600,
				Description: "Maximum number of seconds the cluster autoscaler waits for pod termination when trying to scale down a node",
			},
		},
	}
}

func openIDConnectConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"issuer_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "URL of the provider which allows the API server to discover public signing keys",
			},
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A client id that all tokens must be issued for",
			},
			"username_claim": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "JWT claim to use as the user name",
			},
			"username_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Prefix prepended to username",
			},
			"groups_claim": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "JWT claim to use as the user's group",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"groups_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Prefix prepended to group claims",
			},
			"required_claim": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Multiple key=value pairs that describes a required claim in the ID Token",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}
