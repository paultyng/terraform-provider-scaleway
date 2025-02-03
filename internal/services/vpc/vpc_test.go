package vpc_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	vpcSDK "github.com/scaleway/scaleway-sdk-go/api/vpc/v2"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/vpc"
)

func TestAccVPC_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckVPCDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_vpc" "vpc01" {
					  name = "test-vpc"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCExists(tt, "scaleway_vpc.vpc01"),
					resource.TestCheckResourceAttr("scaleway_vpc.vpc01", "name", "test-vpc"),
					resource.TestCheckResourceAttr("scaleway_vpc.vpc01", "is_default", "false"),
					resource.TestCheckResourceAttrSet("scaleway_vpc.vpc01", "created_at"),
					resource.TestCheckResourceAttrSet("scaleway_vpc.vpc01", "updated_at"),
					resource.TestCheckResourceAttr("scaleway_vpc.vpc01", "enable_routing", "true"),
				),
			},
		},
	})
}

func TestAccVPC_WithRegion(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckVPCDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_vpc" "vpc01" {
					  name = "test-vpc"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCExists(tt, "scaleway_vpc.vpc01"),
					resource.TestCheckResourceAttr("scaleway_vpc.vpc01", "region", "fr-par"),
				),
			},
			{
				Config: `
					resource "scaleway_vpc" "vpc01" {
					  name   = "test-vpc"
					  region = "nl-ams"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCExists(tt, "scaleway_vpc.vpc01"),
					resource.TestCheckResourceAttr("scaleway_vpc.vpc01", "region", "nl-ams"),
				),
			},
		},
	})
}

func TestAccVPC_WithTags(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckVPCDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_vpc" "vpc01" {
					  name = "test-vpc"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCExists(tt, "scaleway_vpc.vpc01"),
					resource.TestCheckNoResourceAttr("scaleway_vpc.vpc01", "tags"),
				),
			},
			{
				Config: `
					resource "scaleway_vpc" "vpc01" {
					  name = "test-vpc"
					  tags = ["terraform-test", "vpc"]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCExists(tt, "scaleway_vpc.vpc01"),
					resource.TestCheckResourceAttr("scaleway_vpc.vpc01", "tags.#", "2"),
					resource.TestCheckResourceAttr("scaleway_vpc.vpc01", "tags.0", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_vpc.vpc01", "tags.1", "vpc"),
				),
			},
		},
	})
}

func TestAccVPC_DisableRouting(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckVPCDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_vpc" "vpc01" {
					  name = "test-vpc-disable-routing"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCExists(tt, "scaleway_vpc.vpc01"),
					resource.TestCheckResourceAttr("scaleway_vpc.vpc01", "enable_routing", "true"),
				),
			},
			{
				Config: `
					resource "scaleway_vpc" "vpc01" {
					  name           = "test-vpc-disable-routing"
					  enable_routing = false
					}
				`,
				ExpectError: regexp.MustCompile("routing cannot be disabled on this VPC"),
			},
		},
	})
}

func testAccCheckVPCExists(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		vpcAPI, region, ID, err := vpc.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = vpcAPI.GetVPC(&vpcSDK.GetVPCRequest{
			VpcID:  ID,
			Region: region,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckVPCDestroy(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_vpc" {
				continue
			}

			vpcAPI, region, ID, err := vpc.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = vpcAPI.GetVPC(&vpcSDK.GetVPCRequest{
				VpcID:  ID,
				Region: region,
			})

			if err == nil {
				return fmt.Errorf("VPC (%s) still exists", rs.Primary.ID)
			}

			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}
