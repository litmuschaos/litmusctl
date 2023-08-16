### litmusctl Syntax
`litmusctl` has a syntax to use as follows:

```shell
litmusctl [command] [TYPE] [flags]
```
* Command: refers to what you do want to perform (connect, create, get and config)
* Type: refers to the feature type you are performing a command against (chaos-delegate, project etc.)
* Flags: It takes some additional information for resource operations. For example, `--installation-mode` allows you to specify an installation mode.

Litmusctl is using the `.litmusconfig` config file to manage multiple accounts
1. If the --config flag is set, then only the given file is loaded. The flag may only be set once and no merging takes place.
2. Otherwise, the ${HOME}/.litmusconfig file is used, and no merging takes place.

Litmusctl supports both interactive and non-interactive(flag based) modes.
> Only `litmusctl connect chaos-delegate`  command needs --non-interactive flag, other commands don't need this flag to be in non-interactive mode. If mandatory flags aren't passed, then litmusctl takes input in an interactive mode.

### Installation modes
Litmusctl can install a Chaos Delegate in two different modes.
* cluster mode: With this mode, the Chaos Delegate can run the chaos in any namespace. It installs appropriate cluster roles and cluster role bindings to achieve this mode. It can be enabled by passing a flag `--installation-mode=cluster`

* namespace mode: With this mode, the Chaos Delegate can run the chaos in its namespace. It installs appropriate roles and role bindings to achieve this mode. It can be enabled by passing a flag `--installation-mode=namespace`

Note: With namespace mode, the user needs to create the namespace to install the Chaos Delegate as a prerequisite.

### Minimal steps to connect a Chaos Delegate

* To setup an account with litmusctl
```shell
litmusctl config set-account --endpoint="" --username="" --password=""
```

* To create an Chaos Delegate with an existing project
> Note: To get `project-id`. Apply `litmusctl get projects`

```shell
litmusctl connect chaos-delegate --name="" --project-id="" --non-interactive
```

### Flags for `connect chaos-delegate` command
<table>
    <tr>
        <th>Flag</th>
        <th>Short Flag</th>
        <th>Type</th>
        <th>Description</th>
    <tr>
        <td>--description</td>
        <td></td>
        <td>String</td>
        <td>Set the Chaos Delegate description (default "---")</td>
    </tr>
    <tr>
        <td>--name</td>
        <td></td>
        <td>String</td>
        <td>Set the name of Chaos Delegate which should be unique</td>
    </tr>
    <tr>
        <td>--skip-ssl</td>
        <td></td>
        <td>Boolean</td>
        <td>Set whether Chaos Delegate will skip ssl/tls check (can be used for self-signed certs, if cert is not provided in portal) (default false)</td>
    </tr>
    <tr>
        <td>--chaos-delegate-type</td>
        <td></td>
        <td>String</td>
        <td>Set the chaos-delegate-type to external for external Chaos Delegates | Supported=external/internal (default "external")</td>
    </tr>
    <tr>
        <td>--installation-mode</td>
        <td></td>
        <td>String</td>
        <td>Set the installation mode for the kind of Chaos Delegate | Supported=cluster/namespace (default "cluster")</td>
    </tr>
    <tr>
        <td>--kubeconfig</td>
        <td>-k</td>
        <td>String</td>
        <td>Set to pass kubeconfig file if it is not in the default location ($HOME/.kube/config)</td>
    </tr>
    <tr>
        <td>--namespace</td>
        <td></td>
        <td>String</td>
        <td>Set the namespace for the Chaos Delegate installation (default "litmus")</td>
    </tr>
    <tr>
        <td>--node-selector</td>
        <td></td>
        <td>String</td>
        <td>Set the node-selector for Chaos Delegate components | Format: key1=value1,key2=value2)</td>
    </tr>
    <tr>
        <td>--non-interactive</td>
        <td>-n</td>
        <td>String</td>
        <td>Set it to true for non interactive mode | Note: Always set the boolean flag as --non-interactive=Boolean</td>
    </tr>
    <tr>
        <td>--ns-exists</td>
        <td></td>
        <td>Boolean</td>
        <td>Set the --ns-exists=false if the namespace mentioned in the --namespace flag is not existed else set it to --ns-exists=true | Note: Always set the boolean flag as --ns-exists=Boolean</td>
    </tr>
    <tr>
        <td>--platform-name</td>
        <td></td>
        <td>String</td>
        <td>Set the platform name. Supported- AWS/GKE/Openshift/Rancher/Others (default "Others")</td>
    </tr>
    <tr>
        <td>--sa-exists</td>
        <td></td>
        <td>Boolean</td>
        <td>Set the --sa-exists=false if the service-account mentioned in the --service-account flag is not existed else set it to --sa-exists=true | Note: Always set the boolean flag as --sa-exists=Boolean"</td>
    </tr>
    <tr>
        <td>--service-account</td>
        <td></td>
        <td>String</td>
        <td>Set the service account to be used by the Chaos Delegate (default "litmus")</td>
    </tr>
    <tr>
        <td>--config</td>
        <td></td>
        <td>String</td>
        <td>config file (default is $HOME/.litmusctl)</td>
    </tr>
