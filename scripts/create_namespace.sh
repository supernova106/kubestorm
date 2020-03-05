#!/bin/bash

kubectl apply -f - <<EOF
kind: Namespace
apiVersion: v1
metadata:
  name: kubestorm
  labels:
    role: kubestorm 
    owner: foo
EOF