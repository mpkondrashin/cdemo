# Showcase Trend AI Vision One IPS

## Table of contents
- [Requirements](#requirements)
- [Deployment](#deployment)
- [Demo](#demo)
- [Cleanup](#cleanup)

## Requirements

- AWS account with permissions to create VPC, EC2, Network Firewall, S3, and CloudWatch resources
- [Optional] Configured AWS CLI with appropriate credentials and region (e.g. us-east-1) if CLI commands will be used within the caurse of the demo.

## Deployment

### Deploy CloudFormation Template

<img src="infrastructure-composer-ips.yaml.png" alt="CloudFormation"/>

#### Using CLI:
Run the following command to deploy the stack:
```bash
aws cloudformation deploy \
  --stack-name ips-demo \
  --template-file ips.yaml \
  --capabilities CAPABILITY_NAMED_IAM
```

And wait for the stack to be created (this may take a few minutes).

Note: Any other name can be used for the stack, but it should be consistent with the rest of the documentation.

#### Using Console:
1. Go to AWS CloudFormation Console
2. Click "Create stack"
3. Upload the `ips.yaml` template
4. Fill in the parameters
5. Click "Create stack"
6. Wait for the stack to be created (this may take a few minutes)

On the AWS CloudFormation console chechk outputs of the stack to get the S3 bucket name and other information.

### Configure Network Firewall
1. On the AWS console Marketplace search for `TrendAI Vision One™ Cloud IPS` -> `View purchase options` Button -> `Subscribe` Button.
2. Go to VPC
3. On the left side go to Network Firewall -> Firewall policies
4. Create a new firewall policy
5. Add the TrendAI Vision One IPS rule group to the policy

### Upload pcaps.zip
1. Go to S3
2. Upload the `pcaps.zip` file to the S3 bucket created by the CloudFormation stack

## Demo

### Generate alers
1. On AWS console go to EC2
2. Select the instance ```ips-demo-attacker```
3. Push "Connect" button
4. Go to "SSM Session Manager" and push "Connect" button 
5. In the session manager console run the following command:
```bash
sudo replay_pcaps.sh
```
If ```pcaps.zip``` file has password, you will be prompted to enter it. The password is `infected`.

### Check logs

1. On AWS console go to CloudWatch
2. Go to `Log Management` -> `Log Groups`
3. Find the log group with the name `/aws/network-firewall/ips-demo/flow` and `/aws/network-firewall/ips-demo/alert` 
4. Click on the log group
5. Click on the log stream
6. You should see the logs from the IPS

## Cleanup

### Remove Firewall policy
1. Go to Network Firewall -> Firewall policies
2. Delete the firewall policy

### Delete CloudFormation Stack
Use the following command
```bash
aws cloudformation delete-stack --stack-name ips-demo
```
or delete it using AWS Console:
1. Go to AWS CloudFormation Console
2. Choose previously created stack
3. Choose "Delete stack"
4. Confirm the deletion
