---
subcategory: "Object Storage"
page_title: "Scaleway: scaleway_object_bucket_acl"
---

# Resource: scaleway_object_bucket_acl

Creates and manages Scaleway object storage bucket ACL.
For more information, see [the documentation](https://www.scaleway.com/en/docs/storage/object/concepts/#access-control-list-(acl)).

-> **Note:** `terraform destroy`  does not delete the Object Bucket ACL but does remove the resource from Terraform state.

-> **Note:** [Account identifiers](https://docs.aws.amazon.com/general/latest/gr/acct-identifiers.html) is not supported by scaleway.

## Example Usage

```terraform
resource "scaleway_object_bucket" "some_bucket" {
  name = "some-unique-name"
}

resource "scaleway_object_bucket_acl" "main" {
  bucket = scaleway_object_bucket.main.id
  acl = "private"
}
```

## Example Usage with Grants

```terraform
resource "scaleway_object_bucket" "main" {
    name = "your-bucket"
}

resource "scaleway_object_bucket_acl" "main" {
    bucket = scaleway_object_bucket.main.id
    access_control_policy {
      grant {
        grantee {
            id   = "<project-id>:<project-id>"
            type = "CanonicalUser"
        }
        permission = "FULL_CONTROL"
      }
    
      grant {
        grantee {
          id   = "<project-id>"
          type = "CanonicalUser"
        }
        permission = "WRITE"
      }
    
      owner {
        id = "<project-id>"
      }
    }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) The name of the bucket, or its Terraform ID.
* `acl` - (Optional) The canned ACL you want to apply to the bucket.
* `access_control_policy` - (Optional, Conflicts with acl) A configuration block that sets the ACL permissions for an object per grantee documented below.
* `expected_bucket_owner` - (Optional, Forces new resource) The project ID of the expected bucket owner.
* `region` - (Optional) The [region](https://developers.scaleway.com/en/quickstart/#region-definition) in which the bucket should be created.
* `project_id` - (Defaults to [provider](../index.md#arguments-reference) `project_id`) The ID of the project the bucket is associated with.

~> **Important:** The `project_id` attribute has a particular behavior with s3 products because the s3 API is scoped by project.
If you are using a project different from the default one, you have to specify the `project_id` for every child resource of the bucket,
like bucket ACLs. Otherwise, Terraform will try to create the child resource with the default project ID and you will get a 403 error.


## The ACL

Please check the [canned ACL](https://docs.aws.amazon.com/AmazonS3/latest/userguide/acl_overview.html#canned-acl)

## The Access Control policy

The `access_control_policy` configuration block supports the following arguments:

* `grant` - (Required) Set of grant configuration blocks documented below.
* `owner` - (Required) Configuration block of the bucket owner's display name and ID documented below.

## The Grant

The `grant` configuration block supports the following arguments:

* `grantee` - (Required) Configuration block for the project being granted permissions documented below.
* `permission` - (Required) Logging permissions assigned to the grantee for the bucket.

## The permission

The following list shows each access policy permissions supported.

`READ`, `WRITE`, `READ_ACP`, `WRITE_ACP`, `FULL_CONTROL`

For more information about ACL permissions in the S3 bucket, see [ACL permissions](https://docs.aws.amazon.com/AmazonS3/latest/userguide/acl-overview.html).

## The owner

The `owner` configuration block supports the following arguments:

* `id` - (Required) The ID of the project owner.
* `display_name` - (Optional) The display name of the owner.

## the grantee

The `grantee` configuration block supports the following arguments:

* `id` - (Required) The canonical user ID of the grantee.
* `type` - (Required) Type of grantee. Valid values: CanonicalUser.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The `region`,`bucket` and `acl` separated by (`/`).

## Import

Bucket ACLs can be imported using the `{region}/{bucketName}/{acl}` identifier, e.g.

```bash
$ terraform import scaleway_object_bucket_acl.some_bucket fr-par/some-bucket/private
```

~> **Important:** The `project_id` attribute has a particular behavior with s3 products because the s3 API is scoped by project.
If you are using a project different from the default one, you have to specify the project ID at the end of the import command.

```bash
$ terraform import scaleway_object_bucket_acl.some_bucket fr-par/some-bucket/private@xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxx
```
