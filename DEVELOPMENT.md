# Litmusctl Local Development Setup Guide

## Introduction

Welcome to the local development setup guide for **`litmusctl`**. This guide will walk you through the steps required to set up and run **`litmusctl`** on your local machine.

## Important Note

Before running **`litmusctl`**, make sure you have a Chaos Centre running. Ensure that the Chaos Centre version is compatible with the **`litmusctl`** version you are using.

## Prerequisites

Before you begin, ensure that you have the following prerequisites installed on your machine:

- [Go programming language](https://golang.org/doc/install) (version or later)
- [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- Kubeconfig - `litmusctl` needs the kubeconfig of the k8s cluster where we need to connect litmus Chaos Delegates. The CLI currently uses the default path of kubeconfig i.e. `~/.kube/config`.

## Clone the Repository

```bash
git clone https://github.com/litmuschaos/litmusctl.git

cd litmusctl
```

## **Install Dependencies**

```bash
go mod download
```

## **Configuration**

Before running **`litmusctl`**, update the following configuration paths in the **`pkg/utils/constants.go`**

From this

```go
// Graphql server API path
GQLAPIPath = "/api/query"

// Auth server API path
AuthAPIPath = "/auth"
```

To this

```go
// Graphql server API path
GQLAPIPath = "/query"

// Auth server API path
AuthAPIPath = ""
```

## **Running `litmusctl`**

Execute the following command to run **`litmusctl`** locally:

```bash
go run main.go <command> <subcommand> <subcommand> [options and parameters]
```

## **Testing `litmusctl`**

To run tests, use the following command:

```bash
go test ./...
```

## **Contributing Guidelines**

If you wish to contribute to **`litmusctl`**, please follow our [contributing guidelines](https://github.com/litmuschaos/litmus/blob/master/CONTRIBUTING.md). Your contributions are valuable, and adhering to these guidelines ensures a smooth and collaborative development process.

## **Troubleshooting**

If you encounter any issues during setup, refer to our [troubleshooting guide](https://docs.litmuschaos.io/docs/troubleshooting) or reach out to our community for assistance. We're here to help you overcome any obstacles and ensure a successful setup.

## **Additional Information**

For more details on using **`litmusctl`**, refer to our [official documentation](https://docs.litmuschaos.io/). This documentation provides comprehensive information to help you make the most out of **`litmusctl`**.

Thank you for setting up **`litmusctl`** locally! Feel free to explore and contribute to the project. Your involvement is crucial to the success of the **`litmusctl`** community.

Let the chaos begin! 🚀🔥
