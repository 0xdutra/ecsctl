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
completion      Generate the autocompletion script for the specified shell
deploy          Commands to deploy ECS infrastructure
elb             Commands to manage Elastic Load Balancer
help            Help about any command
services        Commands to manage ECS services
target-group    Commands to manage target groups
task-definition Commands to manage ECS task definitions
```

<hr>


## Examples

- [Cluster](#Cluster)
  - [Create](#Create)
  - [Delete](#Delete)
  - [List](#List)

- [Load balancer](#Load-balancer)
	- [List](#List)
  - [Describe](#Describe)
  - [Delete](#Delete)

- [Services](#Services)
  - [Events](#Service-events)
	- [Status](#Service-status)
	- [Update capacity](#Updating-capacity)
	- [Update task definition](#Updating-task-definition)
	- [Example - Deploying HTTP service on ECS using AWS Fargate](#Deploying-HTTP-service-on-Amazon-ECS-using-AWS-Fargate)

- [Target group](#Target-group)
	- [List](#list)
	- [Describe](#Describe)
	- [Delete](#Delete)

- [Task definition](#Task-definition)
	- [Edit](#Edit)

<hr>

## Cluster

### Create

```sh
ecsctl cluster create --name <cluster name>
```

### Delete

```sh
ecsctl cluster delete --name <cluster name>
```

### List

```sh
ecsctl cluster list
```

<hr>

## Load balancer

### List

```sh
ecsctl elb list
```

### Describe

```sh
ecsctl elb describe --name <elb name>
```

### Delete

```sh
ecsctl elb delete --arn <elb arn>
```


<hr>

## Services

### Servoce events

```sh
./ecsctl services events --cluster <cluster name> --service <service name>
```

### Service status

```sh
ecsctl services status --service <service name> --cluster <cluster name>
```

### Updating capacity

```sh
ecsctl services update-capacity --service <service name> --cluster <cluster name> --min 10 --max 20 --desired 10
```

### Updating task definition

*Tip*: Updating the service to a specific task definition revision, use `taskname:revision`, see a example:

```sh
ecsctl services update --service <service name> --cluster <cluster name> --task-definition <example:10> --force-new-deployment
```

To update to the latest task definition, use only the task definition name.

```sh
ecsctl services update --service <service name> --task-definition <example> --cluster <cluster name> --force-new-deployment
```

### Deploying HTTP service on Amazon ECS using AWS Fargate

### Target group

```yaml
# target_group.yaml

targetGroup:
  name: httpd
  port: 80
  protocol: HTTP
  vpcId: "" # your VPC id
  targetType: ip
  healthCheckEnabled: true
  healthCheckIntervalSeconds: 30
  healthCheckPath: /
```

```sh
ecsctl deploy -c target_group.yaml -r target-group
```

### Load Balancer

*Tip: you can get TargetGroupArn using*:

```sh
ecsctl target-group describe --name httpd
```

```yaml
# load_balancer.yaml

loadBalancers:
  name: httpd
  subnets: [] # The subnets must be in the VPC used in the target group
  type: application
  scheme: internal

  listener:
    defaultActions:
      targetGroupArn: "" # The ARN of the target group
      type: forward
    Port: 80
    protocol: HTTP
```

```sh
ecsctl deploy -c load_balancer.yaml -r load-balancer
```

### Task definition

```yaml
# task_definition.yaml

taskDefinition:
  name: httpd
  executionRoleArn: ecsTaskExecutionRole
  containerDefinitions:
    name: httpd
    image: httpd
    cpu: 0
    memory: 512
    memoryReservation: 512
    portsMappings:
      hostPort: 80
      protocol: TCP
      containerPort: 80
    logConfiguration:
      logDriver: awslogs
      options:
        awslogs-group: ecs/httpd-log-group # CloudWatch log group is manually created
        awslogs-region: us-east-1
        awslogs-stream-prefix: httpd
  memory: 512
  taskRoleArn: ecsTaskExecutionRole
  requiresCompatibilities:
    - FARGATE
  family: httpd
  cpu: 256
  networkMode: "awsvpc"
```

```sh
ecsctl deploy -c task_definition.yaml -r task-definition
```

### Service

*Tip: you can get TargetGroupArn using*:

```sh
ecsctl target-group describe --name httpd
```

```yaml
# service.yaml

service:
  serviceName: httpd
  cluster: default
  taskDefinition: httpd
  desiredCount: 2
  enableECSmanagedTags: false
  enableExecuteCommand: false
  healthCheckGracePeriodSeconds: 0
  launchType: FARGATE

  deploymentConfiguration:
    maximumPercent: 200
    minimumHealthyPercent: 100

  deploymentController:
    type: ECS

  loadBalancer:
    containerName: httpd
    containerPort: 80
    targetGroupArn: "" # The ARN of the target group

  awsVpcConfiguration:
    assignPublicIp: ENABLED
    securityGroups: [] # The security group is optional, if is empty, the default security group is used
    subnets: [] # The subnets must be in the VPC used in the target group

  schedulingStrategy: REPLICA
```

```sh
ecsctl deploy -c service.yaml -r service
```

### Getting service status

```sh
ecsctl services status --cluster default --service httpd
```

<hr>

## Target group

### List

```sh
ecsctl target-group list
```

#### Describe

```sh
ecsctl target-group describe --name <target group>
```

### Delete

```sh
ecsctl target-group delete --arn <target group>
```

<hr>

## Task definition

### Edit

```sh
ecsctl task-definition edit --name <task definition name>
```
