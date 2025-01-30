package scwconfig_test

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/meta"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/provider"
	"github.com/stretchr/testify/require"
)

const ciAccessKey = "SCWXXXXXXXXXXXXXFAKE"

func TestAccDataSourceConfig_ActiveProfile(t *testing.T) {
	// Skip if we are running coverage tests on CI
	if accessKey := os.Getenv("SCW_ACCESS_KEY"); accessKey == ciAccessKey {
		t.Skip("Skipping TestAccDataSourceConfig_ActiveProfile")
	}
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	ctx := context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.PreCheck(t) },
		ProviderFactories: func() map[string]func() (*schema.Provider, error) {
			_ = os.Unsetenv("SCW_PROFILE")
			_ = os.Setenv("SCW_CONFIG_PATH", "./testfixture/test_config.yaml")
			metaDefault, err := meta.NewMeta(ctx, &meta.Config{
				TerraformVersion: "terraform-tests",
				HTTPClient:       tt.Meta.HTTPClient(),
			})
			require.NoError(t, err)
			return map[string]func() (*schema.Provider, error){
				"default": func() (*schema.Provider, error) {
					return provider.Provider(&provider.Config{Meta: metaDefault})(), nil
				},
			}
		}(),
		Steps: []resource.TestStep{
			{
				Config: `
					data "scaleway_config" "main" {
						provider = "default"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.scaleway_config.main", "access_key", "SCWXXXXXXXXXXXXXXXXX"),
					resource.TestCheckResourceAttr("data.scaleway_config.main", "access_key_source", meta.CredentialsSourceActiveProfile),
					resource.TestCheckResourceAttr("data.scaleway_config.main", "secret_key", "01234567-abcd-effe-dcba-012345678910"),
					resource.TestCheckResourceAttr("data.scaleway_config.main", "secret_key_source", meta.CredentialsSourceActiveProfile),
					resource.TestCheckResourceAttr("data.scaleway_config.main", "project_id", "11111111-2222-3333-4444-555555555555"),
					resource.TestCheckResourceAttr("data.scaleway_config.main", "project_id_source", meta.CredentialsSourceActiveProfile),
					resource.TestCheckResourceAttr("data.scaleway_config.main", "region", "nl-ams"),
					resource.TestCheckResourceAttr("data.scaleway_config.main", "region_source", meta.CredentialsSourceActiveProfile),
					resource.TestCheckResourceAttr("data.scaleway_config.main", "zone", "nl-ams-1"),
					resource.TestCheckResourceAttr("data.scaleway_config.main", "zone_source", meta.CredentialsSourceActiveProfile),
				),
			},
		},
	})
}

func TestAccDataSourceConfig_OtherProfile(t *testing.T) {
	// Skip if we are running coverage tests on CI
	if accessKey := os.Getenv("SCW_ACCESS_KEY"); accessKey == ciAccessKey {
		t.Skip("Skipping TestAccDataSourceConfig_OtherProfile")
	}
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	ctx := context.Background()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { acctest.PreCheck(t) },
		ProviderFactories: func() map[string]func() (*schema.Provider, error) {
			_ = os.Unsetenv("SCW_PROFILE")
			_ = os.Setenv("SCW_CONFIG_PATH", "./testfixture/test_config.yaml")
			_ = os.Setenv("SCW_PROFILE", "other")
			metaOther, err := meta.NewMeta(ctx, &meta.Config{
				TerraformVersion: "terraform-tests",
				HTTPClient:       tt.Meta.HTTPClient(),
			})
			require.NoError(t, err)
			return map[string]func() (*schema.Provider, error){
				"other": func() (*schema.Provider, error) {
					return provider.Provider(&provider.Config{Meta: metaOther})(), nil
				},
			}
		}(),
		Steps: []resource.TestStep{
			{
				Config: `
					data "scaleway_config" "main" {
						provider = "other"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.scaleway_config.main", "access_key", "SCWYYYYYYYYYYYYYYYYY"),
					resource.TestCheckResourceAttr("data.scaleway_config.main", "access_key_source", meta.CredentialsSourceActiveProfile),
					resource.TestCheckResourceAttr("data.scaleway_config.main", "secret_key", "99999999-9999-9999-9999-999999999999"),
					resource.TestCheckResourceAttr("data.scaleway_config.main", "secret_key_source", meta.CredentialsSourceActiveProfile),
					resource.TestCheckResourceAttr("data.scaleway_config.main", "project_id", "99999999-9999-9999-9999-999999999999"),
					resource.TestCheckResourceAttr("data.scaleway_config.main", "project_id_source", meta.CredentialsSourceActiveProfile),
					resource.TestCheckResourceAttr("data.scaleway_config.main", "region", "fr-par"),
					resource.TestCheckResourceAttr("data.scaleway_config.main", "region_source", meta.CredentialsSourceActiveProfile),
					resource.TestCheckResourceAttr("data.scaleway_config.main", "zone", "fr-par-1"),
					resource.TestCheckResourceAttr("data.scaleway_config.main", "zone_source", meta.CredentialsSourceActiveProfile),
				),
			},
		},
	})
}

