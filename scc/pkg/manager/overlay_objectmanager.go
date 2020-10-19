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
    "strings"
    "encoding/json"
    "encoding/base64"
    "github.com/onap/multicloud-k8s/src/orchestrator/pkg/infra/db"
    "github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/module"
    "github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/resource"
    pkgerrors "github.com/pkg/errors"
)

const DEFAULT_MARK = "30"
const VTI_MODE = "VTI-based"
const PUBKEY_AUTH = "pubkey"
const FORCECRYPTOPROPOSAL = "0"
const DEFAULT_CONN = "Connection"
const DEFAULT_UPDOWN = "/etc/updown"
const CONN_TYPE = "tunnel"
const MODE = "start"
const OVERLAYIP = "overlayip"
const HUBTOHUB = "hub-to-hub"
const HUBTODEVICE = "hub-to-device"
const DEVICETODEVICE = "device-to-device"

type OverlayObjectKey struct {
    OverlayName string `json:"overlay-name"`
}

// OverlayObjectManager implements the ControllerObjectManager
type OverlayObjectManager struct {
    BaseObjectManager
}

func NewOverlayObjectManager() *OverlayObjectManager {
    return &OverlayObjectManager{
        BaseObjectManager {
            storeName:  StoreName,
            tagMeta:    "overlay",
            depResManagers: []ControllerObjectManager {},
            ownResManagers: []ControllerObjectManager {},
        },
    }
}

func (c *OverlayObjectManager) IsOperationSupported(oper string) bool {
    return true
}

func (c *OverlayObjectManager) CreateEmptyObject() module.ControllerObject {
    return &module.OverlayObject{}
}

func (c *OverlayObjectManager) GetStoreKey(m map[string]string, t module.ControllerObject, isCollection bool) (db.Key, error) {
    key := OverlayObjectKey{""}

    if isCollection == true {
        return key, nil
    }

    to := t.(*module.OverlayObject)
    meta_name := to.Metadata.Name
    res_name := m[OverlayResource]

    if res_name != "" {
        if meta_name != "" && res_name != meta_name {
            return key, pkgerrors.New("Resource name unmatched metadata name")
        }

        key.OverlayName = res_name
    } else {
        if meta_name == "" {
            return key, pkgerrors.New("Unable to find resource name")
        }

        key.OverlayName = meta_name
    }

    return key, nil;
}

func (c *OverlayObjectManager) ParseObject(r io.Reader) (module.ControllerObject, error) {
    var v module.OverlayObject
    err := json.NewDecoder(r).Decode(&v)

    return &v, err
}

func (c *OverlayObjectManager) CreateObject(m map[string]string, t module.ControllerObject) (module.ControllerObject, error) {
    // for rsync test
     resutil := NewResUtil()

     deviceObject := module.OverlayObject{
                Metadata: module.ObjectMetaData{"local", "", "", ""}, 
                Specification: module.OverlayObjectSpec{}}
     resutil.AddResource(&deviceObject, "Create", &resource.FileResource{"mycm", "ConfigMap", "mycm.yaml"})

     err2 := resutil.Deploy("YAML")

     if err2 != nil {
         return c.CreateEmptyObject(), pkgerrors.Wrap(err2, "Unable to create the object: fail to deploy resource")
     }

    // Create a issuer each overlay
    to := t.(*module.OverlayObject)
    overlay_name := to.Metadata.Name
    cu, err := GetCertUtil()
    if err != nil {
        log.Println(err)
    } else {
        // create overlay ca
        _, err := cu.CreateCertificate(c.CertName(overlay_name), NameSpaceName, RootCAIssuerName, true)
        if err == nil {
            // create overlay issuer
            _, err := cu.CreateCAIssuer(c.IssuerName(overlay_name), NameSpaceName, c.CertName(overlay_name))
            if err != nil {
                log.Println("Failed to create overlay[" + overlay_name +"] issuer: " + err.Error())
            }    
        } else {
            log.Println("Failed to create overlay[" + overlay_name +"] certificate: " + err.Error())
        }
    }

    // DB Operation
    t, err = GetDBUtils().CreateObject(c, m, t)

    return t, err
}

func (c *OverlayObjectManager) GetObject(m map[string]string) (module.ControllerObject, error) {
    // DB Operation
    t, err := GetDBUtils().GetObject(c, m)

    return t, err
}

func (c *OverlayObjectManager) GetObjects(m map[string]string) ([]module.ControllerObject, error) {
    // DB Operation
    t, err := GetDBUtils().GetObjects(c, m)

    return t, err
}

func (c *OverlayObjectManager) UpdateObject(m map[string]string, t module.ControllerObject) (module.ControllerObject, error) {
    // DB Operation
    t, err := GetDBUtils().UpdateObject(c, m, t)

    return t, err
}

