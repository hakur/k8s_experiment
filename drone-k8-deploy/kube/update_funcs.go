package kube

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	jobv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
	"strings"
)

var gcPolicy metav1.DeletionPropagation = "Background"
var DeletePolicy = &metav1.DeleteOptions{PropagationPolicy: &gcPolicy}

// UpdateResourceList 创建列表资源
func UpdateResourceList(confList *[]interface{}) (err error) {
	for k, v := range *confList {
		err := UpdateResource(v)
		if err != nil {
			logrus.Error("序列["+strconv.Itoa(k)+"]", err.Error())
			err = errors.New("序列[" + strconv.Itoa(k) + "]" + " " + err.Error())
		}
	}
	return err
}

// UpdateResource 创建资源  注意这个函数传入的是指针且会修改传入的指针变量的值
func UpdateResource(obj interface{}) (err error) {
	switch cfg := obj.(type) {
	case *appsv1.StatefulSet:
		err = updateStatefulSet(cfg)
	case *appsv1.Deployment:
		err = updateDeployment(cfg)
	case *jobv1.Job:
		err = updateJob(cfg)
	case *corev1.Service:
		err = updateService(cfg)
	case *corev1.ConfigMap:
		updateConfigMap(cfg)
	case *extv1beta1.Ingress:
		updateIngress(cfg)
	default:
		fmt.Println(cfg)
		err = errors.New("找不到支持的资源类型,可能传入的变量不是一层指针")
	}

	return err
}

func getImageTag(imageUrl string) string {
	arr := strings.Split(imageUrl, ":")
	if len(imageUrl) == 2 {
		return arr[1]
	}
	return "latest"
}
func tryCreateResource(obj interface{}) (err error) {
	client := GetDefaultClient()

	switch cfg := obj.(type) {
	case *appsv1.StatefulSet:
		o, err := client.AppsV1().StatefulSets(cfg.ObjectMeta.Namespace).Create(cfg)
		if err != nil {
			return err
		} else {
			*cfg = *o
		}
	case *appsv1.Deployment:
		o, err := client.AppsV1().Deployments(cfg.ObjectMeta.Namespace).Create(cfg)
		if err != nil {
			return err
		} else {
			*cfg = *o
		}

	case *jobv1.Job:
		o, err := client.BatchV1().Jobs(cfg.ObjectMeta.Namespace).Create(cfg)
		if err != nil {
			return err
		} else {
			*cfg = *o
		}
	case *corev1.Service:
		o, err := client.CoreV1().Services(cfg.ObjectMeta.Namespace).Create(cfg)
		if err != nil {
			return err
		} else {
			*cfg = *o
		}
	case *corev1.ConfigMap:
		o, err := client.CoreV1().ConfigMaps(cfg.ObjectMeta.Namespace).Create(cfg)
		if err != nil {
			return err
		} else {
			*cfg = *o
		}
	case *extv1beta1.Ingress:
		o, err := client.ExtensionsV1beta1().Ingresses(cfg.ObjectMeta.Namespace).Create(cfg)
		if err != nil {
			return err
		} else {
			*cfg = *o
		}
	default:
		fmt.Println(cfg)
		err = errors.New("找不到支持的资源类型,可能传入的变量不是一层指针")
	}
	if err != nil {
		return fmt.Errorf("resource not found try to create it but got error -> %s", err.Error())
	}
	return nil
}
