package kube

import (
	extv1beta1 "k8s.io/api/extensions/v1beta1"
)

func updateIngress(cfg *extv1beta1.Ingress) (err error) {
	client := GetDefaultClient()
	o, err := client.ExtensionsV1beta1().Ingresses(cfg.ObjectMeta.Namespace).Update(cfg)
	if err != nil {
		if err = tryCreateResource(cfg); err != nil {
			return err
		}
	} else {
		*cfg = *o
	}
	return nil
}