func TestAccDataSourceConfig_MixedProfile(t *testing.T) {
	// Skip if we are running coverage tests on CI
	if accessKey := os.Getenv("SCW_ACCESS_KEY"); accessKey == ciAccessKey {
		t.Skip("Skipping TestAccDataSourceConfig_MixedProfile")
	}
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	ctx := context.Background()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { acctest.PreCheck(t) },
		ProviderFactories: func() map[string]func() (*schema.Provider, error) {
			_ = os.Unsetenv("SCW_PROFILE")
			_ = os.Setenv("SCW_CONFIG_PATH", "./testfixture/test_config.yaml")
			_ = os.Setenv("SCW_PROFILE", "incomplete")
			_ = os.Setenv("SCW_DEFAULT_PROJECT_ID", "77777777-7777-7777-7777-777777777777")
			metaMixed, err := meta.NewMeta(ctx, &meta.Config{
				TerraformVersion: "terraform-tests",
				HTTPClient:       tt.Meta.HTTPClient(),
			})
			require.NoError(t, err)
			return map[string]func() (*schema.Provider, error){
				"mixed": func() (*schema.Provider, error) {
					return provider.Provider(&provider.Config{Meta: metaMixed})(), nil
				},
			}
		}(),
		Steps: []resource.TestStep{
			{
				Config: `
					data "scaleway_config" "main" {
						provider = "mixed"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.scaleway_config.main", "access_key", "SCW11111111111111111"),
					resource.TestCheckResourceAttr("data.scaleway_config.main", "access_key_source", meta.CredentialsSourceActiveProfile),
					resource.TestCheckResourceAttr("data.scaleway_config.main", "secret_key", "11111111-1111-1111-1111-111111111111"),
					resource.TestCheckResourceAttr("data.scaleway_config.main", "secret_key_source", meta.CredentialsSourceActiveProfile),
					resource.TestCheckResourceAttr("data.scaleway_config.main", "project_id", "77777777-7777-7777-7777-777777777777"),
					resource.TestCheckResourceAttr("data.scaleway_config.main", "project_id_source", meta.CredentialsSourceEnvironment),
					resource.TestCheckResourceAttr("data.scaleway_config.main", "region", "nl-ams"),
					resource.TestCheckResourceAttr("data.scaleway_config.main", "region_source", meta.CredentialsSourceActiveProfile),
					resource.TestCheckResourceAttr("data.scaleway_config.main", "zone", "pl-waw-1"),
					resource.TestCheckResourceAttr("data.scaleway_config.main", "zone_source", meta.CredentialsSourceActiveProfile),
				),
			},
		},
	})
}
