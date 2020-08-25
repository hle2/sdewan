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
    "github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/module"
    "github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/resource"
//    pkgerrors "github.com/pkg/errors"
)

type DeployResource struct {
    Action string
    Resource resource.ISdewanResource
}

type DeployResources struct {
    Resources []DeployResource
}

type ResUtil struct {
    resmap map[module.ControllerObject]*DeployResources
}

func NewResUtil() *ResUtil {
    return &ResUtil{
        resmap: make(map[module.ControllerObject]*DeployResources),
    }
}

func (d *ResUtil) AddResource(device module.ControllerObject, action string, resource resource.ISdewanResource) error {
    if d.resmap[device] == nil {
        d.resmap[device] = &DeployResources{Resources: []DeployResource{}}
    }

    d.resmap[device].Resources = append(d.resmap[device].Resources, DeployResource{Action: action, Resource: resource,})
    return nil
}

func (d *ResUtil) Deploy(format string) error {
    return nil
}