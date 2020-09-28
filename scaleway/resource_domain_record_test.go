package scaleway

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("scaleway_domain_record", &resource.Sweeper{
		Name: "scaleway_domain_record",
		//F:    testSweepDomainRecord,
	})
}

func TestAccScalewayDomainRecord(t *testing.T) {
	//resource.ParallelTest(t, resource.TestCase{
	//	PreCheck:     func() { testAccPreCheck(t) },
	//	Providers:    testAccProviders,
	//	CheckDestroy: testAccCheckScalewayAccountSSHKeyDestroy,
	//	Steps:        []resource.TestStep{
	//{
	//	Config: `
	//		resource "scaleway_account_ssh_key" "main" {
	//			name 	   = "` + name + `"
	//			public_key = "` + SSHKey + `"
	//		}
	//	`,
	//	Check: resource.ComposeTestCheckFunc(
	//		testAccCheckScalewayAccountSSHKeyExists("scaleway_account_ssh_key.main"),
	//		resource.TestCheckResourceAttr("scaleway_account_ssh_key.main", "name", name),
	//		resource.TestCheckResourceAttr("scaleway_account_ssh_key.main", "public_key", SSHKey),
	//	),
	//},
	//{
	//	Config: `
	//		resource "scaleway_account_ssh_key" "main" {
	//			name 	   = "` + name + `-updated"
	//			public_key = "` + SSHKey + `"
	//		}
	//	`,
	//	Check: resource.ComposeTestCheckFunc(
	//		testAccCheckScalewayAccountSSHKeyExists("scaleway_account_ssh_key.main"),
	//		resource.TestCheckResourceAttr("scaleway_account_ssh_key.main", "name", name+"-updated"),
	//		resource.TestCheckResourceAttr("scaleway_account_ssh_key.main", "public_key", SSHKey),
	//	),
	//},
	//},
	//})
}
