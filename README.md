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

### Example

Creating a fargate httpd service

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

### Getting service status

```sh
ecsctl services status --cluster default --service httpd
```
