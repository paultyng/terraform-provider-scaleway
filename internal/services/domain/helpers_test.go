package domain_test

import (
	"testing"

	domainSDK "github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/domain"
	"github.com/stretchr/testify/assert"
)

func TestExpandContact(t *testing.T) {
	tests := []struct {
		name        string
		contactMap  map[string]interface{}
		expected    *domainSDK.Contact
		expectError bool
	}{
		{
			name: "minimal valid contact",
			contactMap: map[string]interface{}{
				"phone_number":   "123456789",
				"legal_form":     "individual",
				"firstname":      "John",
				"lastname":       "Doe",
				"email":          "john.doe@example.com",
				"address_line_1": "123 Main St",
				"zip":            "75001",
				"city":           "Paris",
				"country":        "FR",
			},
			expected: &domainSDK.Contact{
				PhoneNumber:  "123456789",
				LegalForm:    domainSDK.ContactLegalForm("individual"),
				Firstname:    "John",
				Lastname:     "Doe",
				Email:        "john.doe@example.com",
				AddressLine1: "123 Main St",
				Zip:          "75001",
				City:         "Paris",
				Country:      "FR",
			},
		},
		{
			name: "full contact with extensions",
			contactMap: map[string]interface{}{
				"phone_number":   "987654321",
				"legal_form":     "corporate",
				"firstname":      "Jane",
				"lastname":       "Doe",
				"email":          "jane.doe@example.com",
				"address_line_1": "456 Secondary St",
				"zip":            "10001",
				"city":           "New York",
				"country":        "US",
				"extension_fr": map[string]interface{}{
					"mode": "individual",
					"individual_info": map[string]interface{}{
						"whois_opt_in": true,
					},
				},
			},
			expected: &domainSDK.Contact{
				PhoneNumber:  "987654321",
				LegalForm:    domainSDK.ContactLegalForm("corporate"),
				Firstname:    "Jane",
				Lastname:     "Doe",
				Email:        "jane.doe@example.com",
				AddressLine1: "456 Secondary St",
				Zip:          "10001",
				City:         "New York",
				Country:      "US",
				ExtensionFr: &domainSDK.ContactExtensionFR{
					Mode:           domainSDK.ContactExtensionFRModeIndividual,
					IndividualInfo: &domainSDK.ContactExtensionFRIndividualInfo{WhoisOptIn: true},
				},
			},
		},
		{
			name:        "nil input map",
			contactMap:  nil,
			expected:    nil,
			expectError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := domain.ExpandContact(tt.contactMap)

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExpandNewContact(t *testing.T) {
	tests := []struct {
		name        string
		contactMap  map[string]interface{}
		expected    *domainSDK.NewContact
		expectError bool
	}{
		{
			name: "minimal valid new contact",
			contactMap: map[string]interface{}{
				"phone_number":   "123456789",
				"legal_form":     "individual",
				"firstname":      "John",
				"lastname":       "Doe",
				"email":          "john.doe@example.com",
				"address_line_1": "123 Main St",
				"zip":            "75001",
				"city":           "Paris",
				"country":        "FR",
			},
			expected: &domainSDK.NewContact{
				PhoneNumber:  "123456789",
				LegalForm:    domainSDK.ContactLegalForm("individual"),
				Firstname:    "John",
				Lastname:     "Doe",
				Email:        "john.doe@example.com",
				AddressLine1: "123 Main St",
				Zip:          "75001",
				City:         "Paris",
				Country:      "FR",
			},
		},
		{
			name: "new contact with optional fields",
			contactMap: map[string]interface{}{
				"phone_number":                "987654321",
				"legal_form":                  "corporate",
				"firstname":                   "Jane",
				"lastname":                    "Doe",
				"email":                       "jane.doe@example.com",
				"address_line_1":              "456 Secondary St",
				"zip":                         "10001",
				"city":                        "New York",
				"country":                     "US",
				"company_name":                "Acme Inc.",
				"email_alt":                   "jane.alt@example.com",
				"fax_number":                  "+123456789",
				"address_line_2":              "Suite 101",
				"vat_identification_code":     "VAT123",
				"company_identification_code": "C123",
				"state":                       "NY",
				"whois_opt_in":                true,
				"resale":                      true,
			},
			expected: &domainSDK.NewContact{
				PhoneNumber:               "987654321",
				LegalForm:                 domainSDK.ContactLegalForm("corporate"),
				Firstname:                 "Jane",
				Lastname:                  "Doe",
				Email:                     "jane.doe@example.com",
				AddressLine1:              "456 Secondary St",
				Zip:                       "10001",
				City:                      "New York",
				Country:                   "US",
				CompanyName:               scw.StringPtr("Acme Inc."),
				EmailAlt:                  scw.StringPtr("jane.alt@example.com"),
				FaxNumber:                 scw.StringPtr("+123456789"),
				AddressLine2:              scw.StringPtr("Suite 101"),
				VatIDentificationCode:     scw.StringPtr("VAT123"),
				CompanyIDentificationCode: scw.StringPtr("C123"),
				State:                     scw.StringPtr("NY"),
				WhoisOptIn:                true,
				Resale:                    true,
			},
		},
		{
			name: "new contact with extensions",
			contactMap: map[string]interface{}{
				"phone_number":   "654987321",
				"legal_form":     "individual",
				"firstname":      "Alice",
				"lastname":       "Smith",
				"email":          "alice.smith@example.com",
				"address_line_1": "789 Tertiary Ave",
				"zip":            "30301",
				"city":           "Atlanta",
				"country":        "US",
				"extension_fr": map[string]interface{}{
					"mode": "individual",
					"individual_info": map[string]interface{}{
						"whois_opt_in": true,
					},
				},
			},
			expected: &domainSDK.NewContact{
				PhoneNumber:  "654987321",
				LegalForm:    domainSDK.ContactLegalForm("individual"),
				Firstname:    "Alice",
				Lastname:     "Smith",
				Email:        "alice.smith@example.com",
				AddressLine1: "789 Tertiary Ave",
				Zip:          "30301",
				City:         "Atlanta",
				Country:      "US",
				ExtensionFr: &domainSDK.ContactExtensionFR{
					Mode:           domainSDK.ContactExtensionFRModeIndividual,
					IndividualInfo: &domainSDK.ContactExtensionFRIndividualInfo{WhoisOptIn: true},
				},
			},
		},
		{
			name:        "nil input map",
			contactMap:  nil,
			expected:    nil,
			expectError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := domain.ExpandNewContact(tt.contactMap)

			assert.Equal(t, tt.expected, result)
		})
	}
}
