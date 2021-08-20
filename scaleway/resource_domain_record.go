package scaleway

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	domain "github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func resourceScalewayDomainRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScalewayDomainRecordCreate,
		ReadContext:   resourceScalewayDomainRecordRead,
		UpdateContext: resourceScalewayDomainRecordUpdate,
		DeleteContext: resourceScalewayDomainRecordDelete,
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(defaultDomainRecordTimeout),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 0,
		Schema: map[string]*schema.Schema{
			"dns_zone": {
				Type:        schema.TypeString,
				Description: "The zone you want to add the record in",
				Required:    true,
				ForceNew:    true,
			},
			"keep_empty_zone": {
				Type:        schema.TypeBool,
				Description: "When destroy a ressource record, if a zone have only NS, delete the zone",
				Optional:    true,
				Default:     false,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the record",
				ForceNew:    true,
				Required:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "The type of the record",
				ValidateFunc: validation.StringInSlice([]string{
					domain.RecordTypeA.String(),
					domain.RecordTypeAAAA.String(),
					domain.RecordTypeALIAS.String(),
					domain.RecordTypeCNAME.String(),
					domain.RecordTypeMX.String(),
					domain.RecordTypeNS.String(),
					domain.RecordTypePTR.String(),
					domain.RecordTypeSRV.String(),
					domain.RecordTypeTXT.String(),
					domain.RecordTypeTLSA.String(),
					domain.RecordTypeCAA.String(),
				}, false),
				ForceNew: true,
				Required: true,
			},
			"data": {
				Type:        schema.TypeString,
				Description: "The data of the record",
				Required:    true,
			},
			"ttl": {
				Type:         schema.TypeInt,
				Description:  "The ttl of the record",
				Optional:     true,
				Default:      3600,
				ValidateFunc: validation.IntAtLeast(60),
			},
			"priority": {
				Type:         schema.TypeInt,
				Description:  "The priority of the record",
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"geo_ip": {
				Type:          schema.TypeList,
				Description:   "Return record based on client localisation",
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"view", "http_service", "weighted"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"matches": {
							Type:        schema.TypeList,
							Description: "The list of matches",
							MinItems:    1,
							Required:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"countries": {
										Type:        schema.TypeList,
										Optional:    true,
										MinItems:    1,
										Description: "List of countries (eg: FR for France, US for the United States, GB for Great Britain...). List of all countries code: https://api.scaleway.com/domain-private/v2beta1/countries",
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringLenBetween(2, 2),
										},
									},
									"continents": {
										Type:        schema.TypeList,
										Optional:    true,
										MinItems:    1,
										Description: "List of continents (eg: EU for Europe, NA for North America, AS for Asia...). List of all continents code: https://api.scaleway.com/domain-private/v2beta1/continents",
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringLenBetween(2, 2),
										},
									},
									"data": {
										Type:        schema.TypeString,
										Description: "The data of the match result",
										Required:    true,
									},
								},
							},
						},
					},
				},
			},
			"http_service": {
				Type:          schema.TypeList,
				Description:   "Return record based on client localisation",
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"geo_ip", "view", "weighted"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ips": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.IsIPAddress,
							},
							Required:    true,
							MinItems:    1,
							Description: "IPs to check",
						},
						"must_contain": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Text to search",
						},
						"url": {
							Type:         schema.TypeString,
							ValidateFunc: validation.IsURLWithHTTPorHTTPS,
							Required:     true,
							Description:  "URL to match the must_contain text to validate an IP",
						},
						"user_agent": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "User-agent used when checking the URL",
						},
						"strategy": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Strategy to return an IP from the IPs list",
							ValidateFunc: validation.StringInSlice([]string{
								domain.RecordHTTPServiceConfigStrategyRandom.String(),
								domain.RecordHTTPServiceConfigStrategyHashed.String(),
							}, false),
						},
					},
				},
			},
			"view": {
				Type:          schema.TypeList,
				Description:   "Return record based on client subnet",
				Optional:      true,
				ConflictsWith: []string{"geo_ip", "http_service", "weighted"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet": {
							Type:         schema.TypeString,
							Description:  "The subnet of the view",
							Required:     true,
							ValidateFunc: validation.IsCIDR,
						},
						"data": {
							Type:        schema.TypeString,
							Description: "The data of the view record",
							Required:    true,
						},
					},
				},
			},
			"weighted": {
				Type:          schema.TypeList,
				Description:   "Return record based on weight",
				Optional:      true,
				ConflictsWith: []string{"geo_ip", "http_service", "view"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:         schema.TypeString,
							Description:  "The weighted IP",
							Required:     true,
							ValidateFunc: validation.IsIPAddress,
						},
						"weight": {
							Type:         schema.TypeInt,
							Description:  "The weight of the IP",
							Required:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},
					},
				},
			},
		},
	}
}

func resourceScalewayDomainRecordCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	domainAPI := newDomainAPI(meta)

	dnsZone := d.Get("dns_zone").(string)
	geoIP, okGeoIP := d.GetOk("geo_ip")
	record := &domain.Record{
		Data:              d.Get("data").(string),
		Name:              d.Get("name").(string),
		TTL:               uint32(d.Get("ttl").(int)),
		Type:              domain.RecordType(d.Get("type").(string)),
		Priority:          uint32(d.Get("priority").(int)),
		GeoIPConfig:       expandDomainGeoIPConfig(d.Get("data").(string), geoIP, okGeoIP),
		HTTPServiceConfig: expandDomainHTTPService(d.GetOk("http_service")),
		WeightedConfig:    expandDomainWeighted(d.GetOk("weighted")),
		ViewConfig:        expandDomainView(d.GetOk("view")),
		Comment:           nil,
	}
	_, err := domainAPI.UpdateDNSZoneRecords(&domain.UpdateDNSZoneRecordsRequest{
		DNSZone: dnsZone,
		Changes: []*domain.RecordChange{
			{
				Add: &domain.RecordChangeAdd{
					Records: []*domain.Record{record},
				},
			},
		},
		ReturnAllRecords: scw.BoolPtr(false),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceScalewayDomainRecordRead(ctx, d, meta)
}

func resourceScalewayDomainRecordRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	domainAPI := newDomainAPI(meta)
	var record *domain.Record

	// check if this is an inline import. Like: "terraform import scaleway_domain_record.www subdomain.domain.tld/11111111-1111-1111-1111-111111111111"
	if strings.Contains(d.Id(), "/") {
		tab := strings.SplitN(d.Id(), "/", -1)
		if len(tab) != 2 {
			return diag.FromErr(fmt.Errorf("cant parse record id: %s", d.Id()))
		}

		dnsZone := tab[0]
		id := tab[1]

		res, err := domainAPI.ListDNSZoneRecords(&domain.ListDNSZoneRecordsRequest{
			DNSZone: dnsZone,
		}, scw.WithAllPages())

		if err != nil {
			if is404Error(err) {
				d.SetId("")
				return nil
			}
			return diag.FromErr(err)
		}

		for _, r := range res.Records {
			if r.ID == id {
				record = r
				break
			}
		}
	} else {
		res, err := domainAPI.ListDNSZoneRecords(&domain.ListDNSZoneRecordsRequest{
			DNSZone: d.Get("dns_zone").(string),
			Name:    d.Get("name").(string),
			Type:    domain.RecordType(d.Get("type").(string)),
		}, scw.WithAllPages())

		if err != nil {
			if is404Error(err) {
				d.SetId("")
				return nil
			}
			return diag.FromErr(err)
		}

		for _, r := range res.Records {
			if flattenDomainData(r.Data, r.Type).(string) == d.Get("data").(string) {
				record = r
				break
			}
		}
	}

	if record == nil {
		d.SetId("")
		return nil
	}

	d.SetId(record.ID)
	_ = d.Set("data", flattenDomainData(record.Data, record.Type))
	_ = d.Set("ttl", int(record.TTL))
	_ = d.Set("priority", int(record.Priority))
	_ = d.Set("geo_ip", flattenDomainGeoIP(record.GeoIPConfig))
	_ = d.Set("http_service", flattenDomainHTTPService(record.HTTPServiceConfig))
	_ = d.Set("weighted", flattenDomainWeighted(record.WeightedConfig))
	_ = d.Set("view", flattenDomainView(record.ViewConfig))

	return nil
}

func resourceScalewayDomainRecordUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	domainAPI := newDomainAPI(meta)

	if d.HasChangesExcept("dns_zone", "keep_empty_zone", "name", "type") {
		id := d.Id()

		geoIP, okGeoIP := d.GetOk("geo_ip")
		_, err := domainAPI.UpdateDNSZoneRecords(&domain.UpdateDNSZoneRecordsRequest{
			DNSZone: d.Get("dns_zone").(string),
			Changes: []*domain.RecordChange{
				{
					Set: &domain.RecordChangeSet{
						ID: &id,
						Records: []*domain.Record{
							{
								Name:              d.Get("name").(string),
								Data:              d.Get("data").(string),
								Priority:          uint32(d.Get("priority").(int)),
								TTL:               uint32(d.Get("ttl").(int)),
								Type:              domain.RecordType(d.Get("type").(string)),
								GeoIPConfig:       expandDomainGeoIPConfig(d.Get("data").(string), geoIP, okGeoIP),
								HTTPServiceConfig: expandDomainHTTPService(d.GetOk("http_service")),
								WeightedConfig:    expandDomainWeighted(d.GetOk("weighted")),
								ViewConfig:        expandDomainView(d.GetOk("view")),
							},
						},
					},
				},
			},
			ReturnAllRecords: scw.BoolPtr(false),
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceScalewayDomainRecordRead(ctx, d, meta)
}

func resourceScalewayDomainRecordDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	domainAPI := newDomainAPI(meta)

	id := d.Id()
	_, err := domainAPI.UpdateDNSZoneRecords(&domain.UpdateDNSZoneRecordsRequest{
		DNSZone: d.Get("dns_zone").(string),
		Changes: []*domain.RecordChange{
			{
				Delete: &domain.RecordChangeDelete{
					ID: &id,
				},
			},
		},
		ReturnAllRecords: scw.BoolPtr(false),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	// if the zone have only NS records, then delete the all zone
	if !d.Get("keep_empty_zone").(bool) {
		res, err := domainAPI.ListDNSZoneRecords(&domain.ListDNSZoneRecordsRequest{
			DNSZone: d.Get("dns_zone").(string),
		}, scw.WithAllPages())

		if err != nil {
			if is404Error(err) {
				return nil
			}
			return diag.FromErr(err)
		}

		hasRecords := false
		for _, r := range res.Records {
			if r.Type != domain.RecordTypeNS {
				hasRecords = true
				break
			}
		}

		if !hasRecords {
			_, err = domainAPI.DeleteDNSZone(&domain.DeleteDNSZoneRequest{
				DNSZone: d.Get("dns_zone").(string),
			})

			if err != nil {
				if is404Error(err) {
					return nil
				}
				return diag.FromErr(err)
			}
		}
	}

	return nil
}
