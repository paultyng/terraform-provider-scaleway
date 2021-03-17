package scaleway

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
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
			Default: schema.DefaultTimeout(defaultLbLbTimeout),
		},
		SchemaVersion: 0,
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
					lb.ProtocolTCP.String(),
					lb.ProtocolHTTP.String(),
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
					lb.ForwardPortAlgorithmRoundrobin.String(),
					lb.ForwardPortAlgorithmLeastconn.String(),
					lb.ForwardPortAlgorithmFirst.String(),
				}, false),
				Default:     lb.ForwardPortAlgorithmRoundrobin.String(),
				Optional:    true,
				Description: "Load balancing algorithm",
			},
			"sticky_sessions": {
				Type: schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					lb.StickySessionsTypeNone.String(),
					lb.StickySessionsTypeCookie.String(),
					lb.StickySessionsTypeTable.String(),
				}, false),
				Default:     lb.StickySessionsTypeNone.String(),
				Optional:    true,
				Description: "Load balancing algorithm",
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
				Default:     false,
				Deprecated:  "Please use proxy_protocol instead",
			},
			"proxy_protocol": {
				Type:        schema.TypeString,
				Description: "Type of PROXY protocol to enable",
				Optional:    true,
				Default:     flattenLbProxyProtocol(lb.ProxyProtocolProxyProtocolNone).(string),
				ValidateFunc: validation.StringInSlice([]string{
					flattenLbProxyProtocol(lb.ProxyProtocolProxyProtocolNone).(string),
					flattenLbProxyProtocol(lb.ProxyProtocolProxyProtocolV1).(string),
					flattenLbProxyProtocol(lb.ProxyProtocolProxyProtocolV2).(string),
					flattenLbProxyProtocol(lb.ProxyProtocolProxyProtocolV2Ssl).(string),
					flattenLbProxyProtocol(lb.ProxyProtocolProxyProtocolV2SslCn).(string),
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
			},
			"health_check_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Port the HC requests will be send to. Default to `forward_port`",
			},
			"health_check_max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     2,
				Description: "Number of allowed failed HC requests before the backend server is marked down",
			},
			"health_check_tcp": {
				Type:          schema.TypeList,
				MaxItems:      1,
				ConflictsWith: []string{"health_check_http", "health_check_https"},
				Optional:      true,
				Computed:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{},
				},
			},
			"health_check_http": {
				Type:          schema.TypeList,
				MaxItems:      1,
				ConflictsWith: []string{"health_check_tcp", "health_check_https"},
				Optional:      true,
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
					},
				},
			},
			"health_check_https": {
				Type:          schema.TypeList,
				MaxItems:      1,
				ConflictsWith: []string{"health_check_tcp", "health_check_http"},
				Optional:      true,
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
					},
				},
			},
			"on_marked_down_action": {
				Type: schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					"none",
					lb.OnMarkedDownActionShutdownSessions.String(),
				}, false),
				Default:     "none",
				Optional:    true,
				Description: "Modify what occurs when a backend server is marked down",
			},
		},
	}
}

func resourceScalewayLbBackendCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	lbAPI := lbAPI(meta)

	region, LbID, err := parseRegionalID(d.Get("lb_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	healthCheckPort := d.Get("health_check_port").(int)
	if healthCheckPort == 0 {
		healthCheckPort = d.Get("forward_port").(int)
	}

	createReq := &lb.CreateBackendRequest{
		Region:                   region,
		LBID:                     LbID,
		Name:                     expandOrGenerateString(d.Get("name"), "lb-bkd"),
		ForwardProtocol:          expandLbProtocol(d.Get("forward_protocol")),
		ForwardPort:              int32(d.Get("forward_port").(int)),
		ForwardPortAlgorithm:     expandLbForwardPortAlgorithm(d.Get("forward_port_algorithm")),
		StickySessions:           expandLbStickySessionsType(d.Get("sticky_sessions")),
		StickySessionsCookieName: d.Get("sticky_sessions_cookie_name").(string),
		HealthCheck: &lb.HealthCheck{
			Port:            int32(healthCheckPort),
			CheckMaxRetries: int32(d.Get("health_check_max_retries").(int)),
			CheckTimeout:    expandDuration(d.Get("health_check_timeout")),
			CheckDelay:      expandDuration(d.Get("health_check_delay")),
			TCPConfig:       expandLbHCTCP(d.Get("health_check_tcp")),
			HTTPConfig:      expandLbHCHTTP(d.Get("health_check_http")),
			HTTPSConfig:     expandLbHCHTTPS(d.Get("health_check_https")),
		},
		ServerIP:           expandStrings(d.Get("server_ips")),
		SendProxyV2:        d.Get("send_proxy_v2").(bool),
		ProxyProtocol:      expandLbProxyProtocol(d.Get("proxy_protocol")),
		TimeoutServer:      expandDuration(d.Get("timeout_server")),
		TimeoutConnect:     expandDuration(d.Get("timeout_connect")),
		TimeoutTunnel:      expandDuration(d.Get("timeout_tunnel")),
		OnMarkedDownAction: expandLbBackendMarkdownAction(d.Get("on_marked_down_action")),
	}

	if definedRegion, ok := d.GetOk("region"); ok {
		region = scw.Region(definedRegion.(string))
		createReq.Region = region
	}

	res, err := lbAPI.CreateBackend(createReq, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newRegionalIDString(region, res.ID))

	return resourceScalewayLbBackendRead(ctx, d, meta)
}

func resourceScalewayLbBackendRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	lbAPI, region, ID, err := lbAPIWithRegionAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := lbAPI.GetBackend(&lb.GetBackendRequest{
		Region:    region,
		BackendID: ID,
	}, scw.WithContext(ctx))

	if err != nil {
		if is404Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	_ = d.Set("lb_id", newRegionalIDString(region, res.LB.ID))
	_ = d.Set("name", res.Name)
	_ = d.Set("forward_protocol", flattenLbProtocol(res.ForwardProtocol))
	_ = d.Set("forward_port", res.ForwardPort)
	_ = d.Set("forward_port_algorithm", flattenLbForwardPortAlgorithm(res.ForwardPortAlgorithm))
	_ = d.Set("sticky_sessions", flattenLbStickySessionsType(res.StickySessions))
	_ = d.Set("sticky_sessions_cookie_name", res.StickySessionsCookieName)
	_ = d.Set("server_ips", res.Pool)
	_ = d.Set("send_proxy_v2", res.SendProxyV2)
	_ = d.Set("proxy_protocol", flattenLbProxyProtocol(res.ProxyProtocol))
	_ = d.Set("timeout_server", flattenDuration(res.TimeoutServer))
	_ = d.Set("timeout_connect", flattenDuration(res.TimeoutConnect))
	_ = d.Set("timeout_tunnel", flattenDuration(res.TimeoutTunnel))
	_ = d.Set("health_check_port", res.HealthCheck.Port)
	_ = d.Set("health_check_max_retries", res.HealthCheck.CheckMaxRetries)
	_ = d.Set("health_check_timeout", flattenDuration(res.HealthCheck.CheckTimeout))
	_ = d.Set("health_check_delay", flattenDuration(res.HealthCheck.CheckDelay))
	_ = d.Set("on_marked_down_action", flattenLbBackendMarkdownAction(res.OnMarkedDownAction))
	_ = d.Set("health_check_tcp", flattenLbHCTCP(res.HealthCheck.TCPConfig))
	_ = d.Set("health_check_http", flattenLbHCHTTP(res.HealthCheck.HTTPConfig))
	_ = d.Set("health_check_https", flattenLbHCHTTPS(res.HealthCheck.HTTPSConfig))

	return nil
}

func resourceScalewayLbBackendUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	lbAPI, region, ID, err := lbAPIWithRegionAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	req := &lb.UpdateBackendRequest{
		Region:                   region,
		BackendID:                ID,
		Name:                     d.Get("name").(string),
		ForwardProtocol:          expandLbProtocol(d.Get("forward_protocol")),
		ForwardPort:              int32(d.Get("forward_port").(int)),
		ForwardPortAlgorithm:     expandLbForwardPortAlgorithm(d.Get("forward_port_algorithm")),
		StickySessions:           expandLbStickySessionsType(d.Get("sticky_sessions")),
		StickySessionsCookieName: d.Get("sticky_sessions_cookie_name").(string),
		SendProxyV2:              d.Get("send_proxy_v2").(bool),
		ProxyProtocol:            expandLbProxyProtocol(d.Get("proxy_protocol")),
		TimeoutServer:            expandDuration(d.Get("timeout_server")),
		TimeoutConnect:           expandDuration(d.Get("timeout_connect")),
		TimeoutTunnel:            expandDuration(d.Get("timeout_tunnel")),
		OnMarkedDownAction:       expandLbBackendMarkdownAction(d.Get("on_marked_down_action")),
	}

	_, err = lbAPI.UpdateBackend(req, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	// Update Health Check
	updateHCRequest := &lb.UpdateHealthCheckRequest{
		Region:          region,
		BackendID:       ID,
		Port:            int32(d.Get("health_check_port").(int)),
		CheckMaxRetries: int32(d.Get("health_check_max_retries").(int)),
		CheckTimeout:    expandDuration(d.Get("health_check_timeout")),
		CheckDelay:      expandDuration(d.Get("health_check_delay")),
		HTTPConfig:      expandLbHCHTTP(d.Get("health_check_http")),
		HTTPSConfig:     expandLbHCHTTPS(d.Get("health_check_https")),
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
	_, err = lbAPI.SetBackendServers(&lb.SetBackendServersRequest{
		Region:    region,
		BackendID: ID,
		ServerIP:  expandStrings(d.Get("server_ips")),
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceScalewayLbBackendRead(ctx, d, meta)
}

func resourceScalewayLbBackendDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	lbAPI, region, ID, err := lbAPIWithRegionAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = lbAPI.DeleteBackend(&lb.DeleteBackendRequest{
		Region:    region,
		BackendID: ID,
	}, scw.WithContext(ctx))

	if err != nil && !is404Error(err) {
		return diag.FromErr(err)
	}

	return nil
}
