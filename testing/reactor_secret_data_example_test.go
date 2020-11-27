package testing

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func ExampleSecretDataReactor() {
	clientset := fake.NewSimpleClientset()
	clientset.PrependReactor("create", "secrets", SecretDataReactor)

	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "my-secret",
		},
		StringData: map[string]string{
			"my-key": "my-value",
		},
		Type: v1.SecretTypeOpaque,
	}

	createdSecret, _ := clientset.CoreV1().Secrets(secret.Namespace).Create(ctx, secret, metav1.CreateOptions{})

	fmt.Printf("secret.Data[\"my-key\"]: %s", string(createdSecret.Data["my-key"])) //remember: the secrets StringData field was set, not the Data field

	// Output:
	// secret.Data["my-key"]: my-value
}
