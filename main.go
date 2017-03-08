package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/terraform-provider-bigip/bigip"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: bigip.Provider,
	})
}
