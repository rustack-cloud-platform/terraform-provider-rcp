package main

import (
	"context"
	"flag"
	"log"

	// "github.com/pilat/rustack_terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/pilat/terraform-provider-rustack/rustack_terraform"
)

func main() {
	// plugin.Serve(&plugin.ServeOpts{
	// 	ProviderFunc: rustack_terraform.Provider})

	var debugMode bool
	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return rustack_terraform.Provider()
		},
	}

	if debugMode {
		err := plugin.Debug(context.Background(), "pilat/rustack", opts)
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}

	plugin.Serve(opts)
}
