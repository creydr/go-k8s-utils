package testing

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

var ctx = context.Background()

func ExampleGenerateNameReactor() {
	mockStringRandFunc() // only for testing
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

	createdPod, _ := clientset.CoreV1().Pods(pod.Namespace).Create(ctx, pod, metav1.CreateOptions{})

	fmt.Printf("Name: %s", createdPod.Name) //remember: the pods Name field was not set, only GenerateName

	// Output:
	// Name: testpod-hm7shz2o

	unmockStringRandFunc()
}

var stringRandOrigFunc = stringRandFunc
var stringRandMockFunc = func(n int) string {
	return "hm7shz2o"
}

func mockStringRandFunc() {
	stringRandFunc = stringRandMockFunc
}

func unmockStringRandFunc() {
	stringRandMockFunc = stringRandOrigFunc
}
