package domain_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	domainSDK "github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/logging"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/domain"
)

func TestAccDomainZone_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	testDNSZone := "test-zone"
	logging.L.Debugf("TestAccScalewayDomainZone_Basic: test dns zone: %s, with domain: %s", testDNSZone, acctest.TestDomain)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckDomainZoneDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_domain_zone" "test" {
						domain    = "%s"
						subdomain = "%s"
					}
				`, acctest.TestDomain, testDNSZone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainZoneExists(tt, "scaleway_domain_zone.test"),
					resource.TestCheckResourceAttr("scaleway_domain_zone.test", "subdomain", testDNSZone),
					resource.TestCheckResourceAttr("scaleway_domain_zone.test", "domain", acctest.TestDomain),
					resource.TestCheckResourceAttr("scaleway_domain_zone.test", "status", "active"),
				),
			},
		},
	})
}

func testAccCheckDomainZoneExists(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		domainAPI := domain.NewDomainAPI(tt.Meta)
		listDNSZones, err := domainAPI.ListDNSZones(&domainSDK.ListDNSZonesRequest{
			DNSZones: []string{fmt.Sprintf("%s.%s", rs.Primary.Attributes["subdomain"], rs.Primary.Attributes["domain"])},
		})
		if err != nil {
			return err
		}

		if len(listDNSZones.DNSZones) == 0 {
			return fmt.Errorf("zone (%s) not found in: %s",
				rs.Primary.Attributes["subdomain"],
				rs.Primary.Attributes["domain"],
			)
		}

		return nil
	}
}

func testAccCheckDomainZoneDestroy(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_domain_zone" {
				continue
			}

			// check if the zone still exists
			domainAPI := domain.NewDomainAPI(tt.Meta)
			listDNSZones, err := domainAPI.ListDNSZones(&domainSDK.ListDNSZonesRequest{
				DNSZones: []string{fmt.Sprintf("%s.%s", rs.Primary.Attributes["subdomain"], rs.Primary.Attributes["domain"])},
			})

			if httperrors.Is403(err) { // forbidden: subdomain not found
				return nil
			}

			if err != nil {
				return err
			}

			if listDNSZones.TotalCount > 0 {
				return fmt.Errorf("zone %s still exist for domain: %s",
					rs.Primary.Attributes["subdomain"],
					rs.Primary.Attributes["domain"])
			}
			return nil
		}

		return nil
	}
}
