package k8s_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/regional"
)

func TestAccDataSourceVersion_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckK8SClusterDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					data "scaleway_k8s_version" "by_name" {
						name = "1.26.2"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SVersionExists(tt, "data.scaleway_k8s_version.by_name"),
					resource.TestCheckResourceAttrSet("data.scaleway_k8s_version.by_name", "name"),
					resource.TestCheckResourceAttr("data.scaleway_k8s_version.by_name", "available_cnis.#", "4"),
					resource.TestCheckResourceAttr("data.scaleway_k8s_version.by_name", "available_cnis.0", "cilium"),
					resource.TestCheckResourceAttr("data.scaleway_k8s_version.by_name", "available_cnis.1", "calico"),
					resource.TestCheckResourceAttr("data.scaleway_k8s_version.by_name", "available_cnis.2", "kilo"),
					resource.TestCheckResourceAttr("data.scaleway_k8s_version.by_name", "available_cnis.3", "none"),
					resource.TestCheckResourceAttr("data.scaleway_k8s_version.by_name", "available_container_runtimes.#", "1"),
					resource.TestCheckResourceAttr("data.scaleway_k8s_version.by_name", "available_container_runtimes.0", "containerd"),
					resource.TestCheckResourceAttr("data.scaleway_k8s_version.by_name", "available_feature_gates.#", "5"),
					resource.TestCheckResourceAttr("data.scaleway_k8s_version.by_name", "available_feature_gates.0", "HPAScaleToZero"),
					resource.TestCheckResourceAttr("data.scaleway_k8s_version.by_name", "available_feature_gates.1", "GRPCContainerProbe"),
					resource.TestCheckResourceAttr("data.scaleway_k8s_version.by_name", "available_feature_gates.2", "ReadWriteOncePod"),
					resource.TestCheckResourceAttr("data.scaleway_k8s_version.by_name", "available_feature_gates.3", "ValidatingAdmissionPolicy"),
					resource.TestCheckResourceAttr("data.scaleway_k8s_version.by_name", "available_feature_gates.4", "CSINodeExpandSecret"),
				),
			},
		},
	})
}

func TestAccDataSourceVersion_Latest(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckK8SClusterDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					data "scaleway_k8s_version" "latest" {
						name = "latest"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckK8SVersionExists(tt, "data.scaleway_k8s_version.latest"),
					resource.TestCheckResourceAttrSet("data.scaleway_k8s_version.latest", "name"),
					resource.TestCheckResourceAttr("data.scaleway_k8s_version.latest", "name", testAccK8SClusterGetLatestK8SVersion(tt)),
				),
			},
		},
	})
}

func testAccCheckK8SVersionExists(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		region, name, err := regional.ParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		k8sAPI := k8s.NewAPI(tt.Meta.ScwClient())
		_, err = k8sAPI.GetVersion(&k8s.GetVersionRequest{
			Region:      region,
			VersionName: name,
		})
		if err != nil {
			return err
		}

		return nil
	}
}
