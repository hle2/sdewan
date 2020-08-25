package main

import (
    "github.com/onap/multicloud-k8s/src/orchestrator/pkg/infra/db"
    "log"
    "math/rand"
    "time"
    "io/ioutil"
    "encoding/base64"
    pkgerrors "github.com/pkg/errors"

    mtypes "github.com/onap/multicloud-k8s/src/orchestrator/pkg/module/types"
)

type Cluster struct {
    Metadata mtypes.Metadata `json:"metadata"`
}

type ClusterContent struct {
    Kubeconfig string `json:"kubeconfig"`
}

type ClusterKey struct {
    ClusterProviderName string `json:"provider"`
    ClusterName         string `json:"cluster"`
}

func registerCluster(provider_name string, cluster_name string, kubeconfig_file string) error {
    var q ClusterContent
    var p Cluster

    content, err := ioutil.ReadFile(kubeconfig_file)
    q.Kubeconfig = base64.StdEncoding.EncodeToString(content)
    key := ClusterKey{
        ClusterProviderName: provider_name,
        ClusterName:         cluster_name,
    }

    p.Metadata.Name = cluster_name

    err = db.DBconn.Insert("cluster", key, nil, "clustermetadata", p)
    if err != nil {
        return pkgerrors.Wrap(err, "Creating DB Entry")
    }

    err = db.DBconn.Insert("cluster", key, nil, "clustercontent", q)
    if err != nil {
        return pkgerrors.Wrap(err, "Creating DB Entry")
    }

    return nil
}

func main() {
    rand.Seed(time.Now().UnixNano())

    // Initialize the mongodb
    err := db.InitializeDatabaseConnection("scc")
    if err != nil {
        log.Println("Unable to initialize database connection...")
        log.Println(err)
        log.Fatalln("Exiting...")
    }

    provider_name := "akraino_scc"
    cluster_name := "local"
    // Register cluster kubeconfig
    registerCluster(provider_name, cluster_name, "admin.conf")    
}
