# kuberactl

This package provides command line interface to connect agents to Kubera services.

## Requirements

The kuberactl CLI requires the following things:

- Kubeconfig - Kuberactl needs the kubeconfig of the k8s cluster where we need to connect Kubera agents. By default the CLI uses the kubeconfig present in `~/.kube/config`. If the kubeconfig is present in some other location or if we want to use some different kubeconfig altogether, it can passed in the command via `--kubeconfig` flag.

## Installation

**Linux**

To install the latest version of kuberactl CLI follow the below steps:

- Download the latest kuberactl binary from - `https://github.com/mayadata-io/kuberactl/releases/`
- Untar the binary
```shell
$ tar -xvzf kuberactl_<Version>_<OS>_<ARCH>.tar.gz
```
- Move the kuberactl binary to /usr/local/bin
```shell
$ sudo mv kuberactl /usr/local/bin/
```

## Basic Commands

Kuberactl CLI command has the following structure:
```shell
$ kuberactl <command> <subcommand> <subcommand> [options and parameters]
```

To get the version of the kuberactl CLI:
```shell
$ kuberactl version
```

To register Kubera Chaos agent:
```shell
$ kuberactl chaos agent register
```