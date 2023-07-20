package scaleway

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var terraformBetaEnabled = os.Getenv(scw.ScwEnableBeta) != ""

// ProviderConfig config can be used to provide additional config when creating provider.
type ProviderConfig struct {
	// Meta can be used to override Meta that will be used by the provider.
	// This is useful for tests.
	Meta *Meta
}

// DefaultProviderConfig return default ProviderConfig struct
func DefaultProviderConfig() *ProviderConfig {
	return &ProviderConfig{}
}

func addBetaResources(provider *schema.Provider) {
	if !terraformBetaEnabled {
		return
	}
	betaResources := map[string]*schema.Resource{}
	betaDataSources := map[string]*schema.Resource{}
	for resourceName, resource := range betaResources {
		provider.ResourcesMap[resourceName] = resource
	}
	for resourceName, resource := range betaDataSources {
		provider.DataSourcesMap[resourceName] = resource
	}
}

// Provider returns a terraform.ResourceProvider.
func Provider(config *ProviderConfig) plugin.ProviderFunc {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"access_key": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The Scaleway access key.",
				},
				"secret_key": {
					Type:         schema.TypeString,
					Optional:     true, // To allow user to use deprecated `token`.
					Description:  "The Scaleway secret Key.",
					ValidateFunc: validationUUID(),
				},
				"profile": {
					Type:        schema.TypeString,
					Optional:    true, // To allow user to use `access_key`, `secret_key`, `project_id`...
					Description: "The Scaleway profile to use.",
				},
				"project_id": {
					Type:         schema.TypeString,
					Optional:     true, // To allow user to use organization instead of project
					Description:  "The Scaleway project ID.",
					ValidateFunc: validationUUID(),
				},
				"organization_id": {
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "The Scaleway organization ID.",
					ValidateFunc: validationUUID(),
				},
				"region": regionSchema(),
				"zone":   zoneSchema(),
				"api_url": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The Scaleway API URL to use.",
				},
			},

			ResourcesMap: map[string]*schema.Resource{
				"scaleway_account_project":                     resourceScalewayAccountProject(),
				"scaleway_account_ssh_key":                     resourceScalewayAccountSSKKey(),
				"scaleway_apple_silicon_server":                resourceScalewayAppleSiliconServer(),
				"scaleway_baremetal_server":                    resourceScalewayBaremetalServer(),
				"scaleway_cockpit":                             resourceScalewayCockpit(),
				"scaleway_cockpit_token":                       resourceScalewayCockpitToken(),
				"scaleway_cockpit_grafana_user":                resourceScalewayCockpitGrafanaUser(),
				"scaleway_container_namespace":                 resourceScalewayContainerNamespace(),
				"scaleway_container_cron":                      resourceScalewayContainerCron(),
				"scaleway_container_domain":                    resourceScalewayContainerDomain(),
				"scaleway_container_trigger":                   resourceScalewayContainerTrigger(),
				"scaleway_domain_record":                       resourceScalewayDomainRecord(),
				"scaleway_domain_zone":                         resourceScalewayDomainZone(),
				"scaleway_flexible_ip":                         resourceScalewayFlexibleIP(),
				"scaleway_function":                            resourceScalewayFunction(),
				"scaleway_function_cron":                       resourceScalewayFunctionCron(),
				"scaleway_function_domain":                     resourceScalewayFunctionDomain(),
				"scaleway_function_namespace":                  resourceScalewayFunctionNamespace(),
				"scaleway_function_token":                      resourceScalewayFunctionToken(),
				"scaleway_function_trigger":                    resourceScalewayFunctionTrigger(),
				"scaleway_iam_api_key":                         resourceScalewayIamAPIKey(),
				"scaleway_iam_application":                     resourceScalewayIamApplication(),
				"scaleway_iam_group":                           resourceScalewayIamGroup(),
				"scaleway_iam_policy":                          resourceScalewayIamPolicy(),
				"scaleway_instance_user_data":                  resourceScalewayInstanceUserData(),
				"scaleway_instance_image":                      resourceScalewayInstanceImage(),
				"scaleway_instance_ip":                         resourceScalewayInstanceIP(),
				"scaleway_instance_ip_reverse_dns":             resourceScalewayInstanceIPReverseDNS(),
				"scaleway_instance_volume":                     resourceScalewayInstanceVolume(),
				"scaleway_instance_security_group":             resourceScalewayInstanceSecurityGroup(),
				"scaleway_instance_security_group_rules":       resourceScalewayInstanceSecurityGroupRules(),
				"scaleway_instance_server":                     resourceScalewayInstanceServer(),
				"scaleway_instance_snapshot":                   resourceScalewayInstanceSnapshot(),
				"scaleway_iam_ssh_key":                         resourceScalewayIamSSKKey(),
				"scaleway_instance_placement_group":            resourceScalewayInstancePlacementGroup(),
				"scaleway_instance_private_nic":                resourceScalewayInstancePrivateNIC(),
				"scaleway_iot_hub":                             resourceScalewayIotHub(),
				"scaleway_iot_device":                          resourceScalewayIotDevice(),
				"scaleway_iot_route":                           resourceScalewayIotRoute(),
				"scaleway_iot_network":                         resourceScalewayIotNetwork(),
				"scaleway_k8s_cluster":                         resourceScalewayK8SCluster(),
				"scaleway_k8s_pool":                            resourceScalewayK8SPool(),
				"scaleway_lb":                                  resourceScalewayLb(),
				"scaleway_lb_acl":                              resourceScalewayLbACL(),
				"scaleway_lb_ip":                               resourceScalewayLbIP(),
				"scaleway_lb_backend":                          resourceScalewayLbBackend(),
				"scaleway_lb_certificate":                      resourceScalewayLbCertificate(),
				"scaleway_lb_frontend":                         resourceScalewayLbFrontend(),
				"scaleway_lb_route":                            resourceScalewayLbRoute(),
				"scaleway_registry_namespace":                  resourceScalewayRegistryNamespace(),
				"scaleway_tem_domain":                          resourceScalewayTemDomain(),
				"scaleway_container":                           resourceScalewayContainer(),
				"scaleway_container_token":                     resourceScalewayContainerToken(),
				"scaleway_rdb_acl":                             resourceScalewayRdbACL(),
				"scaleway_rdb_database":                        resourceScalewayRdbDatabase(),
				"scaleway_rdb_database_backup":                 resourceScalewayRdbDatabaseBackup(),
				"scaleway_rdb_instance":                        resourceScalewayRdbInstance(),
				"scaleway_rdb_privilege":                       resourceScalewayRdbPrivilege(),
				"scaleway_rdb_user":                            resourceScalewayRdbUser(),
				"scaleway_rdb_read_replica":                    resourceScalewayRdbReadReplica(),
				"scaleway_redis_cluster":                       resourceScalewayRedisCluster(),
				"scaleway_object":                              resourceScalewayObject(),
				"scaleway_object_bucket":                       resourceScalewayObjectBucket(),
				"scaleway_object_bucket_acl":                   resourceScalewayObjectBucketACL(),
				"scaleway_object_bucket_lock_configuration":    resourceObjectLockConfiguration(),
				"scaleway_object_bucket_policy":                resourceScalewayObjectBucketPolicy(),
				"scaleway_object_bucket_website_configuration": ResourceBucketWebsiteConfiguration(),
				"scaleway_mnq_namespace":                       resourceScalewayMNQNamespace(),
				"scaleway_mnq_credential":                      resourceScalewayMNQCredential(),
				"scaleway_mnq_queue":                           resourceScalewayMNQQueue(),
				"scaleway_secret":                              resourceScalewaySecret(),
				"scaleway_secret_version":                      resourceScalewaySecretVersion(),
				"scaleway_vpc":                                 resourceScalewayVPC(),
				"scaleway_vpc_public_gateway":                  resourceScalewayVPCPublicGateway(),
				"scaleway_vpc_gateway_network":                 resourceScalewayVPCGatewayNetwork(),
				"scaleway_vpc_public_gateway_dhcp":             resourceScalewayVPCPublicGatewayDHCP(),
				"scaleway_vpc_public_gateway_dhcp_reservation": resourceScalewayVPCPublicGatewayDHCPReservation(),
				"scaleway_vpc_public_gateway_ip":               resourceScalewayVPCPublicGatewayIP(),
				"scaleway_vpc_public_gateway_ip_reverse_dns":   resourceScalewayVPCPublicGatewayIPReverseDNS(),
				"scaleway_vpc_public_gateway_pat_rule":         resourceScalewayVPCPublicGatewayPATRule(),
				"scaleway_vpc_private_network":                 resourceScalewayVPCPrivateNetwork(),
				"scaleway_webhosting":                          resourceScalewayWebhosting(),
			},

			DataSourcesMap: map[string]*schema.Resource{
				"scaleway_account_project":                     dataSourceScalewayAccountProject(),
				"scaleway_account_ssh_key":                     dataSourceScalewayAccountSSHKey(),
				"scaleway_availability_zones":                  DataSourceAvailabilityZones(),
				"scaleway_baremetal_offer":                     dataSourceScalewayBaremetalOffer(),
				"scaleway_baremetal_option":                    dataSourceScalewayBaremetalOption(),
				"scaleway_baremetal_os":                        dataSourceScalewayBaremetalOs(),
				"scaleway_baremetal_server":                    dataSourceScalewayBaremetalServer(),
				"scaleway_cockpit":                             dataSourceScalewayCockpit(),
				"scaleway_cockpit_plan":                        dataSourceScalewayCockpitPlan(),
				"scaleway_domain_record":                       dataSourceScalewayDomainRecord(),
				"scaleway_domain_zone":                         dataSourceScalewayDomainZone(),
				"scaleway_container_namespace":                 dataSourceScalewayContainerNamespace(),
				"scaleway_container":                           dataSourceScalewayContainer(),
				"scaleway_function":                            dataSourceScalewayFunction(),
				"scaleway_function_namespace":                  dataSourceScalewayFunctionNamespace(),
				"scaleway_iam_application":                     dataSourceScalewayIamApplication(),
				"scaleway_flexible_ip":                         dataSourceScalewayFlexibleIP(),
				"scaleway_iam_group":                           dataSourceScalewayIamGroup(),
				"scaleway_iam_ssh_key":                         dataSourceScalewayIamSSHKey(),
				"scaleway_iam_user":                            dataSourceScalewayIamUser(),
				"scaleway_instance_ip":                         dataSourceScalewayInstanceIP(),
				"scaleway_instance_private_nic":                dataSourceScalewayInstancePrivateNIC(),
				"scaleway_instance_security_group":             dataSourceScalewayInstanceSecurityGroup(),
				"scaleway_instance_server":                     dataSourceScalewayInstanceServer(),
				"scaleway_instance_servers":                    dataSourceScalewayInstanceServers(),
				"scaleway_instance_image":                      dataSourceScalewayInstanceImage(),
				"scaleway_instance_volume":                     dataSourceScalewayInstanceVolume(),
				"scaleway_instance_snapshot":                   dataSourceScalewayInstanceSnapshot(),
				"scaleway_iot_hub":                             dataSourceScalewayIotHub(),
				"scaleway_iot_device":                          dataSourceScalewayIotDevice(),
				"scaleway_k8s_cluster":                         dataSourceScalewayK8SCluster(),
				"scaleway_k8s_pool":                            dataSourceScalewayK8SPool(),
				"scaleway_k8s_version":                         dataSourceScalewayK8SVersion(),
				"scaleway_lb":                                  dataSourceScalewayLb(),
				"scaleway_lbs":                                 dataSourceScalewayLbs(),
				"scaleway_lb_acls":                             dataSourceScalewayLbACLs(),
				"scaleway_lb_backend":                          dataSourceScalewayLbBackend(),
				"scaleway_lb_backends":                         dataSourceScalewayLbBackends(),
				"scaleway_lb_certificate":                      dataSourceScalewayLbCertificate(),
				"scaleway_lb_frontend":                         dataSourceScalewayLbFrontend(),
				"scaleway_lb_frontends":                        dataSourceScalewayLbFrontends(),
				"scaleway_lb_ip":                               dataSourceScalewayLbIP(),
				"scaleway_lb_ips":                              dataSourceScalewayLbIPs(),
				"scaleway_lb_route":                            dataSourceScalewayLbRoute(),
				"scaleway_lb_routes":                           dataSourceScalewayLbRoutes(),
				"scaleway_marketplace_image":                   dataSourceScalewayMarketplaceImage(),
				"scaleway_object_bucket":                       dataSourceScalewayObjectBucket(),
				"scaleway_object_bucket_policy":                dataSourceScalewayObjectBucketPolicy(),
				"scaleway_rdb_acl":                             dataSourceScalewayRDBACL(),
				"scaleway_rdb_instance":                        dataSourceScalewayRDBInstance(),
				"scaleway_rdb_database":                        dataSourceScalewayRDBDatabase(),
				"scaleway_rdb_database_backup":                 dataSourceScalewayRDBDatabaseBackup(),
				"scaleway_rdb_privilege":                       dataSourceScalewayRDBPrivilege(),
				"scaleway_redis_cluster":                       dataSourceScalewayRedisCluster(),
				"scaleway_registry_namespace":                  dataSourceScalewayRegistryNamespace(),
				"scaleway_tem_domain":                          dataSourceScalewayTemDomain(),
				"scaleway_secret":                              dataSourceScalewaySecret(),
				"scaleway_secret_version":                      dataSourceScalewaySecretVersion(),
				"scaleway_registry_image":                      dataSourceScalewayRegistryImage(),
				"scaleway_vpc":                                 dataSourceScalewayVPC(),
				"scaleway_vpcs":                                dataSourceScalewayVPCs(),
				"scaleway_vpc_public_gateway":                  dataSourceScalewayVPCPublicGateway(),
				"scaleway_vpc_gateway_network":                 dataSourceScalewayVPCGatewayNetwork(),
				"scaleway_vpc_public_gateway_dhcp":             dataSourceScalewayVPCPublicGatewayDHCP(),
				"scaleway_vpc_public_gateway_dhcp_reservation": dataSourceScalewayVPCPublicGatewayDHCPReservation(),
				"scaleway_vpc_public_gateway_ip":               dataSourceScalewayVPCPublicGatewayIP(),
				"scaleway_vpc_private_network":                 dataSourceScalewayVPCPrivateNetwork(),
				"scaleway_vpc_public_gateway_pat_rule":         dataSourceScalewayVPCPublicGatewayPATRule(),
				"scaleway_webhosting":                          dataSourceScalewayWebhosting(),
				"scaleway_webhosting_offer":                    dataSourceScalewayWebhostingOffer(),
			},
		}

		addBetaResources(p)

		p.ConfigureContextFunc = func(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
			terraformVersion := p.TerraformVersion

			// If we provide meta in config use it. This is useful for tests
			if config.Meta != nil {
				return config.Meta, nil
			}

			meta, err := buildMeta(ctx, &metaConfig{
				providerSchema:   data,
				terraformVersion: terraformVersion,
			})
			if err != nil {
				return nil, diag.FromErr(err)
			}
			return meta, nil
		}

		return p
	}
}

