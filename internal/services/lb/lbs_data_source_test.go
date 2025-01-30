package lb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func TestAccDataSourceLbs_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isLbDestroyed(tt),
		Steps: []resource.TestStep{
			{
				// Create one IP first because its POST request cannot be matched correctly
				// There is no difference between the two IPS
				Config: `
					resource scaleway_lb_ip ip1 {
					}
				`,
			},
			{
				Config: `
					resource scaleway_lb_ip ip1 {
					}

					resource scaleway_lb_ip ip2 {
					}

					resource scaleway_lb lb1 {
						ip_id = scaleway_lb_ip.ip1.id
						name  = "tf-lb-datasource0"
						description = "a description"
						type = "LB-S"
						tags = [ "terraform-test", "data_scaleway_lbs", "basic" ]
					}
					`,
			},
			{
				Config: `
					resource scaleway_lb_ip ip1 {
					}

					resource scaleway_lb_ip ip2 {
					}

					resource scaleway_lb lb1 {
						ip_id = scaleway_lb_ip.ip1.id
						name  = "tf-lb-datasource0"
						description = "a description"
						type = "LB-S"
						tags = [ "terraform-test", "data_scaleway_lbs", "basic" ]
					}

					resource scaleway_lb lb2 {
						ip_id = scaleway_lb_ip.ip2.id
						name  = "tf-lb-datasource1"
						description = "a description"
						type = "LB-S"
						tags = [ "terraform-test", "data_scaleway_lbs", "basic" ]
					}
					`,
			},
			{
				Config: `
					resource scaleway_lb_ip ip1 {
					}

					resource scaleway_lb_ip ip2 {
					}

					resource scaleway_lb lb1 {
						ip_id = scaleway_lb_ip.ip1.id
						name  = "tf-lb-datasource0"
						description = "a description"
						type = "LB-S"
						tags = [ "terraform-test", "data_scaleway_lbs", "basic" ]
					}

					resource scaleway_lb lb2 {
						ip_id = scaleway_lb_ip.ip2.id
						name  = "tf-lb-datasource1"
						description = "a description"
						type = "LB-S"
						tags = [ "terraform-test", "data_scaleway_lbs", "basic" ]
					}

					data "scaleway_lbs" "lbs_by_name" {
						name  = "tf-lb-datasource"
						depends_on = [scaleway_lb.lb1, scaleway_lb.lb2]
					}

					data "scaleway_lbs" "lbs_by_tags" {
						tags = [ "terraform-test", "data_scaleway_lbs" ]
						depends_on = [scaleway_lb.lb1, scaleway_lb.lb2]
					}

					data "scaleway_lbs" "lbs_by_name_other_zone" {
						name  = "tf-lb-datasource"
						zone = "fr-par-2"
						depends_on = [scaleway_lb.lb1, scaleway_lb.lb2]
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.scaleway_lbs.lbs_by_name", "lbs.0.id"),
					resource.TestCheckResourceAttrSet("data.scaleway_lbs.lbs_by_name", "lbs.1.id"),
					resource.TestCheckResourceAttrSet("data.scaleway_lbs.lbs_by_name", "lbs.0.name"),
					resource.TestCheckResourceAttrSet("data.scaleway_lbs.lbs_by_name", "lbs.1.name"),

					resource.TestCheckResourceAttrSet("data.scaleway_lbs.lbs_by_tags", "lbs.0.id"),
					resource.TestCheckResourceAttrSet("data.scaleway_lbs.lbs_by_tags", "lbs.1.id"),

					resource.TestCheckNoResourceAttr("data.scaleway_lbs.lbs_by_name_other_zone", "lbs.0.id"),
				),
			},
		},
	})
}
