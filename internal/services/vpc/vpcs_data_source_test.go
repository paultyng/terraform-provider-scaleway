package vpc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	instancechecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/instance/testfuncs"
)

func TestAccDataSourceVPCs_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_vpc vpc01 {
						name = "tf-vpc-datasource0"
						tags = [ "terraform-test", "data_scaleway_vpcs", "basic" ]
					}`,
			},
			{
				Config: `
					resource scaleway_vpc vpc01 {
						name = "tf-vpc-datasource0"
						tags = [ "terraform-test", "data_scaleway_vpcs", "basic" ]
					}
				
					resource scaleway_vpc vpc02 {
						name = "tf-vpc-datasource1"
						tags = [ "terraform-test", "data_scaleway_vpcs", "basic" ]
					}`,
			},
			{
				Config: `
					resource scaleway_vpc vpc01 {
						name = "tf-vpc-datasource0"
						tags = [ "terraform-test", "data_scaleway_vpcs", "basic" ]
					}
				
					resource scaleway_vpc vpc02 {
						name = "tf-vpc-datasource1"
						tags = [ "terraform-test", "data_scaleway_vpcs", "basic" ]
					}
					data scaleway_vpcs vpcs_by_name {
						name = "tf-vpc-datasource"
					}
					
					data scaleway_vpcs vpcs_by_tag {
						tags = ["data_scaleway_vpcs", "terraform-test"]
					}

					data scaleway_vpcs vpcs_by_name_other_region {
						name = "tf-vpc-datasource"
						region = "nl-ams"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.scaleway_vpcs.vpcs_by_tag", "vpcs.0.id"),
					resource.TestCheckResourceAttrSet("data.scaleway_vpcs.vpcs_by_tag", "vpcs.1.id"),

					resource.TestCheckResourceAttrSet("data.scaleway_vpcs.vpcs_by_name", "vpcs.0.id"),
					resource.TestCheckResourceAttrSet("data.scaleway_vpcs.vpcs_by_name", "vpcs.1.id"),

					resource.TestCheckNoResourceAttr("data.scaleway_vpcs.vpcs_by_name_other_region", "vpcs.0.id"),
				),
			},
		},
	})
}
