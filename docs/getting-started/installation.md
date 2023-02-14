Before starting with installation, make sure you meet all the [requirements](./requirements.md).

## Install Kubitect CLI tool

After all requirements are met, download the Kubitect command line tool.
```sh
curl -o kubitect.tar.gz -L https://dl.kubitect.io/linux/amd64/latest
```

!!! note "Note"
    
    The download URL is a combination of the operating system type, system architecture and version of Kubitect (`https://dl.kubitect.io/<os>/<arch>/<version>`).

    All releases can be found on [GitHub release page](https://github.com/MusicDin/kubitect/releases).

Unpack `tar.gz` file.
```sh
tar -xzf kubitect.tar.gz
```

Install Kubitect command line tool by placing the Kubitect binary file in `/usr/local/bin` directory.
```sh
sudo mv kubitect /usr/local/bin/
```

Verify the installation by checking the Kubitect version.
```sh
kubitect --version

# kubitect version v2.2.0
```

!!! tip "Tip"

    If you are using Kubitect for the first time, we strongly recommend you to take a look at the [getting started](./getting-started.md) tutorial.

## Enable shell autocomplete

To load completions in your current shell session (`bash`):
```sh
source <(kubitect completion bash)
```

To load completions for every new session, execute once:
```sh
kubitect completion bash > /etc/bash_completion.d/kubitect
```

!!! tip "Tip"

    For all supported shells run: `kubitect completion -h`

    For shell specific instructions run:
    <code>
    kubitect completion <i>shell</i> -h
    </code>