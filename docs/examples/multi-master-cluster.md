This example shows how to use Kubitect to set up a Kubernetes cluster with **3 master and 3 worker nodes**.
Configuring multiple master nodes provides control plane redundancy, meaning that the control plane can continue to operate normally if a certain number of master nodes fail.
Since Kubitect deploys clusters with a *stacked control plane* (see [Kubernetes.io - Stacked etcd topology](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/ha-topology/#stacked-etcd-topology) for more information), this means that there is no downtime as long as (n/2)+1 master nodes are available.



<div align=center>
  <img
    class="mobile-w-100"
    src="/assets/images/topology-3m3w1lb-arch.png" 
    alt="Architecture of the cluster with 3 master and 3 worker nodes"
    width="75%">
</div>

!!! note "Note"

    In this example we skip the explanation of some common configurations (hosts, network, node template, ...), as they are already explained in the [Getting started (step-by-step)](../../getting-started/getting-started) guide.

## Step 1: Cluster configuration

Each worker node stores only a single control plane IP address. 
Therefore, when creating clusters with multiple master nodes, we need to make sure that all of them are reachable at the same IP address, otherwise all workers would send requests to only one of the master nodes.
This problem can be solved by placing a load balancer in front of the control plane, and instructing it to distribute traffic to all master nodes in the cluster, as shown in the figure below.

<div align=center>
  <img
    class="mobile-w-100"
    src="/assets/images/topology-3m3w1lb-base.png" 
    alt="Scheme of load balancing between control plane nodes"
    width="75%">
</div>

To create such cluster, all we need to do, is specify desired node instances and one load balancer.
Control plane will be accessible on the load balancer's IP address.


```yaml title="multi-master.yaml" 
cluster:
  ...
  nodes:
    loadBalancer:
      instances:
        - id: 1
          ip: 192.168.113.100
    master:
      instances: # (1)!
        - id: 1
          ip: 192.168.113.10
        - id: 2
          ip: 192.168.113.11
        - id: 3
          ip: 192.168.113.12
    worker:
      instances:
        - id: 1
          ip: 192.168.113.20
        - id: 2
          ip: 192.168.113.21
        - id: 3
          ip: 192.168.113.22
```

1. Size of the control plane (number of master nodes) must be odd.

Kubitect detects the load balancer instance in the configuration and installs the *HAProxy* load balancer on an additional virtual machine.
By default, the load balancer is configured to distribute traffic received on port 6443 (Kubernetes API server port) to all specified master nodes.

??? abstract "Final cluster configuration <i class="click-tip"></i>"

    ```yaml title="multi-master.yaml" 
    hosts:
      - name: localhost
        connection:
          type: local

    cluster:
      name: local-k8s-cluster
      network:
        mode: nat
        cidr: 192.168.113.0/24
      nodeTemplate:
        user: k8s
        updateOnBoot: true
        ssh:
          addToKnownHosts: true
        os:
          distro: ubuntu
      nodes:
        loadBalancer:
          instances:
            - id: 1
              ip: 192.168.113.100
        master:
          instances:
            - id: 1
              ip: 192.168.113.10
            - id: 2
              ip: 192.168.113.11
            - id: 3
              ip: 192.168.113.12
        worker:
          instances:
            - id: 1
              ip: 192.168.113.20
            - id: 2
              ip: 192.168.113.21
            - id: 3
              ip: 192.168.113.22

    kubernetes:
      version: v1.23.7
      networkPlugin: calico
      dnsMode: coredns
      kubespray:
        version: v2.19.0
    ```

## Step 2: Applying the configuration

Apply the cluster:
```sh
kubitect apply --config multi-master.yaml
```
