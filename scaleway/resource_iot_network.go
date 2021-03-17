package scaleway

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	iot "github.com/scaleway/scaleway-sdk-go/api/iot/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func resourceScalewayIotNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScalewayIotNetworkCreate,
		ReadContext:   resourceScalewayIotNetworkRead,
		DeleteContext: resourceScalewayIotNetworkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 0,
		Schema: map[string]*schema.Schema{
			"hub_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "The ID of the hub on which this network will be created",
				DiffSuppressFunc: diffSuppressFuncLocality,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the network",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The type of the network",
				ValidateFunc: validation.StringInSlice([]string{
					iot.NetworkNetworkTypeSigfox.String(),
					iot.NetworkNetworkTypeRest.String(),
				}, false),
			},
			"topic_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The prefix that will be prepended to all topics for this Network",
			},
			// Computed elements
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of the creation of the network",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The endpoint to use when interacting with the network",
			},
			"secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The endpoint key to keep secret",
				Sensitive:   true,
			},
		},
	}
}

func resourceScalewayIotNetworkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iotAPI, region, err := iotAPIWithRegion(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	////
	// Create network
	////

	req := &iot.CreateNetworkRequest{
		Region: region,
		Name:   expandOrGenerateString(d.Get("name"), "network"),
		Type:   iot.NetworkNetworkType(d.Get("type").(string)),
		HubID:  expandID(d.Get("hub_id")),
	}

	if definedRegion, ok := d.GetOk("region"); ok {
		region = scw.Region(definedRegion.(string))
		req.Region = region
	}

	if topicPrefix, ok := d.GetOk("topic_prefix"); ok {
		req.TopicPrefix = topicPrefix.(string)
	}

	res, err := iotAPI.CreateNetwork(req, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newRegionalIDString(region, res.Network.ID))

	// Secret key cannot be retreived later
	_ = d.Set("secret", res.Secret)

	return resourceScalewayIotNetworkRead(ctx, d, meta)
}

func resourceScalewayIotNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iotAPI, region, networkID, err := iotAPIWithRegionAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	////
	// Read Network
	////
	network, err := iotAPI.GetNetwork(&iot.GetNetworkRequest{
		Region:    region,
		NetworkID: networkID,
	}, scw.WithContext(ctx))
	if err != nil {
		if is404Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	_ = d.Set("name", network.Name)
	_ = d.Set("type", network.Type.String())
	_ = d.Set("endpoint", network.Endpoint)
	_ = d.Set("hub_id", newRegionalID(region, network.HubID).String())
	_ = d.Set("created_at", network.CreatedAt.String())
	_ = d.Set("topic_prefix", network.TopicPrefix)

	return nil
}

func resourceScalewayIotNetworkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iotAPI, region, networkID, err := iotAPIWithRegionAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	////
	// Delete Network
	////
	err = iotAPI.DeleteNetwork(&iot.DeleteNetworkRequest{
		Region:    region,
		NetworkID: networkID,
	}, scw.WithContext(ctx))
	if err != nil {
		if !is404Error(err) {
			return diag.FromErr(err)
		}
	}

	return nil
}
