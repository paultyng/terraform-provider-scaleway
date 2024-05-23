package instance

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	block "github.com/scaleway/scaleway-sdk-go/api/block/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality/zonal"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/meta"
)

type BlockAndInstanceAPI struct {
	*instance.API
	blockAPI *block.API
}

type GetUnknownVolumeRequest struct {
	VolumeID string
	Zone     scw.Zone
}

type UnknownVolume struct {
	Zone     scw.Zone
	ID       string
	Name     string
	Size     scw.Size
	ServerID *string

	// IsBlockVolume is true if volume is managed by block API
	IsBlockVolume bool

	InstanceVolumeType instance.VolumeVolumeType
}

// VolumeTemplateUpdate return a VolumeServerTemplate for an UpdateServer request
func (volume *UnknownVolume) VolumeTemplateUpdate() *instance.VolumeServerTemplate {
	template := &instance.VolumeServerTemplate{
		ID:   scw.StringPtr(volume.ID),
		Name: &volume.Name, // name is ignored by the API, any name will work here
	}
	if volume.IsBlockVolume {
		template.Name = nil
		template.VolumeType = volume.InstanceVolumeType
	}

	return template
}

// IsLocal returns true if the volume is a local volume
func (volume *UnknownVolume) IsLocal() bool {
	return !volume.IsBlockVolume && volume.InstanceVolumeType == instance.VolumeVolumeTypeLSSD
}

// IsAttached returns true if the volume is attached to a server
func (volume *UnknownVolume) IsAttached() bool {
	return volume.ServerID != nil && *volume.ServerID != ""
}

func (api *BlockAndInstanceAPI) GetUnknownVolume(req *GetUnknownVolumeRequest, opts ...scw.RequestOption) (*UnknownVolume, error) {
	getVolumeResponse, err := api.API.GetVolume(&instance.GetVolumeRequest{
		Zone:     req.Zone,
		VolumeID: req.VolumeID,
	}, opts...)
	notFoundErr := &scw.ResourceNotFoundError{}
	if err != nil && !errors.As(err, &notFoundErr) {
		return nil, err
	}

	if getVolumeResponse != nil {
		vol := &UnknownVolume{
			Zone:               getVolumeResponse.Volume.Zone,
			ID:                 getVolumeResponse.Volume.ID,
			Name:               getVolumeResponse.Volume.Name,
			Size:               getVolumeResponse.Volume.Size,
			IsBlockVolume:      false,
			InstanceVolumeType: getVolumeResponse.Volume.VolumeType,
		}
		if getVolumeResponse.Volume.Server != nil {
			vol.ServerID = &getVolumeResponse.Volume.Server.ID
		}

		return vol, nil
	}

	blockVolume, err := api.blockAPI.GetVolume(&block.GetVolumeRequest{
		Zone:     req.Zone,
		VolumeID: req.VolumeID,
	}, opts...)
	if err != nil {
		return nil, err
	}

	vol := &UnknownVolume{
		Zone:               blockVolume.Zone,
		ID:                 blockVolume.ID,
		Name:               blockVolume.Name,
		Size:               blockVolume.Size,
		IsBlockVolume:      true,
		InstanceVolumeType: instance.VolumeVolumeTypeSbsVolume,
	}
	for _, ref := range blockVolume.References {
		if ref.ProductResourceType == "instance_server" {
			vol.ServerID = &ref.ProductResourceID
		}
	}

	return vol, nil
}

// newAPIWithZone returns a new instance API and the zone for a Create request
func instanceAndBlockAPIWithZone(d *schema.ResourceData, m interface{}) (*BlockAndInstanceAPI, scw.Zone, error) {
	instanceAPI := instance.NewAPI(meta.ExtractScwClient(m))
	blockAPI := block.NewAPI(meta.ExtractScwClient(m))

	zone, err := meta.ExtractZone(d, m)
	if err != nil {
		return nil, "", err
	}

	return &BlockAndInstanceAPI{
		API:      instanceAPI,
		blockAPI: blockAPI,
	}, zone, nil
}

// NewAPIWithZoneAndID returns an instance API with zone and ID extracted from the state
func instanceAndBlockAPIWithZoneAndID(m interface{}, zonedID string) (*BlockAndInstanceAPI, scw.Zone, string, error) {
	instanceAPI := instance.NewAPI(meta.ExtractScwClient(m))
	blockAPI := block.NewAPI(meta.ExtractScwClient(m))

	zone, ID, err := zonal.ParseID(zonedID)
	if err != nil {
		return nil, "", "", err
	}

	return &BlockAndInstanceAPI{
		API:      instanceAPI,
		blockAPI: blockAPI,
	}, zone, ID, nil
}
