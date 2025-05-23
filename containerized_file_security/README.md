# Containerize File Security Showcase

This repository contains a showcase of how to containerize file security using Trend Micro VisionOne File Security.

## Setup

### Get Token
You need to have a Trend Micro VisionOne File Security registration token. Login into Vision One Console and navigate to Clout Security -> File Security -> Containerized Scanner and press "Get registration token" button. Save the token for later use.

### Setup AWS Infrastructure

Open AWS Console and go to Cloud Shell.

Run the following commands:
```shell
git clone https://github.com/mpkondrashin/cdemo.git
cd cdemo/containerized_file_security
./prepare_demo.sh <registration token>
```
If registration token is not provided, it will be asked interactively.

Wait for the script to finish. It can take up to 10 minutes or more.

## Verify

Check prepare_demo.log file for the following final line:
```
Demo Setup Completed
```

Check that File Security is up and running:
```shell
kubectl get pods -n visionone-filesecurity
```

Check that ubuntu pod is up and running:
```shell
kubectl get pods
```

## Perform Demo Scan

Login into the pod:
```shell
kubectl exec -it ubuntu-server -- /bin/bash
```

Run the following command:
```shell
./tmfs scan file:./virus.com --tls=false --endpoint=my-release-visionone-filesecurity-scanner.visionone-filesecurity.svc.cluster.local:50051
```

## Uninstall

Initiate EKS cluster deletion on Cloud Shell:
```shell
aws cloudformation delete-stack --stack-name demo-eks-cluster
```

(Optional) Wait until deletion is complete:

```shell
aws cloudformation wait stack-delete-complete --stack-name demo-eks-cluster
```