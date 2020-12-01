# Load balancing traffic to ingress controller

This example will show how to set up load balancer to route all traffic coming on LB's port `80` (HTTP) to ingress controller.

[Kubernetes nginx ingress controller](https://github.com/kubernetes/ingress-nginx) is used in this example.

**IMPORTANT:** This example shows neither best practices nor secure configuration! It's meant to be simplistic and straightforward.

## Configuring load balancer

In Kubernetes service of type `NodePort` can by default take values from 30000-32767,
we will for the sake of this example expose ingress controller on port `30080` for HTTP and `30443` for HTTPS, though we won't cover HTTPS configuration.

Now that we know on which port ingress controller will listen for HTTP traffic, we can configure our HAProxy configuration.

Place the following code where comment `Place custom configurations here` is located in [HAProxy template file](../../templates/haproxy.tpl):
```bash

# Load balancing to ingress controller configuration #

frontend ingress-controller
        # Floating IP (192.168.113.200) will be probably diffrent in your configuration.
        bind            192.168.113.200:80
        mode            tcp
        option          tcplog
        default_backend ingress-controller

backend ingress-controller
        mode tcp
        balance roundrobin
        # Route traffic to your master nodes on ingress controller port
        server          k8s-master-0 192.168.113.10:30080
        server          k8s-master-1 192.168.113.11:30080
        server          k8s-master-2 192.168.113.12:30080
```

Now it's time to **initialize your cluster**.

*Note: If the cluster is already initialized, you can also [configure LBs by SSH-ing](../load-balancer.md#modifying-load-balancers-configuration-over-ssh) into each of them and applying this changes.*

## Installing ingress-controller

Now we are going to install ingress controller and expose it as [NodePort](https://kubernetes.io/docs/concepts/services-networking/service/#nodeport).
We will expose it on port `30080` for HTTP traffic (*we will also expose it on port `30443` for HTTPS*).

Modify a version of ingress controller and download YAML configuration:
<pre>
wget https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v<b>0.35.0</b>/deploy/static/provider/baremetal/deploy.yaml -O ingress-controller-deployment.yaml
</pre>

Edit downloaded `ingress-controller-deployment.yaml` file:
```yaml
...
# Find this section of the code by searching for "NodePort"
spec:
  type: NodePort
  ports:
    - name: http
      port: 80
      protocol: TCP
      targetPort: http
      # Add custom node port for HTTP
      nodePort: 30080
    - name: https
      port: 443
      protocol: TCP
      targetPort: https
      # Add custom node port for HTTPS
      nodePort: 30443
...
```

Apply ingress controller configuration:
```bash
kubectl apply -f ingress-controller-deployment.yaml
```

That's it. Now all HTTP traffic will be load balanced to master nodes on ingress controller's port.