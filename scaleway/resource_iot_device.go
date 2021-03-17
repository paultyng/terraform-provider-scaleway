package scaleway

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	iot "github.com/scaleway/scaleway-sdk-go/api/iot/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func resourceScalewayIotDevice() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScalewayIotDeviceCreate,
		ReadContext:   resourceScalewayIotDeviceRead,
		UpdateContext: resourceScalewayIotDeviceUpdate,
		DeleteContext: resourceScalewayIotDeviceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 0,
		Schema: map[string]*schema.Schema{
			"hub_id": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "The ID of the hub on which this device will be created",
				DiffSuppressFunc: diffSuppressFuncLocality,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the device",
				ForceNew:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the device",
			},
			"allow_insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Allow plain and server-authenticated SSL connections in addition to mutually-authenticated ones",
			},
			"allow_multiple_connections": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Allow multiple connections",
			},
			"message_filters": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Rules to authorize or deny the device to publish/subscribe to specific topics",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"publish": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Rule to restrict topics the device can publish to",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"policy": {
										Type:         schema.TypeString,
										Optional:     true,
										Description:  "Publish message filter policy",
										Default:      iot.DeviceMessageFiltersRulePolicyReject.String(),
										RequiredWith: []string{"message_filters.0.publish.0.topics"},
										ValidateFunc: validation.StringInSlice([]string{
											iot.DeviceMessageFiltersRulePolicyAccept.String(),
											iot.DeviceMessageFiltersRulePolicyReject.String(),
										}, false),
									},
									"topics": {
										Type:         schema.TypeList,
										Optional:     true,
										Description:  "List of topics in the set",
										RequiredWith: []string{"message_filters.0.publish.0.policy"},
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"subscribe": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Rule to restrict topics the device can subscribe to",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"policy": {
										Type:         schema.TypeString,
										Optional:     true,
										Description:  "Subscribe message filter policy",
										Default:      iot.DeviceMessageFiltersRulePolicyReject.String(),
										RequiredWith: []string{"message_filters.0.subscribe.0.topics"},
										ValidateFunc: validation.StringInSlice([]string{
											iot.DeviceMessageFiltersRulePolicyAccept.String(),
											iot.DeviceMessageFiltersRulePolicyReject.String(),
										}, false),
									},
									"topics": {
										Type:         schema.TypeList,
										Optional:     true,
										Description:  "List of topics in the set",
										RequiredWith: []string{"message_filters.0.subscribe.0.policy"},
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			// Computed elements
			"region": regionSchema(),
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of the creation of the device",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of the last update of the device",
			},
			"certificate": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Certificate section of the device",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crt": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "X509 PEM encoded certificate of the device",
						},
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "X509 PEM encoded key of the device",
							Sensitive:   true,
						},
					},
				},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the device",
			},
			"last_activity_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of last MQTT activity of the device",
			},
			"is_connected": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The MQTT connection status of the device",
			},
		},
	}
}

func resourceScalewayIotDeviceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iotAPI, region, err := iotAPIWithRegion(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	////
	// Create device
	////

	req := &iot.CreateDeviceRequest{
		Region: region,
		HubID:  expandID(d.Get("hub_id")),
		Name:   expandOrGenerateString(d.Get("name"), "device"),
	}

	if definedRegion, ok := d.GetOk("region"); ok {
		region = scw.Region(definedRegion.(string))
		req.Region = region
	}

	if allowInsecure, ok := d.GetOk("allow_insecure"); ok {
		req.AllowInsecure = allowInsecure.(bool)
	}

	if allowMultipleConn, ok := d.GetOk("allow_multiple_connections"); ok {
		req.AllowMultipleConnections = allowMultipleConn.(bool)
	}

	if description, ok := d.GetOk("description"); ok {
		req.Description = scw.StringPtr(description.(string))
	}

	if _, ok := d.GetOk("message_filters"); ok {
		mf := &iot.DeviceMessageFilters{}

		fqfn := "message_filters.0"
		if _, ok := d.GetOk(fmt.Sprintf("%s.publish", fqfn)); ok {
			fqfnS := fmt.Sprintf("%s.publish.0", fqfn)
			mfSet := iot.DeviceMessageFiltersRule{}

			if policy, ok := d.GetOk(fmt.Sprintf("%s.policy", fqfnS)); ok {
				mfSet.Policy = iot.DeviceMessageFiltersRulePolicy(policy.(string))
			}
			if topics, ok := d.GetOk(fmt.Sprintf("%s.topics", fqfnS)); ok {
				mfSet.Topics = scw.StringsPtr(expandStringsOrEmpty(topics))
			}

			mf.Publish = &mfSet
		}

		if _, ok := d.GetOk(fmt.Sprintf("%s.subscribe", fqfn)); ok {
			fqfnP := fmt.Sprintf("%s.subscribe.0", fqfn)
			mfSet := iot.DeviceMessageFiltersRule{}

			if policy, ok := d.GetOk(fmt.Sprintf("%s.policy", fqfnP)); ok {
				mfSet.Policy = iot.DeviceMessageFiltersRulePolicy(policy.(string))
			}
			if topics, ok := d.GetOk(fmt.Sprintf("%s.topics", fqfnP)); ok {
				mfSet.Topics = scw.StringsPtr(expandStringsOrEmpty(topics))
			}

			mf.Subscribe = &mfSet
		}

		req.MessageFilters = mf
	}

	res, err := iotAPI.CreateDevice(req, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newRegionalIDString(region, res.Device.ID))

	// Certificate and Key cannot be retreived later
	cert := map[string]interface{}{
		"crt": res.Certificate.Crt,
		"key": res.Certificate.Key,
	}
	_ = d.Set("certificate", []map[string]interface{}{cert})

	return resourceScalewayIotDeviceRead(ctx, d, meta)
}

func resourceScalewayIotDeviceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iotAPI, region, deviceID, err := iotAPIWithRegionAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	////
	// Read Device
	////
	device, err := iotAPI.GetDevice(&iot.GetDeviceRequest{
		Region:   region,
		DeviceID: deviceID,
	}, scw.WithContext(ctx))
	if err != nil {
		if is404Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	_ = d.Set("name", device.Name)
	_ = d.Set("status", device.Status)
	_ = d.Set("hub_id", newRegionalID(region, device.HubID).String())
	_ = d.Set("created_at", device.CreatedAt.String())
	_ = d.Set("updated_at", device.UpdatedAt.String())
	_ = d.Set("last_activity_at", device.LastActivityAt.String())
	_ = d.Set("allow_insecure", device.AllowInsecure)
	_ = d.Set("allow_multiple_connections", device.AllowMultipleConnections)
	_ = d.Set("is_connected", device.IsConnected)
	_ = d.Set("description", device.Description)

	mf := map[string]interface{}{}
	mfHasNonDefaultChange := false

	// We need to set the message filters only in case when we already set a value or we got non default value

	// In case of already set value
	if _, ok := d.GetOk("message_filters.0"); ok {
		mfHasNonDefaultChange = true
	}
	// In case of non default change
	if device.MessageFilters.Publish.Policy != iot.DeviceMessageFiltersRulePolicyReject ||
		(device.MessageFilters.Publish.Topics != nil && len(*device.MessageFilters.Publish.Topics) != 0) {
		mfHasNonDefaultChange = true
	}
	if device.MessageFilters.Subscribe.Policy != iot.DeviceMessageFiltersRulePolicyReject ||
		(device.MessageFilters.Subscribe.Topics != nil && len(*device.MessageFilters.Subscribe.Topics) != 0) {
		mfHasNonDefaultChange = true
	}

	if mfHasNonDefaultChange {
		mf["publish"] = []map[string]interface{}{{
			"policy": device.MessageFilters.Publish.Policy,
			"topics": *device.MessageFilters.Publish.Topics,
		}}

		mf["subscribe"] = []map[string]interface{}{{
			"policy": device.MessageFilters.Subscribe.Policy,
			"topics": *device.MessageFilters.Subscribe.Topics,
		}}
	}

	if mfHasNonDefaultChange {
		_ = d.Set("message_filters", []map[string]interface{}{mf})
	}

	return nil
}

func resourceScalewayIotDeviceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iotAPI, region, hubID, err := iotAPIWithRegionAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	////
	// Update Device
	////
	updateRequest := &iot.UpdateDeviceRequest{
		Region:   region,
		DeviceID: hubID,
	}

	if d.HasChange("allow_insecure") {
		updateRequest.AllowInsecure = scw.BoolPtr(d.Get("allow_insecure").(bool))
	}

	if d.HasChange("message_filters") {
		fqfn := "message_filters.0"
		mf := &iot.DeviceMessageFilters{}
		updateRequest.MessageFilters = mf

		if d.HasChange(fmt.Sprintf("%s.publish", fqfn)) {
			fqfnS := fmt.Sprintf("%s.publish.0", fqfn)
			mfSet := iot.DeviceMessageFiltersRule{}
			mf.Publish = &mfSet

			mfSet.Policy = iot.DeviceMessageFiltersRulePolicy(
				d.Get(fmt.Sprintf("%s.policy", fqfnS)).(string))

			mfSet.Topics = scw.StringsPtr(
				expandStringsOrEmpty(d.Get(fmt.Sprintf("%s.topics", fqfnS))))
		}

		if d.HasChange(fmt.Sprintf("%s.subscribe", fqfn)) {
			fqfnP := fmt.Sprintf("%s.subscribe.0", fqfn)
			mfSet := iot.DeviceMessageFiltersRule{}
			mf.Subscribe = &mfSet

			mfSet.Policy = iot.DeviceMessageFiltersRulePolicy(
				d.Get(fmt.Sprintf("%s.policy", fqfnP)).(string))

			mfSet.Topics = scw.StringsPtr(
				expandStringsOrEmpty(d.Get(fmt.Sprintf("%s.topics", fqfnP))))
		}
	}

	if d.HasChange("hub_id") {
		updateRequest.HubID = scw.StringPtr(d.Get("hub_id").(string))
	}

	if d.HasChange("allow_multiple_connections") {
		updateRequest.AllowMultipleConnections = scw.BoolPtr(d.Get("allow_multiple_connections").(bool))
	}

	if d.HasChange("description") {
		updateRequest.Description = scw.StringPtr(d.Get("description").(string))
	}

	_, err = iotAPI.UpdateDevice(updateRequest, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceScalewayIotDeviceRead(ctx, d, meta)
}

func resourceScalewayIotDeviceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iotAPI, region, deviceID, err := iotAPIWithRegionAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	////
	// Delete Device
	////
	err = iotAPI.DeleteDevice(&iot.DeleteDeviceRequest{
		Region:   region,
		DeviceID: deviceID,
	}, scw.WithContext(ctx))
	if err != nil {
		if !is404Error(err) {
			return diag.FromErr(err)
		}
	}

	return nil
}
