package instance_test

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	instanceSDK "github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/zonal"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/instance"
	instancechecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/instance/testfuncs"
)

func TestAccServer_Minimal1(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_server" "base" {
					  image = "ubuntu_focal"
					  type  = "DEV1-S"
					
					  tags = [ "terraform-test", "scaleway_instance_server", "minimal" ]
					}`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.base"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "image", "ubuntu_focal"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "type", "DEV1-S"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "root_volume.0.delete_on_termination", "true"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "root_volume.0.size_in_gb", "20"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.base", "root_volume.0.volume_id"),
					serverHasNewVolume(tt, "scaleway_instance_server.base"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "enable_dynamic_ip", "false"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "tags.0", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "tags.1", "scaleway_instance_server"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "tags.2", "minimal"),
				),
			},
			{
				// Image label such as ubuntu_focal
				Config: `
					resource "scaleway_instance_server" "base" {
					  image = "ubuntu_focal"
					  type  = "DEV1-S"
					
					  tags = [ "terraform-test", "scaleway_instance_server", "minimal" ]
					}`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.base"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "image", "ubuntu_focal"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "type", "DEV1-S"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "root_volume.0.delete_on_termination", "true"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "root_volume.0.size_in_gb", "20"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.base", "root_volume.0.volume_id"),
					serverHasNewVolume(tt, "scaleway_instance_server.base"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "tags.0", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "tags.1", "scaleway_instance_server"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "tags.2", "minimal"),
				),
			},
		},
	})
}

func TestAccServer_Minimal2(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_server" "base" {
					  image = "ubuntu_focal"
					  type  = "DEV1-S"
					}`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.base"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "image", "ubuntu_focal"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "type", "DEV1-S"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "root_volume.0.delete_on_termination", "true"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "root_volume.0.size_in_gb", "20"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.base", "root_volume.0.volume_id"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "enable_dynamic_ip", "false"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "state", "started"),
				),
			},
			{
				Config: `
					resource "scaleway_instance_server" main {
					  image = "ubuntu_focal"
					  type  = "DEV1-S"
					  root_volume {
						volume_type = "l_ssd"
					  }
					}`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.main"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "image", "ubuntu_focal"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "type", "DEV1-S"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "root_volume.0.delete_on_termination", "true"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "root_volume.0.size_in_gb", "20"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "root_volume.0.volume_type", "l_ssd"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.main", "root_volume.0.volume_id"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "enable_dynamic_ip", "false"),
				),
			},
			{
				Config: `
					resource "scaleway_instance_server" main2 {
					  image = "ubuntu_focal"
					  type  = "DEV1-S"
					  root_volume {
						volume_type = "l_ssd"
						size_in_gb  = 20
					  }
					}`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.main2"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main2", "image", "ubuntu_focal"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main2", "type", "DEV1-S"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main2", "root_volume.0.delete_on_termination", "true"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main2", "root_volume.0.size_in_gb", "20"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main2", "root_volume.0.volume_type", "l_ssd"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.main2", "root_volume.0.volume_id"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main2", "enable_dynamic_ip", "false"),
				),
			},
		},
	})
}

func TestAccServer_RootVolume1(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	serverID := ""

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_server" "base" {
						image = "ubuntu_focal"
						type  = "DEV1-S"
						root_volume {
							size_in_gb = 10
							delete_on_termination = true
						}
						tags = [ "terraform-test", "scaleway_instance_server", "root_volume" ]
					}`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.base"),
					serverHasNewVolume(tt, "scaleway_instance_server.base"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "root_volume.0.size_in_gb", "10"),
					acctest.CheckResourceIDPersisted("scaleway_instance_server.base", &serverID),
				),
			},
			{
				Config: `
					resource "scaleway_instance_server" "base" {
						image = "ubuntu_focal"
						type  = "DEV1-S"
						root_volume {
							size_in_gb = 20
							delete_on_termination = true
						}
						tags = [ "terraform-test", "scaleway_instance_server", "root_volume" ]
					}`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.base"),
					serverHasNewVolume(tt, "scaleway_instance_server.base"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "root_volume.0.size_in_gb", "20"),
					acctest.CheckResourceIDChanged("scaleway_instance_server.base", &serverID), // Server should have been re-created as l_ssd cannot be resized.
				),
			},
		},
	})
}

func TestAccServer_RootVolume_Boot(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_server" "base" {
						image = "ubuntu_focal"
						type  = "DEV1-S"
						state = "stopped"
						root_volume {
							boot = true
							delete_on_termination = true
						}
						tags = [ "terraform-test", "scaleway_instance_server", "root_volume" ]
					}`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.base"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "root_volume.0.boot", "true"),
					serverHasNewVolume(tt, "scaleway_instance_server.base"),
				),
			},
			{
				Config: `
					resource "scaleway_instance_server" "base" {
						image = "ubuntu_focal"
						type  = "DEV1-S"
						state = "stopped"
						root_volume {
							boot = false
							delete_on_termination = true
						}
						tags = [ "terraform-test", "scaleway_instance_server", "root_volume" ]
					}`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.base"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "root_volume.0.boot", "false"),
					serverHasNewVolume(tt, "scaleway_instance_server.base"),
				),
			},
		},
	})
}

