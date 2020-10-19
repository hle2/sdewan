package main

import (
    "log"
    "encoding/json"
    "github.com/onap/multicloud-k8s/src/orchestrator/pkg/infra/db"
    "github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/manager"
    "github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/module"
    "github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/resource"
)

func createConnection(end1 string, end2 string) error {
    var publicIp []string
    publicIp = append(publicIp, "1.1.1.1")
    m1 := module.HubObject{
        Metadata: module.ObjectMetaData{end1, end1 + "-" + "object 1", "", ""},
        Specification: module.HubObjectSpec{publicIp, "1.1.1.1", ""}}
    m2 := module.HubObject{
        Metadata: module.ObjectMetaData{end2, end2 + "-" + "object 2", "", ""},
        Specification: module.HubObjectSpec{publicIp, "1.1.1.1", ""}}

    cend1 := module.NewConnectionEnd(&m1, "127.0.0.1")
    cend1.AddResource(&resource.FileResource{"mycm", "ConfigMap", "mycm.yaml"}, false)
    cend1.AddResource(&resource.FileResource{"mycm", "ConfigMap", "mycm.yaml"}, false)
    cend2 := module.NewConnectionEnd(&m2, "127.0.0.1")
    cend2.AddResource(&resource.FileResource{"mycm1", "ConfigMap", "mycm.yaml"}, false)
    cend2.AddResource(&resource.FileResource{"mycm2", "ConfigMap", "mycm.yaml"}, false)
    co := module.NewConnectionObject(cend1, cend2)
    cm := manager.GetConnectionManager()

    _, err := cm.UpdateObject("overlay", co)
    return err
}

func main() {
    // create database and context database
    err := db.InitializeDatabaseConnection("scc_test")
    if err != nil {
        log.Println("Unable to initialize database connection...")
        log.Println(err)
        log.Fatalln("Exiting...")
    }

    createConnection("hub1", "hub2")
    createConnection("hub1", "hub3")
    createConnection("hub1", "hub4")
    createConnection("hub2", "hub3")
    createConnection("hub3", "hub4")

    cm := manager.GetConnectionManager()
    cn, err := cm.GetObject("overlay", "Hub.hub2", "Hub.hub1")
    if err != nil {
        log.Println(err)
    } else {
        p_data, _ := json.Marshal(cn)
        log.Println(string(p_data))
    }

    cns, err := cm.GetObjects("overlay", "Hub.hub2")
    if err != nil {
        log.Println(err)
    } else {
        p_data, _ := json.Marshal(cns)
        log.Println(string(p_data))
    }

    cns, err = cm.GetObjects("overlay", "Hub.hub1")
    if err != nil {
        log.Println(err)
    } else {
        p_data, _ := json.Marshal(cns)
        log.Println(string(p_data))
    }
}
