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

By default, Kubitect creates and manages clusters in the Kubitect's home directory (`~/.kubitect`) and pulls the required source code from the official git repository.
However, this approach can be inconvenient for active development, as all changes must be pushed to the git repository before they can be used.

To address this, the Kubitect CLI tool includes a `--local` option that can be used with most commands, such as `apply` or `destroy`. 
When the `--local` option is used, a cluster is created within the current directory, and the source code from the current directory is used instead of being pulled from the remote repository.

```sh
kubitect apply --local
```

The resulting cluster will be created in `./.kubitect/clusters/local-<cluster-name>` directory.

</div>
