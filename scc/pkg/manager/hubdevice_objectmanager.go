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

type HubDeviceObjectKey struct {
    OverlayName string `json:"overlay-name"`
    HubName string `json:"hub-name"`
    DeviceName string `json:"device-name"`
}

// HubDeviceObjectManager implements the ControllerObjectManager
type HubDeviceObjectManager struct {
    BaseObjectManager
}

func NewHubDeviceObjectManager() *HubDeviceObjectManager {
    return &HubDeviceObjectManager{
        BaseObjectManager {
            storeName:  StoreName,
            tagMeta:    "hubdevice",
            depResManagers: []ControllerObjectManager {},
            ownResManagers: []ControllerObjectManager {},
        },
    }
}

func (c *HubDeviceObjectManager) IsOperationSupported(oper string) bool {
    if oper == "POST" || oper == "DELETE" {
        return true
    }
    return false
}

func (c *HubDeviceObjectManager) CreateEmptyObject() module.ControllerObject {
    return &module.HubDeviceObject{}
}

func (c *HubDeviceObjectManager) GetStoreKey(m map[string]string, t module.ControllerObject, isCollection bool) (db.Key, error) {
    overlay_name := m[OverlayResource]
    hub_name := m[HubResource]
    device_name := m[DeviceResource]
    key := HubDeviceObjectKey{
        OverlayName: overlay_name,
        HubName: hub_name,
        DeviceName: device_name,
    }

    return key, nil;
}

func (c *HubDeviceObjectManager) ParseObject(r io.Reader) (module.ControllerObject, error) {
    var v module.HubDeviceObject
    err := json.NewDecoder(r).Decode(&v)

    return &v, err
}

func (c *HubDeviceObjectManager) CreateObject(m map[string]string, t module.ControllerObject) (module.ControllerObject, error) {
    // Todo: setup hub-device connection
    return c.CreateEmptyObject(), pkgerrors.New("Not implemented")
}

func (c *HubDeviceObjectManager) GetObject(m map[string]string) (module.ControllerObject, error) {
    return c.CreateEmptyObject(), pkgerrors.New("Not implemented")
}

func (c *HubDeviceObjectManager) GetObjects(m map[string]string) ([]module.ControllerObject, error) {
    return []module.ControllerObject{}, pkgerrors.New("Not implemented")
}

func (c *HubDeviceObjectManager) UpdateObject(m map[string]string, t module.ControllerObject) (module.ControllerObject, error) {
    return c.CreateEmptyObject(), pkgerrors.New("Not implemented")
}

func (c *HubDeviceObjectManager) DeleteObject(m map[string]string) error {
    // Todo: delete hub-device connection
    return pkgerrors.New("Not implemented")
}