func TestAccServer_RootVolume_ID(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_volume" "server_volume" {
					  type       = "b_ssd"
					  name       = "tf_tests_rootvolume"
					  size_in_gb = 10
					}

					resource "scaleway_instance_server" "base" {
						type  = "DEV1-S"
						state = "stopped"
						root_volume {
							volume_id = scaleway_instance_volume.server_volume.id
							boot = true
							delete_on_termination = false
						}
						tags = [ "terraform-test", "scaleway_instance_server", "root_volume" ]
					}`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.base"),
					resource.TestCheckResourceAttrPair("scaleway_instance_server.base", "root_volume.0.volume_id", "scaleway_instance_volume.server_volume", "id"),
					resource.TestCheckResourceAttrPair("scaleway_instance_server.base", "root_volume.0.size_in_gb", "scaleway_instance_volume.server_volume", "size_in_gb"),
				),
			},
		},
	})
}

func TestAccServer_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				// DEV1-M
				Config: `
					data "scaleway_marketplace_image" "ubuntu" {
					  instance_type   = "DEV1-M"
					  label         = "ubuntu_focal"
					}
					
					resource "scaleway_instance_server" "base" {
					  name  = "test"
					  image = "${data.scaleway_marketplace_image.ubuntu.id}"
					  type  = "DEV1-M"
					
					  tags = [ "terraform-test", "scaleway_instance_server", "basic" ]
					}`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.base"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "type", "DEV1-M"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "name", "test"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "tags.0", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "tags.1", "scaleway_instance_server"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "tags.2", "basic"),
				),
			},
			{
				// DEV1-S
				Config: `
					data "scaleway_marketplace_image" "ubuntu" {
					  instance_type   = "DEV1-S"
					  label         = "ubuntu_focal"
					}
					
					resource "scaleway_instance_server" "base" {
					  name  = "test"
					  image = "${data.scaleway_marketplace_image.ubuntu.id}"
					  type  = "DEV1-S"
					  replace_on_type_change  = true
					
					  tags = [ "terraform-test", "scaleway_instance_server", "basic" ]
					}`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.base"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "type", "DEV1-S"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "name", "test"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "tags.0", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "tags.1", "scaleway_instance_server"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "tags.2", "basic"),
				),
			},
		},
	})
}

func TestAccServer_State1(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				// started
				Config: `
					data "scaleway_marketplace_image" "ubuntu" {
					  instance_type = "DEV1-S"
					  label         = "ubuntu_focal"
					}
					
					resource "scaleway_instance_server" "base" {
					  image = "${data.scaleway_marketplace_image.ubuntu.id}"
					  type  = "DEV1-S"
					  state = "started"
					  tags  = [ "terraform-test", "scaleway_instance_server", "state" ]
					}`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.base"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "state", "started"),
				),
			},
			{
				// standby
				Config: `
					data "scaleway_marketplace_image" "ubuntu" {
					  instance_type = "DEV1-S"
					  label         = "ubuntu_focal"
					}
					
					resource "scaleway_instance_server" "base" {
					  image = "${data.scaleway_marketplace_image.ubuntu.id}"
					  type  = "DEV1-S"
					  state = "standby"
					  tags  = [ "terraform-test", "scaleway_instance_server", "state" ]
					}`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.base"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "state", "standby"),
				),
			},
			{
				// stopped
				Config: `
					data "scaleway_marketplace_image" "ubuntu" {
					  instance_type = "DEV1-S"
					  label         = "ubuntu_focal"
					}
					
					resource "scaleway_instance_server" "base" {
					  image = "${data.scaleway_marketplace_image.ubuntu.id}"
					  type  = "DEV1-S"
					  state = "stopped"
					  tags  = [ "terraform-test", "scaleway_instance_server", "state" ]
					}`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.base"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "state", "stopped"),
				),
			},
		},
	})
}

func TestAccServer_State2(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				// stopped
				Config: `
					data "scaleway_marketplace_image" "ubuntu" {
					  instance_type = "DEV1-S"
					  label         = "ubuntu_focal"
					}
					
					resource "scaleway_instance_server" "base" {
					  image = "${data.scaleway_marketplace_image.ubuntu.id}"
					  type  = "DEV1-S"
					  state = "stopped"
					  tags  = [ "terraform-test", "scaleway_instance_server", "state" ]
					}`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.base"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "state", "stopped"),
				),
			},
			{
				// standby
				Config: `
					data "scaleway_marketplace_image" "ubuntu" {
					  instance_type = "DEV1-S"
					  label         = "ubuntu_focal"
					}
					
					resource "scaleway_instance_server" "base" {
					  image = "${data.scaleway_marketplace_image.ubuntu.id}"
					  type  = "DEV1-S"
					  state = "standby"
					  tags  = [ "terraform-test", "scaleway_instance_server", "state" ]
					}`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.base"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "state", "standby"),
				),
			},
		},
	})
}

func TestAccServer_UserData_WithCloudInitAtStart(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
				resource "scaleway_instance_server" "base" {
					image = "ubuntu_focal"
					type  = "DEV1-S"

					user_data = {
				   		foo   = "bar"
						cloud-init =  <<EOF
#cloud-config
apt_update: true
apt_upgrade: true
EOF 
				 	}
				}`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.base"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "user_data.foo", "bar"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "user_data.cloud-init", "#cloud-config\napt_update: true\napt_upgrade: true\n"),
				),
			},
		},
	})
}

func TestAccServer_UserData_WithoutCloudInitAtStart(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				// Without cloud-init
				Config: `
					resource "scaleway_instance_server" "base" {
						image = "ubuntu_focal"
						type  = "DEV1-S"
						root_volume {
							size_in_gb = 20
						}
						tags  = [ "terraform-test", "scaleway_instance_server", "user_data" ]
					}`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.base"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "user_data.%", "0"),
				),
			},
			{
				// With cloud-init
				Config: `
					resource "scaleway_instance_server" "base" {
						image = "ubuntu_focal"
						type  = "DEV1-S"
						tags  = [ "terraform-test", "scaleway_instance_server", "user_data" ]
						root_volume {
							size_in_gb = 20
						}
						user_data = {
							cloud-init = <<EOF
#cloud-config
apt_update: true
apt_upgrade: true
EOF
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.base"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "user_data.cloud-init", "#cloud-config\napt_update: true\napt_upgrade: true\n"),
				),
			},
		},
	})
}

func TestAccServer_AdditionalVolumes(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				// With additional local
				Config: `
					resource "scaleway_instance_volume" "local" {
						size_in_gb = 10
						type = "l_ssd"
					}

					resource "scaleway_instance_server" "base" {
						image = "ubuntu_focal"
						type = "DEV1-S"
						
						root_volume {
							size_in_gb = 10
						}

						tags = [ "terraform-test", "scaleway_instance_server", "additional_volume_ids" ]

						additional_volume_ids = [
							scaleway_instance_volume.local.id
						]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.base"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "root_volume.0.size_in_gb", "10"),
				),
			},
			{
				// With additional local and block
				Config: `
					resource "scaleway_instance_volume" "local" {
						size_in_gb = 10
						type = "l_ssd"
					}

					resource "scaleway_instance_volume" "block" {
						size_in_gb = 10
						type = "b_ssd"
					}

					resource "scaleway_instance_server" "base" {
						image = "ubuntu_focal"
						type = "DEV1-S"
						
						root_volume {
							size_in_gb = 10
						}

						tags = [ "terraform-test", "scaleway_instance_server", "additional_volume_ids" ]

						additional_volume_ids = [
							scaleway_instance_volume.local.id,
							scaleway_instance_volume.block.id
						]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isVolumePresent(tt, "scaleway_instance_volume.block"),
					isServerPresent(tt, "scaleway_instance_server.base"),
					resource.TestCheckResourceAttr("scaleway_instance_volume.block", "size_in_gb", "10"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "root_volume.0.size_in_gb", "10"),
				),
			},
		},
	})
}

func TestAccServer_AdditionalVolumesDetach(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			isVolumeDestroyed(tt),
			instancechecks.IsServerDestroyed(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: `
					variable "zone" {
						type    = string
						default = "fr-par-1"
					}

					resource "scaleway_instance_volume" "main" {
  						type       = "b_ssd"
  						name       = "foobar"
  						size_in_gb = 1
					}

					resource "scaleway_instance_server" "main" {
						type  = "DEV1-S"
  						image = "ubuntu_focal"
  						name  = "foobar"

						enable_dynamic_ip = true

						additional_volume_ids = [scaleway_instance_volume.main.id]
					}
				`,
			},
			{
				Config: `
					variable "zone" {
						type    = string
						default = "fr-par-1"
					}

					resource "scaleway_instance_volume" "main" {
  						type       = "b_ssd"
  						name       = "foobar"
  						size_in_gb = 1
					}

					resource "scaleway_instance_server" "main" {
						type  = "DEV1-S"
  						image = "ubuntu_focal"
  						name  = "foobar"

						enable_dynamic_ip = true

						additional_volume_ids = []
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isVolumePresent(tt, "scaleway_instance_volume.main"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "additional_volume_ids.#", "0"),
				),
			},
		},
	})
}

