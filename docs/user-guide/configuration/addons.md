[tag 2.1.0]: https://github.com/MusicDin/kubitect/releases/tag/v2.1.0
[tag 2.2.0]: https://github.com/MusicDin/kubitect/releases/tag/v2.2.0

<div markdown="1" class="text-center">
# Addons
</div>

<div markdown="1" class="text-justify">

## Configuration

### Kubespray addons

:material-tag-arrow-up-outline: [v2.1.0][tag 2.1.0]

Kubespray provides a variety of configurable addons to enhance the functionality of Kubernetes.
Some popular addons include the [Ingress-NGINX controller](https://kubernetes.github.io/ingress-nginx/) and [MetalLB](https://metallb.io/).

Kubespray addons can be configured under the `addons.kubespray` property.
It's important to note that the Kubespray addons are configured in the same as they would be for Kubespray itself, as Kubitect copies the provided configuration into Kubespray's group variables during cluster creation.

The full range of available addons can be explored in the [Kubespray addons sample](https://github.com/kubernetes-sigs/kubespray/blob/master/inventory/sample/group_vars/k8s_cluster/addons.yml), which is available on GitHub.
Most addons are also documented in the [official Kubespray documentation](https://kubespray.io/).

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

[Rook](https://rook.io) is an orchestration tool that integrates [Ceph](https://ceph.io) with Kubernetes.
Ceph is a highly reliable and scalable storage solution, and Rook simplifies its management by automating the deployment, scaling and management of Ceph clusters.

To enable Rook in Kubitect, set `addons.rook.enabled` property to true.

```yaml
addons:
  rook:
    enabled: true
```

Note that Rook is deployed only on worker nodes.
When a cluster is created without worker nodes, Kubitect attempts to install Rook on the master nodes.
In addition to enabling the Rook addon,  **at least one [data disk](../cluster-nodes#data-disks)** must be attached to a node suitable for Rook deployment.
If Kubitect determines that no data disks are available for Rook, it will skip installing Rook.

#### Node selector

The node selector is a dictionary of node labels used to determine which nodes are eligible for Rook deployment.
If a node does not match all of the specified node labels, Rook resources cannot be deployed on that node and disks attached to that node are not used for distributed storage.

```yaml
addons:
  rook:
    nodeSelector:
      rook: true
```

#### Version

By default, Kubitect uses the latest (master) version of Rook.
If you want to use a specific version of Rook, you can set the `addons.rook.version` property to the desired version.

```yaml
addons:
  rook:
    version: v1.11.3
```

</div>
