<div markdown="1" class="text-center">
# Local development
</div>

<div markdown="1" class="text-justify">

This document shows how to build a CLI tool manually and how to use the project without creating any files outside the project's directory.

## Prerequisites

+ [Git](https://git-scm.com/)
+ [Go 1.18](https://go.dev/dl/) or greater

## Step 1: Clone the project

First, clone the project.
```sh
git clone https://github.com/MusicDin/kubitect
```

Afterwards, move into the cloned project.
```sh
cd kubitect
```

## Step 2: Build Kubitect CLI tool

The Kubitect CLI tool can be manually built using Go. 
Running the following command will produce a `kubitect` binary file.
```sh
go build .
```

To make the binary file globally accessible, move it to the `/usr/local/bin/` directory.
```sh
sudo mv kubitect /usr/local/bin/kubitect
```

## Step 3: Local development

By default, Kubitect creates and manages clusters in the Kubitect's home directory (`~/.kubitect`).
However, for development purposes, it is often more convenient to have all resources created in the current directory.

If you want to create a new cluster in the current directory, you can use the `--local` flag when applying the configuration. 
When you create a cluster using the `--local` flag, its name will be prefixed with *local*. 
This prefix is added to prevent any conflicts that might arise when creating new virtual resources.

```sh
kubitect apply --local
```

The resulting cluster will be created in `./.kubitect/clusters/local-<cluster-name>` directory.

</div>