// Meta contains config and SDK clients used by resources.
//
// This meta value is passed into all resources.
type Meta struct {
	// scwClient is the Scaleway SDK client.
	scwClient *scw.Client
	// httpClient can be either a regular http.Client used to make real HTTP requests
	// or it can be a http.Client used to record and replay cassettes which is useful
	// to replay recorded interactions with APIs locally
	httpClient *http.Client
}

type metaConfig struct {
	providerSchema      *schema.ResourceData
	terraformVersion    string
	forceZone           scw.Zone
	forceProjectID      string
	forceOrganizationID string
	forceAccessKey      string
	forceSecretKey      string
	httpClient          *http.Client
}

// providerConfigure creates the Meta object containing the SDK client.
func buildMeta(ctx context.Context, config *metaConfig) (*Meta, error) {
	////
	// Load Profile
	////
	profile, err := loadProfile(ctx, config.providerSchema)
	if err != nil {
		return nil, err
	}
	if config.forceZone != "" {
		region, err := config.forceZone.Region()
		if err != nil {
			return nil, err
		}
		profile.DefaultRegion = scw.StringPtr(region.String())
		profile.DefaultZone = scw.StringPtr(config.forceZone.String())
	}
	if config.forceProjectID != "" {
		profile.DefaultProjectID = scw.StringPtr(config.forceProjectID)
	}
	if config.forceOrganizationID != "" {
		profile.DefaultOrganizationID = scw.StringPtr(config.forceOrganizationID)
	}
	if config.forceAccessKey != "" {
		profile.AccessKey = scw.StringPtr(config.forceAccessKey)
	}
	if config.forceSecretKey != "" {
		profile.SecretKey = scw.StringPtr(config.forceSecretKey)
	}

	// TODO validated profile

	////
	// Create scaleway SDK client
	////
	opts := []scw.ClientOption{
		scw.WithUserAgent(fmt.Sprintf("terraform-provider/%s terraform/%s", version, config.terraformVersion)),
		scw.WithProfile(profile),
	}

	httpClient := &http.Client{Transport: newRetryableTransport(http.DefaultTransport)}
	if config.httpClient != nil {
		httpClient = config.httpClient
	}
	opts = append(opts, scw.WithHTTPClient(httpClient))

	scwClient, err := scw.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	return &Meta{
		scwClient:  scwClient,
		httpClient: httpClient,
	}, nil
}

