package secret_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	secretSDK "github.com/scaleway/scaleway-sdk-go/api/secret/v1beta1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/secret"
)

func TestAccSecret_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	secretName := "secretNameBasic"
	updatedName := "secretNameBasicUpdated"
	secretDescription := "secret description"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckSecretDestroy(tt),
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
					testAccCheckSecretExists(tt, "scaleway_secret.main"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "name", secretName),
					resource.TestCheckResourceAttr("scaleway_secret.main", "description", secretDescription),
					resource.TestCheckResourceAttr("scaleway_secret.main", "status", secretSDK.SecretStatusReady.String()),
					resource.TestCheckResourceAttr("scaleway_secret.main", "tags.0", "devtools"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "tags.1", "provider"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "tags.2", "terraform"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "tags.#", "3"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "ephemeral_policy.#", "0"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "type", "opaque"),
					resource.TestCheckResourceAttrSet("scaleway_secret.main", "updated_at"),
					resource.TestCheckResourceAttrSet("scaleway_secret.main", "created_at"),
					acctest.CheckResourceAttrUUID("scaleway_secret.main", "id"),
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
					testAccCheckSecretExists(tt, "scaleway_secret.main"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "name", updatedName),
					resource.TestCheckResourceAttr("scaleway_secret.main", "description", "update description"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "tags.0", "devtools"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "tags.#", "1"),
					acctest.CheckResourceAttrUUID("scaleway_secret.main", "id"),
				),
			},
			{
				Config: fmt.Sprintf(`
				resource "scaleway_secret" "main" {
				  name        = "%s"
				}
				`, secretName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretExists(tt, "scaleway_secret.main"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "name", secretName),
					resource.TestCheckResourceAttr("scaleway_secret.main", "description", ""),
					resource.TestCheckResourceAttr("scaleway_secret.main", "tags.#", "0"),
					acctest.CheckResourceAttrUUID("scaleway_secret.main", "id"),
				),
			},
		},
	})
}

func TestAccSecret_Path(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckSecretDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
				resource "scaleway_secret" "main" {
				  name = "test-secret-path-secret"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretExists(tt, "scaleway_secret.main"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "name", "test-secret-path-secret"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "path", "/"),
					acctest.CheckResourceAttrUUID("scaleway_secret.main", "id"),
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
					testAccCheckSecretExists(tt, "scaleway_secret.main"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "name", "test-secret-path-secret"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "path", "/test-secret-path"),
					acctest.CheckResourceAttrUUID("scaleway_secret.main", "id"),
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
					testAccCheckSecretExists(tt, "scaleway_secret.main"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "name", "test-secret-path-secret"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "path", "/test-secret-path-change"),
					acctest.CheckResourceAttrUUID("scaleway_secret.main", "id"),
				),
			},
			{
				Config: `
				resource "scaleway_secret" "main" {
				  name = "test-secret-path-secret"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretExists(tt, "scaleway_secret.main"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "name", "test-secret-path-secret"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "path", "/"),
					acctest.CheckResourceAttrUUID("scaleway_secret.main", "id"),
				),
			},
		},
	})
}

func TestAccSecret_Protected(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckSecretDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
				resource "scaleway_secret" "main" {
					name = "test-secret-protected-secret"
					protected = true
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretExists(tt, "scaleway_secret.main"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "name", "test-secret-protected-secret"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "path", "/"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "protected", "true"),
					acctest.CheckResourceAttrUUID("scaleway_secret.main", "id"),
				),
			},
			{
				Destroy: true,
				Config: `
				resource "scaleway_secret" "main" {
					name = "test-secret-protected-secret"
					protected = true
				}
				`,
				ExpectError: regexp.MustCompile(secret.ErrCannotDeleteProtectedSecret.Error()),
			},
			{
				Config: `
				resource "scaleway_secret" "main" {
					name = "test-secret-protected-secret"
					protected = false
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretExists(tt, "scaleway_secret.main"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "name", "test-secret-protected-secret"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "path", "/"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "protected", "false"),
					acctest.CheckResourceAttrUUID("scaleway_secret.main", "id"),
				),
			},
		},
	})
}

func TestAccSecret_EphemeralPolicy(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckSecretDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
				resource "scaleway_secret" "main" {
					name = "test-secret-policy-secret"
					ephemeral_policy {
						ttl = "30m"
						action = "disable"
					}
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretExists(tt, "scaleway_secret.main"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "name", "test-secret-policy-secret"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "path", "/"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "ephemeral_policy.#", "1"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "ephemeral_policy.0.ttl", "30m0s"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "ephemeral_policy.0.action", "disable"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "ephemeral_policy.0.expires_once_accessed", "false"),
					acctest.CheckResourceAttrUUID("scaleway_secret.main", "id"),
				),
			},
			{
				Config: `
				resource "scaleway_secret" "main" {
					name = "test-secret-policy-secret"
					ephemeral_policy {
						ttl = "5h"
						action = "delete"
						expires_once_accessed = true
					}
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretExists(tt, "scaleway_secret.main"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "name", "test-secret-policy-secret"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "path", "/"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "ephemeral_policy.#", "1"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "ephemeral_policy.0.ttl", "5h0m0s"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "ephemeral_policy.0.action", "delete"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "ephemeral_policy.0.expires_once_accessed", "true"),
					acctest.CheckResourceAttrUUID("scaleway_secret.main", "id"),
				),
			},
			{
				Config: `
				resource "scaleway_secret" "main" {
					name = "test-secret-policy-secret"
					ephemeral_policy {
						action = "delete"
						expires_once_accessed = true
					}
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretExists(tt, "scaleway_secret.main"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "name", "test-secret-policy-secret"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "path", "/"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "ephemeral_policy.#", "1"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "ephemeral_policy.0.ttl", ""),
					resource.TestCheckResourceAttr("scaleway_secret.main", "ephemeral_policy.0.action", "delete"),
					resource.TestCheckResourceAttr("scaleway_secret.main", "ephemeral_policy.0.expires_once_accessed", "true"),
					acctest.CheckResourceAttrUUID("scaleway_secret.main", "id"),
				),
			},
		},
	})
}

func testAccCheckSecretExists(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		api, region, id, err := secret.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = api.GetSecret(&secretSDK.GetSecretRequest{
			SecretID: id,
			Region:   region,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckSecretDestroy(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_secret" {
				continue
			}

			api, region, id, err := secret.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = api.GetSecret(&secretSDK.GetSecretRequest{
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
