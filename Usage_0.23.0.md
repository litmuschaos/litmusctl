> Notes:
>
> - For litmusctl v0.23.0 or latest

### litmusctl Syntax

`litmusctl` has a syntax to use as follows:

```shell
litmusctl [command] [TYPE] [flags]
```

- Command: refers to what you do want to perform (connect, create, get and config)
- Type: refers to the feature type you are performing a command against (chaos-infra, project etc.)
- Flags: It takes some additional information for resource operations. For example, `--installation-mode` allows you to specify an installation mode.

Litmusctl is using the `.litmusconfig` config file to manage multiple accounts

1. If the --config flag is set, then only the given file is loaded. The flag may only be set once and no merging takes place.
2. Otherwise, the ${HOME}/.litmusconfig file is used, and no merging takes place.

Litmusctl supports both interactive and non-interactive(flag based) modes.

> Only `litmusctl connect chaos-infra` command needs --non-interactive flag, other commands don't need this flag to be in non-interactive mode. If mandatory flags aren't passed, then litmusctl takes input in an interactive mode.

### Steps to connect a Chaos Infrastucture

- To setup an account with litmusctl

```shell
litmusctl config set-account
```

Next, you need to enter ChaosCenter details to login into your ChaosCenter account. Fields to be filled in:

**ChaosCenter URL:** Enter the URL used to access the ChaosCenter.

> Example, https://preview.litmuschaos.io/

**Username:** Enter your ChaosCenter username.
**Password:** Enter your ChaosCenter password.

```
Host endpoint where litmus is installed: https://preview.litmuschaos.io/
Username [Default: admin]: admin

Password:
account.username/admin configured
```

- To connect a Chaos Infrastructure in a cluster mode

```shell
litmusctl connect chaos-infra
```

There will be a list of existing projects displayed on the terminal. Select the desired project by entering the sequence number indicated against it.

```
Project list:
1.  Project-Admin

Select a project [Range: 1-1]: 1
```

Next, select the installation mode based on your requirement by entering the sequence number indicated against it.

Litmusctl can install a Chaos Infrastructure in two different modes.

- cluster mode: With this mode, the Chaos Infrastructure can run the chaos in any namespace. It installs appropriate cluster roles and cluster role bindings to achieve this mode.

- namespace mode: With this mode, the Chaos Infrastructure can run the chaos in its namespace. It installs appropriate roles and role bindings to achieve this mode.

Note: With namespace mode, the user needs to create the namespace to install the Chaos Infrastructure as a prerequisite.

```
Installation Modes:
1. Cluster
2. Namespace

Select Mode [Default: cluster] [Range: 1-2]: 1

🏃 Running prerequisites check....
🔑 clusterrole ✅
🔑 clusterrolebinding ✅
🌟 Sufficient permissions. Installing the Chaos Infrastructure...

```

Next, enter the details of the new Chaos infrastructure.

Fields to be filled in <br />

<table>
    <th>Field</th>
    <th>Description</th>
    <tr>
        <td>Chaos Infrastructure Name:</td>
        <td>Enter a name of the Chaos Infrastructure which needs to be unique across the project</td>
    </tr>
    <tr>
        <td>Chaos Infrastructure Description:</td>
        <td>Fill in details about the Chaos Infrastructure</td>
    </tr>
    <tr>
        <td>Chaos EnvironmentID :</td>
        <td>Fill in details about the Chaos Environment ID. The Environment Should be already existing.</td>
    </tr>
    <tr>
        <td>Skip SSL verification</td>
        <td>Choose whether Chaos Infrastructure will skip SSL/TLS verification</td>
    </tr>
    <tr>
        <td>Node Selector:</td>
        <td>To deploy the Chaos Infrastructure on a particular node based on the node selector labels</td>
    </tr>
    <tr>
        <td>Platform Name:</td>
        <td>Enter the platform name on which this Chaos Infrastructure is hosted. For example, AWS, GCP, Rancher etc.</td>
    </tr>
    <tr>
        <td>Enter the namespace:</td>
        <td>You can either enter an existing namespace or enter a new namespace. In cases where the namespace does not exist, litmusctl creates it for you</td>
    </tr>
    <tr>
        <td>Enter service account:</td>
        <td>You can either enter an existing or new service account</td>
    </tr>
</table>

