package scaleway

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccScalewayDataSourceLb_Basic(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayLbIPDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_lb_ip main {
					}

					resource scaleway_lb main {
					    ip_id = scaleway_lb_ip.main.id
						name = "test-lb"
						type = "LB-S"
					}
					
					data "scaleway_lb" "test" {
						lb_id = "${scaleway_lb.main.id}"
					}
					
					data "scaleway_lb" "test2" {
						name = "${scaleway_lb.main.name}"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayLbExists(tt, "data.scaleway_lb.test"),
					testAccCheckScalewayLbExists(tt, "data.scaleway_lb.test2"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_lb.test", "name",
						"scaleway_lb.main", "name"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_lb.test2", "id",
						"scaleway_lb.main", "id"),
				),
			},
		},
	})
}
