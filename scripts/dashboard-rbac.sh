#!/bin/sh
USER="$1"
NAMESPACE="$2"

echo "Creating new service account for kubernetes dashboard..."
kubectl --kubeconfig=../config/admin.conf -n ${NAMESPACE} create serviceaccount ${USER} >/dev/null 2>&1

echo "Add permissions to created service account"
kubectl --kubeconfig=../config/admin.conf create clusterrolebinding ${USER}-cluster-admin --clusterrole=cluster-admin --serviceaccount=${NAMESPACE}:${USER} >/dev/null 2>&1

echo "Printing token"
kubectl --kubeconfig=../config/admin.conf -n ${NAMESPACE} describe secret $(kubectl -n ${NAMESPACE} get secret | grep ${USER}-token | awk '{print $1}') | awk '$1=="token:"{print $2}'
