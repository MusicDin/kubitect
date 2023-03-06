/*
Package kubitect is a CLI tool that helps you manage multiple Kubernetes clusters.

# Installation

A valid installation of Go 1.18 or greater is required. The following example installs the latest stable version of
the kubitect cli. Replace latest with a specific version tag to install other versions.

	go install github.com/MusicDin/kubitect/cmd/kubitect@latest

You can also download the binary from the https://github.com/MusicDin/kubitect/releases page and add
it to the $PATH environment variable (or folder). Each release contains the binary for all the supported platforms.

# Usage

After installation the `kubitect` command should be available for usage.

example:

	kubitect --help

Output:

	Kubitect is a CLI tool that helps you manage multiple Kubernetes clusters.

	Usage:

		kubitect [command]

	Cluster Management Commands:

		apply       Create, scale or upgrade the cluster
		destroy     Destroy the cluster

	Support Commands:

		export      Export specific configuration file
		list        Lists clusters

	Other Commands:

		completion  Generate the autocompletion script for the specified shell
		help        Help about any command

	Flags:

		-h, --help      help for kubitect
		-v, --version   version for kubitect

	Use "kubitect [command] --help" for more information about a command.
*/
package main
