package register

import (
	"context"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	namespace  = "default"
	deployName = "hello-deployment"
	podName    = "hello"
)

func int32Ptr(i int32) *int32 { return &i }

var deployment = appsv1.Deployment{
	ObjectMeta: metav1.ObjectMeta{
		Name: deployName,
	},
	Spec: appsv1.DeploymentSpec{
		Replicas: int32Ptr(1),
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app": podName,
			},
		},
		Template: apiv1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"app": podName,
				},
			},
			Spec: apiv1.PodSpec{
				Containers: []apiv1.Container{
					{
						Name:  "nginx",
						Image: "nginx:alpine",
						Ports: []apiv1.ContainerPort{
							{
								Name:          "http",
								Protocol:      apiv1.ProtocolTCP,
								ContainerPort: 80,
							},
						},
						Command: []string{
							"nginx",
							"-g",
							"daemon off;",
						},
					},
				},
			},
		},
	},
}

func TestK8sRegister_Register(t *testing.T) {
	check := assert.New(t)
	clientSet, err := getClientSet()
	check.Nil(err)

	register := NewK8sRegister(clientSet)
	_, err = clientSet.AppsV1().Deployments(namespace).Create(context.Background(), &deployment, metav1.CreateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	// 删除资源
	defer func() {
		err = clientSet.AppsV1().Deployments(namespace).Delete(context.Background(), deployName, metav1.DeleteOptions{})
		check.Nil(err)
	}()
	err = register.Start()
	if err != nil {
	}

	time.Sleep(time.Second)
	pod, err := clientSet.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: "app=hello"})
	if err != nil {
		t.Fatal(err)
	}
	if len(pod.Items) < 0 {
		t.Fatal("fetch resouce failed")
	}
	err = os.Setenv("HOSTNAME", pod.Items[0].Name)
	check.Nil(err)
	hello := Service{
		ID:       "0",
		Name:     "hello",
		Version:  "v1",
		EndPoint: []string{"http://127.0.0.1:80"},
	}

	err = register.Register(context.Background(), &hello)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Duration(1) * time.Second)
	// 根据service name获取资源
	services, err := register.GetService(context.Background(), "hello")
	check.Nil(err)
	check.Equal("hello", services[0].Name)
	err = register.UnRegister(context.Background(), &hello)
	check.Nil(err)
	time.Sleep(time.Second)
	services, err = register.GetService(context.Background(), "hello")
	check.Nil(err)
	check.Equal(0, len(services))
}

func getClientSet() (*kubernetes.Clientset, error) {
	// 如果是在集群中，则使用集群的内部配置
	config, err := rest.InClusterConfig()
	if err != nil {
		homeDir := homedir.HomeDir()
		configDir := filepath.Join(homeDir, ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", configDir)
		if err != nil {
			return nil, err
		}
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientSet, nil
}

func TestInformer(t *testing.T) {
	check := assert.New(t)
	client, err := getClientSet()
	check.Nil(err)
	watch, err := client.AppsV1().Deployments(getNamespaces()).Watch(context.Background(), metav1.ListOptions{})
	check.Nil(err)
	for {
		select {
		case res := <-watch.ResultChan():
			t.Log(res.Type)
			// 显示变化的类型 ADDED MODIFIED
			t.Log(res.Type)
			// 变更的对象
			t.Log(res.Object)
		}
	}
}
