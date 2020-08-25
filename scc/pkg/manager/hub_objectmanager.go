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
	pkgerrors "github.com/pkg/errors"
)

type HubObjectKey struct {
	OverlayName	string `json:"overlay-name"`
    HubName string `json:"hub-name"`
}

// HubObjectManager implements the ControllerObjectManager
type HubObjectManager struct {
    BaseObjectManager
}

func NewHubObjectManager() *HubObjectManager {
	return &HubObjectManager{
        BaseObjectManager {
            storeName:  StoreName,
            tagMeta:    "hub",
            depResManagers: []ControllerObjectManager {},
            ownResManagers: []ControllerObjectManager {},
        },
	}
}

func (c *HubObjectManager) IsOperationSupported(oper string) bool {
	return true
}

func (c *HubObjectManager) CreateEmptyObject() module.ControllerObject {
    return &module.HubObject{}
}

func (c *HubObjectManager) GetStoreKey(m map[string]string, t module.ControllerObject, isCollection bool) (db.Key, error) {
    overlay_name := m[OverlayResource]
    key := HubObjectKey{
        OverlayName: overlay_name,
        HubName: "",
    }

    if isCollection == true {
        return key, nil
    }

    to := t.(*module.HubObject)
    meta_name := to.Metadata.Name
    res_name := m[HubResource]

    if res_name != "" {
        if meta_name != "" && res_name != meta_name {
            return key, pkgerrors.New("Resource name unmatched metadata name")
        } 

        key.HubName = res_name
    } else {
        if meta_name == "" {
            return key, pkgerrors.New("Unable to find resource name")  
        }

        key.HubName = meta_name
    }

    return key, nil;
}

func (c *HubObjectManager) ParseObject(r io.Reader) (module.ControllerObject, error) {
	var v module.HubObject
    err := json.NewDecoder(r).Decode(&v)

    return &v, err
}

func (c *HubObjectManager) CreateObject(m map[string]string, t module.ControllerObject) (module.ControllerObject, error) {
	// DB Operation
    err :=  GetDBUtils().checkDep(c, m)
    if err != nil {
        return c.CreateEmptyObject(), pkgerrors.Wrap(err, "Unable to create the object")
    }

    resutil := NewResUtil()

    // Todo: call resutil.AddResource() to add to-be-deployed resources

    err = resutil.Deploy("YAML")

    if err != nil {
        return c.CreateEmptyObject(), pkgerrors.Wrap(err, "Unable to create the object: fail to deploy resource")
    }

    t, err = GetDBUtils().CreateObject(c, m, t)

    return t, err
}

func (c *HubObjectManager) GetObject(m map[string]string) (module.ControllerObject, error) {
	// DB Operation
    t, err := GetDBUtils().GetObject(c, m)

    return t, err
}

func (c *HubObjectManager) GetObjects(m map[string]string) ([]module.ControllerObject, error) {
    // DB Operation
    t, err := GetDBUtils().GetObjects(c, m)

    return t, err
}

func (c *HubObjectManager) UpdateObject(m map[string]string, t module.ControllerObject) (module.ControllerObject, error) {
	// DB Operation
    t, err := GetDBUtils().UpdateObject(c, m, t)

    return t, err
}

func (c *HubObjectManager) DeleteObject(m map[string]string) error {
	// DB Operation
    err := GetDBUtils().DeleteObject(c, m)

    return err
}
