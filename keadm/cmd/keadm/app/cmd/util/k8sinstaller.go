/*
Copyright 2019 The Kubeedge Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"

	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	crdclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

//K8SInstTool embedes Common struct and contains the default K8S version and
//a flag depicting if host is an edge or cloud node
//It implements ToolsInstaller interface
type K8SInstTool struct {
	Common
}

//InstallTools sets the OS interface, checks if K8S installation is required or not.
//If required then install the said version.
func (ks *K8SInstTool) InstallTools() error {
	ks.SetOSInterface(GetOSInterface())

	cloudCoreRunning, err := ks.IsKubeEdgeProcessRunning(KubeCloudBinaryName)
	if err != nil {
		return err
	}
	if cloudCoreRunning {
		return fmt.Errorf("CloudCore is already running on this node, please run reset to clean up first")
	}

	err = ks.IsK8SComponentInstalled(ks.KubeConfig, ks.Master)
	if err != nil {
		return err
	}

	fmt.Println("Kubernetes version verification passed, KubeEdge installation will start...")

	err = installCRDs(ks.KubeConfig, ks.Master)
	if err != nil {
		return err
	}

	return nil
}

func installCRDs(kubeConfig, master string) error {
	config, err := BuildConfig(kubeConfig, master)
	if err != nil {
		return fmt.Errorf("Failed to build config, err: %v", err)
	}

	crdClient, err := crdclient.NewForConfig(config)
	if err != nil {
		return err
	}

	// Todo: need to add the crds ro release package
	// create the dir for kubeedge crd
	err = os.MkdirAll(KubeEdgeCrdPath+"/devices", os.ModePerm)
	if err != nil {
		return fmt.Errorf("not able to create %s folder path", KubeEdgeLogPath)
	}
	for _, crdFile := range []string{"devices/devices_v1alpha1_device.yaml",
		"devices/devices_v1alpha1_devicemodel.yaml"} {
		//Download the tar from repo
		dwnldURL := fmt.Sprintf("cd %s && wget -k --no-check-certificate --progress=bar:force %s/%s", KubeEdgeCrdPath+"/devices", KubeEdgeCRDDownloadURL, crdFile)
		_, err := runCommandWithShell(dwnldURL)
		if err != nil {
			return err
		}

		// not found err, create crd from crd file
		err = createKubeEdgeCRD(crdClient, KubeEdgeCrdPath+"/"+crdFile)
		if err != nil && !apierrors.IsAlreadyExists(err) {
			return err
		}
	}

	// Todo: need to add the crds ro release package
	// create the dir for kubeedge crd
	err = os.MkdirAll(KubeEdgeCrdPath+"/reliablesyncs", os.ModePerm)
	if err != nil {
		return fmt.Errorf("not able to create %s folder path", KubeEdgeLogPath)
	}
	for _, crdFile := range []string{"reliablesyncs/cluster_objectsync_v1alpha1.yaml",
		"reliablesyncs/objectsync_v1alpha1.yaml"} {
		//Download the tar from repo
		dwnldURL := fmt.Sprintf("cd %s && wget -k --no-check-certificate --progress=bar:force %s/%s", KubeEdgeCrdPath+"/reliablesyncs", KubeEdgeCRDDownloadURL, crdFile)
		_, err := runCommandWithShell(dwnldURL)
		if err != nil {
			return err
		}

		// not found err, create crd from crd file
		err = createKubeEdgeCRD(crdClient, KubeEdgeCrdPath+"/"+crdFile)
		if err != nil && !apierrors.IsAlreadyExists(err) {
			return err
		}
	}

	return nil
}

func createKubeEdgeCRD(clientset crdclient.Interface, crdFile string) error {
	content, err := ioutil.ReadFile(crdFile)
	if err != nil {
		return fmt.Errorf("read crd yaml error: %v", err)
	}

	kubeEdgeCRD := &apiextensionsv1beta1.CustomResourceDefinition{}
	err = yaml.Unmarshal(content, kubeEdgeCRD)
	if err != nil {
		return fmt.Errorf("unmarshal tfjobCRD error: %v", err)
	}

	_, err = clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Create(kubeEdgeCRD)

	return err
}

//TearDown shoud uninstall K8S, but it is not required either for cloud or edge node.
//It is defined so that K8SInstTool implements ToolsInstaller interface
func (ks *K8SInstTool) TearDown() error {
	return nil
}
