# Litmusctl

![GitHub Workflow](https://github.com/litmuschaos/litmusctl/actions/workflows/push.yml/badge.svg?branch=master)
[![GitHub stars](https://img.shields.io/github/stars/litmuschaos/litmusctl?style=social)](https://github.com/litmuschaos/litmusctl/stargazers)
[![GitHub Release](https://img.shields.io/github/release/litmuschaos/litmusctl.svg?style=flat)]()

The Litmuschaos command-line tool, litmusctl, allows you to manage litmuschaos's agent plane. You can use litmusctl to connect Chaos Delegates, create project, schedule Chaos Scenarios, disconnect Chaos Delegates and manage multiple litmuschaos accounts.

## Usage

For more information including a complete list of litmusctl operations, see the litmusctl reference documentation.

* For 0.23.0 or latest: <a href="https://github.com/litmuschaos/litmusctl/blob/master/Usage_0.23.0.md">Click here</a>
* For v0.12.0 to v0.22.0: <a href="https://github.com/litmuschaos/litmusctl/blob/master/Usage_interactive.md">Click here</a>
* For v0.2.0 or earlier && compatible with Litmus-2.0.0-Beta8 or earlier: <a href="https://github.com/litmuschaos/litmusctl/blob/master/Usage_v0.2.0.md">Click here</a>

## Requirements

The litmusctl CLI requires the following things:

- kubeconfig - litmusctl needs the kubeconfig of the k8s cluster where we need to connect litmus Chaos Delegates. The CLI currently uses the default path of kubeconfig i.e. `~/.kube/config`.
- kubectl- litmusctl is using kubectl under the hood to apply the manifest. To install kubectl, follow: [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl)

## Compatibility matrix

To check compatibility of litmusctl with Chaos Center

<table>
  <th>litmusctl version</th>
  <th>Lowest Chaos Center supported version</th>
  <th>Highest Chaos Center supported version</th>
 <tr>
    <td>1.19.0</td>
    <td>3.0.0</td>
    <td>3.22.0</td>
 </tr>
 <tr>
    <td>1.18.0</td>
    <td>3.0.0</td>
    <td>3.21.0</td>
 </tr>
 <tr>
    <td>1.17.0</td>
    <td>3.0.0</td>
    <td>3.20.0</td>
 </tr>
 <tr>
    <td>1.16.0</td>
    <td>3.0.0</td>
    <td>3.19.0</td>
 </tr>
 <tr>
    <td>1.15.0</td>
    <td>3.0.0</td>
    <td>3.18.0</td>
 </tr>
 <tr>
    <td>1.14.0</td>
    <td>3.0.0</td>
    <td>3.15.0</td>
 </tr>
 <tr>
    <td>1.13.0</td>
    <td>3.0.0</td>
    <td>3.14.0</td>
 </tr>
 <tr>
    <td>1.12.0</td>
    <td>3.0.0</td>
    <td>3.13.0</td>
 </tr>
 <tr>
    <td>1.11.0</td>
    <td>3.0.0</td>
    <td>3.12.0</td>
 </tr>
 <tr>
    <td>1.10.0</td>
    <td>3.0.0</td>
    <td>3.11.0</td>
 </tr>
 <tr>
    <td>1.9.0</td>
    <td>3.0.0</td>
    <td>3.10.0</td>
 </tr>
 <tr>
    <td>1.8.0</td>
    <td>3.0.0</td>
    <td>3.9.1</td>
 </tr>
 <tr>
    <td>1.7.0</td>
    <td>3.0.0</td>
    <td>3.8.0</td>
 </tr>
 <tr>
    <td>1.6.0</td>
    <td>3.0.0</td>
    <td>3.7.0</td>
 </tr>
 <tr>
    <td>1.5.0</td>
    <td>3.0.0</td>
    <td>3.6.0</td>
 </tr>
 <tr>
    <td>1.4.0</td>
    <td>3.0.0</td>
    <td>3.5.0</td>
 </tr>
 <tr>
    <td>1.3.0</td>
    <td>3.0.0</td>
    <td>3.4.0</td>
  </tr>
 <tr>
    <td>1.2.0</td>
    <td>3.0.0</td>
    <td>3.3.0</td>
  </tr> 
 <tr>
    <td>1.1.0</td>
    <td>3.0.0</td>
    <td>3.2.0</td>
  </tr> 
 <tr>
    <td>1.0.0</td>
    <td>3.0.0</td>
    <td>3.1.0</td>
  </tr>
</table>

## Installation

To install the latest version of litmusctl follow the below steps:

<table>
  <th>Platforms</th>
  <th>1.19.0</th>
  <th>1.18.0</th>
  <th>1.17.0</th>
  <th>1.16.0</th>
  <th>1.15.0</th>
  <th>1.14.0</th>
  <th>1.13.0</th>
  <th>1.12.0</th>
  <th>master(Unreleased)</th>
  <tr>
    <td>litmusctl-darwin-amd64 (MacOS)</td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-darwin-amd64-1.19.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-darwin-amd64-1.18.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-darwin-amd64-1.17.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-darwin-amd64-1.16.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-darwin-amd64-1.15.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-darwin-amd64-1.14.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-darwin-amd64-1.13.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-darwin-amd64-1.12.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-darwin-amd64-master.tar.gz">Click here</a></td>
  </tr>
  <tr>
    <td>litmusctl-linux-386</td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-386-1.19.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-386-1.18.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-386-1.17.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-386-1.16.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-386-1.15.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-386-1.14.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-386-1.13.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-386-1.12.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-386-master.tar.gz">Click here</a></td>
  </tr>
  <tr>
    <td>litmusctl-linux-amd64</td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-amd64-1.19.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-amd64-1.18.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-amd64-1.17.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-amd64-1.16.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-amd64-1.15.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-amd64-1.14.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-amd64-1.13.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-amd64-1.12.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-amd64-master.tar.gz">Click here</a></td>
  </tr>
  <tr>
    <td>litmusctl-linux-arm</td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-arm-1.19.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-arm-1.18.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-arm-1.17.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-arm-1.16.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-arm-1.15.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-arm-1.14.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-arm-1.13.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-arm-1.12.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-arm-master.tar.gz">Click here</a></td>
  </tr>
  <tr>
    <td>litmusctl-linux-arm64</td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-arm64-1.19.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-arm64-1.18.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-arm64-1.17.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-arm64-1.16.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-arm64-1.15.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-arm64-1.14.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-arm64-1.13.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-arm64-1.12.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-arm64-master.tar.gz">Click here</a></td>
  </tr>
  <tr>
    <td>litmusctl-windows-386</td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-windows-386-1.19.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-windows-386-1.18.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-windows-386-1.17.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-windows-386-1.16.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-windows-386-1.15.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-windows-386-1.14.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-windows-386-1.13.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-windows-386-1.12.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-windows-386-master.tar.gz">Click here</a></td>
  </tr>
   <tr>
    <td>litmusctl-windows-amd64</td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-windows-amd64-1.19.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-windows-amd64-1.18.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-windows-amd64-1.17.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-windows-amd64-1.16.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-windows-amd64-1.15.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-windows-amd64-1.14.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-windows-amd64-1.13.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-windows-amd64-1.12.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-windows-amd64-master.tar.gz">Click here</a></td>
  </tr>
  <tr>
    <td>litmusctl-windows-arm</td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-windows-arm-1.19.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-windows-arm-1.18.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-windows-arm-1.17.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-windows-arm-1.16.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-windows-arm-1.15.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-windows-arm-1.14.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-windows-arm-1.13.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-windows-arm-1.12.0.tar.gz">Click here</a></td>
    <td><a href="https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-windows-arm-master.tar.gz">Click here</a></td>
  </tr>
</table>

### Linux/MacOS

- Extract the binary

```shell
tar -zxvf litmusctl-<OS>-<ARCH>-<VERSION>.tar.gz
```

- Provide necessary permissions

```shell
chmod +x litmusctl
```

- Move the litmusctl binary to /usr/local/bin/litmusctl. Note: Make sure to use root user or use sudo as a prefix

```shell
mv litmusctl /usr/local/bin/litmusctl
```

- You can run the litmusctl command in Linux/macOS:

```shell
litmusctl <command> <subcommand> <subcommand> [options and parameters]
```

### Windows

- Extract the binary from the zip using WinZip or any other extraction tool.

- You can run the litmusctl command in windows:

```shell
litmusctl.exe <command> <subcommand> <subcommand> [options and parameters]
```

- To check the version of the litmusctl:

```shell
litmusctl version
```

## Verbose Mode

litmusctl supports a `--verbose` (or `-v`) flag to the enable detailed logging during command execution. This is useful for debugging and understanding what's happening behind the scenes.

### Usage

Add the `--verbose` or `-v` flag to any litmusctl command:

```shell
litmusctl <command> --verbose
# or
litmusctl <command> -v
```

### What Gets Logged

When verbose mode is enabled, litmusctl will print additional information including:

- API endpoints being accessed
- HTTP request methods and payloads
- Response status codes
- Internal processing steps
- Project and resource details during operations

### Examples

```shell
# Get projects with verbose logging
litmusctl get projects --verbose

# Create a project with detailed logs
litmusctl create project --name my-project -v

# Connect chaos infrastructure with debug information
litmusctl connect chaos-infra --verbose --name my-infra --non-interactive
```

### Default Behavior

By default, litmusctl operates in minimal logging mode, showing only:
- Error messages
- Warning messages
- User-facing output

Use `--verbose` when you need to troubleshoot issues or understand the internal operations of the CLI.

## Development Guide

You can find the local setup guide for **`litmusctl`** [here](DEVELOPMENT.md).

---
