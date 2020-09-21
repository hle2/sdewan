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

func NewRouter_test(
    testObjectClient manager.ControllerObjectManager) *mux.Router {

    router := mux.NewRouter()

    // Test
    if testObjectClient == nil {
        testObjectClient = manager.NewTestObjectManager()
    }

    testObjectHandler := ControllerHandler{client: testObjectClient}
    lcRouter := router.PathPrefix("/v2").Subrouter()
    lcRouter.HandleFunc(
        "/tests",
        testObjectHandler.createHandler).Methods("POST")
    lcRouter.HandleFunc(
        "/tests",
        testObjectHandler.getsHandler).Methods("GET")
    lcRouter.HandleFunc(
        "/tests/{test-name}",
        testObjectHandler.getHandler).Methods("GET")
    lcRouter.HandleFunc(
        "/tests/{test-name}",
        testObjectHandler.deleteHandler).Methods("DELETE")
    lcRouter.HandleFunc(
        "/tests/{test-name}",
        testObjectHandler.updateHandler).Methods("PUT")
    
    return router
}