```
Enter the details of the Chaos Infrastructure:

Chaos Infrastructure Name: New-Chaos-infrastructure

Chaos Infrastructure Description: This is a new Chaos Infrastructure

Chaos EnvironmentID: test-infra-environment

Do you want Chaos Infrastructure to skip SSL/TLS check (Y/N) (Default: N): n

Do you want NodeSelector to be added in the Chaos Infrastructure deployments (Y/N) (Default: N): N

Platform List:
1. AWS
2. GKE
3. Openshift
4. Rancher
5. Others

Select a platform [Default: Others] [Range: 1-5]: 5

Enter the namespace (new or existing namespace) [Default: litmus]:
👍 Continuing with litmus namespace
```

Once, all these steps are implemented you will be able to see a summary of all the entered fields.
After verification of these details, you can proceed with the connection of the Chaos infra by entering Y. The process of connection might take up to a few seconds.

```
Enter service account [Default: litmus]:

📌 Summary
Chaos Infra Name: test4
Chaos EnvironmentID: test
Chaos Infra Description:
Chaos Infra SSL/TLS Skip: false
Platform Name: Others
Namespace:  litmuwrq (new)
Service Account:  litmus (new)


Installation Mode: cluster

🤷 Do you want to continue with the above details? [Y/N]: Y
👍 Continuing Chaos Infrastructure connection!!
Applying YAML:
https://preview.litmuschaos.io/api/file/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbHVzdGVyX2lkIjoiMDUyZmFlN2UtZGM0MS00YmU4LWJiYTgtMmM4ZTYyNDFkN2I0In0.i31QQDG92X5nD6P_-7TfeAAarZqLvUTFfnAghJYXPiM.yaml

💡 Connecting Chaos Infrastructure to ChaosCenter.
🏃 Chaos Infrastructure is running!!

🚀 Chaos Infrastructure Connection Successful!! 🎉
👉 Litmus Chaos Infrastructure can be accessed here: https://preview.litmuschaos.io/targets
```

#### Verify the new Chaos Infrastructure Connection\*\*

To verify, if the connection process was successful you can view the list of connected Chaos Infrastructures from the Targets section on your ChaosCenter and ensure that the connected Chaos Infrastructure is in Active State.

---

### Steps to create a Chaos Experiment

- To setup an account with litmusctl

```shell
litmusctl config set-account --endpoint="" --username="" --password=""
```

- To create a Chaos Experiment by passing a manifest file
  > Note:
  >
  > - To get `project-id`, apply `litmusctl get projects`
  > - To get `chaos-infra-id`, apply `litmusctl get chaos-infra --project-id=""`

```shell
litmusctl create chaos-experiment -f custom-chaos-experiment.yml --project-id="" --chaos-infra-id=""
```

- To Save the Chaos Experiment:

```shell
litmusctl save chaos-experiment -f custom-litmus-experiment.yaml
```

> Note:
>
> - Experiment Name can also be passed through the Manifest file

```shell
Enter the Project ID: eb7fc0a0-5878-4454-a9db-b67d283713bc
Enter the Chaos Infra ID: e7eb0386-085c-49c2-b550-8d85b58fd
Experiment Description:

🚀 Chaos Experiment/experiment-1 successfully created 🎉
```

- To Run a chaos Experiment:

```shell
litmusctl run chaos-experiment

Enter the Project ID: eb7fc0a0-5878-4454-a9db-b67d283713bc

Enter the Chaos Experiment ID: test_exp

🚀 Chaos Experiment running successfully 🎉
```

### Additional commands

- To view the current configuration of `.litmusconfig`, type:

```shell
litmusctl config view
```

**Output:**

```
accounts:
- users:
  - expires_in: "1626897027"
    token: eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjY4OTcwMjcsInJvbGUiOiJhZG1pbiIsInVpZCI6ImVlODZkYTljLTNmODAtNGRmMy04YzQyLTExNzlhODIzOTVhOSIsInVzZXJuYW1lIjoiYWRtaW4ifQ.O_hFcIhxP4rhyUN9NEVlQmWesoWlpgHpPFL58VbJHnhvJllP5_MNPbrRMKyFvzW3hANgXK2u8437u
    username: admin
  - expires_in: "1626944602"
    token: eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjY5NDQ2MDIsInJvbGUiOiJ1c2VyIiwidWlkIjoiNjFmMDY4M2YtZWY0OC00MGE1LWIzMjgtZTU2ZDA2NjM1MTE4IiwidXNlcm5hbWUiOiJyYWoifQ.pks7xjkFdJD649RjCBwQuPF1_QMoryDWixSKx4tPAqXI75ns4sc-yGhMdbEvIZ3AJSvDaqTa47XTC6c8R
    username: litmus-user
  endpoint: https://preview.litmuschaos.io
  serverEndpoint: https://preview.litmuschaos.io
apiVersion: v1
current-account: https://preview.litmuschaos.io
current-user: litmus-user
kind: Config
```

