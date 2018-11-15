package kube

import (
	corev1 "k8s.io/api/core/v1"
)

func updateConfigMap(cfg *corev1.ConfigMap) (err error) {
	client := GetDefaultClient()
	o, err := client.CoreV1().ConfigMaps(cfg.ObjectMeta.Namespace).Update(cfg)
	if err != nil {
		if err = tryCreateResource(cfg); err != nil {
			return err
		}
	} else {
		*cfg = *o
	}
	return nil
}
