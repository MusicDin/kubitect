# Getting Started

The aim of this short tutorial is to install Kubernetes cluster on a local machine. At the end of the tutorial, you will have a functional cluster that consist of a single master node and three worker nodes.

> :scroll: **Note:** For successful installation of the Kubernetes cluster some [requirements](/docs/requirements.md) need to be met.

> :bulb: **Tip:** If you encounter any issues during the installation, 
> please refer to the [troubleshooting](docs/troubleshooting.md) page first.

## Quick start

### Step 1 - Download and install `tkk` tool

```
wget -O tkk https://raw.githubusercontent.com/MusicDin/terraform-kvm-kubespray/<BRANCH>/scripts/tkk.sh

sudo install tkk /usr/local/bin
```

### Step 2 - Initialize cluster configuration

Run the following command to create the default cluster.
Cluster will be created in `~/.tkk/clusters/default/` directory.

```
tkk create cluster
```

> :scroll: **Note**: Using a `--name` option, cluster name can be changed. This way multiple clusters can be created.

### Step 3 - Deploy the cluster


Prepare the configuration file (default is cluster.yaml) and apply it to create the cluster.


Cluster variables are defined within [cluster.yml](cluster.yml) file. 
You can review or modify them before cluster creation.
Out of the box, variables are prepared to install the Kubernetes cluster with one master and three worker nodes on a local machine (localhost).

> :warning: **Tip:** If some variable is not understandable, please refere to the [configuration](/docs/configuration.md) explanation.

Apply the configuration to set up the cluster.
```
tkk apply --config cluster.yaml
```

> :scroll: **Note:** *The installation process can take up to 20 minutes depending on the configuration.*


### Step 4 - Test the cluster

All configuration files will be generated in `config/` directory,
and one of them will be `admin.conf` which is actually a `kubeconfig` file.

Test if the cluster works by displaying all cluster nodes.

```
kubectl --kubeconfig=config/admin.conf get nodes
```


## What's next?

[Learn how to manage the created cluster](./cluster-management.md)