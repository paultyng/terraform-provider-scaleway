package object_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/acctest"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/object"
	objectchecks "github.com/scaleway/terraform-provider-scaleway/v2/internal/services/object/testfuncs"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/types"
)

const (
	s3ACLGranteeAllUsers           = "AllUsers"
	s3ACLGranteeAuthenticatedUsers = "AuthenticatedUsers"
)

func TestAccObjectBucketACL_Basic(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	testBucketName := sdkacctest.RandomWithPrefix("tf-tests-scw-object-acl-basic")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      objectchecks.IsBucketDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_object_bucket" "main" {
						name = "%s"
						region = "%s"
					}

					resource "scaleway_object_bucket_acl" "main" {
						bucket = scaleway_object_bucket.main.id
						acl = "private"
					}
					`, testBucketName, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.main", true),
					resource.TestCheckResourceAttr("scaleway_object_bucket_acl.main", "bucket", testBucketName),
					resource.TestCheckResourceAttr("scaleway_object_bucket_acl.main", "acl", "private"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_object_bucket" "main" {
						name = "%s"
						region = "%s"
					}

					resource "scaleway_object_bucket_acl" "main" {
						bucket = scaleway_object_bucket.main.id
						acl = "public-read"
					}
					`, testBucketName, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.main", true),
					resource.TestCheckResourceAttr("scaleway_object_bucket_acl.main", "bucket", testBucketName),
					resource.TestCheckResourceAttr("scaleway_object_bucket_acl.main", "acl", "public-read"),
				),
			},
		},
	})
}

