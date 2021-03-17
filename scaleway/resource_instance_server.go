package scaleway

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/api/marketplace/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	scwvalidation "github.com/scaleway/scaleway-sdk-go/validation"
)

func resourceScalewayInstanceServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScalewayInstanceServerCreate,
		ReadContext:   resourceScalewayInstanceServerRead,
		UpdateContext: resourceScalewayInstanceServerUpdate,
		DeleteContext: resourceScalewayInstanceServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(defaultInstanceServerWaitTimeout),
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
				Required:         true,
				ForceNew:         true,
				Description:      "The UUID or the label of the base image used by the server",
				DiffSuppressFunc: diffSuppressFuncLocality,
			},
			"type": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "The instance type of the server", // TODO: link to scaleway pricing in the doc
				DiffSuppressFunc: diffSuppressFuncIgnoreCase,
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
				DiffSuppressFunc: diffSuppressFuncLocality,
				Description:      "The security group the server is attached to",
			},
			"placement_group_id": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: diffSuppressFuncLocality,
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
						"size_in_gb": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							ForceNew:    true, // todo: don't force new but stop server and create new volume instead
							Description: "Size of the root volume in gigabytes",
						},
						"delete_on_termination": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Force deletion of the root volume on instance termination",
						},
						"volume_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Volume ID of the root volume",
						},
					},
				},
			},
			"additional_volume_ids": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateFunc:     validationUUIDorUUIDWithLocality(),
					DiffSuppressFunc: diffSuppressFuncLocality,
				},
				Optional:    true,
				Description: "The additional volumes attached to the server",
			},
			"enable_ipv6": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Determines if IPv6 is enabled for the server",
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
			},
			"ip_id": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "The ID of the reserved IP for the server",
				DiffSuppressFunc: diffSuppressFuncLocality,
			},
			"ipv6_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The default public IPv6 address routed to the server.",
			},
			"ipv6_gateway": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IPv6 gateway address",
			},
			"ipv6_prefix_length": {
				Type:        schema.TypeInt,
				Computed:    true,
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
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The boot type of the server",
				Default:     instance.BootTypeLocal,
				ValidateFunc: validation.StringInSlice([]string{
					instance.BootTypeLocal.String(),
					instance.BootTypeRescue.String(),
					instance.BootTypeBootscript.String(),
				}, false),
			},
			"bootscript_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "ID of the target bootscript (set boot_type to bootscript)",
				ValidateFunc: validationUUID(),
			},
			"cloud_init": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The cloud init script associated with this server",
				ValidateFunc: validation.StringLenBetween(0, 127998),
			},
			"user_data": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The user data associated with the server", // TODO: document reserved keys (`cloud-init`)
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"zone":            zoneSchema(),
			"organization_id": organizationIDSchema(),
			"project_id":      projectIDSchema(),
		},
	}
}

func resourceScalewayInstanceServerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	instanceAPI, zone, err := instanceAPIWithZone(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	////
	// Create the server
	////

	commercialType := d.Get("type").(string)

	imageUUID := expandZonedID(d.Get("image")).ID
	if !scwvalidation.IsUUID(imageUUID) {
		marketPlaceAPI := marketplace.NewAPI(meta.(*Meta).scwClient)
		imageUUID, err = marketPlaceAPI.GetLocalImageIDByLabel(&marketplace.GetLocalImageIDByLabelRequest{
			CommercialType: commercialType,
			Zone:           zone,
			ImageLabel:     imageUUID,
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not get image '%s': %s", newZonedID(zone, imageUUID), err))
		}
	}

	req := &instance.CreateServerRequest{
		Zone:              zone,
		Name:              expandOrGenerateString(d.Get("name"), "srv"),
		Project:           expandStringPtr(d.Get("project_id")),
		Image:             imageUUID,
		CommercialType:    commercialType,
		EnableIPv6:        d.Get("enable_ipv6").(bool),
		SecurityGroup:     expandStringPtr(expandZonedID(d.Get("security_group_id")).ID),
		DynamicIPRequired: scw.BoolPtr(d.Get("enable_dynamic_ip").(bool)),
		Tags:              expandStrings(d.Get("tags")),
	}

	if definedZone, ok := d.GetOk("zone"); ok {
		zone = scw.Zone(definedZone.(string))
		req.Zone = zone
	}

	if bootScriptID, ok := d.GetOk("bootscript_id"); ok {
		req.Bootscript = expandStringPtr(bootScriptID)
	}

	if bootType, ok := d.GetOk("boot_type"); ok {
		bootType := instance.BootType(bootType.(string))
		req.BootType = &bootType
	}

	if ipID, ok := d.GetOk("ip_id"); ok {
		req.PublicIP = expandStringPtr(expandZonedID(ipID).ID)
	}

	if placementGroupID, ok := d.GetOk("placement_group_id"); ok {
		req.PlacementGroup = expandStringPtr(expandZonedID(placementGroupID).ID)
	}

	serverType := getServerType(instanceAPI, req.Zone, req.CommercialType)
	if serverType == nil {
		return diag.FromErr(fmt.Errorf("could not find a server type associated with %s", req.CommercialType))
	}

	req.Volumes = make(map[string]*instance.VolumeTemplate)
	if size, ok := d.GetOk("root_volume.0.size_in_gb"); ok {
		req.Volumes["0"] = &instance.VolumeTemplate{
			Size:       scw.Size(uint64(size.(int)) * gb),
			VolumeType: instance.VolumeVolumeTypeLSSD,
		}
	} else {
		// We had a local root volume if it is not already present
		req.Volumes["0"] = &instance.VolumeTemplate{
			Name:       newRandomName("vol"),
			VolumeType: instance.VolumeVolumeTypeLSSD,
			Size:       serverType.VolumesConstraint.MinSize,
		}
	}

	if raw, ok := d.GetOk("additional_volume_ids"); ok {
		for i, volumeID := range raw.([]interface{}) {
			// We have to get the volume to know whether it is a local or a block volume
			vol, err := instanceAPI.GetVolume(&instance.GetVolumeRequest{
				Zone:     zone,
				VolumeID: expandZonedID(volumeID).ID,
			})
			if err != nil {
				return diag.FromErr(err)
			}
			req.Volumes[strconv.Itoa(i+1)] = &instance.VolumeTemplate{
				ID:         vol.Volume.ID,
				Name:       vol.Volume.Name,
				VolumeType: vol.Volume.VolumeType,
				Size:       vol.Volume.Size,
			}
		}
	}

	// Validate total local volume sizes.
	if err = validateLocalVolumeSizes(req.Volumes, serverType, req.CommercialType); err != nil {
		return diag.FromErr(err)
	}

	// Sanitize the volume map to respect API schemas
	req.Volumes = sanitizeVolumeMap(req.Name, req.Volumes)

	res, err := instanceAPI.CreateServer(req, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newZonedID(zone, res.Server.ID).String())

	////
	// Set user data
	////
	userDataRequests := &instance.SetAllServerUserDataRequest{
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
		err = instanceAPI.SetAllServerUserData(userDataRequests)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	targetState, err := serverStateExpand(d.Get("state").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	err = reachState(ctx, instanceAPI, zone, res.Server.ID, targetState)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceScalewayInstanceServerRead(ctx, d, meta)
}

func resourceScalewayInstanceServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	instanceAPI, zone, ID, err := instanceAPIWithZoneAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	////
	// Read Server
	////
	response, err := instanceAPI.GetServer(&instance.GetServerRequest{
		Zone:     zone,
		ServerID: ID,
	}, scw.WithContext(ctx))
	if err != nil {
		if is404Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	state, err := serverStateFlatten(response.Server.State)
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("state", state)
	_ = d.Set("zone", string(zone))
	_ = d.Set("name", response.Server.Name)
	_ = d.Set("boot_type", response.Server.BootType)
	_ = d.Set("bootscript_id", response.Server.Bootscript.ID)
	_ = d.Set("type", response.Server.CommercialType)
	_ = d.Set("tags", response.Server.Tags)
	_ = d.Set("security_group_id", newZonedID(zone, response.Server.SecurityGroup.ID).String())
	_ = d.Set("enable_ipv6", response.Server.EnableIPv6)
	_ = d.Set("enable_dynamic_ip", response.Server.DynamicIPRequired)
	_ = d.Set("organization_id", response.Server.Organization)
	_ = d.Set("project_id", response.Server.Project)

	// Image could be empty in an import context.
	image := expandRegionalID(d.Get("image").(string))
	if response.Server.Image != nil && (image.ID == "" || scwvalidation.IsUUID(image.ID)) {
		// TODO: If image is a label, check that response.Server.Image.ID match the label.
		// It could be useful if the user edit the image with another tool.
		_ = d.Set("image", newZonedID(zone, response.Server.Image.ID).String())
	}

	if response.Server.PlacementGroup != nil {
		_ = d.Set("placement_group_id", newZonedID(zone, response.Server.PlacementGroup.ID).String())
		_ = d.Set("placement_group_policy_respected", response.Server.PlacementGroup.PolicyRespected)
	}

	if response.Server.PrivateIP != nil {
		_ = d.Set("private_ip", flattenStringPtr(response.Server.PrivateIP))
	}

	if response.Server.PublicIP != nil {
		_ = d.Set("public_ip", response.Server.PublicIP.Address.String())
		d.SetConnInfo(map[string]string{
			"type": "ssh",
			"host": response.Server.PublicIP.Address.String(),
		})
		if !response.Server.PublicIP.Dynamic {
			_ = d.Set("ip_id", newZonedID(zone, response.Server.PublicIP.ID).String())
		} else {
			_ = d.Set("ip_id", "")
		}
	} else {
		_ = d.Set("public_ip", "")
		_ = d.Set("ip_id", "")
		d.SetConnInfo(nil)
	}

	if response.Server.IPv6 != nil {
		_ = d.Set("ipv6_address", response.Server.IPv6.Address.String())
		_ = d.Set("ipv6_gateway", response.Server.IPv6.Gateway.String())
		prefixLength, err := strconv.Atoi(response.Server.IPv6.Netmask)
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
	for i, volume := range orderVolumes(response.Server.Volumes) {
		if i == 0 {
			rootVolume := map[string]interface{}{}

			vs, ok := d.Get("root_volume").([]map[string]interface{})
			if ok && len(vs) > 0 {
				rootVolume = vs[0]
			}

			rootVolume["volume_id"] = newZonedID(zone, volume.ID).String()
			rootVolume["size_in_gb"] = int(uint64(volume.Size) / gb)
			_, rootVolumeAttributeSet := d.GetOk("root_volume") // Related to https://github.com/hashicorp/terraform-plugin-sdk/issues/142
			rootVolume["delete_on_termination"] = d.Get("root_volume.0.delete_on_termination").(bool) || !rootVolumeAttributeSet

			_ = d.Set("root_volume", []map[string]interface{}{rootVolume})
		} else {
			additionalVolumesIDs = append(additionalVolumesIDs, newZonedID(zone, volume.ID).String())
		}
	}
	_ = d.Set("additional_volume_ids", additionalVolumesIDs)

	////
	// Read server user data
	////
	allUserData, _ := instanceAPI.GetAllServerUserData(&instance.GetAllServerUserDataRequest{
		Zone:     zone,
		ServerID: ID,
	}, scw.WithContext(ctx))

	userData := make(map[string]interface{})
	for key, value := range allUserData.UserData {
		userDataValue, err := ioutil.ReadAll(value)
		if err != nil {
			return diag.FromErr(err)
		}
		//if key != "cloud-init" {
		userData[key] = string(userDataValue)
		//	} else {
		//_ = d.Set("cloud_init", string(userDataValue))
		//}
	}
	if len(userData) > 0 {
		_ = d.Set("user_data", userData)
	}

	return nil
}

func resourceScalewayInstanceServerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	instanceAPI, zone, ID, err := instanceAPIWithZoneAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	wantedState := d.Get("state").(string)
	isStopped := wantedState == InstanceServerStateStopped

	var warnings diag.Diagnostics

	////
	// Construct UpdateServerRequest
	////
	updateRequest := &instance.UpdateServerRequest{
		Zone:     zone,
		ServerID: ID,
	}

	if d.HasChange("name") {
		updateRequest.Name = expandStringPtr(d.Get("name"))
	}

	if d.HasChange("tags") {
		updateRequest.Tags = scw.StringsPtr(expandStrings(d.Get("tags")))
	}

	if d.HasChange("security_group_id") {
		updateRequest.SecurityGroup = &instance.SecurityGroupTemplate{
			ID:   expandZonedID(d.Get("security_group_id")).ID,
			Name: newRandomName("sg"), // this value will be ignored by the API
		}
	}

	if d.HasChange("enable_ipv6") {
		updateRequest.EnableIPv6 = scw.BoolPtr(d.Get("enable_ipv6").(bool))
	}

	if d.HasChange("enable_dynamic_ip") {
		updateRequest.DynamicIPRequired = scw.BoolPtr(d.Get("enable_dynamic_ip").(bool))
	}

	volumes := map[string]*instance.VolumeTemplate{}

	if raw, ok := d.GetOk("additional_volume_ids"); d.HasChange("additional_volume_ids") && ok {
		volumes["0"] = &instance.VolumeTemplate{
			ID:   expandZonedID(d.Get("root_volume.0.volume_id")).ID,
			Name: newRandomName("vol"), // name is ignored by the API, any name will work here
		}

		for i, volumeID := range raw.([]interface{}) {
			volumeHasChange := d.HasChange("additional_volume_ids." + strconv.Itoa(i))
			// local volumes can only be added when the instance is stopped
			if volumeHasChange && !isStopped {
				volumeResp, err := instanceAPI.GetVolume(&instance.GetVolumeRequest{
					Zone:     zone,
					VolumeID: expandZonedID(volumeID).ID,
				})
				if err != nil {
					return diag.FromErr(err)
				}

				// We must be able to tell whether a volume is already present in the server or not
				if volumeResp.Volume.Server != nil {
					if volumeResp.Volume.VolumeType == instance.VolumeVolumeTypeLSSD && volumeResp.Volume.Server.ID != "" {
						return diag.FromErr(fmt.Errorf("instance must be stopped to change local volumes"))
					}
				}
			}
			volumes[strconv.Itoa(i+1)] = &instance.VolumeTemplate{
				ID:   expandZonedID(volumeID).ID,
				Name: newRandomName("vol"), // name is ignored by the API, any name will work here
			}
		}

		updateRequest.Volumes = &volumes
	}

	if d.HasChange("placement_group_id") {
		placementGroupID := expandZonedID(d.Get("placement_group_id")).ID
		if placementGroupID == "" {
			updateRequest.PlacementGroup = &instance.NullableStringValue{Null: true}
		} else {
			if !isStopped {
				return diag.FromErr(fmt.Errorf("instance must be stopped to change placement group"))
			}
			updateRequest.PlacementGroup = &instance.NullableStringValue{Value: placementGroupID}
		}
	}

	////
	// Update reserved IP
	////
	if d.HasChange("ip_id") {
		server, err := instanceAPI.GetServer(&instance.GetServerRequest{
			Zone:     zone,
			ServerID: ID,
		})
		if err != nil {
			return diag.FromErr(err)
		}
		newIPID := expandZonedID(d.Get("ip_id")).ID

		// If an IP is already attached and it's not a dynamic IP we detach it.
		if server.Server.PublicIP != nil && !server.Server.PublicIP.Dynamic {
			_, err = instanceAPI.UpdateIP(&instance.UpdateIPRequest{
				Zone:   zone,
				IP:     server.Server.PublicIP.ID,
				Server: &instance.NullableStringValue{Null: true},
			})
			if err != nil {
				return diag.FromErr(err)
			}
		}

		// If a new IP is provided, we attach it to the server
		if newIPID != "" {
			_, err = instanceAPI.UpdateIP(&instance.UpdateIPRequest{
				Zone:   zone,
				IP:     newIPID,
				Server: &instance.NullableStringValue{Value: ID},
			}, scw.WithContext(ctx))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChanges("boot_type") {
		bootType := instance.BootType(d.Get("boot_type").(string))
		updateRequest.BootType = &bootType
		if !isStopped {
			warnings = append(warnings, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "instance may need to be rebooted to use the new boot type",
			})
		}
	}

	if d.HasChanges("bootscript_id") {
		updateRequest.Bootscript = expandStringPtr(d.Get("bootscript_id").(string))
		if !isStopped {
			warnings = append(warnings, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "instance may need to be rebooted to use the new bootscript",
			})
		}
	}

	////
	// Update server user data
	////
	if d.HasChanges("user_data") {
		userDataRequests := &instance.SetAllServerUserDataRequest{
			Zone:     zone,
			ServerID: ID,
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
					Summary:  "instance may need to be rebooted to use the new cloud init config",
				})
			}
		}

		err := instanceAPI.SetAllServerUserData(userDataRequests)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	////
	// Apply changes
	////

	targetState, err := serverStateExpand(d.Get("state").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	// reach expected state
	err = reachState(ctx, instanceAPI, zone, ID, targetState)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = instanceAPI.UpdateServer(updateRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	return append(warnings, resourceScalewayInstanceServerRead(ctx, d, meta)...)
}

func resourceScalewayInstanceServerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	instanceAPI, zone, ID, err := instanceAPIWithZoneAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// reach stopped state
	err = reachState(ctx, instanceAPI, zone, ID, instance.ServerStateStopped)
	if is404Error(err) {
		return nil
	}
	if err != nil {
		return diag.FromErr(err)
	}

	err = instanceAPI.DeleteServer(&instance.DeleteServerRequest{
		Zone:     zone,
		ServerID: ID,
	}, scw.WithContext(ctx))

	if err != nil && !is404Error(err) {
		return diag.FromErr(err)
	}

	// Related to https://github.com/hashicorp/terraform-plugin-sdk/issues/142
	_, rootVolumeAttributeSet := d.GetOk("root_volume")
	if d.Get("root_volume.0.delete_on_termination").(bool) || !rootVolumeAttributeSet {
		err = instanceAPI.DeleteVolume(&instance.DeleteVolumeRequest{
			Zone:     zone,
			VolumeID: expandZonedID(d.Get("root_volume.0.volume_id")).ID,
		})
		if err != nil && !is404Error(err) {
			return diag.FromErr(err)
		}
	}

	return nil
}
