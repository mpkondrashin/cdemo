# File Security and Artifact Scanner Demo

## Setup

Go to AWS CloudFormation and run the following command:
```shell
    git clone https://github.com/mpkondrashin/cdemo.git
    cd cdemo/file_security
```

Upload ```file_security.yaml``` template using CloudFormation.

Required parameters:
- KeyName: Name of an existing EC2 KeyPair to enable SSH access
- InstanceType: EC2 instance type
- SSHLocation: The IP address range that can be used to SSH to the EC2 instance

To get your IP address, run the following command:
```shell
curl ifconfig.me
```
or use ```https://www.myip.com/```.

After the CloudFormation template finishes running check its output to get the public IP address of the newly created EC2 instance.

Log into the EC2 instance and run the following command:
```shell
sudo bash
cd /root
```

## Artifact Scanner (TMAS) Demo

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

## File Scanner (TMFS) Demo 

Get API Key from Vison One (or use the same on as for TMAS) and run the following command:
```shell
export TMFS_API_KEY=your_api_key
```

### Showcase file scan

Run the following command
```shell
 ./tmfs scan file:./virus.com --region=eu-central-1 
 ```

For a malicious file, you will get the following output:
```json
{
    "scannerVersion": "1.0.0-173",
    "schemaVersion": "1.0.0",
    "scanResult":1,
    "scanId": "84bcd5a3-ac24-4cc5-874c-78f2ca663e07",
    "scanTimestamp": "2025-09-02T17:22:48.927Z",
    "fileName": "./virus.com",
    "foundMalwares":[{"fileName": "virus.com", "malwareName": "Eicar_test_file"}],
    "fileSHA1":"cf8bd9dfddff007f75adf4c2be48005cea317c62",
    "fileSHA256":"131f95c51cc819465fa1797f6ccacf9d494aaaff46fa3eac73ae63ffbdfd8267"
}
```

### Showcase folder scan

Run the following command
```shell
./tmfs scan dir:./libraries --region=eu-central-1 
```

### Showcase SDK

go to ```/root/cdemo/file_security/fscan``` and run the following commands:
```shell
source /etc/profile.d/go.sh
go build
./fscan
```

Source code is in the file main.go


