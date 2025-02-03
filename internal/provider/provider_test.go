package provider_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/meta"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/provider"
	iamchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/iam/testfuncs"
	instancechecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/instance/testfuncs"
	"github.com/stretchr/testify/require"
)

func TestAccProvider_InstanceIPZones(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	ctx := context.Background()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { acctest.PreCheck(t) },
		ProviderFactories: func() map[string]func() (*schema.Provider, error) {
			metaProd, err := meta.NewMeta(ctx, &meta.Config{
				TerraformVersion: "terraform-tests",
				ForceZone:        scw.ZoneFrPar2,
				HTTPClient:       tt.Meta.HTTPClient(),
			})
			require.NoError(t, err)

			metaDev, err := meta.NewMeta(ctx, &meta.Config{
				TerraformVersion: "terraform-tests",
				ForceZone:        scw.ZoneFrPar1,
				HTTPClient:       tt.Meta.HTTPClient(),
			})
			require.NoError(t, err)

			return map[string]func() (*schema.Provider, error){
				"prod": func() (*schema.Provider, error) {
					return provider.Provider(&provider.Config{Meta: metaProd})(), nil
				},
				"dev": func() (*schema.Provider, error) {
					return provider.Provider(&provider.Config{Meta: metaDev})(), nil
				},
			}
		}(),
		CheckDestroy: iamchecks.CheckSSHKeyDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_instance_ip dev {
					  provider = "dev"
					}

					resource scaleway_instance_ip prod {
					  provider = "prod"
					}
`,
				Check: resource.ComposeTestCheckFunc(
					instancechecks.CheckIPExists(tt, "scaleway_instance_ip.prod"),
					instancechecks.CheckIPExists(tt, "scaleway_instance_ip.dev"),
					resource.TestCheckResourceAttr("scaleway_instance_ip.prod", "zone", "fr-par-2"),
					resource.TestCheckResourceAttr("scaleway_instance_ip.dev", "zone", "fr-par-1"),
				),
			},
		},
	})
}

func TestAccProvider_SSHKeys(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	SSHKeyName := "TestAccScalewayProvider_SSHKeys"
	SSHKey := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEEYrzDOZmhItdKaDAEqJQ4ORS2GyBMtBozYsK5kiXXX opensource@scaleway.com"

	ctx := context.Background()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { acctest.PreCheck(t) },
		ProviderFactories: func() map[string]func() (*schema.Provider, error) {
			metaProd, err := meta.NewMeta(ctx, &meta.Config{
				TerraformVersion: "terraform-tests",
				HTTPClient:       tt.Meta.HTTPClient(),
			})
			require.NoError(t, err)

			metaDev, err := meta.NewMeta(ctx, &meta.Config{
				TerraformVersion: "terraform-tests",
				HTTPClient:       tt.Meta.HTTPClient(),
			})
			require.NoError(t, err)

			return map[string]func() (*schema.Provider, error){
				"prod": func() (*schema.Provider, error) {
					return provider.Provider(&provider.Config{Meta: metaProd})(), nil
				},
				"dev": func() (*schema.Provider, error) {
					return provider.Provider(&provider.Config{Meta: metaDev})(), nil
				},
			}
		}(),
		CheckDestroy: iamchecks.CheckSSHKeyDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_account_ssh_key" "prod" {
						provider   = "prod"
						name 	   = "%[1]s"
						public_key = "%[2]s"
					}

					resource "scaleway_account_ssh_key" "dev" {
						provider   = "dev"
						name 	   = "%[1]s"
						public_key = "%[2]s"
					}
				`, SSHKeyName, SSHKey),
				Check: resource.ComposeTestCheckFunc(
					iamchecks.CheckSSHKeyExists(tt, "scaleway_account_ssh_key.prod"),
					iamchecks.CheckSSHKeyExists(tt, "scaleway_account_ssh_key.dev"),
				),
			},
		},
	})
}
