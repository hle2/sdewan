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

package resource

import (
        "log"
        "strings"
)

type Connection struct {
        Name           string
        ConnectionType string
        Mode           string
        LocalSourceIp  string
        LocalUpDown    string
        LocalFirewall  string
        RemoteSubnet   string
        RemoteSourceIp string
        RemoteUpDown   string
        RemoteFirewall string
        CryptoProposal []string
        Mark           string
        IfId           string
}

type IpsecResource struct {
        Name                 string
        Type                 string
        Remote               string
        AuthenticationMethod string
        CryptoProposal       []string
        LocalIdentifier      string
        RemoteIdentifier     string
        ForceCryptoProposal  string
        PresharedKey         string
        PublicCert      string
        PrivateCert     string
        SharedCA             string
        Connections          Connection
}

func (c *IpsecResource) GetName() string {
        return c.Name
}

func (c *IpsecResource) GetType() string {
        return c.Type
}

func (c *IpsecResource) ToYaml() string {
        p := strings.Join(c.CryptoProposal, ",")
        pr := strings.Join(c.Connections.CryptoProposal, ",")
        if c.AuthenticationMethod == "pubkey" {
            return `apiVersion: batch.sdewan.akraino.org/v1alpha1 
            kind: IpsecHost
            metadata:
              name:` +  c.Name + `
              namespace: default
              labels:
                sdewanPurpose:` + c.Name + `
            spec:
              name:` + c.Name + `
              type:` + c.Type + `
              remote:` + c.Remote + `
              authentication_method: `+ c.AuthenticationMethod +`
              local_public_cert:` + c.PublicCert + `
              local_private_cert:` + c.PrivateCert + `
              shared_ca:` + c.SharedCA + `
              local_identifier:` + c.LocalIdentifier + `
              force_crypto_proposal:` + c.ForceCryptoProposal + `
              crypto_proposal:` + p + `
              connections: 
              - name:` + c.Connections.Name + `
                conn_type:` + c.Connections.ConnectionType + `
                mode:` +  c.Connections.Mode + `
                mark:` +  c.Connections.Mark + `
                local_updown: /etc/updown
                crypto_proposal:` + pr + `
                 `
        } else if c.AuthenticationMethod == "psk" {
            return `apiVersion: batch.sdewan.akraino.org/v1alpha1 
            kind: IpsecHost
            metadata:
              name:` +  c.Name + `
              namespace: default
                labels:
                  sdewanPurpose:` + c.Name + `
            spec:
              name:` + c.Name + `
              type:` + c.Type + `
              remote:` + c.Remote + `
              authentication_method:` + c.AuthenticationMethod + `
              pre_shared_key:` + c.PresharedKey + `
              local_identifier:` + c.LocalIdentifier + `
              force_crypto_proposal:` + c.ForceCryptoProposal + `
              crypto_proposal:` + p + `
              connections: 
              - name:` + c.Connections.Name + `
                conn_type:` + c.Connections.ConnectionType + `
                mode:` + c.Connections.Mode + `
                mark:` + c.Connections.Mark + `
                local_updown: /etc/updown
                crypto_proposal:` + pr + `
                 `
        } else {
                log.Println("Unsupported authentication method.")
                return "Error"
        }
}
