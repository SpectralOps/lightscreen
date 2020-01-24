---
id: usage
title: Usage
sidebar_label: Usage
---


Lightscreen is a configurable and extensible admission controller for Kubernetes. It is built in Go, and takes into account that extensions are also built in Go with no scripting language or abstraction between you and your extension; you can use the full power of Go.

You can also use Lightscreen's codebase to build your own custom admission controllers (for example if you want to add a scripting language for rules) or mount it in your existing controller or service.


## Quick Start

Run lightscreen with the container on your local development environment for experimenting:

```
$ lightscreen --help
usage: lightscreen [<flags>]

Flags:
      --help                Show context-sensitive help (also try --help-long and --help-man).
  -c, --config="lightscreen.yaml"
                            Action mapping configuration file
      --host="0.0.0.0"      HTTP host
      --port=443            HTTP port
      --metrics=":8080"     Metrics HTTP port
      --check="input.json"
  -p, --production          Run in production mode
      --certs="self-certs"  Certs dir
```
