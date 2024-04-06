#!/bin/sh
set -eu

# Check input arguments.
if [ "${1:-}" = "" ] || [ "${2:-}" = "" ] || [ "${3:-}" = "" ] || [ "${4:-}" = "" ]; then
    echo "Usage: ${0} <cluster_name> <distro> <network_plugin> <k8s_version>"
    exit 1
fi

CLUSTER="${1}"
DISTRO="${2}"
NETWORK_PLUGIN="${3}"
K8S_VERSION="${4}"

echo "==> DEPLOY: Cluster (Single Node) ${DISTRO}/${NETWORK_PLUGIN}/${K8S_VERSION}"

# Create config.
cat <<-EOF > config.yaml
hosts:
  - name: localhost
    connection:
      type: local

cluster:
  name: ${CLUSTER}
  network:
    mode: nat
    cidr: 192.168.113.0/24
  nodeTemplate:
    user: k8s
    updateOnBoot: true
    cpuMode: host-passthrough
    ssh:
      addToKnownHosts: true
    os:
      distro: ${DISTRO}
  nodes:
    master:
      default:
        cpu: 2
        ram: 4
        mainDiskSize: 32
      instances:
        - id: 1
          ip: 192.168.113.10

kubernetes:
  version: ${K8S_VERSION}
  networkPlugin: ${NETWORK_PLUGIN}
EOF

echo "Config:"
echo "---"
cat config.yaml
echo "---"

# Apply config and export kubeconfig.
mkdir -p "${HOME}/.kube"
kubitect apply --config config.yaml
kubitect export kubeconfig --cluster "${CLUSTER}" > "${HOME}/.kube/config"

echo "==> DEBUG: Cluster info"
kubectl cluster-info
kubectl get nodes
