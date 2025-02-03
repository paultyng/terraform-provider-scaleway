package vpcgw_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	vpcgwSDK "github.com/scaleway/scaleway-sdk-go/api/vpcgw/v1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/domain"
	instancechecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/instance/testfuncs"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/vpcgw"
)

func TestAccVPCPublicGatewayIPReverseDns_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	testDNSZone := "tf-reverse-vpcgw." + acctest.TestDomain
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsIPDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_vpc_public_gateway_ip" "main" {}
					
					resource "scaleway_domain_record" "tf_A" {
						dns_zone = %[1]q
						name     = ""
						type     = "A"
                        data     = "${scaleway_vpc_public_gateway_ip.main.address}"
						ttl      = 3600
						priority = 1
					}

					resource "scaleway_vpc_public_gateway_ip_reverse_dns" "main" {
					    gateway_ip_id   = scaleway_vpc_public_gateway_ip.main.id
					    reverse         = %[1]q
						depends_on      = [scaleway_domain_record.tf_A]
					}
				`, testDNSZone),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_ip_reverse_dns.main", "reverse", testDNSZone),
				),
			},
			{
				Config: `
					resource "scaleway_vpc_public_gateway_ip" "main" {}
				`,
				Check: testAccCheckVPCPublicGatewayIPDefaultReverse(tt, "scaleway_vpc_public_gateway_ip.main"),
			},
		},
	})
}

func testAccCheckVPCPublicGatewayIPDefaultReverse(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		api, zone, ID, err := vpcgw.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		ip, err := api.GetIP(&vpcgwSDK.GetIPRequest{
			IPID: ID,
			Zone: zone,
		})
		if err != nil {
			return err
		}

		if *ip.Reverse != domain.FindDefaultReverse(ip.Address.String()) {
			return fmt.Errorf("reverse should be the same, %v is different than %v", *ip.Reverse, ip.Address.String())
		}

		return nil
	}
}
