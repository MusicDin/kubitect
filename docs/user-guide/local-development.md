<h1 align="center">Local development</h1>

This tutorial shows how to use the project without creating any files outside the project's directory and how to build a CLI tool manually.

## Clone the project

First, you have to clone the project.
```sh
git clone https://github.com/MusicDin/terraform-kvm-kubespray
```

Afterwards, move into the cloned project.
```sh
cd terraform-kvm-kubespray
```

## Install Kubitect CLI tool

Kubitect CLI tool is implemented in Go using [cobra](https://github.com/spf13/cobra) library.
The tool can either be installed from already built versions available on GitHub or you can build it manually.

To manually build the CLI tool, first change to the `cli` directory.
```sh
cd cli
```

Now, using build the tool using go.
```sh
go build .
```

This will create a `cli` binary file, which can be moved into `/usr/local/bin/` directory to use it globaly.
```sh
sudo mv cli /usr/bin/local/
```

## Local development

By default, Kubitect creates and manages clusters located in the Kubitect home directory (`~/.kubitect`).
Although this pattern is very useful for everyday use, it can be somewhat inconvenient if you are actively making changes to the project, as each change must be committed to the Git repository. 
For this very reason, the Kubitect CLI tool has the `--local` option, which replaces the project's home directory with the path of the current directory.
This way, the source code from the current directory is used to create a cluster and all cluster-related files are created in the current directory.
This option can be used with most actions, such as `apply` or `destroy`.

```sh
kubitect apply --local
```