func TestAccServer_WithPlacementGroup(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_placement_group" "ha" {
						policy_mode = "enforced"
						policy_type = "max_availability"
					}
					
					resource "scaleway_instance_server" "base" {
						count = 3
						image = "ubuntu_focal"
						type  = "DEV1-S"
						placement_group_id = "${scaleway_instance_placement_group.ha.id}"
						tags  = [ "terraform-test", "scaleway_instance_server", "placement_group" ]
					}`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.base.0"),
					isServerPresent(tt, "scaleway_instance_server.base.1"),
					isServerPresent(tt, "scaleway_instance_server.base.2"),
					isPlacementGroupPresent(tt, "scaleway_instance_placement_group.ha"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base.0", "placement_group_policy_respected", "true"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base.1", "placement_group_policy_respected", "true"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base.2", "placement_group_policy_respected", "true"),
				),
			},
		},
	})
}

func TestAccServer_Ipv6(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_ip" "ip" {
						type = "routed_ipv6"
					}

					resource "scaleway_instance_server" "server01" {
						image = "ubuntu_focal"
		  				type  = "PLAY2-PICO"
						ip_ids = [scaleway_instance_ip.ip.id]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.server01"),
					acctest.CheckResourceAttrIPv6("scaleway_instance_server.server01", "public_ips.0.address"),
				),
			},
			{
				Config: `
					resource "scaleway_instance_server" "server01" {
						image = "ubuntu_focal"
		  				type  = "PLAY2-PICO"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.server01"),
					resource.TestCheckResourceAttr("scaleway_instance_server.server01", "ipv6_address", ""),
					resource.TestCheckResourceAttr("scaleway_instance_server.server01", "ipv6_gateway", ""),
					resource.TestCheckResourceAttr("scaleway_instance_server.server01", "ipv6_prefix_length", "0"),
					resource.TestCheckResourceAttr("scaleway_instance_server.server01", "public_ips.#", "0"),
				),
			},
		},
	})
}

func TestAccServer_Basic2(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					data "scaleway_marketplace_image" "ubuntu" {
					  instance_type   = "DEV1-M"
					  label         = "ubuntu_focal"
					}

					resource "scaleway_instance_server" "server01" {
						type  = "DEV1-S"
						image = data.scaleway_marketplace_image.ubuntu.id
						state = "stopped"
					}
				`,
			},
			{
				ResourceName: "scaleway_instance_server.server01",
				ImportState:  true,
			},
		},
	})
}

func TestAccServer_WithReservedIP(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_ip" "first" {}
					resource "scaleway_instance_server" "base" {
						image = "ubuntu_focal"
						type  = "DEV1-S"
						ip_id = scaleway_instance_ip.first.id
						tags  = [ "terraform-test", "scaleway_instance_server", "reserved_ip" ]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.base"),
					resource.TestCheckResourceAttrPair("scaleway_instance_ip.first", "address", "scaleway_instance_server.base", "public_ip"), // public_ip is deprecated
					resource.TestCheckResourceAttrPair("scaleway_instance_ip.first", "id", "scaleway_instance_server.base", "ip_id"),
				),
			},
			{
				Config: `
					resource "scaleway_instance_ip" "first" {}
					resource "scaleway_instance_ip" "second" {}
					resource "scaleway_instance_server" "base" {
						image = "ubuntu_focal"
						type  = "DEV1-S"
						ip_id = scaleway_instance_ip.second.id
						tags  = [ "terraform-test", "scaleway_instance_server", "reserved_ip" ]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.base"),
					isIPAttachedToServer(tt, "scaleway_instance_ip.second", "scaleway_instance_server.base"),
					resource.TestCheckResourceAttrPair("scaleway_instance_ip.second", "address", "scaleway_instance_server.base", "public_ip"),
					resource.TestCheckResourceAttrPair("scaleway_instance_ip.second", "id", "scaleway_instance_server.base", "ip_id"),
				),
			},
			{
				Config: `
					resource "scaleway_instance_ip" "first" {}
					resource "scaleway_instance_ip" "second" {}
					resource "scaleway_instance_server" "base" {
						image = "ubuntu_focal"
						type  = "DEV1-S"
						tags  = [ "terraform-test", "scaleway_instance_server", "reserved_ip" ]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.base"),
					serverHasNoIPAssigned(tt, "scaleway_instance_server.base"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "public_ip", ""),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "ip_id", ""),
				),
			},
			{
				Config: `
					resource "scaleway_instance_ip" "first" {}
					resource "scaleway_instance_ip" "second" {}
					resource "scaleway_instance_server" "base" {
						image = "ubuntu_focal"
						type  = "DEV1-S"
						enable_dynamic_ip = true
						tags  = [ "terraform-test", "scaleway_instance_server", "reserved_ip" ]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.base"),
					serverHasNoIPAssigned(tt, "scaleway_instance_server.base"),
					acctest.CheckResourceAttrIPv4("scaleway_instance_server.base", "public_ip"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "ip_id", ""),
				),
			},
		},
	})
}

func isServerPresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		instanceAPI, zone, ID, err := instance.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = instanceAPI.GetServer(&instanceSDK.GetServerRequest{ServerID: ID, Zone: zone})
		if err != nil {
			return err
		}

		return nil
	}
}

func arePrivateNICsPresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		instanceAPI, zone, ID, err := instance.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		res, err := instanceAPI.ListPrivateNICs(&instanceSDK.ListPrivateNICsRequest{ServerID: ID, Zone: zone})
		if err != nil {
			return err
		}

		privateNetworksOnServer := make(map[string]struct{})
		// build current private networks on server
		for _, key := range res.PrivateNics {
			privateNetworksOnServer[key.PrivateNetworkID] = struct{}{}
		}

		privateNetworksToCheckOnSchema := make(map[string]struct{})
		// build terraform private networks
		for key, value := range rs.Primary.Attributes {
			if strings.Contains(key, "pn_id") {
				privateNetworksToCheckOnSchema[locality.ExpandID(value)] = struct{}{}
			}
		}

		// check if private networks are present on server
		for pnKey := range privateNetworksToCheckOnSchema {
			if _, exist := privateNetworksOnServer[pnKey]; !exist {
				return errors.New("private network does not exist")
			}
		}

		return nil
	}
}

// serverHasNewVolume tests if volume name is generated by terraform
// It is useful as volume should not be set in request when creating an instanceSDK from an image
func serverHasNewVolume(_ *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		rootVolumeName, ok := rs.Primary.Attributes["root_volume.0.name"]
		if !ok {
			return errors.New("instanceSDK root_volume has no name")
		}

		if strings.HasPrefix(rootVolumeName, "tf") {
			return fmt.Errorf("root volume name is generated by provider, should be generated by api (%s)", rootVolumeName)
		}

		return nil
	}
}

func TestAccServer_AlterTags(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_server" "base" {
						type  = "DEV1-L"
						image = "ubuntu_focal"
						state = "stopped"
						tags = [ "front", "web" ]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.base"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "tags.0", "front"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "tags.1", "web"),
				),
			},
			{
				Config: `
					resource "scaleway_instance_server" "base" {
						type  = "DEV1-L"
						state = "stopped"
						image = "ubuntu_focal"
						tags = [ "front" ]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.base"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "tags.0", "front"),
				),
			},
		},
	})
}

func TestAccServer_WithDefaultRootVolumeAndAdditionalVolume(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_volume" "data" {
						size_in_gb = 100
						type = "b_ssd"
					}

					resource "scaleway_instance_server" "main" {
						type = "DEV1-S"
						image = "ubuntu-bionic"
						root_volume {
							delete_on_termination = false
					  	}
						additional_volume_ids = [ scaleway_instance_volume.data.id ]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.main"),
				),
			},
		},
	})
}

func TestAccServer_ServerWithBlockNonDefaultZone(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_volume" "main" {
						type       = "b_ssd"
						name       = "main"
						size_in_gb = 1
						zone       = "nl-ams-1"
					}

					resource "scaleway_instance_server" "main" {
						zone              = "nl-ams-1"
						image             = "ubuntu_focal"
						type              = "DEV1-S"
						root_volume {
							delete_on_termination = true
							size_in_gb            = 20
						}
						additional_volume_ids = [scaleway_instance_volume.main.id]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.main"),
				),
			},
		},
	})
}

func TestAccServer_PrivateNetwork(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_vpc_private_network internal {
						name = "private_network_instance"
					}

					resource "scaleway_instance_server" "base" {
						image = "ubuntu_focal"
						type  = "DEV1-S"
						zone = "fr-par-2"

						private_network {
							pn_id = scaleway_vpc_private_network.internal.id
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					arePrivateNICsPresent(tt, "scaleway_instance_server.base"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "private_network.#", "1"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.base", "private_network.0.pn_id"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.base", "private_network.0.mac_address"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.base", "private_network.0.status"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.base", "private_network.0.zone"),
					resource.TestCheckResourceAttrPair("scaleway_instance_server.base", "private_network.0.pn_id",
						"scaleway_vpc_private_network.internal", "id"),
				),
			},
			{
				Config: `
					resource scaleway_vpc_private_network pn01 {
						name = "private_network_instance"
					}

					resource "scaleway_instance_server" "base" {
						image = "ubuntu_focal"
						type  = "DEV1-S"
						zone  = "fr-par-1"

						private_network {
							pn_id = scaleway_vpc_private_network.pn01.id
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					arePrivateNICsPresent(tt, "scaleway_instance_server.base"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "private_network.#", "1"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.base", "private_network.0.pn_id"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.base", "private_network.0.mac_address"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.base", "private_network.0.status"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.base", "private_network.0.zone"),
					resource.TestCheckResourceAttrPair("scaleway_instance_server.base", "private_network.0.pn_id",
						"scaleway_vpc_private_network.pn01", "id"),
				),
			},
			{
				Config: `
					resource scaleway_vpc_private_network pn01 {
						name = "private_network_instance"
					}

					resource scaleway_vpc_private_network pn02 {
						name = "private_network_instance_02"
					}

					resource "scaleway_instance_server" "base" {
						image = "ubuntu_focal"
						type  = "DEV1-S"

						private_network {
							pn_id = scaleway_vpc_private_network.pn02.id
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					arePrivateNICsPresent(tt, "scaleway_instance_server.base"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "private_network.#", "1"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.base", "private_network.0.pn_id"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.base", "private_network.0.mac_address"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.base", "private_network.0.status"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.base", "private_network.0.zone"),
					resource.TestCheckResourceAttrPair("scaleway_instance_server.base", "private_network.0.pn_id",
						"scaleway_vpc_private_network.pn02", "id"),
				),
			},
			{
				Config: `
					resource scaleway_vpc_private_network pn01 {
						name = "private_network_instance"
					}

					resource scaleway_vpc_private_network pn02 {
						name = "private_network_instance_02"
					}

					resource "scaleway_instance_server" "base" {
						image = "ubuntu_focal"
						type  = "DEV1-S"

						private_network {
							pn_id = scaleway_vpc_private_network.pn02.id
						}

						private_network {
							pn_id = scaleway_vpc_private_network.pn01.id
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					arePrivateNICsPresent(tt, "scaleway_instance_server.base"),
					resource.TestCheckResourceAttrPair("scaleway_instance_server.base",
						"private_network.0.pn_id",
						"scaleway_vpc_private_network.pn02", "id"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "private_network.#", "2"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.base", "private_network.0.pn_id"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.base", "private_network.0.mac_address"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.base", "private_network.0.status"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.base", "private_network.0.zone"),
					resource.TestCheckResourceAttrPair("scaleway_instance_server.base", "private_network.1.pn_id",
						"scaleway_vpc_private_network.pn01", "id"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.base", "private_network.1.pn_id"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.base", "private_network.1.mac_address"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.base", "private_network.1.status"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.base", "private_network.1.zone"),
				),
			},
			{
				Config: `
					resource scaleway_vpc_private_network pn01 {
						name = "private_network_instance"
					}

					resource scaleway_vpc_private_network pn02 {
						name = "private_network_instance_02"
					}

					resource "scaleway_instance_server" "base" {
						image = "ubuntu_focal"
						type  = "DEV1-S"	
					}`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.base"),
					resource.TestCheckResourceAttr("scaleway_instance_server.base", "private_network.#", "0"),
				),
			},
		},
	})
}

