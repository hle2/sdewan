package test

import (
    "testing"
    "io/ioutil"
    "flag"
    "encoding/json"
    "encoding/base64"
    "fmt"
    "os"
    "github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/module"
    "github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/manager"
)

var BaseUrl string

func TestMain(m *testing.M) {
    servIp := flag.String("ip", "127.0.0.1", "SDEWAN Central Controller IP Address")
    flag.Parse()
    OverlayUrl := "http://" + *servIp + ":9015/scc/v1/" + manager.OverlayCollection
    ProposalUrl := OverlayUrl + "/overlay1/" + manager.ProposalCollection
    BaseUrl = OverlayUrl + "/overlay1/" + manager.HubCollection

    kube_config_A, err := ioutil.ReadFile("localhost.conf")
    if err != nil {
            fmt.Println(err)
    }
    encoded_config_a := base64.StdEncoding.EncodeToString([]byte(kube_config_A))
    /*
    kube_config_B, err := ioutil.ReadFile("admin.conf")
    encoded_config_b := base64.StdEncoding.EncodeToString([]byte(kube_config_B))
    */

    var publicIpA []string
    var publicIpB []string
    publicIpA = append(publicIpA, "192.168.121.23")
    publicIpB = append(publicIpB, "192.168.121.28")

    var object1 = module.OverlayObject{
        Metadata: module.ObjectMetaData{"overlay1", "", "", ""}, 
        Specification: module.OverlayObjectSpec{}}
    var objecta = module.ProposalObject{
        Metadata: module.ObjectMetaData{"proposal1", "", "", ""}, 
        Specification: module.ProposalObjectSpec{"aes128", "sha256", "modp3072"}}
    var objectb = module.ProposalObject{
        Metadata: module.ObjectMetaData{"proposal2", "", "", ""}, 
        Specification: module.ProposalObjectSpec{"aes256", "sha256", "modp3072"}}
    var object2 = module.HubObject{
        Metadata: module.ObjectMetaData{"huba", "", "", ""},
        Specification: module.HubObjectSpec{publicIpA, "192.168.121.23", encoded_config_a}}
    /*var object3 = module.HubObject{
        Metadata: module.ObjectMetaData{"hubB", "", "", ""}, 
        Specification: module.HubObjectSpec{publicIpB, "192.168.121.28", encoded_config_b}}*/

    createControllerObject(OverlayUrl, &object1, &module.OverlayObject{})
    createControllerObject(ProposalUrl, &objecta, &module.ProposalObject{})
    createControllerObject(ProposalUrl, &objectb, &module.ProposalObject{})
    createControllerObject(BaseUrl, &object2, &module.HubObject{})
    //createControllerObject(BaseUrl, &object3, &module.HubObject{})

    var ret = m.Run()

    deleteControllerObject(BaseUrl, "huba")
    deleteControllerObject(ProposalUrl, "proposal2")
    deleteControllerObject(ProposalUrl, "proposal1")
    deleteControllerObject(OverlayUrl, "overlay1")

    os.Exit(ret)
}

func TestGetObjects(t *testing.T) {
    url := BaseUrl
    res, err := callRest("GET", url, "")
    if err != nil {
        printError(err)
        t.Errorf("Test case GetObjects: can not get Objects")
        return
    }

    var objs []module.HubObject
    err = json.Unmarshal([]byte(res), &objs)

    if len(objs) == 0 {
        fmt.Printf("Test case GetObjects: no object found")
        return
    }

    p_data, _ := json.Marshal(objs)
    fmt.Printf("%s\n", string(p_data))
}


func TestGetObject(t *testing.T) {
    tcases := []struct {
        name string
        object_name string
        expectedErr bool
        expectedErrCode int
    }{
        {
            name: "Normal",
            object_name: "huba",
        },
        {
            name: "GetFoolName",
            object_name: "foo_name",
            expectedErr: true,
            expectedErrCode: 500,
        },
    }

    for _, tcase := range tcases {
        _, err := getControllerObject(BaseUrl, tcase.object_name, &module.HubObject{})
        handleError(t, err, tcase.name, tcase.expectedErr, tcase.expectedErrCode)
    }
}

func TestCreateObject(t *testing.T) {
    var publicIp []string
    publicIp = append(publicIp, "1.1.1.1")
    kube_config_B, err := ioutil.ReadFile("admin.conf")
    if err != nil {
            fmt.Println(err)
    }
    encoded_config_b := base64.StdEncoding.EncodeToString([]byte(kube_config_B))

    tcases := []struct {
        name string
        obj module.HubObject
        expectedErr bool
        expectedErrCode int
    }{
        {
            name: "EmptyName",
            obj: module.HubObject{
                Metadata: module.ObjectMetaData{"", "object 1", "", ""},
                Specification: module.HubObjectSpec{publicIp, "1.1.1.1", string(encoded_config_b)}},
            expectedErr: true,
            expectedErrCode: 422,
        },
    }

    for _, tcase := range tcases {
        _, err := createControllerObject(BaseUrl, &tcase.obj, &module.HubObject{})
        handleError(t, err, tcase.name, tcase.expectedErr, tcase.expectedErrCode)
    }
}

func TestCreateObjectPass(t *testing.T) {
    var publicIp []string
    publicIp = append(publicIp, "192.168.121.28")
    kube_config_B, err := ioutil.ReadFile("admin.conf")
    if err != nil {
            fmt.Println(err)
    }
    encoded_config_b := base64.StdEncoding.EncodeToString([]byte(kube_config_B))

    tcases := []struct {
        name string
        obj module.HubObject
        expectedErr bool
        expectedErrCode int
    }{
        {
            name: "Normal",
            obj: module.HubObject{
                Metadata: module.ObjectMetaData{"hubtest", "object 4", "", ""},
                Specification: module.HubObjectSpec{publicIp, "192.168.121.28", string(encoded_config_b)}},
        },
    }

    for _, tcase := range tcases {
        _, err := createControllerObject(BaseUrl, &tcase.obj, &module.HubObject{})
        handleError(t, err, tcase.name, tcase.expectedErr, tcase.expectedErrCode)
    }

    gcases := []struct {
        name string
        object_name string
        expectedErr bool
        expectedErrCode int
    }{
        {
            name: "NormalGet",
            object_name: "hubtest",
        },
    }

    for _, gcase := range gcases {
        _, err := getControllerObject(BaseUrl, gcase.object_name, &module.HubObject{})
        handleError(t, err, gcase.name, gcase.expectedErr, gcase.expectedErrCode)
    }
}

