# kuberactl

This package provides command line interface to connect agents to Kubera services.

## Requirements

The kuberactl CLI requires the following things:

- Kubeconfig - Kuberactl needs the kubeconfig of the k8s cluster where we need to connect Kubera agents. The CLI currently uses the default path of kubeconfig i.e. `~/.kube/config`.

## Installation

**Linux**

To install the latest version of kuberactl CLI follow the below steps:

- Download the latest kuberactl binary from - `http://asset.mayadata.io/kuberactl/latest/kuberactl_latest_Linux_x86_64.tar.gz`
- Untar the binary
```shell
$ tar -xvzf kuberactl_latest_Linux_x86_64.tar.gz
```
- Move the kuberactl binary to /usr/local/bin
```shell
$ sudo mv kuberactl /usr/local/bin/
```

> NOTE: Kuberactl binaries for master and development branches are available in .zip format and can be downloaded from - `http://asset.mayadata.io/kuberactl/<branch-name>/kuberactl_linux_amd64.zip`, `http://asset.mayadata.io/kuberactl/<branch-name>/kuberactl_windows_amd64.zip` and `http://asset.mayadata.io/kuberactl/<branch-name>/kuberactl_darwin_amd64.zip` respectively for Linux, Windows and Darwin.

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
To register Kubera Propel agent:
```shell
$ kuberactl propel agent register
```

## Development workflow

#### Initial setup

**Fork in the cloud**
1. Visit https://github.com/mayadata-io/kuberactl.
2. Click `Fork` button (top right) to establish a cloud-based fork.

**Clone fork to local host**
Place mayadata-io/kuberactl code in any directory using the following cloning procedure -

```
mkdir path/to/directory/mayadata-io
cd mayadata-io

# Note: Here $user is your GitHub profile name
git clone https://github.com/$user/kuberactl.git

# Configure remote upstream
cd path/to/directory/mayadata-io/kuberactl
git remote add upstream https://github.com/mayadata-io/kuberactl.git

# Never push to upstream master
git remote set-url --push upstream no_push

# Confirm that your remotes make sense
git remote -v
```

**Create a new feature branch to work on your issue**
```
$ git checkout -b <branch-name>
Switched to a new branch '<branch-name>'
```

**Make your changes and test them**
Once the changes are done, you can build the binary and test them using the following command:
```
Get into the cloned directory
$ cd path/to/directory/mayadata-io/kuberactl

$ go install
Note: This will build a binary at `~/go/bin`
```
Test your changes by running the necessary `kuberactl` commands.