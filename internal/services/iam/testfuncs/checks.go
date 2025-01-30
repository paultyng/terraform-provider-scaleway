package iamtestfuncs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	iam2 "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/iam"
)

func CheckSSHKeyDestroy(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_iam_ssh_key" {
				continue
			}

			iamAPI := iam.NewAPI(tt.Meta)

			_, err := iamAPI.GetSSHKey(&iam2.GetSSHKeyRequest{
				SSHKeyID: rs.Primary.ID,
			})

			// If no error resource still exist
			if err == nil {
				return fmt.Errorf("SSH key (%s) still exists", rs.Primary.ID)
			}

			// Unexpected api error we return it
			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}

func CheckSSHKeyExists(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		iamAPI := iam.NewAPI(tt.Meta)

		_, err := iamAPI.GetSSHKey(&iam2.GetSSHKeyRequest{
			SSHKeyID: rs.Primary.ID,
		})
		if err != nil {
			return err
		}

		return nil
	}
}
