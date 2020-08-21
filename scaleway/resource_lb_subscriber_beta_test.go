package scaleway

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccScalewayLbSubscriberBeta(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckScalewayLbBetaDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_lb_subscriber_beta sub01 {
						name = "test-sub"
					  	email_config {
							email = "foo@bar.com"
					  	}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_lb_subscriber_beta.sub01", "name", "test-sub"),
					resource.TestCheckResourceAttr("scaleway_lb_subscriber_beta.sub01", "email_config.#", "1"),
					resource.TestCheckResourceAttr("scaleway_lb_subscriber_beta.sub01", "email_config.0.email", "foo@bar.com"),
				),
			},
			{
				Config: `
					resource scaleway_lb_subscriber_beta sub01 {
						name = "test-sub"
					  	email_config {
							email = "zoo@bar.com"
					  	}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_lb_subscriber_beta.sub01", "name", "test-sub"),
					resource.TestCheckResourceAttr("scaleway_lb_subscriber_beta.sub01", "email_config.#", "1"),
					resource.TestCheckResourceAttr("scaleway_lb_subscriber_beta.sub01", "email_config.0.email", "zoo@bar.com"),
				),
			},
			{
				Config: `
					resource scaleway_lb_subscriber_beta sub01 {
						name = "test-sub"
					  	webhook_config {
							uri = "http://www.google.com/"
					  	}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_lb_subscriber_beta.sub01", "name", "test-sub"),
					resource.TestCheckResourceAttr("scaleway_lb_subscriber_beta.sub01", "webhook_config.#", "1"),
					resource.TestCheckResourceAttr("scaleway_lb_subscriber_beta.sub01", "webhook_config.0.uri", "http://www.google.com/"),
				),
			},
		},
	})
}
