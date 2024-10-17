package instance

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	block "github.com/scaleway/scaleway-sdk-go/api/block/v1alpha1"
	instanceSDK "github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	ipamAPI "github.com/scaleway/scaleway-sdk-go/api/ipam/v1"
	"github.com/scaleway/scaleway-sdk-go/api/marketplace/v2"
	"github.com/scaleway/scaleway-sdk-go/scw"
	scwvalidation "github.com/scaleway/scaleway-sdk-go/validation"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/cdf"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/dsf"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/regional"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/zonal"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/meta"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/account"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/ipam"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/vpc"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/types"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/verify"
)

func ResourceServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceInstanceServerCreate,
		ReadContext:   ResourceInstanceServerRead,
		UpdateContext: ResourceInstanceServerUpdate,
		DeleteContext: ResourceInstanceServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create:  schema.DefaultTimeout(DefaultInstanceServerWaitTimeout),
			Read:    schema.DefaultTimeout(DefaultInstanceServerWaitTimeout),
			Update:  schema.DefaultTimeout(DefaultInstanceServerWaitTimeout),
			Delete:  schema.DefaultTimeout(DefaultInstanceServerWaitTimeout),
			Default: schema.DefaultTimeout(DefaultInstanceServerWaitTimeout),
		},
		SchemaVersion: 0,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the server",
			},
			"image": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "The UUID or the label of the base image used by the server",
				DiffSuppressFunc: dsf.Locality,
				ExactlyOneOf:     []string{"image", "root_volume.0.volume_id"},
			},
			"type": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "The instanceSDK type of the server", // TODO: link to scaleway pricing in the doc
				DiffSuppressFunc: dsf.IgnoreCase,
			},
			"replace_on_type_change": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Delete and re-create server if type change",
			},
			"tags": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "The tags associated with the server",
			},
			"security_group_id": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: dsf.Locality,
				Description:      "The security group the server is attached to",
			},
			"placement_group_id": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: dsf.Locality,
				Description:      "The placement group the server is attached to",
			},
			"placement_group_policy_respected": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "True when the placement group policy is respected",
			},
			"root_volume": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Root volume attached to the server on creation",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the root volume",
						},
						"size_in_gb": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Size of the root volume in gigabytes",
						},
						"volume_type": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ForceNew:         true,
							Description:      "Volume type of the root volume",
							ValidateDiagFunc: verify.ValidateEnum[instanceSDK.VolumeVolumeType](),
						},
						"delete_on_termination": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Force deletion of the root volume on instanceSDK termination",
						},
						"boot": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Set the volume where the boot the server",
						},
						"volume_id": {
							Type:         schema.TypeString,
							Computed:     true,
							Optional:     true,
							Description:  "Volume ID of the root volume",
							ExactlyOneOf: []string{"image", "root_volume.0.volume_id"},
						},
						"sbs_iops": {
							Type:        schema.TypeInt,
							Computed:    true,
							Optional:    true,
							Description: "SBS Volume IOPS, only with volume_type as sbs_volume",
						},
					},
				},
			},
			"additional_volume_ids": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: verify.IsUUIDorUUIDWithLocality(),
					DiffSuppressFunc: dsf.Locality,
				},
				Optional:    true,
				Description: "The additional volumes attached to the server",
			},
			"enable_ipv6": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Determines if IPv6 is enabled for the server",
				Deprecated:  "Please use a scaleway_instance_ip with a `routed_ipv6` type",
				DiffSuppressFunc: func(_, _, _ string, d *schema.ResourceData) bool {
					// routed_ip enabled servers already support enable_ipv6. Let's ignore this argument if it is.
					routedIPEnabled := types.GetBool(d, "routed_ip_enabled")
					if routedIPEnabled == nil || routedIPEnabled.(bool) {
						return true
					}

					return false
				},
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Scaleway internal IP address of the server",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The public IPv4 address of the server",
				Deprecated:  "Use public_ips instead",
			},
			"ip_id": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "The ID of the reserved IP for the server",
				DiffSuppressFunc: dsf.Locality,
				ConflictsWith:    []string{"ip_ids"},
			},
			"ip_ids": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"ip_id"},
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					Description:      "ID of the reserved IP for the server",
					DiffSuppressFunc: dsf.Locality,
				},
			},
			"ipv6_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The default public IPv6 address routed to the server.",
				Deprecated:  "Please use a scaleway_instance_ip with a `routed_ipv6` type",
			},
			"ipv6_gateway": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IPv6 gateway address",
				Deprecated:  "Please use a scaleway_instance_ip with a `routed_ipv6` type",
			},
			"ipv6_prefix_length": {
				Type:        schema.TypeInt,
				Computed:    true,
				Deprecated:  "Please use a scaleway_instance_ip with a `routed_ipv6` type",
				Description: "The IPv6 prefix length routed to the server.",
			},
			"enable_dynamic_ip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable dynamic IP on the server",
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     InstanceServerStateStarted,
				Description: "The state of the server should be: started, stopped, standby",
				ValidateFunc: validation.StringInSlice([]string{
					InstanceServerStateStarted,
					InstanceServerStateStopped,
					InstanceServerStateStandby,
				}, false),
			},
			"boot_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "The boot type of the server",
				Default:          instanceSDK.BootTypeLocal,
				ValidateDiagFunc: verify.ValidateEnum[instanceSDK.BootType](),
			},
			"bootscript_id": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "ID of the target bootscript (set boot_type to bootscript)",
				ValidateDiagFunc: verify.IsUUID(),
				Deprecated:       "bootscript is not supported anymore.",
			},
			"cloud_init": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The cloud init script associated with this server",
				ValidateFunc: validation.StringLenBetween(0, 127998),
			},
			"user_data": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Description: "The user data associated with the server", // TODO: document reserved keys (`cloud-init`)
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"private_network": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    8,
				Description: "List of private network to connect with your instanceSDK",
				Elem: &schema.Resource{
					Timeouts: &schema.ResourceTimeout{
						Default: schema.DefaultTimeout(defaultInstancePrivateNICWaitTimeout),
					},
					Schema: map[string]*schema.Schema{
						"pn_id": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: verify.IsUUIDorUUIDWithLocality(),
							Description:      "The Private Network ID",
							DiffSuppressFunc: dsf.Locality,
						},
						// Computed
						"mac_address": {
							Type:        schema.TypeString,
							Description: "MAC address of the NIC",
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The private NIC state",
						},
						"pnic_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the NIC",
						},
						"zone": zonal.Schema(),
					},
				},
			},
			"public_ips": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "List of public IPs attached to your instanceSDK",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the IP",
						},
						"address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP Address",
						},
					},
				},
			},
			"private_ips": {
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
			"routed_ip_enabled": {
				Type:        schema.TypeBool,
				Description: "If server supports routed IPs, default to true",
				Optional:    true,
				Computed:    true,
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					if i == nil {
						return nil
					}
					if !i.(bool) {
						return diag.Diagnostics{{
							Severity:      diag.Error,
							Summary:       "NAT IPs are not supported anymore",
							Detail:        "Remove explicit disabling, enable it or downgrade terraform.\nLearn more about migration: https://www.scaleway.com/en/docs/compute/instances/how-to/migrate-routed-ips/",
							AttributePath: path,
						}}
					}

					return nil
				},
				Deprecated: "Routed IP is the default configuration, it should always be true",
			},
			"zone":            zonal.Schema(),
			"organization_id": account.OrganizationIDSchema(),
			"project_id":      account.ProjectIDSchema(),
		},
		CustomizeDiff: customdiff.All(
			cdf.LocalityCheck(
				"placement_group_id",
				"additional_volume_ids.#",
				"ip_id",
			),
			customDiffInstanceServerType,
			customDiffInstanceServerImage,
		),
	}
}