func TestAccObjectBucketACL_Grantee(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	testBucketName := sdkacctest.RandomWithPrefix("tf-tests-scw-object-acl-grantee")

	ownerID := "105bdce1-64c0-48ab-899d-868455867ecf"
	ownerIDChild := "50ab77d5-56bd-4981-a118-4e0fa5309b59"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      objectchecks.IsBucketDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_object_bucket" "main" {
						name = "%[1]s"
						region = "%[3]s"
					}

					resource "scaleway_object_bucket_acl" "main" {
						bucket = scaleway_object_bucket.main.id
						access_control_policy {
						  	grant {
								grantee {
									id   = "%[2]s"
									type = "CanonicalUser"
								}
								permission = "FULL_CONTROL"
						  	}

						  	grant {
								grantee {
							  		id   = "%[2]s"
							  		type = "CanonicalUser"
								}
								permission = "WRITE"
						  	}

						  	owner {
								id = "%[2]s"
						  	}
						}
					}
					`, testBucketName, ownerID, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.main", true),
					resource.TestCheckResourceAttr("scaleway_object_bucket_acl.main", "bucket", testBucketName),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_object_bucket" "main" {
						name = "%[1]s"
						region = "%[4]s"
					}
			
					resource "scaleway_object_bucket_acl" "main" {
						bucket = scaleway_object_bucket.main.id
						access_control_policy {
						  	grant {
								grantee {
									id   = "%[2]s"
									type = "CanonicalUser"
								}
								permission = "FULL_CONTROL"
						  	}
			
						  	grant {
								grantee {
							  		id   = "%[2]s"
							  		type = "CanonicalUser"
								}
								permission = "WRITE"
						  	}
			
							grant {
								grantee {
								  	id   = "%[3]s"
								  	type = "CanonicalUser"
								}
								permission = "FULL_CONTROL"
							}
			
						  	owner {
								id = "%[2]s"
						  	}
						}
					}
				`, testBucketName, ownerID, ownerIDChild, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.main", true),
					resource.TestCheckResourceAttr("scaleway_object_bucket_acl.main", "bucket", testBucketName),
				),
			},
			{
				ResourceName:      "scaleway_object_bucket_acl.main",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccObjectBucketACL_GranteeWithOwner(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	testBucketName := sdkacctest.RandomWithPrefix("tf-tests-scw-object-acl-owner")
	ownerID := "105bdce1-64c0-48ab-899d-868455867ecf"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      objectchecks.IsBucketDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_object_bucket" "main" {
						name = "%[1]s"
						region = "%[3]s"
					}

					resource "scaleway_object_bucket_acl" "main" {
						bucket = scaleway_object_bucket.main.id
						expected_bucket_owner = "%[2]s"
						access_control_policy {
						  grant {
							grantee {
								id   = "%[2]s"
								type = "CanonicalUser"
							}
							permission = "FULL_CONTROL"
						  }

						  grant {
							grantee {
							  id   = "%[2]s"
							  type = "CanonicalUser"
							}
							permission = "WRITE"
						  }

						  owner {
							id = "%[2]s"
						  }
						}
					}
					`, testBucketName, ownerID, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.main", true),
					resource.TestCheckResourceAttr("scaleway_object_bucket_acl.main", "bucket", testBucketName),
				),
			},
		},
	})
}

func TestAccObjectBucketACL_WithBucketName(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	testBucketName := sdkacctest.RandomWithPrefix("tf-tests-scw-object-acl-name")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      objectchecks.IsBucketDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_object_bucket" "main" {
						name = "%s"
						region = "%s"
					}

					resource "scaleway_object_bucket_acl" "main" {
						bucket = scaleway_object_bucket.main.name
						acl = "public-read"

					}
					`, testBucketName, objectTestsMainRegion),
				ExpectError: regexp.MustCompile("api error NoSuchBucket: The specified bucket does not exist"),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_object_bucket" "main" {
						name = "%s"
						region = "%s"
					}
			
					resource "scaleway_object_bucket_acl" "main" {
						bucket = scaleway_object_bucket.main.name
						acl = "public-read"
						region = "%[2]s"
					}
					`, testBucketName, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.main", true),
					resource.TestCheckResourceAttr("scaleway_object_bucket_acl.main", "bucket", testBucketName),
					resource.TestCheckResourceAttr("scaleway_object_bucket_acl.main", "acl", "public-read"),
				),
			},
		},
	})
}

func TestAccObjectBucketACL_Remove(t *testing.T) {
	tt := acctest.NewTestTools(t)
	defer tt.Cleanup()
	testBucketName := sdkacctest.RandomWithPrefix("tf-tests-scw-object-acl-remove")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      objectchecks.IsBucketDestroyed(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_object_bucket" "main" {
						name = "%s"
						region = "%s"
						acl = "authenticated-read"
					}
					`, testBucketName, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.main", true),
					resource.TestCheckResourceAttr("scaleway_object_bucket.main", "acl", "authenticated-read"),
					testAccObjectBucketACLCheck(tt, "scaleway_object_bucket.main", "authenticated-read"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_object_bucket" "main" {
						name = "%s"
						region = "%s"
						acl = "authenticated-read"
					}

					resource "scaleway_object_bucket_acl" "main" {
						bucket = scaleway_object_bucket.main.id
						acl = "public-read"
					}
					`, testBucketName, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.main", true),
					resource.TestCheckResourceAttr("scaleway_object_bucket_acl.main", "bucket", testBucketName),
					resource.TestCheckResourceAttr("scaleway_object_bucket.main", "acl", "authenticated-read"),
					resource.TestCheckResourceAttr("scaleway_object_bucket_acl.main", "acl", "public-read"),
					testAccObjectBucketACLCheck(tt, "scaleway_object_bucket.main", "public-read"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_object_bucket" "main" {
						name = "%s"
						region = "%s"
						acl = "authenticated-read"
					}
					`, testBucketName, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.main", true),
					resource.TestCheckResourceAttr("scaleway_object_bucket.main", "acl", "authenticated-read"),
					testAccObjectBucketACLCheck(tt, "scaleway_object_bucket.main", "private"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "scaleway_object_bucket" "main" {
						name = "%s"
						region = "%s"
					}
					`, testBucketName, objectTestsMainRegion),
				Check: resource.ComposeTestCheckFunc(
					objectchecks.CheckBucketExists(tt, "scaleway_object_bucket.main", true),
					resource.TestCheckResourceAttr("scaleway_object_bucket.main", "acl", "private"),
					testAccObjectBucketACLCheck(tt, "scaleway_object_bucket.main", "private"),
				),
			},
		},
	})
}

func testAccObjectBucketACLCheck(tt *acctest.TestTools, name string, expectedACL string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := context.Background()
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource not found: %s", name)
		}

		bucketRegion := rs.Primary.Attributes["region"]
		s3Client, err := object.NewS3ClientFromMeta(ctx, tt.Meta, bucketRegion)
		if err != nil {
			return err
		}
		if rs.Primary.ID == "" {
			return errors.New("no ID is set")
		}

		bucketName := rs.Primary.Attributes["name"]
		actualACL, err := s3Client.GetBucketAcl(ctx, &s3.GetBucketAclInput{
			Bucket: types.ExpandStringPtr(bucketName),
		})
		if err != nil {
			return fmt.Errorf("could not get ACL for bucket %s: %v", bucketName, err)
		}

		errs := s3ACLAreEqual(expectedACL, actualACL)
		if len(errs) > 0 {
			return fmt.Errorf("unexpected result: %w", errors.Join(errs...))
		}
		return nil
	}
}

func s3ACLAreEqual(expected string, actual *s3.GetBucketAclOutput) (errs []error) {
	ownerID := *object.NormalizeOwnerID(actual.Owner.ID)
	grantsMap := make(map[string]string)
	for _, actualACL := range actual.Grants {
		if actualACL.Permission == "" {
			return append(errs, errors.New("grant has no permission"))
		}
		if actualACL.Grantee.ID != nil {
			grantsMap[string(actualACL.Permission)] = *object.NormalizeOwnerID(actualACL.Grantee.ID)
		} else {
			groupURI := strings.Split(*actualACL.Grantee.URI, "/")
			grantsMap[string(actualACL.Permission)] = groupURI[len(groupURI)-1]
		}
	}

	switch expected {
	case "private":
		if len(grantsMap) != 1 {
			errs = append(errs, fmt.Errorf("expected 1 grant, but got %d", len(grantsMap)))
			return errs
		}
		if grantsMap["FULL_CONTROL"] != ownerID {
			errs = append(errs, fmt.Errorf("expected FULL_CONTROL to be granted to owner (%s), instead got %q", ownerID, grantsMap["FULL_CONTROL"]))
		}

	case "public-read":
		if len(grantsMap) != 2 {
			errs = append(errs, fmt.Errorf("expected 2 grants, but got %d", len(grantsMap)))
			return errs
		}
		if grantsMap["FULL_CONTROL"] != ownerID {
			errs = append(errs, fmt.Errorf("expected FULL_CONTROL to be granted to owner (%s), instead got %q", ownerID, grantsMap["FULL_CONTROL"]))
		}
		if grantsMap["READ"] != s3ACLGranteeAllUsers {
			errs = append(errs, fmt.Errorf("expected READ to be granted to %q, instead got %q", s3ACLGranteeAllUsers, grantsMap["READ"]))
		}

	case "public-read-write":
		if len(grantsMap) != 3 {
			errs = append(errs, fmt.Errorf("expected 3 grants, but got %d", len(grantsMap)))
			return errs
		}
		if grantsMap["FULL_CONTROL"] != ownerID {
			errs = append(errs, fmt.Errorf("expected FULL_CONTROL to be granted to owner (%s), instead got %q", ownerID, grantsMap["FULL_CONTROL"]))
		}
		if grantsMap["READ"] != s3ACLGranteeAllUsers {
			errs = append(errs, fmt.Errorf("expected READ to be granted to %q, instead got %q", s3ACLGranteeAllUsers, grantsMap["READ"]))
		}
		if grantsMap["WRITE"] != s3ACLGranteeAllUsers {
			errs = append(errs, fmt.Errorf("expected WRITE to be granted to %q, instead got %q", s3ACLGranteeAllUsers, grantsMap["WRITE"]))
		}

	case "authenticated-read":
		if len(grantsMap) != 2 {
			errs = append(errs, fmt.Errorf("expected 2 grants, but got %d", len(grantsMap)))
			return errs
		}
		if grantsMap["FULL_CONTROL"] != ownerID {
			errs = append(errs, fmt.Errorf("expected FULL_CONTROL to be granted to owner (%s), instead got %q", ownerID, grantsMap["FULL_CONTROL"]))
		}
		if grantsMap["READ"] != s3ACLGranteeAuthenticatedUsers {
			errs = append(errs, fmt.Errorf("expected READ to be granted to %q, instead got %q", s3ACLGranteeAuthenticatedUsers, grantsMap["READ"]))
		}
	}
	return errs
}
