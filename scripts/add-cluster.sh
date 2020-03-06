#!/bin/bash

CLUSTER_NAME="${1}"

if [ -z $CLUSTER_NAME ]; then
    echo "CLUSTER_NAME is not given!"
    exit 1
fi

SERVER_URL=$(export CURRENT_CONTEXT=$(kubectl config current-context) && export CURRENT_CLUSTER=$(kubectl config view -o go-template="{{\$curr_context := \"$CURRENT_CONTEXT\" }}{{range .contexts}}{{if eq .name \$curr_context}}{{.context.cluster}}{{end}}{{end}}") && echo $(kubectl config view -o go-template="{{\$cluster_context := \"$CURRENT_CLUSTER\"}}{{range .clusters}}{{if eq .name \$cluster_context}}{{.cluster.server}}{{end}}{{end}}"))
CA_DATA=$(echo $(kubectl get secret -n kube-system -o go-template='{{index .data "ca.crt" }}' $(kubectl get sa kubestorm-user -n kube-system -o go-template="{{range .secrets}}{{.name}}{{end}}")))
TOKEN=$(echo $(kubectl get secret -n kube-system -o go-template='{{index .data "token" }}' $(kubectl get sa kubestorm-user -n kube-system -o go-template="{{range .secrets}}{{.name}}{{end}}")) | base64 --decode)

curl --location --request POST "http://localhost:8080/v1/auth/${CLUSTER_NAME}" \
--form "server=${SERVER_URL}" \
--form "token=${TOKEN}" \
--form "serverCADataString=${CA_DATA}"