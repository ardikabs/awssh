# Usage
## AWS Credentials
You can select one of the followings:
1. export `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`
2. use `AWS_PROFILE` from aws shared-credentials `~/.aws/credentials`
3. if you are using `AWS_PROFILE` from `~/.aws/config`. Then you need to export `export AWS_SDK_LOAD_CONFIG=1`

Set the AWS region either in your AWS credentials or environment variables (`AWS_DEFAULT_REGION` or `AWS_REGION`) or define it from `awssh` flags (`--region <aws-project-region>`)

## Environment Variables
To using `awssh` you can setup your configuration from environment variables as follows:
* `AWSSH_DEBUG`: Enabled debug mode for `awssh`. Default to `0` (false).
* `AWSSH_TAGS`: List of EC2 tags key-value pair. Default to `"Name=*"`.
* `AWSSH_SSH_USERNAME`: An EC2 ssh username. Default to `ec2-user`.
* `AWSSH_SSH_PORT`: An EC2 ssh port. Default to `22`.
* `AWSSH_SSH_OPTS`: An additional ssh options. Default to `"-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/nul -o ConnectTimeout=5"`

## Examples
### How-to
```bash
$ awssh --help

awssh is a simple CLI providing an ssh access to EC2 utilizing ec2-instance-connect

Usage:
  awssh [flags]

Examples:

  # List all of the EC2 instances given by the credentials
  awssh --region=ap-southeast-1

  # Select EC2 instance with instance-id
  awssh i-0387e016c47c6170c

  # Select EC2 instance given with selected tags
  awssh --tags "Environment=production,Project=jenkins,Owner=SRE"

  # Use an additional ssh options
  awssh --tags "Environment=staging,ProductDomain=bastion" --ssh-username=centos --ssh-port=2222 --ssh-opts="-o ServerAliveInterval=60s"

  # Use public ip to connect to the EC2 instance
  awssh --use-public-ip

Available Commands:
  help        Help about any command
  version     Print the version number of awssh

Flags:
  -d, --debug                 Enabled debug mode
  -h, --help                  help for awssh
      --region string         Default AWS region to be used. Either set AWS_REGION or AWS_DEFAULT_REGION (default "ap-southeast-1")
  -o, --ssh-opts string       An additional ssh options (default "-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null")
  -p, --ssh-port string       An EC2 instance ssh port (default "22")
  -u, --ssh-username string   EC2 SSH username (default "ec2-user")
  -t, --tags string           EC2 tags key-value pair (default "Name=*")
      --use-public-ip         Use public IP to access the EC2 instance
```
### Debug Mode
```bash
$ awssh --debug

DEBUG  Region: ap-southeast-1
DEBUG  Use the following filters to filter EC2 instances: [{
  Name: "instance-state-name",
  Values: ["running"]
} {
  Name: "tag:Name",
  Values: ["*"]
}]
DEBUG  Filter EC2 instances with the following tags: Name=*
DEBUG  Found 8 EC2 instances on region ap-southeast-1

Use the arrow keys to navigate: ↓ ↑ → ←  and / toggles search
Select an instance:
  » machine-a i-07f02a0bfd0952301 (10.0.115.250)
    machine-b i-0387e016c47c6170c (10.0.122.145)
    machine-c i-0ebf8d454d9fd4c5e (10.0.127.160)
    ssh-jumper i-08c76965ce9ee0828 (10.0.17.233 / 54.169.42.125)
    master-1.masters.k8s.kops.internal i-0b2566fcc894c1bd1 (10.0.132.143)
    master-2.masters.k8s.kops.internal i-05c0309be99c8a097 (10.0.148.154)
    master-3.masters.k8s.kops.internal i-052913ef86123d500 (10.0.186.183)
    nodes-a.nodes.k8s.kops.internal i-07fc020d8c7f50e27 (10.0.172.143)

nodes-a.nodes.k8s.kops.internal i-07fc020d8c7f50e27
DEBUG  Select EC2 instance 'nodes-a.nodes.k8s.kops.internal' (i-07fc020d8c7f50e27)
DEBUG  Use existing ssh-rsa keypair from ssh-agent (SHA256:2ISinysBKLIbWburvJesabZQaj1uzDkMouCoS45mlf4)
DEBUG  Sending SSH Public Key for EC2 instance 'nodes-a.nodes.k8s.kops.internal' (i-07fc020d8c7f50e27)
DEBUG  Establish an SSH connection to the EC2 instance target 'nodes-a.nodes.k8s.kops.internal' (i-07fc020d8c7f50e27)
Running command: ssh -l ec2-user -p 22 10.0.172.143 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o ConnectTimeout=3
Warning: Permanently added '10.10.5.100' (ECDSA) to the list of known hosts.
Last login: Sun Aug 16 17:01:52 2020 from ip-10-0-172-143.ap-southeast-1.compute.internal
.
.
.
[ec2-user@ip-10-0-172-143 ~]$
[ec2-user@ip-10-0-172-143 ~]$ logout
Connection to 10.0.172.143 closed.
```

### Select EC2 Instances with Tags
```bash
$ awssh --tags "Environment=production,ManagedBy=kops"

Use the arrow keys to navigate: ↓ ↑ → ←  and / toggles search
Select an instance:
  » master-1.masters.k8s.kops.internal i-0b2566fcc894c1bd1 (10.0.132.143)
    master-2.masters.k8s.kops.internal i-05c0309be99c8a097 (10.0.148.154)
    master-3.masters.k8s.kops.internal i-052913ef86123d500 (10.0.186.183)
    nodes-a.nodes.k8s.kops.internal i-07fc020d8c7f50e27 (10.0.172.143)

nodes-a.nodes.k8s.kops.internal i-07fc020d8c7f50e27
Running command: ssh -l ec2-user -p 22 10.0.172.143 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o ConnectTimeout=3
Warning: Permanently added '10.10.5.100' (ECDSA) to the list of known hosts.
Last login: Sun Aug 16 17:01:52 2020 from ip-10-0-172-143.ap-southeast-1.compute.internal
.
.
.
[ec2-user@ip-10-0-172-143 ~]$
[ec2-user@ip-10-0-172-143 ~]$ logout
Connection to 10.0.172.143 closed.
```

### Select EC2 Instances with InstanceID
```bash
$ awssh i-07fc020d8c7f50e27

Running command: ssh -l ec2-user -p 22 10.0.172.143 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o ConnectTimeout=3
Warning: Permanently added '10.10.5.100' (ECDSA) to the list of known hosts.
Last login: Sun Aug 16 17:01:52 2020 from ip-10-0-172-143.ap-southeast-1.compute.internal
.
.
.
[ec2-user@ip-10-0-172-143 ~]$
[ec2-user@ip-10-0-172-143 ~]$ logout
Connection to 10.0.172.143 closed.
```