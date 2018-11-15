package kube

import (
	jobv1 "k8s.io/api/batch/v1"
)

func updateJob(cfg *jobv1.Job) (err error) {
	client := GetDefaultClient()
	tag := getImageTag(cfg.Spec.Template.Spec.Containers[0].Image)
	if tag == "latest" {
		client.BatchV1().Jobs(cfg.ObjectMeta.Namespace).Delete(cfg.Name, DeletePolicy)
	}
	o, err := client.BatchV1().Jobs(cfg.ObjectMeta.Namespace).Update(cfg)
	if err != nil {
		if err = tryCreateResource(cfg); err != nil {
			return err
		}
	} else {
		*cfg = *o
	}
	return nil
}
