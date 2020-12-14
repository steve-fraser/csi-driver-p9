# CSI driver example

After the P9 CSI Driver is deployed in your cluster, you can follow this documentation to quickly deploy some examples. 

You can use P9 CSI Driver to provision Persistent Volumes statically or dynamically. Please read [Kubernetes Persistent Volumes documentation](https://kubernetes.io/docs/concepts/storage/persistent-volumes/) for more information about Static and Dynamic provisioning.

Please refer to [driver parameters](../../docs/driver-parameters.md) for more detailed usage.

## Prerequisite

- [Set up a P9 Server on a Kubernetes cluster](./p9-provisioner/README.md)
- [Install P9 CSI Driver](../../docs/install-csi-driver.md)

## Storage Class Usage (Dynamic Provisioning)

- Follow the folling command to create a `StorageClass`, and then `PersistentVolume` and `PersistentVolumeClaim` dynamically.

```bash
# create StorageClass
kubectl create -f https://raw.githubusercontent.com/kubernetes-csi/csi-driver-p9/master/deploy/example/storageclass-p9.yaml

# create PVC
kubectl create -f https://raw.githubusercontent.com/kubernetes-csi/csi-driver-p9/master/deploy/example/pvc-p9-csi-dynamic.yaml
```

## PV/PVC Usage (Static Provisioning)

- Follow the folling command to create `PersistentVolume` and `PersistentVolumeClaim` statically.

```bash
# create PV
kubectl create -f https://raw.githubusercontent.com/kubernetes-csi/csi-driver-p9/master/deploy/example/pv-p9-csi.yaml

# create PVC
kubectl create -f https://raw.githubusercontent.com/kubernetes-csi/csi-driver-p9/master/deploy/example/pvc-p9-csi-static.yaml
```
