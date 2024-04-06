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

echo "==> TEST: Cluster readiness"

startTime=$(date +%s)
nodes=$(kubectl get nodes | awk 'NR>1 {print $1}')

for node in $nodes; do
    while :; do
        isReady=$(kubectl get node "${node}" \
            -o jsonpath='{.status.conditions[?(@.type=="Ready")].status}'
        )

        if [ "${isReady}" = "True" ]; then
            echo "===> PASS: Node ${node} is ready."
            break
        fi

        currentTime=$(date +%s)
        elapsedTime=$((currentTime - timeStart))

        if [ "${elapsedTime}" -gt "${TIMEOUT}" ]; then
            echo "FAIL: Node ${node} is NOT READY after ${TIMEOUT} seconds!"
            kubectl get nodes
            break
        fi

        sleep 10
    done
done

echo "==> TEST: Running pods"

startTime=$(date +%s)

while :; do
    failedPods=$(kubectl get pods \
        --all-namespaces \
        --field-selector="status.phase!=Succeeded,status.phase!=Running" \
        --output custom-columns="NAMESPACE:metadata.namespace,POD:metadata.name,STATUS:status.phase"
    )

    if [ "$(echo "${failedPods}" | awk 'NR>1')" = "" ]; then
        echo "===> PASS: All pods are running."
        break
    fi

    currentTime=$(date +%s)
    elapsedTime=$((currentTime - startTime))

    if [ "${elapsedTime}" -gt "${TIMEOUT}" ]; then
        echo "==> FAIL: Pods not running after ${TIMEOUT} seconds!"
        echo "${failedPods}"
        break
    fi

    sleep 10
done

echo "==> TEST: DNS"
kubectl run dns-test --image=busybox:1.28.4 --restart=Never -- sleep 180
kubectl wait --for=condition=Ready pod/dns-test --timeout=60s

kubectl exec dns-test -- nslookup kubernetes.default
echo "===> PASS: Local lookup (kubernetes.default)."

kubectl exec dns-test -- nslookup kubitect.io
echo "===> PASS: External lookup (kubitect.io)."

# All tests have passed.
FAIL=0
