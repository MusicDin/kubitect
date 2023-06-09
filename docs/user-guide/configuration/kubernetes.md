[tag 2.0.0]: https://github.com/MusicDin/kubitect/releases/tag/v2.0.0
[tag 2.2.0]: https://github.com/MusicDin/kubitect/releases/tag/v2.2.0
[tag 3.0.0]: https://github.com/MusicDin/kubitect/releases/tag/v3.0.0

<div markdown="1" class="text-center">
# Kubernetes configuration
</div>

<div markdown="1" class="text-justify">

The Kubernetes section of the configuration file contains properties that are specific to Kubernetes, such as the Kubernetes version and network plugin.

## Configuration

### Kubernetes version

:material-tag-arrow-up-outline: [v3.0.0][tag 3.0.0]
&ensp;
:octicons-file-symlink-file-24: Default: `v1.25.6`

By default, the Kubernetes cluster will be deployed using version `v1.25.6`, but you can specify a different version if necessary.


```yaml
kubernetes:
  version: v1.24.7
```

The supported Kubernetes versions include `v1.24`, `v1.25`, and `v1.26`.

### Kubernetes network plugin

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]
&ensp;
:octicons-file-symlink-file-24: Default: `calico`

The `calico` network plugin is deployed by default in a Kubernetes cluster.
However, there are multiple supported network plugins available to choose from:

- `calico`
- `cilium`
- `flannel`
- `kube-router`
- `weave`

```yaml
kubernetes:
  networkPlugin: flannel
```

The following table shows the compatibility matrix of supported network plugins and Kubernetes versions:

| Kubernetes Version |      Calico      |      Cilium      |      Flannel     |    KubeRouter    |       Weave      |
|--------------------|:----------------:|:----------------:|:----------------:|:----------------:|:----------------:|
| **1.24**           | :material-check: | :material-check: | :material-check: | :material-check: | :material-check: |
| **1.25**           | :material-check: | :material-check: | :material-check: | :material-check: | :material-check: |
| **1.26**           | :material-check: | :material-check: | :material-check: | :material-check: | :material-check: |

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

Kubitect offers the option to automatically copy the Kubeconfig file to the `~/.kube/config` path.
By default, this feature is disabled to prevent overwriting an existing file.

```yaml
kubernetes:
  other:
    copyKubeconfig: true
```

### Auto renew control plane certificates

:material-tag-arrow-up-outline: [v2.2.0][tag 2.2.0]
&ensp;
:octicons-file-symlink-file-24: Default: `false`

Control plane certificates are renewed every time the cluster is upgraded, and their validity period is one year.
However, in rare cases, clusters that are not upgraded frequently may experience issues.
To address this, you can enable the automatic renewal of control plane certificates on the first Monday of each month by setting the `autoRenewCertificates` property to `true`.

```yaml
kubernetes:
  other:
    autoRenewCertificates: true
```

</div>
