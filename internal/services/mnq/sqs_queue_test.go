package mnq_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	accountSDK "github.com/scaleway/scaleway-sdk-go/api/account/v3"
	mnqSDK "github.com/scaleway/scaleway-sdk-go/api/mnq/v1beta1"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/httperrors"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/meta"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/provider"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/account"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/mnq"
	"github.com/stretchr/testify/require"
)

func TestAccSQSQueue_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	ctx := context.Background()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      isSQSQueueDestroyed(ctx, tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_account_project main {
						name = "tf_tests_mnq_sqs_queue_basic"
					}

					resource scaleway_mnq_sqs main {
						project_id = scaleway_account_project.main.id
					}

					resource scaleway_mnq_sqs_credentials main {
						project_id = scaleway_mnq_sqs.main.project_id
						permissions {
							can_manage = true
						}
					}

					resource scaleway_mnq_sqs_queue main {
						project_id = scaleway_mnq_sqs.main.project_id
						name = "test-mnq-sqs-queue-basic"
						sqs_endpoint = scaleway_mnq_sqs.main.endpoint
						access_key = scaleway_mnq_sqs_credentials.main.access_key
						secret_key = scaleway_mnq_sqs_credentials.main.secret_key
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isSQSQueuePresent(ctx, tt, "scaleway_mnq_sqs_queue.main"),
					acctest.CheckResourceAttrUUID("scaleway_mnq_sqs_queue.main", "id"),
					resource.TestCheckResourceAttr("scaleway_mnq_sqs_queue.main", "name", "test-mnq-sqs-queue-basic"),
				),
			},
			{
				Config: `
					resource scaleway_account_project main {
						name = "tf_tests_mnq_sqs_queue_basic"
					}

					resource scaleway_mnq_sqs main {
						project_id = scaleway_account_project.main.id
					}

					resource scaleway_mnq_sqs_credentials main {
						project_id = scaleway_mnq_sqs.main.project_id
						permissions {
							can_manage = true
						}
					}

					resource scaleway_mnq_sqs_queue main {
						project_id = scaleway_mnq_sqs.main.project_id
						name = "test-mnq-sqs-queue-basic"
						sqs_endpoint = scaleway_mnq_sqs.main.endpoint
						access_key = scaleway_mnq_sqs_credentials.main.access_key
						secret_key = scaleway_mnq_sqs_credentials.main.secret_key

						message_max_age = 720
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					isSQSQueuePresent(ctx, tt, "scaleway_mnq_sqs_queue.main"),
					acctest.CheckResourceAttrUUID("scaleway_mnq_sqs_queue.main", "id"),
					resource.TestCheckResourceAttr("scaleway_mnq_sqs_queue.main", "message_max_age", "720"),
				),
			},
		},
	})
}

func TestAccSQSQueue_DefaultProject(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()

	ctx := context.Background()

	accountAPI := accountSDK.NewProjectAPI(tt.Meta.ScwClient())
	projectID := ""
	project, err := accountAPI.CreateProject(&accountSDK.ProjectAPICreateProjectRequest{
		Name: "tf_tests_mnq_sqs_queue_default_project",
	})
	require.NoError(t, err)

	projectID = project.ID

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { acctest.PreCheck(t) },
		ProviderFactories: func() map[string]func() (*schema.Provider, error) {
			metaProd, err := meta.NewMeta(ctx, &meta.Config{
				TerraformVersion: "terraform-tests",
				HTTPClient:       tt.Meta.HTTPClient(),
			})
			require.NoError(t, err)

			return map[string]func() (*schema.Provider, error){
				"scaleway": func() (*schema.Provider, error) {
					return provider.Provider(&provider.Config{Meta: metaProd})(), nil
				},
			}
		}(),
		CheckDestroy: resource.ComposeTestCheckFunc(
			isSQSQueueDestroyed(ctx, tt),
			func(_ *terraform.State) error {
				return accountAPI.DeleteProject(&accountSDK.ProjectAPIDeleteProjectRequest{
					ProjectID: projectID,
				})
			},
		),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource scaleway_mnq_sqs main {
						project_id = "%1s"
					}

					resource scaleway_mnq_sqs_credentials main {
						project_id = scaleway_mnq_sqs.main.project_id
						permissions {
							can_manage = true
						}
					}

					resource scaleway_mnq_sqs_queue main {
						project_id = scaleway_mnq_sqs.main.project_id
						name = "test-mnq-sqs-queue-basic"
						access_key = scaleway_mnq_sqs_credentials.main.access_key
						secret_key = scaleway_mnq_sqs_credentials.main.secret_key
					}
				`, projectID),
				Check: resource.ComposeTestCheckFunc(
					isSQSQueuePresent(ctx, tt, "scaleway_mnq_sqs_queue.main"),
					acctest.CheckResourceAttrUUID("scaleway_mnq_sqs_queue.main", "id"),
					resource.TestCheckResourceAttr("scaleway_mnq_sqs_queue.main", "name", "test-mnq-sqs-queue-basic"),
					resource.TestCheckResourceAttr("scaleway_mnq_sqs_queue.main", "project_id", projectID),
				),
			},
		},
	})
}

