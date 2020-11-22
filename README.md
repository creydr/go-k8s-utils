# go-k8s-utils

This repository contains utils for the work with [Kubernetes](https://kubernetes.io/), in specific with the [go-client](https://github.com/kubernetes/client-go) library.

## Testing

This package contains utils which are useful for testing (e.g. with the fake-client).

### Reactors

Having the [ReactionFunc](https://godoc.org/k8s.io/client-go/testing#ReactionFunc) signature, the following helper exist:

* `GenerateNameReactor`: setting the `ObjectMeta.Name` field, based on `ObjectMeta.GenerateName` (as `ObjectMeta.Name` is not set automatically by the fake-client, if only `ObjectMeta.GenerateName` is set).
* `SecretDataReactor`: setting the `Secret.Data` field based on `Secret.StringData` (as `Secret.Data` is not set automatically by the fake-client, if only `Secret.StringData` is set).

Examples:

```golang
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
```

```golang
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

fmt.Printf("Secrets data: %+v", createdSecret.Data) //remember: the secrets StringData field was set, not the Data field
```