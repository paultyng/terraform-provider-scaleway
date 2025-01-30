package baremetal_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/zonal"
	baremetalchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/baremetal/testfuncs"
)

func TestAccDataSourceOption_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      baremetalchecks.CheckServerDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					data "scaleway_baremetal_option" "by_name" {
						name = "Remote Access"
					}
					
					data "scaleway_baremetal_option" "by_id" {
						option_id = "931df052-d713-4674-8b58-96a63244c8e2"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isOptionPresent(tt, "data.scaleway_baremetal_option.by_id"),
					isOptionPresent(tt, "data.scaleway_baremetal_option.by_name"),

					resource.TestCheckResourceAttrSet("data.scaleway_baremetal_option.by_name", "name"),
					resource.TestCheckResourceAttrSet("data.scaleway_baremetal_option.by_name", "option_id"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_option.by_name", "manageable", "true"),

					resource.TestCheckResourceAttrSet("data.scaleway_baremetal_option.by_id", "name"),
					resource.TestCheckResourceAttrSet("data.scaleway_baremetal_option.by_id", "option_id"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_option.by_id", "manageable", "true"),
				),
			},
		},
	})
}

func isOptionPresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		zone, ID, err := zonal.ParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		api := baremetal.NewAPI(tt.Meta.ScwClient())
		_, err = api.GetOption(&baremetal.GetOptionRequest{
			OptionID: ID,
			Zone:     zone,
		})
		if err != nil {
			return err
		}

		return nil
	}
}
