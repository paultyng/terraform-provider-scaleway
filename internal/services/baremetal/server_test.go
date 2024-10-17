package baremetal_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	baremetalSDK "github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/baremetal"
	baremetalchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/baremetal/testfuncs"
)

const SSHKeyBaremetal = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIM7HUxRyQtB2rnlhQUcbDGCZcTJg7OvoznOiyC9W6IxH opensource@scaleway.com"

func TestAccServer_Basic(t *testing.T) {
	// t.Skip("Skipping Baremetal Server test as no stock is available currently")
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	SSHKeyName := "TestAccServer_Basic"
	name := "TestAccServer_Basic"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      baremetalchecks.CheckServerDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "scaleway_baremetal_os" "my_os" {
					  zone    = "fr-par-1"
					  name    = "Ubuntu"
					  version = "22.04 LTS (Jammy Jellyfish)"
					}

					resource "scaleway_iam_ssh_key" "main" {
						name 	   = "%s"
						public_key = "%s"
					}
					
					resource "scaleway_baremetal_server" "base" {
						name        = "%s"
						zone        = "fr-par-1"
						description = "test a description"
						offer       = "EM-A115X-SSD"
						os    = data.scaleway_baremetal_os.my_os.os_id
					
						tags = [ "terraform-test", "scaleway_baremetal_server", "minimal" ]
						ssh_key_ids = [ scaleway_iam_ssh_key.main.id ]
					}
				`, SSHKeyName, SSHKeyBaremetal, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaremetalServerExists(tt, "scaleway_baremetal_server.base"),
					resource.TestCheckResourceAttr("scaleway_baremetal_server.base", "name", name),
					resource.TestCheckResourceAttr("scaleway_baremetal_server.base", "offer_id", "fr-par-1/f7241870-c383-4fa2-bbca-5189600df5c4"),
					resource.TestCheckResourceAttr("scaleway_baremetal_server.base", "os", "fr-par-1/96e5f0f2-d216-4de2-8a15-68730d877885"),
					resource.TestCheckResourceAttr("scaleway_baremetal_server.base", "description", "test a description"),
					resource.TestCheckResourceAttr("scaleway_baremetal_server.base", "tags.0", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_baremetal_server.base", "tags.1", "scaleway_baremetal_server"),
					resource.TestCheckResourceAttr("scaleway_baremetal_server.base", "tags.2", "minimal"),
					acctest.CheckResourceAttrUUID("scaleway_baremetal_server.base", "ssh_key_ids.0"),
				),
			},
			{
				// Trigger a reinstall and update tags
				Config: fmt.Sprintf(`
					data "scaleway_baremetal_os" "my_os" {
					  zone    = "fr-par-1"
					  name    = "Ubuntu"
					  version = "22.04 LTS (Jammy Jellyfish)"
					}

					resource "scaleway_iam_ssh_key" "main" {
						name 	   = "%s"
						public_key = "%s"
					}
					
					resource "scaleway_baremetal_server" "base" {
						name        = "%s"
						zone        = "fr-par-1"
						description = "test a description"
						offer       = "EM-A115X-SSD"
						os          = data.scaleway_baremetal_os.my_os.os_id
					
						tags = [ "terraform-test", "scaleway_baremetal_server", "minimal", "edited" ]
						ssh_key_ids = [ scaleway_iam_ssh_key.main.id ]
					}
				`, SSHKeyName, SSHKeyBaremetal, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaremetalServerExists(tt, "scaleway_baremetal_server.base"),
					resource.TestCheckResourceAttr("scaleway_baremetal_server.base", "name", name),
					resource.TestCheckResourceAttr("scaleway_baremetal_server.base", "offer_id", "fr-par-1/f7241870-c383-4fa2-bbca-5189600df5c4"),
					resource.TestCheckResourceAttr("scaleway_baremetal_server.base", "os", "fr-par-1/96e5f0f2-d216-4de2-8a15-68730d877885"),
					resource.TestCheckResourceAttr("scaleway_baremetal_server.base", "description", "test a description"),
					resource.TestCheckResourceAttr("scaleway_baremetal_server.base", "tags.#", "4"),
					resource.TestCheckResourceAttr("scaleway_baremetal_server.base", "tags.0", "terraform-test"),
					resource.TestCheckResourceAttr("scaleway_baremetal_server.base", "tags.1", "scaleway_baremetal_server"),
					resource.TestCheckResourceAttr("scaleway_baremetal_server.base", "tags.2", "minimal"),
					resource.TestCheckResourceAttr("scaleway_baremetal_server.base", "tags.3", "edited"),
					acctest.CheckResourceAttrUUID("scaleway_baremetal_server.base", "ssh_key_ids.0"),
				),
			},
		},
	})
}

func TestAccServer_RequiredInstallConfig(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      baremetalchecks.CheckServerDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_baremetal_server" "base" {
						name        = "TestAccServer_RequiredInstallConfig"
						zone        = "fr-par-1"
						offer       = "EM-A115X-SSD"
						os          = "7e865c16-1a63-4dc7-8181-dabc020fc21b" // Proxmox

						ssh_key_ids = []
					}`,
				ExpectError: regexp.MustCompile("attribute is required"),
			},
		},
	})
}

