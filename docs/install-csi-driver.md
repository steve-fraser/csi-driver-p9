# Install P9 CSI driver master version on a kubernetes cluster

If you have already installed Helm, you can also use it to install P9 CSI driver. Please see [Installation with Helm](../charts/README.md).

## Install with kubectl
 - remote install
```console
curl -skSL https://raw.githubusercontent.com/kubernetes-csi/csi-driver-p9/master/deploy/install-driver.sh | bash -s master --
```

 - local install
```console
git clone https://github.com/kubernetes-csi/csi-driver-p9.git
cd csi-driver-p9
./deploy/install-driver.sh master local
```

- check pods status:
```console
kubectl -n kube-system get pod -o wide -l app=csi-p9-controller
kubectl -n kube-system get pod -o wide -l app=csi-p9-node
```

example output:

```console
NAME                                       READY   STATUS    RESTARTS   AGE     IP             NODE
csi-p9-controller-56bfddd689-dh5tk       4/4     Running   0          35s     10.240.0.19    k8s-agentpool-22533604-0
csi-p9-controller-56bfddd689-8pgr4       4/4     Running   0          35s     10.240.0.35    k8s-agentpool-22533604-1
csi-p9-node-cvgbs                        3/3     Running   0          35s     10.240.0.35    k8s-agentpool-22533604-1
csi-p9-node-dr4s4                        3/3     Running   0          35s     10.240.0.4     k8s-agentpool-22533604-0
```

- clean up P9 CSI driver
```console
curl -skSL https://raw.githubusercontent.com/kubernetes-csi/csi-driver-p9/master/deploy/uninstall-driver.sh | bash -s master --
```