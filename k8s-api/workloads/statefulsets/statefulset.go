package statefulsets

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/util/retry"
)

type StatefulSetClient struct {
	client v1.StatefulSetInterface
	ns     string
}

func NewStatefulSetClient(clientset *kubernetes.Clientset, namespace string) *StatefulSetClient {
	return &StatefulSetClient{
		client: clientset.AppsV1().StatefulSets(namespace),
		ns:     namespace,
	}
}

func (s *StatefulSetClient) Create(name string) (*appsv1.StatefulSet, error) {
	ss := &appsv1.StatefulSet{}

	return s.client.Create(context.TODO(), ss, metav1.CreateOptions{})
}

func (s *StatefulSetClient) Update(name string) error {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		ss, getErr := s.client.Get(context.TODO(), name, metav1.GetOptions{})
		if getErr != nil {
			return getErr
		}

		// ss.Spec.Replicas

		_, updateErr := s.client.Update(context.TODO(), ss, metav1.UpdateOptions{})
		return updateErr
	})

	return retryErr
}

func (s *StatefulSetClient) Delete(name string) error {
	return s.client.Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (s *StatefulSetClient) Get(name string) (*appsv1.StatefulSet, error) {
	return s.client.Get(context.TODO(), name, metav1.GetOptions{})
}

func (s *StatefulSetClient) List(labels ...string) ([]appsv1.StatefulSet, error) {
	var sts *appsv1.StatefulSetList
	var err error

	if len(labels) == 0 {
		sts, err = s.client.List(context.TODO(), metav1.ListOptions{})
	} else {
		sts, err = s.client.List(context.TODO(), metav1.ListOptions{
			LabelSelector: labels[0],
		})
	}

	if err != nil {
		return nil, err
	}
	return sts.Items, err
}
