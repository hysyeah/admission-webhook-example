#!/bin/bash

SERVICE=pod-webhook
NAMESPACE=default

cfssl gencert -initca ../ca-csr.json | cfssljson -bare ca
cfssl gencert \
  -ca=../ca.crt \
  -ca-key=../ca.key \
  -config=../ca-config.json \
  -hostname=127.0.0.1,$SERVICE,$SERVICE.kube-system,$SERVICE.$NAMESPACE,$SERVICE.$NAMESPACE.svc \
  -profile=default \
  ../webhook-csr.json | cfssljson -bare $SERVICE

kubectl create secret tls $SERVICE-tls -n $NAMESPACE \
  --cert=$SERVICE.pem \
  --key=$SERVICE-key.pem

base64 -w 0 ../ca.crt > ca.crt.base64
echo "Updating deployment caBundle"
sed -i -E "s/(caBundle: )(.*)/\1$(cat ca.crt.base64)/" ../../deployment.yaml