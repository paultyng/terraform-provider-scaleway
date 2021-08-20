# Testing the module

## Unit testing

You can test the provider, by running `make test`.

```sh
make test
```

## Acceptance testing

Acceptance test are made to test the terraform module with real API calls so they will create real resources that will be invoiced.
But in order to run faster tests and avoid bad surprise at the end of the month we shipped the project with mocks (recorded with [go-vcr](https://github.com/dnaeon/go-vcr)).

### Running the acceptance tests with mocks

By default, mocks are used during acceptance tests.

```sh
make testacc
```

### Running the acceptance tests on real resources

:warning: This will cost money.

```sh
export TF_UPDATE_CASSETTES=true
make testacc
```

It's also required to have Scaleway environment vars available:

```sh
export SCW_ACCESS_KEY=SCWXXXXXXXXXXXXXXXXX
export SCW_SECRET_KEY=XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX
export SCW_DEFAULT_PROJECT_ID=XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX
```

For testing the domain API, it will use the first available domain in your domains list. You need to have a valid domain.

You can force the test domain with an environment var:

```sh
export TF_TEST_DOMAIN=your-domain.tld
```

To ease debugging you can also set:
```sh
export TF_LOG=DEBUG
export SCW_DEBUG=1
```

Running a single test:
```sh
TF_UPDATE_CASSETTES=true TF_LOG=DEBUG SCW_DEBUG=1 TF_ACC=1 go test ./scaleway -v -run=TestAccScalewayDataSourceRDBInstance_Basic -timeout=120m -parallel=10
```
