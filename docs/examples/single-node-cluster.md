<div markdown="1" class="text-center">
# Single node cluster
</div>

<div markdown="1" class="text-justify">

This example shows how to setup a single node Kubernetes cluster using Kubitect.

<div class="text-center">
  <img
    class="mobile-w-75"
    src="/assets/images/topology-1m-arch.png" 
    alt="Architecture of a single node cluster"
    width="50%">
</div>

!!! note "Note"

    In this example we skip the explanation of some common configurations (hosts, network, node template, ...), as they are already explained in the [Getting started (step-by-step)](../../getting-started/getting-started) guide.

## Step 1: Create the configuration

If you want to initialize a cluster with only one node,
specify a single master node in the cluster configuration file:

```yaml title="single-node.yaml" 
cluster:
  ...
  nodes:
    master:
      instances:
        - id: 1
          ip: 192.168.113.10 # (1)!
```

1.  Static IP address of the node. 
    If the `ip` property is omitted, the DHCP lease is requested when the cluster is created.

??? abstract "Final cluster configuration <i class="click-tip"></i>"

    ```yaml title="single-node.yaml" 
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
        master:
          default:
            ram: 4
            cpu: 2
            mainDiskSize: 32
          instances:
            - id: 1
              ip: 192.168.113.10

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
kubitect apply --config single-node.yaml
```

Your master node now also becomes a worker node.

</div>
