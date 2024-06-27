package edgeservices

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scaleway/scaleway-sdk-go/api/edge_services/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/locality"
	"github.com/scaleway/terraform-provider-scaleway/v2/internal/types"
)

func expandEdgeServicesScalewayS3BackendConfig(raw interface{}) *edge_services.ScalewayS3BackendConfig {
	if raw == nil || len(raw.([]interface{})) != 1 {
		return nil
	}
	rawMap := raw.([]interface{})[0].(map[string]interface{})
	return &edge_services.ScalewayS3BackendConfig{
		BucketName:   types.ExpandStringPtr(rawMap["bucket_name"].(string)),
		BucketRegion: types.ExpandStringPtr(rawMap["bucket_region"].(string)),
		IsWebsite:    types.ExpandBoolPtr(rawMap["is_website"]),
	}
}

func flattenEdgeServicesScalewayS3BackendConfig(s3backend *edge_services.ScalewayS3BackendConfig) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"bucket_name":   types.FlattenStringPtr(s3backend.BucketName),
			"bucket_region": types.FlattenStringPtr(s3backend.BucketRegion),
			"is_website":    types.FlattenBoolPtr(s3backend.IsWebsite),
		},
	}
}

func expandEdgeServicesPurge(raw interface{}) []*edge_services.PurgeRequest {
	if raw == nil {
		return nil
	}

	purgeRequests := []*edge_services.PurgeRequest(nil)
	for _, pr := range raw.(*schema.Set).List() {
		rawPr := pr.(map[string]interface{})
		purgeRequest := &edge_services.PurgeRequest{}
		purgeRequest.PipelineID = rawPr["pipeline_id"].(string)
		purgeRequest.Assets = types.ExpandStringsPtr(rawPr["assets"])
		purgeRequest.All = types.ExpandBoolPtr(rawPr["all"])

		purgeRequests = append(purgeRequests, purgeRequest)
	}
	return purgeRequests
}

func expandEdgeServicesTLSSecrets(raw interface{}, region scw.Region) []*edge_services.TLSSecret {
	secrets := []*edge_services.TLSSecret(nil)
	rawSecrets := raw.([]interface{})
	for _, rawSecret := range rawSecrets {
		mapSecret := rawSecret.(map[string]interface{})
		secret := &edge_services.TLSSecret{
			SecretID: locality.ExpandID(mapSecret["secret_id"]),
			Region:   region,
		}
		secrets = append(secrets, secret)
	}
	return secrets
}

func flattenEdgeServicesTLSSecrets(secrets []*edge_services.TLSSecret) interface{} {
	if len(secrets) == 0 || secrets == nil {
		return nil
	}

	secretsI := []map[string]interface{}(nil)
	for _, secret := range secrets {
		secretMap := map[string]interface{}{
			"secret_id": secret.SecretID,
			"region":    secret.Region.String(),
		}
		secretsI = append(secretsI, secretMap)
	}
	return secretsI
}

func wrapSecretsInConfig(secrets []*edge_services.TLSSecret) *edge_services.TLSSecretsConfig {
	return &edge_services.TLSSecretsConfig{
		TLSSecrets: secrets,
	}
}
