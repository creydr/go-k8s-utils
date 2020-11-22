package testing

import (
	"strings"
	"testing"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGenerateNameReactor(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	clientset.PrependReactor("create", "pods", GenerateNameReactor)

	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    "default",
			GenerateName: "testpod-",
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "nginx",
					Image: "nginx",
				},
			},
		},
	}

	createdPod, err := clientset.CoreV1().Pods(pod.Namespace).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		t.Errorf("GenerateNameReactor() Test setup failed on creating pod. Error: %v", err)
		return
	}

	// should update pods name
	if !strings.HasPrefix(createdPod.Name, pod.GenerateName) || len(createdPod.Name) < len(pod.GenerateName) {
		t.Errorf("GenerateNameReactor() Did not update name of pod. Expected name starting with %s, actual name: %s", pod.GenerateName, createdPod.Name)
		return
	}

	// should work with multiple reactors
	clientset.PrependReactor("create", "secrets", SecretDataReactor)

	createdPod2, err := clientset.CoreV1().Pods(pod.Namespace).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			t.Errorf("GenerateNameReactor() Did not create a unique ID for the name. Error: %s", err)
		} else {
			t.Errorf("GenerateNameReactor() Test setup failed on creating pod. Error: %v", err)
		}
		return
	}

	if !strings.HasPrefix(createdPod2.Name, pod.GenerateName) || len(createdPod2.Name) < len(pod.GenerateName) {
		t.Errorf("GenerateNameReactor() Did not update name of pod. Expected name starting with %s, actual name: %s", pod.GenerateName, createdPod2.Name)
		return
	}
}
