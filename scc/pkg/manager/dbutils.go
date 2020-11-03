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
    "encoding/base64"
    "github.com/onap/multicloud-k8s/src/orchestrator/pkg/infra/db"
    "github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/module"
    pkgerrors "github.com/pkg/errors"
    mtypes "github.com/onap/multicloud-k8s/src/orchestrator/pkg/module/types"
)

const PROVIDERNAME = "akraino_scc"

type Cluster struct {
    Metadata mtypes.Metadata `json:"metadata"`
}

type ClusterContent struct {
    Kubeconfig string `json:"kubeconfig"`
}

type ClusterKey struct {
    ClusterProviderName string `json:"provider"`
    ClusterName         string `json:"cluster"`
}

type DBUtils struct {
}

var dbutils = DBUtils{}

func GetDBUtils() *DBUtils {
    return &dbutils
}

func (d *DBUtils) CheckDep(c ControllerObjectManager, m map[string]string) error {
    depsResManagers := c.GetDepResManagers()
    for _, mgr := range depsResManagers {
        _, err := d.GetObject(mgr, m)
        if err != nil {
            return pkgerrors.New("Fail to find " + mgr.GetStoreMeta())
        }
    }

    return nil
}

func (d *DBUtils) CreateObject(c ControllerObjectManager, 
    m map[string]string, 
    t module.ControllerObject) (module.ControllerObject, error) {
 
    key, _ := c.GetStoreKey(m, t, false)
    err := db.DBconn.Insert(c.GetStoreName(), key, nil, c.GetStoreMeta(), t)
    if err != nil {
        return c.CreateEmptyObject(), pkgerrors.New("Unable to create the object")
    }

    return t, nil
}

func (d *DBUtils) GetObject(c ControllerObjectManager, 
    m map[string]string) (module.ControllerObject, error) {

    key, err := c.GetStoreKey(m, c.CreateEmptyObject(), false)
    if err != nil {
        return c.CreateEmptyObject(), err
    }

    
    value, err := db.DBconn.Find(c.GetStoreName(), key, c.GetStoreMeta())
    if err != nil {
        return c.CreateEmptyObject(), pkgerrors.Wrap(err, "Get Resource")
    }

    
    if value != nil {
        r := c.CreateEmptyObject()
        err = db.DBconn.Unmarshal(value[0], r)
        if err != nil {
            return c.CreateEmptyObject(), pkgerrors.Wrap(err, "Unmarshaling value")
        }
        return r, nil
    }

    return c.CreateEmptyObject(), pkgerrors.New("No Object")
}

func (d *DBUtils) GetObjects(c ControllerObjectManager,
    m map[string]string) ([]module.ControllerObject, error) {
    
        
    key, err := c.GetStoreKey(m, c.CreateEmptyObject(), true)
    if err != nil {
        return []module.ControllerObject{}, err
    }

    
    values, err := db.DBconn.Find(c.GetStoreName(), key, c.GetStoreMeta())
    if err != nil {
        return []module.ControllerObject{}, pkgerrors.Wrap(err, "Get Overlay Objects")
    }

    
    var resp []module.ControllerObject
    for _, value := range values {
        t := c.CreateEmptyObject()
        err = db.DBconn.Unmarshal(value, t)
        if err != nil {
            return []module.ControllerObject{}, pkgerrors.Wrap(err, "Unmarshaling values")
        }
        resp = append(resp, t)
    }

    return resp, nil
}

func (d *DBUtils) UpdateObject(c ControllerObjectManager,
    m map[string]string, t module.ControllerObject) (module.ControllerObject, error) {

    key, err := c.GetStoreKey(m, t, false)
    if err != nil {
        return c.CreateEmptyObject(), err
    }

    err = db.DBconn.Insert(c.GetStoreName(), key, nil, c.GetStoreMeta(), t)
    if err != nil {
        return c.CreateEmptyObject(), pkgerrors.Wrap(err, "Updating DB Entry")
    }
    return t, nil
}

func (d *DBUtils) DeleteObject(c ControllerObjectManager, m map[string]string) error {
    key, err := c.GetStoreKey(m, c.CreateEmptyObject(), false)
    if err != nil {
        return err
    }

    err = db.DBconn.Remove(c.GetStoreName(), key)
    if err != nil {
        return pkgerrors.Wrap(err, "Delete Object")
    }

    return nil
}

func (d *DBUtils) RegisterDevice(cluster_name string, kubeconfig_file string) error {
    var q ClusterContent
    var p Cluster

    //content, err := ioutil.ReadFile(kubeconfig_file)
    q.Kubeconfig = base64.StdEncoding.EncodeToString(kubeconfig_file)
    key := ClusterKey{
        ClusterProviderName: PROVIDERNAME,
        ClusterName:         cluster_name,
    }

    p.Metadata.Name = cluster_name

    err = db.DBconn.Insert("cluster", key, nil, "clustermetadata", p)
    if err != nil {
        return pkgerrors.Wrap(err, "Creating DB Entry")
    }

    err = db.DBconn.Insert("cluster", key, nil, "clustercontent", q)
    if err != nil {
        return pkgerrors.Wrap(err, "Creating DB Entry")
    }

    return nil
}
