package domain

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	domain "github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
	"github.com/scaleway/scaleway-sdk-go/api/std"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/meta"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/transport"
)

const (
	defaultWaitDomainOrderDomainRetryInterval = 10 * time.Minute
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
		extension := expandContactExtension(extFr, "fr")
		if extension != nil {
			contact.ExtensionFr = extension.(*domain.ContactExtensionFR)
		}
	}
	if extEu, ok := contactMap["extension_eu"].(map[string]interface{}); ok && len(extEu) > 0 {
		extension := expandContactExtension(extEu, "eu")
		if extension != nil {
			contact.ExtensionEu = extension.(*domain.ContactExtensionEU)
		}
	}

	if extNl, ok := contactMap["extension_nl"].(map[string]interface{}); ok && len(extNl) > 0 {
		extension := expandContactExtension(extNl, "nl")
		if extension != nil {
			contact.ExtensionNl = extension.(*domain.ContactExtensionNL)
		}
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

	if extFr, ok := contactMap["extension_fr"].(map[string]interface{}); ok && len(extFr) > 0 {
		extension := expandContactExtension(extFr, "fr")
		if extension != nil {
			contact.ExtensionFr = extension.(*domain.ContactExtensionFR)
		}
	}
	if extEu, ok := contactMap["extension_eu"].(map[string]interface{}); ok && len(extEu) > 0 {
		extension := expandContactExtension(extEu, "eu")
		if extension != nil {
			contact.ExtensionEu = extension.(*domain.ContactExtensionEU)
		}
	}

	if extNl, ok := contactMap["extension_nl"].(map[string]interface{}); ok && len(extNl) > 0 {
		extension := expandContactExtension(extNl, "nl")
		if extension != nil {
			contact.ExtensionNl = extension.(*domain.ContactExtensionNL)
		}
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

func ExtractDomainFromID(id string) (string, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid ID format, expected 'projectID/domainName', got: %s", id)
	}
	return parts[1], nil
}

func flattenContact(contact *domain.Contact) []map[string]interface{} {
	if contact == nil {
		return nil
	}

	flattened := map[string]interface{}{
		"phone_number":                contact.PhoneNumber,
		"legal_form":                  string(contact.LegalForm),
		"firstname":                   contact.Firstname,
		"lastname":                    contact.Lastname,
		"email":                       contact.Email,
		"address_line_1":              contact.AddressLine1,
		"zip":                         contact.Zip,
		"city":                        contact.City,
		"country":                     contact.Country,
		"company_name":                contact.CompanyName,
		"email_alt":                   contact.EmailAlt,
		"fax_number":                  contact.FaxNumber,
		"address_line_2":              contact.AddressLine2,
		"vat_identification_code":     contact.VatIDentificationCode,
		"company_identification_code": contact.CompanyIDentificationCode,
		"lang":                        string(contact.Lang),
		"resale":                      contact.Resale,
		"state":                       contact.State,
		"whois_opt_in":                contact.WhoisOptIn,
	}

	if contact.ExtensionFr != nil {
		flattened["extension_fr"] = flattenContactExtensionFR(contact.ExtensionFr)
	}
	if contact.ExtensionEu != nil {
		flattened["extension_eu"] = flattenContactExtensionEU(contact.ExtensionEu)
	}
	if contact.ExtensionNl != nil {
		flattened["extension_nl"] = flattenContactExtensionNL(contact.ExtensionNl)
	}

	return []map[string]interface{}{flattened}
}

func flattenContactExtensionFR(ext *domain.ContactExtensionFR) []map[string]interface{} {
	if ext == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"mode":                 string(ext.Mode),
			"individual_info":      flattenContactExtensionFRIndividualInfo(ext.IndividualInfo),
			"duns_info":            flattenContactExtensionFRDunsInfo(ext.DunsInfo),
			"association_info":     flattenContactExtensionFRAssociationInfo(ext.AssociationInfo),
			"trademark_info":       flattenContactExtensionFRTrademarkInfo(ext.TrademarkInfo),
			"code_auth_afnic_info": flattenContactExtensionFRCodeAuthAfnicInfo(ext.CodeAuthAfnicInfo),
		},
	}
}

func flattenContactExtensionFRIndividualInfo(info *domain.ContactExtensionFRIndividualInfo) []map[string]interface{} {
	if info == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"whois_opt_in": info.WhoisOptIn,
		},
	}
}

