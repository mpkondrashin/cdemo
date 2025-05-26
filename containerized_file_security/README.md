# Containerized File Security Showcase

This repository contains a showcase of how to containerize file security using Trend Micro VisionOne File Security.

## Setup

### Get Token
You need a Trend Micro VisionOne File Security registration token. Log into Vision One Console, navigate to Clout Security -> File Security -> Containerized Scanner, and press the "Get registration token" button. Save the token for later use.

### Setup AWS Infrastructure

Open the AWS Console and go to Cloud Shell.

Run the following commands:
```shell
git clone https://github.com/mpkondrashin/cdemo.git
cd cdemo/containerized_file_security
./prepare_demo.sh <registration token>
```
If a registration token is not provided, it will be requested interactively.

Wait for the script to finish. It can take up to 10 minutes or more.

## Verify

Check the ```prepare_demo.log``` file for the following final line:

```
Demo Setup Completed
```

Check that File Security is up and running:
```shell
kubectl get pods -n visionone-filesecurity
```

Check that the Ubuntu pod is up and running:

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

Run the following commands:
```shell
./uninstall.sh
```

The script is waiting for EKS cluster deletion to complete. It can take up to 10 minutes or more. Meanwhile the ```cdemo``` folder can be deleted.

## Known Issues

### Token safety
Token is written to prepare_demo.log file. Delete it after successful setup

### Wrong token symptoms

if ```kubectl get pods -n visionone-filesecurity``` shows status ```Init:CrashLoopBackOff``` for ```my-release-visionone-filesecurity-backend-communicator``` pod, it could signal that your Token is invalid. 

### Limitations

Only one demo can be run at a time as it creates EKS cluster with fixed name. If you want to run another demo simultaneously, you need to use another region.