func TestAccServer_WithoutInstallConfig(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      baremetalchecks.CheckServerDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					data "scaleway_baremetal_offer" "my_offer" {
					  zone = "fr-par-1"
					  name = "EM-A115X-SSD"
					}

					resource "scaleway_baremetal_server" "base" {
					  name 			             = "TestAccScalewayBaremetalServer_WithoutInstallConfig"
                      zone     			         = "fr-par-1"
					  offer     				 = data.scaleway_baremetal_offer.my_offer.offer_id
					  install_config_afterward   = true
					}`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaremetalServerExists(tt, "scaleway_baremetal_server.base"),
					resource.TestCheckResourceAttr("scaleway_baremetal_server.base", "name", "TestAccScalewayBaremetalServer_WithoutInstallConfig"),
					resource.TestCheckResourceAttr("scaleway_baremetal_server.base", "offer_id", "fr-par-1/f7241870-c383-4fa2-bbca-5189600df5c4"),
					resource.TestCheckNoResourceAttr("scaleway_baremetal_server.base", "os"),
				),
			},
		},
	})
}

func TestAccServer_CreateServerWithOption(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	SSHKeyName := "TestAccScalewayBaremetalServer_CreateServerWithOption"
	name := "TestAccScalewayBaremetalServer_CreateServerWithOption"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      baremetalchecks.CheckServerDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				data "scaleway_baremetal_os" "my_os" {
				  zone    = "fr-par-1"
				  name    = "Ubuntu"
				  version = "22.04 LTS (Jammy Jellyfish)"
				}
				
				data "scaleway_baremetal_offer" "my_offer" {
				  zone = "fr-par-1"
				  name = "EM-A115X-SSD"
				}
				
				data "scaleway_baremetal_option" "private_network" {
				  zone = "fr-par-1"
				  name = "Private Network"
				}
				
				resource "scaleway_iam_ssh_key" "base" {
				  name       = "%s"
				  public_key = "%s"
				}
				
				resource "scaleway_baremetal_server" "base" {
				  name  = "%s"
				  zone  = "fr-par-1"
				  offer = data.scaleway_baremetal_offer.my_offer.offer_id
				  os    = data.scaleway_baremetal_os.my_os.os_id
				
				  ssh_key_ids = [scaleway_iam_ssh_key.base.id]
				  options {
					id = data.scaleway_baremetal_option.private_network.option_id
				  }
				}
				`, SSHKeyName, SSHKeyBaremetal, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaremetalServerExists(tt, "scaleway_baremetal_server.base"),
					testAccCheckBaremetalServerHasOptions(tt, "scaleway_baremetal_server.base"),
					resource.TestCheckResourceAttrPair("scaleway_baremetal_server.base", "options.0.id", "data.scaleway_baremetal_option.private_network", "option_id"),
					resource.TestCheckResourceAttr("scaleway_baremetal_server.base", "ips.#", "2"),
					resource.TestCheckResourceAttr("scaleway_baremetal_server.base", "ipv4.#", "1"),
					resource.TestCheckResourceAttr("scaleway_baremetal_server.base", "ipv4.0.version", "IPv4"),
					resource.TestCheckResourceAttr("scaleway_baremetal_server.base", "ipv6.#", "1"),
					resource.TestCheckResourceAttr("scaleway_baremetal_server.base", "ipv6.0.version", "IPv6"),
				),
			},
		},
	})
}