func flattenContactExtensionFRDunsInfo(info *domain.ContactExtensionFRDunsInfo) []map[string]interface{} {
	if info == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"duns_id":  info.DunsID,
			"local_id": info.LocalID,
		},
	}
}

func flattenContactExtensionFRAssociationInfo(info *domain.ContactExtensionFRAssociationInfo) []map[string]interface{} {
	if info == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"publication_jo":      info.PublicationJo.Format(time.RFC3339),
			"publication_jo_page": info.PublicationJoPage,
		},
	}
}

func flattenContactExtensionFRTrademarkInfo(info *domain.ContactExtensionFRTrademarkInfo) []map[string]interface{} {
	if info == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"trademark_inpi": info.TrademarkInpi,
		},
	}
}

func flattenContactExtensionFRCodeAuthAfnicInfo(info *domain.ContactExtensionFRCodeAuthAfnicInfo) []map[string]interface{} {
	if info == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"code_auth_afnic": info.CodeAuthAfnic,
		},
	}
}

func flattenContactExtensionEU(ext *domain.ContactExtensionEU) []map[string]interface{} {
	if ext == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"european_citizenship": ext.EuropeanCitizenship,
		},
	}
}

func flattenContactExtensionNL(ext *domain.ContactExtensionNL) []map[string]interface{} {
	if ext == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"legal_form":                     string(ext.LegalForm),
			"legal_form_registration_number": ext.LegalFormRegistrationNumber,
		},
	}
}

func flattenTLD(tld *domain.Tld) []map[string]interface{} {
	if tld == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"name":                  tld.Name,
			"dnssec_support":        tld.DnssecSupport,
			"duration_in_years_min": tld.DurationInYearsMin,
			"duration_in_years_max": tld.DurationInYearsMax,
			"idn_support":           tld.IDnSupport,
			"offers":                flattenTldOffers(tld.Offers),
			"specifications":        tld.Specifications,
		},
	}
}

func flattenTldOffers(offers map[string]*domain.TldOffer) []map[string]interface{} {
	if offers == nil {
		return nil
	}

	flattenedOffers := []map[string]interface{}{}
	for _, offer := range offers {
		flattenedOffers = append(flattenedOffers, map[string]interface{}{
			"action":         offer.Action,
			"operation_path": offer.OperationPath,
			"price": map[string]interface{}{
				"currency_code": offer.Price.CurrencyCode,
				"units":         strconv.Itoa(int(offer.Price.Units)),
				"nanos":         strconv.Itoa(int(offer.Price.Nanos)),
			},
		})
	}

	return flattenedOffers
}
func flattenExternalDomainRegistrationStatus(status *domain.DomainRegistrationStatusExternalDomain) map[string]interface{} {
	if status == nil {
		return nil
	}

	return map[string]interface{}{
		"validation_token": status.ValidationToken,
	}
}

func flattenDNSZones(dnsZones []*domain.DNSZone) []map[string]interface{} {
	if dnsZones == nil {
		return nil
	}

	var zones []map[string]interface{}
	for _, zone := range dnsZones {
		zones = append(zones, map[string]interface{}{
			"domain":     zone.Domain,
			"subdomain":  zone.Subdomain,
			"ns":         zone.Ns,
			"ns_default": zone.NsDefault,
			"ns_master":  zone.NsMaster,
			"status":     zone.Status,
			"message":    zone.Message,
			"updated_at": zone.UpdatedAt.Format(time.RFC3339),
			"project_id": zone.ProjectID,
		})
	}

	return zones
}

func flattenDomainRegistrationStatusTransfer(transferStatus *domain.DomainRegistrationStatusTransfer) map[string]interface{} {
	if transferStatus == nil {
		return nil
	}

	return map[string]interface{}{
		"status":             string(transferStatus.Status),
		"vote_current_owner": transferStatus.VoteCurrentOwner,
		"vote_new_owner":     transferStatus.VoteNewOwner,
	}
}

func waitForOrderDomain(ctx context.Context, api *domain.RegistrarAPI, domainName string, timeout time.Duration) (*domain.Domain, error) {
	retryInterval := defaultWaitDomainOrderDomainRetryInterval
	if transport.DefaultWaitRetryInterval != nil {
		retryInterval = *transport.DefaultWaitRetryInterval
	}
	return api.WaitForOrderDomain(&domain.WaitForOrderDomainRequest{
		Domain:        domainName,
		Timeout:       scw.TimeDurationPtr(timeout),
		RetryInterval: &retryInterval,
	}, scw.WithContext(ctx))
}

