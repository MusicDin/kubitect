#!/bin/sh
set -eu

if [ "${1:-}" = "" ]; then
    echo "Usage: ${0} <cluster_name>"
    exit 1
fi

CLUSTER="${1}"

echo "==> DESTROY: Cluster ${CLUSTER}"
kubitect destroy --cluster "${CLUSTER}"

if [ -d "${HOME}/.kubitect/clusters/${CLUSTER}" ]; then
    echo "==> FAIL: Cluster directory still exists"
    exit 1
fi

echo "==> PASS"
