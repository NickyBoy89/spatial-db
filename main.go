package main

import (
	"github.com/NickyBoy89/spatial-db/connector"
	"github.com/NickyBoy89/spatial-db/loading"
	"github.com/NickyBoy89/spatial-db/visualization"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := cobra.Command{
		Use: "spatialdb",
	}

	rootCmd.AddCommand(visualization.VisualizeCommand)
	rootCmd.AddCommand(connector.ProxyPortCommand)

	loadCmd := &cobra.Command{
		Use:   "load",
		Short: "Loads save files into the database's format",
	}

	loadCmd.AddCommand(loading.LoadSaveDirCommand)

	rootCmd.AddCommand(loadCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
