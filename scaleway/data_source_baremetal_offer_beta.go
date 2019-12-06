package scaleway

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	baremetal "github.com/scaleway/scaleway-sdk-go/api/baremetal/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func dataSourceScalewayBaremetalOfferBeta() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScalewayBaremetalOfferBetaRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Exact name of the desired offer",
				ConflictsWith: []string{"offer_id"},
			},
			"offer_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "ID of the desired offer",
				ConflictsWith: []string{"name"},
			},
			"include_disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Include disabled offers",
			},
			"zone": zoneSchema(),

			"bandwidth": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Available Bandwidth with the offer",
			},
			"commercial_range": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Commercial range of the offer",
			},
			"cpu": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "CPU specifications of the offer",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CPU name",
						},
						"core_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of cores",
						},
						"frequency": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Frequency of the CPU",
						},
						"thread_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of threads",
						},
					},
				},
			},
			"disk": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Disk specifications of the offer",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of disk",
						},
						"capacity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Capacity of the disk in byte",
						},
					},
				},
			},
			"memory": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Memory specifications of the offer",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of memory",
						},
						"capacity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Capacity of the memory in byte",
						},
						"frequency": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Frequency of the memory",
						},
						"ecc": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "True if error-correcting code is available on this memory",
						},
					},
				},
			},
			"stock": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Stock status for this offer",
			},
		},
	}
}

func dataSourceScalewayBaremetalOfferBetaRead(d *schema.ResourceData, m interface{}) error {
	meta := m.(*Meta)
	baremetalApi, fallBackZone, err := getBaremetalAPIWithZone(d, meta)
	if err != nil {
		return err
	}

	zone, offerID, _ := parseZonedID(datasourceNewZonedID(d.Get("offer_id"), fallBackZone))
	res, err := baremetalApi.ListOffers(&baremetal.ListOffersRequest{
		Zone: zone,
	}, scw.WithAllPages())
	if err != nil {
		return err
	}

	matches := []*baremetal.Offer(nil)
	for _, offer := range res.Offers {
		if offer.Name == d.Get("name") || offer.ID == offerID {
			if !offer.Enable && !d.Get("include_disabled").(bool) {
				return fmt.Errorf("offer %s (%s) found in zone %s but is disabled. Add allow_disabled=true in your terraform config to use it.", offer.Name, offer.Name, zone)
			}
			matches = append(matches, offer)
		}
	}
	if len(matches) == 0 {
		return fmt.Errorf("no offer found with the name %s in zone %s", d.Get("name"), zone)
	}
	if len(matches) > 1 {
		return fmt.Errorf("%d offers found with the same name %s in zone %s", len(matches), d.Get("name"), zone)
	}

	offer := matches[0]
	zonedID := datasourceNewZonedID(offer.ID, zone)
	d.SetId(zonedID)
	d.Set("offer_id", zonedID)
	d.Set("zone", zone)

	if err != nil {
		return err
	}

	d.Set("name", offer.Name)
	d.Set("include_disabled", !offer.Enable)
	d.Set("bandwidth", offer.Bandwidth)
	d.Set("commercial_range", offer.CommercialRange)
	d.Set("cpu", flattenBaremetalCPUs(offer.CPU))
	d.Set("disk", flattenBaremetalDisks(offer.Disk))
	d.Set("memory", flattenBaremetalMemory(offer.Memory))
	d.Set("stock", offer.Stock.String())

	return nil
}
