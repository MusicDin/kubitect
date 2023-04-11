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

# kubitect version v3.0.1
```

## Enable shell autocomplete

!!! tip "Tip"

    To list all supported shells, run: `kubitect completion -h`

    For shell specific instructions run:
    <code>
    kubitect completion <em>shell</em> -h
    </code>

=== "Bash"

    This script depends on the `bash-completion` package.
    If it is not installed already, you can install it via your OS's package manager.

    To load completions in your current shell session:

    ```bash
    source <(kubitect completion bash)
    ```

    To load completions for every new session, execute once:

    **Linux:**

    ```bash
    kubitect completion bash > /etc/bash_completion.d/kubitect
    ```

    **macOS:**

    ```bash
    kubitect completion bash > $(brew --prefix)/etc/bash_completion.d/kubitect
    ```

=== "Zsh"

    If shell completion is not already enabled in your environment you will need to enable it.
    You can execute the following once:

    ```zsh
    echo "autoload -U compinit; compinit" >> ~/.zshrc
    ```

    To load completions in your current shell session:

    ```zsh
    source <(kubitect completion zsh); compdef _kubitect kubitect
    ```

    To load completions for every new session, execute once:

    **Linux:**

    ```zsh
    kubitect completion zsh > "${fpath[1]}/_kubitect"
    ```

    **macOS:**
    ```zsh
    kubitect completion zsh > $(brew --prefix)/share/zsh/site-functions/_kubitect
    ```



</div>
