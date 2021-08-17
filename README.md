# Litmusctl
[![BCH compliance](https://bettercodehub.com/edge/badge/litmuschaos/litmusctl?branch=master)](https://bettercodehub.com/)
![GitHub Workflow](https://github.com/litmuschaos/litmusctl/actions/workflows/push.yml/badge.svg?branch=master)
[![GitHub stars](https://img.shields.io/github/stars/litmuschaos/litmusctl?style=social)](https://github.com/litmuschaos/litmusctl/stargazers)
[![GitHub Release](https://img.shields.io/github/release/litmuschaos/litmusctl.svg?style=flat)]()

The Litmuschaos command-line tool, litmusctl, allows you to manage litmuschaos's agent plane. You can use litmusctl to create agents, project, and manage multiple litmuschaos accounts. 

## Usage
For more information including a complete list of litmusctl operations, see the litmusctl reference documentation.
* For v0.3.0 or latest  && compatible with Litmus-2.0.0-Beta9 or latest: 
  * Non-Interactive mode: <a href="https://github.com/litmuschaos/litmusctl/blob/master/Usage.md">Click here</a>
  * Interactive mode: <a href="https://github.com/litmuschaos/litmusctl/blob/master/Usage_interactive.md">Click here</a>
* For v0.2.0 or earlier && compatible with Litmus-2.0.0-Beta8 or earlier: <a href="https://github.com/litmuschaos/litmusctl/blob/master/Usage_v0.2.0.md">Click here</a>

## Requirements

The litmusctl CLI requires the following things:

- kubeconfig - litmusctl needs the kubeconfig of the k8s cluster where we need to connect litmus agents. The CLI currently uses the default path of kubeconfig i.e. `~/.kube/config`.
- kubectl- litmusctl is using kubectl under the hood to apply the manifest. To install kubectl, follow:  [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl)


## Installation

To install the latest version of litmusctl follow the below steps:

<table>
  <th>Platforms</th>
  <th>v0.3.0</th>
  <th>v0.2.0 </th>
  <th>master(Unreleased)</th>
  <tr>
    <td>litmusctl-darwin-amd64 (MacOS)</td>
    <td><a href="https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-darwin-amd64-v0.3.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-darwin-amd64-v0.2.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-darwin-amd64-master.tar.gz">Click here</a></td>
  </tr>
  <tr>
    <td>litmusctl-linux-386</td>
    <td><a href="https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-linux-386-v0.3.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-linux-386-v0.2.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-linux-386-master.tar.gz">Click here</a></td>
  </tr>
  <tr>
    <td>litmusctl-linux-amd64</td>
    <td><a href="https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-linux-amd64-v0.3.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-linux-amd64-v0.2.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-linux-amd64-master.tar.gz">Click here</a></td>
  </tr>
  <tr>
    <td>litmusctl-linux-arm</td>
    <td><a href="https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-linux-arm-v0.3.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-linux-arm-v0.2.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-linux-arm-master.tar.gz">Click here</a></td>
  </tr>
  <tr>
    <td>litmusctl-linux-arm64</td>
    <td><a href="https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-linux-arm64-v0.3.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-linux-arm64-v0.2.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-linux-arm64-master.tar.gz">Click here</a></td>
  </tr>
  <tr>
    <td>litmusctl-windows-386</td>
    <td><a href="https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-windows-386-v0.3.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-windows-386-v0.2.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-windows-386-master.tar.gz">Click here</a></td>
  </tr>
   <tr>
    <td>litmusctl-windows-amd64</td>
    <td><a href="https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-windows-amd64-v0.3.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-windows-amd64-v0.2.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-windows-amd64-master.tar.gz">Click here</a></td>
  </tr>
  <tr>
    <td>litmusctl-windows-arm</td>
    <td><a href="https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-windows-arm-v0.3.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-windows-arm-v0.2.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-windows-arm-master.tar.gz">Click here</a></td>
  </tr>
</table>

### Linux/MacOS

* Extract the binary

```shell
tar -zxvf litmusctl-<OS>-<ARCH>-<VERSION>.tar.gz
```

* Provide necessary permissions

```shell
chmod +x litmusctl
```

* Move the litmusctl binary to /usr/local/bin/litmusctl. Note: Make sure to use root user or use sudo as a prefix

```shell
mv litmusctl /usr/local/bin/litmusctl
```

* You can run the litmusctl command in Linux/macOS:

```shell
litmusctl <command> <subcommand> <subcommand> [options and parameters]
```

### Windows

* Extract the binary from the zip using WinZip or any other extraction tool.

* You can run the litmusctl command in windows:

```shell
litmusctl.exe <command> <subcommand> <subcommand> [options and parameters]
```

* To check the version of the litmusctl:

```shell
litmusctl version
```

----
