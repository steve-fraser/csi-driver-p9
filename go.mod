module github.com/kubernetes-csi/csi-driver-p9

go 1.13

require (
	github.com/container-storage-interface/spec v1.3.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.4.2
	github.com/kubernetes-csi/csi-lib-utils v0.9.0
	github.com/kubernetes-csi/external-snapshotter/v2 v2.0.0-20200617021606-4800ca72d403
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.7.1
	github.com/pborman/uuid v1.2.0
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.5.1
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	google.golang.org/grpc v1.29.0
	google.golang.org/protobuf v1.25.0 // indirect
	k8s.io/api v0.19.0
	k8s.io/apimachinery v0.19.0
	k8s.io/client-go v0.19.0
	k8s.io/klog/v2 v2.2.0
	k8s.io/kubernetes v1.19.0
	k8s.io/utils v0.0.0-20200729134348-d5654de09c73
)

replace k8s.io/api => k8s.io/api v0.19.0

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.19.0

replace k8s.io/apimachinery => k8s.io/apimachinery v0.19.0

replace k8s.io/apiserver => k8s.io/apiserver v0.19.0

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.19.0

replace k8s.io/client-go => k8s.io/client-go v0.19.0

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.19.0

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.19.0

replace k8s.io/code-generator => k8s.io/code-generator v0.19.0

replace k8s.io/component-base => k8s.io/component-base v0.19.0

replace k8s.io/cri-api => k8s.io/cri-api v0.19.0

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.19.0

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.19.0

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.19.0

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.19.0

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.19.0

replace k8s.io/kubectl => k8s.io/kubectl v0.19.0

replace k8s.io/kubelet => k8s.io/kubelet v0.19.0

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.19.0

replace k8s.io/metrics => k8s.io/metrics v0.19.0

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.19.0

replace k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.19.0

replace k8s.io/sample-controller => k8s.io/sample-controller v0.19.0
