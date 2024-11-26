# wonjinsin/go-web-boilerplate(magmar)

Simple rest api server with [Echo framework](https://github.com/labstack/echo)

[![License MIT](https://img.shields.io/badge/License-MIT-green.svg)](http://opensource.org/licenses/MIT)
[![DB Version](https://img.shields.io/badge/DB-Redis-red)](https://redis.io/)
[![DB Version](https://img.shields.io/badge/DB-Mysql-blue)](https://www.mysql.com/)
[![Go Report Card](https://goreportcard.com/badge/github.com/StarpTech/go-web)](https://goreportcard.com/report/github.com/wonjinsin/go-web-boilerplate)

## Features

- [Gorm](https://github.com/go-gorm/gorm) : For mysql ORM
- [Go-redis](https://github.com/go-redis/redis/v8) : Go redis client
- [Zap](https://github.com/uber-go/zap) : Go leveled Logging
- [Viper](https://github.com/spf13/viper) : Config Setting
- [Ginko](https://onsi.github.io/ginkgo) : BDD Testing Framework with [Gomega](https://onsi.github.io/gomega)
- [Makefile]() : go build, test, vendor using Make

## Project structure

DDD(Domain Driven Design) pattern with controller, model, servie, repository

## Getting started

### Initial action

```
$ make all && make build && make start
```

### Build vendors

```
$ make vendor
```

### Build and start

```
$ make build && bin/magmar
```

### Test

```
$ make vet && make lint && make test
// or
$ make test-all
```

### Clean

```
$ make clean
```
