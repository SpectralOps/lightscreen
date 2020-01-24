---
id: admission
title: Admission
sidebar_label: Admission
---


Lightscreen comes with a few built-in actions (an action is a pluggable behavior that you can add to Lightscreen).  These directly relate to Lightscreen's purpose serving as a screening admission controller -- deciding which image moves into your cluster and which images stays out.

You can configure Lightscreen's behavior with an `lightscreen.yaml`.

## Image SHA Resolver

In order to reason about an image in terms of admission, we need to precisely understand what contents this image has. A way to uniquely identify image content is already built into the container model -- the image SHA signature. However, in Kubernetes, a Pod admission request only contains an image name (it's logical name, for example: nginx).

And so, Lightscreen contains an image resolving action that you can configure. This action is considered a mutating action.

By default the resolver will take an image name (e.g. 'nginx') and pull its SHA from Docker's public registry. The resolving provider uses [`crane`](https://github.com/google/go-containerregistry/tree/master/cmd/crane), a modular registry access library.

Here's how to configure `crane` in your `lightscreen.yaml`:

```yaml
mutations:
  - type: resolve_sha
    finder: crane
```

If you don't want to go out to the internet to resolve an image name to a SHA, and/or you already have a job that builds this kind of `image:SHA` mapping, you can include the mapping directly and use the `dict` finder:


```yaml
mutations:
  - type: resolve_sha
    finder: dict
    finder_config:
        nginx: 3fb
```

## Image Admission

Another action that's part of Lightscreen screening functionality is _image admission_. This action is considered a validation action.

It simply takes an image name, and either admits or rejects it based on a whitelist. Here's how to configure it:

```yaml
validations:
  - type: admit_sha
    admit:
      "library/nginx@sha256:a8517b1d89209c88eeb48709bc06d706c261062813720a352a8e4f8d96635d9d": true
```

# Screening with Lightscreen

Typically, and with the two actions specified above, you want to use them in combination:

```yaml
mutations:
  - type: resolve_sha
    finder: crane
    finder_config:
        foo: bar
validations:
  - type: admit_sha
    admit:
      "library/nginx@sha256:a8517b1d89209c88eeb48709bc06d706c261062813720a352a8e4f8d96635d9d": true
```

This goes in your `lightscreen.yaml`. In production, your `lightscreen.yaml` file should be generated from a configmap (see our section about [using a ConfigMap](using-configmap.md) to dynamically configure Lightscreen).

# Extending Lightscreen's Admission Model

Have more ideas regarding admitting an image into your cluster? All you need to do is break this idea into `Action`s. An action that changes a request is a _mutating action_ and an action that returns a logical boolean from a request is considered a _validating action_.

However, to implement an action you just use a single interface:

```Go
type Action interface {
	Name() string
	Run(context.Context, *unstructured.Unstructured) error
}
```

You can follow a few patterns to access your request in a Kubernetes `Unstructured` type. The built-in `Unstructured` interface which lets you drill down into the type that's contained, or use a library like [gabs](https://github.com/Jeffail/gabs) which slurps this object and gives you a better interface in terms of developer ergonomics:

```Go
o, err := gabs.Consume(p.Object)
img := container.Path("image").Data()
```