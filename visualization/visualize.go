package visualization

import "github.com/spf13/cobra"

func init() {
	VisualizeCommand.AddCommand(VisualizeChunkCommand)
}

var VisualizeCommand = &cobra.Command{
	Use:     "visualize",
	Aliases: []string{"vis", "show"},
	Short:   "Visualize part of the data spatially",
}