func isSQSQueuePresent(ctx context.Context, tt *acctest.TestTools, n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		region, _, queueName, err := mnq.DecomposeMNQID(rs.Primary.ID)
		if err != nil {
			return err
		}

		sqsClient, err := mnq.NewSQSClient(ctx, tt.Meta.HTTPClient(), region.String(), rs.Primary.Attributes["sqs_endpoint"], rs.Primary.Attributes["access_key"], rs.Primary.Attributes["secret_key"])
		if err != nil {
			return err
		}

		_, err = sqsClient.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: &queueName,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func isSQSQueueDestroyed(ctx context.Context, tt *acctest.TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_mnq_sqs_queue" {
				continue
			}

			region, projectID, queueName, err := mnq.DecomposeMNQID(rs.Primary.ID)
			if err != nil {
				return err
			}

			// Project may have been deleted, check for it first
			// Checking for Queue first may lead to an AccessDenied if project has been deleted
			accountAPI := account.NewProjectAPI(tt.Meta)
			_, err = accountAPI.GetProject(&accountSDK.ProjectAPIGetProjectRequest{
				ProjectID: projectID,
			})
			if err != nil {
				if httperrors.Is404(err) {
					return nil
				}

				return err
			}

			mnqAPI := mnqSDK.NewSqsAPI(tt.Meta.ScwClient())
			sqsInfo, err := mnqAPI.GetSqsInfo(&mnqSDK.SqsAPIGetSqsInfoRequest{
				Region:    region,
				ProjectID: projectID,
			})
			if err != nil {
				return err
			}

			// SQS may be disabled for project, this means the queue does not exist
			if sqsInfo.Status == mnqSDK.SqsInfoStatusDisabled {
				return nil
			}

			sqsClient, err := mnq.NewSQSClient(ctx, tt.Meta.HTTPClient(), region.String(), rs.Primary.Attributes["sqs_endpoint"], rs.Primary.Attributes["access_key"], rs.Primary.Attributes["secret_key"])
			if err != nil {
				return err
			}

			_, err = sqsClient.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
				QueueName: aws.String(queueName),
			})
			if err != nil {
				if mnq.IsAWSErrorCode(err, mnq.AWSErrNonExistentQueue) || mnq.IsAWSErrorCode(err, "AccessDeniedException") {
					return nil
				}

				return fmt.Errorf("failed to get queue url: %s", err)
			}

			if err == nil {
				return fmt.Errorf("mnq sqs queue (%s) still exists", rs.Primary.ID)
			}

			if !httperrors.Is404(err) {
				return err
			}
		}

		return nil
	}
}
