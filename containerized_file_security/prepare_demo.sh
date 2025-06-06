#!/bin/bash
set -e
set -x

if [ -z "$1" ]; then
  read -p "Enter the registration token: " REGISTRATION_TOKEN
else
  REGISTRATION_TOKEN="$1"
fi

echo "[CDEMO] Start logging"
exec > >(tee -a prepare_demo.log) 2>&1

echo "[CDEMO] Create EKS cluster"
aws cloudformation create-stack \
  --stack-name demo-eks-cluster \
  --template-body file://cf_eks.yaml \
  --capabilities CAPABILITY_IAM CAPABILITY_NAMED_IAM

echo "[CDEMO] Wait for EKS cluster to be created. It can take up to 10 minutes or more"
aws cloudformation wait stack-create-complete --stack-name demo-eks-cluster

echo "[CDEMO] Configure kubectl to use the EKS cluster"
aws eks update-kubeconfig --name demo-eks-cluster

CLUSTER_NAME=$(aws cloudformation describe-stacks --stack-name demo-eks-cluster \
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

echo "[CDEMO] Install File Security (https://trendmicro.github.io/visionone-file-security-helm/)"
kubectl create namespace visionone-filesecurity
kubectl create secret generic token-secret --from-literal=registration-token="$REGISTRATION_TOKEN" -n visionone-filesecurity
kubectl create secret generic device-token-secret -n visionone-filesecurity

helm repo add visionone-filesecurity https://trendmicro.github.io/visionone-file-security-helm/
helm repo update
curl -o public-key.asc https://trendmicro.github.io/visionone-file-security-helm/public-key.asc
gpg --import public-key.asc || true
gpg --export >~/.gnupg/pubring.gpg
helm pull --verify visionone-filesecurity/visionone-filesecurity
helm install my-release visionone-filesecurity/visionone-filesecurity -n visionone-filesecurity

echo "[CDEMO] This is it, File Security is up and running. Check it with the command"
kubectl get pods -n visionone-filesecurity

echo "[CDEMO] Install test pod"

kubectl apply -f ubuntu_file_security.yaml

echo " [CDEMO] Login into the pod"
echo "Run the following commands:"
echo "kubectl exec -it ubuntu-server -- /bin/bash"
echo "./tmfs scan file:./virus.com --tls=false --endpoint=my-release-visionone-filesecurity-scanner.visionone-filesecurity.svc.cluster.local:50051"

echo "[CDEMO] Demo Setup Completed. Check the README.md for more information."

