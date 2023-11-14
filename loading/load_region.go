package loading

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"git.nicholasnovak.io/nnovak/spatial-db/world"
	"github.com/Tnze/go-mc/save"
	"github.com/Tnze/go-mc/save/region"
	"github.com/spf13/cobra"
)

type SectorPos struct {
	x, y int
}

var outputDir string

func init() {
	LoadRegionFileCommand.Flags().StringVar(&outputDir, "output", ".", "The output directory for the files")
}

var LoadRegionFileCommand = &cobra.Command{
	Use: "load",
	RunE: func(cmd *cobra.Command, args []string) error {
		regionFile, err := region.Open(args[0])
		if err != nil {
			return err
		}
		defer regionFile.Close()

		validSectors := []SectorPos{}

		for i := 0; i < 32; i++ {
			for j := 0; j < 32; j++ {
				if regionFile.ExistSector(i, j) {
					validSectors = append(validSectors, SectorPos{i, j})
				}
			}
		}

		for _, sectorPos := range validSectors {
			data, err := regionFile.ReadSector(sectorPos.x, sectorPos.y)
			if err != nil {
				return err
			}

			var chunk save.Chunk
			if err := chunk.Load(data); err != nil {
				return err
			}

			var chunkData world.ChunkData
			chunkData.FromMCAChunk(chunk)

			outfile, err := os.OpenFile(filepath.Join(outputDir, fmt.Sprintf("p.%d.%d.chunk", chunkData.Pos.X, chunkData.Pos.Z)), os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				log.Fatal(err)
			}
			if err := json.NewEncoder(outfile).Encode(chunkData); err != nil {
				log.Fatal(err)
			}
			outfile.Close()
		}

		return nil
	},
}
