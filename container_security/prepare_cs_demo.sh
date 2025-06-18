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


echo "[CDEMO] Start logging"
exec > >(tee -a prepare_cs_demo.log) 2>&1

echo "[CDEMO] Create EKS cluster"
aws cloudformation create-stack \
  --stack-name csdemo-eks-cluster \
  --template-body file://eks.yaml \
  --capabilities CAPABILITY_IAM CAPABILITY_NAMED_IAM

echo "[CDEMO] Wait for EKS cluster to be created. It can take up to 10 minutes or more"
aws cloudformation wait stack-create-complete --stack-name csdemo-eks-cluster

echo "[CDEMO] Configure kubectl to use the EKS cluster"
aws eks update-kubeconfig --name csdemo-eks-cluster

CLUSTER_NAME=$(aws cloudformation describe-stacks --stack-name csdemo-eks-cluster \
  --query "Stacks[0].Outputs[?OutputKey=='ClusterName'].OutputValue" --output text)
echo "[CDEMO] Cluster name: $CLUSTER_NAME"

echo "[CDEMO] Update kubectl to use the EKS cluster"
aws eks update-kubeconfig --name $CLUSTER_NAME

echo "[CDEMO] Check EKS cluster"
kubectl get nodes
kubectl get pods --all-namespaces

echo "[CDEMO] Install Helm"
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
chmod 700 get_helm.sh
./get_helm.sh

echo "[CDEMO] Install Python"
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
API_TOKEN=$API_TOKEN python3 vone_overrides.py setup $REGION

echo "[CDEMO] Install Container Security ()"
helm install \
    trendmicro \
    --namespace trendmicro-system --create-namespace \
    --values overrides.yaml \
    https://github.com/trendmicro/cloudone-container-security-helm/archive/master.tar.gz

echo "[CDEMO] Lets wait for some time for Container Security to be ready"
sleep 30

echo "[CDEMO] This is it, Container Security should be up and running. Check it with the command"
kubectl get pods -n trendmicro-system

echo "[CDEMO] Install test pod"

kubectl apply -f ubuntu_cs.yaml

echo " [CDEMO] Login into the pod"
echo "Run the following commands:"
echo "kubectl exec -it ubuntu-server -- /bin/bash"

echo "[CDEMO] Demo Setup Completed. Check the README.md for more information."
