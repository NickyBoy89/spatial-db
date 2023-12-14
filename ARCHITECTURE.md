# Code Architecture

The project is broken up into several folders that make up modules

* `loading`
* `server`
* `storage`
* `world`

Any other folders are helper code for running tests, or to aid in development, and do not directly connect to the database implementation.

## Loading

This module contains all the functionality to load existing Minecraft world saves into SpatialDB files, so that they can be used for benchmarking. This is accessed thorough the `spatialdb` command-line utility, and is documented more completely in the README.

## Server

This module contains several implementations for the different servers that I analyze in my project. Currently, this contains the following:

* `inmemory_server.go` contains a reference in-memory server
* `simple_server.go` contains the main implementation of the disk-backed server that most of my project focuses on
* `hashserver.go` contains a simple hash-map backed server implementation

## Storage

This module contains all the logic that handles file operations within the database, and all logic from reading into chunk data files

* `unity_file.go` contains the implementation of unity files

## World

This module contains all the data structures that model the chunk system, as well as the voxel grid.
