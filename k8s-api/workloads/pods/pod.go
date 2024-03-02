package pods

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type PodClient struct {
	client v1.PodInterface
	ns     string
}

func NewPodClient(clientset *kubernetes.Clientset, namespace string) *PodClient {
	return &PodClient{
		client: clientset.CoreV1().Pods(namespace),
		ns:     namespace,
	}
}

func (p *PodClient) Create(name string) {}

func (p *PodClient) Update(name string) {}

func (p *PodClient) Delete(name string) error {
	return p.client.Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (p *PodClient) Get(name string) (*corev1.Pod, error) {
	return p.client.Get(context.TODO(), name, metav1.GetOptions{})
}

func (p *PodClient) List(labels ...string) ([]corev1.Pod, error) {
	var podList *corev1.PodList
	var err error

	if len(labels) == 0 {
		podList, err = p.client.List(context.TODO(), metav1.ListOptions{})
	} else {
		podList, err = p.client.List(context.TODO(), metav1.ListOptions{
			LabelSelector: labels[0],
		})
	}

	if err != nil {
		return nil, err
	}
	return podList.Items, nil
}
