package scaleway

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccScalewayIP_Count(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckScalewayIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckScalewayIPConfig_Count,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayIPExists("scaleway_ip.base.0"),
					testAccCheckScalewayIPExists("scaleway_ip.base.1"),
				),
			},
		},
	})
}

func TestAccScalewayIP_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckScalewayIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckScalewayIPConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayIPExists("scaleway_ip.base"),
				),
			},
			{
				Config: testAccCheckScalewayIPConfig_Reverse,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayIPExists("scaleway_ip.base"),
					resource.TestCheckResourceAttr(
						"scaleway_ip.base", "reverse", "www.google.de"),
				),
			},
			{
				Config: testAccCheckScalewayIPAttachConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayServerExists("scaleway_server.base"),
					testAccCheckScalewayIPExists("scaleway_ip.base"),
					testAccCheckScalewayIPAttachment("scaleway_ip.base", func(serverID string) bool {
						return serverID != ""
					}, "attachment failed"),
					resource.TestCheckResourceAttr(
						"scaleway_ip.base", "reverse", ""),
				),
			},
			{
				Config: testAccCheckScalewayIPConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayIPExists("scaleway_ip.base"),
					testAccCheckScalewayIPAttachment("scaleway_ip.base", func(serverID string) bool {
						return serverID == ""
					}, "detachment failed"),
				),
			},
		},
	})
}

func testAccCheckScalewayIPDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Meta).deprecatedClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "scaleway" {
			continue
		}

		_, err := client.GetIP(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("IP still exists")
		}
	}

	return nil
}

func testAccCheckScalewayIPExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No IP ID is set")
		}

		client := testAccProvider.Meta().(*Meta).deprecatedClient
		ip, err := client.GetIP(rs.Primary.ID)

		if err != nil {
			return err
		}

		if ip.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		return nil
	}
}

func testAccCheckScalewayIPAttachment(n string, check func(string) bool, msg string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No IP ID is set")
		}

		client := testAccProvider.Meta().(*Meta).deprecatedClient
		ip, err := client.GetIP(rs.Primary.ID)

		if err != nil {
			return err
		}

		var serverID = ""
		if ip.Server != nil {
			serverID = ip.Server.Identifier
		}
		if !check(serverID) {
			return fmt.Errorf("IP check failed: %q", msg)
		}

		return nil
	}
}

var testAccCheckScalewayIPConfig_Reverse = `
resource "scaleway_ip" "base" {
  reverse = "www.google.de"
}
`

var testAccCheckScalewayIPConfig = `
resource "scaleway_ip" "base" {}
`

var testAccCheckScalewayIPConfig_Count = `
resource "scaleway_ip" "base" {
  count = 2
}
`

var testAccCheckScalewayIPAttachConfig = `
data "scaleway_image" "ubuntu" {
  architecture = "x86_64"
  name         = "Ubuntu Bionic"
  most_recent  = true
}

resource "scaleway_server" "base" {
  name = "test"

  image = "${data.scaleway_image.ubuntu.id}"
  type = "DEV1-S"
}

resource "scaleway_ip" "base" {
  server = "${scaleway_server.base.id}"
}
`