//gocyclo:ignore
func ResourceInstanceServerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api, zone, err := instanceAndBlockAPIWithZone(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	////
	// Create the server
	////

	commercialType := d.Get("type").(string)

	imageUUID := locality.ExpandID(d.Get("image"))
	if imageUUID != "" && !scwvalidation.IsUUID(imageUUID) {
		// Replace dashes with underscores ubuntu-focal -> ubuntu_focal
		imageLabel := formatImageLabel(imageUUID)

		marketPlaceAPI := marketplace.NewAPI(meta.ExtractScwClient(m))
		image, err := marketPlaceAPI.GetLocalImageByLabel(&marketplace.GetLocalImageByLabelRequest{
			CommercialType: commercialType,
			Zone:           zone,
			ImageLabel:     imageLabel,
			Type:           volumeTypeToMarketplaceFilter(d.Get("root_volume.0.volume_type")),
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not get image '%s': %s", zonal.NewID(zone, imageLabel), err))
		}
		imageUUID = image.ID
	}

	req := &instanceSDK.CreateServerRequest{
		Zone:              zone,
		Name:              types.ExpandOrGenerateString(d.Get("name"), "srv"),
		Project:           types.ExpandStringPtr(d.Get("project_id")),
		Image:             imageUUID,
		CommercialType:    commercialType,
		SecurityGroup:     types.ExpandStringPtr(zonal.ExpandID(d.Get("security_group_id")).ID),
		DynamicIPRequired: scw.BoolPtr(d.Get("enable_dynamic_ip").(bool)),
		Tags:              types.ExpandStrings(d.Get("tags")),
		RoutedIPEnabled:   types.ExpandBoolPtr(types.GetBool(d, "routed_ip_enabled")),
	}

	enableIPv6, ok := d.GetOk("enable_ipv6")
	if ok {
		req.EnableIPv6 = scw.BoolPtr(enableIPv6.(bool)) //nolint:staticcheck
	}

	if bootType, ok := d.GetOk("boot_type"); ok {
		bootType := instanceSDK.BootType(bootType.(string))
		req.BootType = &bootType
	}

	if ipID, ok := d.GetOk("ip_id"); ok {
		req.PublicIP = types.ExpandStringPtr(zonal.ExpandID(ipID).ID) //nolint:staticcheck
	}

	if ipIDs, ok := d.GetOk("ip_ids"); ok {
		req.PublicIPs = types.ExpandSliceIDsPtr(ipIDs)
	}

	if placementGroupID, ok := d.GetOk("placement_group_id"); ok {
		req.PlacementGroup = types.ExpandStringPtr(zonal.ExpandID(placementGroupID).ID)
	}

	serverType := getServerType(ctx, api.API, req.Zone, req.CommercialType)
	if serverType == nil {
		return diag.Diagnostics{{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("could not find a server type associated with %s in zone %s", req.CommercialType, req.Zone),
			Detail:        "Ensure that the server type is correct, and that it does exist in this zone.",
			AttributePath: cty.GetAttrPath("type"),
		}}
	}

	req.Volumes = make(map[string]*instanceSDK.VolumeServerTemplate)
	rootVolume := d.Get("root_volume.0").(map[string]any)

	req.Volumes["0"] = prepareRootVolume(rootVolume, serverType, req.Image).VolumeTemplate()
	if raw, ok := d.GetOk("additional_volume_ids"); ok {
		for i, volumeID := range raw.([]interface{}) {
			// We have to get the volume to know whether it is a local or a block volume
			volumeTemplate, err := instanceServerAdditionalVolumeTemplate(api, zone, volumeID.(string))
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to get additional volume: %w", err))
			}
			req.Volumes[strconv.Itoa(i+1)] = volumeTemplate
		}
	}

	// Validate total local volume sizes.
	if err = validateLocalVolumeSizes(req.Volumes, serverType, req.CommercialType); err != nil {
		return diag.FromErr(err)
	}

	res, err := api.CreateServer(req, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(zonal.NewID(zone, res.Server.ID).String())

	_, err = waitForServer(ctx, api.API, zone, res.Server.ID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	////
	// Configure Block Volume
	////
	var diags diag.Diagnostics

	if iops, ok := d.GetOk("root_volume.0.sbs_iops"); ok {
		updateDiags := ResourceInstanceServerUpdateRootVolumeIOPS(ctx, api, zone, res.Server.ID, types.ExpandUint32Ptr(iops))
		if len(updateDiags) > 0 {
			diags = append(diags, updateDiags...)
		}
	}

	////
	// Set user data
	////
	userDataRequests := &instanceSDK.SetAllServerUserDataRequest{
		Zone:     zone,
		ServerID: res.Server.ID,
		UserData: make(map[string]io.Reader),
	}

	if rawUserData, ok := d.GetOk("user_data"); ok {
		for key, value := range rawUserData.(map[string]interface{}) {
			userDataRequests.UserData[key] = bytes.NewBufferString(value.(string))
		}
	}

	// cloud init script is set in user data
	if cloudInit, ok := d.GetOk("cloud_init"); ok {
		userDataRequests.UserData["cloud-init"] = bytes.NewBufferString(cloudInit.(string))
	}

	if len(userDataRequests.UserData) > 0 {
		_, err := waitForServer(ctx, api.API, zone, res.Server.ID, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}

		err = api.SetAllServerUserData(userDataRequests)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	targetState, err := serverStateExpand(d.Get("state").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	err = reachState(ctx, api, zone, res.Server.ID, targetState)
	if err != nil {
		return diag.FromErr(err)
	}

	////
	// Private Network
	////
	if rawPNICs, ok := d.GetOk("private_network"); ok {
		vpcAPI, err := vpc.NewAPI(m)
		if err != nil {
			return diag.FromErr(err)
		}
		pnRequest, err := preparePrivateNIC(ctx, rawPNICs, res.Server, vpcAPI)
		if err != nil {
			return diag.FromErr(err)
		}
		// compute attachment
		for _, q := range pnRequest {
			_, err := waitForServer(ctx, api.API, zone, res.Server.ID, d.Timeout(schema.TimeoutCreate))
			if err != nil {
				return diag.FromErr(err)
			}

			pn, err := api.CreatePrivateNIC(q, scw.WithContext(ctx))
			if err != nil {
				return diag.FromErr(err)
			}
			tflog.Debug(ctx, fmt.Sprintf("private network created (ID: %s, status: %s)", pn.PrivateNic.ID, pn.PrivateNic.State))

			_, err = waitForPrivateNIC(ctx, api.API, zone, res.Server.ID, pn.PrivateNic.ID, d.Timeout(schema.TimeoutCreate))
			if err != nil {
				return diag.FromErr(err)
			}

			_, err = waitForMACAddress(ctx, api.API, zone, res.Server.ID, pn.PrivateNic.ID, d.Timeout(schema.TimeoutCreate))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return append(diags, ResourceInstanceServerRead(ctx, d, m)...)
}

func errorCheck(err error, message string) bool {
	return strings.Contains(err.Error(), message)
}

//gocyclo:ignore
func ResourceInstanceServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api, zone, id, err := instanceAndBlockAPIWithZoneAndID(m, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	server, err := waitForServer(ctx, api.API, zone, id, d.Timeout(schema.TimeoutRead))
	if err != nil {
		if errorCheck(err, "is not found") {
			log.Printf("[WARN] instance %s not found droping from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	////
	// Read Server
	////

	if err == nil {
		state, err := serverStateFlatten(server.State)
		if err != nil {
			return diag.FromErr(err)
		}
		_ = d.Set("state", state)
		_ = d.Set("zone", string(zone))
		_ = d.Set("name", server.Name)
		_ = d.Set("boot_type", server.BootType)

		_ = d.Set("type", server.CommercialType)
		if len(server.Tags) > 0 {
			_ = d.Set("tags", server.Tags)
		}
		_ = d.Set("security_group_id", zonal.NewID(zone, server.SecurityGroup.ID).String())
		// EnableIPv6 is deprecated
		_ = d.Set("enable_ipv6", server.EnableIPv6) //nolint:staticcheck
		_ = d.Set("enable_dynamic_ip", server.DynamicIPRequired)
		_ = d.Set("organization_id", server.Organization)
		_ = d.Set("project_id", server.Project)
		_ = d.Set("routed_ip_enabled", server.RoutedIPEnabled) //nolint:staticcheck

		// Image could be empty in an import context.
		image := regional.ExpandID(d.Get("image").(string))
		if server.Image != nil && (image.ID == "" || scwvalidation.IsUUID(image.ID)) {
			_ = d.Set("image", zonal.NewID(zone, server.Image.ID).String())
		}

		if server.PlacementGroup != nil {
			_ = d.Set("placement_group_id", zonal.NewID(zone, server.PlacementGroup.ID).String())
			_ = d.Set("placement_group_policy_respected", server.PlacementGroup.PolicyRespected)
		}

		if server.PrivateIP != nil {
			_ = d.Set("private_ip", types.FlattenStringPtr(server.PrivateIP))
		}

		if _, hasIPID := d.GetOk("ip_id"); server.PublicIP != nil && hasIPID { //nolint:staticcheck
			if !server.PublicIP.Dynamic { //nolint:staticcheck
				_ = d.Set("ip_id", zonal.NewID(zone, server.PublicIP.ID).String()) //nolint:staticcheck
			} else {
				_ = d.Set("ip_id", "")
			}
		} else {
			_ = d.Set("ip_id", "")
		}

		if server.PublicIP != nil { //nolint:staticcheck
			_ = d.Set("public_ip", server.PublicIP.Address.String()) //nolint:staticcheck
			d.SetConnInfo(map[string]string{
				"type": "ssh",
				"host": server.PublicIP.Address.String(), //nolint:staticcheck
			})
		} else {
			_ = d.Set("public_ip", "")
			d.SetConnInfo(nil)
		}

		if len(server.PublicIPs) > 0 {
			_ = d.Set("public_ips", flattenServerPublicIPs(server.Zone, server.PublicIPs))
		} else {
			_ = d.Set("public_ips", []interface{}{})
		}

		if _, hasIPIDs := d.GetOk("ip_ids"); hasIPIDs {
			_ = d.Set("ip_ids", flattenServerIPIDs(server.PublicIPs))
		} else {
			_ = d.Set("ip_ids", []interface{}{})
		}

		if server.IPv6 != nil { //nolint:staticcheck
			_ = d.Set("ipv6_address", server.IPv6.Address.String()) //nolint:staticcheck
			_ = d.Set("ipv6_gateway", server.IPv6.Gateway.String()) //nolint:staticcheck
			prefixLength, err := strconv.Atoi(server.IPv6.Netmask)  //nolint:staticcheck
			if err != nil {
				return diag.FromErr(err)
			}
			_ = d.Set("ipv6_prefix_length", prefixLength)
		} else {
			_ = d.Set("ipv6_address", nil)
			_ = d.Set("ipv6_gateway", nil)
			_ = d.Set("ipv6_prefix_length", nil)
		}

		var additionalVolumesIDs []string
		for i, volume := range sortVolumeServer(server.Volumes) {
			if i == 0 {
				rootVolume := map[string]interface{}{}

				vs, ok := d.Get("root_volume").([]map[string]interface{})
				if ok && len(vs) > 0 {
					rootVolume = vs[0]
				}

				vol, err := api.GetUnknownVolume(&GetUnknownVolumeRequest{
					VolumeID: volume.ID,
					Zone:     volume.Zone,
				})
				if err != nil {
					return diag.FromErr(fmt.Errorf("failed to read instance volume %s: %w", volume.ID, err))
				}

				rootVolume["volume_id"] = zonal.NewID(zone, vol.ID).String()
				if vol.Size != nil {
					rootVolume["size_in_gb"] = int(uint64(*vol.Size) / gb)
				} else {
					rootVolume["size_in_gb"] = int(uint64(volume.Size) / gb)
				}
				if vol.IsBlockVolume() {
					rootVolume["sbs_iops"] = types.FlattenUint32Ptr(vol.Iops)
				}
				_, rootVolumeAttributeSet := d.GetOk("root_volume") // Related to https://github.com/hashicorp/terraform-plugin-sdk/issues/142
				rootVolume["delete_on_termination"] = d.Get("root_volume.0.delete_on_termination").(bool) || !rootVolumeAttributeSet
				rootVolume["volume_type"] = volume.VolumeType
				rootVolume["boot"] = volume.Boot
				rootVolume["name"] = volume.Name

				_ = d.Set("root_volume", []map[string]interface{}{rootVolume})
			} else {
				additionalVolumesIDs = append(additionalVolumesIDs, zonal.NewID(zone, volume.ID).String())
			}
		}

		_ = d.Set("additional_volume_ids", additionalVolumesIDs)
		if len(additionalVolumesIDs) > 0 {
			_ = d.Set("additional_volume_ids", additionalVolumesIDs)
		}
		////
		// Read server user data
		////
		allUserData, _ := api.GetAllServerUserData(&instanceSDK.GetAllServerUserDataRequest{
			Zone:     zone,
			ServerID: id,
		}, scw.WithContext(ctx))

		userData := make(map[string]interface{})
		for key, value := range allUserData.UserData {
			userDataValue, err := io.ReadAll(value)
			if err != nil {
				return diag.FromErr(err)
			}
			// if key != "cloud-init" {
			userData[key] = string(userDataValue)
			//	} else {
			// _ = d.Set("cloud_init", string(userDataValue))
			// }
		}
		_ = d.Set("user_data", userData)

		////
		// Read server private networks
		////
		ph, err := newPrivateNICHandler(api.API, id, zone)
		if err != nil {
			return diag.FromErr(err)
		}

		// set private networks
		err = ph.set(d)
		if err != nil {
			return diag.FromErr(err)
		}

		var privateNICIDs []string
		for _, nic := range ph.privateNICsMap {
			privateNICIDs = append(privateNICIDs, nic.ID)
		}

		var allPrivateIPs []map[string]interface{}
		resourceType := ipamAPI.ResourceTypeInstancePrivateNic
		region, err := zone.Region()
		if err != nil {
			return diag.FromErr(err)
		}
		for _, nicID := range privateNICIDs {
			opts := &ipam.GetResourcePrivateIPsOptions{
				ResourceType: &resourceType,
				ResourceID:   &nicID,
			}
			privateIPs, err := ipam.GetResourcePrivateIPs(ctx, m, region, opts)
			if err != nil {
				return diag.FromErr(err)
			}
			if privateIPs != nil {
				allPrivateIPs = append(allPrivateIPs, privateIPs...)
			}
		}
		_ = d.Set("private_ips", allPrivateIPs)

		return nil
	}
	return nil
}

//gocyclo:ignore
func ResourceInstanceServerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api, zone, id, err := instanceAndBlockAPIWithZoneAndID(m, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	wantedState := d.Get("state").(string)
	isStopped := wantedState == InstanceServerStateStopped

	var warnings diag.Diagnostics

	server, err := waitForServer(ctx, api.API, zone, id, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.FromErr(err)
	}
	////
	// Construct UpdateServerRequest
	////
	serverShouldUpdate := false
	updateRequest := &instanceSDK.UpdateServerRequest{
		Zone:     zone,
		ServerID: server.ID,
	}

	if d.HasChange("name") {
		serverShouldUpdate = true
		updateRequest.Name = types.ExpandStringPtr(d.Get("name"))
	}

	if d.HasChange("tags") {
		serverShouldUpdate = true
		updateRequest.Tags = types.ExpandUpdatedStringsPtr(d.Get("tags"))
	}

	if d.HasChange("security_group_id") {
		serverShouldUpdate = true
		updateRequest.SecurityGroup = &instanceSDK.SecurityGroupTemplate{
			ID:   zonal.ExpandID(d.Get("security_group_id")).ID,
			Name: types.NewRandomName("sg"), // this value will be ignored by the API
		}
	}

	if d.HasChange("enable_ipv6") {
		serverShouldUpdate = true
		updateRequest.EnableIPv6 = scw.BoolPtr(d.Get("enable_ipv6").(bool))
	}

	if d.HasChange("enable_dynamic_ip") {
		serverShouldUpdate = true
		updateRequest.DynamicIPRequired = scw.BoolPtr(d.Get("enable_dynamic_ip").(bool))
	}

	if d.HasChanges("additional_volume_ids", "root_volume") {
		volumes, err := instanceServerVolumesUpdate(ctx, d, api, zone, isStopped)
		if err != nil {
			return diag.FromErr(err)
		}
		serverShouldUpdate = true
		updateRequest.Volumes = &volumes
	}

	if d.HasChange("placement_group_id") {
		serverShouldUpdate = true
		placementGroupID := zonal.ExpandID(d.Get("placement_group_id")).ID
		if placementGroupID == "" {
			updateRequest.PlacementGroup = &instanceSDK.NullableStringValue{Null: true}
		} else {
			if !isStopped {
				return diag.FromErr(errors.New("instanceSDK must be stopped to change placement group"))
			}
			updateRequest.PlacementGroup = &instanceSDK.NullableStringValue{Value: placementGroupID}
		}
	}

	////
	// Update reserved IP
	////
	if d.HasChange("ip_id") && !instanceIPHasMigrated(d) {
		server, err := waitForServer(ctx, api.API, zone, id, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}

		ipID := zonal.ExpandID(d.Get("ip_id")).ID
		// If an IP is already attached, and it's not a dynamic IP we detach it.
		if server.PublicIP != nil && !server.PublicIP.Dynamic { //nolint:staticcheck
			_, err = api.UpdateIP(&instanceSDK.UpdateIPRequest{
				Zone:   zone,
				IP:     server.PublicIP.ID, //nolint:staticcheck
				Server: &instanceSDK.NullableStringValue{Null: true},
			})
			if err != nil {
				return diag.FromErr(err)
			}
			// we wait to ensure to not detach the new ip.
			_, err := waitForServer(ctx, api.API, zone, id, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.FromErr(err)
			}
		}
		// If a new IP is provided, we attach it to the server
		if ipID != "" {
			_, err := waitForServer(ctx, api.API, zone, id, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.FromErr(err)
			}

			_, err = api.UpdateIP(&instanceSDK.UpdateIPRequest{
				Zone:   zone,
				IP:     ipID,
				Server: &instanceSDK.NullableStringValue{Value: id},
			}, scw.WithContext(ctx))
			if err != nil {
				return diag.FromErr(err)
			}

			_, err = waitForServer(ctx, api.API, zone, id, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("ip_ids") {
		err := ResourceInstanceServerUpdateIPs(ctx, d, api.API, zone, id)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("boot_type") {
		bootType := instanceSDK.BootType(d.Get("boot_type").(string))
		serverShouldUpdate = true
		updateRequest.BootType = &bootType
		if !isStopped {
			warnings = append(warnings, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "instanceSDK may need to be rebooted to use the new boot type",
			})
		}
	}

	////
	// Update server user data
	////
	if d.HasChanges("user_data") {
		userDataRequests := &instanceSDK.SetAllServerUserDataRequest{
			Zone:     zone,
			ServerID: id,
			UserData: make(map[string]io.Reader),
		}

		if allUserData, ok := d.GetOk("user_data"); ok {
			userDataMap := allUserData.(map[string]interface{})
			for key, value := range userDataMap {
				userDataRequests.UserData[key] = bytes.NewBufferString(value.(string))
			}
			if !isStopped && d.HasChange("user_data.cloud-init") {
				warnings = append(warnings, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "instanceSDK may need to be rebooted to use the new cloud init config",
				})
			}
		}

		_, err := waitForServer(ctx, api.API, zone, id, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}

		err = api.SetAllServerUserData(userDataRequests)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	////
	// Update server private network
	////
	if d.HasChanges("private_network") {
		ph, err := newPrivateNICHandler(api.API, id, zone)
		if err != nil {
			diag.FromErr(err)
		}
		if raw, ok := d.GetOk("private_network"); ok {
			// retrieve all current private network interfaces
			for index := range raw.([]interface{}) {
				pnKey := fmt.Sprintf("private_network.%d.pn_id", index)
				if d.HasChange(pnKey) {
					o, n := d.GetChange(pnKey)
					if !cmp.Equal(n, o) {
						_, err := waitForServer(ctx, api.API, zone, id, d.Timeout(schema.TimeoutUpdate))
						if err != nil {
							return diag.FromErr(err)
						}

						err = ph.detach(ctx, o, d.Timeout(schema.TimeoutUpdate))
						if err != nil {
							return diag.FromErr(err)
						}
						err = ph.attach(ctx, n, d.Timeout(schema.TimeoutUpdate))
						if err != nil {
							return diag.FromErr(err)
						}
					}
				}
			}
		} else {
			// retrieve old private network config
			o, _ := d.GetChange("private_network")
			for _, raw := range o.([]interface{}) {
				pn, pnExist := raw.(map[string]interface{})
				if pnExist {
					_, err := waitForServer(ctx, api.API, zone, id, d.Timeout(schema.TimeoutUpdate))
					if err != nil {
						return diag.FromErr(err)
					}

					err = ph.detach(ctx, pn["pn_id"], d.Timeout(schema.TimeoutUpdate))
					if err != nil {
						return diag.FromErr(err)
					}
				}
			}
		}
	}
	////
	// Apply changes
	////

	if d.HasChange("state") {
		targetState, err := serverStateExpand(d.Get("state").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		// reach expected state
		err = reachState(ctx, api, zone, id, targetState)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if serverShouldUpdate {
		_, err = api.UpdateServer(updateRequest)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	_, err = waitForServer(ctx, api.API, zone, id, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("type") {
		err := ResourceInstanceServerMigrate(ctx, d, api, zone, id)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("routed_ip_enabled") {
		err := ResourceInstanceServerEnableRoutedIP(ctx, d, api.API, zone, id)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("root_volume.0.sbs_iops") {
		warnings = append(warnings, ResourceInstanceServerUpdateRootVolumeIOPS(ctx, api, zone, id, types.ExpandUint32Ptr(d.Get("root_volume.0.sbs_iops")))...)
	}

	return append(warnings, ResourceInstanceServerRead(ctx, d, m)...)
}

func ResourceInstanceServerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api, zone, id, err := instanceAndBlockAPIWithZoneAndID(m, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	// detach eip to ensure to free eip even if instanceSDK won't stop
	if ipID, ok := d.GetOk("ip_id"); ok {
		_, err := api.UpdateIP(&instanceSDK.UpdateIPRequest{
			Zone:   zone,
			IP:     zonal.ExpandID(ipID).ID,
			Server: &instanceSDK.NullableStringValue{Null: true},
		})
		if err != nil {
			log.Print("[WARN] Failed to detach eip of server")
		}
	}
	// Remove instanceSDK from placement group to free it even if instanceSDK won't stop
	if _, ok := d.GetOk("placement_group_id"); ok {
		_, err := api.UpdateServer(&instanceSDK.UpdateServerRequest{
			Zone:           zone,
			PlacementGroup: &instanceSDK.NullableStringValue{Null: true},
			ServerID:       id,
		})
		if err != nil {
			log.Print("[WARN] Failed remove server from instanceSDK group")
		}
	}
	// reach stopped state
	err = reachState(ctx, api, zone, id, instanceSDK.ServerStateStopped)
	if httperrors.Is404(err) {
		return nil
	}
	if err != nil {
		return diag.FromErr(err)
	}

	// Delete private-nic if managed by instance_server resource
	if raw, ok := d.GetOk("private_network"); ok {
		ph, err := newPrivateNICHandler(api.API, id, zone)
		if err != nil {
			return diag.FromErr(err)
		}

		for index := range raw.([]interface{}) {
			pnKey := fmt.Sprintf("private_network.%d.pn_id", index)
			pn := d.Get(pnKey)
			err := ph.detach(ctx, pn, d.Timeout(schema.TimeoutDelete))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	_, err = waitForServer(ctx, api.API, zone, id, d.Timeout(schema.TimeoutDelete))
	if err != nil && !httperrors.Is404(err) {
		return diag.FromErr(err)
	}

	err = api.DeleteServer(&instanceSDK.DeleteServerRequest{
		Zone:     zone,
		ServerID: id,
	}, scw.WithContext(ctx))
	if err != nil && !httperrors.Is404(err) {
		return diag.FromErr(err)
	}

	_, err = waitForServer(ctx, api.API, zone, id, d.Timeout(schema.TimeoutDelete))
	if err != nil && !httperrors.Is404(err) {
		return diag.FromErr(err)
	}

	// Related to https://github.com/hashicorp/terraform-plugin-sdk/issues/142
	_, rootVolumeAttributeSet := d.GetOk("root_volume")
	if d.Get("root_volume.0.delete_on_termination").(bool) || !rootVolumeAttributeSet {
		volumeID, volumeExist := d.GetOk("root_volume.0.volume_id")
		if !volumeExist {
			return diag.Errorf("volume ID not found")
		}
		err = api.DeleteVolume(&instanceSDK.DeleteVolumeRequest{
			Zone:     zone,
			VolumeID: locality.ExpandID(volumeID),
		})
		if err != nil && !httperrors.Is404(err) {
			return diag.FromErr(err)
		}
	}

	return nil
}

func instanceServerCanMigrate(api *instanceSDK.API, server *instanceSDK.Server, requestedType string) error {
	var localVolumeSize scw.Size

	for _, volume := range server.Volumes {
		if volume.VolumeType == instanceSDK.VolumeServerVolumeTypeLSSD {
			localVolumeSize += volume.Size
		}
	}

	serverType, err := api.GetServerType(&instanceSDK.GetServerTypeRequest{
		Zone: server.Zone,
		Name: requestedType,
	})
	if err != nil {
		return err
	}

	if serverType.VolumesConstraint != nil &&
		(localVolumeSize > serverType.VolumesConstraint.MaxSize) ||
		(localVolumeSize < serverType.VolumesConstraint.MinSize) {
		return fmt.Errorf("local volume total size does not respect type constraint, expected beteween (%dGB, %dGB), got %sGB",
			serverType.VolumesConstraint.MinSize/scw.GB,
			serverType.VolumesConstraint.MaxSize/scw.GB,
			localVolumeSize/scw.GB)
	}

	return nil
}

func customDiffInstanceServerType(_ context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	if !diff.HasChange("type") || diff.Id() == "" {
		return nil
	}

	if diff.Get("replace_on_type_change").(bool) {
		return diff.ForceNew("type")
	}

	instanceAPI, zone, id, err := NewAPIWithZoneAndID(meta, diff.Id())
	if err != nil {
		return err
	}

	_, newValue := diff.GetChange("type")
	newType := newValue.(string)

	resp, err := instanceAPI.GetServer(&instanceSDK.GetServerRequest{
		Zone:     zone,
		ServerID: id,
	})
	if err != nil {
		return fmt.Errorf("failed to check server type change: %w", err)
	}

	err = instanceServerCanMigrate(instanceAPI, resp.Server, newType)
	if err != nil {
		return fmt.Errorf("cannot change server type: %w", err)
	}

	return nil
}

func customDiffInstanceServerImage(ctx context.Context, diff *schema.ResourceDiff, m interface{}) error {
	if diff.Get("image") == "" || !diff.HasChange("image") || diff.Id() == "" {
		return nil
	}

	// We get the server to fetch the UUID of the image
	instanceAPI, zone, id, err := NewAPIWithZoneAndID(m, diff.Id())
	if err != nil {
		return err
	}
	server, err := instanceAPI.GetServer(&instanceSDK.GetServerRequest{
		Zone:     zone,
		ServerID: id,
	}, scw.WithContext(ctx))
	if err != nil {
		return err
	}

	// If 'image' field is defined by the user and server.Image is empty, we should create a new server
	if server.Server.Image == nil {
		return diff.ForceNew("image")
	}

	// We get the image as it is defined by the user
	image := regional.ExpandID(diff.Get("image").(string))
	if scwvalidation.IsUUID(image.ID) {
		if image.ID == zonal.ExpandID(server.Server.Image.ID).ID {
			return nil
		}
	}

	// If image is a label, we check that server.Image.ID matches the label in case the user has edited
	// the image with another tool.
	marketplaceAPI := marketplace.NewAPI(meta.ExtractScwClient(m))
	if err != nil {
		return err
	}
	marketplaceImage, err := marketplaceAPI.GetLocalImage(&marketplace.GetLocalImageRequest{
		LocalImageID: server.Server.Image.ID,
	}, scw.WithContext(ctx))
	if err != nil {
		// If UUID is not in marketplace, then it's an image change
		if httperrors.Is404(err) {
			return diff.ForceNew("image")
		}
		return err
	}
	if marketplaceImage.Label != image.ID {
		return diff.ForceNew("image")
	}
	return nil
}

func ResourceInstanceServerMigrate(ctx context.Context, d *schema.ResourceData, api *BlockAndInstanceAPI, zone scw.Zone, id string) error {
	server, err := waitForServer(ctx, api.API, zone, id, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return fmt.Errorf("failed to wait for server before changing server type: %w", err)
	}
	beginningState := server.State

	err = reachState(ctx, api, zone, id, instanceSDK.ServerStateStopped)
	if err != nil {
		return fmt.Errorf("failed to stop server before changing server type: %w", err)
	}

	_, err = api.UpdateServer(&instanceSDK.UpdateServerRequest{
		Zone:           zone,
		ServerID:       id,
		CommercialType: types.ExpandStringPtr(d.Get("type")),
	})
	if err != nil {
		return errors.New("failed to change server type server")
	}

	err = reachState(ctx, api, zone, id, beginningState)
	if err != nil {
		return fmt.Errorf("failed to start server after changing server type: %w", err)
	}

	return nil
}

func ResourceInstanceServerEnableRoutedIP(ctx context.Context, d *schema.ResourceData, instanceAPI *instanceSDK.API, zone scw.Zone, id string) error {
	server, err := waitForServer(ctx, instanceAPI, zone, id, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	_, err = instanceAPI.ServerAction(&instanceSDK.ServerActionRequest{
		Zone:     server.Zone,
		ServerID: server.ID,
		Action:   "enable_routed_ip",
	})
	if err != nil {
		return fmt.Errorf("failed to enable routed ip: %w", err)
	}

	_, err = waitForServer(ctx, instanceAPI, zone, id, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	return nil
}

func ResourceInstanceServerUpdateIPs(ctx context.Context, d *schema.ResourceData, instanceAPI *instanceSDK.API, zone scw.Zone, id string) error {
	server, err := waitForServer(ctx, instanceAPI, zone, id, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	schemaIPs := d.Get("ip_ids").([]interface{})
	requestedIPs := make(map[string]bool, len(schemaIPs))

	// Gather request IPs in a map
	for _, rawIP := range schemaIPs {
		requestedIPs[locality.ExpandID(rawIP)] = false
	}

	// Detach all IPs that are not requested and set to true the one that are already attached
	for _, ip := range server.PublicIPs {
		_, isRequested := requestedIPs[ip.ID]
		if isRequested {
			requestedIPs[ip.ID] = true
		} else {
			_, err := instanceAPI.UpdateIP(&instanceSDK.UpdateIPRequest{
				Zone: zone,
				IP:   ip.ID,
				Server: &instanceSDK.NullableStringValue{
					Null: true,
				},
			})
			if err != nil {
				return fmt.Errorf("failed to detach IP: %w", err)
			}
		}
	}

	// Attach all remaining IPs that are not attached
	for ipID, isAttached := range requestedIPs {
		if isAttached {
			continue
		}
		_, err := instanceAPI.UpdateIP(&instanceSDK.UpdateIPRequest{
			Zone: zone,
			IP:   ipID,
			Server: &instanceSDK.NullableStringValue{
				Value: server.ID,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to attach IP: %w", err)
		}
	}

	return nil
}

func ResourceInstanceServerUpdateRootVolumeIOPS(ctx context.Context, api *BlockAndInstanceAPI, zone scw.Zone, serverID string, iops *uint32) diag.Diagnostics {
	res, err := api.GetServer(&instanceSDK.GetServerRequest{
		Zone:     zone,
		ServerID: serverID,
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	rootVolume, exists := res.Server.Volumes["0"]
	if exists {
		_, err := api.blockAPI.UpdateVolume(&block.UpdateVolumeRequest{
			Zone:     zone,
			VolumeID: rootVolume.ID,
			PerfIops: iops,
		}, scw.WithContext(ctx))
		if err != nil {
			return diag.Diagnostics{{
				Severity:      diag.Warning,
				Summary:       "Failed to update root_volume iops",
				Detail:        err.Error(),
				AttributePath: cty.GetAttrPath("root_volume.0.sbs_iops"),
			}}
		}
	} else {
		return diag.Diagnostics{{
			Severity:      diag.Warning,
			Summary:       "Failed to find root_volume",
			Detail:        "Failed to update root_volume IOPS",
			AttributePath: cty.GetAttrPath("root_volume.0.sbs_iops"),
		}}
	}

	return nil
}

// instanceServerVolumesUpdate updates root_volume size and returns the list of volumes templates that should be updated for the server.
// It uses root_volume and additional_volume_ids to build the volumes templates.
func instanceServerVolumesUpdate(ctx context.Context, d *schema.ResourceData, api *BlockAndInstanceAPI, zone scw.Zone, serverIsStopped bool) (map[string]*instanceSDK.VolumeServerTemplate, error) {
	volumes := map[string]*instanceSDK.VolumeServerTemplate{}
	raw, hasAdditionalVolumes := d.GetOk("additional_volume_ids")

	if d.HasChange("root_volume.0.size_in_gb") {
		err := api.ResizeUnknownVolume(&ResizeUnknownVolumeRequest{
			VolumeID: zonal.ExpandID(d.Get("root_volume.0.volume_id")).ID,
			Zone:     zone,
			Size:     scw.SizePtr(scw.Size(d.Get("root_volume.0.size_in_gb").(int)) * scw.GB),
		}, scw.WithContext(ctx))
		if err != nil {
			return nil, err
		}
	}

	volumes["0"] = &instanceSDK.VolumeServerTemplate{
		ID:   scw.StringPtr(zonal.ExpandID(d.Get("root_volume.0.volume_id")).ID),
		Name: scw.StringPtr(types.NewRandomName("vol")), // name is ignored by the API, any name will work here
		Boot: types.ExpandBoolPtr(d.Get("root_volume.0.boot")),
	}

	if !hasAdditionalVolumes {
		raw = []interface{}{} // Set an empty list if not volumes exist
	}

	for i, volumeID := range raw.([]interface{}) {
		volumeHasChange := d.HasChange("additional_volume_ids." + strconv.Itoa(i))
		volume, err := api.GetUnknownVolume(&GetUnknownVolumeRequest{
			VolumeID: zonal.ExpandID(volumeID).ID,
			Zone:     zone,
		}, scw.WithContext(ctx))
		if err != nil {
			return nil, fmt.Errorf("failed to get updated volume: %w", err)
		}

		// local volumes can only be added when the server is stopped
		if volumeHasChange && !serverIsStopped && volume.IsLocal() && volume.IsAttached() {
			return nil, errors.New("instanceSDK must be stopped to change local volumes")
		}
		volumes[strconv.Itoa(i+1)] = volume.VolumeTemplate()
	}

	return volumes, nil
}
