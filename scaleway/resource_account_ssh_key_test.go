package scaleway

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	account "github.com/scaleway/scaleway-sdk-go/api/account/v2alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func init() {
	resource.AddTestSweepers("scaleway_account_ssh_key", &resource.Sweeper{
		Name: "scaleway_account_ssh_key",
		F:    testSweepAccountSSHKey,
	})
}

func testSweepAccountSSHKey(_ string) error {
	return sweepZones([]scw.Zone{scw.ZoneFrPar1}, func(scwClient *scw.Client, zone scw.Zone) error {
		accountAPI := account.NewAPI(scwClient)

		l.Debugf("sweeper: destroying the SSH keys")

		listSSHKeys, err := accountAPI.ListSSHKeys(&account.ListSSHKeysRequest{}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error listing SSH keys in sweeper: %s", err)
		}

		for _, sshKey := range listSSHKeys.SSHKeys {
			err := accountAPI.DeleteSSHKey(&account.DeleteSSHKeyRequest{
				SSHKeyID: sshKey.ID,
			})
			if err != nil {
				return fmt.Errorf("error deleting SSH key in sweeper: %s", err)
			}
		}

		return nil
	})
}

func TestAccScalewayAccountSSHKey_basic(t *testing.T) {
	name := "tf-test-account-ssh-key-basic"
	SSHKey := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEEYrzDOZmhItdKaDAEqJQ4ORS2GyBMtBozYsK5kiXXX opensource@scaleway.com"
	tt := NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayAccountSSHKeyDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_account_ssh_key" "main" {
						name 	   = "` + name + `"
						public_key = "` + SSHKey + `"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayAccountSSHKeyExists(tt, "scaleway_account_ssh_key.main"),
					resource.TestCheckResourceAttr("scaleway_account_ssh_key.main", "name", name),
					resource.TestCheckResourceAttr("scaleway_account_ssh_key.main", "public_key", SSHKey),
				),
			},
			{
				Config: `
					resource "scaleway_account_ssh_key" "main" {
						name 	   = "` + name + `-updated"
						public_key = "` + SSHKey + `"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayAccountSSHKeyExists(tt, "scaleway_account_ssh_key.main"),
					resource.TestCheckResourceAttr("scaleway_account_ssh_key.main", "name", name+"-updated"),
					resource.TestCheckResourceAttr("scaleway_account_ssh_key.main", "public_key", SSHKey),
				),
			},
		},
	})
}

func TestAccScalewayAccountSSHKey_WithNewLine(t *testing.T) {
	name := "tf-test-account-ssh-key-newline"
	SSHKey := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIDjfkdWCwkYlVQMDUfiZlVrmjaGOfBYnmkucssae8Iup opensource@scaleway.com"
	tt := NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayAccountSSHKeyDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_account_ssh_key" "main" {
						name 	   = "` + name + `"
						public_key = "\n\n` + SSHKey + `\n\n"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayAccountSSHKeyExists(tt, "scaleway_account_ssh_key.main"),
					resource.TestCheckResourceAttr("scaleway_account_ssh_key.main", "name", name),
					resource.TestCheckResourceAttr("scaleway_account_ssh_key.main", "public_key", SSHKey),
				),
			},
		},
	})
}

func testAccCheckScalewayAccountSSHKeyDestroy(tt *TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_account_ssh_key" {
				continue
			}

			accountAPI := accountAPI(tt.Meta)

			_, err := accountAPI.GetSSHKey(&account.GetSSHKeyRequest{
				SSHKeyID: rs.Primary.ID,
			})

			// If no error resource still exist
			if err == nil {
				return fmt.Errorf("SSH key (%s) still exists", rs.Primary.ID)
			}

			// Unexpected api error we return it
			if !is404Error(err) {
				return err
			}
		}

		return nil
	}
}

func testAccCheckScalewayAccountSSHKeyExists(tt *TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		accountAPI := accountAPI(tt.Meta)

		_, err := accountAPI.GetSSHKey(&account.GetSSHKeyRequest{
			SSHKeyID: rs.Primary.ID,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func TestAccScalewayAccountSSHKey_ChangeResourceName(t *testing.T) {
	name := "TestAccScalewayAccountSSHKey_ChangeResourceName"
	SSHKey := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAICJEoOOgQBLJPs4g/XcPTKT82NywNPpxeuA20FlOPlpO opensource@scaleway.com"
	tt := NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayAccountSSHKeyDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_account_ssh_key" "main" {
						name 	   = "` + name + `"
						public_key = "\n\n` + SSHKey + `\n\n"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayAccountSSHKeyExists(tt, "scaleway_account_ssh_key.main"),
					resource.TestCheckResourceAttr("scaleway_account_ssh_key.main", "name", name),
					resource.TestCheckResourceAttr("scaleway_account_ssh_key.main", "public_key", SSHKey),
				),
			},
			{
				Config: `
					resource "scaleway_account_ssh_key" "main2" {
						name 	   = "` + name + `"
						public_key = "\n\n` + SSHKey + `\n\n"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayAccountSSHKeyExists(tt, "scaleway_account_ssh_key.main"),
					resource.TestCheckResourceAttr("scaleway_account_ssh_key.main", "name", name),
					resource.TestCheckResourceAttr("scaleway_account_ssh_key.main", "public_key", SSHKey),
				),
			},
		},
	})
}
