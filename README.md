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
```

<hr>

### cluster subcommands

```sh
create-cluster - Create ECS cluster
delete-cluster - Delete ECS cluster
list-clusters  - List all ECS clusters
```

### cluster examples

```sh
ecsctl cluster create-cluster --name example01
ecsctl cluster delete-cluster --name example01
```

<hr>

### task definition subcommands

```sh
create-task-definition - Create task definition
```

### task definition examples

```sh
ecsctl task-definition create-task-definition --input-json examples/task_definition_example.json
```

