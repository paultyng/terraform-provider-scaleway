package scaleway_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
)

func TestAccScalewayDataSourceIamSSHKey_Basic(t *testing.T) {
	dataSourceIamSSHKey := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAILHy/M5FVm5ydLGcal3e5LNcfTalbeN7QL/ZGCvDEdqJ foobar@example.com"
	sshKeyName := "tf-test-ds-iam-ssh-key-basic"
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayIamSSHKeyDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_iam_ssh_key" "main" {
						name 	   = "%s"
						public_key = "%s"
					}
					`, sshKeyName, dataSourceIamSSHKey),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_iam_ssh_key" "main" {
						name 	   = "%s"
						public_key = "%s"
					}
					
					data "scaleway_iam_ssh_key" "prod" {
						name = "${scaleway_iam_ssh_key.main.name}"
					}
					
					data "scaleway_iam_ssh_key" "stg" {
						ssh_key_id = "${scaleway_iam_ssh_key.main.id}"
					}`, sshKeyName, dataSourceIamSSHKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayIamSSHKeyExists(tt, "data.scaleway_iam_ssh_key.prod"),
					resource.TestCheckResourceAttr("data.scaleway_iam_ssh_key.prod", "name", sshKeyName),
					resource.TestCheckResourceAttr("data.scaleway_iam_ssh_key.prod", "public_key", dataSourceIamSSHKey),
					testAccCheckScalewayIamSSHKeyExists(tt, "data.scaleway_iam_ssh_key.stg"),
					resource.TestCheckResourceAttr("data.scaleway_iam_ssh_key.stg", "name", sshKeyName),
					resource.TestCheckResourceAttr("data.scaleway_iam_ssh_key.stg", "public_key", dataSourceIamSSHKey),
				),
			},
		},
	})
}
