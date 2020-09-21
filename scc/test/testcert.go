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
// Based on Code: https://github.com/johandry/klient

package main

import (
    "log"
    "context"
    "time"

    cmmeta "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/apimachinery/pkg/runtime/schema"
    "k8s.io/client-go/kubernetes/scheme"
    "k8s.io/client-go/kubernetes"
    "k8s.io/apimachinery/pkg/util/wait"
    v1 "k8s.io/api/core/v1"
    corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
    certmanagerversioned "github.com/jetstack/cert-manager/pkg/client/clientset/versioned"
    certmanagerv1beta1 "github.com/jetstack/cert-manager/pkg/client/clientset/versioned/typed/certmanager/v1beta1"
    v1beta1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1beta1"
)

type KubernetesClient struct {
    Context string
    KubeConfig string
}

func NewClient(context string, kubeConfig string) *KubernetesClient {
    return &KubernetesClient{
        Context: context,
        KubeConfig: kubeConfig,
    }
}

func (c *KubernetesClient) ToRESTConfig() (*rest.Config, error) {
    // From: k8s.io/kubectl/pkg/cmd/util/kubectl_match_version.go > func setKubernetesDefaults()
    config, err := c.ToRawKubeConfigLoader().ClientConfig()
    if err != nil {
        return nil, err
    }

    if config.GroupVersion == nil {
        config.GroupVersion = &schema.GroupVersion{Group: "", Version: "v1"}
    }
    if config.APIPath == "" {
        config.APIPath = "/api"
    }
    if config.NegotiatedSerializer == nil {
        // This codec config ensures the resources are not converted. Therefore, resources
        // will not be round-tripped through internal versions. Defaulting does not happen
        // on the client.
        config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
    }

    rest.SetKubernetesDefaults(config)
    return config, nil
}

// ToRawKubeConfigLoader creates a client using the following rules:
// 1. builds from the given kubeconfig path, if not empty
// 2. use the in cluster factory if running in-cluster
// 3. gets the factory from KUBECONFIG env var
// 4. Uses $HOME/.kube/factory
// It's required to implement the interface genericclioptions.RESTClientGetter
func (c *KubernetesClient) ToRawKubeConfigLoader() clientcmd.ClientConfig {
    loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
    loadingRules.DefaultClientConfig = &clientcmd.DefaultClientConfig
    if len(c.KubeConfig) != 0 {
        loadingRules.ExplicitPath = c.KubeConfig
    }
    configOverrides := &clientcmd.ConfigOverrides{
        ClusterDefaults: clientcmd.ClusterDefaults,
    }
    if len(c.Context) != 0 {
        configOverrides.CurrentContext = c.Context
    }

    return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
}

func (c *KubernetesClient) GetCMClients() (certmanagerv1beta1.CertmanagerV1beta1Interface, corev1.CoreV1Interface, error) {
    config, err := c.ToRESTConfig()
    if err != nil {
        return nil, nil, err
    }

    cmclientset, err := certmanagerversioned.NewForConfig(config)
    if err != nil {
        return nil, nil, err
    }

    k8sclientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        return nil, nil, err
    }

    return cmclientset.CertmanagerV1beta1(), k8sclientset.CoreV1(), nil
}

type CertUtil struct {
    client certmanagerv1beta1.CertmanagerV1beta1Interface
    k8sclient corev1.CoreV1Interface
}

var certutil = CertUtil{}

func GetCertUtil() (*CertUtil, error) {
    var err error
    if certutil.client == nil || certutil.k8sclient == nil {
//        certutil.client = client.NewClient("", "")
        certutil.client, certutil.k8sclient, err = NewClient("", "").GetCMClients()
    }

    return &certutil, err
}

func (c *CertUtil) CreateNamespace(name string) (*v1.Namespace, error) {
    ns, err := c.k8sclient.Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
    
        if err == nil {
    
            
        return ns, nil
    
        }

    log.Println("Create Namespace: " + name)
    return c.k8sclient.Namespaces().Create(context.TODO(), &v1.Namespace{
        ObjectMeta: metav1.ObjectMeta{
    
            
            
        Name: name,
    
            
        },
    }, metav1.CreateOptions{})
}

