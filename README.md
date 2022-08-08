<p align="center">
  <h3 align="center">ecsctl</h3>
  <p align="center">ecsctl is a tool for management Elastic Container Service.</p>

  <p align="center">
    <a href="https://twitter.com/0xdutra">
      <img src="https://img.shields.io/badge/twitter-@0xdutra-blue.svg">
    </a>
    <a href="https://opensource.org/licenses/BSD-2-Clause">
      <img src="https://img.shields.io/badge/License-BSD-green.svg">
    </a>
    <a href="https://github.com/0xdutra/ecsctl/actions/workflows/golangci-lint.yml">
        <img src="https://github.com/0xdutra/ecsctl/actions/workflows/golangci-lint.yml/badge.svg">
    </a>
  </p>
</p>

<hr>

### All commands available

```sh
cluster         - Commands to manage ECS clusters
task-definition - Commands to manage task definitions
services        - Commands to manage ECS services
```

<hr>

### Cluster subcommands

```sh
create-cluster - Create ECS cluster
delete-cluster - Delete ECS cluster
list-clusters  - List all ECS clusters
```

### Cluster examples

```sh
ecsctl cluster create-cluster --name example01
ecsctl cluster delete-cluster --name example01
```

<hr>

### Task definition subcommands

```sh
register-task-definition - Register task definition
edit-task-definition     - Edit a task definition using a text editor
```

### Task definition examples

```sh
ecsctl task-definition register-task-definition --input-json examples/task_definition_example.json
```

```sh
ecsctl task-definition edit-task-definition --name ecsctl-apache-example --revision 1 --editor nano
```

<hr>

### Service subcommands

```sh
create-service     - Commands to create ECS services
describe-services  - Commands to describe ECS services
list-services      - Commands to list services in your ECS cluster
update-capacity    - Commands to update services capacity
```

### Service examples

Updating capacity

```sh
ecsctl services update-capacity --min 2 --max 2 --desired 2 --service-name <service name> --cluster-name <cluster name>
```

Creating a service using a JSON manifest

```sh
ecsctl services create-service --input-json examples/service_example.json
```

Listing services

```sh
ecsctl services list-services --cluster-name <cluster name>
```

Describe service informations

```sh
ecsctl services describe-services --service-arn <service arn> --cluster-name <cluster name>
```
