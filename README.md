# litmusctl

This package provides command line interface to connect agents to Litmus services.

## Requirements

The litmusctl CLI requires the following things:

- Kubeconfig - litmusctl needs the kubeconfig of the k8s cluster where we need to connect litmus agents. The CLI currently uses the default path of kubeconfig i.e. `~/.kube/config`.

## Installation

**Linux**

To install the latest version of litmusctl CLI follow the below steps:

- Download the latest litmusctl binary from - `http://asset.mayadata.io/litmusctl/latest/litmusctl_latest_Linux_x86_64.tar.gz`
- Untar the binary

```shell
$ tar -xvzf litmusctl_latest_Linux_x86_64.tar.gz
```

- Move the litmusctl binary to /usr/local/bin

```shell
$ sudo mv litmusctl /usr/local/bin/
```

> NOTE: litmusctl binaries for master and development branches are available in .zip format and can be downloaded from - `http://asset.mayadata.io/litmusctl/<branch-name>/litmusctl_linux_amd64.zip`, `http://asset.mayadata.io/litmusctl/<branch-name>/litmusctl_windows_amd64.zip` and `http://asset.mayadata.io/litmusctl/<branch-name>/litmusctl_darwin_amd64.zip` respectively for Linux, Windows and Darwin.

## Basic Commands

litmusctl CLI command has the following structure:

```shell
$ litmusctl <command> <subcommand> <subcommand> [options and parameters]
```

To get the version of the litmusctl CLI:

```shell
$ litmusctl version
```

To register Litmus Chaos agent:

```shell
$ litmusctl chaos agent register
```

To register Litmus Propel agent:

```shell
$ litmusctl propel agent register
```

## Development workflow

#### Initial setup

**Fork in the cloud**

1. Visit https://github.com/litmuschaos/litmusctl.
2. Click `Fork` button (top right) to establish a cloud-based fork.

**Clone fork to local host**
Place mayadata-io/litmusctl code in any directory using the following cloning procedure -

```
mkdir path/to/directory/mayadata-io
cd mayadata-io

# Note: Here $user is your GitHub profile name
git clone https://github.com/$user/litmusctl.git

# Configure remote upstream
cd path/to/directory/mayadata-io/litmusctl
git remote add upstream https://github.com/litmuschaos/litmusctl.git

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
$ cd path/to/directory/mayadata-io/litmusctl

$ go install
Note: This will build a binary at `~/go/bin`
```

Test your changes by running the necessary `litmusctl` commands.
