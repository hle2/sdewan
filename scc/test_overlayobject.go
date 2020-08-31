package main

import (
    "flag"
    "os"
    "io/ioutil"
    "net/http"
    "bytes"
    "encoding/json"
    "github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/module"
    "github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/manager"
)

var BaseUrl string

func callRest(method string, url string, request string) (string, error) {
    client := &http.Client{}
    req_body := bytes.NewBuffer([]byte(request))
    req, _ := http.NewRequest(method, url, req_body)

    req.Header.Set("Cache-Control", "no-cache")
    
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)
    if resp.StatusCode >= 400 {
        return "", nil 
    }

    return string(body), nil
}

func createControllerObject(baseUrl string, obj module.ControllerObject, retObj module.ControllerObject) (module.ControllerObject, error) {
    url := baseUrl
    obj_str, _ := json.Marshal(obj)

    res, err := callRest("POST", url, string(obj_str))
    if err != nil {
        return retObj, err
    }

    err = json.Unmarshal([]byte(res), retObj)
    if err != nil {
        return retObj, err
    }

    return retObj, nil
}

func getControllerObject(baseUrl string, name string, retObj module.ControllerObject) (module.ControllerObject, error) {
    url := baseUrl + "/" + name

    res, err := callRest("GET", url, "")
    if err != nil {
         return retObj, err
    }

    err = json.Unmarshal([]byte(res), retObj)
    if err != nil {
        return retObj, err
    }

    return retObj, nil
}

func updateControllerObject(baseUrl string, name string, obj module.ControllerObject, retObj module.ControllerObject) (module.ControllerObject, error) {
    url := baseUrl + "/" + name
    obj_str, _ := json.Marshal(obj)

    res, err := callRest("PUT", url, string(obj_str))
    if err != nil {
        return retObj, err
    }

    err = json.Unmarshal([]byte(res), retObj)
    if err != nil {
        return retObj, err
    }

    return retObj, nil
}

func deleteControllerObject(baseUrl string, name string) (bool, error) {
    url := baseUrl + "/" + name

    _, err := callRest("DELETE", url, "")
    if err != nil {
        return false, err
    }

    return true, nil
}

func main() {
    servIp := flag.String("ip", "127.0.0.1", "SDEWAN Central Controller IP Address")
    flag.Parse()
    BaseUrl = "http://" + *servIp + ":9015/scc/v1/" + manager.OverlayCollection

    var object1 = module.OverlayObject{
        Metadata: module.ObjectMetaData{"overlay1", "", "", ""}, 
        Specification: module.OverlayObjectSpec{"caid1"}}

    createControllerObject(BaseUrl, &object1, &module.OverlayObject{})

    var ret = 1 // m.Run()

    deleteControllerObject(BaseUrl, "overlay1")

    os.Exit(ret)
}