- To get an overview of the accounts available within `.litmusconfig`, use the `config get-accounts` command:

```shell
litmusctl config get-accounts
```

**Output:**

```
CURRENT  ENDPOINT                         USERNAME  EXPIRESIN
         https://preview.litmuschaos.io   admin     2021-07-22 01:20:27 +0530 IST
*        https://preview.litmuschaos.io   raj       2021-07-22 14:33:22 +0530 IST
```

- To alter the current account use the `use-account` command:

```shell
litmusctl config use-account

Host endpoint where litmus is installed: https://preview.litmuschaos.io

Username: admin

✅ Successfully set the current account to 'account-name' at 'URL'
```

- To create a project, apply the following command :

```shell
litmusctl create project

Enter a project name: new

Project 'project-name' created successfully!🎉
```

- To create a new Environment, apply the following command :

```shell
litmusctl create environment

Enter the Project ID: eb7fc0a0-5878-4454-a9db-b67d283713bc

Enter the Environment Name: test2

🚀 New Chaos Environment creation successful!! 🎉
```

- To view all the projects with the user, use the `get projects` command.

```shell
litmusctl get projects
```

**Output:**

```
PROJECT ID                                PROJECT NAME       CREATEDAT
50addd40-8767-448c-a91a-5071543a2d8e      Developer Project  2021-07-21 14:38:51 +0530 IST
7a4a259a-1ae5-4204-ae83-89a8838eaec3      DevOps Project     2021-07-21 14:39:14 +0530 IST
Press Enter to show the next page (or type 'q' to quit): q
```

- To get an overview of the Chaos Infrastructures available within a project, issue the following command.

```shell
litmusctl get chaos-infra

Enter the Project ID: 50addd40-8767-448c-a91a-5071543a2d8e
```

**Output:**

```
CHAOS Infrastructure ID                      CHAOS Infrastructure NAME    STATUS
55ecc7f2-2754-43aa-8e12-6903e4c6183a   chaos-infra-1            ACTIVE
13dsf3d1-5324-54af-4g23-5331g5v2364f   chaos-infra-2            INACTIVE
```

- To disconnect an Chaos Infrastructure, issue the following command..

```shell
litmusctl disconnect chaos-infra <chaos-infra-id> --project-id=""
```

**Output:**

```
🚀 Chaos Infrastructure successfully disconnected.
```

- To list the created Chaos Experiments within a project, issue the following command.

Using Flag :

```shell
litmusctl get chaos-experiment --project-id=""
```

Using UI :

```shell
Enter the Project ID: "project-id"
Select an output format:
table
json
yaml
```

**Output:**

```
    CHAOS Experiment ID                         CHAOS Experiment NAME           CHAOS Experiment TYPE     NEXT SCHEDULE CHAOS INFRASTRUCTURE ID             CHAOS Experiment NAME  LAST UPDATED BY
9433b48c-4ab7-4544-8dab-4a7237619e09 custom-chaos-experiment-1627980541     Non Cron Chaos Experiment None          f9799723-29f1-454c-b830-ae8ba7ee4c30    Self-infra-infra    admin

Showing 1 of 1 Chaos Experiments
```

- To list all the Chaos Experiment runs within a project, issue the following command.

```shell
litmusctl get chaos-experiment-runs  --project-id=""
```

- To list all the Chaos Experiment runs within a specific experiment, issue the following command.

```shell
litmusctl get chaos-experiment-runs  --project-id="" --experiment-id=""
```

- To list the Chaos Experiment run with a specific experiment-run-id , issue the following command.

```shell
litmusctl get chaos-experiment-runs  --project-id="" --experiment-run-id=""
```

**Output:**

