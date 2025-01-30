package lb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	lbchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/lb/testfuncs"
)

func TestAccDataSourceBackend_Basic(t *testing.T) {
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
						name = "data-test-lb-backend"
						type = "LB-S"
					}

					resource "scaleway_lb_backend" "main" {
					  lb_id            = scaleway_lb.main.id
					  name             = "backend01"
					  forward_protocol = "http"
					  forward_port     = "80"
					}
					
					data "scaleway_lb_backend" "byID" {
						backend_id = "${scaleway_lb_backend.main.id}"
					}
					
					data "scaleway_lb_backend" "byName" {
						name = "${scaleway_lb_backend.main.name}"
						lb_id = "${scaleway_lb.main.id}"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.scaleway_lb_backend.byID", "name",
						"scaleway_lb_backend.main", "name"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_lb_backend.byName", "id",
						"scaleway_lb_backend.main", "id"),
				),
			},
		},
	})
}
