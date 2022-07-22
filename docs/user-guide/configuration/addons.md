[tag 2.1.0]: https://github.com/MusicDin/kubitect/releases/tag/v2.1.0

<h1 align="center">Plugins</h1>

## Configuration

### Kubespray plugins

:material-tag-arrow-up-outline: [v2.1.0][tag 2.1.0]

Kubespray provides many useful configurable plugins, such as [Ingress-NGINX controller](https://kubernetes.github.io/ingress-nginx/), [MetalLB](https://metallb.io/), and so on.

Kubespray plugins can be configured in Kubitect under the `plugins.kubespray` property.
Configuration of the Kubespray plugins is exactly the same as the default Kubespray addons configuration, as Kubitect only copies provided configuration into Kubespray's group variables during the cluster creation.

All available Kubespray plugins can be found in the [Kubespray's addons sample](https://github.com/kubernetes-sigs/kubespray/blob/master/inventory/sample/group_vars/k8s_cluster/addons.yml), while most of them are documented in [Kubespray's official documentation](https://kubespray.io/).

```yaml
plugins:
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