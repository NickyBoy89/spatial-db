package loading

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/NickyBoy89/spatial-db/storage"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var (
	saveOutputDir string
)

func init() {
	LoadSaveDirCommand.Flags().StringVar(&saveOutputDir, "output", ".", "Where to place the converted save files")
}

var LoadSaveDirCommand = &cobra.Command{
	Use:   "worldsave <save-directory>",
	Short: "Loads all the regions in the specified world's save directory",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// First, fetch a list of all the files that are in the specified directory
		regionFiles, err := os.ReadDir(args[0])
		if err != nil {
			return err
		}

		log.Infof("Loading save directory of %s", args[0])

		u, err := storage.CreateUnityFile(saveOutputDir)
		if err != nil {
			return err
		}
		defer u.Close()
		defer u.WriteMetadataFile(saveOutputDir + ".metadata")

		for regionIndex, regionFile := range regionFiles {
			if regionFile.IsDir() {
				continue
			}

			if strings.HasSuffix(regionFile.Name(), ".mcr") {
				log.Warnf("The file %s is in the MCRegion format, skipping", regionFile.Name())
				continue
			}

			log.Infof("Converting region file %s, [%d/%d]",
				regionFile.Name(),
				regionIndex,
				len(regionFiles),
			)

			filePath := filepath.Join(args[0], regionFile.Name())

			// Load each region file
			chunks, err := LoadRegionFile(filePath)
			if err != nil {
				return err
			}

			// Save each chunk to a separate file
			for _, chunk := range chunks {
				if err := u.WriteChunk(chunk); err != nil {
					return err
				}
			}
		}

		return nil
	},
}
