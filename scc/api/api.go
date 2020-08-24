/*
Copyright 2020 Intel Corporation.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
    "github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/manager"
    "github.com/gorilla/mux"
)

// NewRouter creates a router that registers the various urls that are
// supported

func createHandlerMapping(
    objectClient manager.ControllerObjectManager,
    router *mux.Router,
    collections string,
    resource string ) {
    objectHandler := ControllerHandler{client: objectClient}
    if objectClient.IsOperationSupported("POST") == true {
        router.HandleFunc(
            "/" + collections,
            objectHandler.createHandler).Methods("POST")
    }

    if objectClient.IsOperationSupported("GETS") == true {
        router.HandleFunc(
            "/" + collections,
            objectHandler.getsHandler).Methods("GET")
    }

    if objectClient.IsOperationSupported("GET") == true {
        router.HandleFunc(
            "/" + collections + "/{" + resource + "}",
            objectHandler.getHandler).Methods("GET")
    }

    if objectClient.IsOperationSupported("DELETE") == true {
        router.HandleFunc(
            "/" + collections + "/{" + resource + "}",
            objectHandler.deleteHandler).Methods("DELETE")
    }

    if objectClient.IsOperationSupported("PUT") == true {
        router.HandleFunc(
            "/" + collections + "/{" + resource + "}",
            objectHandler.updateHandler).Methods("PUT")
    }
}

func NewRouter(
    overlayObjectClient manager.ControllerObjectManager,
    proposalObjectClient manager.ControllerObjectManager,
    hubObjectClient manager.ControllerObjectManager,
    hubConnObjectClient manager.ControllerObjectManager,
    hubDeviceObjectClient manager.ControllerObjectManager,
    deviceObjectClient manager.ControllerObjectManager,
    deviceConnObjectClient manager.ControllerObjectManager,
    ipRangeObjectClient manager.ControllerObjectManager) *mux.Router {

    router := mux.NewRouter()
    ver := "v1"

    // router
    verRouter := router.PathPrefix("/scc/" + ver).Subrouter()
    olRouter := verRouter.PathPrefix("/" + manager.OverlayCollection + "/{" + manager.OverlayResource + "}").Subrouter()
    hubRouter := olRouter.PathPrefix("/" + manager.HubCollection + "/{" + manager.HubResource + "}").Subrouter()
    devRouter := olRouter.PathPrefix("/" + manager.DeviceCollection + "/{" + manager.DeviceResource + "}").Subrouter()

    // overlay API
    if overlayObjectClient == nil {
         overlayObjectClient = manager.NewOverlayObjectManager()
    }
    createHandlerMapping(overlayObjectClient, verRouter, manager.OverlayCollection, manager.OverlayResource)

    // proposal API
    if proposalObjectClient == nil {
         proposalObjectClient = manager.NewProposalObjectManager()
    }
    createHandlerMapping(proposalObjectClient, olRouter, manager.ProposalCollection, manager.ProposalResource)

    // hub API
    if hubObjectClient == nil {
         hubObjectClient = manager.NewHubObjectManager()
    }
    createHandlerMapping(hubObjectClient, olRouter, manager.HubCollection, manager.HubResource)

    // hub-connection API
    if hubConnObjectClient == nil {
         hubConnObjectClient = manager.NewHubConnObjectManager()
    }
    createHandlerMapping(hubConnObjectClient, hubRouter, manager.ConnectionCollection, manager.ConnectionResource)

    // hub-device API
    if hubDeviceObjectClient == nil {
         hubDeviceObjectClient = manager.NewHubDeviceObjectManager()
    }
    createHandlerMapping(hubDeviceObjectClient, hubRouter, manager.DeviceCollection, manager.DeviceResource)

    // device API
    if deviceObjectClient == nil {
         deviceObjectClient = manager.NewDeviceObjectManager()
    }
    createHandlerMapping(deviceObjectClient, olRouter, manager.DeviceCollection, manager.DeviceResource)

    // device-connection API
    if deviceConnObjectClient == nil {
         deviceConnObjectClient = manager.NewDeviceConnObjectManager()
    }
    createHandlerMapping(deviceConnObjectClient, devRouter, manager.ConnectionCollection, manager.ConnectionResource)

    // iprange API
    if ipRangeObjectClient == nil {
         ipRangeObjectClient = manager.NewIPRangeObjectManager()
    }
    createHandlerMapping(ipRangeObjectClient, olRouter, manager.IPRangeCollection, manager.IPRangeResource)

    // Add depedency
    overlayObjectClient.AddOwnResManager(proposalObjectClient)
    overlayObjectClient.AddOwnResManager(hubObjectClient)
    overlayObjectClient.AddOwnResManager(deviceObjectClient)
    overlayObjectClient.AddOwnResManager(ipRangeObjectClient)
    hubObjectClient.AddOwnResManager(hubDeviceObjectClient)
    deviceObjectClient.AddOwnResManager(hubDeviceObjectClient)

    proposalObjectClient.AddDepResManager(overlayObjectClient)
    hubObjectClient.AddDepResManager(overlayObjectClient)
    deviceObjectClient.AddDepResManager(overlayObjectClient)
    ipRangeObjectClient.AddDepResManager(overlayObjectClient)
    hubDeviceObjectClient.AddDepResManager(hubObjectClient)
    hubDeviceObjectClient.AddDepResManager(deviceObjectClient)
    hubConnObjectClient.AddDepResManager(hubObjectClient)
    deviceConnObjectClient.AddDepResManager(deviceObjectClient)

    return router
}
