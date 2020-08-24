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

type IPRangeObjectKey struct {
	OverlayName	string `json:"overlay-name"`
    IPRangeName string `json:"iprange-name"`
}

// IPRangeObjectManager implements the ControllerObjectManager
type IPRangeObjectManager struct {
    BaseObjectManager
}

func NewIPRangeObjectManager() *IPRangeObjectManager {
	return &IPRangeObjectManager{
        BaseObjectManager {
            storeName:  StoreName,
            tagMeta:    "iprange",
            depResManagers: []ControllerObjectManager {},
            ownResManagers: []ControllerObjectManager {},
        },
	}
}

func (c *IPRangeObjectManager) IsOperationSupported(oper string) bool {
	return true
}

func (c *IPRangeObjectManager) CreateEmptyObject() module.ControllerObject {
    return &module.IPRangeObject{}
}

func (c *IPRangeObjectManager) GetStoreKey(m map[string]string, t module.ControllerObject, isCollection bool) (db.Key, error) {
    overlay_name := m[OverlayResource]
    key := IPRangeObjectKey{
        OverlayName: overlay_name,
        IPRangeName: "",
    }

    if isCollection == true {
        return key, nil
    }

    to := t.(*module.IPRangeObject)
    meta_name := to.Metadata.Name
    res_name := m[IPRangeResource]

    if res_name != "" {
        if meta_name != "" && res_name != meta_name {
            return key, pkgerrors.New("Resource name unmatched metadata name")
        } 

        key.IPRangeName = res_name
    } else {
        if meta_name == "" {
            return key, pkgerrors.New("Unable to find resource name")  
        }

        key.IPRangeName = meta_name
    }

    return key, nil;
}

func (c *IPRangeObjectManager) ParseObject(r io.Reader) (module.ControllerObject, error) {
	var v module.IPRangeObject
    err := json.NewDecoder(r).Decode(&v)

    return &v, err
}

func (c *IPRangeObjectManager) CreateObject(m map[string]string, t module.ControllerObject) (module.ControllerObject, error) {
	// DB Operation
    t, err := GetDBUtils().CreateObject(c, m, t)

    return t, err
}

func (c *IPRangeObjectManager) GetObject(m map[string]string) (module.ControllerObject, error) {
	// DB Operation
    t, err := GetDBUtils().GetObject(c, m)

    return t, err
}

func (c *IPRangeObjectManager) GetObjects(m map[string]string) ([]module.ControllerObject, error) {
    // DB Operation
    t, err := GetDBUtils().GetObjects(c, m)

    return t, err
}

func (c *IPRangeObjectManager) UpdateObject(m map[string]string, t module.ControllerObject) (module.ControllerObject, error) {
	// DB Operation
    t, err := GetDBUtils().UpdateObject(c, m, t)

    return t, err
}

func (c *IPRangeObjectManager) DeleteObject(m map[string]string) error {
	// DB Operation
    err := GetDBUtils().DeleteObject(c, m)

    return err
}