func TestAccServer_Migrate(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_server" "main" {
						image = "ubuntu_jammy"
						type  = "PRO2-XXS"
						
					}`,
				Check: resource.ComposeTestCheckFunc(
					arePrivateNICsPresent(tt, "scaleway_instance_server.main"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "type", "PRO2-XXS"),
				),
			},
			{
				Config: `
					resource "scaleway_instance_server" "main" {
						image = "ubuntu_jammy"
						type  = "PRO2-XS"
						
					}`,
				Check: resource.ComposeTestCheckFunc(
					arePrivateNICsPresent(tt, "scaleway_instance_server.main"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "type", "PRO2-XS"),
				),
			},
			{
				Config: `
					resource "scaleway_instance_server" "main" {
						image = "ubuntu_jammy"
						type  = "PRO2-XXS"
						
					}`,
				Check: resource.ComposeTestCheckFunc(
					arePrivateNICsPresent(tt, "scaleway_instance_server.main"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "type", "PRO2-XXS"),
				),
			},
		},
	})
}

func TestAccServer_MigrateInvalidLocalVolumeSize(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_server" "main" {
						image = "ubuntu_jammy"
						type  = "DEV1-L"
						
					}`,
				Check: resource.ComposeTestCheckFunc(
					arePrivateNICsPresent(tt, "scaleway_instance_server.main"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "type", "DEV1-L"),
				),
			},
			{
				Config: `
					resource "scaleway_instance_server" "main" {
						image = "ubuntu_jammy"
						type  = "DEV1-S"
						
					}`,
				Check: resource.ComposeTestCheckFunc(
					arePrivateNICsPresent(tt, "scaleway_instance_server.main"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "type", "DEV1-S"),
				),
				ExpectError: regexp.MustCompile("cannot change server type"),
				PlanOnly:    true,
			},
		},
	})
}

func TestAccServer_CustomDiffImage(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_server" "main" {
						image = "ubuntu_jammy"
						type = "DEV1-S"
						state = "stopped"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.main"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "image", "ubuntu_jammy"),
				),
			},
			{
				Config: `
					resource "scaleway_instance_server" "main" {
						image = "ubuntu_jammy"
						type = "DEV1-S"
						state = "stopped"
					}
					resource "scaleway_instance_server" "copy" {
						image = "ubuntu_jammy"
						type = "DEV1-S"
						state = "stopped"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.main"),
					isServerPresent(tt, "scaleway_instance_server.copy"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "image", "ubuntu_jammy"),
					resource.TestCheckResourceAttr("scaleway_instance_server.copy", "image", "ubuntu_jammy"),
					resource.TestCheckResourceAttrPair("scaleway_instance_server.main", "id", "scaleway_instance_server.copy", "id"),
				),
				ResourceName: "scaleway_instance_server.copy",
				ImportState:  true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					return state.RootModule().Resources["scaleway_instance_server.main"].Primary.ID, nil
				},
				ImportStatePersist: true,
			},
			{
				Config: `
					data "scaleway_marketplace_image" "jammy" {
						label = "ubuntu_jammy"
					}
					resource "scaleway_instance_server" "main" {
						image = data.scaleway_marketplace_image.jammy.id
						type = "DEV1-S"
						state = "stopped"
					}
					resource "scaleway_instance_server" "copy" {
						image = "ubuntu_jammy"
						type = "DEV1-S"
						state = "stopped"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.main"),
					isServerPresent(tt, "scaleway_instance_server.copy"),
					resource.TestCheckResourceAttrPair("scaleway_instance_server.main", "image", "data.scaleway_marketplace_image.jammy", "id"),
					resource.TestCheckResourceAttrPair("scaleway_instance_server.main", "id", "scaleway_instance_server.copy", "id"),
				),
			},
			{
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.main"),
					isServerPresent(tt, "scaleway_instance_server.copy"),
					resource.TestCheckResourceAttrPair("scaleway_instance_server.main", "image", "data.scaleway_marketplace_image.jammy", "id"),
					resource.TestCheckResourceAttrPair("scaleway_instance_server.main", "id", "scaleway_instance_server.copy", "id"),
				),
			},
			{
				Config: `
					data "scaleway_marketplace_image" "focal" {
						label = "ubuntu_focal"
					}
					resource "scaleway_instance_server" "main" {
						image = data.scaleway_marketplace_image.focal.id
						type = "DEV1-S"
						state = "stopped"
					}
					resource "scaleway_instance_server" "copy" {
						image = "ubuntu_jammy"
						type = "DEV1-S"
						state = "stopped"
					}
				`,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_instance_server.main"),
					resource.TestCheckResourceAttrPair("scaleway_instance_server.main", "image", "data.scaleway_marketplace_image.focal", "id"),
					resource.TestCheckResourceAttr("scaleway_instance_server.copy", "image", "ubuntu_jammy"),
					serverIDsAreDifferent("scaleway_instance_server.main", "scaleway_instance_server.copy"),
				),
			},
		},
	})
}

