<div markdown="1" class="text-center">
# Multi-worker cluster
</div>

<div markdown="1" class="text-justify">

This example demonstrates how to use Kubitect to set up a Kubernetes cluster consisting of **one master and three worker nodes**.
The final topology of the deployed Kubernetes cluster is shown in the figure below.

<div class="text-center">
  <img
    class="mobile-w-100"
    src="../../assets/images/topology-1m3w-arch.png"
    alt="Architecture of the cluster with 1 master and 3 worker nodes"
    width="75%">
</div>

!!! note "Note"

    This example skips the explanation of some common configurations such as hosts, network, and node template, as they are already covered in detail in the [Getting started (step-by-step)](../../getting-started/getting-started) guide.

!!! preset "Preset available"

    To export the preset configuration, run:
    <code>
      kubitect export preset <b>example-multi-worker</b>
    </code>

## Step 1: Cluster configuration

You can easily create a cluster with multiple worker nodes by specifying them in the configuration file.
For this example, we have included three worker nodes, but you can add as many as you like to suit your needs.

```yaml title="multi-worker.yaml"
cluster:
  ...
  nodes:
    master:
      instances:
        - id: 1
          ip: 192.168.113.10 # (1)!
    worker:
      instances:
        - id: 1
          ip: 192.168.113.21
        - id: 7
          ip: 192.168.113.27
        - id: 99
```

1.  Static IP address of the node.
    If the `ip` property is omitted, the DHCP lease is requested when the cluster is created.

??? abstract "Final cluster configuration <i class="click-tip"></i>"

    ```yaml title="multi-worker.yaml"
    hosts:
      - name: localhost
        connection:
          type: local

    cluster:
      name: k8s-cluster
      network:
        mode: nat
        cidr: 192.168.113.0/24
      nodeTemplate:
        user: k8s
        updateOnBoot: true
        ssh:
          addToKnownHosts: true
        os:
          distro: ubuntu22
      nodes:
        master:
          instances:
            - id: 1
              ip: 192.168.113.10
        worker:
          instances:
            - id: 1
              ip: 192.168.113.21
            - id: 7
              ip: 192.168.113.27
            - id: 99

    kubernetes:
      version: v1.27.5
      networkPlugin: calico
    ```

## Step 2: Applying the configuration

To deploy a cluster, apply the configuration file:

```sh
kubitect apply --config multi-worker.yaml
```

</div>
