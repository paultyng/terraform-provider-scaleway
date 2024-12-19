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

const domainName = "test-basic-domain71.com"

func TestAccOrderDomain_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckOrderDomainDestroy(tt),
		Steps: []resource.TestStep{
			{
				// Test initial configuration
				Config: fmt.Sprintf(`
				resource "scaleway_domain_order_domain" "test" {
					domain_name       = "%s"
					duration_in_years = %d

					owner_contact {
						firstname       = "John"
						lastname        = "DOE"
						email           = "john.doe@example.com"
						phone_number    = "+1.23456789"
						address_line_1  = "123 Main Street"
						city            = "Paris"
						zip             = "75001"
						country         = "FR"
						legal_form      = "individual"
						vat_identification_code  = "FR12345678901"
						company_identification_code = "123456789"
					}
				}`, domainName, 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_domain_order_domain.test", "domain_name", domainName),
					resource.TestCheckResourceAttr("scaleway_domain_order_domain.test", "duration_in_years", "1"),
					resource.TestCheckResourceAttrSet("scaleway_domain_order_domain.test", "id"),
					resource.TestCheckResourceAttrSet("scaleway_domain_order_domain.test", "owner_contact.0.firstname"),
				),
			},
			{
				// Update the administrative and technical contacts
				Config: fmt.Sprintf(`
				resource "scaleway_domain_order_domain" "test" {
					domain_name       = "%s"
					duration_in_years = %d

					owner_contact {
						firstname       = "John"
						lastname        = "DOE"
						email           = "john.doe@example.com"
						phone_number    = "+1.23456789"
						address_line_1  = "123 Main Street"
						city            = "Paris"
						zip             = "75001"
						country         = "FR"
						legal_form      = "individual"
						vat_identification_code  = "FR12345678901"
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
						vat_identification_code  = "FR12345678901"
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
						vat_identification_code  = "FR12345678901"
						company_identification_code = "123456789"
					}
				}`, domainName, 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_domain_order_domain.test", "administrative_contact.0.firstname", "Jane"),
					resource.TestCheckResourceAttr("scaleway_domain_order_domain.test", "technical_contact.0.firstname", "Tech"),
				),
			},
		},
	})
}

func TestAccOrderDomain_Multiple_Domains(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	domainName1 := "test-basic-domain250.com"
	domainName2 := "test-basic-domain251.com"
	domainName3 := "test-basic-domain252.com"
	domainName4 := "test-basic-domain253.com"
	domainName5 := "test-basic-domain254.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckOrderDomainDestroy(tt),
		Steps: []resource.TestStep{
			{
				// Test initial configuration
				Config: fmt.Sprintf(`
					variable "domain_orders" {
					  type = list(object({
						domain_name       = string
						duration_in_years = number
					  }))
					  default = [
						{ domain_name = "%s", duration_in_years = 1 },
						{ domain_name = "%s", duration_in_years = 1 },
						{ domain_name = "%s", duration_in_years = 1 },
						{ domain_name = "%s", duration_in_years = 1 }
					  ]
					}
					
					resource "scaleway_domain_order_domain" "first" {
					  domain_name       = "test-basic-domain503.com"
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
					
					resource "scaleway_domain_order_domain" "test" {
					  for_each = { for i, domain in var.domain_orders : i => domain }
					
					  domain_name       = each.value.domain_name
					  duration_in_years = each.value.duration_in_years
					
					  owner_contact_id = scaleway_domain_order_domain.first.owner_contact_id
					
					   depends_on = [
						 scaleway_domain_order_domain.first,
    					 lookup(scaleway_domain_order_domain, each.key - 1, null)
						]
					}`, domainName1, domainName2, domainName3, domainName4),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_domain_order_domain.test[\"0\"]", "domain_name", domainName1),
					resource.TestCheckResourceAttr("scaleway_domain_order_domain.test[\"1\"]", "domain_name", domainName2),
					resource.TestCheckResourceAttr("scaleway_domain_order_domain.test[\"2\"]", "domain_name", domainName3),
					resource.TestCheckResourceAttr("scaleway_domain_order_domain.test[\"3\"]", "domain_name", domainName4),
					resource.TestCheckResourceAttr("scaleway_domain_order_domain.test[\"4\"]", "domain_name", domainName5),
				),
			},
		},
	})
}

func testAccCheckOrderDomainDestroy(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_order_domain" {
				continue
			}

			registrarAPI := domain.NewRegistrarDomainAPI(tt.Meta)
			domainName, err := domain.ExtractDomainFromID(rs.Primary.ID)
			domainResp, err := registrarAPI.GetDomain(&domainSDK.RegistrarAPIGetDomainRequest{
				Domain: domainName,
			})

			if err != nil {
				if httperrors.Is404(err) {
					return nil
				}
				return fmt.Errorf("failed to get domain details: %w", err)
			}

			if domainResp.AutoRenewStatus != domainSDK.DomainFeatureStatusDisabled {
				return fmt.Errorf("expected auto-renew to be 'disabled' for domain %s, got %s",
					domainName, domainResp.AutoRenewStatus)
			}
			return fmt.Errorf("domain %s still exists and auto-renew is not 'disabled'", domainName)
		}

		return nil
	}
}