```
CHAOS EXPERIMENT RUN ID			STATUS		RESILIENCY SCORE	CHAOS EXPERIMENT ID	    CHAOS EXPERIMENT NAME	TARGET CHAOS INFRA	UPDATED AT			UPDATED BY
8ceb712c-1ed4-40e6-adc4-01f78d281506    Running		0.00             9433b48c-4ab7-4544-8dab-4a7237619e09 custom-chaos-experiment-1627980541 Self-Chaos-Infra   June 1 2022, 10:28:02 pm admin

Showing 1 of 1 Chaos Experiments runs
```

- To describe a particular Chaos Experiment, issue the following command.

Using Flag :

```shell
litmusctl describe chaos-experiment <chaos-experiment-id> --project-id=""
```

Using UI :

```shell
litmusctl describe chaos-experiment
Enter the Project ID: "project-id"
Enter the Chaos Experiment ID: "chaos-experiment-id"
Select an output format :
yaml
json
```

**Output:**

```
apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:
    creationTimestamp: null
    labels:
        cluster_id: f9799723-29f1-454c-b830-ae8ba7ee4c30
        subject: custom-chaos-experiment-1627980541
        workflow_id: 9433b48c-4ab7-4544-8dab-4a7237619e09
        workflows.argoproj.io/controller-instanceid: f9799723-29f1-454c-b830-ae8ba7ee4c30
    name: custom-chaos-experiment-1627980541
    namespace: litmus
spec:
...
```

- To delete a particular Chaos Experiment, issue the following commands.

Using Flag :

```shell
litmusctl delete chaos-experiment <chaos-experiment-id> --project-id=""
```

Using UI :

```shell
litmusctl delete chaos-experiment
Enter the Project ID: "project-id"
Enter the Chaos Experiment ID: "chaos-experiment-id"
Are you sure you want to delete this Chaos Experiment? (y/n): y
```

**Output:**

```
🚀 Chaos Experiment successfully deleted.
```

- To get the Chaos Environment, issue the following command.

Using Flag :

```shell
litmusctl get chaos-environment --project-id="" --environment-id=""
```

Using UI :

```shell
litmusctl get chaos-environment
Enter the Project ID: "project-id"
Enter the Environment ID: "chaos-experiment-id"
```

**Output:**

```
CHAOS ENVIRONMENT ID	 shivamenv
CHAOS ENVIRONMENT NAME	 shivamenv
CHAOS ENVIRONMENT Type	 NON_PROD
CREATED AT		 55908-04-03 16:42:51 +0530 IST
CREATED BY		 admin
UPDATED AT		 55908-04-03 16:42:51 +0530 IST
UPDATED BY		 admin
CHAOS INFRA IDs	 d99c7d14-56ef-4836-8537-423f28ceac4e
```

- To list the Chaos Environments, issue the following command.

Using Flag :

```shell
litmusctl list chaos-environments --project-id=""
```

Using UI :

```shell
litmusctl list chaos-environment
Enter the Project ID: "project-id"
```

**Output:**

```
CHAOS ENVIRONMENT ID	CHAOS ENVIRONMENT NAME		CREATED AT			                  CREATED BY
testenv			        testenv				        55985-01-15 01:42:33 +0530 IST	      admin
shivamnewenv		    shivamnewenv			    55962-10-01 15:05:45 +0530 IST	      admin
newenvironmenttest	    newenvironmenttest		    55912-12-01 10:55:23 +0530 IST	      admin
shivamenv		        shivamenv			        55908-04-03 16:42:51 +0530 IST	      admin

```
---

## Flag details

<table>
    <th>Flag</th>
    <th>Short Flag</th>
    <th>Type</th>
    <th>Description</th>
    <tr>
        <td>--cacert</td>
        <td></td>
        <td>String</td>
        <td>custom ca certificate used by litmusctl for communicating with portal</td>
    </tr>
    <tr>
        <td>--config</td>
        <td></td>
        <td>String</td>
        <td>config file (default is $HOME/.litmusctl)</td>
    </tr>
    <tr>
        <td>--skipSSL</td>
        <td></td>
        <td>Boolean</td>
        <td>litmusctl will skip ssl/tls verification while communicating with portal</td>
    </tr>
    <tr>
        <td>--help</td>
        <td>-h</td>
        <td></td>
        <td>help for litmusctl</td>
    </tr>
</table>

For more information related to flags, Use `litmusctl --help`.
