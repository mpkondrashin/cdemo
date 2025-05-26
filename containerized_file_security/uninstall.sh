#!/bin/bash
set -e
set -x

echo "[CDEMO] Start uninstallation"
exec > >(tee -a uninstall_demo.log) 2>&1

echo "[CDEMO] Delete the test pod"
kubectl delete -f ubuntu_file_security.yaml || true

echo "[CDEMO] Uninstall Helm release"
helm uninstall my-release -n visionone-filesecurity || true

echo "[CDEMO] Delete Kubernetes secrets"
kubectl delete secret token-secret -n visionone-filesecurity || true
kubectl delete secret device-token-secret -n visionone-filesecurity || true

echo "[CDEMO] Delete Kubernetes namespace"
kubectl delete namespace visionone-filesecurity || true

echo "[CDEMO] Remove Helm chart and GPG key"
rm -f visionone-filesecurity-*.tgz
rm -f public-key.asc
rm -f get_helm.sh

echo "[CDEMO] Remove imported GPG keys"
# Extract key ID from the file
KEY_ID=$(gpg --with-colons public-key.asc | awk -F: '/^pub/ {print $5}')

if [ -z "$KEY_ID" ]; then
  echo "No public key found in public-key.asc"
else
    echo "Found key ID: $KEY_ID"
    gpg --batch --yes --delete-key "$KEY_ID"
    echo "Public key $KEY_ID deleted."
fi

echo "[CDEMO] Delete EKS cluster and CloudFormation stack"
aws cloudformation delete-stack --stack-name demo-eks-cluster

echo "[CDEMO] Wait for stack deletion to complete"
aws cloudformation wait stack-delete-complete --stack-name demo-eks-cluster

echo "[CDEMO] Uninstallation completed"
