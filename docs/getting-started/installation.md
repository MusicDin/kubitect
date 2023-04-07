<div markdown="1" class="text-center">
# Installation
</div>

<div markdown="1" class="text-justify">

## Install Kubitect CLI tool


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

!!! note "Note"
    
    The download URL is a combination of the operating system type, system architecture and version of Kubitect (`https://dl.kubitect.io/<os>/<arch>/<version>`).

    All releases can be found on [GitHub release page](https://github.com/MusicDin/kubitect/releases).

<!--
=== "Release page"

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

        !!! note "Note"
            
            The download URL is a combination of the operating system type, system architecture and version of Kubitect (`https://dl.kubitect.io/<os>/<arch>/<version>`).

            All releases can be found on [GitHub release page](https://github.com/MusicDin/kubitect/releases).

=== "Go packages"

    !!! quote ""

        Install Kubitect from Go packages.

        ```sh
        go install github.com/MusicDin/kubitect/cmd/kubitect@latest
        ```
-->

Verify the installation by checking the Kubitect version.
```sh
kubitect --version

# kubitect version v3.0.0
```

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
