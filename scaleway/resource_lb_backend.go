package scaleway

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	lbSDK "github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func resourceScalewayLbBackend() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScalewayLbBackendCreate,
		ReadContext:   resourceScalewayLbBackendRead,
		UpdateContext: resourceScalewayLbBackendUpdate,
		DeleteContext: resourceScalewayLbBackendDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create:  schema.DefaultTimeout(defaultLbLbTimeout),
			Read:    schema.DefaultTimeout(defaultLbLbTimeout),
			Update:  schema.DefaultTimeout(defaultLbLbTimeout),
			Delete:  schema.DefaultTimeout(defaultLbLbTimeout),
			Default: schema.DefaultTimeout(defaultLbLbTimeout),
		},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{Version: 0, Type: lbUpgradeV1SchemaType(), Upgrade: lbUpgradeV1SchemaUpgradeFunc},
		},
		Schema: map[string]*schema.Schema{
			"lb_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The load-balancer ID",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the backend",
			},
			"forward_protocol": {
				Type: schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					lbSDK.ProtocolTCP.String(),
					lbSDK.ProtocolHTTP.String(),
				}, false),
				Required:    true,
				Description: "Backend protocol",
			},
			"forward_port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "User sessions will be forwarded to this port of backend servers",
			},
			"forward_port_algorithm": {
				Type: schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					lbSDK.ForwardPortAlgorithmRoundrobin.String(),
					lbSDK.ForwardPortAlgorithmLeastconn.String(),
					lbSDK.ForwardPortAlgorithmFirst.String(),
				}, false),
				Default:     lbSDK.ForwardPortAlgorithmRoundrobin.String(),
				Optional:    true,
				Description: "Load balancing algorithm",
			},
			"sticky_sessions": {
				Type: schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					lbSDK.StickySessionsTypeNone.String(),
					lbSDK.StickySessionsTypeCookie.String(),
					lbSDK.StickySessionsTypeTable.String(),
				}, false),
				Default:     lbSDK.StickySessionsTypeNone.String(),
				Optional:    true,
				Description: "The type of sticky sessions",
			},
			"sticky_sessions_cookie_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Cookie name for for sticky sessions",
			},
			"server_ips": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsIPAddress,
				},
				Optional:    true,
				Description: "Backend server IP addresses list (IPv4 or IPv6)",
			},
			"send_proxy_v2": {
				Type:        schema.TypeBool,
				Description: "Enables PROXY protocol version 2",
				Optional:    true,
				Computed:    true,
				Deprecated:  "Please use proxy_protocol instead",
			},
			"proxy_protocol": {
				Type:        schema.TypeString,
				Description: "Type of PROXY protocol to enable",
				Optional:    true,
				Default:     flattenLbProxyProtocol(lbSDK.ProxyProtocolProxyProtocolNone).(string),
				ValidateFunc: validation.StringInSlice([]string{
					flattenLbProxyProtocol(lbSDK.ProxyProtocolProxyProtocolNone).(string),
					flattenLbProxyProtocol(lbSDK.ProxyProtocolProxyProtocolV1).(string),
					flattenLbProxyProtocol(lbSDK.ProxyProtocolProxyProtocolV2).(string),
					flattenLbProxyProtocol(lbSDK.ProxyProtocolProxyProtocolV2Ssl).(string),
					flattenLbProxyProtocol(lbSDK.ProxyProtocolProxyProtocolV2SslCn).(string),
				}, false),
			},
			// Timeouts
			"timeout_server": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: diffSuppressFuncDuration,
				ValidateFunc:     validateDuration(),
				Description:      "Maximum server connection inactivity time",
			},
			"timeout_connect": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: diffSuppressFuncDuration,
				ValidateFunc:     validateDuration(),
				Description:      "Maximum initial server connection establishment time",
			},
			"timeout_tunnel": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: diffSuppressFuncDuration,
				ValidateFunc:     validateDuration(),
				Description:      "Maximum tunnel inactivity time",
			},

			// Health Check
			"health_check_timeout": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: diffSuppressFuncDuration,
				ValidateFunc:     validateDuration(),
				Default:          "30s",
				Description:      "Timeout before we consider a HC request failed",
			},
			"health_check_delay": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: diffSuppressFuncDuration,
				ValidateFunc:     validateDuration(),
				Default:          "60s",
				Description:      "Interval between two HC requests",
				Deprecated:       "Please use the `health_check` block instead",
			},
			"health_check_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Port the HC requests will be send to. Default to `forward_port`",
				Deprecated:  "Please use the `health_check` block instead",
			},
			"health_check_max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     2,
				Description: "Number of allowed failed HC requests before the backend server is marked down",
				Deprecated:  "Please use the `health_check` block instead",
			},
			"health_check_tcp": {
				Type:          schema.TypeList,
				MaxItems:      1,
				ConflictsWith: []string{"health_check_http", "health_check_https", "health_check"},
				Optional:      true,
				Computed:      true,
				Deprecated:    "Health check TCP is block is deprecated. Please use the health_check block",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{},
				},
			},
			"health_check_http": {
				Type:          schema.TypeList,
				MaxItems:      1,
				ConflictsWith: []string{"health_check_tcp", "health_check_https", "health_check"},
				Optional:      true,
				Deprecated:    "Health check HTTP is block is deprecated. Please use the health_check block",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uri": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The HTTP endpoint URL to call for HC requests",
						},
						"method": {
							Type:        schema.TypeString,
							Default:     "GET",
							Optional:    true,
							Description: "The HTTP method to use for HC requests",
						},
						"code": {
							Type:        schema.TypeInt,
							Default:     200,
							Optional:    true,
							Description: "The expected HTTP status code",
						},
						"host_header": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The HTTP host header to use for HC requests",
						},
					},
				},
			},
			"health_check_https": {
				Type:          schema.TypeList,
				MaxItems:      1,
				ConflictsWith: []string{"health_check_tcp", "health_check_http", "health_check"},
				Optional:      true,
				Deprecated:    "Health check HTTPS is block is deprecated. Please use the health_check block",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uri": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The HTTPS endpoint URL to call for HC requests",
						},
						"method": {
							Type:        schema.TypeString,
							Default:     "GET",
							Optional:    true,
							Description: "The HTTP method to use for HC requests",
						},
						"code": {
							Type:        schema.TypeInt,
							Default:     200,
							Optional:    true,
							Description: "The expected HTTP status code",
						},
						"host_header": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The HTTP host header to use for HC requests",
						},
						"sni": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The SNI to use for HC requests over SSL",
						},
					},
				},
			},
			"health_check": {
				Type:          schema.TypeList,
				MaxItems:      1,
				ConflictsWith: []string{"health_check_tcp", "health_check_http", "health_check_https"},
				Optional:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: diffSuppressFuncIgnoreCase,
							Description:      "Protocol the load balancer uses when performing health checks. Must be either LDAP, MySQL, PgSQL, HTTP, or HTTPS. Defaults to TCP.",
							ValidateFunc: validation.StringInSlice([]string{
								lbSDK.ProtocolTCP.String(),
								lbSDK.ProtocolHTTP.String(),
								"mysql",
								"pgsql",
								"https",
								"ldap",
							}, false),
						},
						"database_user": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The user to use for the health check on protocol mysql or pgsql",
						},
						"check_timeout": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: diffSuppressFuncDuration,
							ValidateFunc:     validateDuration(),
							Default:          "30s",
							Description:      "Timeout before we consider a HC request failed",
						},
						"check_delay": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: diffSuppressFuncDuration,
							ValidateFunc:     validateDuration(),
							Default:          "60s",
							Description:      "The time to wait between two consecutive health checks",
						},
						"transient_check_delay": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     "10",
							Description: "The time to wait between two consecutive health checks when a backend server is in a transient state",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The Port the HC requests will be send to. Default to `forward_port`",
						},
						"check_send_proxy": {
							Type:        schema.TypeBool,
							Description: "It defines whether proxy protocol should be activated for the health check",
							Optional:    true,
							Default:     false,
						},
						"max_retries": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     2,
							Description: "Number of allowed failed HC requests before the backend server is marked down",
						},
						"uri": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The HTTPS endpoint URL to call for HC requests",
						},
						"method": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The HTTP method to use for HC requests",
						},
						"code": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The expected HTTP status code",
						},
						"host_header": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The HTTP host header to use for HC requests",
						},
						"sni": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The SNI to use for HC requests over SSL",
						},
					},
				},
			},
			"on_marked_down_action": {
				Type: schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					"none",
					lbSDK.OnMarkedDownActionShutdownSessions.String(),
				}, false),
				Default:     "none",
				Optional:    true,
				Description: "Modify what occurs when a backend server is marked down",
			},
			"failover_host": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `Scaleway S3 bucket website to be served in case all backend servers are down

**NOTE** : Only the host part of the Scaleway S3 bucket website is expected.
E.g. 'failover-website.s3-website.fr-par.scw.cloud' if your bucket website URL is 'https://failover-website.s3-website.fr-par.scw.cloud/'.`,
			},
			"ssl_bridging": {
				Type:        schema.TypeBool,
				Description: "Enables SSL between load balancer and backend servers",
				Optional:    true,
				Default:     false,
			},
			"ignore_ssl_server_verify": {
				Type:        schema.TypeBool,
				Description: "Specifies whether the Load Balancer should check the backend server’s certificate before initiating a connection",
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func resourceScalewayLbBackendCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	lbAPI, _, err := lbAPIWithZone(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	// parse lb_id. It will be forced to a zoned lb
	zone, lbID, err := parseZonedID(d.Get("lb_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	healthCheckPort := d.Get("health_check_port").(int)
	if healthCheckPort == 0 {
		healthCheckPort = d.Get("forward_port").(int)
	}

	_, err = waitForLB(ctx, lbAPI, zone, lbID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		if is403Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	healthCheckoutTimeout, err := expandDuration(d.Get("health_check_timeout"))
	if err != nil {
		return diag.FromErr(err)
	}
	healthCheckDelay, err := expandDuration(d.Get("health_check_delay"))
	if err != nil {
		return diag.FromErr(err)
	}
	timeoutServer, err := expandDuration(d.Get("timeout_server"))
	if err != nil {
		return diag.FromErr(err)
	}
	timeoutConnect, err := expandDuration(d.Get("timeout_connect"))
	if err != nil {
		return diag.FromErr(err)
	}
	timeoutTunnel, err := expandDuration(d.Get("timeout_tunnel"))
	if err != nil {
		return diag.FromErr(err)
	}

	createReq := &lbSDK.ZonedAPICreateBackendRequest{
		Zone:                     zone,
		LBID:                     lbID,
		Name:                     expandOrGenerateString(d.Get("name"), "lb-bkd"),
		ForwardProtocol:          expandLbProtocol(d.Get("forward_protocol")),
		ForwardPort:              int32(d.Get("forward_port").(int)),
		ForwardPortAlgorithm:     expandLbForwardPortAlgorithm(d.Get("forward_port_algorithm")),
		StickySessions:           expandLbStickySessionsType(d.Get("sticky_sessions")),
		StickySessionsCookieName: d.Get("sticky_sessions_cookie_name").(string),
		ServerIP:                 expandStrings(d.Get("server_ips")),
		ProxyProtocol:            expandLbProxyProtocol(d.Get("proxy_protocol")),
		TimeoutServer:            timeoutServer,
		TimeoutConnect:           timeoutConnect,
		TimeoutTunnel:            timeoutTunnel,
		OnMarkedDownAction:       expandLbBackendMarkdownAction(d.Get("on_marked_down_action")),
		FailoverHost:             expandStringPtr(d.Get("failover_host")),
		SslBridging:              expandBoolPtr(getBool(d, "ssl_bridging")),
		IgnoreSslServerVerify:    expandBoolPtr(getBool(d, "ignore_ssl_server_verify")),
	}

	if hc, exist := d.GetOk("health_check"); exist {
		createReq.HealthCheck, err = expandLbHealthCheck(hc)
		if err != nil {
			return diag.FromErr(err)
		}

		if createReq.HealthCheck.Port == 0 {
			createReq.HealthCheck.Port = int32(healthCheckPort)
		}
	} else {
		createReq.HealthCheck = &lbSDK.HealthCheck{
			Port:            int32(healthCheckPort),
			CheckMaxRetries: int32(d.Get("health_check_max_retries").(int)),
			CheckTimeout:    healthCheckoutTimeout,
			CheckDelay:      healthCheckDelay,
			TCPConfig:       expandLbHCTCP(d.Get("health_check_tcp")),
			HTTPConfig:      expandLbHCHTTP(d.Get("health_check_http")),
			HTTPSConfig:     expandLbHCHTTPS(d.Get("health_check_https")),
		}
	}

	// deprecated attribute
	createReq.SendProxyV2 = expandBoolPtr(getBool(d, "send_proxy_v2"))

	res, err := lbAPI.CreateBackend(createReq, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = waitForLB(ctx, lbAPI, zone, res.LB.ID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		if is403Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.SetId(newZonedIDString(zone, res.ID))

	return resourceScalewayLbBackendRead(ctx, d, meta)
}

func resourceScalewayLbBackendRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	lbAPI, zone, ID, err := lbAPIWithZoneAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	backend, err := lbAPI.GetBackend(&lbSDK.ZonedAPIGetBackendRequest{
		Zone:      zone,
		BackendID: ID,
	}, scw.WithContext(ctx))
	if err != nil {
		if is403Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	_ = d.Set("lb_id", newZonedIDString(zone, backend.LB.ID))
	_ = d.Set("name", backend.Name)
	_ = d.Set("forward_protocol", flattenLbProtocol(backend.ForwardProtocol))
	_ = d.Set("forward_port", backend.ForwardPort)
	_ = d.Set("forward_port_algorithm", flattenLbForwardPortAlgorithm(backend.ForwardPortAlgorithm))
	_ = d.Set("sticky_sessions", flattenLbStickySessionsType(backend.StickySessions))
	_ = d.Set("sticky_sessions_cookie_name", backend.StickySessionsCookieName)
	_ = d.Set("server_ips", backend.Pool)
	_ = d.Set("proxy_protocol", flattenLbProxyProtocol(backend.ProxyProtocol))
	_ = d.Set("timeout_server", flattenDuration(backend.TimeoutServer))
	_ = d.Set("timeout_connect", flattenDuration(backend.TimeoutConnect))
	_ = d.Set("timeout_tunnel", flattenDuration(backend.TimeoutTunnel))
	_ = d.Set("health_check_port", backend.HealthCheck.Port)
	_ = d.Set("health_check_max_retries", backend.HealthCheck.CheckMaxRetries)
	_ = d.Set("health_check_timeout", flattenDuration(backend.HealthCheck.CheckTimeout))
	_ = d.Set("health_check_delay", flattenDuration(backend.HealthCheck.CheckDelay))
	_ = d.Set("on_marked_down_action", flattenLbBackendMarkdownAction(backend.OnMarkedDownAction))

	if _, ok := d.GetOk("health_check"); ok {
		_ = d.Set("health_check", flattenLbHealthCheck(backend.HealthCheck))
		_ = d.Set("health_check_tcp", nil)
		_ = d.Set("health_check_http", nil)
		_ = d.Set("health_check_https", nil)
	} else {
		_ = d.Set("health_check_tcp", flattenLbHCTCP(backend.HealthCheck.TCPConfig))
		_ = d.Set("health_check_http", flattenLbHCHTTP(backend.HealthCheck.HTTPConfig))
		_ = d.Set("health_check_https", flattenLbHCHTTPS(backend.HealthCheck.HTTPSConfig))
	}

	_ = d.Set("send_proxy_v2", flattenBoolPtr(backend.SendProxyV2))
	_ = d.Set("failover_host", backend.FailoverHost)
	_ = d.Set("ssl_bridging", flattenBoolPtr(backend.SslBridging))
	_ = d.Set("ignore_ssl_server_verify", flattenBoolPtr(backend.IgnoreSslServerVerify))

	_, err = waitForLB(ctx, lbAPI, zone, backend.LB.ID, d.Timeout(schema.TimeoutRead))
	if err != nil {
		if is403Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}

//gocyclo:ignore
func resourceScalewayLbBackendUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	lbAPI, zone, ID, err := lbAPIWithZoneAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, lbID, err := parseZonedID(d.Get("lb_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = waitForLB(ctx, lbAPI, zone, lbID, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		if is403Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	timeoutServer, err := expandDuration(d.Get("timeout_server"))
	if err != nil {
		return diag.FromErr(err)
	}
	timeoutConnect, err := expandDuration(d.Get("timeout_connect"))
	if err != nil {
		return diag.FromErr(err)
	}
	timeoutTunnel, err := expandDuration(d.Get("timeout_tunnel"))
	if err != nil {
		return diag.FromErr(err)
	}

	req := &lbSDK.ZonedAPIUpdateBackendRequest{
		Zone:                     zone,
		BackendID:                ID,
		Name:                     d.Get("name").(string),
		ForwardProtocol:          expandLbProtocol(d.Get("forward_protocol")),
		ForwardPort:              int32(d.Get("forward_port").(int)),
		ForwardPortAlgorithm:     expandLbForwardPortAlgorithm(d.Get("forward_port_algorithm")),
		StickySessions:           expandLbStickySessionsType(d.Get("sticky_sessions")),
		StickySessionsCookieName: d.Get("sticky_sessions_cookie_name").(string),
		ProxyProtocol:            expandLbProxyProtocol(d.Get("proxy_protocol")),
		TimeoutServer:            timeoutServer,
		TimeoutConnect:           timeoutConnect,
		TimeoutTunnel:            timeoutTunnel,
		OnMarkedDownAction:       expandLbBackendMarkdownAction(d.Get("on_marked_down_action")),
		FailoverHost:             expandStringPtr(d.Get("failover_host")),
		SslBridging:              expandBoolPtr(getBool(d, "ssl_bridging")),
		IgnoreSslServerVerify:    expandBoolPtr(getBool(d, "ignore_ssl_server_verify")),
	}

	// deprecated
	req.SendProxyV2 = expandBoolPtr(getBool(d, "send_proxy_v2"))

	_, err = lbAPI.UpdateBackend(req, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	healthCheckoutTimeout, err := expandDuration(d.Get("health_check_timeout"))
	if err != nil {
		return diag.FromErr(err)
	}
	healthCheckDelay, err := expandDuration(d.Get("health_check_delay"))
	if err != nil {
		return diag.FromErr(err)
	}
	// Update Health Check
	updateHCRequest := &lbSDK.ZonedAPIUpdateHealthCheckRequest{
		Zone:            zone,
		BackendID:       ID,
		Port:            int32(d.Get("health_check_port").(int)),
		CheckMaxRetries: int32(d.Get("health_check_max_retries").(int)),
		CheckTimeout:    healthCheckoutTimeout,
		CheckDelay:      healthCheckDelay,
		HTTPConfig:      expandLbHCHTTP(d.Get("health_check_http")),
		HTTPSConfig:     expandLbHCHTTPS(d.Get("health_check_https")),
	}

	if d.HasChanges("health_check") {
		hc, err := expandLbHealthCheck(d.Get("health_check"))
		if err != nil {
			return diag.FromErr(err)
		}

		if hc.Port == 0 {
			hc.Port = int32(d.Get("forward_port").(int))
		}

		updateHealthCheckChanges(hc, updateHCRequest)
	}

	// As this is the default behaviour if no other HC type are present we enable TCP
	if updateHCRequest.HTTPConfig == nil && updateHCRequest.HTTPSConfig == nil {
		updateHCRequest.TCPConfig = expandLbHCTCP(d.Get("health_check_tcp"))
	}

	_, err = lbAPI.UpdateHealthCheck(updateHCRequest, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	// Update Backend servers
	_, err = lbAPI.SetBackendServers(&lbSDK.ZonedAPISetBackendServersRequest{
		Zone:      zone,
		BackendID: ID,
		ServerIP:  expandStrings(d.Get("server_ips")),
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = waitForLB(ctx, lbAPI, zone, lbID, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		if is403Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return resourceScalewayLbBackendRead(ctx, d, meta)
}

func resourceScalewayLbBackendDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	lbAPI, zone, ID, err := lbAPIWithZoneAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, lbID, err := parseZonedID(d.Get("lb_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = waitForLB(ctx, lbAPI, zone, lbID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		if is403Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = lbAPI.DeleteBackend(&lbSDK.ZonedAPIDeleteBackendRequest{
		Zone:      zone,
		BackendID: ID,
	}, scw.WithContext(ctx))

	if err != nil && !is404Error(err) {
		return diag.FromErr(err)
	}

	_, err = waitForLB(ctx, lbAPI, zone, lbID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		if is403Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}
