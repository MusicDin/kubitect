#!/bin/sh
set -eu

# Check input arguments.
if [ "${1:-}" = "" ] || [ "${2:-}" = "" ] || [ "${3:-}" = "" ] || [ "${4:-}" = "" ] || [ "${5:-}" = "" ]; then
    echo "Usage: ${0} <cluster_name> <distro> <network_plugin> <k8s_version> <k8s_manager>"
    exit 1
fi

CLUSTER="${1}"
DISTRO="${2}"
NETWORK_PLUGIN="${3}"
K8S_VERSION="${4}"
K8S_MANAGER="${5}"

echo "==> DEPLOY: Cluster ${DISTRO}/${NETWORK_PLUGIN}/${K8S_VERSION}"

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
    loadBalancer:
      default:
        cpu: 1
        ram: 1
        mainDiskSize: 8
      instances:
        - id: 1
          ip: 192.168.113.100
    master:
      default:
        cpu: 2
        ram: 2
        mainDiskSize: 16
      instances:
        - id: 1
          ip: 192.168.113.10
    worker:
      default:
        cpu: 2
        ram: 2
        mainDiskSize: 16
      instances:
        - id: 1
          ip: 192.168.113.20
        - id: 2
          ip: 192.168.113.21

kubernetes:
  manager: ${K8S_MANAGER}
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
