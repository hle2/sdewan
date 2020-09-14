package main

import (
    "log"
    "strconv"
    "github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/module"
)

func allocate(r *module.IPRangeObject, name string) bool {
    ip, err := r.Allocate(name)
    if err != nil {
        log.Println(err)
        return false
    } else {
        log.Println("Allocate: " + ip)
        return true
    }
}

func main() {
    r := module.IPRangeObject {
        Metadata: module.ObjectMetaData {"Range1", "", "", ""},
        Specification: module.IPRangeObjectSpec {"192.168.0.2", 20, 40},
        Status: module.IPRangeObjectStatus {
            Data: make(map[int]string),
        },
    }

    i := 1
    name := "dev" + strconv.Itoa(i)
    log.Println(r)
    for allocate(&r, name) {
        i += 1
        name = "dev" + strconv.Itoa(i)
    }
    log.Println(r)

    var err error
    fa := []int{17, 23, 40, 10, 50}
    for _, ip := range fa {
        err = r.Free("192.168.0." + strconv.Itoa(ip))
        if err != nil {
            log.Println(err)
        }
    }

    log.Println(r)
    for allocate(&r, name) {
        i += 1
        name = "dev" + strconv.Itoa(i)
    }
    log.Println(r)
}
