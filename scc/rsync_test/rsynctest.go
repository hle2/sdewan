package main

import (
    "github.com/onap/multicloud-k8s/src/orchestrator/pkg/appcontext"

    rsyncclient "github.com/onap/multicloud-k8s/src/orchestrator/pkg/grpc/installappclient"
    contextDb "github.com/onap/multicloud-k8s/src/orchestrator/pkg/infra/contextdb"
    "github.com/onap/multicloud-k8s/src/orchestrator/pkg/infra/db"
    "github.com/onap/multicloud-k8s/src/orchestrator/pkg/infra/rpc"
    controller "github.com/onap/multicloud-k8s/src/orchestrator/pkg/module/controller"
    "log"
    "math/rand"
    "time"
    "io/ioutil"
    "fmt"
    "encoding/json"
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

type resource struct {
    name        string
    filecontent string
}

type contextForCompositeApp struct {
    context            appcontext.AppContext
    ctxval             interface{}
    compositeAppHandle interface{}
}

func InitRsyncClient() bool {
    client := controller.NewControllerClient()

    vals, _ := client.GetControllers()
    found := false
    for _, v := range vals {
        if v.Metadata.Name == "rsync" {
            log.Println("Initializing RPC connection to resource synchronizer")
            rpc.UpdateRpcConn(v.Metadata.Name, v.Spec.Host, v.Spec.Port)
            found = true
            break
        }
    }
    return found
}

func makeAppContextForCompositeApp(p, ca, v, rName string) (contextForCompositeApp, error) {
    // ctxval: context.rtcObj.id
    context := appcontext.AppContext{}
    ctxval, err := context.InitAppContext()
    if err != nil {
        return contextForCompositeApp{}, pkgerrors.Wrap(err, "Error creating AppContext CompositeApp")
    }
    // compositeHandle = context.rtc.cid
    // context.rtc.RtcCreate(): (1) save (cid, id) in etcd  
    compositeHandle, err := context.CreateCompositeApp()
    if err != nil {
        return contextForCompositeApp{}, pkgerrors.Wrap(err, "Error creating CompositeApp handle")
    }
    // (1) set context.rtcObj.meta (2) save (cid + meta +"/", json.Marshal(rtc.meta)) in etcd
    err = context.AddCompositeAppMeta(appcontext.CompositeAppMeta{Project: p, CompositeApp: ca, Version: v, Release: rName})
    if err != nil {
        return contextForCompositeApp{}, pkgerrors.Wrap(err, "Error Adding CompositeAppMeta")
    }

    // return CompositeAppMeta{Project: p, CompositeApp: ca, Version: v, Release: rn}
    _, err = context.GetCompositeAppMeta()

    log.Println(":: The meta data stored in the runtime context :: ")

    // cca := contextForCompositeApp{context: appcontext.AppContext, ctxval: context.rtcObj.id, compositeAppHandle: context.rtc.cid}
    cca := contextForCompositeApp{context: context, ctxval: ctxval, compositeAppHandle: compositeHandle}

    return cca, nil
}

func getResources() ([]resource, error) {
    var resources []resource
    yamlFile, _ := ioutil.ReadFile("mycm.yaml")
    resources = append(resources, resource{name: "mycm+ConfigMap", filecontent: string(yamlFile)})

    return resources, nil
}

func addResourcesToCluster(ct appcontext.AppContext, ch interface{}, resources []resource) error {

    var resOrderInstr struct {
        Resorder []string `json:"resorder"`
    }

    var resDepInstr struct {
        Resdep map[string]string `json:"resdependency"`
    }
    resdep := make(map[string]string)

    for _, resource := range resources {
        resOrderInstr.Resorder = append(resOrderInstr.Resorder, resource.name)
        resdep[resource.name] = "go"
        // rtc.RtcAddResource("<cid>/app/app_name/cluster/clusername/", res.name, res.content)
        // -> save ("<cid>/app/app_name/cluster/clusername/resource/res.name/", res.content) in etcd
        // return ("<cid>/app/app_name/cluster/clusername/resource/res.name/"
        _, err := ct.AddResource(ch, resource.name, resource.filecontent)
        if err != nil {
            cleanuperr := ct.DeleteCompositeApp()
            if cleanuperr != nil {
                log.Printf(":: Error Cleaning up AppContext after add resource failure ::")
            }
            return pkgerrors.Wrapf(err, "Error adding resource ::%s to AppContext", resource.name)
        }
        jresOrderInstr, _ := json.Marshal(resOrderInstr)
        resDepInstr.Resdep = resdep
        jresDepInstr, _ := json.Marshal(resDepInstr)
        // rtc.RtcAddInstruction("<cid>app/app_name/cluster/clusername/", "resource", "order", "{[res.name]}")
        // ->save ("<cid>/app/app_name/cluster/clusername/resource/instruction/order/", "{[res.name]}") in etcd
        // return "<cid>/app/app_name/cluster/clusername/resource/instruction/order/"
        _, err = ct.AddInstruction(ch, "resource", "order", string(jresOrderInstr))
        _, err = ct.AddInstruction(ch, "resource", "dependency", string(jresDepInstr))
        if err != nil {
            cleanuperr := ct.DeleteCompositeApp()
            if cleanuperr != nil {
                log.Printf(":: Error Cleaning up AppContext after add instruction failure ::")
            }
            return pkgerrors.Wrapf(err, "Error adding instruction for resource ::%s to AppContext", resource.name)
        }
    }
    return nil
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
    err := db.InitializeDatabaseConnection("mco")
    if err != nil {
        log.Println("Unable to initialize database connection...")
        log.Println(err)
        log.Fatalln("Exiting...")
    }

    // Initialize contextdb
    err = contextDb.InitializeContextDatabase()
    if err != nil {
        log.Println("Unable to initialize database connection...")
        log.Println(err)
        log.Fatalln("Exiting...")
    }

    InitRsyncClient()

    provider_name := "sdewan"
    cluster_name := "local"
    // Register cluster kubeconfig
    registerCluster(provider_name, cluster_name, "admin.conf")
    
    // Generate Application context
    cca, err := makeAppContextForCompositeApp("sdewan", "app", "1", "1")
    context := cca.context  // appcontext.AppContext
    ctxval := cca.ctxval    // id
    compositeHandle := cca.compositeAppHandle // cid

    var appOrderInstr struct {
        Apporder []string `json:"apporder"`
    }

    var appDepInstr struct {
        Appdep map[string]string `json:"appdependency"`
    }
    appdep := make(map[string]string)

    // Add application
    app_name := "mytestapp"
    appOrderInstr.Apporder = append(appOrderInstr.Apporder, app_name)
    appdep[app_name] = "go"

    // rtc.RtcAddLevel(cid, "app", app_name) -> save ("<cid>app/app_name/", app_name) in etcd
    // apphandle = "<cid>app/app_name/"
    apphandle, err := context.AddApp(compositeHandle, app_name)
    resources, err := getResources()

    // Add cluster
    // err = addClustersToAppContext(listOfClusters, context, apphandle, resources)
    // rtc.RtcAddLevel("<cid>app/app_name/", "cluster", clustername) 
    // -> save ("<cid>app/app_name/cluster/clusername/", clustername) in etcd
    // return "<cid>app/app_name/cluster/clusername/"
    clusterhandle, err := context.AddCluster(apphandle, provider_name+"+"+cluster_name)
    err = addResourcesToCluster(context, clusterhandle, resources)

    jappOrderInstr, _ := json.Marshal(appOrderInstr)
    appDepInstr.Appdep = appdep
    jappDepInstr, _ := json.Marshal(appDepInstr)
    context.AddInstruction(compositeHandle, "app", "order", string(jappOrderInstr))
    context.AddInstruction(compositeHandle, "app", "dependency", string(jappDepInstr))

    // invoke deployment prrocess
    appContextID := fmt.Sprintf("%v", ctxval)
    err = rsyncclient.InvokeInstallApp(appContextID)
    if err != nil {
        log.Println(err)
    }
}