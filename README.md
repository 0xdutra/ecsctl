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
cluster         Commands to manage ECS cluster
elb             Commands to manage Elastic Load Balancer
services        Commands to manage ECS services
targetgroup     Commands to manage target groups
task-definition Commands to manage ECS task definitions
```

<hr>

### Cluster subcommands

```sh
create - Create ECS cluster
delete - Delete ECS cluster
list   - List all ECS clusters
```

### Cluster examples

```sh
ecsctl cluster create --cluster example01
ecsctl cluster list
ecsctl cluster delete --cluster example01
```

<hr>

### Task definition subcommands

```sh
register - Register task definition
edit     - Edit a task definition using a text editor
```

### Task definition examples

```sh
ecsctl task-definition register --input-json examples/task_definition_example.json
```

```sh
ecsctl task-definition edit --name ecsctl-apache-example --revision 1 --editor nano
```

<hr>

### Service subcommands

```sh
create          - Commands to create ECS service
describe        - Commands to describe ECS service
list            - Commands to list service in your ECS cluster
status          - Commands to check the status of service
update          - Commands to update ECS service
update-capacity - Commands to update service capacity
```

### Service examples

Checking status

```
ecsctl services status --service <service name> --cluster <cluster name>
```

Updating service

```sh
ecsctl services update --service <service name> --task-definition <taskdef name> --cluster <cluster name>
```

Updating capacity

```sh
ecsctl services update-capacity --min 2 --max 2 --desired 2 --service <service name> --cluster <cluster name>
```

Creating a service using a JSON manifest

```sh
ecsctl services create --input-json examples/service_example.json
```

Listing services

```sh
ecsctl services list --cluster <cluster name>
```

Describe service informations

```sh
ecsctl services describe --service <service name> --cluster <cluster name>
```
