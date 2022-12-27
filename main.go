package main

import (
	"github.com/mach-composer/mach-composer-plugin-sdk/plugin"

	"github.com/mach-composer/mach-composer-plugin-commercelayer/internal"
)

func main() {
	p := internal.NewCommercelayerPlugin()
	plugin.ServePlugin(p)
}
