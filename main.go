package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/terraform-providers/terraform-provider-scaleway/scaleway"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debuggable", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	if debugMode {
		err := plugin.Debug(context.Background(), "registry.terraform.io/scaleway/scaleway",
			&plugin.ServeOpts{
				ProviderFunc: scaleway.Provider(scaleway.DefaultProviderConfig()),
			})
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		// Catch every panic after this line. This will send an anonymous report on Scaleway's sentry.
		if scaleway.Version != "develop" {
			defer scaleway.RecoverPanicAndSendReport()
		}

		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: scaleway.Provider(scaleway.DefaultProviderConfig()),
		})
	}
}
