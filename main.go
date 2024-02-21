package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/pilat/terraform-provider-rustack/rustack_terraform"
)

func main() {
	var debugMode bool = true
	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return rustack_terraform.Provider()
		},
	}

	if debugMode {
		opts.Debug = true
	}

	plugin.Serve(opts)
}
