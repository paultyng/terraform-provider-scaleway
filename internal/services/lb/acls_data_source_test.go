package lb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func TestAccDataSourceACLs_Basic(t *testing.T) {
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
						inbound_port = 50000
						acl {
							name  = "tf-acl-datasource0"
							action {
								type = "allow"
							}
							match {
								ip_subnet = ["192.168.0.1"]
								http_filter = "acl_http_filter_none"
								http_filter_value = []
								invert = "true"
							}
						}
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
						inbound_port = 50000
						acl {
							name  = "tf-acl-datasource0"
							action {
								type = "allow"
							}
							match {
								ip_subnet = ["192.168.0.1"]
								http_filter = "acl_http_filter_none"
								http_filter_value = []
								invert = "true"
							}
						}
						acl {
							name  = "tf-acl-datasource1"
							action {
								type = "deny"
							}
							match {
								http_filter_value = []
								ip_subnet = ["0.0.0.0/0"]
							}
						}
					}
					data "scaleway_lb_acls" "byFrontID" {
						frontend_id = "${scaleway_lb_frontend.frt01.id}"
						depends_on = [scaleway_lb_frontend.frt01]
					}
					data "scaleway_lb_acls" "byFrontID_and_name" {
						frontend_id = "${scaleway_lb_frontend.frt01.id}"
						name = "tf-acl-datasource" 
						depends_on = [scaleway_lb_frontend.frt01]
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.scaleway_lb_acls.byFrontID", "acls.0.id"),
					resource.TestCheckResourceAttrSet("data.scaleway_lb_acls.byFrontID", "acls.1.id"),

					resource.TestCheckResourceAttrSet("data.scaleway_lb_acls.byFrontID_and_name", "acls.0.id"),
					resource.TestCheckResourceAttrSet("data.scaleway_lb_acls.byFrontID_and_name", "acls.1.id"),
				),
			},
		},
	})
}
