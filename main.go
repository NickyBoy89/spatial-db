package main

import (
	"git.nicholasnovak.io/nnovak/spatial-db/visualization"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := cobra.Command{
		Use: "spatialdb",
	}

	rootCmd.AddCommand(visualization.VisualizeCommand)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
