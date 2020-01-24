## Quick Start

1. Create a signed cert/key pair and store it in a Kubernetes `secret` that will be consumed by sidecar deployment
```
./deployment/local/webhook-create-signed-cert.sh \
    --service lightscreen-svc \
    --secret lightscreen-certs \
    --namespace default
```

2. Patch the `MutatingWebhookConfiguration` by set `caBundle` with correct value from Kubernetes cluster
```
cat deployment/local/webhook.yaml | \
    deployment/local/webhook-patch-ca-bundle.sh > \
    deployment/local/webhook-ca-bundle.yaml
```

3. Deploy resources

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

4. enable

kubectl label namespace default lightscreen=enabled
