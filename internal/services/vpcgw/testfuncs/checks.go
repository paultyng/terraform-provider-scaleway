package vpcgwtestfuncs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	vpcgwSDK "github.com/scaleway/scaleway-sdk-go/api/vpcgw/v1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/vpcgw"
)

func IsGatewayNetworkDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_vpc_gateway_network" {
				continue
			}

			vpcgwNetworkAPI, zone, ID, err := vpcgw.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = vpcgwNetworkAPI.GetGatewayNetwork(&vpcgwSDK.GetGatewayNetworkRequest{
				GatewayNetworkID: ID,
				Zone:             zone,
			})

			if err == nil {
				return fmt.Errorf(
					"VPC gateway network %s still exists",
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

func IsGatewayDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_vpc_public_gateway" {
				continue
			}

			vpcgwAPI, zone, ID, err := vpcgw.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = vpcgwAPI.GetGateway(&vpcgwSDK.GetGatewayRequest{
				GatewayID: ID,
				Zone:      zone,
			})

			if err == nil {
				return fmt.Errorf(
					"VPC public gateway %s still exists",
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

func IsDHCPDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_vpc_public_gateway_dhcp" {
				continue
			}

			vpcgwAPI, zone, ID, err := vpcgw.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = vpcgwAPI.GetDHCP(&vpcgwSDK.GetDHCPRequest{
				DHCPID: ID,
				Zone:   zone,
			})

			if err == nil {
				return fmt.Errorf(
					"VPC public gateway DHCP config %s still exists",
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

func IsIPDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_vpc_public_gateway_ip" {
				continue
			}

			vpcgwAPI, zone, ID, err := vpcgw.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = vpcgwAPI.GetIP(&vpcgwSDK.GetIPRequest{
				IPID: ID,
				Zone: zone,
			})

			if err == nil {
				return fmt.Errorf(
					"VPC public gateway ip %s still exists",
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
