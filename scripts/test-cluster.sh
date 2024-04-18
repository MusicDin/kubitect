#!/bin/sh
set -eu

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

echo "==> DEBUG: Cluster info"
kubectl cluster-info
kubectl get nodes

echo "==> TEST: All nodes ready"
kubectl wait --for=condition=ready node --all --timeout=120s

echo "==> TEST: All pods ready"
set +e
# Wait for all pods to be ready. Retry a few times, as it may happen that
# cluster has no pods when check is ran, which will result in an error.
for i in $(seq 5); do
    podsReadiness=$(kubectl wait --for=condition=ready pods --all -A --timeout=30s)
    if [ "$?" -eq 0 ]; then
        break
    else
        echo "(attempt $i/5) Pods are still not ready. Retrying in 10 seconds..."
        sleep 10
    fi
done
set -e
echo "${podsReadiness}"

echo "==> TEST: DNS"
kubectl apply -f https://k8s.io/examples/admin/dns/dnsutils.yaml
kubectl wait --for=condition=Ready pod/dnsutils --timeout=60s

kubectl exec dnsutils -- nslookup kubernetes.default
echo "===> PASS: Local lookup (kubernetes.default)."

kubectl exec dnsutils -- nslookup kubitect.io
echo "===> PASS: External lookup (kubitect.io)."

# All tests have passed.
FAIL=0
