package scaleway

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceScalewayAccountSSKKey() *schema.Resource {
	return ResourceScalewayIamSSKKey()
}
