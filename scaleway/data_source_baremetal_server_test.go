package scaleway

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccScalewayDataSourceBaremetalServer_Basic(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()

	SSHKeyName := "TestAccScalewayDataSourceBaremetalServer_Basic"
	name := "TestAccScalewayDataSourceBaremetalServer_Basic"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayBaremetalServerDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_account_ssh_key" "main" {
						name       = "%s"
						public_key = "%s"
					}
					
					resource "scaleway_baremetal_server" "main" {
						name        = "%s"
						zone        = "fr-par-2"
						description = "test a description"
						offer       = "EM-A210R-HDD"
						os          = "d17d6872-0412-45d9-a198-af82c34d3c5c"
					
						ssh_key_ids = [ scaleway_account_ssh_key.main.id ]
					}

					data "scaleway_baremetal_server" "by_name" {
						name = "${scaleway_baremetal_server.main.name}"
						zone = "fr-par-2"
					}
					
					data "scaleway_baremetal_server" "by_id" {
						server_id = "${scaleway_baremetal_server.main.id}"
						zone = "fr-par-2"
					}
				`, SSHKeyName, SSHKeyBaremetal, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayBaremetalServerExists(tt, "data.scaleway_baremetal_server.by_id"),
					testAccCheckScalewayBaremetalServerExists(tt, "data.scaleway_baremetal_server.by_name"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_baremetal_server.by_name", "name",
						"scaleway_baremetal_server.main", "name"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_baremetal_server.by_id", "name",
						"scaleway_baremetal_server.main", "name"),
				),
			},
		},
	})
}
