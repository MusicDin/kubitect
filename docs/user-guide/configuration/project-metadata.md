[tag 2.0.0]: https://github.com/MusicDin/kubitect/releases/tag/v2.0.0

<h1 align="center">Project metadata</h1>

The project metadata contains properties that are closely related to the project itself.
These properties are rather fine-grained settings, as the default values should cover most use cases.

## Configuration

### Project URL

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]
&ensp;
:octicons-file-symlink-file-24: Default: `https://github.com/MusicDin/kubitect`

By default, Kubitect clones the source code of the project from the official [GitHub repository](https://github.com/MusicDin/kubitect).
To use a custom repository instead, set the `kubitect.url` property in the configuration file to the desired url.

```yaml
kubitect:
    url: "https://github.com/user/repo"
```

### Project version

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]
&ensp;
:octicons-file-symlink-file-24: Default: *CLI tool version*

If project version is not specified in the configuration file, version of the CLI tool will be used by default (*recommended*).
To set the specific project version, set `kubitect.version` property in the configuration file.

```yaml
kubitect:
    version: v2.2.0 # (1)!
```

1. Version can be either a tag (`v2.2.0`) or a branch name (`master`).

All Kubitect versions can be found on the [release page](https://github.com/MusicDin/kubitect/releases).

!!! warning "Warning"

    A mismatched CLI and project version can lead to unexpected behavior.