package applesilicon_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	applesiliconSDK "github.com/scaleway/scaleway-sdk-go/api/applesilicon/v1alpha1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/applesilicon"
)

func TestAccServer_Basic(t *testing.T) {
	t.Skip("Skipping AppleSilicon test as this kind of server can't be deleted before 24h")
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isServerDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_apple_silicon_server main {
						name = "test-m1"
						type = "M1-M"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isServerPresent(tt, "scaleway_apple_silicon_server.main"),
					resource.TestCheckResourceAttr("scaleway_apple_silicon_server.main", "name", "test-m1"),
					resource.TestCheckResourceAttr("scaleway_apple_silicon_server.main", "type", "M1-M"),
					// Computed
					resource.TestCheckResourceAttrSet("scaleway_apple_silicon_server.main", "ip"),
					resource.TestCheckResourceAttrSet("scaleway_apple_silicon_server.main", "vnc_url"),
					resource.TestCheckResourceAttrSet("scaleway_apple_silicon_server.main", "created_at"),
					resource.TestCheckResourceAttrSet("scaleway_apple_silicon_server.main", "deletable_at"),
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

		asAPI, zone, ID, err := applesilicon.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = asAPI.GetServer(&applesiliconSDK.GetServerRequest{
			ServerID: ID,
			Zone:     zone,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func isServerDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_apple_silicon_server" {
				continue
			}

			asAPI, zone, ID, err := applesilicon.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = asAPI.GetServer(&applesiliconSDK.GetServerRequest{
				ServerID: ID,
				Zone:     zone,
			})

			// If no error resource still exist
			if err == nil {
				return fmt.Errorf("server (%s) still exists", rs.Primary.ID)
			}

			// Unexpected api error we return it
			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}
