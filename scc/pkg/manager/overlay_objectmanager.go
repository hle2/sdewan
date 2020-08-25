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
	"io"
	"encoding/json"
	"github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/infra/db"
	"github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/module"
    "github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/resource"
	pkgerrors "github.com/pkg/errors"
)

type OverlayObjectKey struct {
	OverlayName	string `json:"overlay-name"`
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
	// DB Operation
    err :=  GetDBUtils().checkDep(c, m)
    if err != nil {
        return c.CreateEmptyObject(), pkgerrors.Wrap(err, "Unable to create the object")
    }

    // for test
    resutil := NewResUtil()

    deviceObject := module.OverlayObject{
                Metadata: module.ObjectMetaData{"local", "", "", ""}, 
                Specification: module.OverlayObjectSpec{"caid1"}}
    resutil.AddResource(&deviceObject, "Create", &resource.FileResource{"mycm", "ConfigMap", "mycm.yaml"})

    err = resutil.Deploy("YAML")

    if err != nil {
        return c.CreateEmptyObject(), pkgerrors.Wrap(err, "Unable to create the object: fail to deploy resource")
    }

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
	// DB Operation
    err := GetDBUtils().DeleteObject(c, m)

    return err
}
