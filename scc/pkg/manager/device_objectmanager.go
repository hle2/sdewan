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
    "time"
    //"strconv"
    "encoding/json"
    "encoding/base64"
    "strings"

    "k8s.io/apimachinery/pkg/util/wait"

    "github.com/open-ness/EMCO/src/orchestrator/pkg/infra/db"
    "github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/module"
    //"github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/client"
    "github.com/akraino-edge-stack/icn-sdwan/central-controller/src/scc/pkg/resource"
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

func (c *DeviceObjectManager) GetResourceName() string {
    return DeviceResource
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

    // initial Status
    v.Status.Data = make(map[string]string)
    return &v, err
}

func (c *DeviceObjectManager) PreProcessing(m map[string]string, t module.ControllerObject) error {
    to := t.(*module.DeviceObject)

    ipr_manager := GetManagerset().IPRange
    kubeutil := GetKubeConfigUtil()

    local_public_ips := to.Specification.PublicIps
    kube_config, err := base64.StdEncoding.DecodeString(to.Specification.KubeConfig)
    if err != nil {
        return pkgerrors.Wrap(err, "Fail to decode kubeconfig")
    }

    if len(local_public_ips) > 0{
        // Use public IP as external connection
        to.Status.Mode = 1
        
        kube_config, local_public_ip, err := kubeutil.checkKubeConfigAvail(kube_config, local_public_ips, "6443")
        if err != nil {
            return pkgerrors.Wrap(err, "Fail to verify public ip")
        }

        // Set IP in device
        log.Println("Use public ip " + local_public_ip)
        to.Status.Ip = local_public_ip

        // Set new kubeconfig in device
        to.Specification.KubeConfig = base64.StdEncoding.EncodeToString([]byte(kube_config))
    } else {
        // Use Hub as external connection
        to.Status.Mode = 2

        /*
        // validate hub information
        if to.Specification.ProxyHub == "" {
            return pkgerrors.New("Hub information is missing")
        }

        hm := make(map[string]string)
        hm[OverlayResource] = m[OverlayResource]
        hm[HubResource] = to.Specification.ProxyHub
        proxy_hub, err := hub_manager.GetObject(hm)
        if err != nil {
            return pkgerrors.Wrap(err, "Fail to get ProxyHub " + to.Specification.ProxyHub)
        }
        proxy_hub_obj := proxy_hub.(*module.HubObject)

        if to.Specification.ProxyHubPort == 0 {
            to.Specification.ProxyHubPort, err = proxy_hub_obj.AllocateProxyPort()
            if err != nil {
                return pkgerrors.Wrap(err, "Fail in " + to.Specification.ProxyHub)
            }
        } else {
            if proxy_hub_obj.IsProxyPortUsed(to.Specification.ProxyHubPort) {
                return pkgerrors.New("Proxy port is in-used")
            }
        }
        // update hub object with proxy-port 
        proxy_hub_obj.SetProxyPort(to.Specification.ProxyHubPort, to.Metadata.Name)
        */

        // allocate OIP for device
        overlay_name := m[OverlayResource]
        oip, err := ipr_manager.Allocate(overlay_name, to.Metadata.Name)
        if err != nil {
            return pkgerrors.Wrap(err, "Fail to allocate overlay ip for " + to.Metadata.Name)
        }

        // Set OIP in Device
        log.Println("Use overlay ip " + oip)
        to.Status.Ip = oip

        // Get all proposal resources
        proposal := GetManagerset().Proposal
        proposals, err := proposal.GetObjects(m)
        if len(proposals) == 0 || err != nil {
            log.Println("Missing Proposal in the overlay\n")
            log.Println(err)
            return pkgerrors.New("Error in getting proposals")
        }

        var all_proposal []string
        var proposalresource []*resource.ProposalResource
        for i:= 0 ; i < len(proposals); i++ {
            proposal_obj := proposals[i].(*module.ProposalObject)
            all_proposal = append(all_proposal, proposal_obj.Metadata.Name)
            // pr := resource.ProposalResource{proposals[i].(*module.ProposalObject).Metadata.Name, proposals[i].(*module.ProposalObject).Specification.Encryption, proposals[i].(*module.ProposalObject).Specification.Hash, proposals[i].(*module.ProposalObject).Specification.DhGroup}
            pr := proposal_obj.ToResource()
            proposalresource = append(proposalresource, pr)
        }

        //Extract SCC cert/key
	cu, err := GetCertUtil()
	if err != nil {
		log.Println("Getting certutil error")
	}
        crts, key, err := cu.GetKeypair(SCCCertName, NameSpaceName)
        root_ca := strings.SplitAfter(crts, "-----END CERTIFICATE-----")[1]
        crt := strings.SplitAfter(crts, "-----END CERTIFICATE-----")[0]

        // Build up ipsec resource
        scc_conn := resource.Connection{
            Name: DEFAULT_CONN,
            ConnectionType: CONN_TYPE,
            Mode: MODE,
            Mark: DEFAULT_MARK,
            RemoteSourceIp: oip,
            LocalUpDown: DEFAULT_UPDOWN,
            CryptoProposal: all_proposal,
        }

	scc_ipsec_resource := resource.IpsecResource{
            Name: "localto" + strings.ToLower(strings.Replace(to.Metadata.Name, "-", "", -1)),
            Type: VTI_MODE,
            Remote: ANY,
            AuthenticationMethod: PUBKEY_AUTH,
            PublicCert: base64.StdEncoding.EncodeToString([]byte(crt)),
            PrivateCert: base64.StdEncoding.EncodeToString([]byte(key)),
            SharedCA: base64.StdEncoding.EncodeToString([]byte(root_ca)),
            LocalIdentifier: "CN="+ SCCCertName + "-cert",
            RemoteIdentifier: "CN=" + to.Metadata.Name + "-cert",
            CryptoProposal: all_proposal,
            ForceCryptoProposal: FORCECRYPTOPROPOSAL,
            Connections: scc_conn,
        }

        scc := module.EmptyObject{
            Metadata: module.ObjectMetaData{"local", "", "", ""}}
        // Add and deploy resource
        resutil := NewResUtil()
        resutil.AddResource(&scc, "create", &scc_ipsec_resource)
        _, err = resutil.Deploy("localto" + to.Metadata.Name, "YAML")

        //Reserve ipsec resource to device object
        r := make(map[string]string)
        res_str, err := resource.GetResourceBuilder().ToString(&scc_ipsec_resource)
        r["scc_ipsec_resource"] = res_str
        to.Status.Data = r

        /*
        // Deploy SNAT rule in Hub to enable k8s API access proxy to device
        err = overlay_manager.SetupHubProxy(m, proxy_hub_obj, to, NameSpaceName)
        if err != nil {
            proxy_hub_obj.UnsetProxyPort(to.Specification.ProxyHubPort)
            ipr_manager.Free(overlay_name, oip)
            return pkgerrors.Wrap(err, "Fail to Setup hub proxy for " + to.Metadata.Name)
        }
        */

        // Check device availability
        //hub_ips := []string{proxy_hub_obj.Status.Ip}
        /*local_cnf_ip := kubeclient.getLocalCNFAddress()
        local_device_proxy_port := kubeclient.AllocateProxyPort()
        networks := []string{}

        // DNAT rule
        localZoneResource := resource.FirewallZoneResource{
            Name: "local_" + kubeclient.GetLocalIfName()),
            Network: append(networks, kubeclient.GetLocalIfName()),
            Input: ACCEPT,
            Output: ACCEPT,
            Forward: ACCEPT,
            MASQ: "0",
            MTU_FIX: "1",
        }

        localDnatResource := resource.FirewallDnatResource{
            Name: "localto" + strings.ToLower(strings.Replace(to.Metadata.Name, "-", "", -1)),
            Source: localZoneResource.GetName(),
            SourceDestIP: local_cnf_ip,
            SourceDestPort: strconv.Itoa(local_device_proxy_port),
            DestinationIP: to.Status.Ip,
            DestinationPort: DEFAULT_K8S_API_SERVER_PORT,
            Protocol: BASE_PROTOCOL,
        }
        */
        var ips []string
	ips = append(ips, oip)

        err = wait.PollImmediate(time.Second*5, time.Second*30,
            func() (bool, error) {
                kube_config, _, err := kubeutil.checkKubeConfigAvail(kube_config, ips, DEFAULT_K8S_API_SERVER_PORT)
                if err != nil {
                    log.Println("Waiting for scc connection to be set up.")
                }
                // Set new kubeconfig in device
                // Todo: to set kubeconfig even when timeout
                to.Specification.KubeConfig = base64.StdEncoding.EncodeToString([]byte(kube_config))
                return true, nil
            },
        )

        if err != nil {
            log.Println(err)
            return pkgerrors.Wrap(err, "Fail to connect to scc.")
        }

        /*
        // save proxy hub information
        _, err = GetDBUtils().UpdateObject(hub_manager, hm, proxy_hub_obj)
        */
    }
    return nil

}