func TestAccServer_AddOption(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	SSHKeyName := "TestAccScalewayBaremetalServer_AddOption"
	name := "TestAccScalewayBaremetalServer_AddOption"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      baremetalchecks.CheckServerDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "scaleway_baremetal_os" "by_id" {
					  zone    = "fr-par-1"
					  name    = "Ubuntu"
					  version = "22.04 LTS (Jammy Jellyfish)"
					}
					
					data "scaleway_baremetal_offer" "my_offer" {
					  zone = "fr-par-1"
					  name = "EM-A115X-SSD"
					}
					
					resource "scaleway_iam_ssh_key" "base" {
					  name       = "%s"
					  public_key = "%s"
					}
					
					resource "scaleway_baremetal_server" "base" {
					  name  = "%s"
					  zone  = "fr-par-1"
					  offer = data.scaleway_baremetal_offer.my_offer.offer_id
					  os    = data.scaleway_baremetal_os.by_id.os_id
					
					  ssh_key_ids = [scaleway_iam_ssh_key.base.id]
					}
				`, SSHKeyName, SSHKeyBaremetal, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaremetalServerExists(tt, "scaleway_baremetal_server.base"),
				),
			},
			{
				Config: fmt.Sprintf(`
				data "scaleway_baremetal_os" "my_os" {
				  zone    = "fr-par-1"
				  name    = "Ubuntu"
				  version = "22.04 LTS (Jammy Jellyfish)"
				}
				
				data "scaleway_baremetal_offer" "my_offer" {
				  zone = "fr-par-1"
				  name = "EM-A115X-SSD"
				}
				
				data "scaleway_baremetal_option" "private_network" {
				  zone = "fr-par-1"
				  name = "Private Network"
				}
				
				resource "scaleway_iam_ssh_key" "base" {
				  name       = "%s"
				  public_key = "%s"
				}
				
				resource "scaleway_baremetal_server" "base" {
				  name  = "%s"
				  zone  = "fr-par-1"
				  offer = data.scaleway_baremetal_offer.my_offer.offer_id
				  os    = data.scaleway_baremetal_os.my_os.os_id
				
				  ssh_key_ids = [scaleway_iam_ssh_key.base.id]
				  options {
					id = data.scaleway_baremetal_option.private_network.option_id
				  }
				}
				`, SSHKeyName, SSHKeyBaremetal, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaremetalServerExists(tt, "scaleway_baremetal_server.base"),
					testAccCheckBaremetalServerHasOptions(tt, "scaleway_baremetal_server.base"),
					resource.TestCheckResourceAttrPair("scaleway_baremetal_server.base", "options.0.id", "data.scaleway_baremetal_option.private_network", "option_id"),
				),
			},
		},
	})
}

func TestAccServer_AddTwoOptionsThenDeleteOne(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	SSHKeyName := "TestAccScalewayBaremetalServer_AddTwoOptionsThenDeleteOne"
	name := "TestAccScalewayBaremetalServer_AddTwoOptionsThenDeleteOne"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      baremetalchecks.CheckServerDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "scaleway_baremetal_os" "by_id" {
					  zone    = "fr-par-1"
					  name    = "Ubuntu"
					  version = "22.04 LTS (Jammy Jellyfish)"
					}
					
					data "scaleway_baremetal_offer" "my_offer" {
					  zone = "fr-par-1"
					  name = "EM-A115X-SSD"
					}
					
					resource "scaleway_iam_ssh_key" "base" {
					  name       = "%s"
					  public_key = "%s"
					}
					
					resource "scaleway_baremetal_server" "base" {
					  name  = "%s"
					  zone  = "fr-par-1"
					  offer = data.scaleway_baremetal_offer.my_offer.offer_id
					  os    = data.scaleway_baremetal_os.by_id.os_id
					
					  ssh_key_ids = [scaleway_iam_ssh_key.base.id]
					}
				`, SSHKeyName, SSHKeyBaremetal, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaremetalServerExists(tt, "scaleway_baremetal_server.base"),
				),
			},
			{
				Config: fmt.Sprintf(`
					data "scaleway_baremetal_os" "my_os" {
					  zone    = "fr-par-1"
					  name    = "Ubuntu"
					  version = "22.04 LTS (Jammy Jellyfish)"
					}
					
					data "scaleway_baremetal_offer" "my_offer" {
					  zone = "fr-par-1"
					  name = "EM-A115X-SSD"
					}
					
					data "scaleway_baremetal_option" "remote_access" {
					  zone = "fr-par-1"
					  name = "Remote Access"
					}
					
					data "scaleway_baremetal_option" "private_network" {
					  zone = "fr-par-1"
					  name = "Private Network"
					}
					
					resource "scaleway_iam_ssh_key" "base" {
					  name       = "%s"
					  public_key = "%s"
					}
					
					resource "scaleway_baremetal_server" "base" {
					  name        = "%s"
					  zone        = "fr-par-1"
					  offer       = data.scaleway_baremetal_offer.my_offer.offer_id
					  os          = data.scaleway_baremetal_os.my_os.os_id
					  ssh_key_ids = [scaleway_iam_ssh_key.base.id]
					
					  options {
						id = data.scaleway_baremetal_option.private_network.option_id
					  }
					  options {
						id         = data.scaleway_baremetal_option.remote_access.option_id
						expires_at = "2025-07-06T09:00:00Z"
					  }
					}
				`, SSHKeyName, SSHKeyBaremetal, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaremetalServerExists(tt, "scaleway_baremetal_server.base"),
					testAccCheckBaremetalServerHasOptions(tt, "scaleway_baremetal_server.base"),
					resource.TestCheckTypeSetElemAttrPair("scaleway_baremetal_server.base", "options.*.id", "data.scaleway_baremetal_option.remote_access", "option_id"),
					resource.TestCheckTypeSetElemAttrPair("scaleway_baremetal_server.base", "options.*.id", "data.scaleway_baremetal_option.private_network", "option_id"),
					resource.TestCheckTypeSetElemNestedAttrs("scaleway_baremetal_server.base", "options.*", map[string]string{
						"id":         "fr-par-1/931df052-d713-4674-8b58-96a63244c8e2",
						"expires_at": "2025-07-06T09:00:00Z",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("scaleway_baremetal_server.base", "options.*", map[string]string{
						"id": "fr-par-1/cd4158d7-2d65-49be-8803-c4b8ab6f760c",
					}),
				),
			},
			{
				Config: fmt.Sprintf(`
					data "scaleway_baremetal_os" "my_os" {
					  zone    = "fr-par-1"
					  name    = "Ubuntu"
					  version = "22.04 LTS (Jammy Jellyfish)"
					}
					
					data "scaleway_baremetal_offer" "my_offer" {
					  zone = "fr-par-1"
					  name = "EM-A115X-SSD"
					}
					
					data "scaleway_baremetal_option" "remote_access" {
					  zone = "fr-par-1"
					  name = "Remote Access"
					}
					
					resource "scaleway_iam_ssh_key" "base" {
					  name       = "%s"
					  public_key = "%s"
					}
					
					resource "scaleway_baremetal_server" "base" {
					  name        = "%s"
					  zone        = "fr-par-1"
					  offer       = data.scaleway_baremetal_offer.my_offer.offer_id
					  os          = data.scaleway_baremetal_os.my_os.os_id
					  ssh_key_ids = [scaleway_iam_ssh_key.base.id]
					
					  options {
						id         = data.scaleway_baremetal_option.remote_access.option_id
						expires_at = "2025-07-06T09:00:00Z"
					  }
					}
				`, SSHKeyName, SSHKeyBaremetal, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaremetalServerExists(tt, "scaleway_baremetal_server.base"),
					testAccCheckBaremetalServerHasOptions(tt, "scaleway_baremetal_server.base"),
					resource.TestCheckResourceAttrPair("scaleway_baremetal_server.base", "options.0.id", "data.scaleway_baremetal_option.remote_access", "option_id"),
					resource.TestCheckTypeSetElemNestedAttrs("scaleway_baremetal_server.base", "options.*", map[string]string{
						"id":         "fr-par-1/931df052-d713-4674-8b58-96a63244c8e2",
						"expires_at": "2025-07-06T09:00:00Z",
					}),
				),
			},
		},
	})
}