func (c *OverlayObjectManager) DeleteObject(m map[string]string) error {
    overlay_name := m[OverlayResource]
    cu, err := GetCertUtil()
    if err != nil {
        log.Println(err)
    } else {
        err = cu.DeleteIssuer(c.IssuerName(overlay_name), NameSpaceName)
        if err != nil {
            log.Println("Failed to delete overlay[" + overlay_name +"] issuer: " + err.Error())
        }
        err = cu.DeleteCertificate(c.CertName(overlay_name), NameSpaceName)
        if err != nil {
            log.Println("Failed to delete overlay[" + overlay_name +"] certificate: " + err.Error())
        }
    }

    // DB Operation
    err = GetDBUtils().DeleteObject(c, m)

    return err
}

func (c *OverlayObjectManager) IssuerName(name string) string {
    return name + "-issuer"
}

func (c *OverlayObjectManager) CertName(name string) string {
    return name + "-cert"
}

func (c *OverlayObjectManager) CreateCertificate(oname string, cname string) (string, string, error) {
    cu, err := GetCertUtil()
    if err != nil {
        log.Println(err)
    } else {
        _, err := cu.CreateCertificate(cname, NameSpaceName, c.IssuerName(oname), false)
        if err != nil {
            log.Println("Failed to create overlay[" + oname +"] certificate: " + err.Error())
        } else {
            crts, key, err := cu.GetKeypair(cname, NameSpaceName)
            if err != nil {
                log.Println(err)
                return "", "", err
            } else {
                crt := strings.SplitAfter(crts, "-----END CERTIFICATE-----")[0]
                return crt, key, nil
            }
        }
    }

    return "", "", nil
}

func (c *OverlayObjectManager) DeleteCertificate(cname string) (string, string, error) {
    cu, err := GetCertUtil()
    if err != nil {
        log.Println(err)
    } else {
        err = cu.DeleteCertificate(cname, NameSpaceName)
        if err != nil {
            log.Println("Failed to delete " + cname +" certificate: " + err.Error())
        }
    }

    return "", "", nil
}

func (c *OverlayObjectManager) GetCertificate(oname string) (string, string, error) {
        cu, err := GetCertUtil()
        if err != nil {
                log.Println(err)
        } else {
                cname := c.CertName(oname)
                return cu.GetKeypair(cname, NameSpaceName)
        }
        return "", "", nil
}

