<h1 align="center">Rook cluster</h1>

!!! warning "Important"

    Since the Rook addon is still under development, it may not work as expected.
    Therefore, we would greatly appreciate your feedback.

This example shows how to use Kubitect to set up distributed storage with Rook.
For distributed storage, we add an additional data disk to each virtual machine as shown on the figure below.

<div align=center>
  <img
    class="mobile-w-100"
    src="/assets/images/rook-basic.png" 
    alt="Basic Rook cluster scheme"
    width="75%">
</div>

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

### Step 2: Attaching data disks

After the data resource pool is configured, we are ready to allocate some data disks to the virtual machines.

```yaml title="rook-sample.yaml"
cluster:
  nodes:
    worker:
      instances:
        - id: 1
          dataDisks:
            - name: rook
              pool: rook-pool # (1)
              size: 256
        - id: 2
          dataDisks:
            - name: rook
              pool: rook-pool
              size: 256
        - id: 3
          dataDisks:
            - name: rook
              pool: rook-pool
              size: 256
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
Similarly, all attached data disks are recognized by Rook and used for distributed storage.
This behavior can be restricted with node and data disk selectors.

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