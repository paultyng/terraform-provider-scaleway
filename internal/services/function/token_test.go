package function_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	functionSDK "github.com/scaleway/scaleway-sdk-go/api/function/v1beta1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/function"
)

func TestAccFunctionToken_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	expiresAt := time.Now().Add(time.Hour * 24).Format(time.RFC3339)
	if !*acctest.UpdateCassettes {
		expiresAt = "2024-07-06T10:59:52+02:00"
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckFunctionTokenDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource scaleway_function_namespace main {
						name = "test-function-token-ns"
					}

					resource scaleway_function main {
						namespace_id = scaleway_function_namespace.main.id
						runtime = "node22"
						privacy = "private"
						handler = "handler.handle"
					}

					resource scaleway_function_token namespace {
						namespace_id = scaleway_function_namespace.main.id
						expires_at = "%s"
					}

					resource scaleway_function_token function {
						function_id = scaleway_function.main.id
					}
				`, expiresAt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFunctionTokenExists(tt, "scaleway_function_token.namespace"),
					testAccCheckFunctionTokenExists(tt, "scaleway_function_token.function"),
					acctest.CheckResourceAttrUUID("scaleway_function_token.namespace", "id"),
					acctest.CheckResourceAttrUUID("scaleway_function_token.function", "id"),
					resource.TestCheckResourceAttrSet("scaleway_function_token.namespace", "token"),
					resource.TestCheckResourceAttrSet("scaleway_function_token.function", "token"),
				),
			},
		},
	})
}

func testAccCheckFunctionTokenExists(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		api, region, id, err := function.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = api.GetToken(&functionSDK.GetTokenRequest{
			TokenID: id,
			Region:  region,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckFunctionTokenDestroy(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_function_token" {
				continue
			}

			api, region, id, err := function.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = api.DeleteToken(&functionSDK.DeleteTokenRequest{
				TokenID: id,
				Region:  region,
			})

			if err == nil {
				return fmt.Errorf("function token (%s) still exists", rs.Primary.ID)
			}

			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}
