[tag 2.1.0]: https://github.com/MusicDin/kubitect/releases/tag/v2.1.0

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