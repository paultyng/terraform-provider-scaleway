package scaleway

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	container "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/api/registry/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func init() {
	resource.AddTestSweepers("scaleway_container_namespace", &resource.Sweeper{
		Name:         "scaleway_container_namespace",
		F:            testSweepContainerNamespace,
		Dependencies: []string{"scaleway_container"},
	})
}

func testSweepContainerNamespace(_ string) error {
	return sweepRegions([]scw.Region{scw.RegionFrPar}, func(scwClient *scw.Client, region scw.Region) error {
		containerAPI := container.NewAPI(scwClient)
		l.Debugf("sweeper: destroying the container namespaces in (%s)", region)
		listNamespaces, err := containerAPI.ListNamespaces(
			&container.ListNamespacesRequest{
				Region: region,
			}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("error listing namespaces in (%s) in sweeper: %s", region, err)
		}

		for _, ns := range listNamespaces.Namespaces {
			_, err := containerAPI.DeleteNamespace(&container.DeleteNamespaceRequest{
				NamespaceID: ns.ID,
				Region:      region,
			})
			if err != nil {
				l.Debugf("sweeper: error (%s)", err)

				return fmt.Errorf("error deleting namespace in sweeper: %s", err)
			}
		}

		return nil
	})
}

func TestAccScalewayContainerNamespace_Basic(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayContainerNamespaceDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_container_namespace main {
						name = "test-cr-ns-01"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayContainerNamespaceExists(tt, "scaleway_container_namespace.main"),
					testCheckResourceAttrUUID("scaleway_container_namespace.main", "id"),
				),
			},
			{
				Config: `
					resource scaleway_container_namespace main {
						name = "test-cr-ns-01"
						description = "test container namespace 01"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayContainerNamespaceExists(tt, "scaleway_container_namespace.main"),
					resource.TestCheckResourceAttr("scaleway_container_namespace.main", "description", "test container namespace 01"),
					resource.TestCheckResourceAttr("scaleway_container_namespace.main", "name", "test-cr-ns-01"),
					testCheckResourceAttrUUID("scaleway_container_namespace.main", "id"),
				),
			},
			{
				Config: `
					resource scaleway_container_namespace main {
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
					testAccCheckScalewayContainerNamespaceExists(tt, "scaleway_container_namespace.main"),
					resource.TestCheckResourceAttr("scaleway_container_namespace.main", "description", ""),
					resource.TestCheckResourceAttr("scaleway_container_namespace.main", "name", "test-cr-ns-01"),
					resource.TestCheckResourceAttr("scaleway_container_namespace.main", "environment_variables.test", "test"),
					resource.TestCheckResourceAttr("scaleway_container_namespace.main", "secret_environment_variables.test_secret", "test_secret"),

					testCheckResourceAttrUUID("scaleway_container_namespace.main", "id"),
				),
			},
			{
				Config: `
					resource scaleway_container_namespace main {
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayContainerNamespaceExists(tt, "scaleway_container_namespace.main"),
					resource.TestCheckResourceAttrSet("scaleway_container_namespace.main", "name"),
					resource.TestCheckResourceAttrSet("scaleway_container_namespace.main", "registry_endpoint"),
					resource.TestCheckResourceAttrSet("scaleway_container_namespace.main", "registry_namespace_id"),
				),
			},
			{
				Config: `
					resource scaleway_container_namespace main {
						name = "tf-env-test"
						environment_variables = {
							"test" = "test"
						}
						secret_environment_variables = {
							"test_secret" = "test_secret"
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayContainerNamespaceExists(tt, "scaleway_container_namespace.main"),
					resource.TestCheckResourceAttr("scaleway_container_namespace.main", "name", "tf-env-test"),
					resource.TestCheckResourceAttr("scaleway_container_namespace.main", "environment_variables.test", "test"),
					resource.TestCheckResourceAttr("scaleway_container_namespace.main", "secret_environment_variables.test_secret", "test_secret"),

					testCheckResourceAttrUUID("scaleway_container_namespace.main", "id"),
				),
			},
			{
				Config: `
					resource scaleway_container_namespace main {
						name = "tf-env-test"
						environment_variables = {
							"foo" = "bar"
						}
						secret_environment_variables = {
							"foo_secret" = "bar_secret"
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayContainerNamespaceExists(tt, "scaleway_container_namespace.main"),
					resource.TestCheckResourceAttr("scaleway_container_namespace.main", "name", "tf-env-test"),
					resource.TestCheckResourceAttr("scaleway_container_namespace.main", "environment_variables.foo", "bar"),
					resource.TestCheckResourceAttr("scaleway_container_namespace.main", "secret_environment_variables.foo_secret", "bar_secret"),

					testCheckResourceAttrUUID("scaleway_container_namespace.main", "id"),
				),
			},
		},
	})
}

func TestAccScalewayContainerNamespace_DestroyRegistry(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckScalewayContainerNamespaceDestroy(tt),
			testAccCheckScalewayContainerRegistryDestroy(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_container_namespace main {
						region = "pl-waw"
						name = "test-cr-ns-01"
						destroy_registry = true
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayContainerNamespaceExists(tt, "scaleway_container_namespace.main"),
					testCheckResourceAttrUUID("scaleway_container_namespace.main", "id"),
				),
			},
		},
	})
}

func TestAccScalewayContainerNamespace_ImportWithSecrets(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckScalewayContainerNamespaceDestroy(tt),
			testAccCheckScalewayContainerRegistryDestroy(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_container_namespace main {
						name = "test-cr-ns-import-01"
						secret_environment_variables = {
							foo = "bar"
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayContainerNamespaceExists(tt, "scaleway_container_namespace.main"),
					testCheckResourceAttrUUID("scaleway_container_namespace.main", "id"),
					resource.TestCheckResourceAttr("scaleway_container_namespace.main", "secret_environment_variables.foo", "bar"),
				),
			},
			{
				Config: `
					resource scaleway_container_namespace main {
						name = "test-cr-ns-import-01"
						secret_environment_variables = {
							foo = "bar"
						}
					}

					resource scaleway_container_namespace import {
						name = "test-cr-ns-import-01"
						secret_environment_variables = {
							foo = "bar"
						}
					}
				`,
				ResourceName:       "scaleway_container_namespace.import",
				ImportState:        true,
				ImportStatePersist: true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					return state.RootModule().Resources["scaleway_container_namespace.main"].Primary.ID, nil
				},
			},
			{
				Config: `
					resource scaleway_container_namespace main {
						name = "test-cr-ns-import-01"
						secret_environment_variables = {
							foo = "bar"
						}
					}

					resource scaleway_container_namespace import {
						name = "test-cr-ns-import-01"
						secret_environment_variables = {
							foo = "bar"
						}
					}
				`,
				PlanOnly: true,
			},
		},
	})
}

func testAccCheckScalewayContainerNamespaceExists(tt *TestTools, n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		api, region, id, err := containerAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = api.GetNamespace(&container.GetNamespaceRequest{
			NamespaceID: id,
			Region:      region,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckScalewayContainerNamespaceDestroy(tt *TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_container_namespace" { //nolint:goconst
				continue
			}

			api, region, id, err := containerAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = api.DeleteNamespace(&container.DeleteNamespaceRequest{
				NamespaceID: id,
				Region:      region,
			})

			if err == nil {
				return fmt.Errorf("container namespace (%s) still exists", rs.Primary.ID)
			}

			if !is404Error(err) {
				return err
			}
		}

		return nil
	}
}

func testAccCheckScalewayContainerRegistryDestroy(tt *TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_container_namespace" {
				continue
			}

			api, region, _, err := registryAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = api.DeleteNamespace(&registry.DeleteNamespaceRequest{
				NamespaceID: rs.Primary.Attributes["registry_namespace_id"],
				Region:      region,
			})

			if err == nil {
				return fmt.Errorf("registry namespace (%s) still exists", rs.Primary.Attributes["registry_namespace_id"])
			}

			if !is404Error(err) {
				return err
			}
		}

		return nil
	}
}
