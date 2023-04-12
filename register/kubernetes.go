package register

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	listerv1 "k8s.io/client-go/listers/core/v1"
	"net/url"
	"strconv"
	"time"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"os"
)

const (
	LabelsKeyServiceID      = "z-service-id"
	LabelsKeyServiceName    = "z-service-name"
	LabelsKeyServiceVersion = "z-service-version"

	// AnnotationsKeyProtocolMap is used to define the protocols of the service
	// Through the value of this field, Kratos can obtain the application layer protocol corresponding to the port
	// Example value: {"80": "http", "8081": "grpc"}
	AnnotationsKeyProtocolMap = "z-service-protocols"
)

type K8sRegister struct {
	clientSet       *kubernetes.Clientset
	podLister       listerv1.PodLister
	informerFactory informers.SharedInformerFactory
	stopChan        chan struct{}
}

func NewK8sRegister(clientSet *kubernetes.Clientset) *K8sRegister {
	informerFactory := informers.NewSharedInformerFactory(clientSet, time.Minute*10)
	//podInformer := informerFactory.Core().V1().Pods().Informer()

	// init informer
	podLister := informerFactory.Core().V1().Pods().Lister()
	return &K8sRegister{
		clientSet:       clientSet,
		podLister:       podLister,
		informerFactory: informerFactory,
		stopChan:        make(chan struct{}),
	}
}

func (k *K8sRegister) Start() error {
	// start listen
	k.informerFactory.Start(k.stopChan)
	return nil
}

type protocolMap map[string]string

func (k *K8sRegister) GetService(_ context.Context, serviceName string) ([]*Service, error) {
	pods, err := k.podLister.List(labels.SelectorFromSet(map[string]string{
		LabelsKeyServiceName: serviceName,
	}))
	if err != nil {
		return nil, err
	}
	res := make([]*Service, 0, len(pods))
	for _, pod := range pods {
		if pod.Status.Phase != corev1.PodRunning {
			continue
		}
		service, err := getServiceInstanceFromPod(pod)
		if err != nil {
			return nil, err
		}
		res = append(res, service)
	}
	return res, nil

}

func getServiceInstanceFromPod(pod *corev1.Pod) (*Service, error) {
	podLabels := pod.Labels
	podIp := pod.Status.PodIP
	annotations := pod.GetAnnotations()
	protocolMap := make(map[string]string)
	if s := annotations[AnnotationsKeyProtocolMap]; !isEmptyObjectString(s) {
		err := jsoniter.UnmarshalFromString(s, &protocolMap)
		if err != nil {
			return nil, fmt.Errorf("handler service protocol failed:%v", err)
		}
	}
	containers := pod.Spec.Containers
	var endpoints []string
	for _, container := range containers {
		ports := container.Ports
		for _, port := range ports {
			containerPort := port.ContainerPort
			if protocol, ok := protocolMap[strconv.Itoa(int(containerPort))]; ok {
				endpoint := fmt.Sprintf("%s://%s:%d", protocol, podIp, containerPort)
				endpoints = append(endpoints, endpoint)
			}
		}
	}
	return &Service{
		ID:       podLabels[LabelsKeyServiceID],
		Name:     podLabels[LabelsKeyServiceName],
		IP:       "",
		Port:     0,
		Version:  podLabels[LabelsKeyServiceVersion],
		EndPoint: endpoints,
	}, nil
}

func isEmptyObjectString(s string) bool {
	switch s {
	case "", "[]", "null", "nil", "{}":
		return true
	}
	return false
}

func (k *K8sRegister) Register(ctx context.Context, service *Service) error {
	protocolMap := protocolMap{}
	points := service.EndPoint
	for _, point := range points {
		parse, err := url.Parse(point)
		if err != nil {
			return fmt.Errorf("parse url failed")
		}
		protocolMap[parse.Port()] = parse.Scheme
	}
	toString, err := jsoniter.MarshalToString(protocolMap)
	if err != nil {
		return fmt.Errorf("json parse failed: %v", err)
	}
	patchBytes, err := jsoniter.Marshal(map[string]interface{}{
		"metadata": metav1.ObjectMeta{
			Labels: map[string]string{
				LabelsKeyServiceID:      service.ID,
				LabelsKeyServiceName:    service.Name,
				LabelsKeyServiceVersion: service.Version,
			},
			Annotations: map[string]string{
				AnnotationsKeyProtocolMap: toString,
			},
		},
	})
	_, err = k.clientSet.
		CoreV1().
		Pods(getNamespaces()).
		Patch(ctx, getPodName(), types.StrategicMergePatchType, patchBytes, metav1.PatchOptions{})
	if err != nil {
		return err
	}
	return nil
}

func getPodName() string {
	return os.Getenv("HOSTNAME")
}

var currentNamespace = LoadNamespace()

// ServiceAccountNamespacePath pod中当前namespace的文件
const ServiceAccountNamespacePath = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

func LoadNamespace() string {
	data, err := os.ReadFile(ServiceAccountNamespacePath)
	if err != nil {
		return "default"
	}
	return string(data)
}

func getNamespaces() string {
	return currentNamespace
}

func (k *K8sRegister) UnRegister(ctx context.Context, _ *Service) error {
	return k.Register(ctx, &Service{})
}
