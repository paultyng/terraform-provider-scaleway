package scaleway

import (
	"fmt"

	"github.com/getsentry/raven-go"
)

const (
	dsn = "https://02b954a634fc439dbb1ac16df9cd3da5@sentry.internal.scaleway.com/223"
)

// RecoverPanicAndSendReport will recover error if any, log them, and send them to sentry.
// It must be called with the defer built-in.
func RecoverPanicAndSendReport() {
	e := recover()
	if e == nil {
		return
	}

	sentryClient, err := newSentryClient()
	if err != nil {
		l.Warningf("cannot create sentry client: %s", err)
	}

	err, isError := e.(error)
	if isError {
		l.Debugf("sending sentry report: %s", sentryClient.CaptureErrorAndWait(err, nil))
	} else {
		l.Debugf("sending sentry report: %s", sentryClient.CaptureErrorAndWait(fmt.Errorf("unknownw error: %v", e), nil))
	}

	panic(e) // lintignore:R009
}

// newSentryClient creates a sentry client with build info tags.
func newSentryClient() (*raven.Client, error) {
	client, err := raven.New(dsn)
	if err != nil {
		return nil, err
	}

	client.SetTagsContext(map[string]string{
		"version": Version,
	})

	return client, nil
}
