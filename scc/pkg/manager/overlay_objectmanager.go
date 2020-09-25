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
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package manager

import (
    "io"
    "log"
    "encoding/json"
    "github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/infra/db"
    "github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/module"
    "github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/resource"
    pkgerrors "github.com/pkg/errors"
)

type OverlayObjectKey struct {
    OverlayName string `json:"overlay-name"`
}

// OverlayObjectManager implements the ControllerObjectManager
type OverlayObjectManager struct {
    BaseObjectManager
}

func NewOverlayObjectManager() *OverlayObjectManager {
    return &OverlayObjectManager{
        BaseObjectManager {
            storeName:  StoreName,
            tagMeta:    "overlay",
            depResManagers: []ControllerObjectManager {},
            ownResManagers: []ControllerObjectManager {},
        },
    }
}

func (c *OverlayObjectManager) IsOperationSupported(oper string) bool {
    return true
}

func (c *OverlayObjectManager) CreateEmptyObject() module.ControllerObject {
    return &module.OverlayObject{}
}

func (c *OverlayObjectManager) GetStoreKey(m map[string]string, t module.ControllerObject, isCollection bool) (db.Key, error) {
    key := OverlayObjectKey{""}

    if isCollection == true {
        return key, nil
    }

    to := t.(*module.OverlayObject)
    meta_name := to.Metadata.Name
    res_name := m[OverlayResource]

    if res_name != "" {
        if meta_name != "" && res_name != meta_name {
            return key, pkgerrors.New("Resource name unmatched metadata name")
        } 

        key.OverlayName = res_name
    } else {
        if meta_name == "" {
            return key, pkgerrors.New("Unable to find resource name")  
        }

        key.OverlayName = meta_name
    }

    return key, nil;
}

func (c *OverlayObjectManager) ParseObject(r io.Reader) (module.ControllerObject, error) {
    var v module.OverlayObject
    err := json.NewDecoder(r).Decode(&v)

    return &v, err
}

func (c *OverlayObjectManager) CreateObject(m map[string]string, t module.ControllerObject) (module.ControllerObject, error) {
    // for rsync test
     resutil := NewResUtil()

     deviceObject := module.OverlayObject{
                Metadata: module.ObjectMetaData{"local", "", "", ""}, 
                Specification: module.OverlayObjectSpec{}}
     resutil.AddResource(&deviceObject, "Create", &resource.FileResource{"mycm", "ConfigMap", "mycm.yaml"})

     err2 := resutil.Deploy("YAML")

     if err2 != nil {
         return c.CreateEmptyObject(), pkgerrors.Wrap(err2, "Unable to create the object: fail to deploy resource")
     }

    // Create a issuer each overlay
    to := t.(*module.OverlayObject)
    overlay_name := to.Metadata.Name
    cu, err := GetCertUtil()
    if err != nil {
        log.Println(err)
    } else {
        // create overlay ca
        _, err := cu.CreateCertificate(c.CertName(overlay_name), NameSpaceName, RootCAIssuerName, true)
        if err == nil {
            // create overlay issuer
            _, err := cu.CreateCAIssuer(c.IssuerName(overlay_name), NameSpaceName, c.CertName(overlay_name))
            if err != nil {
                log.Println("Failed to create overlay[" + overlay_name +"] issuer: " + err.Error())
            }    
        } else {
            log.Println("Failed to create overlay[" + overlay_name +"] certificate: " + err.Error())
        }
    }

    // DB Operation
    t, err = GetDBUtils().CreateObject(c, m, t)

    return t, err
}

func (c *OverlayObjectManager) GetObject(m map[string]string) (module.ControllerObject, error) {
    // DB Operation
    t, err := GetDBUtils().GetObject(c, m)

    return t, err
}

func (c *OverlayObjectManager) GetObjects(m map[string]string) ([]module.ControllerObject, error) {
    // DB Operation
    t, err := GetDBUtils().GetObjects(c, m)

    return t, err
}

func (c *OverlayObjectManager) UpdateObject(m map[string]string, t module.ControllerObject) (module.ControllerObject, error) {
    // DB Operation
    t, err := GetDBUtils().UpdateObject(c, m, t)

    return t, err
}

func (c *OverlayObjectManager) DeleteObject(m map[string]string) error {
    overlay_name := m[OverlayResource]
    cu, err := GetCertUtil()
    if err != nil {
        log.Println(err)
    } else {
        err = cu.DeleteIssuer(c.IssuerName(overlay_name), NameSpaceName)
        if err != nil {
            log.Println("Failed to delete overlay[" + overlay_name +"] issuer: " + err.Error())
        }
        err = cu.DeleteCertificate(c.CertName(overlay_name), NameSpaceName)
        if err != nil {
            log.Println("Failed to delete overlay[" + overlay_name +"] certificate: " + err.Error())
        }
    }

    // DB Operation
    err = GetDBUtils().DeleteObject(c, m)

    return err
}

func (c *OverlayObjectManager) IssuerName(name string) string {
    return name + "-issuer"
}

func (c *OverlayObjectManager) CertName(name string) string {
    return name + "-cert"
}

func (c *OverlayObjectManager) CreateCertificate(oname string, cname string) (string, string, error) {
    cu, err := GetCertUtil()
    if err != nil {
        log.Println(err)
    } else {
        _, err := cu.CreateCertificate(cname, NameSpaceName, c.IssuerName(oname), false)
        if err != nil {
            log.Println("Failed to create overlay[" + oname +"] certificate: " + err.Error())
        } else {
            return cu.GetKeypair(cname, NameSpaceName)
        }
    }

    return "", "", nil
}

func (c *OverlayObjectManager) DeleteCertificate(cname string) (string, string, error) {
    cu, err := GetCertUtil()
    if err != nil {
        log.Println(err)
    } else {
        err = cu.DeleteCertificate(cname, NameSpaceName)
        if err != nil {
            log.Println("Failed to delete " + cname +" certificate: " + err.Error())
        }
    }

    return "", "", nil
}