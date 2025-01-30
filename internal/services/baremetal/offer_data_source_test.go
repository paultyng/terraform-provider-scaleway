package baremetal_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	baremetalSDK "github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/zonal"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/baremetal"
)

const (
	OfferName = "EM-B220E-NVME"
	Zone      = "fr-par-1"
	OfferID   = "206ea234-9097-4ae1-af68-6d2be09f47ed"
)

func TestAccDataSourceOffer_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	if !IsOfferAvailable(OfferID, Zone, tt) {
		t.Skip("Offer is out of stock")
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "scaleway_baremetal_offer" "test1" {
						zone = "fr-par-1"
						name = "%s"
					}
					
					data "scaleway_baremetal_offer" "test2" {
						offer_id = data.scaleway_baremetal_offer.test1.offer_id
					}
				`, OfferName),
				Check: resource.ComposeTestCheckFunc(
					isOfferPresent(tt, "data.scaleway_baremetal_offer.test1"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test1", "name", OfferName),
					isOfferPresent(tt, "data.scaleway_baremetal_offer.test2"),
					resource.TestCheckResourceAttrPair("data.scaleway_baremetal_offer.test2", "offer_id", "data.scaleway_baremetal_offer.test1", "offer_id"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "name", OfferName),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "commercial_range", "beryllium"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "include_disabled", "false"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "bandwidth", "1000000000"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "commercial_range", "beryllium"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "stock", "available"), // skipping this as stocks vary too much
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "cpu.0.name", "AMD EPYC 7232P"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "cpu.0.core_count", "8"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "cpu.0.frequency", "3100"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "cpu.0.thread_count", "16"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "disk.0.type", "NVMe"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "disk.0.capacity", "1024209543168"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "disk.1.type", "NVMe"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "disk.1.capacity", "1024209543168"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "memory.0.type", "DDR4"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "memory.0.capacity", "64000000000"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "memory.0.frequency", "2400"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "memory.0.is_ecc", "true"),
				),
			},
		},
	})
}

func TestAccDataSourceOffer_SubscriptionPeriodHourly(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	if !IsOfferAvailable(OfferID, Zone, tt) {
		t.Skip("Offer is out of stock")
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "scaleway_baremetal_offer" "test1" {
						zone = "fr-par-1"
						name = "%s"

						subscription_period = "hourly"
					}
					
					data "scaleway_baremetal_offer" "test2" {
						offer_id = data.scaleway_baremetal_offer.test1.offer_id
					}
				`, OfferName),
				Check: resource.ComposeTestCheckFunc(
					isOfferPresent(tt, "data.scaleway_baremetal_offer.test1"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test1", "name", OfferName),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test1", "subscription_period", "hourly"),
					isOfferPresent(tt, "data.scaleway_baremetal_offer.test2"),
					resource.TestCheckResourceAttrPair("data.scaleway_baremetal_offer.test2", "offer_id", "data.scaleway_baremetal_offer.test1", "offer_id"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "name", OfferName),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "subscription_period", "hourly"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "commercial_range", "beryllium"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "include_disabled", "false"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "bandwidth", "1000000000"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "commercial_range", "beryllium"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "stock", "available"), // skipping this as stocks vary too much
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "cpu.0.name", "AMD EPYC 7232P"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "cpu.0.core_count", "8"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "cpu.0.frequency", "3100"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "cpu.0.thread_count", "16"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "disk.0.type", "NVMe"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "disk.0.capacity", "1024209543168"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "disk.1.type", "NVMe"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "disk.1.capacity", "1024209543168"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "memory.0.type", "DDR4"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "memory.0.capacity", "64000000000"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "memory.0.frequency", "2400"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "memory.0.is_ecc", "true"),
				),
			},
		},
	})
}

func TestAccDataSourceOffer_SubscriptionPeriodMonthly(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	if !IsOfferAvailable(OfferID, Zone, tt) {
		t.Skip("Offer is out of stock")
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "scaleway_baremetal_offer" "test1" {
						zone = "fr-par-1"
						name = "%s"

						subscription_period = "monthly"
					}
					
					data "scaleway_baremetal_offer" "test2" {
						offer_id = data.scaleway_baremetal_offer.test1.offer_id
					}
				`, OfferName),
				Check: resource.ComposeTestCheckFunc(
					isOfferPresent(tt, "data.scaleway_baremetal_offer.test1"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test1", "name", OfferName),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test1", "subscription_period", "monthly"),
					isOfferPresent(tt, "data.scaleway_baremetal_offer.test2"),
					resource.TestCheckResourceAttrPair("data.scaleway_baremetal_offer.test2", "offer_id", "data.scaleway_baremetal_offer.test1", "offer_id"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "name", OfferName),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "subscription_period", "monthly"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "commercial_range", "beryllium"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "include_disabled", "false"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "bandwidth", "1000000000"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "commercial_range", "beryllium"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "stock", "available"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "cpu.0.name", "AMD EPYC 7232P"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "cpu.0.core_count", "8"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "cpu.0.frequency", "3100"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "cpu.0.thread_count", "16"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "disk.0.type", "NVMe"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "disk.0.capacity", "1024209543168"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "disk.1.type", "NVMe"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "disk.1.capacity", "1024209543168"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "memory.0.type", "DDR4"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "memory.0.capacity", "64000000000"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "memory.0.frequency", "2400"),
					resource.TestCheckResourceAttr("data.scaleway_baremetal_offer.test2", "memory.0.is_ecc", "true"),
				),
			},
		},
	})
}

func isOfferPresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		zone, id, err := zonal.ParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		api := baremetalSDK.NewAPI(tt.Meta.ScwClient())
		_, err = baremetal.FindOfferByID(context.Background(), api, zone, id)
		if err != nil {
			return err
		}

		return nil
	}
}
