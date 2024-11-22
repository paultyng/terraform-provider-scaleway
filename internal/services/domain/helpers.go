package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	domain "github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
	"github.com/scaleway/scaleway-sdk-go/api/std"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/meta"
)

// NewDomainAPI returns a new domain API.
func NewDomainAPI(m interface{}) *domain.API {
	return domain.NewAPI(meta.ExtractScwClient(m))
}

func NewRegistrarDomainAPI(m interface{}) *domain.RegistrarAPI {
	return domain.NewRegistrarAPI(meta.ExtractScwClient(m))
}

func getRecordFromTypeAndData(dnsType domain.RecordType, data string, records []*domain.Record) (*domain.Record, error) {
	var currentRecord *domain.Record
	for _, r := range records {
		flattedData := flattenDomainData(strings.ToLower(r.Data), r.Type).(string)
		flattenCurrentData := flattenDomainData(strings.ToLower(data), r.Type).(string)
		if strings.HasPrefix(flattedData, flattenCurrentData) && r.Type == dnsType {
			if currentRecord != nil {
				return nil, errors.New("multiple records found with same type and data")
			}
			currentRecord = r
			break
		}
	}

	if currentRecord == nil {
		return nil, fmt.Errorf("record with type %s and data %s not found", dnsType.String(), data)
	}

	return currentRecord, nil
}

func FindDefaultReverse(address string) string {
	parts := strings.Split(address, ".")
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}
	return strings.Join(parts, "-") + ".instances.scw.cloud"
}

func ExpandContact(contactMap map[string]interface{}) *domain.Contact {
	if contactMap == nil {
		return nil
	}

	contact := &domain.Contact{
		PhoneNumber:  contactMap["phone_number"].(string),
		LegalForm:    domain.ContactLegalForm(contactMap["legal_form"].(string)),
		Firstname:    contactMap["firstname"].(string),
		Lastname:     contactMap["lastname"].(string),
		Email:        contactMap["email"].(string),
		AddressLine1: contactMap["address_line_1"].(string),
		Zip:          contactMap["zip"].(string),
		City:         contactMap["city"].(string),
		Country:      contactMap["country"].(string),
	}

	// Optional fields
	if v, ok := contactMap["company_name"].(string); ok && v != "" {
		contact.CompanyName = v
	}
	if v, ok := contactMap["email_alt"].(string); ok && v != "" {
		contact.EmailAlt = v
	}
	if v, ok := contactMap["fax_number"].(string); ok && v != "" {
		contact.FaxNumber = v
	}
	if v, ok := contactMap["address_line_2"].(string); ok && v != "" {
		contact.AddressLine2 = v
	}
	if v, ok := contactMap["vat_identification_code"].(string); ok && v != "" {
		contact.VatIDentificationCode = v
	}
	if v, ok := contactMap["company_identification_code"].(string); ok && v != "" {
		contact.CompanyIDentificationCode = v
	}
	if v, ok := contactMap["lang"].(string); ok && v != "" {
		contact.Lang = std.LanguageCode(v)
	}
	if v, ok := contactMap["resale"].(bool); ok {
		contact.Resale = v
	}
	if v, ok := contactMap["state"].(string); ok && v != "" {
		contact.State = v
	}
	if v, ok := contactMap["whois_opt_in"].(bool); ok {
		contact.WhoisOptIn = v
	}

	if extFr, ok := contactMap["extension_fr"].(map[string]interface{}); ok && len(extFr) > 0 {
		contact.ExtensionFr = expandContactExtension(extFr, "fr").(*domain.ContactExtensionFR)
	}
	if extEu, ok := contactMap["extension_eu"].(map[string]interface{}); ok && len(extEu) > 0 {
		contact.ExtensionEu = expandContactExtension(extEu, "eu").(*domain.ContactExtensionEU)
	}
	if extNl, ok := contactMap["extension_nl"].(map[string]interface{}); ok && len(extNl) > 0 {
		contact.ExtensionNl = expandContactExtension(extNl, "nl").(*domain.ContactExtensionNL)
	}

	return contact
}

func expandContactExtension(extensionMap map[string]interface{}, extensionType string) interface{} {
	if extensionMap == nil || len(extensionMap) == 0 {
		return nil
	}

	switch extensionType {
	case "fr":
		return &domain.ContactExtensionFR{
			Mode:              parseEnum[domain.ContactExtensionFRMode](extensionMap, "mode", domain.ContactExtensionFRModeModeUnknown),
			IndividualInfo:    parseStruct[domain.ContactExtensionFRIndividualInfo](extensionMap, "individual_info"),
			DunsInfo:          parseStruct[domain.ContactExtensionFRDunsInfo](extensionMap, "duns_info"),
			AssociationInfo:   parseStruct[domain.ContactExtensionFRAssociationInfo](extensionMap, "association_info"),
			TrademarkInfo:     parseStruct[domain.ContactExtensionFRTrademarkInfo](extensionMap, "trademark_info"),
			CodeAuthAfnicInfo: parseStruct[domain.ContactExtensionFRCodeAuthAfnicInfo](extensionMap, "code_auth_afnic_info"),
		}
	case "nl":
		legalFormRegistrationNumber := ""
		if value, ok := extensionMap["legal_form_registration_number"]; ok {
			if str, isString := value.(string); isString {
				legalFormRegistrationNumber = str
			}
		}

		return &domain.ContactExtensionNL{
			LegalForm:                   parseEnum[domain.ContactExtensionNLLegalForm](extensionMap, "legal_form", domain.ContactExtensionNLLegalFormLegalFormUnknown),
			LegalFormRegistrationNumber: legalFormRegistrationNumber,
		}
	case "eu":
		europeanCitizenship := ""
		if value, ok := extensionMap["european_citizenship"]; ok {
			if str, isString := value.(string); isString {
				europeanCitizenship = str
			}
		}
		return &domain.ContactExtensionEU{
			EuropeanCitizenship: europeanCitizenship,
		}
	default:
		return nil
	}
}

