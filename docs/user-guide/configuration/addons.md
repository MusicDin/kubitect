[tag 2.1.0]: https://github.com/MusicDin/kubitect/releases/tag/v2.1.0
[tag 2.2.0]: https://github.com/MusicDin/kubitect/releases/tag/v2.2.0

<h1 align="center">Addons</h1>

## Configuration

### Kubespray addons

:material-tag-arrow-up-outline: [v2.1.0][tag 2.1.0]

Kubespray offers many useful configurable addons, such as the [Ingress-NGINX controller](https://kubernetes.github.io/ingress-nginx/), [MetalLB](https://metallb.io/), and so on.

Kubespray addons can be configured in Kubitect under the `addons.kubespray` property.
The configuration of Kubespray addons is exactly the same as the default configuration of Kubespray addons, since Kubitect simply copies the provided configuration into Kubespray's group variables when the cluster is created.

All available Kubespray addons can be found in the [Kubespray addons sample](https://github.com/kubernetes-sigs/kubespray/blob/master/inventory/sample/group_vars/k8s_cluster/addons.yml), while most of them are documented in the [official Kubespray documentation](https://kubespray.io/).

```yaml
addons:
  kubespray:

    # Nginx ingress controller deployment
    ingress_nginx_enabled: true
    ingress_nginx_namespace: "ingress-nginx"
    ingress_nginx_insecure_port: 80
    ingress_nginx_secure_port: 443

    # MetalLB deployment
    metallb_enabled: true
    metallb_speaker_enabled: true
    metallb_ip_range:
      - "10.10.9.201-10.10.9.254"
    metallb_pool_name: "default"
    metallb_auto_assign: true
    metallb_version: v0.12.1
    metallb_protocol: "layer2"
```

### Rook addon

:material-tag-arrow-up-outline: [v2.2.0][tag 2.2.0]
&ensp;
:material-flask-outline: **Experimental**

[Rook](https://rook.io) is an orchestration tool that allows [Ceph](https://ceph.io), a reliable and scalable storage, to run within a Kubernetes cluster.

In Kubitect, Rook can be enabled by simply setting `addons.rook.enabled` to true.

```yaml
addons:
  rook:
    enabled: true
```

By default, Rook uses all available data disks attached to worker nodes and converts them to distributed storage.
Similarly, all worker nodes are used for Rook deployment.
To restrict on which nodes Rook resources can be deployed or which disks can be used by Rook, the Node Selector can be used together with the Data Disk Selector.

#### Node selector

The node selector is a dictionary of node labels used to determine which nodes are eligible for Rook deployment.
If a node does not match all of the specified node labels, Rook resources cannot be deployed on that node and disks attached to that node are not used for distributed storage.

```yaml
addons:
  rook:
    nodeSelector:
      rook: true
```

#### Data disk selector

!!! Warning "Warning"

    The data disk selector is evaluated during the initial setup of the cluster.
    This prevents additional data disks that are later attached to the virtual machines from being detected by Rook.
    Therefore, changes to the data disk scheme after the initial setup of the cluster can result in unexpected behavior.

In addition to the node selector, Rook can be restricted to use only specific data disks attached to VMs.
The data disk selector is a list of strings, more specifically disk names, that instructs Rook to use only disks with those names.

```yaml
addons:
  rook:
    dataDiskSelector:
      - rook
      - rook-disk
```
