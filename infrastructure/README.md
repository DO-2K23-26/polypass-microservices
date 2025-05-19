# Infrastructure
---

This part detail how the polypass application is hosted.

## Requirements

You should have on you local machine:
- helm v3.13.1
- kubectl v1.29.0
- istioctl v1.24.3

## Vm

A virtual machine on Serdaigle Proxmox cluster has been deployed with:
- 24 go RAM
- 6 Cpu
- Debian 12

## Kubernetes cluster

The cluster has been deployed using [k3s](https://k3s.io/) and the default command:
```sh
curl -sfL https://get.k3s.io | sh -
```

## Service mesh and network

We will use Istio as a service mesh. Istio manage its own gateway service and k3s install by default traefik IngressController.
So we will begin by removing the helm chart:
```sh
helm uninstall traefik -n kube-system
```

The cluster is ready to host Istio. We will use the [default profile](https://istio.io/latest/docs/setup/additional-setup/config-profiles/#deployment-profiles) to install istio.
It will allow the install of [`istio-ingressgateway`](https://istio.io/latest/docs/tasks/traffic-management/ingress/ingress-control/) and therefore we will be able to manage ingress traffic to our services.
You can use the [command](https://istio.io/latest/docs/setup/install/istioctl/):
```sh
istioctl install -f ./istio_operator.yaml
```

## Observability


[Jaeger install](https://istio.io/latest/docs/ops/integrations/jaeger/):
```sh
kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.25/samples/addons/jaeger.yaml
kubectl apply -f jaeger.yaml
```

[Install Otel](https://istio.io/latest/docs/tasks/observability/distributed-tracing/opentelemetry/):
```sh
kubectl create namespace observability
kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.25/samples/open-telemetry/otel.yaml -n observability
kubectl apply -f otel.yaml
```

[Install Prometheus](https://istio.io/latest/docs/ops/integrations/prometheus/#configuration):
```sh
kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.25/samples/addons/prometheus.yaml
```


[Install grafana](https://istio.io/latest/docs/ops/integrations/grafana/):
```sh
kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.25/samples/addons/grafana.yaml
```


Install a mesh vizualizer [Kiali](https://istio.io/latest/docs/ops/integrations/kiali/):
```sh
kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.25/samples/addons/kiali.yaml
```


## Confluent

[Prepare the cluster](https://docs.confluent.io/operator/current/co-prepare.html#quick-easy-deployment-and-relaxed-permissions):
```sh
kubectl create namespace confluent
helm repo add confluentinc https://packages.confluent.io/helm
kubectl apply -f confluent/confluent-for-kubernetes/crds
# As this crd is too long you have to force it see
# https://docs.confluent.io/operator/current/co-troubleshooting.html#issue-an-error-returns-while-applying-a-crd-during-an-upgrade
kubectl apply --server-side=true --force-conflicts -f confluent/confluent-for-kubernetes/crds/platform.confluent.io_kafkas.yaml
```

Install confluent-for-kubernetes:
```sh
helm upgrade --install confluent-operator \
  confluentinc/confluent-for-kubernetes \
  --namespace confluent
```

Install confluent platform (Kafka, Schema registry and zookeeper):
```sh
kubectl apply -f confluent-platform.yaml
```

It can take up to 10 min to deploy.


## Deploy Authzed Operator


[Deploy the operator:](https://github.com/authzed/spicedb-operator?tab=readme-ov-file#getting-started)
```sh
kubectl apply --server-side -f https://github.com/authzed/spicedb-operator/releases/latest/download/bundle.yaml
```

Create a namespace for the authorization service:
```sh
kubectl create namespace authz
```

Deploy the spicedb cluster:
```sh
kubectl apply -f ./authzed/cluster-spicedb.yaml
```