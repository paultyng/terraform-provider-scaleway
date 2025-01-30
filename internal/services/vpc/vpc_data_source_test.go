package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func TestAccDataSourceVPC_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	vpcName := "DataSourceVPC_Basic"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckVPCDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_vpc" "vpc01" {
					  name = "%s"
					}`, vpcName),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_vpc" "vpc01" {
					  name = "%s"
					}

					data "scaleway_vpc" "by_name" {
						name = "${scaleway_vpc.vpc01.name}"
					}

					data "scaleway_vpc" "by_id" {
						vpc_id = "${scaleway_vpc.vpc01.id}"
					}
				`, vpcName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCExists(tt, "scaleway_vpc.vpc01"),
					resource.TestCheckResourceAttrPair("data.scaleway_vpc.by_name", "vpc_id", "scaleway_vpc.vpc01", "id"),
					resource.TestCheckResourceAttrPair("data.scaleway_vpc.by_id", "name", "scaleway_vpc.vpc01", "name"),
				),
			},
		},
	})
}

func TestAccDataSourceVPC_Default(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					data "scaleway_vpc" "default" {
						is_default = true
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.scaleway_vpc.default", "id"),
					resource.TestCheckResourceAttr("data.scaleway_vpc.default", "name", "default"),
					resource.TestCheckResourceAttr("data.scaleway_vpc.default", "is_default", "true"),
					resource.TestCheckResourceAttr("data.scaleway_vpc.default", "tags.0", "default"),
				),
			},
		},
	})
}
