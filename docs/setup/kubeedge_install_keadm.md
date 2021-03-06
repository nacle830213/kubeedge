# Setup from KubeEdge Installer

Please refer to KubeEdge Installer proposal document for details on the motivation of having KubeEdge Installer.
It also explains the functionality of the proposed commands.
[KubeEdge Installer Doc](https://github.com/kubeedge/kubeedge/blob/master/docs/proposals/keadm-scope.md/)

## Limitation

- Currently support of `KubeEdge installer` is available only for Ubuntu OS. CentOS support is in-progress.

## Getting KubeEdge Installer

There are currently two ways to get keadm

- Download from [KubeEdge Release](<https://github.com/kubeedge/kubeedge/releases>)

  1. Go to [KubeEdge Release](<https://github.com/kubeedge/kubeedge/releases>) page and download `keadm-$VERSION-$OS-$ARCH.tar.gz.`.
  2. Untar it at desired location, by executing `tar -xvzf keadm-$VERSION-$OS-$ARCH.tar.gz`.
  3. kubeedge folder is created after execution the command.

- Building from source

  1. Download the source code.
  
      ```shell
      git clone https://github.com/kubeedge/kubeedge.git $GOPATH/src/github.com/ kubeedge/kubeedge
      cd $GOPATH/src/github.com/kubeedge/kubeedge/keadm
      make
      ```

      or

      ```shell
      go get github.com/kubeedge/kubeedge/keadm/cmd/keadm
      ```

  2. Binary `keadm` is available in current path. If you are using `go` get the binary is available in `$GOPATH/bin/`

## Setup Cloud Side (KubeEdge Master Node)

Referring to `KubeEdge Installer Doc`, the command to install KubeEdge cloud component (edge controller) and pre-requisites.
Port 8080, 6443 and 10000 in your cloud component needs to be accessible for your edge nodes.

keadm by default can install Docker, Kubernetes and KubeEdge. It also provide flag by which specific versions can be set.

1. Execute `keadm init` : keadm needs super user rights (or root rights) to run successfully.

    Command flags
    The optional flags with this command are mentioned below

    ```shell
    keadm init --help

    keadm init command bootstraps KubeEdge cloud component. It checks if the pre-requisites are installed already, If not installed, this command will help in download, install and execute on the host.

    Usage:
      keadm init [flags]

    Examples:

    keadm init

    Flags:
        -h, --help                                       help for init
        --docker-version string[="18.06.0"]          Use this key to download and use the required Docker version (default "18.06.0")
        --kubeedge-version string[="0.3.0-beta.0"]   Use this key to download and use the required KubeEdge version (default "0.3.0-beta.0")
        --kubernetes-version string[="1.14.1"]       Use this key to download and use the required Kubernetes version (default "1.14.1"). It will install `kubeadm`, `kubectl` and `kubelet` in this host.
        --pod-network-cidr string[="100.64.0.0/10"]   Use this key to set the Kubernetes pod Network cidr  (default "100.64.0.0/10")
        --kube-config string                          Use this key to set kube-config path, eg: $HOME/.kube/config
    ```

    Command format is

    ```shell
    keadm init --docker-version=<expected version> --kubernetes-version=<expected version>  --kubeedge-version=<expected version>
    ```

  **NOTE:**
  Version mentioned as defaults for Docker and K8S are being tested with.

  **NOTE:**
On the console output, observe the below line

```shell
kubeadm join **192.168.20.134**:6443 --token 2lze16.l06eeqzgdz8sfcvh \
         --discovery-token-ca-cert-hash sha256:1e5c808e1022937474ba264bb54fea42b05eddb9fde2d35c9cad5b83cf5ef9ac  
```

After Kubeedge init, please note the **cloudIP** as highlighted above generated from console output and port is **8080**.

## Setup Edge Side (KubeEdge Worker Node)

Referring to `KubeEdge Installer Doc`, the command to install KubeEdge Edge component (edge core) and pre-requisites

Execute `keadm join <flags>`

 Command flags

  The optional flags with this command are shown in below shell

  ```shell
  $ keadm join --help

  "keadm join" command bootstraps KubeEdge's edge component. It checks if the pre-requisites are installed already, If not installed, this command will help in download, to install the prerequisites. It will help the edge node to connect to the cloud.
  
  Usage:
  keadm join [flags]

  Flags:
      --docker-version string[="18.06.0"]          Use this key to download and use the required Docker version (default "18.06.0")
  -e, --cloudcoreip string                         IP address of KubeEdge cloudcore
  -i, --edgenodeid string                          KubeEdge Node unique identification string, If flag not used then the command will generate a unique id on its own
  -h, --help                                       help for join
  -k, --k8sserverip string                         IP:Port address of K8S API-Server
      --kubeedge-version string[="0.3.0-beta.0"]   Use this key to download and use the required KubeEdge version (default "0.3.0-beta.0")

For this command --cloudcoreip flag is a Mandatory flag
  ```

 Examples:

  ```shell
    keadm join --cloudcoreip=<ip address> --edgenodeid=<unique string as edge identifier>
  ```

 ```shell
  keadm join --cloudcoreip=10.20.30.40 --edgenodeid=testing123 --kubeedge-version=1.1.1 --k8sserverip=50.60.70.80:8080
 ```

 In case, any option is used in a format like as shown for "--docker-version" or "--docker-version=", without a value then default values will be used. Also options like "--docker-version", and "--kubeedge-version", version should be in format like "18.06.3" and "0.2.1".

**IMPORTANT NOTE:** The KubeEdge version used in cloud and edge side should be same.

Sample execution output:

```shell
# ./keadm join --cloudcoreip=192.168.20.50 --edgenodeid=testing123 --k8sserverip=192.168.20.50:8080
Same version of docker already installed in this host
Host has mosquit+ already installed and running. Hence skipping the installation steps !!!
Expected or Default KubeEdge version 0.3.0-beta.0 is already downloaded
kubeedge/
kubeedge/edge/
kubeedge/edge/conf/
kubeedge/edge/conf/modules.yaml
kubeedge/edge/conf/logging.yaml
kubeedge/edge/conf/edge.yaml
kubeedge/edge/edgecore
kubeedge/cloud/
kubeedge/cloud/cloudcore
kubeedge/cloud/conf/
kubeedge/cloud/conf/controller.yaml
kubeedge/cloud/conf/modules.yaml
kubeedge/cloud/conf/logging.yaml
kubeedge/version

KubeEdge Edge Node: testing123 successfully add to kube-apiserver, with operation status: 201 Created
Content {"kind":"Node","apiVersion":"v1","metadata":{"name":"testing123","selfLink":"/api/v1/nodes/testing123","uid":"87d8d7a3-7acd-11e9-b86b-286ed488c645","resourceVersion":"3864","creationTimestamp":"2019-05-20T07:04:37Z","labels":{"name":"edge-node"}},"spec":{"taints":[{"key":"node.kubernetes.io/not-ready","effect":"NoSchedule"}]},"status":{"daemonEndpoints":{"kubeletEndpoint":{"Port":0}},"nodeInfo":{"machineID":"","systemUUID":"","bootID":"","kernelVersion":"","osImage":"","containerRuntimeVersion":"","kubeletVersion":"","kubeProxyVersion":"","operatingSystem":"","architecture":""}}}

KubeEdge edge core is running, For logs visit /etc/kubeedge/kubeedge/edge/
```

**Note**:Cloud IP refers to IP generated ,from the step 1 as highlighted

## Reset KubeEdge Master and Worker nodes

Referring to `KubeEdge Installer Doc`, the command to stop KubeEdge cloud (edge controller). It doesn't uninstall/remove any of the pre-requisites.

Execute `keadm reset`

Command flags

```shell
keadm reset --help

keadm reset command can be executed in both cloud and edge node
In master node it shuts down the cloud processes of KubeEdge
In worker node it shuts down the edge processes of KubeEdge

Usage:
  keadm reset [flags]

Examples:

For master node:
keadm reset

For worker node:
keadm reset --k8sserverip 10.20.30.40:8080


Flags:
  -h, --help                 help for reset
  -k, --k8sserverip string   IP:Port address of cloud components host/VM

```

## Errata

1. If GPG key for docker repo fail to fetch from key server. Please refer [Docker GPG error fix](<https://forums.docker.com/t/gpg-key-for-docker-repo-fail-to-fetch-from-key-server/24253>)

2. After kubeadm init, if you face any errors regarding swap memory and preflight checks please refer  [Kubernetes preflight error fix](<https://github.com/kubernetes/kubeadm/issues/610>)

3. Error in CloudCore

    If you are getting the below error in Cloudcore.log

    ```shell
    E1231 04:37:27.397431   19607 reflector.go:125] github.com/kubeedge/kubeedge/cloud/pkg/devicecontroller/manager/device.go:40: Failed to list *v1alpha1.Device: the server could not find the requested resource (get devices.devices.kubeedge.io)
    E1231 04:37:27.398273   19607 reflector.go:125] github.com/kubeedge/kubeedge/cloud/pkg/devicecontroller/manager/devicemodel.go:40: Failed to list *v1alpha1.DeviceModel: the server could not find the requested resource (get devicemodels.devices.kubeedge.io)
    ```

    browse to the

    ```shell
    cd $GOPATH/src/github.com/kubeedge/kubeedge/build/crds/devices
    ```

    and apply the below

    ```shell
      kubectl create -f devices_v1alpha1_devicemodel.yaml
      kubectl create -f devices_v1alpha1_device.yaml
    ```

    or

    ```shell
     kubectl create -f https://raw.githubusercontent.com/kubeedge/kubeedge/<kubeEdge Version>/build/crds/devices/devices_v1alpha1_device.yaml
     kubectl create -f https://raw.githubusercontent.com/kubeedge/kubeedge/<kubeEdge Version>/build/crds/devices/devices_v1alpha1_devicemodel.yaml
    ```

  Also, create ClusterObjectSync and ObjectSync CRDs which are used in reliable message delivery.

    ```shell
     cd $GOPATH/src/github.com/kubeedge/kubeedge/build/crds/reliablesyncs
     kubectl create -f cluster_objectsync_v1alpha1.yaml
     kubectl create -f objectsync_v1alpha1.yaml
    ```
