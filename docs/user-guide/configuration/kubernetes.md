[tag 2.0.0]: https://github.com/MusicDin/kubitect/releases/tag/v2.0.0
[tag 2.2.0]: https://github.com/MusicDin/kubitect/releases/tag/v2.2.0

<div markdown="1" class="text-center">
# Kubernetes configuration
</div>

<div markdown="1" class="text-justify">

The Kubernetes section of the configuration file contains properties that are closely related to Kubernetes, such as Kubernetes version and network plugin.

## Configuration

### Kubernetes version

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]
&ensp;
:octicons-file-symlink-file-24: Default: `v1.25.6`

Kubernetes version to be deployed.

```yaml
kubernetes:
  version: v1.24.7
```

### Kubernetes network plugin

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]
&ensp;
:octicons-file-symlink-file-24: Default: `calico`

Kubitect supports multiple Kubernetes network plugins.
Currently, the following network plugins are supported:

  - `calico`
  - `canal`
  - `cilium`
  - `flannel`
  - `kube-router`
  - `weave`

If the network plugin is not set in the Kubitect configuration file, `calico` is used by default.

```yaml
kubernetes:
  networkPlugin: calico
```

### Kubernetes DNS mode

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]
&ensp;
:octicons-file-symlink-file-24: Default: `coredns`

Currently, the only DNS mode supported by Kubitect is `coredns`.
Therefore, it is safe to omit this property.

```yaml
kubernetes:
  dnsMode: coredns
```

### Copy kubeconfig

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]
&ensp;
:octicons-file-symlink-file-24: Default: `false`

Kubitect provides option to automatically copy the Kubeconfig file to `~/.kube/config` path. 
By default, this option is disabled, as it can overwrite an existing file.

```yaml
kubernetes:
  other:
    copyKubeconfig: true
```

### Auto renew control plane certificates

:material-tag-arrow-up-outline: [v2.2.0][tag 2.2.0]
&ensp;
:octicons-file-symlink-file-24: Default: `false`

Control plane certificates are valid for 1 year and are renewed each time the cluster is upgraded.
In some rare cases, this can cause clusters that are not upgraded frequently to stop working properly.
Therefore, the control plane certificates can be renewed automatically on the first Monday of each month by setting the `autoRenewCertificates` property to true.

```yaml
kubernetes:
  other:
    autoRenewCertificates: true
```

</div>
