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

func resourceScalewayIotRoute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScalewayIotRouteCreate,
		ReadContext:   resourceScalewayIotRouteRead,
		DeleteContext: resourceScalewayIotRouteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 0,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the route",
			},
			"hub_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "The ID of the route's hub",
				DiffSuppressFunc: diffSuppressFuncLocality,
			},
			"topic": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Topic the route subscribes to (wildcards allowed)",
			},
			"database": {
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: "Database Route parameters",
				ExactlyOneOf: []string{
					iot.RouteRouteTypeDatabase.String(),
					iot.RouteRouteTypeRest.String(),
					iot.RouteRouteTypeS3.String(),
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"query": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "SQL query to be executed ($TOPIC and $PAYLOAD variables are available, see documentation)",
						},
						"host": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The database hostname",
						},
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							ForceNew:    true,
							Description: "The database port",
						},
						"dbname": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The database name",
						},
						"username": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The database username",
						},
						"password": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The database password",
							Sensitive:   true,
						},
					},
				},
			},
			"rest": {
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: "Rest Route parameters",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"verb": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The HTTP Verb used to call REST URI",
							ValidateFunc: validation.StringInSlice([]string{
								iot.RouteRestConfigHTTPVerbGet.String(),
								iot.RouteRestConfigHTTPVerbPost.String(),
								iot.RouteRestConfigHTTPVerbPut.String(),
								iot.RouteRestConfigHTTPVerbPatch.String(),
								iot.RouteRestConfigHTTPVerbDelete.String(),
							}, false),
						},
						"uri": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The URI of the REST endpoint",
						},
						"headers": {
							Type:        schema.TypeMap,
							Required:    true,
							ForceNew:    true,
							Description: "The HTTP call extra headers",
							Elem: &schema.Schema{
								Type:     schema.TypeString,
								ForceNew: true,
							},
						},
					},
				},
			},
			"s3": {
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: "S3 Route parameters",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket_region": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The region of the S3 route's destination bucket",
						},
						"bucket_name": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The name of the S3 route's destination bucket",
						},
						"object_prefix": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The string to prefix object names with",
						},
						"strategy": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "How the S3 route's objects will be created: one per topic or one per message",
							ValidateFunc: validation.StringInSlice([]string{
								iot.RouteS3ConfigS3StrategyPerTopic.String(),
								iot.RouteS3ConfigS3StrategyPerMessage.String(),
							}, false),
						},
					},
				},
			},
			// Computed elements
			"region": regionSchema(),
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of the creation of the IoT Route",
			},
		},
	}
}

func resourceScalewayIotRouteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iotAPI, region, err := iotAPIWithRegion(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	////
	// Create route
	////

	req := &iot.CreateRouteRequest{
		Region: region,
		Name:   expandOrGenerateString(d.Get("name"), "route"),
		HubID:  expandZonedID(d.Get("hub_id")).ID,
		Topic:  d.Get("topic").(string),
	}

	if definedRegion, ok := d.GetOk("region"); ok {
		region = scw.Region(definedRegion.(string))
		req.Region = region
	}

	if _, ok := d.GetOk(iot.RouteRouteTypeS3.String()); ok {
		prefixKey := fmt.Sprintf("%s.0", iot.RouteRouteTypeS3.String())
		req.S3Config = &iot.CreateRouteRequestS3Config{
			BucketRegion: d.Get(fmt.Sprintf("%s.bucket_region", prefixKey)).(string),
			BucketName:   d.Get(fmt.Sprintf("%s.bucket_name", prefixKey)).(string),
			ObjectPrefix: d.Get(fmt.Sprintf("%s.object_prefix", prefixKey)).(string),
			Strategy:     iot.RouteS3ConfigS3Strategy(d.Get(fmt.Sprintf("%s.strategy", prefixKey)).(string)),
		}
	} else if _, ok := d.GetOk(iot.RouteRouteTypeRest.String()); ok {
		prefixKey := fmt.Sprintf("%s.0", iot.RouteRouteTypeRest.String())
		req.RestConfig = &iot.CreateRouteRequestRestConfig{
			Verb:    iot.RouteRestConfigHTTPVerb(d.Get(fmt.Sprintf("%s.verb", prefixKey)).(string)),
			URI:     d.Get(fmt.Sprintf("%s.uri", prefixKey)).(string),
			Headers: extractRestHeaders(d, fmt.Sprintf("%s.headers", prefixKey)),
		}
	} else if _, ok := d.GetOk(iot.RouteRouteTypeDatabase.String()); ok {
		prefixKey := fmt.Sprintf("%s.0", iot.RouteRouteTypeDatabase.String())
		req.DbConfig = &iot.CreateRouteRequestDatabaseConfig{
			Host:     d.Get(fmt.Sprintf("%s.host", prefixKey)).(string),
			Port:     uint32(d.Get(fmt.Sprintf("%s.port", prefixKey)).(int)),
			Dbname:   d.Get(fmt.Sprintf("%s.dbname", prefixKey)).(string),
			Username: d.Get(fmt.Sprintf("%s.username", prefixKey)).(string),
			Password: d.Get(fmt.Sprintf("%s.password", prefixKey)).(string),
			Query:    d.Get(fmt.Sprintf("%s.query", prefixKey)).(string),
		}
	} else {
		return diag.FromErr(fmt.Errorf("no route type have been chosen"))
	}

	res, err := iotAPI.CreateRoute(req, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newRegionalIDString(region, res.ID))

	return resourceScalewayIotRouteRead(ctx, d, meta)
}

func resourceScalewayIotRouteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iotAPI, region, routeID, err := iotAPIWithRegionAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	////
	// Read Route
	////
	response, err := iotAPI.GetRoute(&iot.GetRouteRequest{
		Region:  region,
		RouteID: routeID,
	})
	if err != nil {
		if is404Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	_ = d.Set("region", string(region))
	_ = d.Set("name", response.Name)
	_ = d.Set("hub_id", response.HubID)
	_ = d.Set("topic", response.Topic)
	_ = d.Set("created_at", response.CreatedAt.String())

	switch response.Type {
	case iot.RouteRouteTypeDatabase:
		conf := []map[string]interface{}{{
			"query":    response.DbConfig.Query,
			"host":     response.DbConfig.Host,
			"port":     int(response.DbConfig.Port),
			"dbname":   response.DbConfig.Dbname,
			"username": response.DbConfig.Username,
			// Password is never returned. To avoid password getting erased, take it back.
			"password": d.Get(fmt.Sprintf("%s.0.password", iot.RouteRouteTypeDatabase.String())),
		}}
		_ = d.Set("database", conf)
	case iot.RouteRouteTypeRest:
		conf := []map[string]interface{}{{
			"verb":    response.RestConfig.Verb,
			"uri":     response.RestConfig.URI,
			"headers": response.RestConfig.Headers,
		}}
		_ = d.Set("rest", conf)
	case iot.RouteRouteTypeS3:
		conf := []map[string]interface{}{{
			"bucket_region": response.S3Config.BucketRegion,
			"bucket_name":   response.S3Config.BucketName,
			"object_prefix": response.S3Config.ObjectPrefix,
			"strategy":      response.S3Config.Strategy,
		}}
		_ = d.Set("s3", conf)
	}

	return nil
}

func resourceScalewayIotRouteDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iotAPI, region, routeID, err := iotAPIWithRegionAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	////
	// Delete route
	////
	err = iotAPI.DeleteRoute(&iot.DeleteRouteRequest{
		Region:  region,
		RouteID: routeID,
	})
	if err != nil {
		if is404Error(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}
