# Usage: Litmusctl v0.2.0

> Notes:
>
> - For litmusctl v0.3.0 or earlier
> - Compatible with Litmus 2.0.0-Beta8 or earlier

### Connecting an agent

To connect Litmus Chaos agent:

```shell
litmusctl agent connect
```

Next, you need to enter ChaosCenter details to login into your ChaosCenter account. Fields to be filled in:

**ChaosCenter UI URL:** Enter the URL used to access the ChaosCenter UI.
Example, http://172.17.0.2:31696/

**Username:** Enter your ChaosCenter username.
**Password:** Enter your ChaosCenter password.

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
**Agent Name:** Enter the name of the new agent.

**Agent Description:** Fill in details about the agent.

**Platform Name:** Enter the platform name on which this agent is hosted. For example, AWS, GCP, Rancher etc.

**Enter the namespace:** You can either enter an existing namespace or enter a new namespace. In cases where the namespace does not exist, ChaosCenter creates it for you.

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

💡 Connecting agent to ChaosCenter.
🏃 Agents running!!
🚀 Agent Connection Successful!! 🎉
👉 Litmus agents can be accessed here: http://172.17.0.2:31696/targets
```

To verify, if the connection process was successful you can view the list of connected agents from the Targets section on your ChaosCenter and ensure that the connected agent is in Active State.
