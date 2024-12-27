package registry_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	registrySDK "github.com/scaleway/scaleway-sdk-go/api/registry/v1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/registry"
)

func TestAccNamespace_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isNamespaceDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_registry_namespace cr01 {
						region = "pl-waw"
						name = "test-cr-ns-01"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isNamespacePresent(tt, "scaleway_registry_namespace.cr01"),
					resource.TestCheckResourceAttr("scaleway_registry_namespace.cr01", "name", "test-cr-ns-01"),
					acctest.CheckResourceAttrUUID("scaleway_registry_namespace.cr01", "id"),
				),
			},
			{
				Config: `
					resource scaleway_registry_namespace cr01 {
						region = "pl-waw"
						name = "test-cr-ns-01"
						description = "test registry namespace 01"
						is_public = true
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isNamespacePresent(tt, "scaleway_registry_namespace.cr01"),
					resource.TestCheckResourceAttr("scaleway_registry_namespace.cr01", "description", "test registry namespace 01"),
					resource.TestCheckResourceAttr("scaleway_registry_namespace.cr01", "is_public", "true"),
					acctest.CheckResourceAttrUUID("scaleway_registry_namespace.cr01", "id"),
				),
			},
		},
	})
}

func isNamespacePresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		api, region, id, err := registry.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = api.GetNamespace(&registrySDK.GetNamespaceRequest{
			NamespaceID: id,
			Region:      region,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func isNamespaceDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_registry_namespace" {
				continue
			}

			api, region, id, err := registry.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = api.WaitForNamespace(&registrySDK.WaitForNamespaceRequest{
				NamespaceID: id,
				Region:      region,
			})

			if err == nil {
				return fmt.Errorf("namespace (%s) still exists", rs.Primary.ID)
			}

			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}
