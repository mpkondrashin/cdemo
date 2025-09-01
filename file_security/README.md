# File Security and Artifact Scanner Demo

## Setup

Go to AWS CloudFormation and upload ```file_security.yaml``` template.

Required parameters:
- KeyName: Name of an existing EC2 KeyPair to enable SSH access
- InstanceType: EC2 instance type
- SSHLocation: The IP address range that can be used to SSH to the EC2 instance

To get your IP address, run the following command:
```shell
curl ifconfig.me
```
or use ```https://www.myip.com/```.

After Cloud Formation template will finish to run, login into the EC2 instance and run the following command:
```shell
sudo bash
cd /root
```

## Artifact Scanner

Get API Key from Vison One and run the following command:
```shell
export TMAS_API_KEY=your_api_key
```

### Showcase container image scan

Run the following command
```shell
./tmas scan docker:ubuntu:xenial-20161010 -VMS --region=eu-central-1  > tmas.json
```

```-v``` can be added to get verbose output.

tmas.json will contain inspection results.

To get this file locally, run the following command:
```shell
scp -i <your_key_pair> ubuntu@<your_ec2_ip>:/root/tmas.json .
```

### Showcase exported container image scan

Run the following command
```shell
./tmas scan docker-archive:ubuntu-xenial-20161010.tar -VMS --region=eu-central-1  > tmas_tar.json
```

tmas_tar.json will contain inspection results (same as for previous scan).

### Showcase RPM scan

Run the following command
```shell
./tmas scan file:firefox-78.7.0-1.mga8.x86_64.rpm -VS --region=eu-central-1 
```

### Showcase tar.gz scan

Run the following command
```shell
./tmas scan file:openssl-engine-0.9.6.tar.gz -VS --region=eu-central-1 
```

### Showcase scanning libraries

Run the following command
```shell
./tmas scan folder:libraries -VS --region=eu-central-1 
```
