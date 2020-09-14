/*
 * Copyright 2020 Intel Corporation, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governinog permissions and
 * limitations under the License.
 */

package manager

import (
    "log"
	"io"
	"encoding/json"
	"github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/infra/db"
	"github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/module"
	pkgerrors "github.com/pkg/errors"
)

type ProposalObjectKey struct {
	OverlayName	string `json:"overlay-name"`
    ProposalName string `json:"proposal-name"`
}

// ProposalObjectManager implements the ControllerObjectManager
type ProposalObjectManager struct {
    BaseObjectManager
}

func NewProposalObjectManager() *ProposalObjectManager {
	return &ProposalObjectManager{
        BaseObjectManager {
            storeName:  StoreName,
            tagMeta:    "proposal",
            depResManagers: []ControllerObjectManager {},
            ownResManagers: []ControllerObjectManager {},
        },
	}
}

func (c *ProposalObjectManager) IsOperationSupported(oper string) bool {
	return true
}

func (c *ProposalObjectManager) CreateEmptyObject() module.ControllerObject {
    return &module.ProposalObject{}
}

func (c *ProposalObjectManager) GetStoreKey(m map[string]string, t module.ControllerObject, isCollection bool) (db.Key, error) {
    overlay_name := m[OverlayResource]
    key := ProposalObjectKey{
        OverlayName: overlay_name,
        ProposalName: "",
    }

    if isCollection == true {
        return key, nil
    }

    to := t.(*module.ProposalObject)
    meta_name := to.Metadata.Name
    res_name := m[ProposalResource]

    if res_name != "" {
        if meta_name != "" && res_name != meta_name {
            return key, pkgerrors.New("Resource name unmatched metadata name")
        } 

        key.ProposalName = res_name
    } else {
        if meta_name == "" {
            return key, pkgerrors.New("Unable to find resource name")  
        }

        key.ProposalName = meta_name
    }

    return key, nil;
}

func (c *ProposalObjectManager) ParseObject(r io.Reader) (module.ControllerObject, error) {
	var v module.ProposalObject
    err := json.NewDecoder(r).Decode(&v)

    return &v, err
}

func (c *ProposalObjectManager) CreateObject(m map[string]string, t module.ControllerObject) (module.ControllerObject, error) {
	// for certificate test
    overlay := GetManagerset().Overlay
    to := t.(*module.ProposalObject)
    pname := to.Metadata.Name
    oname := m[OverlayResource]
    log.Println("Create Certificate: " + pname + "-cert")
    crt, key, _:= overlay.CreateCertificate(oname, pname + "-cert")
    log.Println("Crt: \n" + crt)
    log.Println("Key: \n" + key)

    // for ip range test
    iprange := GetManagerset().IPRange
    var nip string
    var ip [5]string
    var err error
    for i:=0; i<5; i++ {
        ip[i], err = iprange.Allocate(oname, "Dev1")
        if err != nil {
            log.Println(err)
        } else {
            log.Println("Allocated IP: " + ip[i])
        }
    }
    
    log.Println("Free IP: " + ip[2])
    iprange.Free(oname, ip[2])

    nip, err = iprange.Allocate(oname, "Dev1")
    if err != nil {
        log.Println(err)
    } else {
        log.Println("Allocated IP: " + nip)
    }

    // DB Operation
    t, err = GetDBUtils().CreateObject(c, m, t)

    return t, err
}

func (c *ProposalObjectManager) GetObject(m map[string]string) (module.ControllerObject, error) {
	// DB Operation
    t, err := GetDBUtils().GetObject(c, m)

    return t, err
}

func (c *ProposalObjectManager) GetObjects(m map[string]string) ([]module.ControllerObject, error) {
    // DB Operation
    t, err := GetDBUtils().GetObjects(c, m)

    return t, err
}

func (c *ProposalObjectManager) UpdateObject(m map[string]string, t module.ControllerObject) (module.ControllerObject, error) {
	// DB Operation
    t, err := GetDBUtils().UpdateObject(c, m, t)

    return t, err
}

func (c *ProposalObjectManager) DeleteObject(m map[string]string) error {
	// for certificate test
    overlay := GetManagerset().Overlay
    pname := m[ProposalResource]
    log.Println("Delete Certificate: " + pname + "-cert")
    overlay.DeleteCertificate(pname + "-cert")

    // for ip range test
    iprange := GetManagerset().IPRange
    oname := m[OverlayResource]
    iprange.FreeAll(oname)

    // DB Operation
    err := GetDBUtils().DeleteObject(c, m)

    return err
}