func serverIDsAreDifferent(nameFirst, nameSecond string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[nameFirst]
		if !ok {
			return fmt.Errorf("resource was not found: %s", nameFirst)
		}
		idFirst := rs.Primary.ID

		rs, ok = s.RootModule().Resources[nameSecond]
		if !ok {
			return fmt.Errorf("resource was not found: %s", nameSecond)
		}
		idSecond := rs.Primary.ID

		if idFirst == idSecond {
			return fmt.Errorf("IDs of both resources were equal when they should not have been (%s and %s)", nameFirst, nameSecond)
		}
		return nil
	}
}

func TestAccServer_IPs(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_ip" "ip1" {
						type = "routed_ipv4"
					}

					resource "scaleway_instance_server" "main" {
						name = "tf-tests-instance-server-ips"
						ip_ids = [scaleway_instance_ip.ip1.id]
						image = "ubuntu_jammy"
						type  = "PRO2-XXS"
						state = "stopped"
					}`,
				Check: resource.ComposeTestCheckFunc(
					arePrivateNICsPresent(tt, "scaleway_instance_server.main"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "routed_ip_enabled", "true"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "public_ips.#", "1"),
					resource.TestCheckResourceAttrPair("scaleway_instance_server.main", "public_ips.0.id", "scaleway_instance_ip.ip1", "id"),
				),
			},
			{
				Config: `
					resource "scaleway_instance_ip" "ip1" {
						type = "routed_ipv4"
					}

					resource "scaleway_instance_ip" "ip2" {
						type = "routed_ipv4"
					}

					resource "scaleway_instance_server" "main" {
						name = "tf-tests-instance-server-ips"
						ip_ids = [scaleway_instance_ip.ip1.id, scaleway_instance_ip.ip2.id]
						image = "ubuntu_jammy"
						type  = "PRO2-XXS"
						state = "stopped"
					}`,
				Check: resource.ComposeTestCheckFunc(
					arePrivateNICsPresent(tt, "scaleway_instance_server.main"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "routed_ip_enabled", "true"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "public_ips.#", "2"),
					resource.TestCheckResourceAttrPair("scaleway_instance_server.main", "public_ips.0.id", "scaleway_instance_ip.ip1", "id"),
					resource.TestCheckResourceAttrPair("scaleway_instance_server.main", "public_ips.1.id", "scaleway_instance_ip.ip2", "id"),
				),
			},
			{
				Config: `
					resource "scaleway_instance_ip" "ip1" {
						type = "routed_ipv4"
					}

					resource "scaleway_instance_ip" "ip2" {
						type = "routed_ipv4"
					}

					resource "scaleway_instance_server" "main" {
						name = "tf-tests-instance-server-ips"
						ip_ids = [scaleway_instance_ip.ip2.id]
						image = "ubuntu_jammy"
						type  = "PRO2-XXS"
						state = "stopped"
					}`,
				Check: resource.ComposeTestCheckFunc(
					arePrivateNICsPresent(tt, "scaleway_instance_server.main"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "routed_ip_enabled", "true"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "public_ips.#", "1"),
					resource.TestCheckResourceAttrPair("scaleway_instance_server.main", "public_ips.0.id", "scaleway_instance_ip.ip2", "id"),
				),
			},
		},
	})
}

func TestAccServer_IPRemoved(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_ip" "main" {}

					resource "scaleway_instance_server" "main" {
						name = "tf-tests-instance-server-ip-removed"
						ip_id = scaleway_instance_ip.main.id
						image = "ubuntu_jammy"
						type  = "PRO2-XXS"
						state = "stopped"
					}`,
				Check: resource.ComposeTestCheckFunc(
					arePrivateNICsPresent(tt, "scaleway_instance_server.main"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "public_ips.#", "1"),
					resource.TestCheckResourceAttrPair("scaleway_instance_server.main", "public_ips.0.id", "scaleway_instance_ip.main", "id"),
					resource.TestCheckResourceAttrPair("scaleway_instance_server.main", "public_ips.0.address", "scaleway_instance_server.main", "public_ip"),
				),
			},
			{
				Config: `
					resource "scaleway_instance_ip" "main" {}

					resource "scaleway_instance_server" "main" {
						name = "tf-tests-instance-server-ip-removed"
						image = "ubuntu_jammy"
						type  = "PRO2-XXS"
						state = "stopped"
					}`,
				Check: resource.ComposeTestCheckFunc(
					arePrivateNICsPresent(tt, "scaleway_instance_server.main"),
					serverHasNoIPAssigned(tt, "scaleway_instance_server.main"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "public_ips.#", "0"),
				),
			},
		},
	})
}

