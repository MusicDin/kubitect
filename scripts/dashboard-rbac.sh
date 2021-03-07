#!/bin/sh
USER="$1"
NAMESPACE="$2"
KUBECONFIG=config/admin.conf

echo "Creating new service account for kubernetes dashboard..."

kubectl --kubeconfig=${KUBECONFIG} -n ${NAMESPACE} create serviceaccount ${USER} >/dev/null 2>&1

echo "Add permissions to created service account"
kubectl --kubeconfig=${KUBECONFIG} create clusterrolebinding ${USER}-cluster-admin --clusterrole=cluster-admin --serviceaccount=${NAMESPACE}:${USER} >/dev/null 2>&1

echo "Printing token"
kubectl --kubeconfig=${KUBECONFIG} -n ${NAMESPACE} describe secret $(kubectl --kubeconfig=${KUBECONFIG} -n ${NAMESPACE} get secret | grep ${USER}-token | awk '{print $1}') | awk '$1=="token:"{print $2}'
