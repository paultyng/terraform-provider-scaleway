package edgeservices

import (
	"errors"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	edgeservices "github.com/scaleway/scaleway-sdk-go/api/edge_services/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/meta"
)

// NewEdgeServicesAPI returns a new edge_services API
func NewEdgeServicesAPI(m interface{}) *edgeservices.API {
	return edgeservices.NewAPI(meta.ExtractScwClient(m))
}

// NewEdgeServicesAPIWithRegion returns a new edge_services API and the region
func NewEdgeServicesAPIWithRegion(d *schema.ResourceData, m interface{}) (*edgeservices.API, scw.Region, error) {
	api := edgeservices.NewAPI(meta.ExtractScwClient(m))

	region, err := meta.ExtractRegion(d, m)
	if err != nil {
		return nil, "", err
	}
	return api, region, err
}

func isStageUsedInPipelineError(err error) bool {
	if err == nil {
		return false
	}

	responseError := &scw.ResponseError{}
	if errors.As(err, &responseError) && responseError.StatusCode == http.StatusBadRequest {
		if strings.Contains(responseError.Message, "operation was rejected because the stage stage is used in a pipeline") {
			return true
		}
	}
	return false
}
