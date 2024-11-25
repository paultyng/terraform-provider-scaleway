package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	domain "github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/account"
)

func ResourceOrderDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrderDomainCreate,
		ReadContext:   resourceOrderDomainsRead,
		UpdateContext: resourceOrderDomainUpdate,
		DeleteContext: resourceOrderDomainDelete,
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
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The domain name to be managed",
			},
			"duration_in_years": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"project_id": account.ProjectIDSchema(),

			"owner_contact_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the owner contact. Either `owner_contact_id` or `owner_contact` must be provided.",
			},

			"owner_contact": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: contactSchema(),
				},
				Description: "Details of the owner contact. Either `owner_contact_id` or `owner_contact` must be provided.",
			},

			"administrative_contact_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"administrative_contact": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: contactSchema(),
				},
				Description: "Details of the administrative contact.",
			},
			"technical_contact_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"technical_contact": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: contactSchema(),
				},
				Description: "Details of the technical contact.",
			},
			//computed
			"auto_renew_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the automatic renewal of the domain.",
			},
			"dnssec_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the DNSSEC configuration of the domain.",
			},
			"epp_code": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of the domain's EPP codes.",
			},
			"expired_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date of expiration of the domain (RFC 3339 format).",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last modification date of the domain (RFC 3339 format).",
			},
			"registrar": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The registrar managing the domain.",
			},
			"is_external": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether Scaleway is the domain's registrar.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the domain.",
			},
			"organization_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Organization ID associated with the domain.",
			},
			"pending_trade": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if a trade is ongoing for the domain.",
			},
			"external_domain_registration_status": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Registration status of an external domain, if applicable.",
			},
			"transfer_registration_status": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Status of the domain transfer, when available.",
			},
			"linked_products": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of Scaleway resources linked to the domain.",
			},
			"tld": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the TLD.",
						},
						"dnssec_support": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether DNSSEC is supported for this TLD.",
						},
						"duration_in_years_min": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum duration (in years) for which this TLD can be registered.",
						},
						"duration_in_years_max": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum duration (in years) for which this TLD can be registered.",
						},
						"idn_support": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether this TLD supports IDN (Internationalized Domain Names).",
						},
						"offers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of the offer action (e.g., create, transfer).",
									},
									"operation_path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Path of the operation associated with the offer.",
									},
									"price": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "Pricing information for the TLD offer.",
									},
								},
							},
							Description: "Available offers for the TLD.",
						},
						"specifications": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Specifications of the TLD.",
						},
					},
				},
				Description: "Details about the TLD (Top-Level Domain).",
			},

			"dns_zones": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeMap},
				Description: "List of DNS zones associated with the domain.",
			},
		},
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			hasOwnerContactID := d.HasChange("owner_contact_id") && d.Get("owner_contact_id").(string) != ""
			hasOwnerContact := d.HasChange("owner_contact") && len(d.Get("owner_contact").([]interface{})) > 0

			if !hasOwnerContactID && !hasOwnerContact {
				return fmt.Errorf("either `owner_contact_id` or `owner_contact` must be provided")
			}

			if hasOwnerContactID && hasOwnerContact {
				return fmt.Errorf("only one of `owner_contact_id` or `owner_contact` can be provided")
			}

			return nil
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
			Required: true,
		},
		"company_identification_code": {
			Type:     schema.TypeString,
			Required: true,
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
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1, // Ensure it's a single-item list if needed
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"mode": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Mode of the French extension.",
					},
					"individual_info": {
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"whois_opt_in": {
									Type:        schema.TypeBool,
									Optional:    true,
									Description: "Whether WHOIS opt-in is enabled.",
								},
							},
						},
						Description: "Information about the individual owning the domain.",
					},
					"duns_info": {
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"duns_id": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "DUNS ID associated with the domain owner.",
								},
								"local_id": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Local ID of the domain owner.",
								},
							},
						},
						Description: "DUNS information for the domain owner.",
					},
					"association_info": {
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"publication_jo": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Publication date in the Official Journal (RFC3339 format).",
								},
								"publication_jo_page": {
									Type:        schema.TypeInt,
									Optional:    true,
									Description: "Page number of the publication in the Official Journal.",
								},
							},
						},
						Description: "Association-specific information for the domain.",
					},
					"trademark_info": {
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"trademark_inpi": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Trademark information from INPI.",
								},
							},
						},
						Description: "Trademark-related information for the domain.",
					},
					"code_auth_afnic_info": {
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"code_auth_afnic": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Authorization code from AFNIC.",
								},
							},
						},
						Description: "Information regarding AFNIC authorization.",
					},
				},
			},
			Description: "Details specific to French domain extensions.",
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

func resourceOrderDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	registrarAPI := NewRegistrarDomainAPI(m)

	projectID := d.Get("project_id").(string)
	domainName := d.Get("domain_name").(string)
	durationInYears := uint32(d.Get("duration_in_years").(int))

	buyDomainsRequest := &domain.RegistrarAPIBuyDomainsRequest{
		Domains:         []string{domainName},
		DurationInYears: durationInYears,
		ProjectID:       projectID,
	}

	ownerContactID := d.Get("owner_contact_id").(string)
	if ownerContactID != "" {
		buyDomainsRequest.OwnerContactID = &ownerContactID
	} else if ownerContacts, ok := d.GetOk("owner_contact"); ok {
		contacts := ownerContacts.([]interface{})
		if len(contacts) > 0 {
			buyDomainsRequest.OwnerContact = ExpandNewContact(contacts[0].(map[string]interface{}))
		}
	}

	adminContactID := d.Get("administrative_contact_id").(string)
	if adminContactID != "" {
		buyDomainsRequest.AdministrativeContactID = &adminContactID
	} else if adminContacts, ok := d.GetOk("administrative_contact"); ok {
		contacts := adminContacts.([]interface{})
		if len(contacts) > 0 {
			buyDomainsRequest.AdministrativeContact = ExpandNewContact(contacts[0].(map[string]interface{}))
		}
	}

	techContactID := d.Get("technical_contact_id").(string)
	if techContactID != "" {
		buyDomainsRequest.TechnicalContactID = &techContactID
	} else if techContacts, ok := d.GetOk("technical_contact"); ok {
		contacts := techContacts.([]interface{})
		if len(contacts) > 0 {
			buyDomainsRequest.TechnicalContact = ExpandNewContact(contacts[0].(map[string]interface{}))
		}
	}

	resp, err := registrarAPI.BuyDomains(buyDomainsRequest, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	err = waitForTaskCompletion(ctx, registrarAPI, resp.TaskID, 3600)
	_, err = waitForOrderDomain(ctx, registrarAPI, domainName, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.ProjectID + "/" + domainName)

	return resourceOrderDomainsRead(ctx, d, m)
}

func resourceOrderDomainsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	registrarAPI := NewRegistrarDomainAPI(m)
	id := d.Id()

	domainName, err := extractDomainFromID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := waitForOrderDomain(ctx, registrarAPI, domainName, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("domain_name", res.Domain); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("organization_id", res.OrganizationID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("project_id", res.ProjectID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("auto_renew_status", string(res.AutoRenewStatus)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("expired_at", res.ExpiredAt.Format(time.RFC3339)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("updated_at", res.UpdatedAt.Format(time.RFC3339)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("registrar", res.Registrar); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_external", res.IsExternal); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("status", string(res.Status)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("pending_trade", res.PendingTrade); err != nil {
		return diag.FromErr(err)
	}

	if res.OwnerContact != nil {
		if err := d.Set("owner_contact", flattenContact(res.OwnerContact)); err != nil {
			return diag.FromErr(err)
		}
	}
	if res.TechnicalContact != nil {
		if err := d.Set("technical_contact", flattenContact(res.TechnicalContact)); err != nil {
			return diag.FromErr(err)
		}
	}
	if res.AdministrativeContact != nil {
		if err := d.Set("administrative_contact", flattenContact(res.AdministrativeContact)); err != nil {
			return diag.FromErr(err)
		}
	}

	if res.Dnssec != nil {
		if err := d.Set("dnssec_status", string(res.Dnssec.Status)); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("epp_code", res.EppCode); err != nil {
		return diag.FromErr(err)
	}
	if res.Tld != nil {
		if err := d.Set("tld", flattenTLD(res.Tld)); err != nil {
			return diag.FromErr(err)
		}
	}
	if res.TransferRegistrationStatus != nil {
		if err := d.Set("transfer_registration_status", flattenDomainRegistrationStatusTransfer(res.TransferRegistrationStatus)); err != nil {
			return diag.FromErr(err)
		}
	}
	if res.ExternalDomainRegistrationStatus != nil {
		if err := d.Set("external_domain_registration_status", flattenExternalDomainRegistrationStatus(res.ExternalDomainRegistrationStatus)); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("linked_products", res.LinkedProducts); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("dns_zones", flattenDNSZones(res.DNSZones)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceOrderDomainUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	registrarAPI := NewRegistrarDomainAPI(m)

	id := d.Id()
	domainName, err := extractDomainFromID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	updateRequest := &domain.RegistrarAPIUpdateDomainRequest{
		Domain: domainName,
	}

	if d.HasChange("administrative_contact_id") {
		administrativeContactID := d.Get("administrative_contact_id").(string)
		updateRequest.AdministrativeContactID = &administrativeContactID
	}

	if d.HasChange("technical_contact_id") {
		technicalContactID := d.Get("technical_contact_id").(string)
		updateRequest.TechnicalContactID = &technicalContactID
	}

	if d.HasChange("administrative_contact") {
		if adminContacts, ok := d.GetOk("administrative_contact"); ok {
			contacts := adminContacts.([]interface{})
			if len(contacts) > 0 {
				updateRequest.AdministrativeContact = ExpandNewContact(contacts[0].(map[string]interface{}))
			}
		}
	}

	if d.HasChange("technical_contact") {
		if techContacts, ok := d.GetOk("technical_contact"); ok {
			contacts := techContacts.([]interface{})
			if len(contacts) > 0 {
				updateRequest.TechnicalContact = ExpandNewContact(contacts[0].(map[string]interface{}))
			}
		}
	}

	_, err = registrarAPI.UpdateDomain(updateRequest, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceOrderDomainsRead(ctx, d, m)
}

func resourceOrderDomainDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	//registrarAPI := NewRegistrarDomainAPI(m)
	//
	//id := d.Id()
	//domainName, err := extractDomainFromID(id)
	//if err != nil {
	//	return diag.FromErr(err)
	//}
	//
	//deleteRequest := &domain.RegistrarAPIDeleteDomainHostRequest{
	//
	//}
	//}
	//
	//err = registrarAPI.DeleteDomain(deleteRequest, scw.WithContext(ctx))
	//if err != nil {
	//	return diag.FromErr(err)
	//}

	d.SetId("")

	return nil
}
