package lb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func TestAccDataSourceRoutes_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isLbDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_lb_ip ip01 {}
					resource scaleway_lb lb01 {
						ip_id = scaleway_lb_ip.ip01.id
						name = "test-lb"
						type = "lb-s"
					}
					resource scaleway_lb_backend bkd01 {
						lb_id = scaleway_lb.lb01.id
						forward_protocol = "tcp"
						forward_port = 80
						proxy_protocol = "none"
					}
					resource scaleway_lb_frontend frt01 {
						lb_id = scaleway_lb.lb01.id
						backend_id = scaleway_lb_backend.bkd01.id
						inbound_port = 80
					}
					resource scaleway_lb_route rt01 {
						frontend_id = scaleway_lb_frontend.frt01.id
						backend_id = scaleway_lb_backend.bkd01.id
						match_sni = "sni.scaleway.com"
					}
				`,
			},
			{
				Config: `
					resource scaleway_lb_ip ip01 {}
					resource scaleway_lb lb01 {
						ip_id = scaleway_lb_ip.ip01.id
						name = "test-lb"
						type = "lb-s"
					}
					resource scaleway_lb_backend bkd01 {
						lb_id = scaleway_lb.lb01.id
						forward_protocol = "tcp"
						forward_port = 80
						proxy_protocol = "none"
					}
					resource scaleway_lb_frontend frt01 {
						lb_id = scaleway_lb.lb01.id
						backend_id = scaleway_lb_backend.bkd01.id
						inbound_port = 80
					}
					resource scaleway_lb_route rt01 {
						frontend_id = scaleway_lb_frontend.frt01.id
						backend_id = scaleway_lb_backend.bkd01.id
						match_sni = "sni.scaleway.com"
					}
					resource scaleway_lb_route rt02 {
						frontend_id = scaleway_lb_frontend.frt01.id
						backend_id = scaleway_lb_backend.bkd01.id
						match_sni = "sni2.scaleway.com"
					}
					data "scaleway_lb_routes" "by_frontendID" {
						frontend_id = "${scaleway_lb_frontend.frt01.id}"
						depends_on = [scaleway_lb_route.rt01, scaleway_lb_route.rt02]
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.scaleway_lb_routes.by_frontendID", "routes.0.id"),
					resource.TestCheckResourceAttrSet("data.scaleway_lb_routes.by_frontendID", "routes.1.id"),
				),
			},
		},
	})
}
