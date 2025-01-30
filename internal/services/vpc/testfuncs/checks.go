package vpctestfuncs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	vpc2 "github.com/scaleway/scaleway-sdk-go/api/vpc/v2"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/vpc"
)

func CheckPrivateNetworkDestroy(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_vpc_private_network" {
				continue
			}

			vpcAPI, region, ID, err := vpc.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}
			_, err = vpcAPI.GetPrivateNetwork(&vpc2.GetPrivateNetworkRequest{
				PrivateNetworkID: ID,
				Region:           region,
			})

			if err == nil {
				return fmt.Errorf(
					"VPC private network %s still exists",
					rs.Primary.ID,
				)
			}
			// Unexpected api error we return it
			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}

func IsPrivateNetworkPresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		vpcAPI, region, ID, err := vpc.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = vpcAPI.GetPrivateNetwork(&vpc2.GetPrivateNetworkRequest{
			PrivateNetworkID: ID,
			Region:           region,
		})
		if err != nil {
			return err
		}

		return nil
	}
}
