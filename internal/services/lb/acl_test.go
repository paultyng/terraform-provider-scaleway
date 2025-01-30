package lb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	lbSDK "github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/lb"
)

func TestAccAcl_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isACLDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_lb_ip ip01 {}
					resource scaleway_lb lb01 {
						ip_id = scaleway_lb_ip.ip01.id
						name = "test-lb-acl"
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
						name = "tf-test"
						inbound_port = 80
						timeout_client = "30s"
						external_acls = true
					}
					resource scaleway_lb_acl acl01 {
						frontend_id = scaleway_lb_frontend.frt01.id
						name = "test-acl-basic"
						description = "a description"
						index = 4
						action {
							type = "allow"
						}
						match {
							ip_subnet = ["192.168.0.1", "192.168.0.2", "192.168.10.0/24"]
							http_filter = "acl_http_filter_none"
							http_filter_value = []
							invert = "true"
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isACLPresent(tt, "scaleway_lb_acl.acl01"),
					resource.TestCheckResourceAttrPair(
						"scaleway_lb_acl.acl01", "frontend_id",
						"scaleway_lb_frontend.frt01", "id"),
					resource.TestCheckResourceAttr("scaleway_lb_acl.acl01", "name", "test-acl-basic"),
					resource.TestCheckResourceAttr("scaleway_lb_acl.acl01", "description", "a description"),
					resource.TestCheckResourceAttr("scaleway_lb_acl.acl01", "index", "4"),
					resource.TestCheckResourceAttr("scaleway_lb_acl.acl01", "action.0.type", "allow"),
					resource.TestCheckResourceAttr("scaleway_lb_acl.acl01", "match.0.ip_subnet.#", "3"),
					resource.TestCheckResourceAttr("scaleway_lb_acl.acl01", "match.0.ip_subnet.0", "192.168.0.1"),
					resource.TestCheckResourceAttr("scaleway_lb_acl.acl01", "match.0.ip_subnet.1", "192.168.0.2"),
					resource.TestCheckResourceAttr("scaleway_lb_acl.acl01", "match.0.ip_subnet.2", "192.168.10.0/24"),
					resource.TestCheckResourceAttr("scaleway_lb_acl.acl01", "match.0.http_filter", "acl_http_filter_none"),
					resource.TestCheckResourceAttr("scaleway_lb_acl.acl01", "match.0.http_filter_value.#", "0"),
					resource.TestCheckResourceAttr("scaleway_lb_acl.acl01", "match.0.invert", "true"),
				),
			},
			{
				Config: `
					resource scaleway_lb_ip ip01 {}
					resource scaleway_lb lb01 {
						ip_id = scaleway_lb_ip.ip01.id
						name = "test-lb-acl"
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
						name = "tf-test"
						inbound_port = 80
						timeout_client = "30s"
						external_acls = true
					}
					resource scaleway_lb_acl acl01 {
						frontend_id = scaleway_lb_frontend.frt01.id
						name = "updated-test-acl-basic"
						description = "updated description"
						index = 3
						action {
							type = "deny"
						}
						match {
							ip_subnet = ["0.0.0.0/0"]
							http_filter = "acl_http_filter_none"
							http_filter_value = []
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isACLPresent(tt, "scaleway_lb_acl.acl01"),
					resource.TestCheckResourceAttrPair(
						"scaleway_lb_acl.acl01", "frontend_id",
						"scaleway_lb_frontend.frt01", "id"),
					resource.TestCheckResourceAttr("scaleway_lb_acl.acl01", "name", "updated-test-acl-basic"),
					resource.TestCheckResourceAttr("scaleway_lb_acl.acl01", "description", "updated description"),
					resource.TestCheckResourceAttr("scaleway_lb_acl.acl01", "index", "3"),
					resource.TestCheckResourceAttr("scaleway_lb_acl.acl01", "action.0.type", "deny"),
					resource.TestCheckResourceAttr("scaleway_lb_acl.acl01", "match.0.ip_subnet.#", "1"),
					resource.TestCheckResourceAttr("scaleway_lb_acl.acl01", "match.0.ip_subnet.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("scaleway_lb_acl.acl01", "match.0.http_filter", "acl_http_filter_none"),
					resource.TestCheckResourceAttr("scaleway_lb_acl.acl01", "match.0.http_filter_value.#", "0"),
					resource.TestCheckResourceAttr("scaleway_lb_acl.acl01", "match.0.invert", "false"),
				),
			},
			{
				Config: `
					resource scaleway_lb_ip ip01 {}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("scaleway_lb_ip.ip01", "id"),
				),
			},
		},
	})
}

func isACLPresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		lbAPI, zone, ID, err := lb.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = lbAPI.GetACL(&lbSDK.ZonedAPIGetACLRequest{
			ACLID: ID,
			Zone:  zone,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func isACLDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_lb_acl" {
				continue
			}

			lbAPI, zone, ID, err := lb.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = lbAPI.GetACL(&lbSDK.ZonedAPIGetACLRequest{
				Zone:  zone,
				ACLID: ID,
			})

			// If no error resource still exist
			if err == nil {
				return fmt.Errorf("LB ACL (%s) still exists", rs.Primary.ID)
			}

			// Unexpected api error we return it
			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}
