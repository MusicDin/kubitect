<h1 align="center">Restricting Rook cluster</h1>

!!! warning "Important"

    Since the Rook addon is still under development, it may not work as expected.
    Therefore, we would greatly appreciate your feedback.

In the [Rook cluster](../rook-cluster) example, we demonstrated how to set up Rook on three worker nodes.
Node and data disk selectors can be used to deploy Rook only on specific nodes and to select specific data disks for distributed storage.
Therefore, this example extends aforementioned example and shows how to use node and data disk selectors.

!!! note "Note"

    In this example, we skip the explanation of how to attach data disks to the virtual machines because this is already explained in the [Rook cluster](../rook-cluster) example.

### Step 1: Label nodes

The node selector restricts on which nodes Rook can be deployed, by selecting only those nodes that match all the specified labels.

In this example, we label all worker nodes with the label `rook`.
To ensure that scaling the cluster does not subsequently affect Rook, we set label value to false by default.
Only the nodes where Rook should be deployed are labeled `#!yaml rook: true`, as shown on the figure below.

<div align=center>
  <img
    class="mobile-w-100"
    src="/assets/images/rook-labels.png" 
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
            rook: true # (1)
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

### Step 2: Enable Rook addon

To enable Rook addon, simply set `addon.rook.enabled` property to true.

```yaml title="rook-sample.yaml"
addons:
  rook:
    enabled: true
```

### Step 3: Node selector

So far we have labeled all worker nodes, but labels are not enough to prevent Rook from being deployed on all worker nodes.
To restrict on which nodes Rook resources can be deployed, we need to configure a node selector.
The node selector is a dictionary of labels and their potential values.

In this example, we want to deploy Rook on the nodes labeled with the label `#!yaml rook: true`.

```yaml title="rook-sample.yaml"
addons:
  rook:
    enabled: true
    nodeSelector:
      rook: true
```

The following figure shows which nodes and data disks are used by Rook when the current configuration is applied.

<div align=center>
  <img
    class="mobile-w-100"
    src="/assets/images/rook-node-selector.png" 
    alt="Cluster scheme of Rook deployment with applied node selector"
    width="85%">
</div>

### Step 4: Data disk selector

!!! warning "Important"

    Use the data disk selector as a last resort.
    The data disk selector disables Rook's automatic detection of the data disk.
    As a result, any further addition of data disks will require manual reconfiguration of Rook.

Data disk selector is a list of data disk names.
It can be used to prevent Rook from consuming disks that do not match exactly one of the given names.

For example, we could limit Rook to use only disks that are named `rook`.

```yaml title="rook-sample.yaml"
addons:
  rook:
    enabled: true
    nodeSelector:
      rook: true
    dataDiskSelector:
      - "rook"
```

The following figure shows which nodes and data disks are used by Rook when the current configuration is applied.

<div align=center>
  <img
    class="mobile-w-100"
    src="/assets/images/rook-disk-selector.png" 
    alt="Cluster scheme of Rook deployment with applied node and data disk selector"
    width="85%">
</div>

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
                - name: test
                  pool: rook-pool
                  size: 32
            - id: 2
              labels:
                rook: true
              dataDisks:
                - name: rook
                  pool: rook-pool
                  size: 256
                - name: backup
                  pool: rook-pool
                  size: 1024
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
        dataDiskSelector:
          - "rook"
    ```

### Step 5: Apply the configuration

```sh
kubitect apply --config rook-sample.yaml
```