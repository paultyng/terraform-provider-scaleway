package scaleway

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	baremetal "github.com/scaleway/scaleway-sdk-go/api/baremetal/v1alpha1"
)

func TestAccScalewayBaremetalServerMinimal1(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckScalewayBaremetalServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckScalewayBaremetalServerConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayBaremetalServerExists("scaleway_baremetal_server_beta.base"),
					resource.TestCheckResourceAttr("scaleway_baremetal_server_beta.base", "name", "namo"),
					resource.TestCheckResourceAttr("scaleway_baremetal_server_beta.base", "type", "9eebce52-f7d5-484f-9437-b234164c4c4b"),
					resource.TestCheckResourceAttr("scaleway_baremetal_server_beta.base", "image_id", "d17d6872-0412-45d9-a198-af82c34d3c5c"),
					resource.TestCheckResourceAttr("scaleway_baremetal_server_beta.base", "description", "test a description"),
					resource.TestCheckResourceAttr("scaleway_baremetal_server_beta.base", "tags.0", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_baremetal_server_beta.base", "tags.1", "scaleway_baremetal_server_beta"),
					resource.TestCheckResourceAttr("scaleway_baremetal_server_beta.base", "tags.2", "minimal"),
				),
			},
		},
	})
}

func testAccCheckScalewayBaremetalServerExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		baremetalAPI, zone, ID, err := getBaremetalAPIWithZoneAndID(testAccProvider.Meta(), rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = baremetalAPI.GetServer(&baremetal.GetServerRequest{ServerID: ID, Zone: zone})
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckScalewayBaremetalServerDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "scaleway_baremetal_server_beta" {
			continue
		}

		baremetalAPI, zone, ID, err := getBaremetalAPIWithZoneAndID(testAccProvider.Meta(), rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = baremetalAPI.GetServer(&baremetal.GetServerRequest{
			ServerID: ID,
			Zone:     zone,
		})

		// If no error resource still exist
		if err == nil {
			return fmt.Errorf("Server (%s) still exists", rs.Primary.ID)
		}

		// Unexpected api error we return it
		if !is404Error(err) {
			return err
		}
	}

	return nil
}

func testAccCheckScalewayBaremetalServerConfigMinimal() string {
	return `
resource "scaleway_baremetal_server_beta" "base" {
  name        = "namo"
  zone		  = "fr-par-2"
  description = "test a description"
  type        = "9eebce52-f7d5-484f-9437-b234164c4c4b"
  image_id    = "d17d6872-0412-45d9-a198-af82c34d3c5c"
  ssh_key_ids = ["2a9845b2-d722-4fc1-a77d-fe763728cc37"]

  tags = [ "terraform-test", "scaleway_baremetal_server_beta", "minimal" ]
}`
}
