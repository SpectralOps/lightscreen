![Build Status](https://travis-ci.org/spectralops/lightscreen.svg?branch=master)

![](media/logo.png)

# Lightscreen

A configurable and flexible admission controller built in Go and extensible with Go.

✅ Fast turnaround by using Go to customize your controller, deploying a single binary
✅ Modular and flexible architecture
✅ Simple API
✅ Built-in Server abstraction that takes care of running a first-class Kubernetes controller for you
✅ Great debugging story with a CLI based `check` command

## Quick Start

We're going to use [kind](https://kind.sigs.k8s.io/) to get a feel for lightscreen.

On macOS, install `kind` and `helm` (for other OS see [here](https://kind.sigs.k8s.io/docs/user/quick-start/)):

```
$ brew install kind helm
```

Now set up a cluster, build lightscreen image, preload it into the cluster and set up a kube dashboard:

```
$ git clone https://github.com/spectralops/lightscreen
$ cd lightscreen
$ make kube-start
```

You can now login to your dashboard if you like (not a must). Login token should be at your console.

Let's set up a lightscreen service on our cluster:


```
$ cd deployment/helm && helm install pandas .
```

With your favorite logger, watch the system logs (I use [kail](https://github.com/boz/kail)):

```
$ kail --ns=kube-system
```

Now try to schedule an nginx instance on your cluster, and watch it being blocked:


```
Error creating: admission webhook "lightscreen.spectralops.io" denied the request: Image library/nginx@sha256:70821e443be75ea38bdf52a974fd2271babd5875b2b1964f05025981c75a6717 is not allowed to be admitted
```

If you want to see how nginx is being admitted successfully into your cluster, take a look at [deployment/helm/files/lightscreen.yaml](deployment/helm/files/lightscreen.yaml), edit it to have the latest SHA and repeat the process.

In the general sense, the Spectral Platform uses Lightscreen as its admission controller, and the Spectral Platform generates and maintains this file automatically, based on successfully scanned containers. If you want to build something similar yourself look at [docs/using-configmaps](docs/using-configmap.md).


## Quick Start (Test Cluster)

If you just want to use Lightscreen in your cluster, you can use the Helm chart and load it in (of course, use a test cluster first to understand how admission works)

```
# use this to set up permission or equivalent:
$ kubectl create clusterrolebinding add-on-cluster-admin --clusterrole=cluster-admin --serviceaccount=kube-system:default

# enable lightscreen in your cluster
$ kubectl label namespace default lightscreen=enabled

# just call this deployment 'panda', pandas are cute
$ git clone https://github.com/spectralops/lightscreen
$ cd deployment/helm && helm install panda .
```


To get a bit more value out of Lightscreen you're going to either use it out of the box for admission control, or use it as a library to build your own custom service. Either way it's best if you have a development story ready:

```
$ git clone https://github.com/spectralops/lightscreen
$ cd lightscreen && make build
$ ./lightscreen --help
usage: lightscreen [<flags>]

Flags:
      --help                Show context-sensitive help (also try --help-long and --help-man).
  -c, --config="lightscreen.yaml"
                            Action mapping configuration file
      --check=CHECK         Check input file for admission
      --host="0.0.0.0"      HTTP host
      --port=443            HTTP port
      --metrics=":8080"     Metrics HTTP port
  -p, --production          Run in production mode
      --certs="self-certs"  Certs dir
```

## Lightscreen is all about your actions


To make full use of Lightscreen you want to build your custom actions. You need to implement the following interface:

```Go
type Action interface {
	Name() string
	Run(context.Context, *unstructured.Unstructured) error
}
```

Then, load it into the Lightscreen server:

```Go
server := admission.NewServer(admission.ServerOptions{
    Config:      *config,
    Development: development,
    Address:     *address,
}, logger)
yourAction := NewAction()
server.Actions.Add(yourAction)
```

And serve:

```Go
server.Serve()
```

You can always start off from our [example](examples/spectral-notary).


### Thanks:

To all [Contributors](https://github.com/spectralops/lightscreen/graphs/contributors) - you make this happen, thanks!

# Copyright

Copyright (c) 2020 [@jondot](http://twitter.com/jondot). See [LICENSE](LICENSE.txt) for further details.