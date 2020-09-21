package main

import (
    "log"
    "os"
    "io/ioutil"
    "sigs.k8s.io/yaml"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/client-go/kubernetes/scheme"
    "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
    pkgerrors "github.com/pkg/errors"
)

func DecodeYAMLFromFile(path string, into runtime.Object) (runtime.Object, error) {
    if _, err := os.Stat(path); err != nil {
        if os.IsNotExist(err) {
            return nil, pkgerrors.New("File " + path + " not found")
        } else {
            return nil, pkgerrors.Wrap(err, "Stat file error")
        }
    }

    rawBytes, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, pkgerrors.Wrap(err, "Read YAML file error")
    }

    decode := scheme.Codecs.UniversalDeserializer().Decode
    obj, _, err := decode(rawBytes, nil, into)
    if err != nil {
        return nil, pkgerrors.Wrap(err, "Deserialize YAML error")
    }

    return obj, nil
}

func DecodeYAMLFromData(data string, into runtime.Object) (runtime.Object, error) {
    decode := scheme.Codecs.UniversalDeserializer().Decode
    obj, _, err := decode([]byte(data), nil, into)
    if err != nil {
        return nil, pkgerrors.Wrap(err, "Deserialize YAML error")
    }

    return obj, nil
}

type KubeConfigUtil struct {
}

var kubeutil = KubeConfigUtil{}

func GetKubeConfigUtil() *KubeConfigUtil {
    return &kubeutil
}

func (c *KubeConfigUtil) toYaml(data *unstructured.Unstructured) (string, error) {
    byte_json, err := data.MarshalJSON()
    if err != nil {
        return "", pkgerrors.Wrap(err, "Fail to generate yaml")
    }

    byte_yaml, err := yaml.JSONToYAML(byte_json)
    if err != nil {
        return "", pkgerrors.Wrap(err, "Fail to generate yaml")
    }

    return string(byte_yaml), nil
}

func (c *KubeConfigUtil) UpdateK8sConfig(conf string, server string) (string, error) {
    conf_us_obj := &unstructured.Unstructured{}
    _, err := DecodeYAMLFromData(conf, conf_us_obj)
    if err == nil {
        conf_obj := conf_us_obj.UnstructuredContent()
        cluster_objs, _, err := unstructured.NestedSlice(conf_obj, "clusters")
        if err == nil {
            if len(cluster_objs) > 0 {
                cluster_obj := cluster_objs[0].(map[string]interface{})
                // remove certificate-authority-data
                unstructured.RemoveNestedField(cluster_obj, "cluster", "certificate-authority-data")
                // add insecure-skip-tls-verify
                err = unstructured.SetNestedField(cluster_obj, true, "cluster", "insecure-skip-tls-verify")
                if err == nil {
                    // set server
                    err = unstructured.SetNestedField(cluster_obj, server, "cluster", "server")
                    if err == nil {
                        err = unstructured.SetNestedSlice(conf_obj, cluster_objs, "clusters")
                        if err == nil {
                            return c.toYaml(conf_us_obj)
                        }
                    }
                }
            } else {
                return "", pkgerrors.New("UpdateK8sConfig: No cluster")
            }
        }
    }

    return "", pkgerrors.Wrap(err, "UpdateK8sConfig")
}

func TestYaml() string {
    return `apiVersion: sdewan.akraino.org/v1alpha1
kind: IpsecProposal
metadata:
  name: proposal1
  namespace: default
  labels:
    sdewanPurpose: cnf
spec:
  encryption_algorithm: sha
  hash_algorithm: sha256
  dh_group: dh`
}

func main() {
    obj := &unstructured.Unstructured{}
    _, err := DecodeYAMLFromData(TestYaml(), obj)
    if err != nil {
        log.Println(err)
    } else {
        log.Println(obj)
    }

    conf_bytes, err := ioutil.ReadFile("./admin.conf")
    if err != nil {
        log.Println(err)
    }

    ku := GetKubeConfigUtil()
    conf_new, err := ku.UpdateK8sConfig(string(conf_bytes), "192.168.1.1:6000")

    if err != nil {
        log.Println(err)
    }

    new_file, err := os.OpenFile("admin_update.yaml", os.O_CREATE | os.O_WRONLY | os.O_TRUNC, 0666)
    _, err = new_file.Write([]byte(conf_new))
    if err != nil {
       log.Println(err)
    }
}
