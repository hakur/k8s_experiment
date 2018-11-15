package kube

import (
	corev1 "k8s.io/api/core/v1"
)

func updateService(cfg *corev1.Service) (err error) {
	client := GetDefaultClient()
	o, err := client.CoreV1().Services(cfg.ObjectMeta.Namespace).Update(cfg)
	if err != nil {
		if err = tryCreateResource(cfg); err != nil {
			return err
		}
	} else {
		*cfg = *o
	}
	return nil
}