func getStatusTasks(ctx context.Context, api *domain.RegistrarAPI, taskID string) (domain.TaskStatus, error) {
	var page int32 = 1
	var pageSize uint32 = 1000
	for {
		listTasksResponse, err := api.ListTasks(&domain.RegistrarAPIListTasksRequest{
			Page:     &page,
			PageSize: &pageSize,
		}, scw.WithContext(ctx))

		if err != nil {
			return "", fmt.Errorf("error retrieving tasks: %w", err)
		}

		for _, task := range listTasksResponse.Tasks {
			if task.ID == taskID {
				return task.Status, nil
			}
		}

		if len(listTasksResponse.Tasks) == 0 || uint32(len(listTasksResponse.Tasks)) < pageSize {
			break
		}

		page++
	}

	return "", fmt.Errorf("task with ID '%s' not found", taskID)
}

func waitForTaskCompletion(ctx context.Context, registrarAPI *domain.RegistrarAPI, taskID string, duration int) error {
	timeout := time.Duration(duration) * time.Second
	return retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		status, err := getStatusTasks(ctx, registrarAPI, taskID)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to retrieve task status: %w", err))
		}

		if status == domain.TaskStatusPending || status == domain.TaskStatusWaitingPayment || status == domain.TaskStatusNew {
			return retry.RetryableError(errors.New("task is not yet complete, retrying"))
		}

		if status == domain.TaskStatusSuccess {
			return nil
		}

		if status == domain.TaskStatusError {
			return retry.NonRetryableError(fmt.Errorf("task failed for domain: %s", taskID)) // Échec
		}

		return retry.NonRetryableError(fmt.Errorf("unexpected task status: %v", status))
	})
}

func ExpandDSRecord(dsRecordList []interface{}) *domain.DSRecord {
	if len(dsRecordList) == 0 || dsRecordList[0] == nil {
		return nil
	}

	dsRecordMap := dsRecordList[0].(map[string]interface{})
	dsRecord := &domain.DSRecord{
		KeyID:     uint32(dsRecordMap["key_id"].(int)),
		Algorithm: domain.DSRecordAlgorithm(dsRecordMap["algorithm"].(string)),
	}

	if digestList, ok := dsRecordMap["digest"].([]interface{}); ok && len(digestList) > 0 {
		digestMap := digestList[0].(map[string]interface{})
		dsRecord.Digest = &domain.DSRecordDigest{
			Type:   domain.DSRecordDigestType(digestMap["type"].(string)),
			Digest: digestMap["digest"].(string),
		}

		if publicKeyList, ok := digestMap["public_key"].([]interface{}); ok && len(publicKeyList) > 0 {
			publicKeyMap := publicKeyList[0].(map[string]interface{})
			dsRecord.Digest.PublicKey = &domain.DSRecordPublicKey{
				Key: publicKeyMap["key"].(string),
			}
		}
	}

	if publicKeyList, ok := dsRecordMap["public_key"].([]interface{}); ok && len(publicKeyList) > 0 {
		publicKeyMap := publicKeyList[0].(map[string]interface{})
		dsRecord.PublicKey = &domain.DSRecordPublicKey{
			Key: publicKeyMap["key"].(string),
		}
	}

	return dsRecord
}

func FlattenDSRecord(dsRecord *domain.DSRecord) []map[string]interface{} {
	if dsRecord == nil {
		return nil
	}

	result := map[string]interface{}{
		"key_id":    dsRecord.KeyID,
		"algorithm": string(dsRecord.Algorithm),
	}

	if dsRecord.Digest != nil {
		digest := map[string]interface{}{
			"type":   string(dsRecord.Digest.Type),
			"digest": dsRecord.Digest.Digest,
		}

		if dsRecord.Digest.PublicKey != nil {
			digest["public_key"] = []map[string]interface{}{
				{
					"key": dsRecord.Digest.PublicKey.Key,
				},
			}
		}

		result["digest"] = []map[string]interface{}{digest}
	}

	if dsRecord.PublicKey != nil {
		result["public_key"] = []map[string]interface{}{
			{
				"key": dsRecord.PublicKey.Key,
			},
		}
	}

	return []map[string]interface{}{result}
}
