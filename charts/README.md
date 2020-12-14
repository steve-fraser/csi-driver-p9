# Installation with Helm 3

Follow this guide to install the P9 Driver for Kubernetes.

## Prerequisites

- [Install Helm 3](https://helm.sh/docs/intro/quickstart/#install-helm)

## Install via `helm install`

```
$ cd charts/latest
$ helm install csi-driver-p9 ./csi-driver-p9 -n kube-system
```
## Install via Helm repository

```
$ helm repo add csi-driver-p9 https://raw.githubusercontent.com/kubernetes-csi/csi-driver-p9/master/charts
$ helm install --name csi-driver-p9 csi-driver-p9/csi-driver-p9 --namespace kube-system
```

### Search for available versions

```
$ helm search repo -l csi-driver-p9
```

### Install a specific version

```
https://raw.githubusercontent.com/kubernetes-csi/csi-driver-p9/master/charts --version v2.0.0
```

## Chart configuration

The following table lists the configurable parameters of the latest P9 CSI Driver chart and their default values.

| Parameter                                         | Description                                                | Default                                                           |
|---------------------------------------------------|------------------------------------------------------------|-------------------------------------------------------------------|
| `image.p9.repository`                            | csi-driver-p9 docker image                                | gcr.io/k8s-staging-sig-storage/p9plugin                          |
| `image.p9.tag`                                   | csi-driver-p9 docker image tag                            | amd64-linux-canary                                                |
| `image.p9.pullPolicy`                            | csi-driver-p9 image pull policy                           | IfNotPresent                                                      |
| `image.csiProvisioner.repository`                 | csi-provisioner docker image                               | k8s.gcr.io/sig-storage/csi-provisioner                            |
| `image.csiProvisioner.tag`                        | csi-provisioner docker image tag                           | v2.0.4                                                            |
| `image.csiProvisioner.pullPolicy`                 | csi-provisioner image pull policy                          | IfNotPresent                                                      |
| `image.livenessProbe.repository`                  | liveness-probe docker image                                | k8s.gcr.io/sig-storage/livenessprobe                              |
| `image.livenessProbe.tag`                         | liveness-probe docker image tag                            | v2.1.0                                                            |
| `image.livenessProbe.pullPolicy`                  | liveness-probe image pull policy                           | IfNotPresent                                                      |
| `image.nodeDriverRegistrar.repository`            | csi-node-driver-registrar docker image                     | k8s.gcr.io/sig-storage/csi-node-driver-registrar                  |
| `image.nodeDriverRegistrar.tag`                   | csi-node-driver-registrar docker image tag                 | v2.0.1                                                            |
| `image.nodeDriverRegistrar.pullPolicy`            | csi-node-driver-registrar image pull policy                | IfNotPresent                                                      |
| `serviceAccount.create`                           | whether create service account of csi-p9-controller       | true                                                              |
| `rbac.create`                                     | whether create rbac of csi-p9-controller                  | true                                                              |
| `controller.replicas`                             | the replicas of csi-p9-controller                         | 2                                                                 |

## Troubleshooting

If there are some errors when using helm to install, follow the steps to debug:

1. Add `--wait -v=5 --debug` in `helm install` command.
2. Then the error pods  can be located.
3. Use `kubectl describe` to acquire more info.
4. Check the related resource of the pod, such as serviceaacount, rbac, etc.