//Set up Connection between objects
//Passing the original map resource, the two objects, connection type("hub-to-hub", "hub-to-device", "device-to-device") and namespace name.
func (c *OverlayObjectManager) SetupConnection(m map[string]string, m1 module.ControllerObject, m2 module.ControllerObject, conntype string, namespace string) error {
    //Get all proposals available in the overlay
    proposal := GetManagerset().Proposal
    proposals, err := proposal.GetObjects(m)
    if len(proposals) == 0 || err != nil {
        log.Println("Missing Proposal in the overlay\n")
        log.Println(err)
        return pkgerrors.New("Error in getting proposals")
    }
    var all_proposals []string
    var proposalresources []resource.ProposalResource
    for i:= 0 ; i < len(proposals); i++ {
            all_proposals = append(all_proposals, proposals[i].(*module.ProposalObject).Metadata.Name)
            pr := resource.ProposalResource{proposals[i].(*module.ProposalObject).Metadata.Name, proposals[i].(*module.ProposalObject).Specification.Encryption, proposals[i].(*module.ProposalObject).Specification.Hash, proposals[i].(*module.ProposalObject).Specification.DhGroup}
            proposalresources = append(proposalresources, pr)
    }

    //Get the overlay cert
    var root_ca string
    root_ca, _, _ = c.GetCertificate(m[OverlayResource])

    var Obj1 module.ControllerObject
    var Obj2 module.ControllerObject
    var obj1_ipsec_resource resource.IpsecResource
    var obj2_ipsec_resource resource.IpsecResource

    switch conntype {
    case HUBTOHUB:
        obj1 := m1.(*module.HubObject)
        obj2 := m2.(*module.HubObject)

        obj1_ip := obj1.Status.Data[PUBLICIP]
        obj2_ip := obj2.Status.Data[PUBLICIP]

        Obj1 = obj1
        Obj2 = obj2

        //Keypair
        obj1_crt, obj1_key, err := GetHubCertificate(obj1.GetCertName(),namespace)
        if err != nil {
            log.Println(err)
        }
        obj2_crt, obj2_key, err := GetHubCertificate(obj2.GetCertName(),namespace)
        if err != nil {
            log.Println(err)
        }
        //IpsecResources
        conn := resource.Connection{
            Name: DEFAULT_CONN,
            ConnectionType: CONN_TYPE,
            Mode: MODE,
            Mark: DEFAULT_MARK,
            LocalUpDown: DEFAULT_UPDOWN,
            CryptoProposal: all_proposals,
        }
        obj1_ipsec_resource = resource.IpsecResource{
            Name: strings.ToLower(strings.Replace(obj1.Metadata.Name, "-", "", -1)) + strings.ToLower(strings.Replace(obj2.Metadata.Name, "-", "", -1)),
            SdewanPurpose: obj1.Metadata.Name,
            Type: VTI_MODE,
            Remote: obj2_ip,
            AuthenticationMethod: PUBKEY_AUTH,
            PublicCert: base64.StdEncoding.EncodeToString([]byte(obj2_crt)),
            PrivateCert: base64.StdEncoding.EncodeToString([]byte(obj2_key)),
            SharedCA: base64.StdEncoding.EncodeToString([]byte(root_ca)),
            LocalIdentifier: obj1_ip,
            CryptoProposal: all_proposals,
            ForceCryptoProposal: FORCECRYPTOPROPOSAL,
            Connections: conn,
        }
        obj2_ipsec_resource = resource.IpsecResource{
            Name: strings.ToLower(strings.Replace(obj1.Metadata.Name, "-", "", -1)) + strings.ToLower(strings.Replace(obj2.Metadata.Name, "-", "", -1)), 
            SdewanPurpose: obj2.Metadata.Name,
            Type: VTI_MODE,
            Remote: obj1_ip,
            AuthenticationMethod: PUBKEY_AUTH,
            PublicCert: base64.StdEncoding.EncodeToString([]byte(obj1_crt)),
            PrivateCert: base64.StdEncoding.EncodeToString([]byte(obj1_key)),
            SharedCA: base64.StdEncoding.EncodeToString([]byte(root_ca)),
            LocalIdentifier: obj2_ip,
            CryptoProposal: all_proposals,
            ForceCryptoProposal: FORCECRYPTOPROPOSAL,
            Connections: conn,
        }
    // Todo: Hub-to-device connection
    case HUBTODEVICE:
    /*    obj1 := m1.(*module.HubOject)
        obj2 := m2.(*module.DeviceOject)

        obj1_ip := obj1.Status.Data[PUBLICIP]
        obj2_ip := obj2.Status.Data[OVERLAYIP]

        //Keypair
        obj1_crt, obj1_key, err := obj1.GetCertificate(namespace)
        if err != nil {
            log.Println(err)
        }
        obj2_crt, obj2_key, err := obj2.GetCertificate(namespace)
        if err != nil {
            log.Println(err)
        }

        //IpsecResources
        obj1_conn := resource.Connection{
            Name: DEFAULT_CONN,
            ConnectionType: CONN_TYPE,
            Mode: MODE,
            Mark: DEFAULT_MARK,
            RemoteSourceIp: 
            CryptoProposal: all_proposals,
        }
        obj2_conn := resource.Connection{
            Name: DEFAULT_CONN,
            ConnectionType: CONN_TYPE,
            Mode: MODE,
            Mark: DEFAULT_MARK,
            LocalSourceIp: "%config" //Need to use const
            CryptoProposal: all_proposals,
        }
        obj1_ipsec_resource := resource.IpsecResource{
            Name: obj1.Metadata.Name + obj2.Metadata.Name + "Conn",
            Type: VTI_MODE,
            Remote: obj2_ip,
            AuthenticationMethod: PUBKEY_AUTH,
            PublicCert: obj2_crt,
            PrivateCert: obj2_key,
            SharedCA: root_ca,
            LocalIdentifier: obj1_ip,
            CryptoProposal: all_proposals,
            ForceCryptoProposal: FORCECRYPTOPROPOSAL,
            Connections: conn,
        }
        obj2_ipsec_resource := resource.IpsecResource{
            Name: obj2.Metadata.Name + obj1.Metadata.Name + "Conn",
            Type: VTI_MODE,
            Remote: obj1_ip,
            AuthenticationMethod: PUBKEY_AUTH,
            PublicCert: obj1_crt,
            PrivateCert: obj1_key,
            SharedCA: root_ca,
            LocalIdentifier: obj2_ip,
            CryptoProposal: all_proposals,
            ForceCryptoProposal: FORCECRYPTOPROPOSAL,
            Connections: conn,
        }
        */
    //Todo: Device-to-device connection
    case DEVICETODEVICE:
    default:
        return pkgerrors.New("Unknown connection type")
    }

    //Add resource
    resutil := NewResUtil()
    for i :=0; i < len(proposalresources); i++ {
        resutil.AddResource(Obj1, "create", &proposalresources[i])
        resutil.AddResource(Obj2, "create", &proposalresources[i])
    }
    resutil.AddResource(Obj1, "create", &obj1_ipsec_resource)
    resutil.AddResource(Obj2, "create", &obj2_ipsec_resource)

    //Deploy resources
    err = resutil.Deploy("YAML")

    if err != nil {
        return pkgerrors.Wrap(err, "Unable to create the object: fail to deploy resource")
    }

    return nil

}