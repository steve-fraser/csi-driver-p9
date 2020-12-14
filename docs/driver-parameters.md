## Driver Parameters
> This plugin driver itself only provides a communication layer between resources in the cluser and the P9 server, you need to bring your own P9 server before using this driver.

### Storage Class Usage (Dynamic Provisioning)
> [`StorageClass` example](../deploy/example/storageclass-p9.yaml)

Name | Meaning | Example Value | Mandatory | Default value
--- | --- | --- | --- | ---
server | P9 Server endpoint | Domain name `p9-server.default.svc.cluster.local` <br>Or IP address `127.0.0.1` | Yes |
share | P9 share path | `/` | Yes |

### PV/PVC Usage (Static Provisioning)
> [`PersistentVolume` example](../deploy/example/pv-p9-csi.yaml)

Name | Meaning | Example Value | Mandatory | Default value
--- | --- | --- | --- | ---
volumeAttributes.source | P9 Server endpoint | Domain name `p9-server.default.svc.cluster.local` <br>Or IP address `127.0.0.1` | Yes |
volumeAttributes.share | P9 share path | `/` |  Yes  |