func TestAccServer_IPsRemoved(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_ip" "main" {
						type = "routed_ipv4"
					}

					resource "scaleway_instance_server" "main" {
						name = "tf-tests-instance-server-ips-removed"
						ip_ids = [scaleway_instance_ip.main.id]
						image = "ubuntu_jammy"
						type  = "PRO2-XXS"
						state = "stopped"
					}`,
				Check: resource.ComposeTestCheckFunc(
					arePrivateNICsPresent(tt, "scaleway_instance_server.main"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "routed_ip_enabled", "true"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "public_ips.#", "1"),
					resource.TestCheckResourceAttrPair("scaleway_instance_server.main", "public_ips.0.id", "scaleway_instance_ip.main", "id"),
				),
			},
			{
				Config: `
					resource "scaleway_instance_ip" "main" {
						type = "routed_ipv4"
					}

					resource "scaleway_instance_server" "main" {
						name = "tf-tests-instance-server-ips-removed"
						image = "ubuntu_jammy"
						type  = "PRO2-XXS"
						state = "stopped"
					}`,
				Check: resource.ComposeTestCheckFunc(
					arePrivateNICsPresent(tt, "scaleway_instance_server.main"),
					serverHasNoIPAssigned(tt, "scaleway_instance_server.main"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "routed_ip_enabled", "true"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "public_ips.#", "0"),
				),
			},
		},
	})
}

func TestAccServer_BlockExternal(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_block_volume" "volume" {
						iops = 5000
						size_in_gb = 10
					}

					resource "scaleway_instance_server" "main" {
						image = "ubuntu_jammy"
						type  = "PLAY2-PICO"
						additional_volume_ids = [scaleway_block_volume.volume.id]
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "type", "PLAY2-PICO"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "additional_volume_ids.#", "1"),
					resource.TestCheckResourceAttrPair("scaleway_instance_server.main", "additional_volume_ids.0", "scaleway_block_volume.volume", "id"),
				),
			},
			{
				Config: `
					resource "scaleway_block_volume" "volume" {
						iops = 5000
						size_in_gb = 10
					}

					resource "scaleway_instance_server" "main" {
						image = "ubuntu_jammy"
						type  = "PLAY2-PICO"
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "type", "PLAY2-PICO"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "additional_volume_ids.#", "0"),
				),
			},
			{
				Config: `
					resource "scaleway_block_volume" "volume" {
						iops = 5000
						size_in_gb = 10
					}

					resource "scaleway_instance_server" "main" {
						image = "ubuntu_jammy"
						type  = "PLAY2-PICO"
						additional_volume_ids = [scaleway_block_volume.volume.id]
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "type", "PLAY2-PICO"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "additional_volume_ids.#", "1"),
					resource.TestCheckResourceAttrPair("scaleway_instance_server.main", "additional_volume_ids.0", "scaleway_block_volume.volume", "id"),
				),
			},
			{
				Config: `
					resource "scaleway_block_volume" "volume" {
						iops = 5000
						size_in_gb = 10
					}

					resource "scaleway_instance_volume" "volume" {
						type = "b_ssd"
						size_in_gb = 10
					}

					resource "scaleway_instance_server" "main" {
						image = "ubuntu_jammy"
						type  = "PLAY2-PICO"
						additional_volume_ids = [scaleway_block_volume.volume.id, scaleway_instance_volume.volume.id]
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "type", "PLAY2-PICO"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "additional_volume_ids.#", "2"),
					resource.TestCheckResourceAttrPair("scaleway_instance_server.main", "additional_volume_ids.0", "scaleway_block_volume.volume", "id"),
					resource.TestCheckResourceAttrPair("scaleway_instance_server.main", "additional_volume_ids.1", "scaleway_instance_volume.volume", "id"),
				),
			},
		},
	})
}

func TestAccServer_BlockExternalRootVolume(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	serverID := ""

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			instancechecks.IsServerDestroyed(tt),
			instancechecks.IsServerRootVolumeDestroyed(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_server" "main" {
						name = "tf-tests-instance-block-external-root-volume"
						image = "ubuntu_jammy"
						type  = "PLAY2-PICO"
						root_volume {
							volume_type = "sbs_volume"
							size_in_gb = 50
							sbs_iops = 15000
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "type", "PLAY2-PICO"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "additional_volume_ids.#", "0"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "root_volume.0.volume_type", string(instanceSDK.VolumeVolumeTypeSbsVolume)),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "root_volume.0.sbs_iops", "15000"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "root_volume.0.size_in_gb", "50"),
					acctest.CheckResourceIDPersisted("scaleway_instance_server.main", &serverID),
				),
			},
			{
				Config: `
					resource "scaleway_instance_server" "main" {
						name = "tf-tests-instance-block-external-root-volume"
						image = "ubuntu_jammy"
						type  = "PLAY2-PICO"
						root_volume {
							volume_type = "sbs_volume"
							size_in_gb = 60
							sbs_iops = 15000
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "type", "PLAY2-PICO"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "additional_volume_ids.#", "0"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "root_volume.0.volume_type", string(instanceSDK.VolumeVolumeTypeSbsVolume)),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "root_volume.0.sbs_iops", "15000"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "root_volume.0.size_in_gb", "60"),
					acctest.CheckResourceIDPersisted("scaleway_instance_server.main", &serverID),
				),
			},
		},
	})
}

func TestAccServer_BlockExternalRootVolumeUpdate(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_server" "main" {
						name = "tf-tests-instance-block-external-root-volume-iops-update"
						image = "ubuntu_jammy"
						type  = "PLAY2-PICO"
						root_volume {
							volume_type = "sbs_volume"
							size_in_gb = 50
							sbs_iops = 5000
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "type", "PLAY2-PICO"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "additional_volume_ids.#", "0"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "root_volume.0.volume_type", string(instanceSDK.VolumeVolumeTypeSbsVolume)),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "root_volume.0.sbs_iops", "5000"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "root_volume.0.size_in_gb", "50"),
				),
			},
			{
				Config: `
					resource "scaleway_instance_server" "main" {
						name = "tf-tests-instance-block-external-root-volume-iops-update"
						image = "ubuntu_jammy"
						type  = "PLAY2-PICO"
						root_volume {
							volume_type = "sbs_volume"
							size_in_gb = 50
							sbs_iops = 15000
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "type", "PLAY2-PICO"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "additional_volume_ids.#", "0"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "root_volume.0.volume_type", string(instanceSDK.VolumeVolumeTypeSbsVolume)),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "root_volume.0.sbs_iops", "15000"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "root_volume.0.size_in_gb", "50"),
				),
			},
		},
	})
}

