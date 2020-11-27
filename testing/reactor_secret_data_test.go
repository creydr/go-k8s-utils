package testing

import (
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestSecretDataReactor(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	clientset.PrependReactor("create", "secrets", SecretDataReactor)

	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "my-secret",
		},
		StringData: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
		Type: v1.SecretTypeOpaque,
	}

	createdSecret, err := clientset.CoreV1().Secrets(secret.Namespace).Create(ctx, secret, metav1.CreateOptions{})
	if err != nil {
		t.Errorf("SecretDataReactor() Test setup failed on creating secret. Error: %v", err)
		return
	}

	// should populate .Data field
	val1, ok1 := createdSecret.Data["key1"]
	val2, ok2 := createdSecret.Data["key2"]
	if !ok1 || !ok2 {
		t.Errorf("SecretDataReactor() Did not populate secret.Data field")
		return
	}

	if string(val1) != secret.StringData["key1"] || string(val2) != secret.StringData["key2"] {
		t.Errorf("SecretDataReactor() Did not set secret.Data field correctly. Got %s and %s, but expected %s and %s", string(val1), secret.StringData["key1"], string(val2), secret.StringData["key2"])
		return
	}

	// should not affect other resource types except secrets
	clientset.PrependReactor("create", "pods", SecretDataReactor)

	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "my-pod",
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

	_, err = clientset.CoreV1().Pods(pod.Namespace).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		t.Errorf("SecretDataReactor() Should not affect other resource types. Err: %s", err)
		return
	}
}
