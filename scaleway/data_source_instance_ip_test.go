package scaleway_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func TestAccScalewayDataSourceInstanceIP_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayInstanceServerDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `resource "scaleway_instance_ip" "ip" {}`,
			},
			{
				Config: `
					resource "scaleway_instance_ip" "ip" {}

					data "scaleway_instance_ip" "ip-from-address" {
						address = "${scaleway_instance_ip.ip.address}"
					}

					data "scaleway_instance_ip" "ip-from-id" {
						id = "${scaleway_instance_ip.ip.id}"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceAttrIP("scaleway_instance_ip.ip", "address"),
					testCheckResourceAttrIP("data.scaleway_instance_ip.ip-from-address", "address"),
					testCheckResourceAttrIP("data.scaleway_instance_ip.ip-from-id", "address"),
					resource.TestCheckResourceAttrPair("scaleway_instance_ip.ip", "address", "data.scaleway_instance_ip.ip-from-address", "address"),
					resource.TestCheckResourceAttrPair("scaleway_instance_ip.ip", "address", "data.scaleway_instance_ip.ip-from-id", "address"),
				),
			},
		},
	})
}
