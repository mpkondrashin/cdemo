#!/bin/bash
set -e
set -x

if [ -z "$1" ]; then
  read -p "Enter the Region: " REGION
else
  REGION="$1"
fi

if [ -z "$2" ]; then
  read -p "Enter the API token: " API_TOKEN
else
  API_TOKEN="$2"
fi

echo "[CDEMO] Start uninstallation"
exec > >(tee -a uninstall_demo.log) 2>&1

echo "[CDEMO] Unregister cluster"
API_TOKEN=$API_TOKEN python3 vone_overrides.py cleanup $REGION

echo "[CDEMO] Remove venv"
rm -rf venv

echo "[CDEMO] Delete the test pod"
kubectl delete -f ubuntu_cs.yaml || true

echo "[CDEMO] Uninstall Helm release"
helm uninstall trendmicro -n trendmicro-system || true

echo "[CDEMO] Delete Kubernetes namespace"
kubectl delete namespace trendmicro-system || true

echo "[CDEMO] Remove Helm chart and GPG key"
rm -f get_helm.sh

echo "[CDEMO] Delete EKS cluster and CloudFormation stack"
aws cloudformation delete-stack --stack-name demo-eks-cluster

echo "[CDEMO] Wait for stack deletion to complete"
aws cloudformation wait stack-delete-complete --stack-name demo-eks-cluster

echo "[CDEMO] Uninstallation completed"
