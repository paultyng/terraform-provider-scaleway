package tem_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	temSDK "github.com/scaleway/scaleway-sdk-go/api/tem/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/tem"
)

func TestAccDomain_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	domainName := "terraform-rs.test.local"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isDomainDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource scaleway_tem_domain cr01 {
						name       = "%s"
						accept_tos = true
					}
				`, domainName),
				Check: resource.ComposeTestCheckFunc(
					isDomainPresent(tt, "scaleway_tem_domain.cr01"),
					resource.TestCheckResourceAttr("scaleway_tem_domain.cr01", "name", domainName),
					resource.TestCheckResourceAttrSet("scaleway_tem_domain.cr01", "dmarc_config"),
					resource.TestCheckResourceAttr("scaleway_tem_domain.cr01", "dmarc_name", "_dmarc.terraform-rs"),
					resource.TestCheckResourceAttr("scaleway_tem_domain.cr01", "last_error", ""), // last_error is deprecated
					acctest.CheckResourceAttrUUID("scaleway_tem_domain.cr01", "id"),
				),
			},
		},
	})
}

func TestAccDomain_Tos(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	domainName := "terraform-rs.test.local"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isDomainDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource scaleway_tem_domain cr01 {
						name       = "%s"
						accept_tos = false
					}
				`, domainName),
				ExpectError: regexp.MustCompile("you must accept"),
			},
		},
	})
}

func TestAccDomain_Autoconfig(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	subDomainName := "test-autoconfig"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isDomainDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`

					resource "scaleway_domain_zone" "test" {
  						domain    = "%s"
  						subdomain = "%s"
					}

					resource scaleway_tem_domain cr01 {
						name       = scaleway_domain_zone.test.id
						accept_tos = true
						autoconfig = true
					}

					resource scaleway_tem_domain_validation valid {
  						domain_id = scaleway_tem_domain.cr01.id
  						region = scaleway_tem_domain.cr01.region
						timeout = 3600
					}

				`, domainNameValidation, subDomainName),
				Check: resource.ComposeTestCheckFunc(
					isDomainPresent(tt, "scaleway_tem_domain.cr01"),
					resource.TestCheckResourceAttr("scaleway_tem_domain.cr01", "name", subDomainName+"."+domainNameValidation),
					resource.TestCheckResourceAttr("scaleway_tem_domain.cr01", "autoconfig", "true"),
					resource.TestCheckResourceAttrSet("scaleway_tem_domain.cr01", "dmarc_config"),
					resource.TestCheckResourceAttr("scaleway_tem_domain.cr01", "dmarc_name", "_dmarc"+"."+subDomainName),
					resource.TestCheckResourceAttr("scaleway_tem_domain.cr01", "last_error", ""), // last_error is deprecated
					acctest.CheckResourceAttrUUID("scaleway_tem_domain.cr01", "id"),
					resource.TestCheckResourceAttr("scaleway_tem_domain_validation.valid", "validated", "true"),
				),
			},
		},
	})
}

func TestAccDomain_AutoconfigUpdate(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	subDomainName := "test-autoconfig-update"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isDomainDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_domain_zone" "test" {
  						domain    = "%s"
  						subdomain = "%s"
					}

					resource scaleway_tem_domain cr01 {
						name       = scaleway_domain_zone.test.id
						accept_tos = true
						autoconfig = false
					}

				`, domainNameValidation, subDomainName),
				Check: resource.ComposeTestCheckFunc(
					isDomainPresent(tt, "scaleway_tem_domain.cr01"),
					resource.TestCheckResourceAttr("scaleway_tem_domain.cr01", "name", subDomainName+"."+domainNameValidation),
					resource.TestCheckResourceAttr("scaleway_tem_domain.cr01", "autoconfig", "false"),
					resource.TestCheckResourceAttrSet("scaleway_tem_domain.cr01", "dmarc_config"),
					resource.TestCheckResourceAttr("scaleway_tem_domain.cr01", "dmarc_name", "_dmarc"+"."+subDomainName),
					resource.TestCheckResourceAttr("scaleway_tem_domain.cr01", "last_error", ""), // last_error is deprecated
					acctest.CheckResourceAttrUUID("scaleway_tem_domain.cr01", "id"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_domain_zone" "test" {
  						domain    = "%s"
  						subdomain = "%s"
					}

					resource scaleway_tem_domain cr01 {
						name       = scaleway_domain_zone.test.id
						accept_tos = true
						autoconfig = true
					}

				`, domainNameValidation, subDomainName),
				Check: resource.ComposeTestCheckFunc(
					isDomainPresent(tt, "scaleway_tem_domain.cr01"),
					resource.TestCheckResourceAttr("scaleway_tem_domain.cr01", "autoconfig", "true"),
				),
			},
		},
	})
}

func isDomainPresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		api, region, id, err := tem.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = api.GetDomain(&temSDK.GetDomainRequest{
			DomainID: id,
			Region:   region,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func isDomainDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_tem_domain" {
				continue
			}

			api, region, id, err := tem.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = api.RevokeDomain(&temSDK.RevokeDomainRequest{
				Region:   region,
				DomainID: id,
			}, scw.WithContext(context.Background()))
			if err != nil {
				return err
			}

			_, err = tem.WaitForDomain(context.Background(), api, region, id, tem.DefaultDomainTimeout)
			if err != nil && !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}
