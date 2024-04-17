#!/bin/sh
set -eu

TIMEOUT=600 # seconds

defer() {
    if [ "${FAIL}" = "1" ]; then
        echo "==> DEBUG: Cluster events"
        kubectl get events --all-namespaces

        echo "==> FAIL"
        exit 1
    fi

    echo "==> PASS"
    exit 0
}

FAIL=1
trap defer EXIT HUP INT TERM

echo "==> TEST: All nodes ready"
kubectl wait --for=condition=ready node --all --timeout=120s

startTime=$(date +%s)
nodes=$(kubectl get nodes | awk 'NR>1 {print $1}')

echo "==> TEST: All pods ready"
kubectl wait --for=condition=ready pods --all -A --timeout=120s

echo "==> TEST: DNS"
kubectl apply -f https://k8s.io/examples/admin/dns/dnsutils.yaml
kubectl wait --for=condition=Ready pod/dnsutils --timeout=60s

kubectl exec dnsutils -- nslookup kubernetes.default
echo "===> PASS: Local lookup (kubernetes.default)."

kubectl exec dnsutils -- nslookup kubitect.io
echo "===> PASS: External lookup (kubitect.io)."

# All tests have passed.
FAIL=0
