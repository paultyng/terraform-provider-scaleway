package iot_test

import (
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	iotSDK "github.com/scaleway/scaleway-sdk-go/api/iot/v1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/iot"
	objectchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/object/testfuncs"
	rdbchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/rdb/testfuncs"
)

const postgreSQLEngineName = "PostgreSQL"

func TestAccRoute_RDB(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	latestEngineVersion := rdbchecks.GetLatestEngineVersion(tt, postgreSQLEngineName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		// Destruction is done via the hub destruction.
		CheckDestroy: resource.ComposeTestCheckFunc(
			isHubDestroyed(tt),
			rdbchecks.IsInstanceDestroyed(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
						resource "scaleway_iot_hub" "minimal" {
							name         = "minimal"
							product_plan = "plan_shared"
						}

						resource "scaleway_rdb_instance" "minimal" {
							name           = "minimal"
							node_type      = "db-dev-s"
							engine         = %q
							is_ha_cluster  = false
							disable_backup = true
							user_name      = "root"
							password       = "T3stP4ssw0rdD0N0tUs3!"
						}

						resource "scaleway_iot_route" "default" {
							name   = "default"
							hub_id = scaleway_iot_hub.minimal.id
							topic  = "#"
							database {
								query  = "SELECT NOW()"
								host   = scaleway_rdb_instance.minimal.endpoint_ip
								port   = scaleway_rdb_instance.minimal.endpoint_port
								dbname = "rdb"
								username = scaleway_rdb_instance.minimal.user_name
								password = scaleway_rdb_instance.minimal.password
							}
						}
						`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					isHubPresent(tt, "scaleway_iot_hub.minimal"),
					isRoutePresent(tt, "scaleway_iot_route.default"),
					resource.TestCheckResourceAttrSet("scaleway_iot_route.default", "id"),
					resource.TestCheckResourceAttrSet("scaleway_iot_route.default", "hub_id"),
					resource.TestCheckResourceAttr("scaleway_iot_route.default", "topic", "#"),
					resource.TestCheckResourceAttr("scaleway_iot_route.default", "database.0.query", "SELECT NOW()"),
					resource.TestCheckResourceAttrSet("scaleway_iot_route.default", "database.0.host"),
					resource.TestCheckResourceAttrSet("scaleway_iot_route.default", "database.0.port"),
					resource.TestCheckResourceAttrSet("scaleway_iot_route.default", "database.0.username"),
					resource.TestCheckResourceAttr("scaleway_iot_route.default", "database.0.password", "T3stP4ssw0rdD0N0tUs3!"),
					resource.TestCheckResourceAttr("scaleway_iot_route.default", "database.0.dbname", "rdb"),
				),
			},
		},
	})
}

func TestAccRoute_S3(t *testing.T) {
	if !*acctest.UpdateCassettes {
		t.Skip("Skipping ObjectStorage test as this kind of resource can't be deleted before 24h")
	}
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	bucketName := sdkacctest.RandomWithPrefix("tf-tests-scaleway-iot-route-s3")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		// Destruction is done via the hub destruction.
		CheckDestroy: resource.ComposeAggregateTestCheckFunc(
			isHubDestroyed(tt),
			objectchecks.IsBucketDestroyed(tt),
		),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
						resource "scaleway_object_bucket" "minimal" {
							name = "%s"
						}

						resource "scaleway_iot_hub" "minimal" {
							name         = "minimal"
							product_plan = "plan_shared"
						}

						resource "scaleway_iot_route" "default" {
							name   = "default"
							hub_id = scaleway_iot_hub.minimal.id
							topic  = "#"

							s3 {
								bucket_region = scaleway_object_bucket.minimal.region
								bucket_name   = scaleway_object_bucket.minimal.name
								object_prefix = "foo"
								strategy      = "per_topic"
							}
							
							depends_on = [scaleway_object_bucket.minimal]
						}
						`, bucketName),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.minimal", true),
					isHubPresent(tt, "scaleway_iot_hub.minimal"),
					isRoutePresent(tt, "scaleway_iot_route.default"),
					resource.TestCheckResourceAttrSet("scaleway_iot_route.default", "id"),
					resource.TestCheckResourceAttrSet("scaleway_iot_route.default", "hub_id"),
					resource.TestCheckResourceAttr("scaleway_iot_route.default", "topic", "#"),
					resource.TestCheckResourceAttr("scaleway_iot_route.default", "s3.0.bucket_region", "fr-par"),
					resource.TestCheckResourceAttrSet("scaleway_iot_route.default", "s3.0.bucket_name"),
					resource.TestCheckResourceAttr("scaleway_iot_route.default", "s3.0.object_prefix", "foo"),
					resource.TestCheckResourceAttr("scaleway_iot_route.default", "s3.0.strategy", "per_topic"),
				),
			},
		},
	})
}

func TestAccRoute_REST(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		// Destruction is done via the hub destruction.
		CheckDestroy: isHubDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: `
						resource "scaleway_iot_hub" "minimal" {
							name         = "minimal"
							product_plan = "plan_shared"
						}

						resource "scaleway_iot_route" "default" {
							name   = "default"
							hub_id = scaleway_iot_hub.minimal.id
							topic  = "#"

							rest {
								verb = "get"
								uri  = "http://scaleway.com"
								headers = {
									X-terraform-test = "inprogress"
								}
							}
						}
						`,
				Check: resource.ComposeTestCheckFunc(
					isHubPresent(tt, "scaleway_iot_hub.minimal"),
					isRoutePresent(tt, "scaleway_iot_route.default"),
					resource.TestCheckResourceAttrSet("scaleway_iot_route.default", "id"),
					resource.TestCheckResourceAttrSet("scaleway_iot_route.default", "hub_id"),
					resource.TestCheckResourceAttr("scaleway_iot_route.default", "topic", "#"),
					resource.TestCheckResourceAttr("scaleway_iot_route.default", "rest.0.verb", "get"),
					resource.TestCheckResourceAttr("scaleway_iot_route.default", "rest.0.uri", "http://scaleway.com"),
					resource.TestCheckResourceAttr("scaleway_iot_route.default", "rest.0.headers.X-terraform-test", "inprogress"),
				),
			},
		},
	})
}

func isRoutePresent(tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		iotAPI, region, routeID, err := iot.NewAPIWithRegionAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = iotAPI.GetRoute(&iotSDK.GetRouteRequest{
			Region:  region,
			RouteID: routeID,
		})
		if err != nil {
			return err
		}

		return nil
	}
}
