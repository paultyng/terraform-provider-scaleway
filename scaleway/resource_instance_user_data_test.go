package scaleway_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func TestAccScalewayInstanceServerUserData_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayInstanceServerDestroy(tt),
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
