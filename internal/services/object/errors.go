package object

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// Error code constants missing from AWS Go SDK:
// https://docs.aws.amazon.com/sdk-for-go/api/service/s3/#pkg-constants

const (
	// ErrCodeNoSuchTagSet tag(s) not found
	ErrCodeNoSuchTagSet = "NoSuchTagSet"
	// ErrCodeNoSuchCORSConfiguration CORS configuration not found
	ErrCodeNoSuchCORSConfiguration = "NoSuchCORSConfiguration"
	// ErrCodeNoSuchLifecycleConfiguration lifecycle configuration rule not found
	ErrCodeNoSuchLifecycleConfiguration = "NoSuchLifecycleConfiguration"
	// ErrCodeAccessDenied action on resource is denied
	ErrCodeAccessDenied = "AccessDenied"
	// ErrCodeBucketNotEmpty bucket is not empty
	ErrCodeBucketNotEmpty = "BucketNotEmpty"
	// ErrCodeNoSuchBucket bucket not found
	ErrCodeNoSuchBucket = "NoSuchBucket"
	// ErrCodeNoSuchBucketPolicy policy not found
	ErrCodeNoSuchBucketPolicy = "NoSuchBucketPolicy"
	// ErrCodeNoSuchWebsiteConfiguration website configuration not found
	ErrCodeNoSuchWebsiteConfiguration = "NoSuchWebsiteConfiguration"
	// ErrCodeObjectLockConfigurationNotFoundError object lock configuration not found
	ErrCodeObjectLockConfigurationNotFoundError = "ObjectLockConfigurationNotFoundError"
	// ErrCodeAuthorizationError authorization error
	ErrCodeAuthorizationError = "AuthorizationError"
	// ErrCodeInternalException internal exception
	ErrCodeInternalException = "InternalException"
	// ErrCodeInternalServiceError internal exception error
	ErrCodeInternalServiceError = "InternalServiceError"
	// ErrCodeInvalidAction invalid action
	ErrCodeInvalidAction = "InvalidAction"
	// ErrCodeInvalidParameterException invalid parameter exception
	ErrCodeInvalidParameterException = "InvalidParameterException"
	// ErrCodeInvalidParameterValue invalid parameter value
	ErrCodeInvalidParameterValue = "InvalidParameterValue"
	// ErrCodeInvalidRequest invalid request
	ErrCodeInvalidRequest = "InvalidRequest"
	// ErrCodeOperationDisabledException operation disabled exception
	ErrCodeOperationDisabledException = "OperationDisabledException"
	// ErrCodeOperationNotPermitted operation not permitted
	ErrCodeOperationNotPermitted = "OperationNotPermitted"
	// ErrCodeUnknownOperationException   unknown operation exception
	ErrCodeUnknownOperationException = "UnknownOperationException"
	// ErrCodeUnsupportedFeatureException = unsupported Feature exception
	ErrCodeUnsupportedFeatureException = "UnsupportedFeatureException"
	// ErrCodeUnsupportedOperation unsupported operation
	ErrCodeUnsupportedOperation = "UnsupportedOperation"
	// ErrCodeValidationError validation error
	ErrCodeValidationError = "ValidationError"
	// ErrCodeValidationException validation exception
	ErrCodeValidationException = "ValidationException"
)

// TimedOut returns true if the error represents a "wait timed out" condition.
// Specifically, TimedOut returns true if the error matches all these conditions:
//   - err is of type resource.TimeoutError
//   - TimeoutError.LastError is nil
func TimedOut(err error) bool {
	// This explicitly does *not* match wrapped TimeoutErrors
	timeoutErr, ok := err.(*retry.TimeoutError) //nolint:errorlint // Explicitly does *not* match wrapped TimeoutErrors
	return ok && timeoutErr.LastError == nil
}

// ErrCodeEquals returns true if the error matches all these conditions:
//   - err is of type scw.Error
//   - Error.Error() equals one of the passed codes
func ErrCodeEquals(err error, codes ...string) bool {
	var scwErr scw.SdkError
	if errors.As(err, &scwErr) {
		for _, code := range codes {
			if scwErr.Error() == code {
				return true
			}
		}
	}
	return false
}

type ServiceErrorCheckFunc func(*testing.T) resource.ErrorCheckFunc

var serviceErrorCheckFunc map[string]ServiceErrorCheckFunc

func ErrorCheck(t *testing.T, endpointIDs ...string) resource.ErrorCheckFunc {
	t.Helper()
	return func(err error) error {
		if err == nil {
			return nil
		}

		for _, endpointID := range endpointIDs {
			if f, ok := serviceErrorCheckFunc[endpointID]; ok {
				ef := f(t)
				err = ef(err)
			}

			if err == nil {
				break
			}
		}

		return err
	}
}
