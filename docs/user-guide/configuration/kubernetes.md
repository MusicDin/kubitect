[tag 2.0.0]: https://github.com/MusicDin/kubitect/releases/tag/v2.0.0
[tag 2.2.0]: https://github.com/MusicDin/kubitect/releases/tag/v2.2.0

<h1 align="center">Kubernetes configuration</h1>

The Kubernetes section of the configuration file contains properties that are closely related to Kubernetes, such as Kubernetes version, network plugin, and DNS mode. 
In addition, the Kubespray project version and URL can also be specified in this section of the Kubitect configuration.

## Configuration

### Kubespray version

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]
&ensp;
:material-alert-circle-outline: Required

As Kubitect relays on the Kubespray for deploying a Kubernetes cluster, a Kubespray project is cloned during a cluster creation.
This property defines the version of the Kubespray to be cloned.
All Kubespray versions can be found on on their GitHub [release page](https://github.com/kubernetes-sigs/kubespray/releases).

```yaml
kubernetes:
  kubespray:
    version: v2.19.0
```

### Kubespray URL

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]
&ensp;
:octicons-file-symlink-file-24: Default: `https://github.com/kubernetes-sigs/kubespray`

By default, Kubespray is cloned from the official GitHub repository.
If there is a need to use a custom forked version of the project, the url to the repository can be specified with this property.

```yaml
kubernetes:
  kubespray:
    url: https://github.com/kubernetes-sigs/kubespray
```

### Kubernetes version

:material-tag-arrow-up-outline: [v2.0.0][tag 2.0.0]
&ensp;
:material-alert-circle-outline: Required

The Kubernetes version must be defined in the Kubitect configuration.
It must be ensured that Kubespray supports the provided Kubernetes version, otherwise the cluster setup will fail.
In addition, Kubernetes version must be prefixed with the `v`.

```yaml
kubernetes:
  version: v1.23.6
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

Two DNS servers are supproted, `coredns` and `kubedns`.
It is highly recommended to use CoreDNS, which has replaced the kube-dns.
If this property is omitted, CoreDNS is used.

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

## Example usage

### Minimal configuration

The minimalistic Kubernetes configuration encompasses setting Kubernetes and Kubesrpay versions.

```yaml
kuberentes:
  version: v1.23.7
  kubespray:
    version: v2.19.0
```
