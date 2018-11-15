package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ParseCephFSParams(o *CephfsParams, params map[string]string) (err error) {
	for k, v := range params {
		switch strings.ToLower(k) {
		case "cluster":
			o.Cluster = v
		case "monitors":
			arr := strings.Split(v, ",")
			for _, m := range arr {
				o.Monitors = append(o.Monitors, m)
			}
		case "adminid":
			o.AdminID = v
		case "adminsecretname":
			o.AdminSecretName = v
		case "adminsecretnamespace":
			o.AdminSecretNamespace = v
		case "claimroot":
			o.ClaimRoot = v
		case "deterministicnames":
			// On error, strconv.ParseBool() returns false; leave that, as it is a perfectly fine default
			o.DeterministicNames, _ = strconv.ParseBool(v)
		default:
			return fmt.Errorf("invalid option %q", k)
		}
	}
	if o.AdminSecretName == "" {
		return fmt.Errorf("adminSecretName is empty")
	}
	if o.AdminSecret, err = fetchSecretFromK8s(o.AdminSecretName, o.AdminSecretNamespace); err != nil {
		return fmt.Errorf("get adminSecret string fialed from [adminSecretNamespace] %s , [adminSecretName] %s , error is %s", o.AdminSecretNamespace, o.AdminSecretName, err.Error())
	}

	if len(o.Monitors) < 1 {
		return fmt.Errorf("at least one monitor is required, storageClass parameters not configed with monitors")
	}
	return nil
}

func fetchSecretFromK8s(secretName, secretNamespace string) (secretBase64 string, err error) {
	secret, err := Client.CoreV1().Secrets(secretNamespace).Get(secretName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	for _, data := range secret.Data {
		return string(data), nil
	}
	return "", fmt.Errorf("no secret found")
}

func mkCephFSDirIfNotExist(dirPath string) (err error) {
	_, err = os.Stat(dirPath)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func mountCephFSRootDir(mntToLocalPath string, opt *CephfsParams) (err error) {
	//mount -t ceph 172.16.0.50:6789:/aa /docker/xxx -o name=admin,secret=AQBbjulZ1GGaGRAAm/iZM+/fBUbZDI3SpQKPmg==

	_, err = os.Stat(mntToLocalPath)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(mntToLocalPath, 0755)
		if err != nil {
			return err
		}
	}

	var monitors string
	for _, v := range opt.Monitors {
		monitors += v + ":/"
	}

	args := []string{"-t", "ceph", monitors, mntToLocalPath, "-o", "name=" + opt.AdminID + ",secret=" + opt.AdminSecret}
	cmd := exec.Command("mount", args...)
	//	fmt.Println(strings.Join(cmd.Args, " "))
	err = cmd.Run()

	if err != nil {
		buf, _ := cmd.Output()
		return fmt.Errorf("%s %s", err.Error(), string(buf))
	}

	return nil
}

func umountCephFSRootDir(mntToLocalPath string) (err error) {
	args := []string{mntToLocalPath}
	cmd := exec.Command("umount", args...)
	err = cmd.Run()

	if err != nil {
		buf, _ := cmd.Output()
		return fmt.Errorf("%s %s", err.Error(), string(buf))
	}

	if err = os.RemoveAll(mntToLocalPath); err != nil {
		return err
	}

	return nil
}

// deletePVCDir : delete pvc directory from
func deletePVCDir(pvcCephPath string, param *CephfsParams) (err error) {
	/*pathArr := strings.Split(pvcDirPath, "/")
	var parentDirPath string
	if len(pathArr) > 1 {
		parentDirPath = strings.Join(pathArr[:len(pathArr)-1], "/")
	} else {
		return fmt.Errorf("could not find parent directory path")
	}*/

	mntPath := cephFSRootMntPath + pvcCephPath
	err = mountCephFSRootDir(mntPath, param)
	if err != nil {
		err = fmt.Errorf("mount cepfs root to local failed -> %s", err.Error())
		return err
	}
	fmt.Println("mntPath is " + mntPath)
	pvDirPath := mntPath + pvcCephPath
	fmt.Println("pvDirPath is " + pvDirPath)
	err = os.RemoveAll(pvDirPath)
	if err != nil {
		err = fmt.Errorf("delete cephfs directory failed -> %s", err.Error())
		return err
	}
	err = umountCephFSRootDir(mntPath)
	if err != nil {
		err = fmt.Errorf("umount <%s> failed -> %s", mntPath, err.Error())
		return err
	}
	return nil
}
