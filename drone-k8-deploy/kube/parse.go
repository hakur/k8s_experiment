package kube

import (
	"errors"
	"fmt"
	//"github.com/sirupsen/logrus"
	"strings"

	//yamlv2 "gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	jobv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
)

/*type resourceKind struct {
	Kind string `yaml:"kind"`
}

func getResourceKind(yamlString string) string {
	t := resourceKind{}
	yamlv2.Unmarshal([]byte(yamlString), &t)
	return t.Kind
}*/

func ParseResources(yamlFileContent string) (*[]interface{}, error) {
	resources := make([]interface{}, 0, 3)
	objs := strings.Split(yamlFileContent, "---")
	for _, v := range objs {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		reader := strings.NewReader(v)
		_, kind, err := scheme.Codecs.UniversalDeserializer().Decode([]byte(v), nil, nil)
		if err != nil {
			return nil, fmt.Errorf("template parse error %s %s", "->", err.Error())
		}

		switch kind.Kind {
		case "Job":
			res := new(jobv1.Job)
			if err := yaml.NewYAMLOrJSONDecoder(reader, 512).Decode(res); err != nil {
				//logrus.Error("解析YAML模板错误", "kind=", kind, "->", err)
				return nil, err
			}
			resources = append(resources, res)
			break
		case "Service":
			res := new(corev1.Service)
			if err := yaml.NewYAMLOrJSONDecoder(reader, 512).Decode(res); err != nil {
				return nil, err
			}
			resources = append(resources, res)
			break
		case "Deployment":
			res := new(appsv1.Deployment)
			if err := yaml.NewYAMLOrJSONDecoder(reader, 512).Decode(res); err != nil {
				return nil, err
			}
			resources = append(resources, res)
			break
		case "Ingress":
			res := new(extv1beta1.Ingress)
			if err := yaml.NewYAMLOrJSONDecoder(reader, 512).Decode(res); err != nil {
				return nil, err
			}

			resources = append(resources, res)
			break
		default:
			println(v)
			err := errors.New("Yaml文件未知的资源类型 kind=" + kind.Kind)
			return nil, err
		}
	}
	return &resources, nil
}
