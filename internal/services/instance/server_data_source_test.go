package instance_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	instancechecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/instance/testfuncs"
)

func TestAccDataSourceServer_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	serverName := "tf-server"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_instance_server" "main" {
						name 	= "%s"
						image = "ubuntu_focal"
						type  = "DEV1-S"
						state = "stopped"
						tags  = [ "terraform-test", "data_scaleway_instance_server", "basic" ]
					}`, serverName),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_instance_server" "main" {
						name 	= "%s"
						image = "ubuntu_focal"
						type  = "DEV1-S"
						state = "stopped"
						tags  = [ "terraform-test", "data_scaleway_instance_server", "basic" ]
					}
					
					data "scaleway_instance_server" "prod" {
						name = "${scaleway_instance_server.main.name}"
					}
					
					data "scaleway_instance_server" "stg" {
						server_id = "${scaleway_instance_server.main.id}"
					}`, serverName),
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "data.scaleway_instance_server.prod"),
					resource.TestCheckResourceAttr("data.scaleway_instance_server.prod", "name", serverName),
					isServerPresent(tt, "data.scaleway_instance_server.stg"),
					resource.TestCheckResourceAttr("data.scaleway_instance_server.stg", "name", serverName),
				),
			},
		},
	})
}
