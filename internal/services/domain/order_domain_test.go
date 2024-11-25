package domain_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func TestAccOrderDomain_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	domainName := "test-basic-domain11.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckOrderDomainDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "scaleway_domain_order_domain" "test" {
				domain_name       = "%s"
				duration_in_years = %d

				owner_contact {
					  firstname       = "John"
					  lastname        = "Doe"
					  email           = "john.doe@example.com"
					  phone_number    = "+123456789"
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
					resource.TestCheckResourceAttr("scaleway_order_domain.test", "domain_name", domainName),
					resource.TestCheckResourceAttr("scaleway_order_domain.test", "duration_in_years", "1"),
					resource.TestCheckResourceAttrSet("scaleway_order_domain.test", "id"),
					resource.TestCheckResourceAttrSet("scaleway_order_domain.test", "auto_renew_status"),
					resource.TestCheckResourceAttrSet("scaleway_order_domain.test", "status"),
					resource.TestCheckResourceAttrSet("scaleway_order_domain.test", "created_at"),
					resource.TestCheckResourceAttrSet("scaleway_order_domain.test", "updated_at"),
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

		}

		return nil
	}
}
