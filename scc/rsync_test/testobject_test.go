package test

import (
//    "reflect"
    "testing"
    "flag"
    "encoding/json"
    "fmt"
    "os"
    "github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/module"
)

var BaseUrl string

func createObject(obj module.TestObject) (bool, error) {
    url := BaseUrl
    obj_str, _ := json.Marshal(obj)

    _, err := callRest("POST", url, string(obj_str))
    if err != nil {
        return false, err
    }

    return true, nil
}

func getObject(name string) (module.TestObject, error) {
    url := BaseUrl + "/" + name

    res, err := callRest("GET", url, "")
    if err != nil {
         return module.TestObject{}, err
    }

    var obj module.TestObject
    err = json.Unmarshal([]byte(res), &obj)
    if err != nil {
        return module.TestObject{}, err
    }

    return obj, nil
}

func updateObject(name string, obj module.TestObject) (module.TestObject, error) {
    url := BaseUrl + "/" + name
    obj_str, _ := json.Marshal(obj)

    res, err := callRest("PUT", url, string(obj_str))
    if err != nil {
        return module.TestObject{}, err
    }

    var ret_obj module.TestObject
    err = json.Unmarshal([]byte(res), &ret_obj)
    if err != nil {
        return module.TestObject{}, err
    }

    return ret_obj, nil
}

func deleteObject(name string) bool {
    url := BaseUrl + "/" + name

    _, err := callRest("DELETE", url, "")
    if err != nil {
        return false
    }

    return true
}

func TestMain(m *testing.M) {
    servIp := flag.String("ip", "127.0.0.1", "SDEWAN Central Hub Controller IP Address")
    flag.Parse()
    BaseUrl = "http://" + *servIp + ":9015/v1/tests"

    var object1 = module.TestObject{
        Metadata: module.ObjectMetaData{"obj1", "object 1", "", ""}, 
        Specification: module.TestObjectSpec{"192.168.1.1", 3000, "myfield2"}}
    var object2 = module.TestObject{
        Metadata: module.ObjectMetaData{"obj2", "object 2", "", ""}, 
        Specification: module.TestObjectSpec{"192.168.1.2", 3001, "myfield2"}}

    createObject(object1)
    createObject(object2)

    var ret = m.Run()

    deleteObject("obj1")
    deleteObject("obj2")

    os.Exit(ret)
}

func TestGetObjects(t *testing.T) {
    url := BaseUrl
    res, err := callRest("GET", url, "")
    if err != nil {
        printError(err)
        t.Errorf("Test case GetObjects: can not get Test Objects")
        return
    }

    var objs []module.TestObject
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
            object_name: "obj1",
        },
        {
            name: "GetFoolName",
            object_name: "foo_name",
            expectedErr: true,
            expectedErrCode: 500,
        },
    }

    for _, tcase := range tcases {
        _, err := getObject(tcase.object_name)
        handleError(t, err, tcase.name, tcase.expectedErr, tcase.expectedErrCode)
    }
}

func TestCreateObject(t *testing.T) {
    tcases := []struct {
        name string
        test_object module.TestObject
        expectedErr bool
        expectedErrCode int
    }{
        {
            name: "InvalidIP",
            test_object: module.TestObject{
                Metadata: module.ObjectMetaData{"obj1", "object 1", "", ""}, 
                Specification: module.TestObjectSpec{"192.168.1.1.1", 3000, "myfield2"}},
            expectedErr: true,
            expectedErrCode: 422,
        },
        {
            name: "InvalidPort",
            test_object: module.TestObject{
                Metadata: module.ObjectMetaData{"obj1", "object 1", "", ""}, 
                Specification: module.TestObjectSpec{"192.168.1.1", 90000, "myfield2"}},
            expectedErr: true,
            expectedErrCode: 422,
        },
        {
            name: "InvalidField",
            test_object: module.TestObject{
                Metadata: module.ObjectMetaData{"obj1", "object 1", "", ""}, 
                Specification: module.TestObjectSpec{"192.168.1.1", 9000, "myfield3"}},
            expectedErr: true,
            expectedErrCode: 422,
        },
    }

    for _, tcase := range tcases {
        _, err := createObject(tcase.test_object)
        handleError(t, err, tcase.name, tcase.expectedErr, tcase.expectedErrCode)
    }
}