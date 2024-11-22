package domain

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	domain "github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/account"
)

func resourceDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainCreate,
		ReadContext:   resourceDomainsRead,
		Timeouts: &schema.ResourceTimeout{
			Create:  schema.DefaultTimeout(defaultDomainRecordTimeout),
			Read:    schema.DefaultTimeout(defaultDomainRecordTimeout),
			Update:  schema.DefaultTimeout(defaultDomainRecordTimeout),
			Delete:  schema.DefaultTimeout(defaultDomainRecordTimeout),
			Default: schema.DefaultTimeout(defaultDomainRecordTimeout),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 0,

		Schema: map[string]*schema.Schema{
			"domain_names": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of domain names to be managed",
			},
			"duration_in_years": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"project_id": account.ProjectIDSchema(),

			"owner_contact_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"owner_contact": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: contactSchema(),
				},
			},
			"administrative_contact_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"administrative_contact": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: contactSchema(),
				},
			},
			"technical_contact_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"technical_contact": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: contactSchema(),
				},
			},
		},
	}
}

func contactSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"legal_form": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Legal form of the contact (e.g., 'individual' or 'organization').",
		},
		"firstname": {
			Type:     schema.TypeString,
			Required: true,
		},
		"lastname": {
			Type:     schema.TypeString,
			Required: true,
		},
		"company_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"email": {
			Type:     schema.TypeString,
			Required: true,
		},
		"email_alt": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"phone_number": {
			Type:     schema.TypeString,
			Required: true,
		},
		"fax_number": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"address_line_1": {
			Type:     schema.TypeString,
			Required: true,
		},
		"address_line_2": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"zip": {
			Type:     schema.TypeString,
			Required: true,
		},
		"city": {
			Type:     schema.TypeString,
			Required: true,
		},
		"country": {
			Type:     schema.TypeString,
			Required: true,
		},
		"vat_identification_code": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"company_identification_code": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"lang": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"resale": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"extension_fr": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"extension_eu": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"whois_opt_in": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"state": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"extension_nl": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func resourceDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	registrarAPI := NewRegistrarDomainAPI(m)

	projectID := d.Get("project_id").(string)
	domainNamesInterface := d.Get("domain_names").([]interface{})
	domains := make([]string, len(domainNamesInterface))
	for i, v := range domainNamesInterface {
		domains[i] = v.(string)
	}
	durationInYears := uint32(d.Get("duration_in_years").(int))

	buyDomainsRequest := &domain.RegistrarAPIBuyDomainsRequest{
		Domains:         domains,
		DurationInYears: durationInYears,
		ProjectID:       projectID,
	}

	ownerContactID := d.Get("owner_contact_id").(string)
	if ownerContactID != "" {

		buyDomainsRequest.OwnerContactID = &ownerContactID
	} else if ownerContact, ok := d.GetOk("owner_contact"); ok {
		buyDomainsRequest.OwnerContact = ExpandNewContact(ownerContact.(map[string]interface{}))
	}

	adminContactID := d.Get("administrative_contact_id").(string)
	if adminContactID != "" {
		buyDomainsRequest.AdministrativeContactID = &adminContactID
	} else if adminContact, ok := d.GetOk("administrative_contact"); ok {
		buyDomainsRequest.AdministrativeContact = ExpandNewContact(adminContact.(map[string]interface{}))
	}

	techContactID := d.Get("technical_contact_id").(string)
	if adminContactID != "" {
		buyDomainsRequest.TechnicalContactID = &techContactID
	} else if techContact, ok := d.GetOk("technical_contact"); ok {
		buyDomainsRequest.TechnicalContact = ExpandNewContact(techContact.(map[string]interface{}))
	}

	resp, err := registrarAPI.BuyDomains(buyDomainsRequest, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.ProjectID + "/" + resp.TaskID)

	return resourceDomainsRead(ctx, d, m)
}

func resourceDomainsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func contactToMap(contact *domain.Contact) map[string]interface{} {
	if contact == nil {
		return nil
	}
	return map[string]interface{}{
		"id":        contact.ID,
		"firstname": contact.Firstname,
		"lastname":  contact.Lastname,
		"email":     contact.Email,
	}
}

func fetchAllDomains(ctx context.Context, registrarAPI *domain.RegistrarAPI, projectID string) ([]*domain.DomainSummary, error) {
	var allDomains []*domain.DomainSummary

	page := int32(1)
	pageSize := uint32(1000)

	for {
		listDomainsRequest := &domain.RegistrarAPIListDomainsRequest{
			ProjectID: &projectID,
			Page:      &page,
			PageSize:  &pageSize,
		}

		domainsResponse, err := registrarAPI.ListDomains(listDomainsRequest, scw.WithContext(ctx))
		if err != nil {
			return nil, err
		}

		allDomains = append(allDomains, domainsResponse.Domains...)

		if len(domainsResponse.Domains) < int(pageSize) {
			break
		}
		page++
	}

	return allDomains, nil
}
