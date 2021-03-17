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
    BaseUrl = OverlayUrl + "/overlay1/" + manager.DeviceCollection
    IprangeUrl := OverlayUrl + "/overlay1/" + manager.IPRangeCollection
    flag := true

    kube_config_B, _ := ioutil.ReadFile("admin.conf")
    encoded_config_b := base64.StdEncoding.EncodeToString([]byte(kube_config_B))

    var publicIpB []string

    var object1 = module.OverlayObject{
        Metadata: module.ObjectMetaData{"overlay1", "", "", ""}, 
        Specification: module.OverlayObjectSpec{}}
    var objecta = module.ProposalObject{
        Metadata: module.ObjectMetaData{"proposal1", "", "", ""}, 
        Specification: module.ProposalObjectSpec{"aes128", "sha256", "modp3072"}}
    var objectb = module.ProposalObject{
        Metadata: module.ObjectMetaData{"proposal2", "", "", ""}, 
        Specification: module.ProposalObjectSpec{"aes256", "sha256", "modp3072"}}
    var device = module.DeviceObject{
	    Metadata: module.ObjectMetaData{"device-a", "", "", ""},
	    Specification: module.DeviceObjectSpec{publicIpB, true, "", 65536, true, false, "sdewan-edge-a", encoded_config_b}}
    var iprange_object1 = module.IPRangeObject{
        Metadata: module.ObjectMetaData{"ipr1", "", "", ""}, 
        Specification: module.IPRangeObjectSpec{"192.168.0.2", 1, 15}}

    if flag {
    createControllerObject(OverlayUrl, &object1, &module.OverlayObject{})
    createControllerObject(ProposalUrl, &objecta, &module.ProposalObject{})
    createControllerObject(ProposalUrl, &objectb, &module.ProposalObject{})
    createControllerObject(IprangeUrl, &iprange_object1, &module.IPRangeObject{})
    createControllerObject(BaseUrl, &device, &module.DeviceObject{})
    }

    var ret = m.Run()

    deleteControllerObject(BaseUrl, "device-a")
    deleteControllerObject(IprangeUrl, "ipr1")
    deleteControllerObject(ProposalUrl, "proposal2")
    deleteControllerObject(ProposalUrl, "proposal1")
    //deleteControllerObject(OverlayUrl, "overlay1")

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

    var objs []module.DeviceObject
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
            object_name: "device-a",
        },
        {
            name: "GetFoolName",
            object_name: "foo_name",
            expectedErr: true,
            expectedErrCode: 500,
        },
    }

    for _, tcase := range tcases {
        _, err := getControllerObject(BaseUrl, tcase.object_name, &module.DeviceObject{})
        handleError(t, err, tcase.name, tcase.expectedErr, tcase.expectedErrCode)
    }
}

func TestCreateObject(t *testing.T) {
    var publicIp []string

    kube_config_B, err := ioutil.ReadFile("admin1.conf")
    if err != nil {
            fmt.Println(err)
    }
    encoded_config_b := base64.StdEncoding.EncodeToString([]byte(kube_config_B))

    tcases := []struct {
        name string
        obj module.DeviceObject
        expectedErr bool
        expectedErrCode int
    }{
        {
            name: "EmptyName",
            obj: module.DeviceObject{
                Metadata: module.ObjectMetaData{"", "object 1", "", ""},
                Specification: module.DeviceObjectSpec{publicIp, true, "", 65536, true, false, "emptyobject", encoded_config_b}},
            expectedErr: true,
            expectedErrCode: 422,
        },
    }

    for _, tcase := range tcases {
        _, err := createControllerObject(BaseUrl, &tcase.obj, &module.DeviceObject{})
        handleError(t, err, tcase.name, tcase.expectedErr, tcase.expectedErrCode)
    }
}

