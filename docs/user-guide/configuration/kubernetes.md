[tag 2.0.0]: https://github.com/MusicDin/kubitect/releases/tag/v2.0.0
[tag 2.2.0]: https://github.com/MusicDin/kubitect/releases/tag/v2.2.0
[tag 3.0.0]: https://github.com/MusicDin/kubitect/releases/tag/v3.0.0
[tag 3.4.0]: https://github.com/MusicDin/kubitect/releases/tag/v3.4.0

<div markdown="1" class="text-center">
# Kubernetes configuration
</div>

<div markdown="1" class="text-justify">

The Kubernetes section of the configuration file contains properties that are specific to Kubernetes, such as the Kubernetes version and network plugin.

## Configuration

### Kubernetes manager

:material-tag-arrow-up-outline: [v3.4.0][tag 3.4.0]
&ensp;
:octicons-file-symlink-file-24: Default: `kubespray`

Specify manager that is used for deploying Kubernetes cluster. Supported values are `kubespray` and `k3s`.

```yaml
kubernetes:
  manager: k3s
```

!!! warning "Warning"

    Support for K3s manager has been added recently, therefore, it may not be fully stable.

### Kubernetes version

:material-tag-arrow-up-outline: [v3.0.0][tag 3.0.0]
&ensp;
:octicons-file-symlink-file-24: Default: `v1.28.6`

By default, the Kubernetes cluster will be deployed using version `v1.33.4`, but you can specify a different version if necessary.

```yaml
kubernetes:
  version: v1.33.5
```

The supported Kubernetes versions include `v1.31`, `v1.32`, and `v1.33`.

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

```yaml
kubernetes:
  networkPlugin: flannel
```

The following table shows the compatibility matrix of supported network plugins and Kubernetes versions:

| Kubernetes Version | Operating system |      Calico      |      Cilium      |      Flannel     |    KubeRouter    |
|--------------------|:-----------------|:----------------:|:----------------:|:----------------:|:----------------:|
| **1.31**           | Ubuntu           | :material-check: | :material-check: | :material-check: | :material-check: |
| **1.31**           | Debian           | :material-check: | :material-check: | :material-check: | :material-check: |
| **1.31**           | CentOS           | :material-check: | :material-check: | :material-check: | :material-check: |
| **1.31**           | RockyLinux       | :material-check: | :material-check: | :material-check: | :material-check: |
| **1.32**           | Ubuntu           | :material-check: | :material-check: | :material-check: | :material-check: |
| **1.32**           | Debian           | :material-check: | :material-check: | :material-check: | :material-check: |
| **1.32**           | CentOS           | :material-check: | :material-check: | :material-check: | :material-check: |
| **1.32**           | RockyLinux       | :material-check: | :material-check: | :material-check: | :material-check: |
| **1.33**           | Ubuntu           | :material-check: | :material-check: | :material-check: | :material-check: |
| **1.33**           | Debian           | :material-check: | :material-check: | :material-check: | :material-check: |
| **1.33**           | CentOS           | :material-check: | :material-check: | :material-check: | :material-check: |
| **1.33**           | RockyLinux       | :material-check: | :material-check: | :material-check: | :material-check: |

!!! note "Note"

    K3s manager currently supports only `flannel` network plugin.

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

### Merge kubeconfig

:material-tag-arrow-up-outline: [v3.4.0][tag 3.4.0]
&ensp;
:octicons-file-symlink-file-24: Default: `false`

Kubitect offers an option to merge the resulting Kubeconfig file with the config on path `~/.kube/config`.
This means that whenever a new cluster is created, it can be selected by context which equals the cluster name.

```yaml
kubernetes:
  other:
    mergeKubeconfig: true
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
