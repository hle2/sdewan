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
    "github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/infra/validation"
    "github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/module"
    "github.com/go-playground/validator/v10"
    pkgerrors "github.com/pkg/errors"
)

type TestObjectKey struct {
    Name string `json:"name"`
}

// TestObjectManager implements the ControllerObjectManager
type TestObjectManager struct {
    storeName       string
    tagMeta         string
}

func NewTestObjectManager() *TestObjectManager {
    object_name := "testobject"
    validate := validation.GetValidator(object_name)
    validate.RegisterValidation("check_field", ValidateMyField)
    validate.RegisterStructValidation(ValidateMyObject, module.TestObject{})

    return &TestObjectManager{
        storeName:  object_name,
        tagMeta:    "testmetadata",
    }
}

func ValidateMyField(fl validator.FieldLevel) bool {
    return fl.Field().String() == "myfield" || fl.Field().String() == "myfield2"
}

func ValidateMyObject(sl validator.StructLevel) {
    obj := sl.Current().Interface().(module.TestObject)
    if obj.Specification.MyField == "myfield" {
        sl.ReportError(obj.Specification.MyField, "myfield", "myfield", "ValidateMyObject", "")
    }
}


func (c *TestObjectManager) IsOperationSupported(oper string) bool {
    return true
}

func (c *TestObjectManager) GetName() string {
    return c.storeName
}

func (c *TestObjectManager) ParseObject(r io.Reader) (module.ControllerObject, error) {
    var v module.TestObject
    err := json.NewDecoder(r).Decode(&v)

    return &v, err
}

func (c *TestObjectManager) CreateObject(m map[string]string, t module.ControllerObject) (module.ControllerObject, error) {
    to := t.(*module.TestObject)
    key := TestObjectKey{
        Name: to.Metadata.Name,
    }

    err := db.DBconn.Insert(c.storeName, key, nil, c.tagMeta, to)
    if err != nil {
        return &module.TestObject{}, pkgerrors.New("Unable to create the test object")
    }
    return t, nil
}

func (c *TestObjectManager) GetObject(m map[string]string) (module.ControllerObject, error) {
    name := m["test-name"]
    if name == "" {
        return &module.TestObject{}, pkgerrors.New("Missing test-name in GET request")
    }

    key := TestObjectKey{
        Name: name,
    }

    value, err := db.DBconn.Find(c.storeName, key, c.tagMeta)
    if err != nil {
        return &module.TestObject{}, pkgerrors.Wrap(err, "Get Resource")
    }
    if value != nil {
        t := module.TestObject{}
        err = db.DBconn.Unmarshal(value[0], &t)
        if err != nil {
            return &module.TestObject{}, pkgerrors.Wrap(err, "Unmarshaling value")
        }
        return &t, nil
    }

    
    return &module.TestObject{}, pkgerrors.New("No Object")
}

func (c *TestObjectManager) GetObjects(m map[string]string) ([]module.ControllerObject, error) {
    key := TestObjectKey{
        Name: "",
    }

    values, err := db.DBconn.Find(c.storeName, key, c.tagMeta)
    if err != nil {
        return []module.ControllerObject{}, pkgerrors.Wrap(err, "Get Test Objects")
    }

    var resp []module.ControllerObject    
    for _, value := range values {
        t := module.TestObject{}
        err = db.DBconn.Unmarshal(value, &t)
        if err != nil {
            return []module.ControllerObject{}, pkgerrors.Wrap(err, "Unmarshaling values")
        }
        resp = append(resp, &t)
    }

    return resp, nil
}

func (c *TestObjectManager) UpdateObject(m map[string]string, t module.ControllerObject) (module.ControllerObject, error) {
    name := m["test-name"]
    if name == "" {
        return &module.TestObject{}, pkgerrors.New("Missing test-name in PUT request")
    }

    to := t.(*module.TestObject)
        key := TestObjectKey{
        Name: to.Metadata.Name,
    }

    err := db.DBconn.Insert(c.storeName, key, nil, c.tagMeta, to)
    if err != nil {
        return &module.TestObject{}, pkgerrors.Wrap(err, "Updating DB Entry")
    }
    return t, nil
}

func (c *TestObjectManager) DeleteObject(m map[string]string) error {
    name := m["test-name"]
    if name == "" {
        return pkgerrors.New("Missing test-name in DELETE request")
    }

    key := TestObjectKey{
        Name: name,
    }

    err := db.DBconn.Remove(c.storeName, key)
    if err != nil {
        return pkgerrors.Wrap(err, "Delete Logical Cloud")
    }

    return nil
}
