# spatial-db

`spatial-db` is a special-purpose database meant to store three-dimensional shapes in voxel worlds

Its original inspiration came from the game Minecraft, but the project expanded to be more general than that. This project was part of my senior Computer Science comprehensive project to graduate from Occidental College

# Requirements

In order to build the project, the requirements are:
* A Go toolchain for at least verion 1.19, but earlier versions may work

# Building the Project

After fufilling all the requirements in the prevoius section, type `go build .` to build the `spatial-db` binary

# Using the Project

The project is available as an importable Go package. A simple quickstart can be found in the following code snippet:

```go
package main

import (
    "github.com/NickyBoy89/spatial-db/server"
    "github.com/NickyBoy89/spatial-db/world"
)

func main() {
    var storageServer server.SimpleServer

    storageServer.SetStorageRoot(".")

    pos := world.BlockPos{X: 0, Y: 0, Z: 0}

    if err := storageServer.ChangeBlockAt(pos, world.Generic); err != nil {
        panic(err)
    }
}
```

# Using the CLI

If the project is built from source using the build instructions, a statically compiled binary is generated. This binary has several commands that are helpful for debugging purposes.

* `spatial-db load worldsave` allows the server to convert pre-existing world saves to a format SpatialDB can use.

# Replicating the Database Results

1. Build the project from source, using the build instructions
2. However, there are a few additional requirements
    1. At least 9GB disk space to hold the world file
    2. At least 12GB RAM to run the in-memory benchmarks, less if only the on-disk implementation
3. The dataset that I chose was the Minecraft map SkyGrid, which can be found on [PlanetMinecraft](https://www.planetminecraft.com/project/skygrid-survival-map/).
    1. Download the resulting zip file and extract those files into a directory
4. Run the built-in conversion script using the compiled binary, with

```bash
./spatial-db load worldsave <path-to-extracted-zip-region-folder> --output "skygrid-save"
```
5. Run the built-in benchmarks using `go test -bench .` to run all tests, or optionally:
    1. `go run -bench=InMemory` to only run SpatialDB's in-memory implementation
    2. `go run -bench=OnDisk` to only run SpatialDB's on disk implementation
    3. `go run -bench=Hashtable` to run a comparison hash-table based implementation
