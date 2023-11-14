package main

import (
	"git.nicholasnovak.io/nnovak/spatial-db/connector"
	"git.nicholasnovak.io/nnovak/spatial-db/loading"
	"git.nicholasnovak.io/nnovak/spatial-db/visualization"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := cobra.Command{
		Use: "spatialdb",
	}

	rootCmd.AddCommand(visualization.VisualizeCommand)
	rootCmd.AddCommand(connector.ProxyPortCommand)
	rootCmd.AddCommand(loading.LoadRegionFileCommand)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