func ExpandNewContact(contactMap map[string]interface{}) *domain.NewContact {
	if contactMap == nil {
		return nil
	}

	contact := &domain.NewContact{
		PhoneNumber:  contactMap["phone_number"].(string),
		LegalForm:    domain.ContactLegalForm(contactMap["legal_form"].(string)),
		Firstname:    contactMap["firstname"].(string),
		Lastname:     contactMap["lastname"].(string),
		Email:        contactMap["email"].(string),
		AddressLine1: contactMap["address_line_1"].(string),
		Zip:          contactMap["zip"].(string),
		City:         contactMap["city"].(string),
		Country:      contactMap["country"].(string),
	}

	if v, ok := contactMap["resale"].(bool); ok {
		contact.Resale = v
	} else {
		contact.Resale = false
	}

	if v, ok := contactMap["whois_opt_in"].(bool); ok {
		contact.WhoisOptIn = v
	} else {
		contact.WhoisOptIn = false
	}

	if v, ok := contactMap["company_name"].(string); ok {
		contact.CompanyName = scw.StringPtr(v)
	}
	if v, ok := contactMap["email_alt"].(string); ok {
		contact.EmailAlt = scw.StringPtr(v)
	}
	if v, ok := contactMap["fax_number"].(string); ok {
		contact.FaxNumber = scw.StringPtr(v)
	}
	if v, ok := contactMap["address_line_2"].(string); ok {
		contact.AddressLine2 = scw.StringPtr(v)
	}
	if v, ok := contactMap["vat_identification_code"].(string); ok {
		contact.VatIDentificationCode = scw.StringPtr(v)
	}
	if v, ok := contactMap["company_identification_code"].(string); ok {
		contact.CompanyIDentificationCode = scw.StringPtr(v)
	}
	if v, ok := contactMap["state"].(string); ok {
		contact.State = scw.StringPtr(v)
	}

	if extFr, ok := contactMap["extension_fr"].(map[string]interface{}); ok {
		contact.ExtensionFr = expandContactExtension(extFr, "fr").(*domain.ContactExtensionFR)
	}
	if extEu, ok := contactMap["extension_eu"].(map[string]interface{}); ok {
		contact.ExtensionEu = expandContactExtension(extEu, "eu").(*domain.ContactExtensionEU)
	}
	if extNl, ok := contactMap["extension_nl"].(map[string]interface{}); ok {
		contact.ExtensionNl = expandContactExtension(extNl, "nl").(*domain.ContactExtensionNL)
	}

	return contact
}

func parseEnum[T ~string](data map[string]interface{}, key string, defaultValue T) T {
	if value, ok := data[key].(string); ok {
		return T(value)
	}
	return defaultValue
}

func parseStruct[T any](data map[string]interface{}, key string) *T {
	if nested, ok := data[key].(map[string]interface{}); ok {
		var result T
		mapToStruct(nested, &result)
		return &result
	}
	return nil
}

func mapToStruct(data map[string]interface{}, target interface{}) {
	switch t := target.(type) {
	case *domain.ContactExtensionFRIndividualInfo:
		if v, ok := data["whois_opt_in"].(bool); ok {
			t.WhoisOptIn = v
		}
	case *domain.ContactExtensionFRDunsInfo:
		if v, ok := data["duns_id"].(string); ok {
			t.DunsID = v
		}
		if v, ok := data["local_id"].(string); ok {
			t.LocalID = v
		}
	case *domain.ContactExtensionFRAssociationInfo:
		if v, ok := data["publication_jo"].(string); ok {
			if parsedTime, err := time.Parse(time.RFC3339, v); err == nil {
				t.PublicationJo = &parsedTime
			}
		}
		if v, ok := data["publication_jo_page"].(float64); ok {
			t.PublicationJoPage = uint32(v)
		}
	case *domain.ContactExtensionFRTrademarkInfo:
		if v, ok := data["trademark_inpi"].(string); ok {
			t.TrademarkInpi = v
		}
	case *domain.ContactExtensionFRCodeAuthAfnicInfo:
		if v, ok := data["code_auth_afnic"].(string); ok {
			t.CodeAuthAfnic = v
		}
	}
}
