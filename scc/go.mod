module github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc

require (
        github.com/open-ness/EMCO/src/orchestrator v0.0.0-00010101000000-000000000000
	go.etcd.io/etcd v3.3.12+incompatible
	google.golang.org/grpc v1.27.1
	k8s.io/client-go v0.19.0
)

replace (
        github.com/open-ness/EMCO/src/orchestrator => ../vendor/github.com/open-ness/EMCO/src/orchestrator
        github.com/open-ness/EMCO/src/rsync => ../rsync
        github.com/open-ness/EMCO/src/monitor => ../monitor
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