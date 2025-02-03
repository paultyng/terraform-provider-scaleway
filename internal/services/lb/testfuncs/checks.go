package lbtestfuncs

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	lb2 "github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/instance"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/lb"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/transport"
)

func IsIPDestroyed(tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_lb_ip" {
				continue
			}

			lbAPI, zone, ID, err := lb.NewAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			lbID, lbExist := rs.Primary.Attributes["lb_id"]
			if lbExist && len(lbID) > 0 {
				retryInterval := lb.DefaultWaitLBRetryInterval

				if transport.DefaultWaitRetryInterval != nil {
					retryInterval = *transport.DefaultWaitRetryInterval
				}

				_, err := lbAPI.WaitForLbInstances(&lb2.ZonedAPIWaitForLBInstancesRequest{
					Zone:          zone,
					LBID:          lbID,
					Timeout:       scw.TimeDurationPtr(instance.DefaultInstanceServerWaitTimeout),
					RetryInterval: &retryInterval,
				}, scw.WithContext(context.Background()))

				// Unexpected api error we return it
				if !httperrors.Is404(err) {
					return err
				}
			}

			err = retry.RetryContext(context.Background(), lb.RetryLbIPInterval, func() *retry.RetryError {
				_, errGet := lbAPI.GetIP(&lb2.ZonedAPIGetIPRequest{
					Zone: zone,
					IPID: ID,
				})
				if httperrors.Is403(errGet) {
					return retry.RetryableError(errGet)
				}

				return retry.NonRetryableError(errGet)
			})

			// If no error resource still exist
			if err == nil {
				return fmt.Errorf("IP (%s) still exists", rs.Primary.ID)
			}

			// Unexpected api error we return it
			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}
