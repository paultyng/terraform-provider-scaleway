package domain_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func TestAccDataSourceDomainZone_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	testDNSZone := "test-zone2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckDomainZoneDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource scaleway_domain_zone main {
						domain    = "%s"
						subdomain = "%s"
					}

					data scaleway_domain_zone test {
						domain    = scaleway_domain_zone.main.domain
						subdomain = scaleway_domain_zone.main.subdomain
					}
				`, acctest.TestDomain, testDNSZone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainZoneExists(tt, "data.scaleway_domain_zone.test"),
					resource.TestCheckResourceAttr("data.scaleway_domain_zone.test", "subdomain", testDNSZone),
					resource.TestCheckResourceAttr("data.scaleway_domain_zone.test", "domain", acctest.TestDomain),
					resource.TestCheckResourceAttr("data.scaleway_domain_zone.test", "status", "active"),
				),
			},
		},
	})
}
