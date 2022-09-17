<div markdown="1" class="text-center">
# Rook cluster
</div>

<div markdown="1" class="text-justify">

!!! warning "Important"

    Since the Rook addon is still under development, it may not work as expected.
    Therefore, any feedback would be greatly appreciated.

This example shows how to use Kubitect to set up **distributed storage with Rook**.
For distributed storage, we add an additional data disk to each virtual machine as shown on the figure below.

<div class="text-center">
  <img
    class="mobile-w-100"
    src="/assets/images/rook-cluster-arch.png" 
    alt="Basic Rook cluster scheme"
    width="75%">
</div>

## Basic setup

### Step 1: Define data resource pool

To configure distributed storage with Rook, the data disks must be attached to the virtual machines.
By default, each data disk is created in a main resource pool.
Optionally, you can configure additional resource pools and associate data disks with them later.

In this example, we define an additional resource pool named 'rook-pool'.
```yaml title="rook-sample.yaml"
hosts:
  - name: localhost
    connection:
      type: local
    dataResourcePools:
      - name: rook-pool
```

### Step 2: Attach data disks

After the data resource pool is configured, we are ready to allocate some data disks to the virtual machines.

```yaml title="rook-sample.yaml"
cluster:
  nodes:
    worker:
      instances:
        - id: 1
          dataDisks:
            - name: rook
              pool: rook-pool # (1)!
              size: 256
        - id: 2
          dataDisks:
            - name: rook
              pool: rook-pool
              size: 256
        - id: 3
        - id: 4
          dataDisks:
            - name: rook
              pool: rook-pool
              size: 256
            - name: test
              pool: rook-pool
              size: 32
```

1. To create data disks in the main resource pool, either omit the pool property or set its value to `main`.

### Step 3: Enable Rook addon

Once the disks are configured, you only need to activate the Rook addon.

```yaml title="rook-sample.yaml"
addons:
  rook:
    enabled: true
```

By default, Rook resources are provisioned on all worker nodes (without any constraints).
This behavior can be restricted with node selectors.

??? abstract "Final cluster configuration <i class="click-tip"></i>"

    ```yaml title="rook-sample.yaml"
    hosts:
      - name: localhost
        connection:
          type: local
        dataResourcePools:
          - name: rook-pool

    cluster:
      name: rook-cluster
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
          instances:
            - id: 1
        worker:
          instances:
            - id: 1
              dataDisks:
                - name: rook
                  pool: rook-pool
                  size: 256
            - id: 2
              dataDisks:
                - name: rook
                  pool: rook-pool
                  size: 256
            - id: 3
            - id: 4
              dataDisks:
                - name: rook
                  pool: rook-pool
                  size: 256
                - name: test
                  pool: rook-pool
                  size: 32

    kubernetes:
      version: v1.23.7
      kubespray:
        version: v2.19.0

    addons:
      rook:
        enabled: true
    ```

### Step 4: Apply the configuration

```sh
kubitect apply --config rook-sample.yaml
```

---

## Node selector

The node selector is a dictionary of labels and their potential values.
The node selector restricts on which nodes Rook can be deployed, by selecting only those nodes that match all the specified labels.

### Step 1: Set node labels

To use the node selector effectively, you should give your nodes custom labels.

In this example, we label all worker nodes with the label `rook`.
To ensure that scaling the cluster does not subsequently affect Rook, we set label's value to false by default.
Only the nodes where Rook should be deployed are labeled `#!yaml rook: true`, as shown in the figure below.

<div class="text-center">
  <img
    class="mobile-w-100"
    src="/assets/images/rook-cluster-labels.png" 
    alt="Cluster scheme with labels to restrict Rook deployment"
    width="85%">
</div>

The following configuration snippet shows how to set a default label and override it for a particular instance.

```yaml title="rook-sample.yaml"
cluster:
  nodes:
    worker:
      default:
        labels:
          rook: false
      instances:
        - id: 1
          labels:
            rook: true # (1)!
        - id: 2
          labels:
            rook: true
        - id: 3
          labels:
            rook: true
        - id: 4
```

1.  By default, the label `#!yaml rook: false` is set for all worker nodes. 
    Setting the label `#!yaml rook: true` for this particular instance overrides the default label.

### Step 2: Configure a node selector

So far we have labeled all worker nodes, but labeling is not enough to prevent Rook from being deployed on all worker nodes.
To restrict on which nodes Rook resources can be deployed, we need to configure a node selector.

We want to deploy Rook on the nodes labeled with the label `#!yaml rook: true`, as shown in the figure below.

<div class="text-center">
  <img
    class="mobile-w-100"
    src="/assets/images/rook-cluster-node-selector.png" 
    alt="Cluster scheme of Rook deployment with applied node selector"
    width="85%">
</div>

The following configuration snippet shows how to configure the node selector mentioned above.

```yaml title="rook-sample.yaml"
addons:
  rook:
    enabled: true
    nodeSelector:
      rook: true
```

??? abstract "Final cluster configuration <i class="click-tip"></i>"

    ```yaml title="rook-sample.yaml"
    hosts:
      - name: localhost
        connection:
          type: local
        dataResourcePools:
          - name: rook-pool

    cluster:
      name: rook-cluster
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
          instances:
            - id: 1
        worker:
          default:
            labels:
              rook: false
          instances:
            - id: 1
              labels:
                rook: true
              dataDisks:
                - name: rook
                  pool: rook-pool
                  size: 256
            - id: 2
              labels:
                rook: true
              dataDisks:
                - name: rook
                  pool: rook-pool
                  size: 256
            - id: 3
              labels:
                rook: true
            - id: 4
              dataDisks:
                - name: rook
                  pool: rook-pool
                  size: 256
                - name: test
                  pool: rook-pool
                  size: 32

    kubernetes:
      version: v1.23.7
      kubespray:
        version: v2.19.0

    addons:
      rook:
        enabled: true
        nodeSelector:
          rook: true
    ```

### Step 3: Apply the configuration

```sh
kubitect apply --config rook-sample.yaml
```

</div>
