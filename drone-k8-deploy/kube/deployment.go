package kube

import (
	appsv1 "k8s.io/api/apps/v1"
)

func updateDeployment(cfg *appsv1.Deployment) error {
	client := GetDefaultClient()
	tag := getImageTag(cfg.Spec.Template.Spec.Containers[0].Image)
	if tag == "latest" {
		client.AppsV1().Deployments(cfg.ObjectMeta.Namespace).Delete(cfg.Name, DeletePolicy)
	}
	o, err := client.AppsV1().Deployments(cfg.ObjectMeta.Namespace).Update(cfg)
	if err != nil {
		if err = tryCreateResource(cfg); err != nil {
			return err
		}
	} else {
		*cfg = *o
	}
	return nil
}
