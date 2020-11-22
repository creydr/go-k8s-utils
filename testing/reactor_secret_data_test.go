package testing

import (
	"encoding/base64"
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

	val1decoded := make([]byte, base64.StdEncoding.DecodedLen(len(val1)))
	lenVal1, _ := base64.StdEncoding.Decode(val1decoded, val1)
	val2decoded := make([]byte, base64.StdEncoding.DecodedLen(len(val2)))
	lenVal2, _ := base64.StdEncoding.Decode(val2decoded, val2)

	if string(val1decoded[:lenVal1]) != secret.StringData["key1"] || string(val2decoded[:lenVal2]) != secret.StringData["key2"] {
		t.Errorf("SecretDataReactor() Did not encode correctly. Got %s and %s, but expected %s and %s", string(val1decoded[:lenVal1]), secret.StringData["key1"], string(val2decoded[:lenVal2]), secret.StringData["key2"])
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
