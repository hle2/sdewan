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
    "log"
    //"errors"
    "strings"
    "encoding/json"
    "encoding/base64"

    "github.com/onap/multicloud-k8s/src/orchestrator/pkg/infra/db"
    "github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/module"
    //"github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/client"
    //"github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/resource"
    pkgerrors "github.com/pkg/errors"
)

const PUBLICIP = "publicip"
const HUBTOHUB = "hub-to-hub"

type HubObjectKey struct {
    OverlayName string `json:"overlay-name"`
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
    overlay := GetManagerset().Overlay
    overlay_name := m[OverlayResource]
    to := t.(*module.HubObject)
    hub_name := to.Metadata.Name

    //Todo: Check if public ip can be used.
    var local_public_ip string
    var config []byte
    config, err := base64.StdEncoding.DecodeString(to.Specification.KubeConfig)
    if err != nil {
            log.Println(err)
            return t, err
    }
    log.Println(to.Specification.KubeConfig)
    log.Println(string(config))

    local_public_ips := to.Specification.PublicIps
    log.Println("public ips: %+v", local_public_ips)
    kubeutil := GetKubeConfigUtil()
    config, local_public_ip, err = kubeutil.checkKubeConfigAvail(config, local_public_ips, "6443")
    if err == nil {
        log.Println("Verified public ip " + local_public_ip)
        stat := make(map[string]string)
        stat[PUBLICIP] = local_public_ip
        to.Status.Data = stat
    } else {
        log.Println(err)
    }

    
    //Create cert for ipsec connection
    log.Println("Create Certificate: " + hub_name + "-cert")
    crt, key, err:= overlay.CreateCertificate(overlay_name, hub_name + "-cert")
    log.Println("Crt: \n" + crt)
    log.Println("Key: \n" + key)
    if err != nil {
        log.Println(err)
        return t, err
    }

    /*
    //Todo: Get all available proposals
    proposal := GetManagerset().Proposal
    proposals, err := proposal.GetObjects(m)
    if len(proposals) == 0 || err != nil {
        log.Println("Missing Proposal in the overlay\n")
        log.Println(err)
        return t, errors.New("Error in getting proposals")
    }
    var all_proposals []string
    for i:= 0 ; i < len(proposals); i++ {
            all_proposals = append(all_proposals, proposals[i].(*module.ProposalObject).Metadata.Name)
    }
    */


    //Get all available hub objects
    hub := GetManagerset().Hub
    hubs, err := hub.GetObjects(m)
    if err != nil {
            log.Println(err)
    }

    if len(hubs) > 0 && err == nil {
        for i := 0; i < len(hubs); i++ {
            err := overlay.SetupConnection(m, t, hubs[i], HUBTOHUB, NameSpaceName)
            if err != nil {
                log.Println("Setup connection with " + hubs[i].(*module.HubObject).Metadata.Name + " failed.")
            }    
            /*
            remote_hub := hubs[i].(*module.HubObject)
            remote_hub_name := remote_hub.Metadata.Name
            remote_public_ip := remote_hub.Status.Data["publicip"]
            //Get RootCA
            root_ca, _, _ := overlay.GetCertificate(overlay_name)

            var remotecrt string
            var remotekey string

            cu, err := GetCertUtil()
            if err != nil {
                log.Println(err)
            } else {
                ready := cu.IsCertReady(remote_hub_name + "-cert", NameSpaceName)
                if ready != true {
                    log.Println("Cert for remote hub is not ready")
                } else {
                    remotecrts, remotekey, err := cu.GetKeypair(remote_hub_name + "-cert", NameSpaceName)
                    remotecrt = strings.SplitAfter(remotecrts, "-----END CERTIFICATE-----")[0]
                    log.Println("Remote crt" + remotecrt)
                    log.Println("Remote key" + remotekey)
                    if err != nil {
                        log.Println(err)
                    }
                }
        }


        conn := resource.Connection{
            Name: "hubConn",
            ConnectionType: "tunnel",
            Mode: "start",
            Mark: default_mark,
            CryptoProposal: all_proposals,
        }

        ipsec_resource_local := resource.IpsecResource{
            Name: hub_name,
            Type: "VTI-based",
            Remote: remote_public_ip,
            AuthenticationMethod: "pubkey",
            PublicCert: remotecrt,
            PrivateCert: remotekey,
            SharedCA: root_ca,
            LocalIdentifier: local_public_ip,
            CryptoProposal: all_proposals,
            ForceCryptoProposal: "0",
            Connections: conn,

        }
        
        ipsec_resource_remote := resource.IpsecResource{
            Name: remote_hub_name,
            Type: "VTI-based",
            Remote: local_public_ip,
            AuthenticationMethod: "pubkey",
            PublicCert: crt,
            PrivateCert: key,
            SharedCA: root_ca,
            LocalIdentifier: remote_public_ip,
            CryptoProposal: all_proposals,
            ForceCryptoProposal: "0",
            Connections: conn,
        }

        resutil := NewResUtil()
        resutil.AddResource(to, "create", &ipsec_resource_local)
        resutil.AddResource(remote_hub, "create", &ipsec_resource_remote)

        err = resutil.Deploy("YAML")

        if err != nil {
            return c.CreateEmptyObject(), pkgerrors.Wrap(err, "Unable to create the object: fail to deploy resource")
        }
        */   
        }
        t, err = GetDBUtils().CreateObject(c, m, t)
    } else {

        t, err = GetDBUtils().CreateObject(c, m, t)
    }

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
    overlay := GetManagerset().Overlay
    hub_name := m[HubResource]
    log.Println("Delete Certificate: " + hub_name + "-cert")
    overlay.DeleteCertificate(hub_name + "-cert")
        // DB Operation
    err := GetDBUtils().DeleteObject(c, m)
    return err

}

func GetHubCertificate(cert_name string, namespace string)(string, string, error){
    cu, err := GetCertUtil()
    if err != nil {
        log.Println(err)
        return "", "", err
    } else {
        ready := cu.IsCertReady(cert_name, namespace)
        if ready != true {
            log.Println("Cert for hub is not ready")
            return "", "", pkgerrors.New("Cert for hub is not ready")
        } else {
            crts, key, err := cu.GetKeypair(cert_name, namespace)
            crt := strings.SplitAfter(crts, "-----END CERTIFICATE-----")[0]
            if err != nil {
                log.Println(err)
                return "", "", err
            }
            return crt, key, nil
        }
    }
}