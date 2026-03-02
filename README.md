# Vision One Cloud Security Demo

This repository contains a set of CloudFormation templates and scripts to showcase the capabilities of Trend Micro Vision One Cloud Security.

## Prerequisites

- Vision One Account

## Setup

- Open AWS Cloud Shell
- Clone this repository
```bash
git clone https://github.com/mpkondrashin/cdemo.git
cd cdemo
```

To showcase Container Security, run the following command:
```bash
cd container_security
./prepare_cs_demo.sh
```

To showcase Containerized File Security, run the following command:
```bash
cd containerized_file_security
./prepare_demo.sh
```

Check README file in each directory for more information.