func TestAccServer_RootVolumeFromExternalSnapshot(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_instance_server" "main" {
						name = "tf-tests-instance-root-volume-from-external-snapshot"
						image = "ubuntu_jammy"
						type  = "PLAY2-PICO"
						root_volume {
							volume_type = "sbs_volume"
							size_in_gb = 50
							sbs_iops = 5000
						}
					}

					resource "scaleway_block_snapshot" "snapshot" {
						volume_id = scaleway_instance_server.main.root_volume.0.volume_id
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "type", "PLAY2-PICO"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "additional_volume_ids.#", "0"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "root_volume.0.volume_type", string(instanceSDK.VolumeVolumeTypeSbsVolume)),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "root_volume.0.sbs_iops", "5000"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "root_volume.0.size_in_gb", "50"),
				),
			},
			{
				Config: `
					resource "scaleway_instance_server" "main" {
						name = "tf-tests-instance-root-volume-from-external-snapshot"
						image = "ubuntu_jammy"
						type  = "PLAY2-PICO"
						root_volume {
							volume_type = "sbs_volume"
							size_in_gb = 50
							sbs_iops = 5000
						}
					}

					resource "scaleway_block_snapshot" "snapshot" {
						volume_id = scaleway_instance_server.main.root_volume.0.volume_id
					}

					resource "scaleway_block_volume" "volume" {
						snapshot_id = scaleway_block_snapshot.snapshot.id
						iops = 5000
					}

					resource "scaleway_instance_server" "from_snapshot" {
						name = "tf-tests-instance-root-volume-from-external-snapshot-2"
						type  = "PLAY2-PICO"
						root_volume {
							volume_type = "sbs_volume"
							volume_id = scaleway_block_volume.volume.id
						}
					}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "type", "PLAY2-PICO"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "additional_volume_ids.#", "0"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "root_volume.0.volume_type", string(instanceSDK.VolumeVolumeTypeSbsVolume)),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "root_volume.0.sbs_iops", "5000"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "root_volume.0.size_in_gb", "50"),
					resource.TestCheckResourceAttrPair("scaleway_instance_server.from_snapshot", "root_volume.0.volume_id", "scaleway_block_volume.volume", "id"),
				),
			},
		},
	})
}

func TestAccServer_PrivateNetworkMissingPNIC(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      instancechecks.IsServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_vpc_private_network pn {}

					resource "scaleway_instance_server" "main" {
						image = "ubuntu_jammy"
						type  = "PLAY2-PICO"
						private_network {
							pn_id = scaleway_vpc_private_network.pn.id
						}
					}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "type", "PLAY2-PICO"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "private_network.#", "1"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.main", "private_network.0.pn_id"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.main", "private_network.0.mac_address"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.main", "private_network.0.status"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.main", "private_network.0.zone"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.main", "private_network.0.pnic_id"),
					resource.TestCheckResourceAttrPair("scaleway_instance_server.main", "private_network.0.pn_id",
						"scaleway_vpc_private_network.pn", "id"),
				),
			},
			{
				Config: `
					resource scaleway_vpc_private_network pn {}

					resource "scaleway_instance_server" "main" {
						image = "ubuntu_jammy"
						type  = "PLAY2-PICO"
						private_network {
							pn_id = scaleway_vpc_private_network.pn.id
						}
					}

					resource scaleway_instance_private_nic pnic {
						private_network_id = scaleway_vpc_private_network.pn.id
						server_id = scaleway_instance_server.main.id
					}
`,
				ResourceName: "scaleway_instance_private_nic.pnic",
				ImportState:  true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					serverID := state.RootModule().Resources["scaleway_instance_server.main"].Primary.ID
					pnicID, exists := state.RootModule().Resources["scaleway_instance_server.main"].Primary.Attributes["private_network.0.pnic_id"]
					if !exists {
						return "", errors.New("private_network.0.pnic_id not found")
					}

					id := serverID + "/" + pnicID

					return id, nil
				},
				ImportStatePersist: true,
			},
			{ // We import private nic as a separate resource to trigger its deletion.
				Config: `
					resource scaleway_vpc_private_network pn {}

					resource "scaleway_instance_server" "main" {
						image = "ubuntu_jammy"
						type  = "PLAY2-PICO"
						private_network {
							pn_id = scaleway_vpc_private_network.pn.id
						}
					}

					resource scaleway_instance_private_nic pnic {
						private_network_id = scaleway_vpc_private_network.pn.id
						server_id = scaleway_instance_server.main.id
					}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "type", "PLAY2-PICO"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "private_network.#", "1"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.main", "private_network.0.pn_id"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.main", "private_network.0.mac_address"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.main", "private_network.0.status"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.main", "private_network.0.zone"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.main", "private_network.0.pnic_id"),
					resource.TestCheckResourceAttrPair("scaleway_instance_server.main", "private_network.0.pn_id",
						"scaleway_vpc_private_network.pn", "id"),
					func(state *terraform.State) error {
						serverPNICID, exists := state.RootModule().Resources["scaleway_instance_server.main"].Primary.Attributes["private_network.0.pnic_id"]
						if !exists {
							return errors.New("private_network.0.pnic_id not found")
						}
						localizedPNICID := state.RootModule().Resources["scaleway_instance_private_nic.pnic"].Primary.ID
						_, pnicID, _, err := zonal.ParseNestedID(localizedPNICID)
						if err != nil {
							return err
						}

						if serverPNICID != pnicID {
							return fmt.Errorf("expected server pnic (%s) to equal standalone pnic id (%s)", serverPNICID, pnicID)
						}

						return nil
					},
				),
			},
			{
				Config: `
					resource scaleway_vpc_private_network pn {}

					resource "scaleway_instance_server" "main" {
						image = "ubuntu_jammy"
						type  = "PLAY2-PICO"
						private_network {
							pn_id = scaleway_vpc_private_network.pn.id
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "type", "PLAY2-PICO"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "private_network.#", "1"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.main", "private_network.0.pn_id"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.main", "private_network.0.mac_address"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.main", "private_network.0.status"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.main", "private_network.0.zone"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.main", "private_network.0.pnic_id"),
					resource.TestCheckResourceAttrPair("scaleway_instance_server.main", "private_network.0.pn_id",
						"scaleway_vpc_private_network.pn", "id"),
				),
				ExpectNonEmptyPlan: true, // pnic get deleted and the plan is not empty after the apply as private_network is now missing
			},
			{
				Config: `
					resource scaleway_vpc_private_network pn {}

					resource "scaleway_instance_server" "main" {
						image = "ubuntu_jammy"
						type  = "PLAY2-PICO"
						private_network {
							pn_id = scaleway_vpc_private_network.pn.id
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "type", "PLAY2-PICO"),
					resource.TestCheckResourceAttr("scaleway_instance_server.main", "private_network.#", "1"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.main", "private_network.0.pn_id"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.main", "private_network.0.mac_address"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.main", "private_network.0.status"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.main", "private_network.0.zone"),
					resource.TestCheckResourceAttrSet("scaleway_instance_server.main", "private_network.0.pnic_id"),
					resource.TestCheckResourceAttrPair("scaleway_instance_server.main", "private_network.0.pn_id",
						"scaleway_vpc_private_network.pn", "id"),
				),
			},
		},
	})
}
