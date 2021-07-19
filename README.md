# Litmusctl

![GitHub Workflow](https://github.com/litmuschaos/litmusctl/actions/workflows/push.yml/badge.svg?branch=master)
[![GitHub stars](https://img.shields.io/github/stars/litmuschaos/litmusctl?style=social)](https://github.com/litmuschaos/litmusctl/stargazers)
[![GitHub Release](https://img.shields.io/github/release/litmuschaos/litmusctl.svg?style=flat)]()  

Litmusctl is a command line interface to manage LitmusPortal services.

## Requirements

The litmusctl CLI requires the following things:

- Kubeconfig - litmusctl needs the kubeconfig of the k8s cluster where we need to connect litmus agents. The CLI currently uses the default path of kubeconfig i.e. `~/.kube/config`.

## Installation

To install the latest version of litmusctl follow the below steps:

- Download the latest litmusctl(master) binary from:

| Platforms                      | Download Link                                                                                           |
| ------------------------------ | ------------------------------------------------------------------------------------------------------- |
| litmusctl-darwin-386 (MacOS)   | [Click here](https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-darwin-386-master.tar.gz)    |
| litmusctl-darwin-amd64 (MacOS) | [Click here](https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-darwin-amd64-master.tar.gz)  |
| litmusctl-linux-386            | [Click here](https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-linux-386-master.tar.gz)     |
| litmusctl-linux-amd64          | [Click here](https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-linux-amd64-master.tar.gz)   |
| litmusctl-linux-arm            | [Click here](https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-linux-arm-master.tar.gz)     |
| litmusctl-linux-arm64          | [Click here](https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-linux-arm64-master.tar.gz)   |
| litmusctl-windows-386          | [Click here](https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-windows-386-master.tar.gz)   |
| litmusctl-windows-amd64        | [Click here](https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-windows-amd64-master.tar.gz) |
| litmusctl-windows-arm          | [Click here](https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-windows-arm-master.tar.gz)   |

<br>

- Download the litmusctl(v0.2.0) binary from:

| Platforms                      | Download Link                                                                                           |
| ------------------------------ | ------------------------------------------------------------------------------------------------------- |
| litmusctl-darwin-386 (MacOS)   | [Click here](https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-darwin-386-v0.2.0.tar.gz)    |
| litmusctl-darwin-amd64 (MacOS) | [Click here](https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-darwin-amd64-v0.2.0.tar.gz)  |
| litmusctl-linux-386            | [Click here](https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-linux-386-v0.2.0.tar.gz)     |
| litmusctl-linux-amd64          | [Click here](https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-linux-amd64-v0.2.0.tar.gz)   |
| litmusctl-linux-arm            | [Click here](https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-linux-arm-v0.2.0.tar.gz)     |
| litmusctl-linux-arm64          | [Click here](https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-linux-arm64-v0.2.0.tar.gz)   |
| litmusctl-windows-386          | [Click here](https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-windows-386-v0.2.0.tar.gz)   |
| litmusctl-windows-amd64        | [Click here](https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-windows-amd64-v0.2.0.tar.gz) |
| litmusctl-windows-arm          | [Click here](https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-windows-arm-v0.2.0.tar.gz)   |

<br>

- Extract the binary

```shell
tar -zxvf litmusctl-<OS>-<ARCH>-<VERSION>.tar.gz
```

- Provide necessary permissions

```shell
chmod +x litmusctl
```

- Move the litmusctl binary to /usr/local/bin/litmusctl

```shell
sudo mv litmusctl /usr/local/bin/litmusctl
```

## Basic Commands

litmusctl CLI command has the following structure:

```shell
$ litmusctl <command> <subcommand> <subcommand> [options and parameters]
```

To get the version of the litmusctl CLI:

```shell
$ litmusctl version
```
