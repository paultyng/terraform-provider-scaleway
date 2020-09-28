package scaleway

import (
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mitchellh/go-homedir"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var mu = sync.Mutex{}

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {

	// Try to migrate config
	_, err := scw.MigrateLegacyConfig()
	if err != nil {
		l.Errorf("cannot migrate configuration: %s", err)
		return nil
	}

	// Load active profile
	var activeProfile *scw.Profile
	scwConfig, err := scw.LoadConfig()
	if err != nil {
		l.Warningf("cannot load configuration: %s", err)
	} else {
		activeProfile, err = scwConfig.GetActiveProfile()
		if err != nil {
			l.Errorf("cannot load configuration: %s", err)
		}
	}

	// load env
	envProfile := scw.LoadEnvProfile()

	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Scaleway access key.",
				DefaultFunc: func() (interface{}, error) {
					// Keep the deprecated behavior
					if accessKey := os.Getenv("SCALEWAY_ACCESS_KEY"); accessKey != "" {
						l.Warningf("SCALEWAY_ACCESS_KEY is deprecated, please use SCW_ACCESS_KEY instead")
						return accessKey, nil
					}

					// Add special case temporary until acceptance test are fixed.
					if accessKey := os.Getenv("SCALEWAY_ACCESS_KEY_"); accessKey != "" {
						l.Warningf("SCALEWAY_ACCESS_KEY_ is deprecated, please use SCW_ACCESS_KEY instead")
						return accessKey, nil
					}
					if envProfile.AccessKey != nil {
						return *envProfile.AccessKey, nil
					}
					if activeProfile != nil && activeProfile.AccessKey != nil {
						return *activeProfile.AccessKey, nil
					}
					return nil, nil
				},
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true, // To allow user to use deprecated `token`.
				Description: "The Scaleway secret Key.",
				DefaultFunc: func() (interface{}, error) {
					// No error is returned here to allow user to use deprecated `token`.
					if envProfile.SecretKey != nil {
						return *envProfile.SecretKey, nil
					}
					if activeProfile != nil && activeProfile.SecretKey != nil {
						return *activeProfile.SecretKey, nil
					}

					// Keep the deprecated behavior from 'token'.
					for _, k := range []string{"SCALEWAY_TOKEN", "SCALEWAY_ACCESS_KEY"} {
						if os.Getenv(k) != "" {
							l.Warningf("%s is deprecated, please use SCW_SECRET_KEY instead", k)
							return os.Getenv(k), nil
						}
					}
					if path, err := homedir.Expand("~/.scwrc"); err == nil {
						scwAPIKey, _, err := readDeprecatedScalewayConfig(path)
						if err != nil {
							// No error is returned here to allow user to use `secret_key`.
							l.Errorf("cannot parse deprecated config file: %s", err)
							return nil, nil
						}
						// Depreciation log is already handled by scw config.
						return scwAPIKey, nil
					}
					// No error is returned here to allow user to use `secret_key`.
					return nil, nil
				},
				ValidateFunc: validationUUID(),
			},
			"organization_id": {
				Type:        schema.TypeString,
				Optional:    true, // To allow user to use deprecated `organization`.
				Description: "The Scaleway organization ID.",
				DefaultFunc: schema.SchemaDefaultFunc(func() (interface{}, error) {
					if envProfile.DefaultOrganizationID != nil {
						return *envProfile.DefaultOrganizationID, nil
					}
					if activeProfile != nil && activeProfile.DefaultOrganizationID != nil {
						return *activeProfile.DefaultOrganizationID, nil
					}

					// Keep the deprecated behavior of 'organization'.
					if organization := os.Getenv("SCALEWAY_ORGANIZATION"); organization != "" {
						l.Warningf("SCALEWAY_ORGANIZATION is deprecated, please use SCW_DEFAULT_ORGANIZATION_ID instead")
						return organization, nil
					}
					if path, err := homedir.Expand("~/.scwrc"); err == nil {
						_, scwOrganization, err := readDeprecatedScalewayConfig(path)
						if err != nil {
							// No error is returned here to allow user to use `organization_id`.
							l.Errorf("cannot parse deprecated config file: %s", err)
							return nil, nil
						}
						return scwOrganization, nil
					}
					// No error is returned here to allow user to use `organization_id`.
					return nil, nil
				}),
				ValidateFunc: validationUUID(),
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Scaleway default region to use for your resources.",
				DefaultFunc: func() (interface{}, error) {
					// Keep the deprecated behavior
					// Note: The deprecated region format conversion is handled in `config.GetDeprecatedClient`.
					if region := os.Getenv("SCALEWAY_REGION"); region != "" {
						l.Warningf("SCALEWAY_REGION is deprecated, please use SCW_DEFAULT_REGION instead")
						return region, nil
					}
					if envProfile.DefaultRegion != nil {
						return *envProfile.DefaultRegion, nil
					}
					if activeProfile != nil && activeProfile.DefaultRegion != nil {
						return *activeProfile.DefaultRegion, nil
					}
					return string(scw.RegionFrPar), nil
				},
				ValidateFunc: validationRegion(),
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Scaleway default zone to use for your resources.",
				DefaultFunc: func() (interface{}, error) {
					if envProfile.DefaultZone != nil {
						return *envProfile.DefaultZone, nil
					}
					if activeProfile != nil && activeProfile.DefaultZone != nil {
						return *activeProfile.DefaultZone, nil
					}
					return nil, nil
				},
				ValidateFunc: validationZone(),
			},
			"api_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Scaleway API URL to use.",
				DefaultFunc: func() (interface{}, error) {
					if envProfile.APIURL != nil {
						return *envProfile.APIURL, nil
					}
					if activeProfile != nil && activeProfile.APIURL != nil {
						return *activeProfile.APIURL, nil
					}
					return nil, nil
				},
			},
			// Deprecated values
			"token": {
				Type:       schema.TypeString,
				Optional:   true, // To allow user to use `secret_key`.
				Deprecated: "Use `secret_key` instead.",
			},
			"organization": {
				Type:       schema.TypeString,
				Optional:   true, // To allow user to use `organization_id`.
				Deprecated: "Use `organization_id` instead.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"scaleway_account_ssh_key":               resourceScalewayAccountSSKKey(),
			"scaleway_baremetal_server":              resourceScalewayBaremetalServer(),
			"scaleway_instance_ip":                   resourceScalewayInstanceIP(),
			"scaleway_instance_ip_reverse_dns":       resourceScalewayInstanceIPReverseDns(),
			"scaleway_instance_volume":               resourceScalewayInstanceVolume(),
			"scaleway_instance_security_group":       resourceScalewayInstanceSecurityGroup(),
			"scaleway_instance_security_group_rules": resourceScalewayInstanceSecurityGroupRules(),
			"scaleway_instance_server":               resourceScalewayInstanceServer(),
			"scaleway_instance_placement_group":      resourceScalewayInstancePlacementGroup(),
			"scaleway_k8s_cluster_beta":              resourceScalewayK8SClusterBeta(),
			"scaleway_k8s_pool_beta":                 resourceScalewayK8SPoolBeta(),
			"scaleway_lb_beta":                       resourceScalewayLbBeta(),
			"scaleway_lb_ip_beta":                    resourceScalewayLbIPBeta(),
			"scaleway_lb_backend_beta":               resourceScalewayLbBackendBeta(),
			"scaleway_lb_certificate_beta":           resourceScalewayLbCertificateBeta(),
			"scaleway_lb_frontend_beta":              resourceScalewayLbFrontendBeta(),
			"scaleway_registry_namespace_beta":       resourceScalewayRegistryNamespaceBeta(),
			"scaleway_rdb_instance_beta":             resourceScalewayRdbInstanceBeta(),
			"scaleway_rdb_user_beta":                 resourceScalewayRdbUserBeta(),
			"scaleway_object_bucket":                 resourceScalewayObjectBucket(),
			"scaleway_user_data":                     resourceScalewayUserData(),
			"scaleway_server":                        resourceScalewayServer(),
			"scaleway_token":                         resourceScalewayToken(),
			"scaleway_ssh_key":                       resourceScalewaySSHKey(),
			"scaleway_ip":                            resourceScalewayIP(),
			"scaleway_ip_reverse_dns":                resourceScalewayIPReverseDNS(),
			"scaleway_security_group":                resourceScalewaySecurityGroup(),
			"scaleway_security_group_rule":           resourceScalewaySecurityGroupRule(),
			"scaleway_volume":                        resourceScalewayVolume(),
			"scaleway_volume_attachment":             resourceScalewayVolumeAttachment(),
			"scaleway_domain_dns_records":            resourceScalewayDomainDNSRecords(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"scaleway_bootscript":              dataSourceScalewayBootscript(),
			"scaleway_image":                   dataSourceScalewayImage(),
			"scaleway_security_group":          dataSourceScalewaySecurityGroup(),
			"scaleway_volume":                  dataSourceScalewayVolume(),
			"scaleway_account_ssh_key":         dataSourceScalewayAccountSSHKey(),
			"scaleway_instance_security_group": dataSourceScalewayInstanceSecurityGroup(),
			"scaleway_instance_server":         dataSourceScalewayInstanceServer(),
			"scaleway_instance_image":          dataSourceScalewayInstanceImage(),
			"scaleway_instance_volume":         dataSourceScalewayInstanceVolume(),
			"scaleway_baremetal_offer":         dataSourceScalewayBaremetalOffer(),
			"scaleway_lb_ip_beta":              dataSourceScalewayLbIPBeta(),
			"scaleway_marketplace_image_beta":  dataSourceScalewayMarketplaceImageBeta(),
			"scaleway_registry_namespace_beta": dataSourceScalewayRegistryNamespaceBeta(),
			"scaleway_registry_image_beta":     dataSourceScalewayRegistryImageBeta(),
		},
	}

	p.ConfigureFunc = func(data *schema.ResourceData) (i interface{}, e error) {
		terraformVersion := p.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(data, terraformVersion)
	}

	return p
}

