<h1 align="center">Installation</h1>

Before starting with installation, make sure you meet all the [requirements](./requirements.md).

## Install Kubitect CLI tool

After all requirements are met, download the Kubitect command line tool.
```sh
curl -o kubitect.tar.gz -L https://github.com/MusicDin/kubitect/releases/download/v2.0.0/kubitect-cli-v2.0.0-linux-amd64.tar.gz
```

Unpack `tar.gz` file.
```sh
tar -xzf kubitect.tar.gz
```

Install Kubitect command line tool by placing the Kubitect binary file in `/usr/local/bin` directory.
```sh
sudo mv kubitect-cli /usr/local/bin/kubitect
```

Verify the installation by checking the Kubitect version.
```sh
kubitect --version

# kubitect version 2.0.0
```

!!! tip "Tip"

    If you are using Kubitect for the first time, we strongly recommend you to take a look at the [getting started](./getting-started.md) tutorial.

## Enable shell autocomplete

For example, to enable automplete for `bash`, run the following command.
```sh
echo 'source <(kubitect completion bash)' >> ~/.bashrc
```

Then reload your shell.