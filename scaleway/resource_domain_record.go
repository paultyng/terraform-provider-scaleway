package scaleway

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	domain "github.com/scaleway/scaleway-sdk-go/api/domain/v2alpha2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func resourceScalewayDomainRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceScalewayDomainRecordCreate,
		Read:   resourceScalewayDomainRecordRead,
		Update: resourceScalewayDomainRecordUpdate,
		Delete: resourceScalewayDomainRecordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the record",
				ForceNew:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "The type of the record",
				ValidateFunc: validation.StringInSlice([]string{
					domain.RecordTypeA.String(),
					domain.RecordTypeAAAA.String(),
					domain.RecordTypeCNAME.String(),
					domain.RecordTypeMX.String(),
					domain.RecordTypeNS.String(),
					domain.RecordTypePTR.String(),
					domain.RecordTypeSRV.String(),
					domain.RecordTypeTXT.String(),
				}, false),
				ForceNew: true,
			},
			"data": {
				Type:        schema.TypeString,
				Description: "The data of the record",
				ForceNew:    true,
			},
			"ttl": {
				Type:        schema.TypeInt,
				Description: "The ttl of the record",
			},
			"dns_zone": {
				Type:        schema.TypeString,
				Description: "The zone you want to add the record in",
				Required:    true,
			},
		},
	}
}

func resourceScalewayDomainRecordCreate(d *schema.ResourceData, m interface{}) error {
	domainAPI := domainAPI(m)

	dnsZone := d.Get("dns_zone").(string)
	name := d.Get("name").(string)
	data := d.Get("data").(string)
	recordType := d.Get("type").(string)
	res, err := domainAPI.UpdateDNSZoneRecords(&domain.UpdateDNSZoneRecordsRequest{
		DNSZone: dnsZone,
		Changes: []*domain.RecordChange{
			{
				Add: &domain.RecordChangeAdd{
					Records: []*domain.Record{
						{
							Data:            data,
							Name:            name,
							TTL:             uint32(d.Get("ttl").(int)),
							Type:            domain.RecordType(recordType),
							Priority:        0,
							Comment:         nil,
							GeoIPConfig:     nil,
							ServiceUpConfig: nil,
							WeightedConfig:  nil,
							ViewConfig:      nil,
						},
					},
				},
			},
		},
		ReturnAllRecords: Bool(false),
	})
	if err != nil {
		return err
	}

	d.SetId(flattenRecordID(dnsZone, name, recordType, data))
	return resourceScalewayAccountSSHKeyRead(d, m)
}

func resourceScalewayDomainRecordRead(d *schema.ResourceData, m interface{}) error {
	domainAPI := domainAPI(m)

	dnsZone, name, recordType, data, err := expandRecordID(d.Id())
	if err != nil {
		return err
	}

	res, err := domainAPI.ListDNSZoneRecords(&domain.ListDNSZoneRecordsRequest{
		DNSZone: dnsZone,
		Name:    name,
		Type:    domain.RecordType(recordType),
	}, scw.WithAllPages())

	if err != nil {
		if is404Error(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	var record *domain.Record
	for _, r := range res.Records {
		if r.Data == data {
			record = r
			break
		}
	}

	if record == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("ttl", record.TTL)
	return nil
}

func resourceScalewayDomainRecordUpdate(d *schema.ResourceData, m interface{}) error {
	domainAPI := domainAPI(m)

	dnsZone, name, recordType, data, err := expandRecordID(d.Id())
	if err != nil {
		return err
	}

	if d.HasChange("ttl") {
		_, err := domainAPI.UpdateDNSZoneRecords(&domain.UpdateDNSZoneRecordsRequest{
			DNSZone: dnsZone,
			Changes: []*domain.RecordChange{
				{
					Set: &domain.RecordChangeSet{
						Data: data,
						Name: name,
						TTL:  uint32(d.Get("ttl").(int)),
						Type: domain.RecordType(recordType),
						Records: []*domain.Record{
							{
								Data:     d.Get("data").(string),
								Name:     d.Get("name").(string),
								Priority: 0,
								TTL:      uint32(d.Get("ttl").(int)),
								Type:     domain.RecordType(d.Get("type").(string)),
							},
						},
					},
				},
			},
			ReturnAllRecords: nil,
		})
		if err != nil {
			return err
		}
	}

	return resourceScalewayDomainRecordRead(d, m)
}

func resourceScalewayDomainRecordDelete(d *schema.ResourceData, m interface{}) error {
	domainAPI := domainAPI(m)

	dnsZone, name, recordType, data, err := expandRecordID(d.Id())
	if err != nil {
		return err
	}

	_, err = domainAPI.UpdateDNSZoneRecords(&domain.UpdateDNSZoneRecordsRequest{
		DNSZone: dnsZone,
		Changes: []*domain.RecordChange{
			{
				Delete: &domain.RecordChangeDelete{
					Data: data,
					Name: name,
					Type: domain.RecordType(recordType),
				},
			},
		},
		ReturnAllRecords: nil,
	})
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func flattenRecordID(zone string, name string, recordType string, data string) string {
	return fmt.Sprintf("%s/%s/%s/%s", zone, name, recordType, data)
}

func expandRecordID(id string) (zone string, name string, recordType string, data string, err error) {
	_, err = fmt.Sscanf(id, "%s/%s/%s", &zone, &name, &recordType, &data)
	return
}
