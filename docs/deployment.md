---
id: deployment
title: Deployment
sidebar_label: Deployment
---

# Set up with Kind

The goal of the development harness using Kind as the kube cluster is to
host Lightscreen as an admission controller and allow you to manually test with some kind of image to run on the cluster, for example `nginx`.

First, populate every new shell with `kind`'s environment for wiring your kubectl command to it:

```
export KUBECONFIG="$(kind get kubeconfig-path --name="kind")"
```

Then to test your admission controller (see how to deploy it below):

```
# note: uses local 'k' aliases. replace with kubectl if needed.

$ kd deployment nginx && k run nginx --image=nginx --replicas=1
```

View logs with [`kail`](https://github.com/boz/kail) to confirm that the `nginx` image has been rejected or accepted.

```
$ kail -n kube-system
```

# Deploying Lightscreen From Scratch

First we prepare a `kind` based kube cluster which includes:

* Kubernetes
* Kube admin/dash
* Provisioned dashboard admin token
* Build a Lightscreen Docker image and preload it to `kind` (so that you won't need network or registry back/forth)
* Start a kube proxy

```
$ make kube-start
```

Now you should have a kube system ready. Deploying lightscreen can happen with either of the following strategies for your convenience:

## 1. Using Helm

If you don't have tiller set up:

```
$ helm init
```

Then a one-time set up for permission is required (you may skip this and use your own authentication best practices):

```
$ kubectl create clusterrolebinding add-on-cluster-admin --clusterrole=cluster-admin --serviceaccount=kube-system:default
```

Lets tell this cluster that we want Lightscreen to be enabled:

```
$ kubectl label namespace default lightscreen=enabled
```

And deploy:

```
$ cd deployment/helm && helm install --name panda
```

## 2. Manual

Create a signed cert/key pair and store it in a Kubernetes `secret` that will be consumed by sidecar deployment

```
./deployment/local/webhook-create-signed-cert.sh \
    --service lightscreen-svc \
    --secret lightscreen-certs \
    --namespace default
```

Patch the `MutatingWebhookConfiguration` by set `caBundle` with correct value from Kubernetes cluster

```
cat deployment/local/webhook.yaml | \
    deployment/local/webhook-patch-ca-bundle.sh > \
    deployment/local/webhook-ca-bundle.yaml
```

Deploy resources

```
kubectl create -f deployment/local/deployment.yaml
kubectl create -f deployment/local/service.yaml
kubectl create -f deployment/local/webhook-ca-bundle.yaml
```

or

```
kubectl apply -f deployment/local/deployment.yaml
kubectl apply -f deployment/local/service.yaml
kubectl apply -f deployment/local/webhook-ca-bundle.yaml
```

Enable

```
kubectl label namespace default lightscreen=enabled
```

## 3. As a Remote Webhook

Set up Lightscreen on an SSL connection and write the address in `deployments/webhook-remote.yaml`. Then apply:

```
kubectl apply -f deployments/webhook-remote.yaml
```

If you want to use this set up to iterate quickly, serving Lightscreen on your local workstation you can use a service like ngrok to tunnel easily.
