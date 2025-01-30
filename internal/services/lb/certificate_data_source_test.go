package lb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	lbchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/lb/testfuncs"
)

func TestAccDataSourceCertificate_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      lbchecks.IsIPDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_lb_ip main {
					}

					resource scaleway_lb main {
					    ip_id = scaleway_lb_ip.main.id
						name = "data-test-lb-cert"
						type = "LB-S"
					}

					resource scaleway_lb_certificate main {
						lb_id = scaleway_lb.main.id
						name = "data-test-lb-cert"
						letsencrypt {
							common_name = "${replace(scaleway_lb.main.ip_address, ".", "-")}.lb.${scaleway_lb.main.region}.scw.cloud"
						}
					}
					
					data "scaleway_lb_certificate" "byID" {
						certificate_id = "${scaleway_lb_certificate.main.id}"
					}
					
					data "scaleway_lb_certificate" "byName" {
						name = "${scaleway_lb_certificate.main.name}"
						lb_id = "${scaleway_lb.main.id}"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.scaleway_lb_certificate.byID", "name",
						"scaleway_lb_certificate.main", "name"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_lb_certificate.byName", "id",
						"scaleway_lb_certificate.main", "id"),
				),
			},
		},
	})
}