func TestAccServer_CreateServerWithPrivateNetwork(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	SSHKeyName := "TestAccScalewayBaremetalServer_CreateServerWithPrivateNetwork"
	name := "TestAccScalewayBaremetalServer_CreateServerWithPrivateNetwork"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			baremetalchecks.CheckServerDestroy(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "scaleway_baremetal_os" "my_os" {
						zone = "fr-par-1"
						name = "Ubuntu"
						version = "22.04 LTS (Jammy Jellyfish)"						
					}

					data "scaleway_baremetal_offer" "my_offer" {
						zone = "fr-par-1"
						name = "EM-A115X-SSD"
					}

					data "scaleway_baremetal_option" "private_network" {
						zone = "fr-par-1"
						name = "Private Network"
					}

					resource "scaleway_vpc_private_network" "pn" {
						name = "baremetal_private_network"
					} 

					resource "scaleway_iam_ssh_key" "base" {
						name 	   = "%s"
						public_key = "%s"
					}
					
					resource "scaleway_baremetal_server" "base" {
						name        = "%s"
						zone        = "fr-par-1"
						offer       = data.scaleway_baremetal_offer.my_offer.offer_id
						os          = data.scaleway_baremetal_os.my_os.os_id
					
						ssh_key_ids = [ scaleway_iam_ssh_key.base.id ]
						options {
						  id = data.scaleway_baremetal_option.private_network.option_id
						}
						private_network {
						  id = scaleway_vpc_private_network.pn.id
						}
					}
				`, SSHKeyName, SSHKeyBaremetal, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaremetalServerExists(tt, "scaleway_baremetal_server.base"),
					testAccCheckBaremetalServerHasPrivateNetwork(tt, "scaleway_baremetal_server.base"),
					resource.TestCheckResourceAttrPair("scaleway_baremetal_server.base", "private_network.0.id", "scaleway_vpc_private_network.pn", "id"),
					resource.TestCheckResourceAttrSet("scaleway_baremetal_server.base", "private_ip.0.id"),
					resource.TestCheckResourceAttrSet("scaleway_baremetal_server.base", "private_ip.0.address"),
				),
			},
		},
	})
}

func TestAccServer_AddPrivateNetwork(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	SSHKeyName := "TestAccScalewayBaremetalServer_AddPrivateNetwork"
	name := "TestAccScalewayBaremetalServer_AddPrivateNetwork"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			baremetalchecks.CheckServerDestroy(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "scaleway_baremetal_os" "my_os" {
						zone = "fr-par-1"
						name = "Ubuntu"
						version = "22.04 LTS (Jammy Jellyfish)"						
					}

					data "scaleway_baremetal_offer" "my_offer" {
						zone = "fr-par-1"
						name = "EM-A115X-SSD"
					}

					data "scaleway_baremetal_option" "private_network" {
						zone = "fr-par-1"
						name = "Private Network"
					}

					resource "scaleway_vpc_private_network" "pn" {
						name = "baremetal_private_network"
					} 

					resource "scaleway_iam_ssh_key" "base" {
						name 	   = "%s"
						public_key = "%s"
					}
					
					resource "scaleway_baremetal_server" "base" {
						name        = "%s"
						zone        = "fr-par-1"
						offer       = data.scaleway_baremetal_offer.my_offer.offer_id
						os          = data.scaleway_baremetal_os.my_os.os_id
					
						ssh_key_ids = [ scaleway_iam_ssh_key.base.id ]
						options {
						  id = data.scaleway_baremetal_option.private_network.option_id
						}
					}
				`, SSHKeyName, SSHKeyBaremetal, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaremetalServerExists(tt, "scaleway_baremetal_server.base"),
				),
			},
			{
				Config: fmt.Sprintf(`
					data "scaleway_baremetal_os" "my_os" {
						zone = "fr-par-1"
						name = "Ubuntu"
						version = "22.04 LTS (Jammy Jellyfish)"						
					}

					data "scaleway_baremetal_offer" "my_offer" {
						zone = "fr-par-1"
						name = "EM-A115X-SSD"
					}

					data "scaleway_baremetal_option" "private_network" {
						zone = "fr-par-1"
						name = "Private Network"
					}

					resource "scaleway_vpc_private_network" "pn" {
						name = "baremetal_private_network"
					} 

					resource "scaleway_iam_ssh_key" "base" {
						name 	   = "%s"
						public_key = "%s"
					}
					
					resource "scaleway_baremetal_server" "base" {
						name        = "%s"
						zone        = "fr-par-1"
						offer       = data.scaleway_baremetal_offer.my_offer.offer_id
						os          = data.scaleway_baremetal_os.my_os.os_id
					
						ssh_key_ids = [ scaleway_iam_ssh_key.base.id ]
						options {
						  id = data.scaleway_baremetal_option.private_network.option_id
						}
						private_network {
						  id = scaleway_vpc_private_network.pn.id
						}
					}
				`, SSHKeyName, SSHKeyBaremetal, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaremetalServerExists(tt, "scaleway_baremetal_server.base"),
					testAccCheckBaremetalServerHasPrivateNetwork(tt, "scaleway_baremetal_server.base"),
					resource.TestCheckResourceAttrPair("scaleway_baremetal_server.base", "private_network.0.id", "scaleway_vpc_private_network.pn", "id"),
				),
			},
		},
	})
}

