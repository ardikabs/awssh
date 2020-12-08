# awssh
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.14-61CFDD.svg?style=flat-square)
[![Go Report Card](https://goreportcard.com/badge/github.com/ardikabs/awssh?style=flat-square)](https://goreportcard.com/report/github.com/ardikabs/awssh)
## Description
`awssh` is a simple CLI providing an ssh access to EC2 utilizing an [EC2 Instance Connect](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Connect-using-EC2-Instance-Connect.html) feature.<br>
The `awssh` is extending the `aws ec2-instance-connect` command to be aware with ssh-agent and/or populate new temporary ssh keypair while trying to establish an ssh connection to the EC2 instance target.

## Diagram Flow
```bash
<user>
    |
<AWS credentials>
    |
<Get list of ec2 instances>
    |
<Create temporary ssh keypair>
    or
<Use existing ssh keypair within ssh-agent>
    |
<Send ssh public key>
    |
<Establish an ssh connection>
```

## Development Guide
### Prerequisites
* Go 1.14 or later

### Setup
* Install Git
* Install Go 1.14 or later
* Clone this repository

### Build and run binary file
To build binary file: `make build`

### Installation Guide
Check the [release page](https://github.com/ardikabs/awssh/releases).
Please read [user guide](USAGE.md) for further use.