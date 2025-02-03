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

func TestAccFunction_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckFunctionDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_function_namespace main {}

					resource scaleway_function main {
						name = "foobar"
						namespace_id = scaleway_function_namespace.main.id
						runtime = "node22"
						privacy = "private"
						handler = "handler.handle"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFunctionExists(tt, "scaleway_function.main"),
					resource.TestCheckResourceAttr("scaleway_function.main", "name", "foobar"),
					resource.TestCheckResourceAttr("scaleway_function.main", "runtime", "node22"),
					resource.TestCheckResourceAttr("scaleway_function.main", "privacy", "private"),
					resource.TestCheckResourceAttr("scaleway_function.main", "handler", "handler.handle"),
					resource.TestCheckResourceAttrSet("scaleway_function.main", "namespace_id"),
					resource.TestCheckResourceAttrSet("scaleway_function.main", "region"),
				),
			},
		},
	})
}

func TestAccFunction_Timeout(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckFunctionDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_function_namespace main {}

					resource scaleway_function main {
						name = "foobar"
						namespace_id = scaleway_function_namespace.main.id
						runtime = "node22"
						privacy = "private"
						handler = "handler.handle"
						timeout = 10
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFunctionExists(tt, "scaleway_function.main"),
					resource.TestCheckResourceAttr("scaleway_function.main", "name", "foobar"),
					resource.TestCheckResourceAttr("scaleway_function.main", "runtime", "node22"),
					resource.TestCheckResourceAttr("scaleway_function.main", "privacy", "private"),
					resource.TestCheckResourceAttr("scaleway_function.main", "handler", "handler.handle"),
					resource.TestCheckResourceAttr("scaleway_function.main", "timeout", "10"),
				),
			},
		},
	})
}

func TestAccFunction_NoName(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckFunctionNamespaceDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_function_namespace main {}

					resource scaleway_function main {
						namespace_id = scaleway_function_namespace.main.id
						runtime = "node22"
						privacy = "private"
						handler = "handler.handle"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFunctionExists(tt, "scaleway_function.main"),
					resource.TestCheckResourceAttrSet("scaleway_function.main", "name"),
				),
			},
		},
	})
}

func TestAccFunction_EnvironmentVariables(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckFunctionDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_function_namespace main {}

					resource scaleway_function main {
						name = "foobar"
						namespace_id = scaleway_function_namespace.main.id
						runtime = "node22"
						privacy = "private"
						handler = "handler.handle"
						environment_variables = {
							"test" = "test"
						}

						secret_environment_variables = {
							"test_secret" = "test_secret"
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFunctionExists(tt, "scaleway_function.main"),
					resource.TestCheckResourceAttr("scaleway_function.main", "environment_variables.test", "test"),
					resource.TestCheckResourceAttr("scaleway_function.main", "secret_environment_variables.test_secret", "test_secret"),
				),
			},
			{
				Config: `
					resource scaleway_function_namespace main {}

					resource scaleway_function main {
						name = "foobar"
						namespace_id = scaleway_function_namespace.main.id
						runtime = "node22"
						privacy = "private"
						handler = "handler.handle"
						environment_variables = {
							"foo" = "bar"
						}
						
						secret_environment_variables = {
							"foo_secret" = "bar_secret"
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFunctionExists(tt, "scaleway_function.main"),
					resource.TestCheckResourceAttr("scaleway_function.main", "environment_variables.foo", "bar"),
					resource.TestCheckResourceAttr("scaleway_function.main", "secret_environment_variables.foo_secret", "bar_secret"),
				),
			},
		},
	})
}

func TestAccFunction_Upload(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckFunctionDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_function_namespace main {}

					resource scaleway_function main {
						name = "foobar"
						namespace_id = scaleway_function_namespace.main.id
						runtime = "go122"
						privacy = "private"
						handler = "Handle"
						zip_file = "testfixture/gofunction.zip"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFunctionExists(tt, "scaleway_function.main"),
				),
			},
		},
	})
}

func TestAccFunction_Deploy(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckFunctionDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_function_namespace main {}

					resource scaleway_function main {
						name = "foobar"
						namespace_id = scaleway_function_namespace.main.id
						runtime = "go122"
						privacy = "private"
						handler = "Handle"
						zip_file = "testfixture/gofunction.zip"
						deploy = true
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFunctionExists(tt, "scaleway_function.main"),
				),
			},
			{
				Config: `
					resource scaleway_function_namespace main {}

					resource scaleway_function main {
						name = "foobar"
						namespace_id = scaleway_function_namespace.main.id
						runtime = "go122"
						privacy = "private"
						handler = "Handle"
						zip_file = "testfixture/gofunction.zip"
						zip_hash = "value"
						deploy = true
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFunctionExists(tt, "scaleway_function.main"),
				),
			},
		},
	})
}

