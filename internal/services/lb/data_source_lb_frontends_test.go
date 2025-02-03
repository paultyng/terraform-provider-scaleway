package lb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func TestAccDataSourceFrontends_Basic(t *testing.T) {
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
						forward_protocol = "http"
						forward_port = 80
						proxy_protocol = "none"
					}
					resource scaleway_lb_frontend frt01 {
						lb_id = scaleway_lb.lb01.id
						backend_id = scaleway_lb_backend.bkd01.id
						name  = "tf-frontend-datasource0"
						inbound_port = 50000
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
						forward_protocol = "http"
						forward_port = 80
						proxy_protocol = "none"
					}
					resource scaleway_lb_frontend frt01 {
						lb_id = scaleway_lb.lb01.id
						backend_id = scaleway_lb_backend.bkd01.id
						name  = "tf-frontend-datasource0"
						inbound_port = 50000
					}
					resource scaleway_lb_frontend frt02 {
						lb_id = scaleway_lb.lb01.id
						backend_id = scaleway_lb_backend.bkd01.id
						name  = "tf-frontend-datasource1"
						inbound_port = 50001
					}
					data "scaleway_lb_frontends" "byLBID" {
						lb_id = "${scaleway_lb.lb01.id}"
						depends_on = [scaleway_lb_frontend.frt01, scaleway_lb_frontend.frt02]
					}
					data "scaleway_lb_frontends" "byLBID_and_name" {
						lb_id = "${scaleway_lb.lb01.id}"
						name = "tf-frontend-datasource" 
						depends_on = [scaleway_lb_frontend.frt01, scaleway_lb_frontend.frt02]
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.scaleway_lb_frontends.byLBID", "frontends.0.id"),
					resource.TestCheckResourceAttrSet("data.scaleway_lb_frontends.byLBID", "frontends.1.id"),

					resource.TestCheckResourceAttrSet("data.scaleway_lb_frontends.byLBID_and_name", "frontends.0.id"),
					resource.TestCheckResourceAttrSet("data.scaleway_lb_frontends.byLBID_and_name", "frontends.1.id"),
				),
			},
		},
	})
}
