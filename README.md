ecs-formation
==========

[![Circle CI](https://circleci.com/gh/stormcat24/ecs-formation.svg?style=shield&circle-token=baf60b45ce2de8c5d11b3e6d77a3a23ebf2d5991)](https://circleci.com/gh/stormcat24/ecs-formation)
[![License: MIT](http://img.shields.io/badge/license-MIT-orange.svg)](LICENSE)

ecs-formation is a tool for defining several Docker continers and clusters on [Amazon EC2 Container Service(ECS)](https://aws.amazon.com/ecs/).

# Features

* Define services on ECS cluster, and Task Definitions.
* Supports YAML definition like docker-compose. Be able to run ecs-formation if copy docker-compose.yml(formerly fig.yml).
* Manage ECS Services and Task Definitions by AWS API.

# Usage

### Installation

ecs-formation is written by Go. Please run `go get`.

```bash
$ go get github.com/stormcat24/ecs-formation
```

### Define environment variables

ecs-formation requires environment variables to run, as follows.

* AWS_ACCESS_KEY: AWS access key
* AWS_SECRET_ACCESS_KEY: AWS secret access key
* AWS_REGION: Target AWS region name

### Make working directory

Make working directory for ecs-formation. This working directory should be managed by Git.

```bash
$ mkdir -p path-to-path/test-ecs-formation
$ mkdir -p path-to-path/test-ecs-formation/task
$ mkdir -p path-to-path/test-ecs-formation/service
$ mkdir -p path-to-path/test-ecs-formation/bluegreen
```

### Make ECS Cluster

You need to create ECS cluster in advance. And also, ECS instance must be join in ECS cluster.

### Define Task Definitions

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

### Define Services on Cluster

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

### Manage Task Definitions

Show update plan.

```bash
(path-to-path/test-ecs-formation $ ecs-formation task plan
```

Apply definition.

```bash
(path-to-path/test-ecs-formation $ ecs-formation task apply
```

### Manage Services on Cluster

Show update plan.

```bash
(path-to-path/test-ecs-formation $ ecs-formation service plan
```

Apply definition.

```bash
(path-to-path/test-ecs-formation $ ecs-formation service apply
```


License
===
See [LICENSE](LICENSE).

Copyright © Akinori Yamada([@stormcat24](https://twitter.com/stormcat24)). All Rights Reserved.
