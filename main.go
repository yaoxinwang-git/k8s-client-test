package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	v1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
	"path/filepath"
)

func main() {
	var kubeConfig *string
	ctx := context.Background()
	//取得kubeconfig的路径
	if home := homedir.HomeDir(); home != "" {
		kubeConfig = flag.String("kubeConfig", filepath.Join("./", ".kube", "config"), "absolute path to the kubeconfig file")
	} else {
		kubeConfig = flag.String("kubeConfig", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		klog.Fatal(err)
		return
	}
	clientSet, err := kubernetes.NewForConfig(config)
	//listNameSpace(clientSet, ctx)
	//createDeploymentFromYaml(clientSet, ctx)
	deleteDeployment(clientSet, ctx)

}

func listNameSpace(clientSet *kubernetes.Clientset, ctx context.Context) {
	//获取namespaces信息
	namespaces := clientSet.CoreV1().Namespaces()
	namespaceList, err := namespaces.List(ctx, metav1.ListOptions{})
	if err != nil {
		klog.Fatal(err)
		return
	}
	for _, namespace := range namespaceList.Items {
		fmt.Println(namespace.Name)
	}
}

func createDeploymentFromYaml(clientSet *kubernetes.Clientset, ctx context.Context) {

	//创建namespace myweb
	myweb := &coreV1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: "myweb"},
	}
	clientSet.CoreV1().Namespaces().Create(ctx, myweb, metav1.CreateOptions{})

	//创建deployment myngx
	deploye := &v1.Deployment{}
	bytes, err := ioutil.ReadFile("yamls/nginx.yaml")
	if err != nil {
		klog.Fatal(err)
	}

	toJSON, err := yaml.ToJSON(bytes)

	if err != nil {
		klog.Fatal(err)
	}

	json.Unmarshal(toJSON, &deploye)
	_, err = clientSet.AppsV1().Deployments("myweb").Create(ctx, deploye, metav1.CreateOptions{})
	if err != nil {
		klog.Fatal(err)
	}
}

func deleteDeployment(clientSet *kubernetes.Clientset, ctx context.Context) {

	err := clientSet.AppsV1().Deployments("myweb").Delete(ctx, "myngx", metav1.DeleteOptions{})
	if err != nil {
		klog.Fatal(err)
	}
}
