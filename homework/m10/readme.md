### 1、通过helm安装prometheus
安装文档：https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack
```shell
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

helm -n prometheus-stack install kube-prometheus-stack prometheus-community/kube-prometheus-stack 
```
安装完成之后通过如下命名查看prometheus
```shell
kubectl get prometheuses -n prometheus-stack
```

### 2、创建服务自动发现的规则
参考文章：https://www.qikqiak.com/post/prometheus-operator-advance/
把下面的yaml保存成文件（prometheus-additional.yaml），用于创建additional-configs
```yaml
prometheus-additional.yaml- job_name: 'kubernetes-endpoints'
  kubernetes_sd_configs:
    - role: endpoints
  relabel_configs:
    - source_labels: [__meta_kubernetes_service_annotation_prometheus_io_scrape]
      action: keep
      regex: true
    - source_labels: [__meta_kubernetes_service_annotation_prometheus_io_scheme]
      action: replace
      target_label: __scheme__
      regex: (https?)
    - source_labels: [__meta_kubernetes_service_annotation_prometheus_io_path]
      action: replace
      target_label: __metrics_path__
      regex: (.+)
    - source_labels: [__address__, __meta_kubernetes_service_annotation_prometheus_io_port]
      action: replace
      target_label: __address__
      regex: ([^:]+)(?::\d+)?;(\d+)
      replacement: $1:$2
    - action: labelmap
      regex: __meta_kubernetes_service_label_(.+)
    - source_labels: [__meta_kubernetes_namespace]
      action: replace
      target_label: kubernetes_namespace
    - source_labels: [__meta_kubernetes_service_name]
      action: replace
      target_label: kubernetes_name
    - source_labels: [__meta_kubernetes_pod_name]
      action: replace
      target_label: kubernetes_pod_name
```
### 3、把规则存入到集群的secret中
```shell
kubectl create secret generic additional-configs --from-file=prometheus-additional.yaml -n  prometheus-stack
```
创建完成之后，该信息会被base64编码后作为prometheus-additional.yaml 这个key对应的值存在 secret中
```shell
kubectl get secret additional-configs -n prometheus-stack -o yaml

apiVersion: v1
data:
  prometheus-additional.yaml: cHJvbWV0aGV1cy1hZGRpdGlvbmFsLnlhbWwtIGpvYl9uYW1lOiAna3ViZXJuZXRlcy1lbmRwb2ludHMnCiAga3ViZXJuZXRlc19zZF9jb25maWdzOgogICAgLSByb2xlOiBlbmRwb2ludHMKICByZWxhYmVsX2NvbmZpZ3M6CiAgICAtIHNvdXJjZV9sYWJlbHM6IFtfX21ldGFfa3ViZXJuZXRlc19zZXJ2aWNlX2Fubm90YXRpb25fcHJvbWV0aGV1c19pb19zY3JhcGVdCiAgICAgIGFjdGlvbjoga2VlcAogICAgICByZWdleDogdHJ1ZQogICAgLSBzb3VyY2VfbGFiZWxzOiBbX19tZXRhX2t1YmVybmV0ZXNfc2VydmljZV9hbm5vdGF0aW9uX3Byb21ldGhldXNfaW9fc2NoZW1lXQogICAgICBhY3Rpb246IHJlcGxhY2UKICAgICAgdGFyZ2V0X2xhYmVsOiBfX3NjaGVtZV9fCiAgICAgIHJlZ2V4OiAoaHR0cHM/KQogICAgLSBzb3VyY2VfbGFiZWxzOiBbX19tZXRhX2t1YmVybmV0ZXNfc2VydmljZV9hbm5vdGF0aW9uX3Byb21ldGhldXNfaW9fcGF0aF0KICAgICAgYWN0aW9uOiByZXBsYWNlCiAgICAgIHRhcmdldF9sYWJlbDogX19tZXRyaWNzX3BhdGhfXwogICAgICByZWdleDogKC4rKQogICAgLSBzb3VyY2VfbGFiZWxzOiBbX19hZGRyZXNzX18sIF9fbWV0YV9rdWJlcm5ldGVzX3NlcnZpY2VfYW5ub3RhdGlvbl9wcm9tZXRoZXVzX2lvX3BvcnRdCiAgICAgIGFjdGlvbjogcmVwbGFjZQogICAgICB0YXJnZXRfbGFiZWw6IF9fYWRkcmVzc19fCiAgICAgIHJlZ2V4OiAoW146XSspKD86OlxkKyk/OyhcZCspCiAgICAgIHJlcGxhY2VtZW50OiAkMTokMgogICAgLSBhY3Rpb246IGxhYmVsbWFwCiAgICAgIHJlZ2V4OiBfX21ldGFfa3ViZXJuZXRlc19zZXJ2aWNlX2xhYmVsXyguKykKICAgIC0gc291cmNlX2xhYmVsczogW19fbWV0YV9rdWJlcm5ldGVzX25hbWVzcGFjZV0KICAgICAgYWN0aW9uOiByZXBsYWNlCiAgICAgIHRhcmdldF9sYWJlbDoga3ViZXJuZXRlc19uYW1lc3BhY2UKICAgIC0gc291cmNlX2xhYmVsczogW19fbWV0YV9rdWJlcm5ldGVzX3NlcnZpY2VfbmFtZV0KICAgICAgYWN0aW9uOiByZXBsYWNlCiAgICAgIHRhcmdldF9sYWJlbDoga3ViZXJuZXRlc19uYW1lCiAgICAtIHNvdXJjZV9sYWJlbHM6IFtfX21ldGFfa3ViZXJuZXRlc19wb2RfbmFtZV0KICAgICAgYWN0aW9uOiByZXBsYWNlCiAgICAgIHRhcmdldF9sYWJlbDoga3ViZXJuZXRlc19wb2RfbmFtZQoK
kind: Secret
...
```
### 4、将该配置加入到声明prometheus的资源对象中，通过additionalScrapeConfigs属性
在spec下增加 additionalScrapeConfigs
参考文档：https://github.com/prometheus-operator/prometheus-operator/blob/main/Documentation/additional-scrape-config.md
```shell
kubectl edit prometheus kube-prometheus-stack-prometheus -n prometheus-stack

...
spec:
  additionalScrapeConfigs:
    key: prometheus-additional.yaml
    name: additional-configs
...
```
如果配置不生效，可以查看operator的日志
```shell
kubectl get po -n prometheus-stack

NAME                                                        READY   STATUS    RESTARTS   AGE
kube-prometheus-stack-operator-547c659bb-tcwml              1/1     Running   0          87m

kubectl logs kube-prometheus-stack-operator-547c659bb-tcwml
```
### 5、RBAC配置
如果查看exporter日志发现有很多forbidden，说明是RBAC的问题，需要给予get的权限， `kube-prometheus-stack-prometheus`是在上方创建prometheus时的`serviceAccountName`
```shell
kubectl get clusterrole kube-prometheus-stack-prometheus -oyaml
```