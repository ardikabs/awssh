# awssh
## Description
`awssh` is a simple CLI providing an ssh access to EC2 utilizing an ec2-instance-connect command.<br>
The `awssh` is extending the ec2-instance-connect command to be aware with ssh-agent and/or populate new temporary ssh keypair while trying to establish an ssh connection to the EC2 instance target.

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
* Go 1.13.5 or later

### Setup
* Install Git
* Install Go 1.13.5 or later
* Clone this repository

### Build and run binary file
To build binary file:
1. Linux: `make build-linux`
2. Windows: `make build-windows`
3. MacOS: `make build-darwin`

### Installation Guide
Check the [release page](https://github.com/ardikabs/awssh/releases).
Please read [user guide](USAGE.md) for further use.