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
var HubUrl string

func TestMain(m *testing.M) {
    servIp := flag.String("ip", "127.0.0.1", "SDEWAN Central Controller IP Address")
    flag.Parse()
    OverlayUrl := "http://" + *servIp + ":9015/scc/v1/" + manager.OverlayCollection
    ProposalUrl := OverlayUrl + "/overlay1/" + manager.ProposalCollection
    DeviceUrl := OverlayUrl + "/overlay1/" + manager.DeviceCollection
    HubUrl = OverlayUrl + "/overlay1/" + manager.HubCollection
    IprangeUrl := OverlayUrl + "/overlay1/" + manager.IPRangeCollection
    CertUrl := OverlayUrl + "/overlay1/" + manager.CertCollection
    BaseUrl := OverlayUrl + "/overlay1/" + manager.HubCollection + "/huba/" + manager.DeviceCollection
    flag := true

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
    var hubdevice_object = module.HubDeviceObject{
	Metadata: module.ObjectMetaData{"hubdeviceconn", "", "", ""},
	Specification: module.HubDeviceObjectSpec{"device-a"}}
    var cert_object1 = module.CertificateObject{
        Metadata: module.ObjectMetaData{"device-a", "", "", ""}}


    if flag {
    createControllerObject(OverlayUrl, &object1, &module.OverlayObject{})
    createControllerObject(ProposalUrl, &objecta, &module.ProposalObject{})
    createControllerObject(ProposalUrl, &objectb, &module.ProposalObject{})
    createControllerObject(IprangeUrl, &iprange_object1, &module.IPRangeObject{})
    createControllerObject(CertUrl, &cert_object1, &module.CertificateObject{})
    log.Println("Preparation ready! ")
    time.Sleep(10)
    createControllerObject(HubUrl, &hub, &module.HubObject{})
    log.Println("Register hub ready! ")
    time.Sleep(10)
    createControllerObject(DeviceUrl, &device, &module.DeviceObject{})
    log.Println("Register device ready! ")
    time.Sleep(60 * time.Second)
    createControllerObject(BaseUrl, &hubdevice_object, &module.HubDeviceObject{})
    log.Println("Register hubdeviceconn ready! ")
    }


    var ret = m.Run()

    deleteControllerObject(BaseUrl, "device-a")
    deleteControllerObject(DeviceUrl, "device-a")
    deleteControllerObject(HubUrl, "huba")
    deleteControllerObject(IprangeUrl, "ipr1")
    deleteControllerObject(ProposalUrl, "proposal2")
    deleteControllerObject(ProposalUrl, "proposal1")
    //deleteControllerObject(OverlayUrl, "overlay1")

    os.Exit(ret)
}

func TestGetObjects(t *testing.T) {
    url := HubUrl + "/huba/connections" 
    res, err := callRest("GET", url, "")
    if err != nil {
        printError(err)
        t.Errorf("Test case GetObjects: can not get Objects")
        return
    }

    var objs []module.ConnectionObject
    err = json.Unmarshal([]byte(res), &objs)

    if len(objs) == 0 {
        fmt.Printf("Test case GetObjects: no object found")
        return
    }

    p_data, _ := json.Marshal(objs)
    fmt.Printf("%s\n", string(p_data))

}

