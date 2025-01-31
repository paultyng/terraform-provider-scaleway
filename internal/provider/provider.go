package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/regional"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/zonal"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/meta"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/account"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/applesilicon"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/az"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/baremetal"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/billing"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/block"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/cockpit"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/container"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/domain"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/flexibleip"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/function"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/iam"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/inference"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/instance"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/iot"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/ipam"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/jobs"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/k8s"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/lb"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/marketplace"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/mnq"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/mongodb"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/object"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/rdb"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/redis"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/registry"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/scwconfig"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/sdb"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/secret"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/tem"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/vpc"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/vpcgw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/webhosting"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/verify"
)

var terraformBetaEnabled = os.Getenv(scw.ScwEnableBeta) != ""

// Config can be used to provide additional config when creating provider.
type Config struct {
	// Meta can be used to override Meta that will be used by the provider.
	// This is useful for tests.
	Meta *meta.Meta
}

// DefaultConfig return default Config struct
func DefaultConfig() *Config {
	return &Config{}
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
func Provider(config *Config) plugin.ProviderFunc {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"access_key": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The Scaleway access key.",
				},
				"secret_key": {
					Type:             schema.TypeString,
					Optional:         true, // To allow user to use deprecated `token`.
					Description:      "The Scaleway secret Key.",
					ValidateDiagFunc: verify.IsUUID(),
				},
				"profile": {
					Type:        schema.TypeString,
					Optional:    true, // To allow user to use `access_key`, `secret_key`, `project_id`...
					Description: "The Scaleway profile to use.",
				},
				"project_id": {
					Type:             schema.TypeString,
					Optional:         true, // To allow user to use organization instead of project
					Description:      "The Scaleway project ID.",
					ValidateDiagFunc: verify.IsUUID(),
				},
				"organization_id": {
					Type:             schema.TypeString,
					Optional:         true,
					Description:      "The Scaleway organization ID.",
					ValidateDiagFunc: verify.IsUUID(),
				},
				"region": regional.Schema(),
				"zone":   zonal.Schema(),
				"api_url": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The Scaleway API URL to use.",
				},
			},

			ResourcesMap: map[string]*schema.Resource{
				"scaleway_account_project":                     account.ResourceProject(),
				"scaleway_account_ssh_key":                     iam.ResourceSSKKey(),
				"scaleway_apple_silicon_server":                applesilicon.ResourceServer(),
				"scaleway_baremetal_server":                    baremetal.ResourceServer(),
				"scaleway_block_snapshot":                      block.ResourceSnapshot(),
				"scaleway_block_volume":                        block.ResourceVolume(),
				"scaleway_cockpit":                             cockpit.ResourceCockpit(),
				"scaleway_cockpit_source":                      cockpit.ResourceCockpitSource(),
				"scaleway_cockpit_grafana_user":                cockpit.ResourceCockpitGrafanaUser(),
				"scaleway_cockpit_token":                       cockpit.ResourceToken(),
				"scaleway_cockpit_alert_manager":               cockpit.ResourceCockpitAlertManager(),
				"scaleway_container":                           container.ResourceContainer(),
				"scaleway_container_cron":                      container.ResourceCron(),
				"scaleway_container_domain":                    container.ResourceDomain(),
				"scaleway_container_namespace":                 container.ResourceNamespace(),
				"scaleway_container_token":                     container.ResourceToken(),
				"scaleway_container_trigger":                   container.ResourceTrigger(),
				"scaleway_domain_record":                       domain.ResourceRecord(),
				"scaleway_domain_domains_registration":         domain.ResourceDomainsRegistration(),
				"scaleway_domain_zone":                         domain.ResourceZone(),
				"scaleway_flexible_ip":                         flexibleip.ResourceIP(),
				"scaleway_flexible_ip_mac_address":             flexibleip.ResourceMACAddress(),
				"scaleway_function":                            function.ResourceFunction(),
				"scaleway_function_cron":                       function.ResourceCron(),
				"scaleway_function_domain":                     function.ResourceDomain(),
				"scaleway_function_namespace":                  function.ResourceNamespace(),
				"scaleway_function_token":                      function.ResourceToken(),
				"scaleway_function_trigger":                    function.ResourceTrigger(),
				"scaleway_iam_api_key":                         iam.ResourceAPIKey(),
				"scaleway_iam_application":                     iam.ResourceApplication(),
				"scaleway_iam_group":                           iam.ResourceGroup(),
				"scaleway_iam_group_membership":                iam.ResourceGroupMembership(),
				"scaleway_iam_policy":                          iam.ResourcePolicy(),
				"scaleway_iam_ssh_key":                         iam.ResourceSSKKey(),
				"scaleway_iam_user":                            iam.ResourceUser(),
				"scaleway_inference_deployment":                inference.ResourceDeployment(),
				"scaleway_instance_image":                      instance.ResourceImage(),
				"scaleway_instance_ip":                         instance.ResourceIP(),
				"scaleway_instance_ip_reverse_dns":             instance.ResourceIPReverseDNS(),
				"scaleway_instance_placement_group":            instance.ResourcePlacementGroup(),
				"scaleway_instance_private_nic":                instance.ResourcePrivateNIC(),
				"scaleway_instance_security_group":             instance.ResourceSecurityGroup(),
				"scaleway_instance_security_group_rules":       instance.ResourceSecurityGroupRules(),
				"scaleway_instance_server":                     instance.ResourceServer(),
				"scaleway_instance_snapshot":                   instance.ResourceSnapshot(),
				"scaleway_instance_user_data":                  instance.ResourceUserData(),
				"scaleway_instance_volume":                     instance.ResourceVolume(),
				"scaleway_iot_device":                          iot.ResourceDevice(),
				"scaleway_iot_hub":                             iot.ResourceHub(),
				"scaleway_iot_network":                         iot.ResourceNetwork(),
				"scaleway_iot_route":                           iot.ResourceRoute(),
				"scaleway_ipam_ip":                             ipam.ResourceIP(),
				"scaleway_ipam_ip_reverse_dns":                 ipam.ResourceIPReverseDNS(),
				"scaleway_job_definition":                      jobs.ResourceDefinition(),
				"scaleway_k8s_cluster":                         k8s.ResourceCluster(),
				"scaleway_k8s_pool":                            k8s.ResourcePool(),
				"scaleway_lb":                                  lb.ResourceLb(),
				"scaleway_lb_acl":                              lb.ResourceACL(),
				"scaleway_lb_backend":                          lb.ResourceBackend(),
				"scaleway_lb_certificate":                      lb.ResourceCertificate(),
				"scaleway_lb_frontend":                         lb.ResourceFrontend(),
				"scaleway_lb_ip":                               lb.ResourceIP(),
				"scaleway_lb_route":                            lb.ResourceRoute(),
				"scaleway_mnq_nats_account":                    mnq.ResourceNatsAccount(),
				"scaleway_mnq_nats_credentials":                mnq.ResourceNatsCredentials(),
				"scaleway_mnq_sns":                             mnq.ResourceSNS(),
				"scaleway_mnq_sns_credentials":                 mnq.ResourceSNSCredentials(),
				"scaleway_mnq_sns_topic":                       mnq.ResourceSNSTopic(),
				"scaleway_mnq_sns_topic_subscription":          mnq.ResourceSNSTopicSubscription(),
				"scaleway_mnq_sqs":                             mnq.ResourceSQS(),
				"scaleway_mnq_sqs_credentials":                 mnq.ResourceSQSCredentials(),
				"scaleway_mnq_sqs_queue":                       mnq.ResourceSQSQueue(),
				"scaleway_mongodb_instance":                    mongodb.ResourceInstance(),
				"scaleway_mongodb_snapshot":                    mongodb.ResourceSnapshot(),
				"scaleway_object":                              object.ResourceObject(),
				"scaleway_object_bucket":                       object.ResourceBucket(),
				"scaleway_object_bucket_acl":                   object.ResourceBucketACL(),
				"scaleway_object_bucket_lock_configuration":    object.ResourceLockConfiguration(),
				"scaleway_object_bucket_policy":                object.ResourceBucketPolicy(),
				"scaleway_object_bucket_website_configuration": object.ResourceBucketWebsiteConfiguration(),
				"scaleway_rdb_acl":                             rdb.ResourceACL(),
				"scaleway_rdb_database":                        rdb.ResourceDatabase(),
				"scaleway_rdb_database_backup":                 rdb.ResourceDatabaseBackup(),
				"scaleway_rdb_instance":                        rdb.ResourceInstance(),
				"scaleway_rdb_privilege":                       rdb.ResourcePrivilege(),
				"scaleway_rdb_read_replica":                    rdb.ResourceReadReplica(),
				"scaleway_rdb_user":                            rdb.ResourceUser(),
				"scaleway_redis_cluster":                       redis.ResourceCluster(),
				"scaleway_registry_namespace":                  registry.ResourceNamespace(),
				"scaleway_sdb_sql_database":                    sdb.ResourceDatabase(),
				"scaleway_secret":                              secret.ResourceSecret(),
				"scaleway_secret_version":                      secret.ResourceVersion(),
				"scaleway_tem_domain":                          tem.ResourceDomain(),
				"scaleway_tem_domain_validation":               tem.ResourceDomainValidation(),
				"scaleway_tem_webhook":                         tem.ResourceWebhook(),
				"scaleway_vpc":                                 vpc.ResourceVPC(),
				"scaleway_vpc_gateway_network":                 vpcgw.ResourceNetwork(),
				"scaleway_vpc_private_network":                 vpc.ResourcePrivateNetwork(),
				"scaleway_vpc_public_gateway":                  vpcgw.ResourcePublicGateway(),
				"scaleway_vpc_public_gateway_dhcp":             vpcgw.ResourceDHCP(),
				"scaleway_vpc_public_gateway_dhcp_reservation": vpcgw.ResourceDHCPReservation(),
				"scaleway_vpc_public_gateway_ip":               vpcgw.ResourceIP(),
				"scaleway_vpc_public_gateway_ip_reverse_dns":   vpcgw.ResourceIPReverseDNS(),
				"scaleway_vpc_public_gateway_pat_rule":         vpcgw.ResourcePATRule(),
				"scaleway_vpc_route":                           vpc.ResourceRoute(),
				"scaleway_webhosting":                          webhosting.ResourceWebhosting(),
			},

			DataSourcesMap: map[string]*schema.Resource{
				"scaleway_account_project":                     account.DataSourceProject(),
				"scaleway_account_ssh_key":                     iam.DataSourceSSHKey(),
				"scaleway_availability_zones":                  az.DataSourceAvailabilityZones(),
				"scaleway_baremetal_offer":                     baremetal.DataSourceOffer(),
				"scaleway_baremetal_option":                    baremetal.DataSourceOption(),
				"scaleway_baremetal_os":                        baremetal.DataSourceOS(),
				"scaleway_baremetal_server":                    baremetal.DataSourceServer(),
				"scaleway_billing_consumptions":                billing.DataSourceConsumptions(),
				"scaleway_billing_invoices":                    billing.DataSourceInvoices(),
				"scaleway_block_snapshot":                      block.DataSourceSnapshot(),
				"scaleway_block_volume":                        block.DataSourceVolume(),
				"scaleway_cockpit":                             cockpit.DataSourceCockpit(),
				"scaleway_cockpit_plan":                        cockpit.DataSourcePlan(),
				"scaleway_cockpit_source":                      cockpit.DataSourceCockpitSource(),
				"scaleway_config":                              scwconfig.DataSourceConfig(),
				"scaleway_container":                           container.DataSourceContainer(),
				"scaleway_container_namespace":                 container.DataSourceNamespace(),
				"scaleway_domain_record":                       domain.DataSourceRecord(),
				"scaleway_domain_zone":                         domain.DataSourceZone(),
				"scaleway_flexible_ip":                         flexibleip.DataSourceFlexibleIP(),
				"scaleway_flexible_ips":                        flexibleip.DataSourceFlexibleIPs(),
				"scaleway_function":                            function.DataSourceFunction(),
				"scaleway_function_namespace":                  function.DataSourceNamespace(),
				"scaleway_iam_application":                     iam.DataSourceApplication(),
				"scaleway_iam_group":                           iam.DataSourceGroup(),
				"scaleway_iam_ssh_key":                         iam.DataSourceSSHKey(),
				"scaleway_iam_user":                            iam.DataSourceUser(),
				"scaleway_iam_api_key":                         iam.DataSourceAPIKey(),
				"scaleway_instance_image":                      instance.DataSourceImage(),
				"scaleway_instance_ip":                         instance.DataSourceIP(),
				"scaleway_instance_placement_group":            instance.DataSourcePlacementGroup(),
				"scaleway_instance_private_nic":                instance.DataSourcePrivateNIC(),
				"scaleway_instance_security_group":             instance.DataSourceSecurityGroup(),
				"scaleway_instance_server":                     instance.DataSourceServer(),
				"scaleway_instance_servers":                    instance.DataSourceServers(),
				"scaleway_instance_snapshot":                   instance.DataSourceSnapshot(),
				"scaleway_instance_volume":                     instance.DataSourceVolume(),
				"scaleway_iot_device":                          iot.DataSourceDevice(),
				"scaleway_iot_hub":                             iot.DataSourceHub(),
				"scaleway_ipam_ip":                             ipam.DataSourceIP(),
				"scaleway_ipam_ips":                            ipam.DataSourceIPs(),
				"scaleway_k8s_cluster":                         k8s.DataSourceCluster(),
				"scaleway_k8s_pool":                            k8s.DataSourcePool(),
				"scaleway_k8s_version":                         k8s.DataSourceVersion(),
				"scaleway_lb":                                  lb.DataSourceLb(),
				"scaleway_lb_acls":                             lb.DataSourceACLs(),
				"scaleway_lb_backend":                          lb.DataSourceBackend(),
				"scaleway_lb_backends":                         lb.DataSourceBackends(),
				"scaleway_lb_certificate":                      lb.DataSourceCertificate(),
				"scaleway_lb_frontend":                         lb.DataSourceFrontend(),
				"scaleway_lb_frontends":                        lb.DataSourceFrontends(),
				"scaleway_lb_ip":                               lb.DataSourceIP(),
				"scaleway_lb_ips":                              lb.DataSourceIPs(),
				"scaleway_lb_route":                            lb.DataSourceRoute(),
				"scaleway_lb_routes":                           lb.DataSourceRoutes(),
				"scaleway_lbs":                                 lb.DataSourceLbs(),
				"scaleway_marketplace_image":                   marketplace.DataSourceImage(),
				"scaleway_mnq_sqs":                             mnq.DataSourceSQS(),
				"scaleway_mnq_sns":                             mnq.DataSourceSNS(),
				"scaleway_mongodb_instance":                    mongodb.DataSourceInstance(),
				"scaleway_object_bucket":                       object.DataSourceBucket(),
				"scaleway_object_bucket_policy":                object.DataSourceBucketPolicy(),
				"scaleway_rdb_acl":                             rdb.DataSourceACL(),
				"scaleway_rdb_database":                        rdb.DataSourceDatabase(),
				"scaleway_rdb_database_backup":                 rdb.DataSourceDatabaseBackup(),
				"scaleway_rdb_instance":                        rdb.DataSourceInstance(),
				"scaleway_rdb_privilege":                       rdb.DataSourcePrivilege(),
				"scaleway_redis_cluster":                       redis.DataSourceCluster(),
				"scaleway_registry_image":                      registry.DataSourceImage(),
				"scaleway_registry_namespace":                  registry.DataSourceNamespace(),
				"scaleway_registry_image_tag":                  registry.DataSourceImageTag(),
				"scaleway_secret":                              secret.DataSourceSecret(),
				"scaleway_secret_version":                      secret.DataSourceVersion(),
				"scaleway_tem_domain":                          tem.DataSourceDomain(),
				"scaleway_vpc":                                 vpc.DataSourceVPC(),
				"scaleway_vpc_gateway_network":                 vpcgw.DataSourceNetwork(),
				"scaleway_vpc_private_network":                 vpc.DataSourcePrivateNetwork(),
				"scaleway_vpc_public_gateway":                  vpcgw.DataSourceVPCPublicGateway(),
				"scaleway_vpc_public_gateway_dhcp":             vpcgw.DataSourceDHCP(),
				"scaleway_vpc_public_gateway_dhcp_reservation": vpcgw.DataSourceDHCPReservation(),
				"scaleway_vpc_public_gateway_ip":               vpcgw.DataSourceIP(),
				"scaleway_vpc_public_gateway_pat_rule":         vpcgw.DataSourcePATRule(),
				"scaleway_vpc_routes":                          vpc.DataSourceRoutes(),
				"scaleway_vpcs":                                vpc.DataSourceVPCs(),
				"scaleway_webhosting":                          webhosting.DataSourceWebhosting(),
				"scaleway_webhosting_offer":                    webhosting.DataSourceOffer(),
			},
		}

		addBetaResources(p)

		p.ConfigureContextFunc = func(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
			terraformVersion := p.TerraformVersion

			// If we provide meta in config use it. This is useful for tests
			if config.Meta != nil {
				return config.Meta, nil
			}

			m, err := meta.NewMeta(ctx, &meta.Config{
				ProviderSchema:   data,
				TerraformVersion: terraformVersion,
			})
			if err != nil {
				return nil, diag.FromErr(err)
			}
			return m, nil
		}

		return p
	}
}

//gocyclo:ignore