func TestAccFunction_HTTPOption(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckFunctionDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_function_namespace main {}

					resource scaleway_function main {
						name = "foobar"
						namespace_id = scaleway_function_namespace.main.id
						runtime = "node22"
						privacy = "private"
						handler = "handler.handle"
						http_option = "enabled"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFunctionExists(tt, "scaleway_function.main"),
					resource.TestCheckResourceAttr("scaleway_function.main", "http_option", functionSDK.FunctionHTTPOptionEnabled.String()),
				),
			},
			{
				Config: `
					resource scaleway_function_namespace main {}

					resource scaleway_function main {
						name = "foobar"
						namespace_id = scaleway_function_namespace.main.id
						runtime = "node22"
						privacy = "private"
						handler = "handler.handle"
						http_option = "redirected"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFunctionExists(tt, "scaleway_function.main"),
					resource.TestCheckResourceAttr("scaleway_function.main", "http_option", functionSDK.FunctionHTTPOptionRedirected.String()),
				),
			},
			{
				Config: `
					resource scaleway_function_namespace main {}

					resource scaleway_function main {
						name = "foobar"
						namespace_id = scaleway_function_namespace.main.id
						runtime = "node22"
						privacy = "private"
						handler = "handler.handle"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFunctionExists(tt, "scaleway_function.main"),
					resource.TestCheckResourceAttr("scaleway_function.main", "http_option", functionSDK.FunctionHTTPOptionEnabled.String()),
				),
			},
		},
	})
}

func TestAccFunction_Sandbox(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckFunctionDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_function_namespace main {}

					resource scaleway_function main {
						name = "foobar"
						namespace_id = scaleway_function_namespace.main.id
						runtime = "node22"
						privacy = "private"
						handler = "handler.handle"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFunctionExists(tt, "scaleway_function.main"),
					resource.TestCheckResourceAttrSet("scaleway_function.main", "sandbox"),
				),
			},
			{
				Config: `
					resource scaleway_function_namespace main {}

					resource scaleway_function main {
						name = "foobar"
						namespace_id = scaleway_function_namespace.main.id
						runtime = "node22"
						privacy = "private"
						handler = "handler.handle"
						sandbox = "v2"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFunctionExists(tt, "scaleway_function.main"),
					resource.TestCheckResourceAttr("scaleway_function.main", "sandbox", functionSDK.FunctionSandboxV2.String()),
				),
			},
			{
				Config: `
					resource scaleway_function_namespace main {}

					resource scaleway_function main {
						name = "foobar"
						namespace_id = scaleway_function_namespace.main.id
						runtime = "node22"
						privacy = "private"
						handler = "handler.handle"
						sandbox = "v1"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFunctionExists(tt, "scaleway_function.main"),
					resource.TestCheckResourceAttr("scaleway_function.main", "sandbox", functionSDK.FunctionSandboxV1.String()),
				),
			},
			{
				Config: `
					resource scaleway_function_namespace main {}

					resource scaleway_function main {
						name = "foobar"
						namespace_id = scaleway_function_namespace.main.id
						runtime = "node22"
						privacy = "private"
						handler = "handler.handle"
						sandbox = "v2"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFunctionExists(tt, "scaleway_function.main"),
					resource.TestCheckResourceAttr("scaleway_function.main", "sandbox", functionSDK.FunctionSandboxV2.String()),
				),
			},
		},
	})
}

func testAccCheckFunctionExists(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		api, region, id, err := function.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = api.GetFunction(&functionSDK.GetFunctionRequest{
			FunctionID: id,
			Region:     region,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckFunctionDestroy(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_function" {
				continue
			}

			api, region, id, err := function.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = api.DeleteFunction(&functionSDK.DeleteFunctionRequest{
				FunctionID: id,
				Region:     region,
			})

			if err == nil {
				return fmt.Errorf("function (%s) still exists", rs.Primary.ID)
			}

			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}
