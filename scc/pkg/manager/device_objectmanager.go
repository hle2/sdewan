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
    "encoding/json"
    "github.com/onap/multicloud-k8s/src/orchestrator/pkg/infra/db"
    "github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/module"
    pkgerrors "github.com/pkg/errors"
)

type DeviceObjectKey struct {
    OverlayName string `json:"overlay-name"`
    DeviceName string `json:"device-name"`
}

// DeviceObjectManager implements the ControllerObjectManager
type DeviceObjectManager struct {
    BaseObjectManager
}

func NewDeviceObjectManager() *DeviceObjectManager {
    return &DeviceObjectManager{
        BaseObjectManager {
            storeName:  StoreName,
            tagMeta:    "device",
            depResManagers: []ControllerObjectManager {},
            ownResManagers: []ControllerObjectManager {},
        },
    }
}

func (c *DeviceObjectManager) IsOperationSupported(oper string) bool {
    return true
}

func (c *DeviceObjectManager) CreateEmptyObject() module.ControllerObject {
    return &module.DeviceObject{}
}

func (c *DeviceObjectManager) GetStoreKey(m map[string]string, t module.ControllerObject, isCollection bool) (db.Key, error) {
    overlay_name := m[OverlayResource]
    key := DeviceObjectKey{
        OverlayName: overlay_name,
        DeviceName: "",
    }

    if isCollection == true {
        return key, nil
    }

    to := t.(*module.DeviceObject)
    meta_name := to.Metadata.Name
    res_name := m[DeviceResource]

    if res_name != "" {
        if meta_name != "" && res_name != meta_name {
            return key, pkgerrors.New("Resource name unmatched metadata name")
        } 

        key.DeviceName = res_name
    } else {
        if meta_name == "" {
            return key, pkgerrors.New("Unable to find resource name")  
        }

        key.DeviceName = meta_name
    }

    return key, nil;
}

func (c *DeviceObjectManager) ParseObject(r io.Reader) (module.ControllerObject, error) {
    var v module.DeviceObject
    err := json.NewDecoder(r).Decode(&v)

    return &v, err
}

func (c *DeviceObjectManager) CreateObject(m map[string]string, t module.ControllerObject) (module.ControllerObject, error) {
    // DB Operation
    t, err := GetDBUtils().CreateObject(c, m, t)

    return t, err
}

func (c *DeviceObjectManager) GetObject(m map[string]string) (module.ControllerObject, error) {
    // DB Operation
    t, err := GetDBUtils().GetObject(c, m)

    return t, err
}

func (c *DeviceObjectManager) GetObjects(m map[string]string) ([]module.ControllerObject, error) {
    // DB Operation
    t, err := GetDBUtils().GetObjects(c, m)

    return t, err
}

func (c *DeviceObjectManager) UpdateObject(m map[string]string, t module.ControllerObject) (module.ControllerObject, error) {
    // DB Operation
    t, err := GetDBUtils().UpdateObject(c, m, t)

    return t, err
}

func (c *DeviceObjectManager) DeleteObject(m map[string]string) error {
    // DB Operation
    err := GetDBUtils().DeleteObject(c, m)

    return err
}
