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

func TestAccOrderDomain_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	domainName := "test-basic-domain56.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckOrderDomainDestroy(tt),
		Steps: []resource.TestStep{
			{
				//test and doc with for each for multiple domains
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
					resource.TestCheckResourceAttrSet("scaleway_domain_order_domain.test", "auto_renew_status"),
					resource.TestCheckResourceAttrSet("scaleway_domain_order_domain.test", "dnssec_status"),
					resource.TestCheckResourceAttr("scaleway_domain_order_domain.test", "epp_code.0", "clientTransferProhibited"),
					resource.TestCheckResourceAttrSet("scaleway_domain_order_domain.test", "expired_at"),
					resource.TestCheckResourceAttrSet("scaleway_domain_order_domain.test", "updated_at"),
					resource.TestCheckResourceAttrSet("scaleway_domain_order_domain.test", "registrar"),
					resource.TestCheckResourceAttrSet("scaleway_domain_order_domain.test", "is_external"),
					resource.TestCheckResourceAttrSet("scaleway_domain_order_domain.test", "status"),
					resource.TestCheckResourceAttrSet("scaleway_domain_order_domain.test", "organization_id"),
					resource.TestCheckResourceAttrSet("scaleway_domain_order_domain.test", "pending_trade"),

					resource.TestCheckResourceAttrSet("scaleway_domain_order_domain.test", "tld.0.name"),
					resource.TestCheckResourceAttrSet("scaleway_domain_order_domain.test", "tld.0.dnssec_support"),
					resource.TestCheckResourceAttrSet("scaleway_domain_order_domain.test", "dns_zones.0.domain"),
					resource.TestCheckResourceAttrSet("scaleway_domain_order_domain.test", "dns_zones.0.status"),
					resource.TestCheckResourceAttrSet("scaleway_domain_order_domain.test", "dns_zones.0.ns_default.0"),
					resource.TestCheckResourceAttr("scaleway_domain_order_domain.test", "linked_products.#", "0"),
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
