package scaleway

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// Provider config can be used to provide additional config when creating provider.
type ProviderConfig struct {
	// Meta can be used to override Meta that will be used by the provider.
	// This is useful for tests.
	Meta *Meta
}

// DefaultProviderConfig return default ProviderConfig struct
func DefaultProviderConfig() *ProviderConfig {
	return &ProviderConfig{}
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
				"project_id": {
					Type:         schema.TypeString,
					Optional:     true, // To allow user to use organization instead of project
					Description:  "The Scaleway project ID.",
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
				"scaleway_account_ssh_key":               resourceScalewayAccountSSKKey(),
				"scaleway_apple_silicon_server":          resourceScalewayAppleSiliconServer(),
				"scaleway_baremetal_server":              resourceScalewayBaremetalServer(),
				"scaleway_domain_record":                 resourceScalewayDomainRecord(),
				"scaleway_instance_ip":                   resourceScalewayInstanceIP(),
				"scaleway_instance_ip_reverse_dns":       resourceScalewayInstanceIPReverseDNS(),
				"scaleway_instance_volume":               resourceScalewayInstanceVolume(),
				"scaleway_instance_security_group":       resourceScalewayInstanceSecurityGroup(),
				"scaleway_instance_security_group_rules": resourceScalewayInstanceSecurityGroupRules(),
				"scaleway_instance_server":               resourceScalewayInstanceServer(),
				"scaleway_instance_snapshot":             resourceScalewayInstanceSnapshot(),
				"scaleway_instance_placement_group":      resourceScalewayInstancePlacementGroup(),
				"scaleway_instance_private_nic":          resourceScalewayInstancePrivateNIC(),
				"scaleway_iot_hub":                       resourceScalewayIotHub(),
				"scaleway_iot_device":                    resourceScalewayIotDevice(),
				"scaleway_iot_route":                     resourceScalewayIotRoute(),
				"scaleway_iot_network":                   resourceScalewayIotNetwork(),
				"scaleway_k8s_cluster":                   resourceScalewayK8SCluster(),
				"scaleway_k8s_pool":                      resourceScalewayK8SPool(),
				"scaleway_lb":                            resourceScalewayLb(),
				"scaleway_lb_ip":                         resourceScalewayLbIP(),
				"scaleway_lb_backend":                    resourceScalewayLbBackend(),
				"scaleway_lb_certificate":                resourceScalewayLbCertificate(),
				"scaleway_lb_frontend":                   resourceScalewayLbFrontend(),
				"scaleway_registry_namespace":            resourceScalewayRegistryNamespace(),
				"scaleway_rdb_acl":                       resourceScalewayRdbACL(),
				"scaleway_rdb_database":                  resourceScalewayRdbDatabase(),
				"scaleway_rdb_instance":                  resourceScalewayRdbInstance(),
				"scaleway_rdb_privilege":                 resourceScalewayRdbPrivilege(),
				"scaleway_rdb_user":                      resourceScalewayRdbUser(),
				"scaleway_object_bucket":                 resourceScalewayObjectBucket(),
				"scaleway_vpc_public_gateway":            resourceScalewayVPCPublicGateway(),
				"scaleway_vpc_public_gateway_ip":         resourceScalewayVPCPublicGatewayIP(),
				"scaleway_vpc_private_network":           resourceScalewayVPCPrivateNetwork(),
			},

			DataSourcesMap: map[string]*schema.Resource{
				"scaleway_account_ssh_key":         dataSourceScalewayAccountSSHKey(),
				"scaleway_instance_security_group": dataSourceScalewayInstanceSecurityGroup(),
				"scaleway_instance_server":         dataSourceScalewayInstanceServer(),
				"scaleway_instance_image":          dataSourceScalewayInstanceImage(),
				"scaleway_instance_volume":         dataSourceScalewayInstanceVolume(),
				"scaleway_baremetal_offer":         dataSourceScalewayBaremetalOffer(),
				"scaleway_rdb_acl":                 dataSourceScalewayRDBACL(),
				"scaleway_rdb_instance":            dataSourceScalewayRDBInstance(),
				"scaleway_rdb_database":            dataSourceScalewayRDBDatabase(),
				"scaleway_rdb_privilege":           dataSourceScalewayRDBPrivilege(),
				"scaleway_k8s_cluster":             dataSourceScalewayK8SCluster(),
				"scaleway_k8s_pool":                dataSourceScalewayK8SPool(),
				"scaleway_lb_ip":                   dataSourceScalewayLbIP(),
				"scaleway_marketplace_image":       dataSourceScalewayMarketplaceImage(),
				"scaleway_registry_namespace":      dataSourceScalewayRegistryNamespace(),
				"scaleway_registry_image":          dataSourceScalewayRegistryImage(),
				"scaleway_vpc_public_gateway":      dataSourceScalewayVPCPublicGateway(),
				"scaleway_vpc_public_gateway_ip":   dataSourceScalewayVPCPublicGatewayIP(),
				"scaleway_vpc_private_network":     dataSourceScalewayVPCPrivateNetwork(),
			},
		}

		p.ConfigureContextFunc = func(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
			terraformVersion := p.TerraformVersion

			// If we provide meta in config use it. This is useful for tests
			if config.Meta != nil {
				return config.Meta, nil
			}

			meta, err := buildMeta(&MetaConfig{
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

type MetaConfig struct {
	providerSchema   *schema.ResourceData
	terraformVersion string
	forceZone        scw.Zone
	httpClient       *http.Client
}

// providerConfigure creates the Meta object containing the SDK client.
func buildMeta(config *MetaConfig) (*Meta, error) {
	////
	// Load Profile
	////
	profile, err := loadProfile(config.providerSchema)
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

func loadProfile(d *schema.ResourceData) (*scw.Profile, error) {
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
		if accessKey, exist := d.GetOk("access_key"); exist {
			providerProfile.AccessKey = scw.StringPtr(accessKey.(string))
		}
		if secretKey, exist := d.GetOk("secret_key"); exist {
			providerProfile.SecretKey = scw.StringPtr(secretKey.(string))
		}
		if projectID, exist := d.GetOk("project_id"); exist {
			providerProfile.DefaultProjectID = scw.StringPtr(projectID.(string))
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
		l.Debugf("guess region from %s zone", zone)
		region, err := zone.Region()
		if err == nil {
			profile.DefaultRegion = scw.StringPtr(region.String())
		} else {
			l.Debugf("cannot guess region: %w", err)
		}
	}
	return profile, nil
}