func (c *CertUtil) DeleteNamespace(name string) error {
    return c.k8sclient.Namespaces().Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (c *CertUtil) GetIssuer(name string, namespace string) (*v1beta1.Issuer, error) {
    return c.client.Issuers(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (c *CertUtil) DeleteIssuer(name string, namespace string) error {
    return c.client.Issuers(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (c *CertUtil) CreateSelfSignedIssuer(name string, namespace string) (*v1beta1.Issuer, error) {
    issuer, err := c.GetIssuer(name, namespace)
    if err == nil {
        return issuer, nil
    }

    // Not existing issuer, create a new one
    return c.client.Issuers(namespace).Create(context.TODO(), &v1beta1.Issuer{
    
            
        ObjectMeta: metav1.ObjectMeta{
    
            
            
        Name: name,
    
            
        },
    
            
        Spec: v1beta1.IssuerSpec{
    
            
            
        IssuerConfig: v1beta1.IssuerConfig{
    
            
            
            
        SelfSigned: &v1beta1.SelfSignedIssuer{
    
            
            
            
        },
    
            
            
        },
        },
    
        }, metav1.CreateOptions{})
}

func (c *CertUtil) CreateCAIssuer(name string, namespace string, caname string) (*v1beta1.Issuer, error) {
    issuer, err := c.GetIssuer(name, namespace)
    if err == nil {
        return issuer, nil
    }

    // Not existing issuer, create a new one
    return c.client.Issuers(namespace).Create(context.TODO(), &v1beta1.Issuer{
    
            
        ObjectMeta: metav1.ObjectMeta{
    
            
            
        Name: name,
    
            
        },
    
            
        Spec: v1beta1.IssuerSpec{
    
            
            
        IssuerConfig: v1beta1.IssuerConfig{
    
            
            
            
        CA: &v1beta1.CAIssuer{
                    SecretName: c.GetCertSecretName(caname),
    
            
            
            
        },
    
            
            
        },
        },
    
        }, metav1.CreateOptions{})
}

func (c *CertUtil) GetCertSecretName(name string) string {
    return name + "-cert-secret"
}

func (c *CertUtil) GetCertificate(name string, namespace string) (*v1beta1.Certificate, error) {
    return c.client.Certificates(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (c *CertUtil) DeleteCertificate(name string, namespace string) error {
    return c.client.Certificates(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (c *CertUtil) CreateCertificate(name string, namespace string, issuer string, isCA bool) (*v1beta1.Certificate, error) {
    cert, err := c.GetCertificate(name, namespace)
    if err == nil {
        return cert, nil
    }

    // Not existing cert, create a new one
    // Todo: add Duration, RenewBefore, DNSNames
    return c.client.Certificates(namespace).Create(context.TODO(), &v1beta1.Certificate{
    
            
        ObjectMeta: metav1.ObjectMeta{
    
            
            
        Name: name,
    
            
        },
    
            
        Spec: v1beta1.CertificateSpec{
            CommonName: name,
            // Duration: duration,
            // RenewBefore: renewBefore,
            // DNSNames: dnsNames,
            SecretName: c.GetCertSecretName(name),
            IssuerRef: cmmeta.ObjectReference {
                Name: issuer,
                Kind: "Issuer",
            },
            IsCA: isCA,
        },
    
        }, metav1.CreateOptions{})
}

func (c *CertUtil) IsCertReady(name string, namespace string) bool {
    err := wait.PollImmediate(time.Second, time.Second*20,
    
            
        func() (bool, error) {
    
            
            
        var err error
            var crt *v1beta1.Certificate
    
            
            
        crt, err = c.GetCertificate(name, namespace)
    
            
            
        if err != nil {
    
            
            
            
        log.Println("Failed to find certificate " + name + ": " + err.Error())
                return false, err
    
            
            
        }
            curConditions := crt.Status.Conditions
            for _, cond := range curConditions {
                if v1beta1.CertificateConditionReady == cond.Type && cmmeta.ConditionTrue == cond.Status {
                    return true, nil
                }
            }
            log.Println("Waiting for Certificate " + name + " to be ready.")
    
            
            
        return false, nil
    
            
        },
    
        )

    if err != nil {
        log.Println(err)
        return false
    }

    return true
}

func (c *CertUtil) GetKeypair(certname string, namespace string) (string, string, error) {
    if c.IsCertReady(certname, namespace) == false {
        return "", "", nil
    }

    secret, err ：= c.k8sclient.Secrets(namespace).Get(
        context.TODO(),
        c.GetCertSecretName(certname),
        metav1.GetOptions{})
    if err != nil {
        log.Println("Failed to get certificate's key pair: " + err.Error())
        return "", "", err
    }

    return string(secret.Data["tls.crt"]), string(secret.Data["tls.key"]), nil
}

func main() {
    cu, err := GetCertUtil()
    if err != nil {
        log.Println(err)
        return
    }

    log.Println("Success to get certmanager utility!")

    // create namespace
    _, err = cu.CreateNamespace("my-system")
    if err != nil {
        log.Println(err)
        return
    }

    log.Println("Success to create namespace!")

    // create a self-signed issuer as root issuer
    issuer_root, err := cu.CreateSelfSignedIssuer("my-root-issuer", "my-system")

    if err != nil {
        log.Println(err)
        return
    } else {
        log.Println(issuer_root)
    }

    // create root cert by root issuer
    cert_root, err := cu.CreateCertificate("my-root-cert", "my-system", "my-root-issuer", true)
    if err != nil {
        log.Println(err)
        return
    } else {
        log.Println(cert_root)
    }

    // create CA issuer based on root cert
    issuer_test, err := cu.CreateCAIssuer("my-test-issuer", "my-system", "my-root-cert")
    if err != nil {
        log.Println(err)
        return
    } else {
        log.Println(issuer_test)
    }

    // create a cert based on issuer_ca
    cert_test, err := cu.CreateCertificate("my-test-cert", "my-system", "my-test-issuer", false)
    if err != nil {
        log.Println(err)
        return
    } else {
        log.Println(cert_test)
    }

    crt, key, err := cu.GetKeypair("my-test-cert", "my-system")
    if err != nil {
        log.Println(err)
        return
    } else {
        log.Println("Crt: \n" + crt)
        log.Println("Key: \n" + key)
    }

    cu.DeleteCertificate("my-test-cert", "my-system")
    cu.DeleteIssuer("my-test-issuer", "my-system")
    cu.DeleteCertificate("my-root-cert", "my-system")
    cu.DeleteIssuer("my-root-issuer", "my-system")
    cu.DeleteNamespace("my-system")
}