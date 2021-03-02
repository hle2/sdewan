package test

import (
    "testing"
    "io/ioutil"
    "flag"
    "encoding/json"
    "encoding/base64"
    "fmt"
    "os"
    "time"
    "log"
    "github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/module"
    "github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/manager"
)

var BaseUrl string
var CertUrl string

func TestMain(m *testing.M) {
    servIp := flag.String("ip", "127.0.0.1", "SDEWAN Central Controller IP Address")
    flag.Parse()
    OverlayUrl := "http://" + *servIp + ":9015/scc/v1/" + manager.OverlayCollection
    ProposalUrl := OverlayUrl + "/overlay1/" + manager.ProposalCollection
    DeviceUrl := OverlayUrl + "/overlay1/" + manager.DeviceCollection
    HubUrl := OverlayUrl + "/overlay1/" + manager.HubCollection
    IprangeUrl := OverlayUrl + "/overlay1/" + manager.IPRangeCollection
    CertUrl := OverlayUrl + "/overlay1/" + manager.CertCollection
    BaseUrl := OverlayUrl + "/overlay1/" + manager.HubCollection + "/huba/" + manager.DeviceCollection
    flag := false 

    kube_config_A, err := ioutil.ReadFile("huba.conf")
    if err != nil {
            fmt.Println(err)
    }
    encoded_config_a := base64.StdEncoding.EncodeToString([]byte(kube_config_A))

    kube_config_B, _ := ioutil.ReadFile("admin.conf")
    encoded_config_b := base64.StdEncoding.EncodeToString([]byte(kube_config_B))

    var publicIpA []string
    var publicIpB []string
    publicIpA = append(publicIpA, "10.10.20.17")
    //publicIpB = append(publicIpB, "10.10.10.16")

    var object1 = module.OverlayObject{
        Metadata: module.ObjectMetaData{"overlay1", "", "", ""}, 
        Specification: module.OverlayObjectSpec{}}
    var objecta = module.ProposalObject{
        Metadata: module.ObjectMetaData{"proposal1", "", "", ""}, 
        Specification: module.ProposalObjectSpec{"aes128", "sha256", "modp3072"}}
    var objectb = module.ProposalObject{
        Metadata: module.ObjectMetaData{"proposal2", "", "", ""}, 
        Specification: module.ProposalObjectSpec{"aes256", "sha256", "modp3072"}}
    var hub = module.HubObject{
        Metadata: module.ObjectMetaData{"huba", "", "", ""},
        Specification: module.HubObjectSpec{publicIpA, "10.10.10.15", encoded_config_a}}
    var device = module.DeviceObject{
	    Metadata: module.ObjectMetaData{"device-a", "", "", ""},
	    Specification: module.DeviceObjectSpec{publicIpB, true, "", 65536, true, false, "sdewan-edge-a", encoded_config_b}}
    var iprange_object1 = module.IPRangeObject{
        Metadata: module.ObjectMetaData{"ipr1", "", "", ""}, 
        Specification: module.IPRangeObjectSpec{"192.168.0.2", 1, 15}}
    var cert_object1 = module.CertificateObject{
        Metadata: module.ObjectMetaData{"device-a", "", "", ""}}
    var hubdevice_object = module.HubDeviceObject{
	    Metadata: module.ObjectMetaData{"hubdeviceconn", "", "", ""}}

    createControllerObject(OverlayUrl, &object1, &module.OverlayObject{})
    createControllerObject(ProposalUrl, &objecta, &module.ProposalObject{})
    createControllerObject(ProposalUrl, &objectb, &module.ProposalObject{})
    createControllerObject(IprangeUrl, &iprange_object1, &module.IPRangeObject{})
    createControllerObject(CertUrl, &cert_object1, &module.CertificateObject{})
    log.Println("Preparation ready! ")
    time.Sleep(10)

    if flag {
    createControllerObject(HubUrl, &hub, &module.HubObject{})
    log.Println("Register hub ready! ")
    time.Sleep(10)
    createControllerObject(DeviceUrl, &device, &module.DeviceObject{})
    log.Println("Register device ready! ")
    time.Sleep(10)
    }

    updateControllerObject(BaseUrl, "device-a", &hubdevice_object, &module.HubDeviceObject{})
    log.Println("Register hubdeviceconn ready! ")
    time.Sleep(10)

    var ret = m.Run()

    deleteControllerObject(BaseUrl, "device-a")
    deleteControllerObject(DeviceUrl, "device-a")
    deleteControllerObject(HubUrl, "huba")
    deleteControllerObject(IprangeUrl, "ipr1")
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

    var objs []module.HubDeviceObject
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
        _, err := getControllerObject(BaseUrl, tcase.object_name, &module.HubDeviceObject{})
        handleError(t, err, tcase.name, tcase.expectedErr, tcase.expectedErrCode)
    }
}

func TestCreateObject(t *testing.T) {
    tcases := []struct {
        name string
        obj module.HubDeviceObject
        expectedErr bool
        expectedErrCode int
    }{
        {
            name: "EmptyName",
            obj: module.HubDeviceObject{
                Metadata: module.ObjectMetaData{"", "object 1", "", ""}},
            expectedErr: true,
            expectedErrCode: 422,
        },
    }

    for _, tcase := range tcases {
        _, err := createControllerObject(BaseUrl, &tcase.obj, &module.HubDeviceObject{})
        handleError(t, err, tcase.name, tcase.expectedErr, tcase.expectedErrCode)
    }
}

func TestCreateObjectPass(t *testing.T) {

    tcases := []struct {
        name string
        obj module.HubDeviceObject
        expectedErr bool
        expectedErrCode int
    }{
        {
            name: "Normal",
            obj: module.HubDeviceObject{
                Metadata: module.ObjectMetaData{"hubdevicetest", "object 2", "", ""}},
        },
    }

    for _, tcase := range tcases {
        _, err := createControllerObject(BaseUrl, &tcase.obj, &module.HubDeviceObject{})
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
            object_name: "hubdevicetest",
        },
    }

    for _, gcase := range gcases {
        _, err := getControllerObject(BaseUrl, gcase.object_name, &module.HubDeviceObject{})
        handleError(t, err, gcase.name, gcase.expectedErr, gcase.expectedErrCode)
    }
}

