package lb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	lbchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/lb/testfuncs"
)

func TestAccDataSourceRoute_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			lbchecks.IsIPDestroyed(tt),
			isLbDestroyed(tt),
		),
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
					
					data "scaleway_lb_route" "byID" {
						route_id = "${scaleway_lb_route.rt01.id}"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.scaleway_lb_route.byID", "name",
						"scaleway_lb_route.rt01", "name"),
				),
			},
		},
	})
}