func (c *DeviceObjectManager) CreateObject(m map[string]string, t module.ControllerObject) (module.ControllerObject, error) {
    err := c.PreProcessing(m, t)
    if err != nil {
        return c.CreateEmptyObject(), err
    }

    overlay_manager := GetManagerset().Overlay

    to := t.(*module.DeviceObject)
    overlay_name := m[OverlayResource]
    device_name := to.Metadata.Name
    
    //Create cert for ipsec connection
    log.Println("Create Certificate: " + device_name + "-cert")
    cert := GetManagerset().Cert
    _, _, _, err = cert.GetOrCreateDC(overlay_name, device_name)
    if err != nil {
        log.Println(err)
        return t, err
    }

    devices, err := c.GetObjects(m)
    if err != nil {
        log.Println(err)
        return t, nil
    }

    for i := 0; i < len(devices); i++ {
        dev :=  devices[i].(*module.DeviceObject)
        if to.Status.Mode == 1 || dev.Status.Mode == 1 {
            err = overlay_manager.SetupConnection(m, to, dev, DEVICETODEVICE, NameSpaceName)
            if err != nil {
                log.Println(err)
            }
        }
    }

    // DB Operation
    t, err = GetDBUtils().CreateObject(c, m, t)

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
    t, err := c.GetObject(m)
    if err != nil {
        return nil
    }

    //overlay_manager := GetManagerset().Overlay
    ipr_manager := GetManagerset().IPRange

    overlay_name := m[OverlayResource]
    device_name := m[DeviceResource]

    // Reset all IpSec connection setup by this device
    /*
    conns, err := conn_manager.GetObjects(overlay_name, module.CreateEndName(t.GetType(), device_name))
    if err != nil {
        log.Println(err)
    } else {
        for i := 0; i < len(conns); i++ {
            conn :=  conns[i].(*module.ConnectionObject)
            err = conn_manager.Undeploy(overlay_name, *conn)
            if err != nil {
                log.Println(err)
            }
        }
    }*/

    to := t.(*module.DeviceObject)
    if to.Status.Mode == 2 {
        // Free OIP
        ipr_manager.Free(overlay_name, to.Status.Ip)

        scc := module.EmptyObject{
            Metadata: module.ObjectMetaData{"local", "", "", ""}}

	resutil := NewResUtil()
        r_str := to.Status.Data["scc_ipsec_resource"]
        r, _ := resource.GetResourceBuilder().ToObject(r_str)
        resutil.AddResource(&scc, "create", r)
	resutil.Undeploy("localto" + device_name, "YAML")

        // Free Hub Proxy port
        /*
        hm := make(map[string]string)
        hm[OverlayResource] = overlay_name
        hm[HubResource] = to.Specification.ProxyHub
        proxy_hub, err := hub_manager.GetObject(hm)
        if err != nil {
            log.Println(err)
        } else {
            proxy_hub_obj := proxy_hub.(*module.HubObject)

            // unset hub object with proxy-port 
            proxy_hub_obj.UnsetProxyPort(to.Specification.ProxyHubPort)
            _, err = GetDBUtils().UpdateObject(hub_manager, hm, proxy_hub_obj)
            if err != nil {
                log.Println(err)
            }
        }*/
    }

    // Delete certificate
    //log.Println("Delete Certificate: " + device_name + "-cert")
    //overlay_manager.DeleteCertificate(device_name + "-cert")

    // DB Operation
    err = GetDBUtils().DeleteObject(c, m)

    return err
}

func GetDeviceCertificate(overlay_name string, device_name string)(string, string, error){
    cert := GetManagerset().Cert
    _, crts, key, err := cert.GetOrCreateDC(overlay_name, device_name)
    if err != nil {
        log.Println("Error in get cert for device ...")
            return "", "", err
    }
    
    crt := strings.SplitAfter(crts, "-----END CERTIFICATE-----")[0]
    return crt, key, nil
}

