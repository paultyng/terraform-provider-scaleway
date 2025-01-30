package domain_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	domainSDK "github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/domain"
)

func TestAccDomainRegistration_SingleDomainWithUpdate(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	singleDomain := "test-single-updates11" + ".com" // à adapter

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckDomainDestroy(tt),
		Steps: []resource.TestStep{
			{
				// Test initial : un seul domaine, owner_contact
				Config: fmt.Sprintf(`
                    resource "scaleway_domain_domains_registration" "test" {
                      domain_names = [ "%s"]
                      duration_in_years = 1

                      owner_contact {
                        firstname                   = "John"
                        lastname                    = "DOE"
                        email                       = "john.doe@example.com"
                        phone_number                = "+1.23456789"
                        address_line_1              = "123 Main Street"
                        city                        = "Paris"
                        zip                         = "75001"
                        country                     = "FR"
                        legal_form                  = "individual"
                        vat_identification_code     = "FR12345678901"
                        company_identification_code = "123456789"
                      }
                    }
                `, singleDomain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_domain_domains_registration.test", "domains_info.0.domain_name", singleDomain),
					resource.TestCheckResourceAttr("scaleway_domain_domains_registration.test", "duration_in_years", "1"),
				),
			},
			// Décommente ceci pour tester la mise à jour (administrative_contact + technical_contact)
			/*
			   {
			       Config: fmt.Sprintf(`
			           resource "scaleway_domain_domains_registration" "test" {
			             domains_info = [
			               { domain_name = "%s" }
			             ]
			             duration_in_years = 1

			             owner_contact {
			               firstname                   = "John"
			               lastname                    = "DOE"
			               email                       = "john.doe@example.com"
			               phone_number                = "+1.23456789"
			               address_line_1              = "123 Main Street"
			               city                        = "Paris"
			               zip                         = "75001"
			               country                     = "FR"
			               legal_form                  = "individual"
			               vat_identification_code     = "FR12345678901"
			               company_identification_code = "123456789"
			             }

			             administrative_contact {
			               firstname       = "Jane"
			               lastname        = "DOE"
			               email           = "jane.doe@example.com"
			               phone_number    = "+1.98765432"
			               address_line_1  = "456 Another Street"
			               city            = "Lyon"
			               zip             = "69002"
			               country         = "FR"
			               legal_form      = "individual"
			               vat_identification_code     = "FR12345678901"
			               company_identification_code = "123456789"
			             }

			             technical_contact {
			               firstname       = "Tech"
			               lastname        = "Support"
			               email           = "tech.support@example.com"
			               phone_number    = "+1.55555555"
			               address_line_1  = "789 Tech Road"
			               city            = "Marseille"
			               zip             = "13001"
			               country         = "FR"
			               legal_form      = "individual"
			               vat_identification_code     = "FR12345678901"
			               company_identification_code = "123456789"
			             }
			           }
			       `, singleDomain),
			       Check: resource.ComposeTestCheckFunc(
			           resource.TestCheckResourceAttr("scaleway_domain_domains_registration.test", "administrative_contact.0.firstname", "Jane"),
			           resource.TestCheckResourceAttr("scaleway_domain_domains_registration.test", "technical_contact.0.firstname", "Tech"),
			       ),
			   },
			*/
		},
	})
}

func TestAccDomainRegistration_MultipleDomainsNoUpdate(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	domainName1 := "test-multi-1.com"
	domainName2 := "test-multi-2.com"
	domainName3 := "test-multi-3.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckDomainDestroy(tt),
		Steps: []resource.TestStep{
			{
				// Création d'une ressource unique avec plusieurs domaines,
				// aucune mise à jour ensuite (1 seule étape).
				Config: fmt.Sprintf(`
                    resource "scaleway_domain_domains_registration" "multi" {
                      domains_info = [
                        { domain_name = "%s" },
                        { domain_name = "%s" },
                        { domain_name = "%s" }
                      ]

                      duration_in_years = 1

                      owner_contact {
                        firstname                   = "John"
                        lastname                    = "DOE"
                        email                       = "john.doe@example.com"
                        phone_number                = "+1.23456789"
                        address_line_1              = "123 Main Street"
                        city                        = "Paris"
                        zip                         = "75001"
                        country                     = "FR"
                        legal_form                  = "individual"
                        vat_identification_code     = "FR12345678901"
                        company_identification_code = "123456789"
                      }
                    }
                `, domainName1, domainName2, domainName3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_domain_domains_registration.multi", "domains_info.0.domain_name", domainName1),
					resource.TestCheckResourceAttr("scaleway_domain_domains_registration.multi", "domains_info.1.domain_name", domainName2),
					resource.TestCheckResourceAttr("scaleway_domain_domains_registration.multi", "domains_info.2.domain_name", domainName3),
				),
			},
		},
	})
}

func testAccCheckDomainDestroy(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_domain_domains_registration" {
				continue
			}

			registrarAPI := domain.NewRegistrarDomainAPI(tt.Meta)

			domainNames, err := domain.ExtractDomainsFromTaskID(nil, rs.Primary.ID, registrarAPI)
			if err != nil {
				return nil
			}

			for _, domainName := range domainNames {
				domainResp, getErr := registrarAPI.GetDomain(&domainSDK.RegistrarAPIGetDomainRequest{
					Domain: domainName,
				})
				if getErr != nil {
					if httperrors.Is404(getErr) {
						continue
					}
					return fmt.Errorf("failed to get domain details for %s: %w", domainName, getErr)
				}

				if domainResp.AutoRenewStatus != domainSDK.DomainFeatureStatusDisabled {
					return fmt.Errorf(
						"domain %s still exists, and auto-renew is not disabled (current: %s)",
						domainName,
						domainResp.AutoRenewStatus,
					)
				}
			}
		}
		return nil
	}
}
