package scaleway_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	secret "github.com/scaleway/scaleway-sdk-go/api/secret/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/logging"
	"github.com/scaleway/terraform-provider-scaleway/v2/scaleway"
)

func init() {
	resource.AddTestSweepers("scaleway_secret", &resource.Sweeper{
		Name: "scaleway_secret",
		F:    testSweepSecret,
	})
}

func testSweepSecret(_ string) error {
	return sweepRegions(scw.AllRegions, func(scwClient *scw.Client, region scw.Region) error {
		secretAPI := secret.NewAPI(scwClient)

		logging.L.Debugf("sweeper: deleting the secrets in (%s)", region)

		listSecrets, err := secretAPI.ListSecrets(&secret.ListSecretsRequest{Region: region}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error listing secrets in (%s) in sweeper: %s", region, err)
		}

		for _, se := range listSecrets.Secrets {
			err := secretAPI.DeleteSecret(&secret.DeleteSecretRequest{
				SecretID: se.ID,
				Region:   region,
			})
			if err != nil {
				logging.L.Debugf("sweeper: error (%s)", err)

				return fmt.Errorf("error deleting secret in sweeper: %s", err)
			}
		}

		return nil
	})
}

func TestAccScalewaySecret_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	secretName := "secretNameBasic"
	updatedName := "secretNameBasicUpdated"
	secretDescription := "secret description"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewaySecretDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "scaleway_secret" "main" {
				  name        = "%s"
				  description = "%s"
				  tags        = ["devtools", "provider", "terraform"]
				}
				`, secretName, secretDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewaySecretExists(tt, "scaleway_secret.main"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "name", secretName),
					resource.TestCheckResourceAttr("scaleway_secret.main", "description", secretDescription),
					resource.TestCheckResourceAttr("scaleway_secret.main", "status", secret.SecretStatusReady.String()),
					resource.TestCheckResourceAttr("scaleway_secret.main", "tags.0", "devtools"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "tags.1", "provider"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "tags.2", "terraform"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "tags.#", "3"),
					resource.TestCheckResourceAttrSet("scaleway_secret.main", "updated_at"),
					resource.TestCheckResourceAttrSet("scaleway_secret.main", "created_at"),
					testCheckResourceAttrUUID("scaleway_secret.main", "id"),
				),
			},
			{
				Config: fmt.Sprintf(`
				resource "scaleway_secret" "main" {
				  name        = "%s"
				  description = "update description"
				  tags        = ["devtools"]
				}
				`, updatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewaySecretExists(tt, "scaleway_secret.main"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "name", updatedName),
					resource.TestCheckResourceAttr("scaleway_secret.main", "description", "update description"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "tags.0", "devtools"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "tags.#", "1"),
					testCheckResourceAttrUUID("scaleway_secret.main", "id"),
				),
			},
			{
				Config: fmt.Sprintf(`
				resource "scaleway_secret" "main" {
				  name        = "%s"
				}
				`, secretName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewaySecretExists(tt, "scaleway_secret.main"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "name", secretName),
					resource.TestCheckResourceAttr("scaleway_secret.main", "description", ""),
					resource.TestCheckResourceAttr("scaleway_secret.main", "tags.#", "0"),
					testCheckResourceAttrUUID("scaleway_secret.main", "id"),
				),
			},
		},
	})
}

func TestAccScalewaySecret_Path(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewaySecretDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
				resource "scaleway_secret" "main" {
				  name = "test-secret-path-secret"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewaySecretExists(tt, "scaleway_secret.main"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "name", "test-secret-path-secret"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "path", "/"),
					testCheckResourceAttrUUID("scaleway_secret.main", "id"),
				),
			},
			{
				Config: `
				resource "scaleway_secret" "main" {
				  name = "test-secret-path-secret"
                  path = "/test-secret-path"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewaySecretExists(tt, "scaleway_secret.main"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "name", "test-secret-path-secret"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "path", "/test-secret-path"),
					testCheckResourceAttrUUID("scaleway_secret.main", "id"),
				),
			},
			{
				Config: `
				resource "scaleway_secret" "main" {
				  name = "test-secret-path-secret"
                  path = "/test-secret-path/"
				}
				`,
				PlanOnly: true,
			},
			{
				Config: `
				resource "scaleway_secret" "main" {
				  name = "test-secret-path-secret"
                  path = "/test-secret-path-change/"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewaySecretExists(tt, "scaleway_secret.main"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "name", "test-secret-path-secret"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "path", "/test-secret-path-change"),
					testCheckResourceAttrUUID("scaleway_secret.main", "id"),
				),
			},
			{
				Config: `
				resource "scaleway_secret" "main" {
				  name = "test-secret-path-secret"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewaySecretExists(tt, "scaleway_secret.main"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "name", "test-secret-path-secret"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "path", "/"),
					testCheckResourceAttrUUID("scaleway_secret.main", "id"),
				),
			},
		},
	})
}

func testAccCheckScalewaySecretExists(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		api, region, id, err := scaleway.SecretAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = api.GetSecret(&secret.GetSecretRequest{
			SecretID: id,
			Region:   region,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckScalewaySecretDestroy(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_secret" {
				continue
			}

			api, region, id, err := scaleway.SecretAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = api.GetSecret(&secret.GetSecretRequest{
				SecretID: id,
				Region:   region,
			})

			if err == nil {
				return fmt.Errorf("secret (%s) still exists", rs.Primary.ID)
			}

			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}
