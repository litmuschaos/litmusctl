# Litmusctl

Litmusctl is a command line interface to manage LitmusPortal services.

## Requirements

The litmusctl CLI requires the following things:

- Kubeconfig - litmusctl needs the kubeconfig of the k8s cluster where we need to connect litmus agents. The CLI currently uses the default path of kubeconfig i.e. `~/.kube/config`.

## Installation

To install the latest version of litmusctl follow the below steps:

- Download the stable litmusctl(v0.1.0) binary from:

| Platforms                      | Download Link                                                                                           |
| ------------------------------ | ------------------------------------------------------------------------------------------------------- |
| litmusctl-darwin-386 (MacOS)   | [Click here](https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-darwin-386-v0.1.0.tar.gz)    |
| litmusctl-darwin-amd64 (MacOS) | [Click here](https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-darwin-amd64-v0.1.0.tar.gz)  |
| litmusctl-linux-386            | [Click here](https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-linux-386-v0.1.0.tar.gz)     |
| litmusctl-linux-amd64          | [Click here](https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-linux-amd64-v0.1.0.tar.gz)   |
| litmusctl-linux-arm            | [Click here](https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-linux-arm-v0.1.0.tar.gz)     |
| litmusctl-linux-arm64          | [Click here](https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-linux-arm64-v0.1.0.tar.gz)   |
| litmusctl-windows-386          | [Click here](https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-windows-386-v0.1.0.tar.gz)   |
| litmusctl-windows-amd64        | [Click here](https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-windows-amd64-v0.1.0.tar.gz) |
| litmusctl-windows-arm          | [Click here](https://litmusctl-bucket.s3-eu-west-1.amazonaws.com/litmusctl-windows-arm-v0.1.0.tar.gz)   |

<br>

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

- Extract the binary

```shell
$ tar -zxvf litmusctl-<OS>-<ARCH>-<VERSION>.tar.gz
```

- Provide necessary permissions

```shell
$ chmod +x litmusctl
```

- Move the litmusctl binary to /usr/local/bin/litmusctl

```shell
$ sudo mv litmusctl /usr/local/bin/litmusctl
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

### Connecting an agent

To connect Litmus Chaos agent:

```shell
$ litmusctl agent connect
```

Next, you need to enter LitmusPortal details to login into your LitmusPortal account. Fields to be filled in:

**LimtusPortal UI URL:** Enter the URL used to access the Litmus Portal UI.
Example, http://172.17.0.2:31696/

**Username:** Enter your LitmusPortal username.
**Password:** Enter your LitmusPortal password.

```shell
🔥 Connecting LitmusChaos agent

📶 Please enter LitmusChaos details --
👉 Host URL where litmus is installed: http://172.17.0.2:31696/
🤔 Username [admin]: admin
🙈 Password:
✅ Login Successful!
```

Upon successful login, there will be a list of exiting projects displayed on the terminal. Select the desired project by entering the sequence number indicated against it.

```shell
✨ Projects List:
1.  abc

🔎 Select Project: 1
```

Next, select the installation mode. In case the selected mode was a Cluster there will be a prerequisites check to verify ClusterRole and ClusterRoleBinding.

```shell
🔌 Installation Modes:
1. Cluster
2. Namespace

👉 Select Mode [cluster]: 1

🏃 Running prerequisites check....
🔑  clusterrole - ✅
🔑  clusterrolebinding - ✅

🌟 Sufficient permissions. Connecting Agent
```

Next, enter the details of the new agent.

Fields to filled in:
**Agent Name:** Enter the name for the new agent.

**Agent Description:** Fill in details about the agent.

**Platform Name:** Enter the platform name on which this agent is hosted. For example, AWS, GCP, Rancher etc.

**Enter the namespace:** You can either enter an existing namespace or enter a new namespace. In cases where the namespace does not exist, LimtusPortal creates it for you.

**Enter service account:** Enter a name for your service account.

```shell
🔗 Enter the details of the agent ----
🤷 Agent Name: my-agent
📘 Agent Description: This is a new agent.
📦 Platform List
1. AWS
2. GKE
3. Openshift
4. Rancher
5. Others
🔎 Select Platform [Others]: 5
📁 Enter the namespace (new or existing) [litmus]: litmus
🔑 Enter service account [litmus]: litmus
```

Once, all these steps are implemented you will be able to see a summary of all the entered fields.
After verification of these details, you can proceed with the connection of the agent by entering Y. The process of connection might take up to a few seconds.

```shell
📌 Summary --------------------------

Agent Name:         my-agent
Agent Description:  This is a new agent.
Platform Name:      Others
Namespace:          litmus
Service Account:    litmus
Installation Mode:  cluster

-------------------------------------

🤷 Do you want to continue with the above details? [Y/N]: Y

💡 Connecting agent to Litmus Portal.
🏃 Agents running!!
🚀 Agent Connection Successful!! 🎉
👉 Litmus agents can be accessed here: http://172.17.0.2:31696/targets
```

To verify, if the connection process was successful you can view the list of connected agents from the Targets section on your LitmusPortal and ensure that the connected agent is in Active State.
