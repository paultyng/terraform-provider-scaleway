package container_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	containerSDK "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/container"
)

func TestAccToken_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	expiresAt := time.Now().Add(time.Hour * 24).Format(time.RFC3339)
	if !*acctest.UpdateCassettes {
		expiresAt = "2023-01-05T14:12:46+01:00"
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isTokenDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource scaleway_container_namespace main {
						name = "test-container-token-ns"
					}

					resource scaleway_container main {
						namespace_id = scaleway_container_namespace.main.id
					}

					resource scaleway_container_token namespace {
						namespace_id = scaleway_container_namespace.main.id
						expires_at = "%s"
					}

					resource scaleway_container_token container {
						container_id = scaleway_container.main.id
					}
				`, expiresAt),
				Check: resource.ComposeTestCheckFunc(
					isTokenPresent(tt, "scaleway_container_token.namespace"),
					isTokenPresent(tt, "scaleway_container_token.container"),
					acctest.CheckResourceAttrUUID("scaleway_container_token.namespace", "id"),
					acctest.CheckResourceAttrUUID("scaleway_container_token.container", "id"),
					resource.TestCheckResourceAttrSet("scaleway_container_token.namespace", "token"),
					resource.TestCheckResourceAttrSet("scaleway_container_token.container", "token"),
				),
			},
		},
	})
}

func isTokenPresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		api, region, id, err := container.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = api.GetToken(&containerSDK.GetTokenRequest{
			TokenID: id,
			Region:  region,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func isTokenDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_container_token" {
				continue
			}

			api, region, id, err := container.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = api.DeleteToken(&containerSDK.DeleteTokenRequest{
				TokenID: id,
				Region:  region,
			})

			if err == nil {
				return fmt.Errorf("container token (%s) still exists", rs.Primary.ID)
			}

			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}
