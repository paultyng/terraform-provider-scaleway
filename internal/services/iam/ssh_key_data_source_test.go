package iam_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	iamchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/iam/testfuncs"
)

func TestAccDataSourceSSHKey_Basic(t *testing.T) {
	dataSourceIamSSHKey := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAILHy/M5FVm5ydLGcal3e5LNcfTalbeN7QL/ZGCvDEdqJ foobar@example.com"
	dataSourceIamSSHKeyWithoutComment := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAILHy/M5FVm5ydLGcal3e5LNcfTalbeN7QL/ZGCvDEdqJ"
	sshKeyName := "tf-test-ds-iam-ssh-key-basic"
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      iamchecks.CheckSSHKeyDestroy(tt),
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
					iamchecks.CheckSSHKeyExists(tt, "data.scaleway_iam_ssh_key.prod"),
					resource.TestCheckResourceAttr("data.scaleway_iam_ssh_key.prod", "name", sshKeyName),
					resource.TestCheckResourceAttr("data.scaleway_iam_ssh_key.prod", "public_key", dataSourceIamSSHKeyWithoutComment),
					iamchecks.CheckSSHKeyExists(tt, "data.scaleway_iam_ssh_key.stg"),
					resource.TestCheckResourceAttr("data.scaleway_iam_ssh_key.stg", "name", sshKeyName),
					resource.TestCheckResourceAttr("data.scaleway_iam_ssh_key.stg", "public_key", dataSourceIamSSHKeyWithoutComment),
				),
			},
		},
	})
}
