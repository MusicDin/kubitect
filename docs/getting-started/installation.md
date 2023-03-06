<div markdown="1" class="text-center">
# Installation
</div>

<div markdown="1" class="text-justify">

Before proceeding with the installation, make sure you meet all of the [requirements](./requirements.md).

## Install Kubitect CLI tool


After all requirements are met, install the Kubitect command line tool.

=== "Release page"

    !!! note "Note"
            
            The download URL is a combination of the operating system type, system architecture and version of Kubitect (`https://dl.kubitect.io/<os>/<arch>/<version>`).

            All releases can be found on [GitHub release page](https://github.com/MusicDin/kubitect/releases).

    !!! quote ""

        Download Kubitect binary file from the release page.
        ```sh
        curl -o kubitect.tar.gz -L https://dl.kubitect.io/linux/amd64/latest
        ```

        Unpack `tar.gz` file.
        ```sh
        tar -xzf kubitect.tar.gz
        ```

        Install Kubitect command line tool by placing the Kubitect binary file in `/usr/local/bin` directory.
        ```sh
        sudo mv kubitect /usr/local/bin/
        ```

=== "Go packages"

    !!! quote ""

        Install from Go packages.

        ```sh
        go install github.com/MusicDin/kubitect/cmd/kubitect@latest
        ```

Verify the installation by checking the Kubitect version.
```sh
kubitect --version

# kubitect version v2.3.0
```

!!! tip "Tip"

    If you are using Kubitect for the first time, we strongly recommend you to take a look at the [getting started](./getting-started.md) tutorial.

<br>

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

    To list all supported shells, run: `kubitect completion -h`

    For shell specific instructions run:
    <code>
    kubitect completion <em>shell</em> -h
    </code>

</div>
