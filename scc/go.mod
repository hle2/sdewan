module github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc

require (
	go.etcd.io/etcd v3.3.12+incompatible
	google.golang.org/grpc v1.27.1
	k8s.io/client-go v0.19.0
)

replace (
        github.com/onap/multicloud-k8s/src/clm => ../clm
        github.com/onap/multicloud-k8s/src/orchestrator => ../orchestrator
        github.com/onap/multicloud-k8s/src/rsync => ../rsync
        k8s.io/api => k8s.io/api v0.19.0
        k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.19.0
        k8s.io/apimachinery => k8s.io/apimachinery v0.19.0
        k8s.io/cli-runtime => k8s.io/cli-runtime v0.19.0
        k8s.io/client-go => k8s.io/client-go v0.19.0
        k8s.io/kubectl => k8s.io/kubectl v0.19.0
        k8s.io/kubernetes => k8s.io/kubernetes v1.14.1
        k8s.io/apiserver => k8s.io/apiserver v0.0.0-20190409021813-1ec86e4da56c
        k8s.io/cloud-provider => k8s.io/cloud-provider v0.0.0-20190409023720-1bc0c81fa51d
)

go 1.14
