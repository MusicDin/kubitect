[tag 2.0.0]: https://github.com/MusicDin/kubitect/releases/tag/v2.0.0

To ensure the configuration will reproduce exactly the same cluster when reapplied, it is recommended to set specific Kubitect version.
Furthermore if the project is forked or cloned, you can also instruct Kubitect to use your repository instead of the original one.

## Configuration

### Project version

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]
&ensp;
:octicons-file-symlink-file-24: Default: *CLI tool version*

If project version is not specified in the configuration file, version of the CLI tool will be used by default.
To set the specific project version, set `kubitect.version` property in the configuration file.

```yaml
kubitect:
    version: v2.1.0
```

All Kubitect versions can be found on the [release page](https://github.com/MusicDin/kubitect/releases).

!!! warning "Warning"

    It is not recommended to use version lower then the CLI tool version.

    To check the CLI tool version, run the following command.

    ```sh
    kubitect --version
    ```

### Project URL

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]
&ensp;
:octicons-file-symlink-file-24: Default: `https://github.com/MusicDin/kubitect`

By default, Kubitect clones the source code of the project from the official [GitHub repository](https://github.com/MusicDin/kubitect).
To use the custom repository instead, set the `kubitect.url` property in the configuration file to the desired url.

```yaml
kubitect:
    url: "https://github.com/user/repo"
```
