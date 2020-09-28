package scaleway

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	domain "github.com/scaleway/scaleway-sdk-go/api/domain/v2alpha2"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func resourceScalewayDomainDNSRecords() *schema.Resource {
	return &schema.Resource{
		Create: resourceScalewayDomainDNSRecordsCreate,
		Read:   resourceScalewayDomainDNSRecordsRead,
		Delete: resourceScalewayDomainDNSRecordsDelete,
		Update: resourceScalewayDomainDNSRecordsCreate,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 0,
		Schema: map[string]*schema.Schema{
			"dns_zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"records": {
				Type:     schema.TypeSet,
				Optional: true,
				Set:      recordHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"data": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ttl": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"organization_id": organizationIDSchema(),
		},
	}
}

func resourceScalewayDomainDNSRecordsCreate(d *schema.ResourceData, m interface{}) error {
	meta := m.(*Meta)
	domainAPI := domain.NewAPI(meta.scwClient)

	recordsToAdd := []*domain.Record(nil)
	recordsSet := d.Get("records").(*schema.Set)
	for _, item := range recordsSet.List() {
		recordsToAdd = append(recordsToAdd, expandDomainDNSRecords(item))
	}

	_, err := domainAPI.UpdateDNSZoneRecords(&domain.UpdateDNSZoneRecordsRequest{
		DNSZone: d.Get("dns_zone").(string),
		Changes: []*domain.RecordChange{
			{
				Clear: &domain.RecordChangeClear{},
			},
			{
				Add: &domain.RecordChangeAdd{
					Records: recordsToAdd,
				},
			},
		},
	})
	if err != nil {
		return err
	}

	d.SetId(d.Get("dns_zone").(string))

	return resourceScalewayDomainDNSRecordsRead(d, m)
}

func resourceScalewayDomainDNSRecordsRead(d *schema.ResourceData, m interface{}) error {
	meta := m.(*Meta)
	domainAPI := domain.NewAPI(meta.scwClient)

	dnsZoneRecordsResponse, err := domainAPI.ListDNSZoneRecords(&domain.ListDNSZoneRecordsRequest{
		DNSZone: d.Id(),
	}, scw.WithAllPages())
	if err != nil {
		return err
	}

	stateRecords := make([]interface{}, 0, len(dnsZoneRecordsResponse.Records))

	onlineNSRecordPatter := regexp.MustCompile(`^ns[0-9]\.online\.net\.$`)
	for _, record := range dnsZoneRecordsResponse.Records {
		if onlineNSRecordPatter.MatchString(record.Data) {
			// Ignore Online NS records.
			continue
		}

		stateRecords = append(stateRecords, flattenDomainDNSRecords(record))
	}

	_ = d.Set("dns_zone", d.Id())
	_ = d.Set("records", schema.NewSet(recordHash, stateRecords))

	return nil
}

func resourceScalewayDomainDNSRecordsDelete(d *schema.ResourceData, m interface{}) error {
	meta := m.(*Meta)
	domainAPI := domain.NewAPI(meta.scwClient)

	_, err := domainAPI.UpdateDNSZoneRecords(&domain.UpdateDNSZoneRecordsRequest{
		DNSZone: d.Get("dns_zone").(string),
		Changes: []*domain.RecordChange{
			{
				Clear: &domain.RecordChangeClear{},
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func recordHash(v interface{}) int {
	userData := v.(map[string]interface{})
	return hashcode.String(fmt.Sprintf("%s%s%s%d", userData["name"], userData["type"], userData["data"], userData["ttl"]))
}

func expandDomainDNSRecords(i interface{}) *domain.Record {
	flatRecord := i.(map[string]interface{})

	record := &domain.Record{
		Name:     flatRecord["name"].(string),
		Type:     domain.RecordType(flatRecord["type"].(string)),
		Data:     flatRecord["data"].(string),
		TTL:      uint32(flatRecord["ttl"].(int)),
		Priority: uint32(flatRecord["priority"].(int)),
	}

	if comment := flatRecord["comment"].(string); comment != "" {
		record.Comment = scw.StringPtr(comment)
	}

	return record
}

func flattenDomainDNSRecords(record *domain.Record) map[string]interface{} {
	return map[string]interface{}{
		"data":     record.Data,
		"name":     record.Name,
		"ttl":      record.TTL,
		"type":     record.Type,
		"priority": record.Priority,
		"comment":  record.Comment,
	}
}