</table>

---

### Steps to create a chaos scenaro

* To setup an account with litmusctl
```shell
litmusctl config set-account --endpoint="" --username="" --password=""
```

* To create a Chaos Scenario by passing a manifest file
> Note:
> * To get `project-id`, apply `litmusctl get projects`
> * To get `chaos-delegate-id`, apply `litmusctl get chaos-delegates --project-id=""`
```shell
litmusctl create chaos-scenario -f custom-chaos-scenario.yml --project-id="" --chaos-delegate-id=""
```

---

### Additional commands

* To view the current configuration of `.litmusconfig`, type:
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
apiVersion: v1
current-account: https://preview.litmuschaos.io
current-user: litmus-user
kind: Config
```

* To get an overview of the accounts available within `.litmusconfig`, use the `config get-accounts` command:

```shell
litmusctl config get-accounts
```

**Output:**

```
CURRENT  ENDPOINT                         USERNAME  EXPIRESIN
         https://preview.litmuschaos.io   admin     2021-07-22 01:20:27 +0530 IST
*        https://preview.litmuschaos.io   raj       2021-07-22 14:33:22 +0530 IST
```

* To alter the current account use the `use-account` command with the --endpoint and --username flags:
```shell
litmusctl config use-account --endpoint="" --username=""
```

* To create a project, apply the following command with the `--name` flag:
```shell
litmusctl create project --name=""
```

* To view all the projects with the user, use the `get projects` command.
```shell
litmusctl get projects
```

**Output:**

```
PROJECT ID                                PROJECT NAME       CREATEDAT
50addd40-8767-448c-a91a-5071543a2d8e      Developer Project  2021-07-21 14:38:51 +0530 IST
7a4a259a-1ae5-4204-ae83-89a8838eaec3      DevOps Project     2021-07-21 14:39:14 +0530 IST
```


* To get an overview of the Chaos Delegates available within a project, issue the following command.
```shell
litmusctl get chaos-delegates --project-id=""
```

**Output:**

```
CHAOS DELEGATE ID                            CHAOS DELEGATE NAME          STATUS     REGISTRATION
55ecc7f2-2754-43aa-8e12-6903e4c6183a   chaos-delegate-1            ACTIVE     REGISTERED
13dsf3d1-5324-54af-4g23-5331g5v2364f   chaos-delegate-2            INACTIVE   NOT REGISTERED
```


* To disconnect a Chaos Delegate, issue the following command..
```shell
litmusctl disconnect chaos-delegate <chaos-delegate-id> --project-id=""
```

**Output:**

```
🚀 Chaos Delegate successfully disconnected.
```


* To list the created Chaos Scenarios within a project, issue the following command.
```shell
litmusctl get chaos-scenarios --project-id=""
```

**Output:**

```
CHAOS SCENARIO ID                         CHAOS SCENARIO NAME                   CHAOS SCENARIO TYPE     NEXT SCHEDULE CHAOS DELEGATE ID                             CHAOS DELEGATE NAME LAST UPDATED BY
9433b48c-4ab7-4544-8dab-4a7237619e09 custom-chaos-scenario-1627980541 Non Cron Chaos Scenario None          f9799723-29f1-454c-b830-ae8ba7ee4c30 Self-Chaos-delegate admin

Showing 1 of 1 Chaos Scenarios
```


* To list all the Chaos Scenario runs within a project, issue the following command.
```shell
litmusctl get chaos-scenario-runs  --project-id=""
```

**Output:**

```
CHAOS SCENARIO RUN ID                      STATUS  RESILIENCY SCORE CHAOS SCENARIO ID                          CHAOS SCENARIO NAME                    TARGET CHAOS DELEGATE LAST RUN                 EXECUTED BY
8ceb712c-1ed4-40e6-adc4-01f78d281506 Running 0.00             9433b48c-4ab7-4544-8dab-4a7237619e09 custom-chaos-scenario-1627980541 Self-Chaos-Delegate   June 1 2022, 10:28:02 pm admin

Showing 1 of 1 Chaos Scenario runs
```


* To describe a particular Chaos Scenario, issue the following command.
```shell
litmusctl describe chaos-scenario <chaos-scenario-id> --project-id=""
```

**Output:**

```
apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:
    creationTimestamp: null
    labels:
        cluster_id: f9799723-29f1-454c-b830-ae8ba7ee4c30
        subject: custom-chaos-scenario_litmus
        workflow_id: 9433b48c-4ab7-4544-8dab-4a7237619e09
        workflows.argoproj.io/controller-instanceid: f9799723-29f1-454c-b830-ae8ba7ee4c30
    name: custom-chaos-scenario-1627980541
    namespace: litmus
spec:
...
```


* To delete a particular Chaos Scenario, issue the following command.
```shell
litmusctl delete chaos-scenario <chaos-scenario-id> --project-id=""
```

**Output:**

```
🚀 Chaos Scenario successfully deleted.
```



For more information related to flags, Use `litmusctl --help`.

----
