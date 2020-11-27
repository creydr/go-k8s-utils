package testing

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ktesting "k8s.io/client-go/testing"
)

// SecretDataReactor sets the secret.Data field based on the values from secret.StringData
func SecretDataReactor(action ktesting.Action) (bool, runtime.Object, error) {
	secret, ok := action.(ktesting.CreateAction).GetObject().(*v1.Secret)
	if !ok {
		return false, nil, fmt.Errorf("SecretDataReactor can only be applied on secrets")
	}

	if len(secret.StringData) > 0 {
		if secret.Data == nil {
			secret.Data = make(map[string][]byte)
		}

		for k, v := range secret.StringData {
			secret.Data[k] = []byte(v)
		}
	}

	return false, nil, nil
}
