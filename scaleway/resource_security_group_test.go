package scaleway

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccScalewaySecurityGroup_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckScalewaySecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckScalewaySecurityGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewaySecurityGroupExists("scaleway_security_group.base"),
					testAccCheckScalewaySecurityGroupAttributes("scaleway_security_group.base"),
					resource.TestCheckResourceAttr("scaleway_security_group.base", "name", "public"),
					resource.TestCheckResourceAttr("scaleway_security_group.base", "description", "public gateway"),
				),
			},
		},
	})
}

func TestAccScalewaySecurityGroup_Stateful(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckScalewaySecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckScalewaySecurityGroupConfig_Stateful,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewaySecurityGroupExists("scaleway_security_group.base"),
					resource.TestCheckResourceAttr("scaleway_security_group.base", "inbound_default_policy", "accept"),
					resource.TestCheckResourceAttr("scaleway_security_group.base", "outbound_default_policy", "drop"),
					resource.TestCheckResourceAttr("scaleway_security_group.base", "stateful", "true"),
				),
			},
			{
				Config: testAccCheckScalewaySecurityGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewaySecurityGroupExists("scaleway_security_group.base"),
					resource.TestCheckResourceAttr("scaleway_security_group.base", "stateful", "false"),
					resource.TestCheckResourceAttr("scaleway_security_group.base", "inbound_default_policy", "accept"),
					resource.TestCheckResourceAttr("scaleway_security_group.base", "outbound_default_policy", "accept"),
				),
			},
		},
	})
}

func testAccCheckScalewaySecurityGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Meta).deprecatedClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "scaleway" {
			continue
		}

		_, err := client.GetSecurityGroup(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Security Group still exists")
		}
	}

	return nil
}

func testAccCheckScalewaySecurityGroupAttributes(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Unknown resource: %s", n)
		}

		client := testAccProvider.Meta().(*Meta).deprecatedClient
		group, err := client.GetSecurityGroup(rs.Primary.ID)
		if err != nil {
			return err
		}

		if group.Name != "public" {
			return fmt.Errorf("Security Group has wrong name")
		}
		if group.Description != "public gateway" {
			return fmt.Errorf("Security Group has wrong description")
		}

		return nil
	}
}

func testAccCheckScalewaySecurityGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Security Group ID is set")
		}

		client := testAccProvider.Meta().(*Meta).deprecatedClient
		group, err := client.GetSecurityGroup(rs.Primary.ID)

		if err != nil {
			return err
		}

		if group.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		return nil
	}
}

var testAccCheckScalewaySecurityGroupConfig = `
resource "scaleway_security_group" "base" {
  name = "public"
  description = "public gateway"
}
`

var testAccCheckScalewaySecurityGroupConfig_Stateful = `
resource "scaleway_security_group" "base" {
  name = "public"
  description = "public gateway"
  stateful = true
  inbound_default_policy = "accept"
  outbound_default_policy = "drop"
}
`
