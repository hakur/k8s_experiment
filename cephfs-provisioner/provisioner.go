package main

import (
	"fmt"

	"github.com/kubernetes-incubator/external-storage/lib/controller"
	"github.com/pborman/uuid"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/apis/core/v1/helper"
)

type CephfsParams struct {
	Cluster              string
	Monitors             []string
	AdminID              string
	AdminSecret          string
	AdminSecretName      string
	AdminSecretNamespace string
	ClaimRoot            string
	DeterministicNames   bool
}

type Provisioner struct {
	EnableQuota bool
}

func NewProvisioner() *Provisioner {
	o := new(Provisioner)
	return o
}

// Provision : create persistenVolume , this func implement from controller.Provisioner
func (t *Provisioner) Provision(options controller.VolumeOptions) (pv *corev1.PersistentVolume, err error) {
	if options.PVC.Spec.Selector != nil {
		return nil, fmt.Errorf("claim Selector is not supported")
	}
	param := new(CephfsParams)
	err = ParseCephFSParams(param, options.Parameters)
	if err != nil {
		return nil, err
	}

	var sharedDirName string
	if param.DeterministicNames {
		sharedDirName = fmt.Sprintf(options.PVC.Name)
	} else {
		// create random share name
		sharedDirName = fmt.Sprintf("kubernetes-dynamic-pvc-%s", uuid.NewUUID())
	}

	// call script to create directory on cephfs,if there is no specific directory,an error will be reported
	mntPath := cephFSRootMntPath + "/" + sharedDirName
	err = mountCephFSRootDir(mntPath, param)
	if err != nil {
		err = fmt.Errorf("mount cepfs root to local failed -> %s", err.Error())
		return nil, err
	}
	pvDirPath := mntPath + "/" + PVBaseDirName + "/" + sharedDirName
	err = mkCephFSDirIfNotExist(pvDirPath)
	if err != nil {
		err = fmt.Errorf("create cephfs directory failed -> %s", err.Error())
		return nil, err
	}
	err = umountCephFSRootDir(mntPath)
	if err != nil {
		err = fmt.Errorf("umount <%s> failed -> %s", mntPath, err.Error())
		return nil, err
	}

	pv = &corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: options.PVName,
			Annotations: map[string]string{
				//provisionerIDAnn: p.identity,
				cephShareAnn: sharedDirName,
			},
		},
		Spec: corev1.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy: options.PersistentVolumeReclaimPolicy,
			AccessModes:                   options.PVC.Spec.AccessModes,
			MountOptions:                  options.MountOptions,
			Capacity: corev1.ResourceList{
				// Quotas are supported by the userspace client(ceph-fuse, libcephfs), or kernel client >= 4.17 but only on mimic clusters.
				// In other cases capacity is meaningless here.
				// If quota is enabled, provisioner will set ceph.quota.max_bytes on volume path.
				corev1.ResourceName(corev1.ResourceStorage): options.PVC.Spec.Resources.Requests[corev1.ResourceName(corev1.ResourceStorage)],
			},
			PersistentVolumeSource: corev1.PersistentVolumeSource{
				CephFS: &corev1.CephFSPersistentVolumeSource{
					Monitors: param.Monitors,
					Path:     "/" + PVBaseDirName + "/" + sharedDirName,
					SecretRef: &corev1.SecretReference{
						Name:      param.AdminSecretName,
						Namespace: "default",
					},
					User: param.AdminID,
				},
			},
		},
	}
	return pv, nil
}

// Delete : delete persistenVolume  , this func implement from controller.Provisioner
func (p *Provisioner) Delete(volume *corev1.PersistentVolume) error {
	if volume.Spec.PersistentVolumeReclaimPolicy == "Delete" {
		// get pvc storageClass info
		//Client.StorageV1beta1().StorageClasses().Get(helper.GetPersistentVolumeClass(volume), metav1.GetOptions{})
		storageClass, err := Client.StorageV1().StorageClasses().Get(helper.GetPersistentVolumeClass(volume), metav1.GetOptions{})
		if err != nil {
			logrus.Error("error when delete cephfs directory", volume.Spec.CephFS.Path, " pvc name is ", volume.Name, " pvc namespace is ", volume.Namespace, " -> ", err)
			return err
		}

		param := new(CephfsParams)
		err = ParseCephFSParams(param, storageClass.Parameters)
		if err != nil {
			logrus.Error("error when delete cephfs directory", volume.Spec.CephFS.Path, " pvc name is ", volume.Name, " pvc namespace is ", volume.Namespace, " -> ", err)
			return err
		}

		if err = deletePVCDir(volume.Spec.CephFS.Path, param); err != nil {
			logrus.Error("error when delete cephfs directory", volume.Spec.CephFS.Path, " pvc name is ", volume.Name, " pvc namespace is ", volume.Namespace, " -> ", err)
			return err
		}
		// delete directory if volume.Spec.PersistentVolumeReclaimPolicy is Delete
	}
	//litter.Dump(volume)
	return nil
}
