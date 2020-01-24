---
id: configmap
title: Using Configmap
sidebar_label: Using Configmap
---


If you use our Helm deployment, you get this for free. A Lightscreen deployment that has its configuration baked into a dynamic configmap:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: lightscreen-config
  labels:
{{ include "lightscreen.labels" . | indent 4 }}
data:
  lightscreen.yaml: {{.Files.Get "files/lightscreen.yaml" | printf "%s" | quote }}
```

In this example, taken directly from the Helm chart, we read the local `lightscreen.yaml` file and put it into a configmap.


Using the same technique you can custom-bake your own workflow. The only condition is that you restart Lightscreen to pick up the new configuration.

In your `lightscreen.yaml` file you can update admission whitelists, add or remove actions, maintain a SHA resolving list and more.
