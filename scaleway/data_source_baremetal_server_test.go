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
	SSHKey := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIM7HUxRyQtB2rnlhQUcbDGCZcTJg7OvoznOiyC9W6IxH opensource@scaleway.com"
	name := "TestAccScalewayDataSourceBaremetalServer_Basic"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayBaremetalServerDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_account_ssh_key" "main" {
						name 	   = "%s"
						public_key = "%s"
					}
					
					resource "scaleway_baremetal_server" "main" {
						name        = "%s"
						zone        = "fr-par-2"
						description = "test a description"
						offer       = "GP-BM1-S"
						os          = "d17d6872-0412-45d9-a198-af82c34d3c5c"
					
						ssh_key_ids = [ scaleway_account_ssh_key.main.id ]
					}

					data "scaleway_baremetal_server" "by_name" {
						name = "${scaleway_baremetal_server.main.name}"
					}
					
					data "scaleway_baremetal_server" "by_id" {
						server_id = "${scaleway_baremetal_server.main.id}"
					}
				`, SSHKeyName, SSHKey, name),
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
