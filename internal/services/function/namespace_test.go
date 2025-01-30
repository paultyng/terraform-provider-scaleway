package function_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	functionSDK "github.com/scaleway/scaleway-sdk-go/api/function/v1beta1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/function"
)

func TestAccFunctionNamespace_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckFunctionNamespaceDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_function_namespace main {
						name = "test-cr-ns-01"
						tags = ["tag1", "tag2"]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFunctionNamespaceExists(tt, "scaleway_function_namespace.main"),
					acctest.CheckResourceAttrUUID("scaleway_function_namespace.main", "id"),
					resource.TestCheckResourceAttr("scaleway_function_namespace.main", "tags.#", "2"),
					resource.TestCheckResourceAttr("scaleway_function_namespace.main", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("scaleway_function_namespace.main", "tags.1", "tag2"),
				),
			},
			{
				Config: `
					resource scaleway_function_namespace main {
						name = "test-cr-ns-01"
						description = "test function namespace 01"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFunctionNamespaceExists(tt, "scaleway_function_namespace.main"),
					resource.TestCheckResourceAttr("scaleway_function_namespace.main", "description", "test function namespace 01"),
					resource.TestCheckResourceAttr("scaleway_function_namespace.main", "name", "test-cr-ns-01"),
					acctest.CheckResourceAttrUUID("scaleway_function_namespace.main", "id"),
					resource.TestCheckResourceAttr("scaleway_function_namespace.main", "tags.#", "0"),
				),
			},
			{
				Config: `
					resource scaleway_function_namespace main {
						name = "test-cr-ns-01"
						environment_variables = {
							"test" = "test"
						}
						secret_environment_variables = {
							"test_secret" = "test_secret"
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFunctionNamespaceExists(tt, "scaleway_function_namespace.main"),
					resource.TestCheckResourceAttr("scaleway_function_namespace.main", "description", ""),
					resource.TestCheckResourceAttr("scaleway_function_namespace.main", "name", "test-cr-ns-01"),
					resource.TestCheckResourceAttr("scaleway_function_namespace.main", "environment_variables.test", "test"),
					resource.TestCheckResourceAttr("scaleway_function_namespace.main", "secret_environment_variables.test_secret", "test_secret"),
					resource.TestCheckResourceAttr("scaleway_function_namespace.main", "tags.#", "0"),

					acctest.CheckResourceAttrUUID("scaleway_function_namespace.main", "id"),
				),
			},
			{
				Config: `
                                        resource scaleway_function_namespace main {
                                                name = "test-cr-ns-01"
                                                environment_variables = {
                                                        "test" = "test"
                                                }
                                                secret_environment_variables = {
                                                        "test_secret" = "test_secret"
                                                }
                                                tags = ["tag1", "tag2"]
                                        }
                               `,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFunctionNamespaceExists(tt, "scaleway_function_namespace.main"),
					resource.TestCheckResourceAttr("scaleway_function_namespace.main", "description", ""),
					resource.TestCheckResourceAttr("scaleway_function_namespace.main", "name", "test-cr-ns-01"),
					resource.TestCheckResourceAttr("scaleway_function_namespace.main", "environment_variables.test", "test"),
					resource.TestCheckResourceAttr("scaleway_function_namespace.main", "secret_environment_variables.test_secret", "test_secret"),
					resource.TestCheckResourceAttr("scaleway_function_namespace.main", "tags.#", "2"),
					resource.TestCheckResourceAttr("scaleway_function_namespace.main", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("scaleway_function_namespace.main", "tags.1", "tag2"),

					acctest.CheckResourceAttrUUID("scaleway_function_namespace.main", "id"),
				),
			},
		},
	})
}

func TestAccFunctionNamespace_NoName(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckFunctionNamespaceDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_function_namespace main {
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFunctionNamespaceExists(tt, "scaleway_function_namespace.main"),
				),
			},
		},
	})
}

func TestAccFunctionNamespace_EnvironmentVariables(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckFunctionNamespaceDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_function_namespace main {
						name = "tf-env-test"
						environment_variables = {
							"test" = "test"
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFunctionNamespaceExists(tt, "scaleway_function_namespace.main"),
					resource.TestCheckResourceAttr("scaleway_function_namespace.main", "name", "tf-env-test"),
					resource.TestCheckResourceAttr("scaleway_function_namespace.main", "environment_variables.test", "test"),

					acctest.CheckResourceAttrUUID("scaleway_function_namespace.main", "id"),
				),
			},
			{
				Config: `
					resource scaleway_function_namespace main {
						name = "tf-env-test"
						environment_variables = {
							"foo" = "bar"
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFunctionNamespaceExists(tt, "scaleway_function_namespace.main"),
					resource.TestCheckResourceAttr("scaleway_function_namespace.main", "name", "tf-env-test"),
					resource.TestCheckResourceAttr("scaleway_function_namespace.main", "environment_variables.foo", "bar"),

					acctest.CheckResourceAttrUUID("scaleway_function_namespace.main", "id"),
				),
			},
		},
	})
}

func testAccCheckFunctionNamespaceExists(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		api, region, id, err := function.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = api.GetNamespace(&functionSDK.GetNamespaceRequest{
			NamespaceID: id,
			Region:      region,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckFunctionNamespaceDestroy(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_function_namespace" {
				continue
			}

			api, region, id, err := function.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = api.DeleteNamespace(&functionSDK.DeleteNamespaceRequest{
				NamespaceID: id,
				Region:      region,
			})

			if err == nil {
				return fmt.Errorf("function namespace (%s) still exists", rs.Primary.ID)
			}

			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}