// providerConfigure creates the Meta object containing the SDK client.
func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	apiKey := ""
	if v, ok := d.Get("secret_key").(string); ok {
		apiKey = v
	} else if v, ok := d.Get("token").(string); ok {
		apiKey = v
	} else {
		if v, ok := d.Get("access_key").(string); ok {
			l.Warningf("you seem to use the access_key instead of secret_key in the provider. This bogus behavior is deprecated, please use the `secret_key` field instead.")
			apiKey = v
		}
	}

	organizationID := d.Get("organization_id").(string)
	if organizationID == "" {
		organizationID = d.Get("organization").(string)
	}

	if apiKey == "" {
		if path, err := homedir.Expand("~/.scwrc"); err == nil {
			scwAPIKey, scwOrganization, err := readDeprecatedScalewayConfig(path)
			if err != nil {
				return nil, fmt.Errorf("error loading credentials from SCW: %s", err)
			}
			apiKey = scwAPIKey
			organizationID = scwOrganization
		}
	}

	rawRegion := d.Get("region").(string)
	region, err := scw.ParseRegion(rawRegion)
	if err != nil {
		return nil, err
	}

	rawZone := d.Get("zone").(string)
	zone, err := scw.ParseZone(rawZone)
	if err != nil {
		return nil, err
	}

	meta := &Meta{
		AccessKey:             d.Get("access_key").(string),
		SecretKey:             apiKey,
		DefaultOrganizationID: organizationID,
		DefaultRegion:         region,
		DefaultZone:           zone,
		TerraformVersion:      terraformVersion,
		APIURL:                d.Get("api_url").(string),
	}

	err = meta.bootstrap()
	if err != nil {
		return nil, err
	}

	// fetch known scaleway server types to support validation in r/server
	if len(commercialServerTypes) == 0 {
		instanceAPI := instance.NewAPI(meta.scwClient)
		availabilityResp, err := instanceAPI.GetServerTypesAvailability(&instance.GetServerTypesAvailabilityRequest{}, scw.WithAllPages())
		if err == nil {
			for k := range availabilityResp.Servers {
				commercialServerTypes = append(commercialServerTypes, k)
			}
			sort.StringSlice(commercialServerTypes).Sort()
		}

		if os.Getenv("DISABLE_SCALEWAY_SERVER_TYPE_VALIDATION") != "" {
			commercialServerTypes = commercialServerTypes[:0]
		}
	}

	return meta, nil
}
