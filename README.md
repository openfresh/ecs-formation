ecs-formation
==========

[![Circle CI](https://circleci.com/gh/stormcat24/ecs-formation.svg?style=shield&circle-token=baf60b45ce2de8c5d11b3e6d77a3a23ebf2d5991)](https://circleci.com/gh/stormcat24/ecs-formation)
[![Language](http://img.shields.io/badge/language-go-brightgreen.svg?style=flat)](https://golang.org/)
[![issues](https://img.shields.io/github/issues/stormcat24/ecs-formation.svg?style=flat)](https://github.com/stormcat24/ecs-formation/issues?state=open)
[![License: MIT](http://img.shields.io/badge/license-MIT-orange.svg)](LICENSE)

ecs-formation is a tool for defining several Docker continers and clusters on [Amazon EC2 Container Service(ECS)](https://aws.amazon.com/ecs/).

# Features

* Define services on ECS cluster, and Task Definitions.
* Supports YAML definition like docker-compose. Be able to run ecs-formation if copy docker-compose.yml(formerly fig.yml).
* Manage ECS Services and Task Definitions by AWS API.

# Usage

### Setup

#### Installation

ecs-formation is written by Go. Please run `go get`.

```bash
$ go get github.com/stormcat24/ecs-formation
```

#### Define environment variables

ecs-formation requires environment variables to run, as follows.

* AWS_ACCESS_KEY: AWS access key
* AWS_SECRET_ACCESS_KEY: AWS secret access key
* AWS_REGION: Target AWS region name

#### Make working directory

Make working directory for ecs-formation. This working directory should be managed by Git.

```bash
$ mkdir -p path-to-path/test-ecs-formation
$ mkdir -p path-to-path/test-ecs-formation/task
$ mkdir -p path-to-path/test-ecs-formation/service
$ mkdir -p path-to-path/test-ecs-formation/bluegreen
```

### Manage Task Definition and Services

#### Make ECS Cluster

You need to create ECS cluster in advance. And also, ECS instance must be join in ECS cluster.

#### Define Task Definitions

Make Task Definitions file in task directory. This file name is used as ECS Task Definition name.

```bash
(path-to-path/test-ecs-formation/task) $ vim test-definition.yml
nginx:
  image: nginx:latest
  ports:
    - 80:80
  environment:
    PARAM1: value1
    PARAM2: value2
  links:
    - api
  memory: 512
  cpu_units: 512
  essential: true

api:
  image: your_namespace/your-api:latest
  ports:
    - 8080:8080
  memory: 1024
  cpu_units: 1024
  essential: true
  links:
    - redis

redis:
  image: redis:latest
  ports:
    - 6379:6379
  memory: 512
  cpu_units: 512
  essential: true
```

#### Define Services on Cluster

Make Service Definition file in cluster directory. This file name must be equal ECS cluster name.

For example, if target cluster name is `test-cluster`, you need to make `test-cluster.yml`.

```bash
(path-to-path/test-ecs-formation/service) $ vim test-cluster.yml
test-service:
  task_definition: test-definition
  desired_count: 1
  role: your-ecs-elb-role
  load_balancers:
    -
      name: test-elb
      container_name: nginx
      container_port: 80
```

#### Keep desired_count at updating service

If you modify value of `desired_count` by AWS Management Console or aws-cli, you'll fear override value of `desired_count` by ecs-formation. This value should be flexibly changed in the operation.

If `keep_desired_count` is `true`, keep current `desired_count` at updating service.

```bash
(path-to-path/test-ecs-formation/service) $ vim test-cluster.yml
test-service:
  task_definition: test-definition
  desired_count: 1
  keep_desired_count: true
```


#### Manage Task Definitions

Show update plan.

```bash
(path-to-path/test-ecs-formation $ ecs-formation task plan
```

Apply definition.

```bash
(path-to-path/test-ecs-formation $ ecs-formation task apply
```

#### Manage Services on Cluster

Show update plan.

```bash
(path-to-path/test-ecs-formation $ ecs-formation service plan
```

Apply definition.

```bash
(path-to-path/test-ecs-formation $ ecs-formation service apply
```

### Blue Green Deployment

ecs-formation supports blue-green deployment.

#### Requirements on ecs-formation

* Requires two ECS cluster. Blue and Green.
* Requires two ELB. Primary ELB and Standby ELB.
* ECS cluster should be built by EC2 Autoscaling group.

#### Define Blue Green Deployment

Make management file of Blue Green Deployment file in bluegreen directory.

```bash
(path-to-path/test-ecs-formation/bluegreen) $ vim test-bluegreen.yml
blue:
  cluster: test-blue
  service: test-service
  autoscaling_group: test-blue-asg
green:
  cluster: test-green
  service: test-service
  autoscaling_group: test-green-asg
primary_elb: test-elb-primary
standby_elb: test-elb-standby
```

Show blue green deployment plan.

```bash
(path-to-path/test-ecs-formation $ ecs-formation bluegreen plan
```

Apply blue green deployment.

```bash
(path-to-path/test-ecs-formation $ ecs-formation bluegreen apply
```

if with `--nodeploy` option, not update services. Only swap ELB on blue and green groups.

```bash
(path-to-path/test-ecs-formation $ ecs-formation bluegreen apply --nodeploy
```

If autoscaling group have several different ELB, you should specify array property of `chain_elb`. ecs-formation can swap `chain_elb` ELB group with main ELB group at the same time.

```Ruby
(path-to-path/test-ecs-formation/bluegreen) $ vim test-bluegreen.yml
blue:
  cluster: test-blue
  service: test-service
  autoscaling_group: test-blue-asg
green:
  cluster: test-green
  service: test-service
  autoscaling_group: test-green-asg
primary_elb: test-elb-primary
standby_elb: test-elb-standby
chain_elb:
  - primary_elb: test-internal-elb-primary
    standby_elb: test-internal-elb-standby
```

### Others
#### Passing custom parameters

You can use custom parameters. Define parameters in yaml file(task, service, bluegreen) as follows.

```Ruby
nginx:
    image: stormcat24/nginx:${NGINX_VERSION}
    ports:
        - 80:${NGINX_PORT}
```

You can set value for these parameters by using `-p` option.

```bash
ecs-formation task -p NGINX_VERSION=1.0 -p NGINX_PORT=80 plan your-web-task
```

#### env_file

You can use `env_file` like docker-compose. https://docs.docker.com/compose/compose-file/#env-file

```Ruby
nginx:
    image: stormcat24/nginx:${NGINX_VERSION}
    ports:
        - 80:${NGINX_PORT}
    env_file:
        - ./test1.env
        - ../test2.env
```

License
===
See [LICENSE](LICENSE).

Copyright © Akinori Yamada([@stormcat24](https://twitter.com/stormcat24)). All Rights Reserved.
