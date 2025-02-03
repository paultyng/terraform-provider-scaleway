package instance_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	instancechecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/instance/testfuncs"
)

func TestAccServerUserData_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
				resource "scaleway_instance_user_data" "main" {
					server_id = scaleway_instance_server.main.id
				   	key = "cloud-init"
					value = <<-EOF
#cloud-config
apt-update: true
apt-upgrade: true
EOF
				}

				resource "scaleway_instance_server" "main" {
					image = "ubuntu_focal"
					type  = "DEV1-S"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_instance_user_data.main", "key", "cloud-init"),
				),
			},
		},
	})
}