//gocyclo:ignore
func loadProfile(ctx context.Context, d *schema.ResourceData) (*scw.Profile, error) {
	config, err := scw.LoadConfig()
	// If the config file do not exist, don't return an error as we may find config in ENV or flags.
	if _, isNotFoundError := err.(*scw.ConfigFileNotFoundError); isNotFoundError {
		config = &scw.Config{}
	} else if err != nil {
		return nil, err
	}

	// By default we set default zone and region to fr-par
	defaultZoneProfile := &scw.Profile{
		DefaultRegion: scw.StringPtr(scw.RegionFrPar.String()),
		DefaultZone:   scw.StringPtr(scw.ZoneFrPar1.String()),
	}

	activeProfile, err := config.GetActiveProfile()
	if err != nil {
		return nil, err
	}
	envProfile := scw.LoadEnvProfile()

	providerProfile := &scw.Profile{}
	if d != nil {
		if profileName, exist := d.GetOk("profile"); exist {
			profileFromConfig, err := config.GetProfile(profileName.(string))
			if err == nil {
				providerProfile = profileFromConfig
			}
		}
		if accessKey, exist := d.GetOk("access_key"); exist {
			providerProfile.AccessKey = scw.StringPtr(accessKey.(string))
		}
		if secretKey, exist := d.GetOk("secret_key"); exist {
			providerProfile.SecretKey = scw.StringPtr(secretKey.(string))
		}
		if projectID, exist := d.GetOk("project_id"); exist {
			providerProfile.DefaultProjectID = scw.StringPtr(projectID.(string))
		}
		if orgID, exist := d.GetOk("organization_id"); exist {
			providerProfile.DefaultOrganizationID = scw.StringPtr(orgID.(string))
		}
		if region, exist := d.GetOk("region"); exist {
			providerProfile.DefaultRegion = scw.StringPtr(region.(string))
		}
		if zone, exist := d.GetOk("zone"); exist {
			providerProfile.DefaultZone = scw.StringPtr(zone.(string))
		}
		if apiURL, exist := d.GetOk("api_url"); exist {
			providerProfile.APIURL = scw.StringPtr(apiURL.(string))
		}
	}

	profile := scw.MergeProfiles(defaultZoneProfile, activeProfile, providerProfile, envProfile)

	// If profile have a defaultZone but no defaultRegion we set the defaultRegion
	// to the one of the defaultZone
	if profile.DefaultZone != nil && *profile.DefaultZone != "" &&
		(profile.DefaultRegion == nil || *profile.DefaultRegion == "") {
		zone := scw.Zone(*profile.DefaultZone)
		tflog.Debug(ctx, fmt.Sprintf("guess region from %s zone", zone))
		region, err := zone.Region()
		if err == nil {
			profile.DefaultRegion = scw.StringPtr(region.String())
		} else {
			tflog.Debug(ctx, fmt.Sprintf("cannot guess region: %s", err.Error()))
		}
	}
	return profile, nil
}