func TestAccServer_AddAnotherPrivateNetwork(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	SSHKeyName := "TestAccScalewayBaremetalServer_AddAnotherPrivateNetwork"
	name := "TestAccScalewayBaremetalServer_AddAnotherPrivateNetwork"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			baremetalchecks.CheckServerDestroy(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "scaleway_baremetal_os" "my_os" {
						zone = "fr-par-1"
						name = "Ubuntu"
						version = "22.04 LTS (Jammy Jellyfish)"						
					}

					data "scaleway_baremetal_offer" "my_offer" {
						zone = "fr-par-1"
						name = "EM-A115X-SSD"
					}

					data "scaleway_baremetal_option" "private_network" {
						zone = "fr-par-1"
						name = "Private Network"
					}

					resource "scaleway_vpc_private_network" "pn" {
						name = "baremetal_private_network"
					}

					resource "scaleway_iam_ssh_key" "base" {
						name 	   = "%s"
						public_key = "%s"
					}
					
					resource "scaleway_baremetal_server" "base" {
						name        = "%s"
						zone        = "fr-par-1"
						offer       = data.scaleway_baremetal_offer.my_offer.offer_id
						os          = data.scaleway_baremetal_os.my_os.os_id
					
						ssh_key_ids = [ scaleway_iam_ssh_key.base.id ]
						options {
						  id = data.scaleway_baremetal_option.private_network.option_id
						}
						private_network {
						  id = scaleway_vpc_private_network.pn.id
						}
					}
				`, SSHKeyName, SSHKeyBaremetal, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaremetalServerExists(tt, "scaleway_baremetal_server.base"),
					testAccCheckBaremetalServerHasPrivateNetwork(tt, "scaleway_baremetal_server.base"),
					resource.TestCheckResourceAttrPair("scaleway_baremetal_server.base", "private_network.0.id", "scaleway_vpc_private_network.pn", "id"),
				),
			},
			{
				Config: fmt.Sprintf(`
					data "scaleway_baremetal_os" "my_os" {
						zone = "fr-par-1"
						name = "Ubuntu"
						version = "22.04 LTS (Jammy Jellyfish)"						
					}

					data "scaleway_baremetal_offer" "my_offer" {
						zone = "fr-par-1"
						name = "EM-A115X-SSD"
					}

					data "scaleway_baremetal_option" "private_network" {
						zone = "fr-par-1"
						name = "Private Network"
					}

					resource "scaleway_vpc_private_network" "pn" {
						name = "baremetal_private_network"
					} 

					resource "scaleway_vpc_private_network" "pn2" {
						name = "baremetal_private_network2"
					} 

					resource "scaleway_iam_ssh_key" "base" {
						name 	   = "%s"
						public_key = "%s"
					}
					
					resource "scaleway_baremetal_server" "base" {
						name        = "%s"
						zone        = "fr-par-1"
						offer       = data.scaleway_baremetal_offer.my_offer.offer_id
						os          = data.scaleway_baremetal_os.my_os.os_id
					
						ssh_key_ids = [ scaleway_iam_ssh_key.base.id ]
						options {
						  id = data.scaleway_baremetal_option.private_network.option_id
						}
						private_network {
						  id = scaleway_vpc_private_network.pn.id
						}
						private_network {
						  id = scaleway_vpc_private_network.pn2.id
						}
					}
				`, SSHKeyName, SSHKeyBaremetal, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaremetalServerExists(tt, "scaleway_baremetal_server.base"),
					testAccCheckBaremetalServerHasPrivateNetwork(tt, "scaleway_baremetal_server.base"),
					resource.TestCheckTypeSetElemAttrPair("scaleway_baremetal_server.base", "private_network.*.id", "scaleway_vpc_private_network.pn", "id"),
					resource.TestCheckTypeSetElemAttrPair("scaleway_baremetal_server.base", "private_network.*.id", "scaleway_vpc_private_network.pn2", "id"),
				),
			},
		},
	})
}

func testAccCheckBaremetalServerExists(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		baremetalAPI, zonedID, err := baremetal.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = baremetalAPI.GetServer(&baremetalSDK.GetServerRequest{
			ServerID: zonedID.ID,
			Zone:     zonedID.Zone,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckBaremetalServerHasOptions(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		baremetalAPI, zonedID, err := baremetal.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		server, err := baremetalAPI.GetServer(&baremetalSDK.GetServerRequest{
			ServerID: zonedID.ID,
			Zone:     zonedID.Zone,
		})
		if err != nil {
			return err
		}

		if len(server.Options) == 0 {
			return fmt.Errorf("server (%s) has no options enabled", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckBaremetalServerHasPrivateNetwork(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		_, zonedID, err := baremetal.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		baremetalPrivateNetworkAPI, _, err := baremetal.NewPrivateNetworkAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		listPrivateNetworks, err := baremetalPrivateNetworkAPI.ListServerPrivateNetworks(&baremetalSDK.PrivateNetworkAPIListServerPrivateNetworksRequest{
			Zone:     zonedID.Zone,
			ServerID: &zonedID.ID,
		})
		if err != nil {
			return err
		}

		if len(listPrivateNetworks.ServerPrivateNetworks) == 0 {
			return fmt.Errorf("server (%s) has no private networks attached to it", rs.Primary.ID)
		}

		return nil
	}
}
