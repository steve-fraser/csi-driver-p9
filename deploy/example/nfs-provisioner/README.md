# Set up a P9 Server on a Kubernetes cluster

After the P9 CSI Driver is deployed in your cluster, you can follow this documentation to quickly deploy some example applications. You can use P9 CSI Driver to provision Persistent Volumes statically or dynamically. Please read Kubernetes Persistent Volumes for more information about Static and Dynamic provisioning.

There are multiple different P9 servers you can use for testing of 
the plugin, the major versions of the protocol v2, v3 and v4 should be supported
by the current implementation. This page will show you how to set up a P9 Server deployment on a Kubernetes cluster.

- To create a P9 provisioner on your Kubernetes cluster, run the following command.

```bash
kubectl create -f https://raw.githubusercontent.com/kubernetes-csi/csi-driver-p9/master/deploy/example/p9-provisioner/p9-server.yaml
```

- During the deployment, a new service `p9-server` will be created which exposes the P9 server endpoint `p9-server.default.svc.cluster.local` and the share path `/`. You can specify `PersistentVolume` or `StorageClass` using these information.

- Deploy the P9 CSI driver, please refer to [install P9 CSI driver](../../../docs/install-csi-driver.md).

- To check if the P9 server is working, we can statically create a PersistentVolume and a PersistentVolumeClaim, and mount it onto a sample pod:

```bash
kubectl create -f https://raw.githubusercontent.com/kubernetes-csi/csi-driver-p9/master/deploy/example/p9-provisioner/nginx-pod.yaml
```

 - Verify if the P9 server is functional, you can check the mount point from the example pod.

 ```bash
kubectl exec nginx-p9-example -- bash -c "findmnt /var/www -o TARGET,SOURCE,FSTYPE"
```

 - The output should look like the following:

 ```bash
TARGET   SOURCE                                 FSTYPE
/var/www p9-server.default.svc.cluster.local:/ p94
```
