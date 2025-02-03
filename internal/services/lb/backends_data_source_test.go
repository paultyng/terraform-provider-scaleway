package lb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func TestAccDataSourceBackends_Basic(t *testing.T) {
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
						name  = "tf-backend-datasource0"
						forward_protocol = "tcp"
						forward_port = 80
						proxy_protocol = "none"
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
						name  = "tf-backend-datasource0"
						forward_protocol = "tcp"
						forward_port = 80
						proxy_protocol = "none"
					}
					resource scaleway_lb_backend bkd02 {
						lb_id = scaleway_lb.lb01.id
						name  = "tf-backend-datasource1"
						forward_protocol = "http"
						forward_port = 80
						proxy_protocol = "none"
					}
					data "scaleway_lb_backends" "byLBID" {
						lb_id = "${scaleway_lb.lb01.id}"
						depends_on = [scaleway_lb_backend.bkd01, scaleway_lb_backend.bkd02]
					}
					data "scaleway_lb_backends" "byLBID_and_name" {
						lb_id = "${scaleway_lb.lb01.id}"
						name = "tf-backend-datasource" 
						depends_on = [scaleway_lb_backend.bkd01, scaleway_lb_backend.bkd02]
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.scaleway_lb_backends.byLBID", "backends.0.id"),
					resource.TestCheckResourceAttrSet("data.scaleway_lb_backends.byLBID", "backends.1.id"),

					resource.TestCheckResourceAttrSet("data.scaleway_lb_backends.byLBID_and_name", "backends.0.id"),
					resource.TestCheckResourceAttrSet("data.scaleway_lb_backends.byLBID_and_name", "backends.1.id"),
				),
			},
		},
	})
}
