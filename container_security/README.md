# Container Security Demo Setup

This repository provides ability to run Trend Micro VisionOne Container Security using one script.

## Setup

### Get API Key
You need a Trend Micro VisionOne API Key. Log into Vision One Console, navigate to Administration -> API Keys, and press the "Add API Key" button. Pick Master Administrator role and provide a name for the key. Pick "Create" button. Copy the API Key for later use.

### Setup AWS Infrastructure

Open the AWS Console and go to Cloud Shell.

Run the following commands:
```shell
git clone https://github.com/mpkondrashin/cdemo.git
cd cdemo/container_security
./prepare_cs_demo.sh
```
The script will request Vision One region and API Key. Wait for the script to finish. It can take up to 10 minutes or more.

## Verify

Check the ```prepare_demo.log``` file for the following final line:

```
Demo Setup Completed
```

Check that Container Security is up and running:
```shell
kubectl get pods -n trendmicro-system
```

Check that the Ubuntu pod is up and running:

```shell
kubectl get pods
```

### Demo

## Configure

Apply to the new cluster the some policy and add runtime rules.

## Simulate Malicious Activity

Login into the pod:
```shell
kubectl exec -it ubuntu-server -- /bin/bash
```

Run commands that configured previously in runtime protection policy.

## Uninstall

Run the following commands:
```shell
./uninstall.sh
```

The script is waiting for EKS cluster deletion to complete. It can take up to 10 minutes or more. Meanwhile the ```cdemo``` folder can be deleted.

If stack deletion fails, you can delete it manually in the AWS Console:
1. Go to CloudFormation in the AWS Console
2. Pick the stack name ```demo-eks-cluster```.
3. Press the "Retry delete" button.
4. Pick "Force delete this entire stack" option if prompted.

## Known Issues

### API Key safety
API Key is written to prepare_demo.log file. Delete it after successful setup

### Limitations

Only one demo can be run at a time as it creates EKS cluster with fixed name. If you want to run another demo simultaneously, you need to use another region.